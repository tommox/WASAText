<template>
  <div v-if="chat" class="w-full flex flex-col flex-grow bg-gray-100">
    <!-- Header Chat -->
     <div class="chat-header">
      <label v-if="type === 'group'" class="profile-picture">
        <img :src="avatarUrl || defaultAvatar" alt="Foto Gruppo" class="profile-img" />
        <input type="file" @change="uploadGroupPhoto" accept="image/*" class="upload-input" />
      </label>
      <img v-else :src="avatarUrl" class="w-10 h-10 rounded-full">
      <span class="ml-3 font-bold text-lg">{{ chatName }}</span>
      <button @click="deleteConversation" class="delete-btn">🗑️</button>
    </div>

    <!-- Lista Messaggi -->
    <div ref="messageContainer" class="chat-messages flex-grow overflow-y-auto p-3">
      <div v-if="loading" class="text-center text-gray-500">Caricamento messaggi...</div>
      <div v-else-if="messages.length === 0" class="text-center text-gray-500">Inizia una nuova conversazione!</div>

      <div 
        v-else v-for="message in messages" 
        :key="message.id" 
        class="mb-2 flex"
        :class="{'justify-end': message.sender === 'me', 'justify-start': message.sender === 'other'}">
        <div 
          class="relative flex flex-col max-w-xs p-3 rounded-3xl shadow-md"
          :class="message.sender === 'me' ? 'bg-blue-500 text-black self-end' : 'bg-gray-200 self-start'">
          <div class="flex items-end">
            <span>{{ message.text }}</span>
            <div class="message-time" @click="openMenu(message.id, message.sender)">{{formatTime(message.timestamp)}}</div>
            <div v-if="showOptions && selectedMessageId === message.id" class="modal-overlay">
              <div class="modal-content">
                <h2>Seleziona un opzione</h2>
                <div class="option-list">
                  <div v-if="selectedMessageSender === 'me'" class="option-item" @click="deleteMessage(selectedMessageId)">Elimina</div>
                </div>
                <button @click="showOptions = false" class="cancel-btn">Chiudi</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Input Messaggio -->
    <div class="p-3 border-t bg-gray-100 flex">
      <input 
        v-model="newMessage" 
        type="text" 
        placeholder="Scrivi un messaggio..." 
        class="flex-grow px-3 py-2 border rounded-lg focus:outline-none"
        @keyup.enter="sendMessage"
      />
      <button @click="sendMessage" class="ml-3 p-2 bg-blue-500 text-white rounded-lg">
        ➤
      </button>
    </div>
  </div>

  <div v-else class="empty-chat w-full flex items-center justify-center flex-grow bg-gray-100">
    Apri o inizia una nuova conversazione
  </div>
</template>

<script>
import axios from "axios";
import defaultAvatar from "@/assets/images/user.png";
import eventBus from "@/eventBus";

export default {
  props: { chat: Object ,type: String },
  data() {
    return { 
      newMessage: "",
      messages: [],
      loading: true,
      avatarUrl: defaultAvatar,
      defaultAvatar,
      showOptions: false,
      selectedMessageId: "",
      selectedMessageSender: ""
    };
  },

  computed: {
    chatName() {
      return this.type === "private" ? this.chat.name : this.chat.group_name;
    }
  },

  methods: {

    async openMenu(messageId, messageSender) {
      this.selectedMessageId = messageId;  
      this.selectedMessageSender = messageSender;
      this.showOptions = true;
    },

    async deleteConversation() {
      if (!this.chat || !this.chat.conversation_id) return;
      const token = localStorage.getItem("token");

      try {
        await axios.delete(`${__API_URL__}/conversations/${this.chat.conversation_id}`, {
          headers: { Authorization: `Bearer ${token}` }
        });

        eventBus.emit("conversationDeleted", this.chat.conversation_id);
        this.$emit("conversationDeleted", this.chat.conversation_id);
      } catch (error) {
        console.error("Errore nell'eliminazione della conversazione:", error);
      }
    },

    async fetchMessages() {
      if (!this.chat) return;
      this.loading = true;
      try {
        if (this.type === "private") {
          await this.fetchPrivateMessages();
        } else if (this.type === "group") {
          await this.fetchGroupMessages();
        }
        this.scrollToBottom();
      } catch (error) {
        console.error("Errore nel caricamento dei messaggi:", error);
      } finally {
        this.loading = false;
      }
    },

    async fetchPrivateMessages() {
      if (!this.chat.conversation_id) return;
      const token = localStorage.getItem("token");
      try {
        const response = await axios.get(`${__API_URL__}/conversations/${this.chat.conversation_id}?type=${this.type}`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        if (Array.isArray(response.data)) {
          this.messages = response.data.map(msg => ({
            id: msg.message_id,
            text: msg.message_content,
            sender: msg.sender_id === Number(token) ? "me" : "other",
            timestamp: new Date(msg.timestamp)
          }));
        } else {
          this.messages = [];
        }
      } catch (error) {
        console.error("Errore nel caricamento dei messaggi privati:", error);
      }
    },

    async fetchGroupMessages() {
      if (!this.chat.group_conversation_id) return;
      const token = localStorage.getItem("token");
      try {
        const response = await axios.get(`${__API_URL__}/conversations/${this.chat.group_conversation_id}?type=${this.type}`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        if (Array.isArray(response.data)) {
          this.messages = response.data.map(msg => ({
            id: msg.message_id,
            text: msg.message_content,
            sender: msg.sender_id === Number(token) ? "me" : "other",
            timestamp: new Date(msg.timestamp)
          }));
        } else {
          this.messages = [];
        }
      } catch (error) {
        console.error("Errore nel caricamento dei messaggi del gruppo:", error);
      }
    },

    async fetchUserPhoto() {
      if (!this.chat || !this.chat.recipient_id) return;
      this.avatarUrl = defaultAvatar; 
      try {
        const response = await axios.get(`${__API_URL__}/users/${this.chat.recipient_id}/photo`, {
          responseType: "blob"
        });

        if (response.data.size === 0) {
          this.avatarUrl = defaultAvatar;
          return;
        }
        const imageUrl = URL.createObjectURL(response.data);
        this.avatarUrl = ""; 
        this.$nextTick(() => {
          this.avatarUrl = imageUrl;
        });
      } catch (error) {
        console.error("Errore nel recupero della foto profilo:", error);
        this.avatarUrl = defaultAvatar;
      }
    },

    async fetchGroupPhoto() {
      if (!this.chat || this.type !== "group") return;
      this.avatarUrl = defaultAvatar; 
      try {
        const response = await axios.get(`${__API_URL__}/groups/${this.chat.group_conversation_id}/photo`, {
          responseType: "blob"
        });
        if (response.data.size === 0) {
          this.avatarUrl = defaultAvatar;
          return;
        }
        const imageUrl = URL.createObjectURL(response.data);
        this.avatarUrl = ""; 
        this.$nextTick(() => {
          this.avatarUrl = imageUrl;
        });
      } catch (error) {
        console.error("Errore nel recupero della foto del gruppo:", error);
        this.avatarUrl = defaultAvatar;
      }
    },

    async uploadGroupPhoto(event) {
      const file = event.target.files[0];
      if (!file) return;
      const token = localStorage.getItem("token");
      const formData = new FormData();
      formData.append("photo", file);
      console.log("1:",this.chat);
      try {
        await axios.put(`${__API_URL__}/groups/${this.chat.group_conversation_id}/photo`, formData, {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "multipart/form-data",
          },
        });
        this.fetchGroupPhoto();
      } catch (error) {
        console.error("Errore nell'upload della foto del gruppo:", error);
      }
    },

    async sendMessage() {
      if (this.newMessage.trim() !== "") {
        const token = localStorage.getItem("token");
        if (!token) {
          console.error("Errore: token non trovato. L'utente deve effettuare il login.");
          alert("Sessione scaduta. Effettua nuovamente il login.");
          this.$router.push("/login");
          return;
        }
        try {
          const response = await axios.post(`${__API_URL__}/messages`, {
            conversation_id: this.chat.conversation_id,
            message_content: this.newMessage
          },
          {
        headers: { Authorization: `Bearer ${token}` } 
          }
        );
          this.messages.push({
            id: response.data.message_id,
            text: this.newMessage,
            sender: "me",
            timestamp: new Date()
          });
          this.scrollToBottom();
          eventBus.emit("newMessage", {
            conversation_id: this.chat.conversation_id,
            lastMessage: this.newMessage,
            type: this.type
          })
          this.newMessage = "";
        } catch (error) {
          console.error("Errore nell'invio del messaggio:", error);
        }
      }
    },

    async sendGroupMessage(){},

    async deleteMessage(selectedMessageId) {
      if (!this.selectedMessageId) {
        console.error("Errore: nessun messaggio selezionato per l'eliminazione.");
        return;
      }
      const token = localStorage.getItem("token");
      try {
        await axios.delete(`${__API_URL__}/messages/${selectedMessageId}`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        this.messages = this.messages.filter(msg => msg.id !== this.selectedMessageId);
        this.selectedMessageId = null;
        this.showOptions = false;
        const lastMessage = this.messages.length > 0 ? this.messages[this.messages.length - 1].text : "Nessun messaggio";
        eventBus.emit("newMessage", {
          conversation_id: this.chat.conversation_id,
          lastMessage: lastMessage
        });
      } catch (error) {
          console.error("Errore nell'eliminazione del messaggio:", error);
      }
    },

    formatTime(timestamp) {
      if (!timestamp) return "";
      const date = new Date(timestamp);
      const hours = date.getHours().toString().padStart(2, "0");
      const minutes = date.getMinutes().toString().padStart(2, "0");
      return `${hours}:${minutes}`;
    },

    scrollToBottom() {
      this.$nextTick(() => {
        const container = this.$refs.messageContainer;
        if (container) {
          container.scrollTop = container.scrollHeight;
        }
      });
    }
  },
  watch: {
    chat: {
      immediate: true,
      deep: true,
      handler() {
        this.fetchMessages();
        if (this.type === "private") {
          this.fetchUserPhoto();
        } else if (this.type === "group") {
          this.fetchGroupPhoto();
        }
      }
    }
  }
};
</script>

<style scoped>
/* Contenitore principale della chat */
.w-full {
  display: flex;
  flex-direction: column;
  height: 100%; 
  width: 100%;
  background-color: #e5ddd5;
  padding: 0;
  margin: 0;
}

/* HEADER CHAT */
.chat-header {
  position: fixed;
  top: 0;
  left: 30%;
  width: 70%;
  height: 60px;
  background-color: #2f814e !important;
  color: white;
  display: flex;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid #ccc;
  z-index: 10;
}

.chat-header img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
}

.chat-header span {
  margin-left: 15px;
  font-size: 1.2rem; 
  font-weight: bold;
  color: white;
}

.delete-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 20px;
  color: white;
}

.delete-btn:hover {
  color: red;
}

.chat-messages {
  margin-top: 45px;
  margin-bottom: 45px; 
  flex-grow: 1;
  overflow-y: auto;
  word-wrap: break-word; 
  word-break: break-word;
  overflow-wrap: break-word; 
  white-space: pre-wrap;
}

.p-3 {
  padding: 1rem;
}

.message-time {
  font-size: 10px;
  color: gray;
  white-space: nowrap;
}

.w-10.h-10 {
  width: 40px;
  height: 40px;
  border-radius: 50%;
}

.font-bold.text-lg {
  font-size: 1.125rem;
  color:rgb(255, 255, 255);
}

/* Lista messaggi */
.flex-grow {
  background-color: #d6dbd6;
  overflow-y: auto;
  padding: 1rem;
}

.inline-block {
  max-width: 70%;
  word-wrap: break-word;
  padding: 0.8rem 1rem;
  box-shadow: 0px 1px 2px rgba(0, 0, 0, 0.1);
}

.bg-blue-500, .bg-gray-200 {
  background-color: #dcf8c6;
  color: black;
  border-radius: 16px;
}

.text-right {
  text-align: right;
}

/* Input messaggio */
.p-3.border-t {
  padding: 0.5rem 1rem;
  border-top: 1px solid #ddd;
  background-color: #f0f0f0;
  position: fixed;
  bottom: 0;
  right: 0;
  width: 70%;
  z-index: 20;
  display: flex;
  align-items: center;
}

input[type="text"] {
  flex-grow: 1;
  padding: 0.6rem 0.8rem;
  font-size: 14px;
  border: 1px solid #ddd;
  border-radius: 20px;
  background-color: white;
  outline: none;
  transition: box-shadow 0.2s ease;
}

input[type="text"]:focus {
  box-shadow: 0 0 5px rgba(0, 149, 246, 0.5);
}

button {
  color: rgb(255, 255, 255);
  padding: 0.6rem 0.8rem;
  font-size: 16px;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s ease;
}

button:hover {
  background-color: #069327;
}

.justify-end {
  display: flex;
  justify-content: flex-end;
}

.justify-start {
  display: flex;
  justify-content: flex-start;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-content {
  background: white;
  padding: 20px;
  border-radius: 10px;
  text-align: center;
  width: 90%;
  max-width: 400px;
}

.option-list {
  max-height: 200px;
  overflow-y: auto;
  border-top: 1px solid #ddd;
  margin-top: 10px;
}

.option-item {
  padding: 15px;
  border-bottom: 1px solid #ddd;
  cursor: pointer;
  transition: background 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  color: black;
}

.option-item:hover {
  background: #f0f0f0;
}

.cancel-btn {
  background-color: #069327;
  color: white;
  padding: 10px 15px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  margin-top: 10px;
}

.cancel-btn:hover {
  background-color: #069327;
}

.profile-picture {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  position: relative;
}

  .profile-picture:hover {
  border: 2px solid #069327; 
}

.upload-input {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: pointer;
}
  
.profile-img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

</style>
