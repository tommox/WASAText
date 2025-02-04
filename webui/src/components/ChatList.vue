<template>
  <div class="chat-list-container">
    <!-- Barra di ricerca + Pulsante "Nuova Conversazione" -->
    <div class="search-bar-container flex items-center border-b bg-gray-100">
      <input 
        v-model="search" 
        type="text" 
        placeholder="Cerca chat" 
        class="flex-grow px-3 py-2 border rounded-lg"
      />
      <button @click="fetchUsers" class="new-chat-btn">
        ➕
      </button>
    </div>

    <!-- Lista Chat -->
    <div class="chat-list">
      <ChatItem 
        v-for="chat in filteredChats" 
        :key="chat.id" 
        :chat="chat" 
        @selectChat="$emit('chatSelected', chat)"
      />
    </div>
  </div>

  <!-- Modale per scegliere un utente -->
  <div v-if="showUserList" class="modal-overlay">
    <div class="modal-content">
      <h2>Seleziona un utente</h2>
      <input v-model="userSearch" type="text" placeholder="Cerca utente..." class="modal-input" />

      <div class="user-list">
        <div 
          v-for="user in filteredUsers" 
          :key="user.id" 
          class="user-item"
          @click="startChat(user)"
        >
          {{ user.name }}
        </div>
      </div>

      <button @click="showUserList = false" class="cancel-btn">Chiudi</button>
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
      search: "", 
      chats: [],
      users: [], // Lista utenti dal database
      showUserList: false,
      userSearch: "", // Per filtrare gli utenti
      loading: true
    };
  },
  computed: {
    filteredChats() {
      return this.chats.filter(chat => 
        chat.name.toLowerCase().includes(this.search.toLowerCase())
      );
    },
    filteredUsers() {
      return this.users.filter(user => 
        user.name.toLowerCase().includes(this.userSearch.toLowerCase())
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
    },

    async fetchUsers() {
      try {
        const response = await axios.get("/users"); // API per ottenere gli utenti dal DB
        this.users = response.data;
        this.showUserList = true; // Mostra il modale
      } catch (error) {
        console.error("Errore nel caricamento degli utenti:", error);
      }
    },

    startChat(user) {
      this.$emit("chatSelected", user);
      this.showUserList = false;
    }
  },
  created() {
    this.fetchChats();
  }
};
</script>

<style scoped>
/* Contenitore della lista chat */
.chat-list-container {
  display: flex;
  flex-direction: column;
  max-height: 100vh;
  overflow: hidden;
}

/* Pulsante nuova conversazione */
.new-chat-container {
  padding: 10px;
  text-align: center;
}

/* Pulsante nuovo chat (piccolo quadrato accanto alla barra di ricerca) */
.new-chat-btn {
  width: 40px;
  height: 40px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 5px;
  font-size: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-left: 10px;
}

.new-chat-btn:hover {
  background-color: #0056b3;
}


/* Lista chat */
.chat-list {
  flex-grow: 1;
  min-height: 0;
  overflow-y: auto;
}

/* Modale per selezionare un utente */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-content {
  background: white;
  padding: 20px;
  border-radius: 10px;
  text-align: center;
  width: 300px;
}

.modal-input {
  width: 100%;
  padding: 10px;
  margin: 10px 0;
  border: 1px solid #ddd;
  border-radius: 5px;
}

/* Lista utenti */
.user-list {
  max-height: 200px;
  overflow-y: auto;
  border-top: 1px solid #ddd;
  margin-top: 10px;
}

.user-item {
  padding: 10px;
  border-bottom: 1px solid #ddd;
  cursor: pointer;
  transition: background 0.3s;
}

.user-item:hover {
  background: #f0f0f0;
}

.cancel-btn {
  background-color: #dc3545;
  color: white;
  padding: 10px 15px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  margin-top: 10px;
}

.cancel-btn:hover {
  background-color: #c82333;
}
</style>
