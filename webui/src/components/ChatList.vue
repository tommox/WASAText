<template>
  <div class="chat-list-container">
    <div class="search-bar-container flex items-center border-b bg-gray-100 p-2">
      <input 
        v-model="search" 
        type="text" 
        placeholder="Cerca chat" 
        class="search-input flex-grow px-4 py-2 border rounded-lg"
      />
      <button @click="showChatOptions = true" class="new-chat-btn">➕</button>
    </div>

    <div class="chat-list">
      <ChatItem 
        v-for="chat in filteredChats" 
        :key="'chat-' + chat.conversation_id" 
        :chat="chat" 
        :type="'private'"
        @selectChat="$emit('chatSelected', chat)"
      />
      <ChatItem 
        v-for="group in filteredGroupChats" 
        :key="'group-' + group.group_conversation_id" 
        :chat="group" 
        :type="'group'"
        @selectChat="$emit('chatSelected', group)"
      />
    </div>
  </div>

  <!-- Modale per la selezione tra Nuova Chat e Nuovo Gruppo -->
  <div v-if="showChatOptions" class="modal-overlay">
    <div class="modal-content">
      <h2>Seleziona un'opzione</h2>
        <div class="user-list">
          <div @click="fetchUsers" class="user-item">Nuova Chat</div>
        <div @click="fetchGroupUsers" class="user-item">Nuovo Gruppo</div>
      </div>
      <button @click="showChatOptions = false" class="cancel-btn">Chiudi</button>
    </div>
  </div>

  <!-- Modale per la lista utenti -->
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

    <!-- Modale per la creazione del gruppo -->
    <div v-if="showGroupUserList" class="modal-overlay">
    <div class="modal-content">
      <h2>Crea un nuovo gruppo</h2>
      <input v-model="groupName" type="text" placeholder="Nome del gruppo" class="modal-input" />
      <div class="user-list">
        <div 
          v-for="user in filteredUsers" 
          :key="user.id" 
          class="user-item"
          :class="{ 'selected': selectedUsers.includes(user.User_id) }"
          @click="toggleUserSelection(user.User_id)"
        >
          {{ user.Nickname }}
        </div>
      </div>
      <button @click="createGroup" class="confirm-btn" :disabled="!canCreateGroup">Crea</button>
      <button @click="closeGroupUserList" class="cancel-btn">Chiudi</button>
    </div>
  </div>
</template>



<script>
import ChatItem from "./ChatItem.vue";
import axios from "axios";
import eventBus from "@/eventBus";
import defaultAvatar from "@/assets/images/user.png";

export default {
  emits: ["chatSelected"],
  components: { ChatItem },
  data() {
    return {
      search: "", 
      chats: [],
      groupChats: [],
      users: [],
      showUserList: false,
      showGroupUserList: false,
      showChatOptions: false,
      userSearch: "",
      groupName: "",
      selectedUsers: []
    };
  },
  created() {
    this.fetchChats();
    this.fetchGroupChats();
    eventBus.on("conversationDeleted", (conversationId) => {
      this.chats = this.chats.filter(chat => chat.conversation_id !== conversationId);
      this.groupChats = this.groupChats.filter(chat => chat.g_conversation_id !== conversationId);
    });
  },
  beforeUnmount() {
    eventBus.off("conversationDeleted");
  },
  computed: {
    filteredChats() {
      return this.chats.filter(chat => 
        chat.conversation_id?.toString().includes(this.search.toLowerCase()) || 
        chat.name?.toLowerCase().includes(this.search.toLowerCase())
      );
    },
    filteredGroupChats() {
      return this.groupChats.filter(group => 
        group.group_conversation_id?.toString().includes(this.search.toLowerCase()) || 
        group.group_name?.toLowerCase().includes(this.search.toLowerCase())
      );
    },
    filteredUsers() {
      return this.users.filter(user => user.Nickname.toLowerCase().includes(this.userSearch.toLowerCase()));
    },
    canCreateGroup() {
      return this.groupName.trim() !== "" && this.selectedUsers.length > 0;
    },
    allChats() {
      return [...this.filteredChats, ...this.filteredGroupChats];
    }
  },
  methods: {

    async fetchGroupUsers() {
      this.showChatOptions = false;
      try {
        const token = localStorage.getItem("token");
        const response = await axios.get(`${__API_URL__}/users`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        if (Array.isArray(response.data)) {
          this.users = response.data.filter(user => user.User_id.toString() !== token);
        } else {
          this.users = [];
        }
        this.showGroupUserList = true; 
      } catch (error) {
        console.error("Errore nel recupero degli utenti per il gruppo:", error);
        alert("Errore nel recupero degli utenti.");
      }
    },


    async fetchChats() {
      const token = localStorage.getItem("token");
      if (!token) return;
      try {
        const response = await axios.get(`${__API_URL__}/conversations`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        const userResponse = await axios.get(`${__API_URL__}/users`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        const allUsers = userResponse.data;
        if (response.data && Array.isArray(response.data.private_conversations)) {
          this.chats = response.data.private_conversations.map(chat => {
            if (!chat.conversation_id) {
              console.error("Conversazione senza ID:", chat);
              return null;
            }
            const isCurrentUserSender = chat.sender_id === parseInt(token);
            const recipientId = isCurrentUserSender ? chat.recipient_id : chat.sender_id;
            const recipient = allUsers.find(user => user.User_id === recipientId);
            return {
              ...chat,
              recipient_id: recipientId,
              name: recipient ? recipient.Nickname : "Utente Sconosciuto",
              avatarUrl: recipient ? `${__API_URL__}/users/${recipientId}/photo` : defaultAvatar,
              lastMessage: chat.last_message_id ? parseInt(chat.last_message_id) : "Nessun messaggio",
            };
          }).filter(chat => chat !== null);
        } else {
          console.error("Formato della risposta inatteso:", response.data);
          this.chats = [];
        }
        console.log("Chats: ", this.chats);
      } catch (error) {
        console.error("Errore nel recupero delle conversazioni private:", error);
      }
    },

    async fetchGroupChats() {
      const token = localStorage.getItem("token");
      if (!token) return;
      try {
        const response = await axios.get(`${__API_URL__}/conversations`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        if (response.data && Array.isArray(response.data.group_conversations)) {
          this.groupChats = response.data.group_conversations.map(group => ({
            group_conversation_id: group.group_conversation_id, 
            group_name: group.group_name || "Gruppo Sconosciuto",
            group_avatarUrl: `${__API_URL__}/groups/${group.group_id}/photo`,
            group_last_message_id: group.last_message_id || "Nessun messaggio",
          }));
        }
      } catch (error) {
        console.error("Errore nel recupero delle conversazioni di gruppo:", error);
      }
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
      this.showChatOptions = false;
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
      this.chats = [...this.chats, newChat];
      this.$emit("chatSelected", newChat);
      this.showChatOptions = false;
      } catch (error) {
        console.error("Errore nell'iniziare la chat:", error);
        alert("Errore: impossibile iniziare la conversazione.");
      }
    },

    toggleUserSelection(userId) {
      if (this.selectedUsers.includes(userId)) {
        this.selectedUsers = this.selectedUsers.filter(id => id !== userId);
      } else {
        this.selectedUsers.push(userId);
      }
    },

    async createGroup() {
      if (!this.canCreateGroup) return;

      try {
        const token = localStorage.getItem("token");
        const response = await axios.post(`${__API_URL__}/groups`, {
          group_name: this.groupName 
        }, { headers: { Authorization: `Bearer ${token}` }
        });

        const groupId = response.data.group_id;
        for (const userId of this.selectedUsers) {
          await axios.post(
            `${__API_URL__}/groups/${groupId}/users/${userId}?state=add`,
            { role: "member" },
            { headers: { Authorization: `Bearer ${token}` } }
          );
        }
        const newGroup = {
          conversation_id: response.data.group_id, 
          name: this.groupName,
          avatarUrl: "default_group_avatar.png",
          lastMessage: ""
        };
        this.chats.push(newGroup);
        this.closeGroupUserList();
      } catch (error) {
        console.error("Errore nella creazione del gruppo:", error);
      }
    },

    closeGroupUserList() {
      this.selectedUsers = [];
      this.groupName = "";
      this.showGroupUserList = false;
      this.showChatOptions = false;
    },

    handleConversationDeleted(conversationId) {
      this.chats = this.chats.filter(chat => chat.conversation_id !== conversationId);
    }
  },
};
</script>

<style>

.chat-list-container {
  display: flex;
  flex-direction: column;
  height: 100%;
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

.selected {
  background-color: lightgreen !important;
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

.confirm-btn {
  background-color: #069327;
  color: white;
  padding: 10px 15px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  margin-top: 10px;
}

.confirm-btn:disabled {
  background-color: grey !important;
  cursor: not-allowed;
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
