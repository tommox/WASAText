<template>
  <div v-if="chat" class="w-full flex flex-col h-screen bg-gray-100">
    <!-- Header Chat -->
    <div class="p-3 border-b bg-gray-200 flex items-center">
      <img :src="chat.avatar || defaultAvatar" class="w-10 h-10 rounded-full">
      <span class="ml-3 font-bold text-lg">{{ chat.name }}</span>
    </div>

    <!-- Lista Messaggi -->
    <div ref="messageContainer" class="flex-grow overflow-y-auto p-3">
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

  <div v-else class="w-full flex items-center justify-center h-screen bg-gray-100 text-gray-500">
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
      this.messages = response.data.messages || [];
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
