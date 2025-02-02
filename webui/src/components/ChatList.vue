<template>
    <div class="w-1/3 bg-white bg-opacity-75 flex flex-col border-r border-black backdrop-blur-md">
      <!-- Logo -->
      <div class="p-4 flex items-center justify-center border-b">
        <img
          src="@/assets/images/WASAText.png"
          alt="WASAText logo"
          style="width: 400px; height: 150px;"
        />
      </div>
  
      <!-- Barra di ricerca -->
      <div class="p-3 flex items-center border-b bg-gray-100">
        <input 
          v-model="search" 
          type="text" 
          placeholder="Cerca chat" 
          class="flex-grow px-3 py-2 border rounded-lg"
        />
      </div>
  
      <!-- Lista Chat -->
      <div v-if="loading" class="flex justify-center items-center py-5">
        <span class="text-gray-500">Caricamento chat...</span>
      </div>
  
      <div v-else-if="chats.length === 0" class="flex justify-center items-center py-5">
        <span class="text-gray-500">Nessuna chat trovata</span>
      </div>
  
      <div v-else class="overflow-y-auto flex-grow">
        <ChatItem 
          v-for="chat in filteredChats" 
          :key="chat.id" 
          :chat="chat" 
          @selectChat="$emit('chatSelected', chat)"
        />
      </div>
    </div>
  </template>
  
  <script>
  import ChatItem from "./ChatItem.vue";
  import axios from "axios";
  
  export default {
    components: { ChatItem },
    data() {
      return {
        search: "",  // ✅ Definiamo search dentro data()
        chats: [],
        loading: true
      };
    },
    computed: {
      filteredChats() {
        return this.chats.filter(chat => 
          chat.name.toLowerCase().includes(this.search.toLowerCase())
        );
      }
    },
    methods: {
      async fetchChats() {
        try {
          const response = await axios.get("/conversations");
          this.chats = response.data;
        } catch (error) {
          console.error("Errore nel caricamento delle conversazioni:", error);
        } finally {
          this.loading = false;
        }
      }
    },
    created() {
      this.fetchChats();
    }
  };
  </script>
  