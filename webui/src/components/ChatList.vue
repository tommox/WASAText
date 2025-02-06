<template>
  <div class="chat-list-container">
    <div class="search-bar-container flex items-center border-b bg-gray-100 p-2">
      <input 
        v-model="search" 
        type="text" 
        placeholder="Cerca chat" 
        class="search-input flex-grow px-4 py-2 border rounded-lg"
      />
      <button @click="fetchUsers" class="new-chat-btn">➕</button>
    </div>

    <div class="chat-list">
      <ChatItem 
        v-for="chat in filteredChats" 
        :key="chat.id" 
        :chat="chat" 
        @selectChat="$emit('chatSelected', chat)"
      />
    </div>
  </div>

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
          {{ user.Nickname }}
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
  emits: ["chatSelected"],
  components: { ChatItem },
  data() {
    return {
      search: "", 
      chats: [],
      users: [],
      showUserList: false,
      userSearch: ""
    };
  },
  computed: {
    filteredChats() {
      return this.chats.filter(chat => chat.id.toString().includes(this.search.toLowerCase()));
    },
    filteredUsers() {
      return this.users.filter(user => user.Nickname.toLowerCase().includes(this.userSearch.toLowerCase()));
    }
  },
  methods: {
    async fetchChats() {
      const token = localStorage.getItem("token"); 
      const response = await axios.get(__API_URL__+"/conversations", {
        headers: { Authorization: `Bearer ${token}` }
      });
      this.chats = response.data ?? [];
    },

    async fetchUsers() {
      const response = await axios.get(__API_URL__+"/users");
      this.users = Array.isArray(response.data) ? response.data : [];
      this.showUserList = true;
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
.chat-list-container {
  display: flex;
  flex-direction: column;
  max-height: 100vh;
  overflow: hidden;
}

.search-input {
  flex-grow: 1;
  width: 80%;
  padding: 10px;
  font-size: 16px;
}

.new-chat-btn {
  width: 36px;
  height: 36px;
  background-color: #ebe9e9;
  color: white;
  border: none;
  border-radius: 50%;
  font-size: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-left: 300px;
  margin-top: -38px;
}

.new-chat-btn:hover {
  background-color: #d6d6d6;
}

.chat-list {
  flex-grow: 1;
  min-height: 0;
  overflow-y: auto;
}

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
</style>
