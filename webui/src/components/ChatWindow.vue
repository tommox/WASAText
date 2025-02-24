<template>
  <div v-if="chat" class="w-full flex flex-col flex-grow bg-gray-100">
    <!-- Header Chat -->
    <div class="chat-header">
      <img :src="avatarUrl" class="w-10 h-10 rounded-full">
      <span class="ml-3 font-bold text-lg">{{ chat.name }}</span>
      <button @click="deleteConversation" class="delete-btn">
        🗑️
      </button>
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
            <span class="message-time">{{ formatTime(message.timestamp) }}</span>
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
  props: { chat: Object },
  data() {
    return { 
      newMessage: "",
      messages: [],
      loading: true,
      avatarUrl: defaultAvatar
    };
  },
  methods: {

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
    if (!this.chat || !this.chat.conversation_id) {
      console.warn("fetchMessages: conversation_id is missing");
      return;
    }
    const token = localStorage.getItem("token");
    this.loading = true;
    try {
      const response = await axios.get(`${__API_URL__}/conversations/${this.chat.conversation_id}`, {
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
      this.scrollToBottom();
    } catch (error) {
      console.error("Errore nel caricamento dei messaggi:", error);
    } finally {
      this.loading = false;
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
            lastMessage: this.newMessage
          })
          this.newMessage = "";
        } catch (error) {
          console.error("Errore nell'invio del messaggio:", error);
        }
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
        this.fetchUserPhoto();
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
}

.p-3 {
  padding: 1rem;
}

.message-time {
  font-size: 10px;
  color: gray;
  margin-left: 8px;
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

.bg-blue-500 {
  background-color: #dcf8c6;
  color: black;
  border-radius: 16px;
}

.bg-gray-200 {
  background-color: white;
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

</style>
