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
      return this.chats.filter(chat => chat.conversation_id.toString().includes(this.search.toLowerCase()));
    },
    filteredUsers() {
      return this.users.filter(user => user.Nickname.toLowerCase().includes(this.userSearch.toLowerCase()));
    }
  },
  methods: {
    async fetchChats() {
    const token = localStorage.getItem("token");
    const response = await axios.get(`${__API_URL__}/conversations`, {
      headers: { Authorization: `Bearer ${token}` }
    });
    const userResponse = await axios.get(`${__API_URL__}/users`, {
      headers: { Authorization: `Bearer ${token}` }
    });
    const allUsers = userResponse.data;
    this.chats = response.data.map(chat => {
      const recipientId = chat.sender_id === parseInt(token) ? chat.recipient_id : chat.sender_id;
      const recipient = allUsers.find(user => user.User_id === recipientId);
      return {
        ...chat,
        name: recipient ? recipient.Nickname : "Utente Sconosciuto"
      };
    });
  },

    async fetchUsers() {
      const token = localStorage.getItem("token");
      const response = await axios.get(__API_URL__ + "/users");
      if (Array.isArray(response.data)) {
        this.users = response.data.filter(user => user.User_id.toString() !== token);
      } else {
        this.users = [];
      }
      this.showUserList = true;
    },

    async startChat(user) {
    
    const existingChat = this.chats.find(
    (chat) => chat.recipient_id === user.User_id || chat.sender_id === user.User_id
  );

  if (existingChat) {
    this.$emit("chatSelected", existingChat);
    this.showUserList = false;
    return;
  }
    
    const token = localStorage.getItem("token");
    
    try {
      const response = await axios.post(
        `${__API_URL__}/conversations/conversation`,
        {recipient_id: user.User_id}, 
        {headers: { Authorization: `Bearer ${token}` }}
      );
      const conversationId = response.data.conversation_id;
      const newChat = {
      conversation_id: conversationId,
      sender_id: parseInt(token),
      recipient_id: user.User_id,
      name: user.Nickname,
      avatar: user.Avatar || null,
      lastMessage: "",
    };

    this.chats.unshift(newChat); 

    this.$emit("chatSelected", newChat);

      this.showUserList = false;
    } catch (error) {
      console.error("Errore nell'iniziare la chat:", error);
      alert("Errore: impossibile iniziare la conversazione.");
    }
  }
},

  created() {
    this.fetchChats();
  }
};
</script>

<style>

.chat-list-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f0f0f0;
}

.search-bar-container {
  display: flex;
  align-items: center;
  padding: 8px 10px;
  background-color: #ffffff;
}

.logout-btn {
  margin-left: auto; 
  padding: 10px 15px;
  background-color: #069327;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.logout-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 10px rgba(0, 0, 0, 0.1);
}

.logout-btn:active {
  transform: translateY(0);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.search-input {
  flex-grow: 1;
  padding: 10px;
  font-size: 14px;
  border: none;
  border-radius: 20px;
  outline: none;
  background-color: #fff;
  margin-right: 10px;
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
  transition: background-color 0.2s;
}

.new-chat-btn:hover {
  background-color: #069327;
}

.chat-list {
  flex-grow: 1;
  overflow-y: auto;
  background-color: #fff;
}

.chat-list::-webkit-scrollbar {
  width: 8px;
}

.chat-list::-webkit-scrollbar-thumb {
  background-color: rgba(0, 0, 0, 0.2);
  border-radius: 4px;
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
  z-index: 999;
}

.modal-content {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  width: 90%;
  max-width: 400px;
}

.modal-input {
  width: 100%;
  padding: 10px;
  margin: 10px 0;
  border: 1px solid #ddd;
  border-radius: 20px;
  outline: none;
}

.user-list {
  max-height: 200px;
  overflow-y: auto;
  border-top: 1px solid #ddd;
  margin-top: 10px;
}

.user-item {
  padding: 15px;
  border-bottom: 1px solid #ddd;
  cursor: pointer;
  transition: background 0.3s;
  display: flex;
  align-items: center;
  gap: 10px;
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
