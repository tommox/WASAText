<template>
  <div v-if="chat" class="w-full flex flex-col flex-grow bg-gray-100">
    <!-- Header Chat -->
    <div class="chat-header">
      <img :src="chat.avatar || defaultAvatar" class="w-10 h-10 rounded-full">
      <span class="ml-3 font-bold text-lg">{{ chat.name }}</span>
    </div>

    <!-- Lista Messaggi -->
    <div ref="messageContainer" class="chat-messages flex-grow overflow-y-auto p-3">
      <div v-if="loading" class="text-center text-gray-500">Caricamento messaggi...</div>
      <div v-else-if="messages.length === 0" class="text-center text-gray-500">Inizia una nuova conversazione!</div>

      <div 
        v-else v-for="message in messages" 
        :key="message.id" 
        class="mb-2"
        :class="{ 'text-right': message.sender === 'me' }"
      >
        <div 
          class="inline-block p-3 rounded-lg shadow-md"
          :class="message.sender === 'me' ? 'bg-blue-500 text-white' : 'bg-gray-200'"
        >
          {{ message.text }}
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

export default {
  props: { chat: Object },
  data() {
    return { 
      newMessage: "",
      messages: [],
      loading: true,
      defaultAvatar,
    };
  },
  methods: {
    async fetchMessages() {
    if (!this.chat || !this.chat.conversation_id) {
      console.warn("fetchMessages: conversation_id is missing");
      return;
    }
    const token = localStorage.getItem("token");
    this.loading = true;
    try {
      const response = await axios.get(`${__API_URL__}/conversations/${this.chat.conversation_id}`, {
        headers: {
          Authorization: `Bearer ${token}`
        }
    });
    this.messages = response.data && Array.isArray(response.data.messages) ? response.data.messages : [];
      this.scrollToBottom();
    } catch (error) {
      console.error("Errore nel caricamento dei messaggi:", error);
    } finally {
      this.loading = false;
    }
},
    async sendMessage() {
      if (this.newMessage.trim() !== "") {
        try {
          const response = await axios.post("/messages", {
            conversation_id: this.chat.id,
            text: this.newMessage
          });

          this.messages.push({
            id: response.data.message_id,
            text: this.newMessage,
            sender: "me"
          });

          this.newMessage = "";
          this.scrollToBottom();
        } catch (error) {
          console.error("Errore nell'invio del messaggio:", error);
        }
      }
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
      handler() {
        this.fetchMessages();
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

.chat-messages {
  margin-top: 60px;
}

.p-3 {
  padding: 1rem;
}

.w-10.h-10 {
  width: 40px;
  height: 40px;
  border-radius: 50%;
}

.font-bold.text-lg {
  font-size: 1.125rem;
  color:white;
}

/* Lista messaggi */
.flex-grow {
  background-color: #d6dbd6;
  overflow-y: auto;
  padding: 1rem;
}

.mb-2 {
  margin-bottom: 0.5rem;
}

.inline-block {
  max-width: 70%;
  word-wrap: break-word;
  padding: 0.8rem 1rem;
  border-radius: 8px;
  box-shadow: 0px 1px 2px rgba(0, 0, 0, 0.1);
}

.bg-blue-500 {
  background-color: #dcf8c6; /* Colore messaggio inviato tipo WhatsApp */
  color: black;
}

.bg-gray-200 {
  background-color: white; /* Colore messaggio ricevuto tipo WhatsApp */
  color: black;
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
</style>
