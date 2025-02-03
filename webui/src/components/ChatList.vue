<template>
  <div class="chat-list-container">
    <!-- Barra di ricerca (più in alto) -->
    <div class="search-bar-container p-2 flex items-center border-b bg-gray-100">
      <input 
        v-model="search" 
        type="text" 
        placeholder="Cerca chat" 
        class="flex-grow px-3 py-2 border rounded-lg"
      />
    </div>
  </div>

  <!-- Lista Chat con comportamento dinamico -->
  <div class="chat-list">
    <ChatItem 
      v-for="chat in filteredChats" 
      :key="chat.id" 
      :chat="chat" 
      @selectChat="$emit('chatSelected', chat)"
    />
  </div>
</template>


<script>
import ChatItem from "./ChatItem.vue";
import axios from "axios";

export default {
  components: { ChatItem },
  data() {
    return {
      search: "", 
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

<style scoped>

/* Stile per la lista di chat */
.chat-list-container {
  display: flex;
  flex-direction: column;
  max-height: 100vh; /* Limita l'altezza alla schermata */
  overflow: hidden; /* Nessuno scroll di default */
}

.chat-list {
  flex-grow: 1; /* Occupa lo spazio disponibile */
  min-height: 0; /* Necessario per gestire il flexbox con overflow */
  overflow-y: auto; /* Scorrimento verticale solo quando necessario */
}

.search-bar-container {
  padding-top: 5px; /* Ridotto rispetto al default */
  padding-bottom: 10px;
}

/* Stile per la barra di ricerca */
input {
  width: 100%;
  padding: 10px;
  border: 2px solid #ddd;
  border-radius: 20px;
  font-size: 16px;
  transition: all 0.3s ease-in-out;
}

input:focus {
  border-color: #007bff;
  box-shadow: 0 0 8px rgba(0, 123, 255, 0.5);
  outline: none;
}

/* Stile per la lista chat */
.chat-item {
  padding: 15px;
  border-bottom: 1px
}
</style>