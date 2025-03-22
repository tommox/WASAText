<template>
  <div class="container-fluid">
    <!-- PANNELLO SINISTRO: Profilo utente e lista chat -->
    <div class="left-panel">
      <!-- Profilo utente -->
      <div class="user-profile">
        <label class="profile-picture">
          <img :src="userImage || defaultAvatar" alt="Foto Profilo" class="profile-img" />
          <input type="file" @change="uploadProfilePicture" accept="image/*" class="upload-input" />
        </label>
        <div class="user-name-container" @click="editNickname">
          <span v-if="!isEditing" class="user-name">{{ nickname }}</span>
          <input
            v-else
            ref="editableNickname"
            v-model="editableNickname"
            class="user-name-input"
            @blur="saveNickname"
            @keyup.enter="saveNickname"
          />
        </div>
        <button @click="logout" class="logout-btn">Logout</button>
      </div>

      <!-- Lista Chat -->
      <div class="chat-list-container">
        <div class="search-bar-container">
          <input v-model="search" type="text" placeholder="Cerca chat" class="search-input" />
          <button @click="showChatOptions = true" class="new-chat-btn">➕</button>
        </div>
        <div class="chat-list">
          <!-- Chat private -->
          <div
            v-for="chat in filteredChats"
            :key="'chat-' + chat.conversation_id"
            class="chat-item"
            @click="selectChat(chat, 'private')">
            <img :src="chat.avatarUrl" alt="Avatar" class="profile-img" />
            <div class="chat-details">
              <div class="chat-name">{{ chat.name }}</div>
              <div class="chat-last-message">{{ chat.lastMessage || 'Nessun messaggio' }}</div>
            </div>
          </div>
          <!-- Chat di gruppo -->
          <div
            v-for="group in filteredGroupChats"
            :key="'group-' + group.group_conversation_id"
            class="chat-item"
            @click="selectChat(group, 'group')"
          >
            <img :src="group.group_avatarUrl" alt="Foto Gruppo" class="profile-img" />
            <div class="chat-details">
              <div class="chat-name">{{ group.group_name }}</div>
              <div class="group-chat-last-message">{{ group.group_lastMessage || 'Nessun messaggio' }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Modali per Nuova Chat / Gruppo -->
      <div v-if="showChatOptions" class="modal-overlay">
        <div class="modal-content">
          <h2>Seleziona un'opzione</h2>
          <div class="user-list">
            <div @click="fetchUsersForChat" class="user-item">Nuova Chat</div>
            <div @click="fetchUsersForGroup" class="user-item">Nuovo Gruppo</div>
          </div>
          <button @click="showChatOptions = false" class="cancel-btn">Chiudi</button>
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
    </div>

    <!-- PANNELLO DESTRO: Finestra chat -->
    <div class="right-panel">
      <div v-if="selectedChat" class="chat-window">
        <!-- Header Chat -->
        <div class="chat-header">
          <div v-if="selectedChatType === 'group'" class="profile-picture">
            <img :src="groupImage || defaultAvatar" alt="Foto Gruppo" class="profile-img" />
            <input type="file" @change="uploadGroupPhoto" accept="image/*" class="upload-input" />
          </div>
          <div v-else>
            <img :src="avatarUrl" alt="Avatar" class="profile-img" />
          </div>
          <span class="chat-header-name">
            {{ selectedChatType === 'private' ? selectedChat.name : selectedChat.group_name }}
          </span>
          <button
            @click="selectedChatType === 'private' ? deleteConversation() : deleteGroupConversation()"
            class="delete-btn">
            🗑️
          </button>
          <div v-if="showGroupErrorModal" class="modal-overlay">
            <div class="modal-content">
              <h2 class="modal-error-title">Errore</h2>
                <p class="modal-error-box">{{ groupErrorMessage }}</p>
              <button class="confirm-btn" @click="showGroupErrorModal = false">OK</button>
            </div>
          </div>
        </div>

        <!-- Lista messaggi -->
        <div ref="messageContainer" class="chat-messages">
          <div v-if="loading" class="loading-text">Caricamento messaggi...</div>
          <div v-else-if="messages.length === 0" class="loading-text">Inizia una nuova conversazione!</div>
          <div
            v-else
            v-for="message in messages"
            :key="message.id"
            class="message"
            :class="{'message-me': message.sender === 'me', 'message-other': message.sender === 'other'}"
          >
            <div class="message-content">
              <strong v-if="selectedChatType === 'group'">
                {{ message.sender === 'me' ? 'Tu' : userMap[message.rawSenderId] || 'Utente' }}
              </strong>
              <br />
              <span>{{ message.text }}</span>
              <div class="message-time" @click="openMessageMenu(message.id, message.sender)">
                {{ formatTime(message.timestamp) }}
              </div>
              <div v-if="showOptions && selectedMessageId === message.id" class="modal-overlay message-options">
                <div class="modal-content">
                  <h2>Seleziona un'opzione</h2>
                  <div class="option-list">
                    <div
                      v-if="selectedMessageSender === 'me'"
                      class="option-item"
                      @click="deleteMessage(selectedMessageId)"
                    >
                      Elimina
                    </div>
                  </div>
                  <button @click="showOptions = false" class="cancel-btn">Chiudi</button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Input per inviare messaggio -->
        <div class="message-input-container">
          <input
            v-model="newMessage"
            type="text"
            placeholder="Scrivi un messaggio..."
            @keyup.enter="sendCurrentMessage"
          />
          <button @click="sendCurrentMessage">➤</button>
        </div>
      </div>
      <div v-else class="empty-chat">
        Apri o inizia una nuova conversazione
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import eventBus from "@/eventBus";
import defaultAvatar from "@/assets/images/user.png";

// Funzione helper per convertire un blob in Base64 (utilizzata per foto profilo e gruppi)
function blobToBase64(blob) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onloadend = () => resolve(reader.result);
    reader.onerror = reject;
    reader.readAsDataURL(blob);
  });
}
//this.userMap = Object.fromEntries(this.users.map(user => [user.User_id, user.Nickname]));
export default {
  name: "HomeView",
  data() {
    return {
      // Profilo utente
      nickname: localStorage.getItem("nickname") || "Utente",
      userImage: defaultAvatar, // Inizializzato a defaultAvatar, verrà aggiornato dal backend
      editableNickname: "",
      isEditing: false,
      showGroupErrorModal: false,
      groupErrorMessage: "",
      // Lista chat
      search: "",
      chats: [],
      groupChats: [],
      userMap: {},
      // Selezione chat
      selectedChat: null,
      selectedChatType: "",
      // Modali per creare chat/gruppo
      showChatOptions: false,
      showUserList: false,
      showGroupUserList: false,
      userSearch: "",
      users: [],
      groupName: "",
      selectedUsers: [],
      // Finestra chat
      messages: [],
      newMessage: "",
      loading: false,
      avatarUrl: defaultAvatar, // Usato per la chat window (utente privato)
      groupImage: defaultAvatar, // Usato per la chat window dei gruppi
      // Opzioni messaggio
      showOptions: false,
      selectedMessageId: null,
      selectedMessageSender: null,
    };
  },
  computed: {
    filteredChats() {
      return this.chats.filter(
        chat =>
          (chat.name && chat.name.toLowerCase().includes(this.search.toLowerCase())) ||
          (chat.conversation_id && chat.conversation_id.toString().includes(this.search.toLowerCase()))
      );
    },
    filteredGroupChats() {
      return this.groupChats.filter(
        group =>
          (group.group_name && group.group_name.toLowerCase().includes(this.search.toLowerCase())) ||
          (group.group_conversation_id && group.group_conversation_id.toString().includes(this.search.toLowerCase()))
      );
    },
    filteredUsers() {
      return this.users.filter(user =>
        user.Nickname.toLowerCase().includes(this.userSearch.toLowerCase())
      );
    },
    canCreateGroup() {
      return this.groupName.trim() !== "" && this.selectedUsers.length > 0;
    },
  },
  methods: {
    // --- Recupero delle conversazioni ---
    async fetchChats() {
      const token = localStorage.getItem("token");
      if (!token) return;
      try {
        // Richiama le conversazioni e gli utenti
        const response = await axios.get(`${__API_URL__}/conversations`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        const userResponse = await axios.get(`${__API_URL__}/users`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        const allUsers = userResponse.data;
        
        // Per ogni conversazione privata, se presente last_message_id, richiedi il contenuto
        let chats = [];
        if (Array.isArray(response.data.private_conversations)) {
          chats = await Promise.all(response.data.private_conversations.map(async (chat) => {
            if (!chat.conversation_id) return null;
            const isCurrentUserSender = chat.sender_id === parseInt(token);
            const recipientId = isCurrentUserSender ? chat.recipient_id : chat.sender_id;
            const recipient = allUsers.find(user => user.User_id === recipientId);
            let avatarUrl = defaultAvatar;
            if (recipient) {
              try {
                const photoResponse = await axios.get(`${__API_URL__}/users/${recipientId}/photo`, {
                  headers: { Authorization: `Bearer ${token}` },
                  responseType: "blob",
                });
                if (photoResponse.data && photoResponse.data.size > 0) {
                  avatarUrl = await blobToBase64(photoResponse.data);
                }
              } catch (e) {
                console.warn("Errore nel recupero della foto utente:", e);
                avatarUrl = defaultAvatar;
              }
            }
            let lastMessage = "Nessun messaggio";
            if (chat.last_message_id) {
              try {
                const msgResponse = await axios.get(
                  `${__API_URL__}/messages/${chat.last_message_id}?type=private`,
                  { headers: { Authorization: `Bearer ${token}` } }
                );
                if (msgResponse.data && msgResponse.data.message_content) {
                  lastMessage = msgResponse.data.message_content;
                }
              } catch (error) {
                console.error("Errore nel recupero dell'ultimo messaggio", error);
              }
            }
            return {
              ...chat,
              recipient_id: recipientId,
              name: recipient ? recipient.Nickname : "Utente Sconosciuto",
              avatarUrl,
              lastMessage
            };
          }));
          chats = chats.filter(chat => chat !== null);
        }
        this.chats = chats;
      } catch (error) {
        console.error("Errore nel recupero delle conversazioni private:", error);
        this.chats = [];
      }
    },
    async fetchGroupChats() {
      const token = localStorage.getItem("token");
      if (!token) return;
      try {
        const response = await axios.get(`${__API_URL__}/conversations`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        let groupChats = [];
        if (Array.isArray(response.data.group_conversations)) {
          groupChats = await Promise.all(response.data.group_conversations.map(async group => {
            const groupConversationId = group.group_conversation_id;
            let lastMessage = "Nessun messaggio";
            if (group.last_message_id) {
              try {
                const msgResponse = await axios.get(
                  `${__API_URL__}/messages/${group.last_message_id}?type=group`,
                  { headers: { Authorization: `Bearer ${token}` } }
                );
                if (msgResponse.data && msgResponse.data.message_content) {
                  lastMessage = msgResponse.data.message_content;
                }
              } catch (error) {
                console.error("Errore nel recupero dell'ultimo messaggio per il gruppo", error);
              }
            }
            return {
              group_conversation_id: groupConversationId,
              group_name: group.group_name || "Gruppo Sconosciuto",
              group_avatarUrl: `${__API_URL__}/groups/${groupConversationId}/photo`,
              group_lastMessage: lastMessage,
            };
          }));
        }
        this.groupChats = groupChats;
      } catch (error) {
        console.error("Errore nel recupero delle conversazioni di gruppo:", error);
        this.groupChats = [];
      }
    },
    async fetchGroupMembers(groupId) {
      const token = localStorage.getItem("token");
      try {
        const response = await axios.get(`${__API_URL__}/groups/${groupId}/users`, {
          headers: { Authorization: `Bearer ${token}` }
        });

        const members = response.data;
        this.userMap = Object.fromEntries(members.map(user => [user.user_id, user.nickname]));
      } catch (error) {
        console.error("Errore nel recupero dei membri del gruppo:", error);
      }
    },
    selectChat(chat, type) {
      this.selectedChat = chat;
      this.selectedChatType = type;
      this.fetchMessages();
      if (type === "private") {
        this.fetchUserPhoto();
      } else if (type === "group") {
        this.fetchGroupPhoto();
        this.fetchGroupMembers(chat.group_conversation_id);
      }
    },
    async fetchMessages() {
      if (!this.selectedChat) return;
      this.loading = true;
      const token = localStorage.getItem("token");
      try {
        if (this.selectedChatType === "private" && this.selectedChat.conversation_id) {
          const response = await axios.get(
            `${__API_URL__}/conversations/${this.selectedChat.conversation_id}?type=private`,
            { headers: { Authorization: `Bearer ${token}` } }
          );
          if (Array.isArray(response.data)) {
            this.messages = response.data.map(msg => ({
              id: msg.message_id,
              text: msg.message_content,
              sender: msg.sender_id === Number(token) ? "me" : "other",
              timestamp: new Date(msg.timestamp),
            }));
          } else {
            this.messages = [];
          }
        } else if (this.selectedChatType === "group" && this.selectedChat.group_conversation_id) {
          const response = await axios.get(
            `${__API_URL__}/conversations/${this.selectedChat.group_conversation_id}?type=group`,
            { headers: { Authorization: `Bearer ${token}` } }
          );
          if (Array.isArray(response.data)) {
            this.messages = response.data.map(msg => ({
              id: msg.message_id,
              text: msg.message_content,
              sender: msg.sender_id === Number(token) ? "me" : "other",
              rawSenderId: msg.sender_id,
              timestamp: new Date(msg.timestamp),
            }));
          } else {
            this.messages = [];
          }
        }
        this.scrollToBottom();
      } catch (error) {
        console.error("Errore nel caricamento dei messaggi:", error);
      } finally {
        this.loading = false;
      }
    },
    async sendCurrentMessage() {
      if (this.newMessage.trim() === "" || !this.selectedChat) return;
      const token = localStorage.getItem("token");
      if (!token) {
        alert("Sessione scaduta. Effettua nuovamente il login.");
        this.$router.push("/login");
        return;
      }
      try {
        let response;
        if (this.selectedChatType === "private") {
          response = await axios.post(
            `${__API_URL__}/messages`,
            {
              conversation_id: this.selectedChat.conversation_id,
              message_content: this.newMessage,
            },
            { headers: { Authorization: `Bearer ${token}` } }
          );
        } else if (this.selectedChatType === "group") {
          response = await axios.post(
            `${__API_URL__}/groups/${this.selectedChat.group_conversation_id}/messages`,
            { message_content: this.newMessage },
            { headers: { Authorization: `Bearer ${token}` } }
          );
        }
        this.messages.push({
          id: response.data.message_id,
          text: this.newMessage,
          sender: "me",
          timestamp: new Date(),
        });
        if (this.selectedChatType === "private") {
          this.selectedChat.lastMessage = this.newMessage;
          const idx = this.chats.findIndex(chat => chat.conversation_id === this.selectedChat.conversation_id);
          if (idx !== -1) {
            this.chats[idx] = {
            ...this.chats[idx],
            lastMessage: this.newMessage,
          };
          }
        } else if (this.selectedChatType === "group") {
          this.selectedChat.lastMessage = this.newMessage;
          const idx = this.groupChats.findIndex(group => group.group_conversation_id === this.selectedChat.group_conversation_id);
          if (idx !== -1) {
            this.groupChats[idx] = {
            ...this.groupChats[idx],
            group_lastMessage: this.newMessage,
          };
          }
        }
        this.scrollToBottom();
        this.fetchChats();
      } catch (error) {
        console.error("Errore nell'invio del messaggio:", error);
      }
      this.newMessage = "";
    },
    async deleteMessage(messageId) {
      if (!messageId) return;
      const token = localStorage.getItem("token");
      try {
        await axios.delete(`${__API_URL__}/messages/${messageId}`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        this.messages = this.messages.filter(msg => msg.id !== messageId);
        this.showOptions = false;
      } catch (error) {
        console.error("Errore nell'eliminazione del messaggio:", error);
      }
      this.fetchChats();
    },
    openMessageMenu(messageId, messageSender) {
      this.selectedMessageId = messageId;
      this.selectedMessageSender = messageSender;
      this.showOptions = true;
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
    },
    async deleteConversation() {
      const token = localStorage.getItem("token");
      if (this.selectedChatType === "private" && this.selectedChat.conversation_id) {
        try {
          await axios.delete(`${__API_URL__}/conversations/${this.selectedChat.conversation_id}`, {
            headers: { Authorization: `Bearer ${token}` },
          });
          this.chats = this.chats.filter(chat => chat.conversation_id !== this.selectedChat.conversation_id);
          this.selectedChat = null;
        } catch (error) {
          console.error("Errore nell'eliminazione della conversazione:", error);
        }
      } else if (this.selectedChatType === "group" && this.selectedChat.group_conversation_id) {
        try {
          await axios.delete(`${__API_URL__}/conversations/${this.selectedChat.group_conversation_id}`, {
            headers: { Authorization: `Bearer ${token}` },
          });
          this.groupChats = this.groupChats.filter(
            group => group.group_conversation_id !== this.selectedChat.group_conversation_id
          );
          this.selectedChat = null;
        } catch (error) {
          console.error("Errore nell'eliminazione della conversazione di gruppo:", error);
        }
      }
    },
    async deleteGroupConversation() {
      const token = localStorage.getItem("token");
      if (!this.selectedChat || !this.selectedChat.group_conversation_id) return
      this.groupErrorMessage = ""; 
      try {
        const groupId = this.selectedChat.group_conversation_id;
        await axios.delete(`${__API_URL__}/groups/${groupId}`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        this.groupChats = this.groupChats.filter(
          group => group.group_conversation_id !== groupId
        );
        this.selectedChat = null;
      } catch (error) {
        if (error.response && error.response.status === 403) {
          this.groupErrorMessage = "Non puoi eliminare il gruppo in quanto non sei admin.";
        } else {
          this.groupErrorMessage = "Si è verificato un errore durante l'eliminazione del gruppo.";
          console.error("Errore nella cancellazione del gruppo:", error);
        }
        this.showGroupErrorModal = true;
      }
    },
    async fetchUserPhoto() {
      if (!this.selectedChat || this.selectedChatType !== "private") return;
      const recipientId = this.selectedChat.recipient_id;
      const token = localStorage.getItem("token");
      try {
        const response = await axios.get(`${__API_URL__}/users/${recipientId}/photo`, {
          headers: { Authorization: `Bearer ${token}` },
          responseType: "blob",
        });
        if (!response.data || response.data.size === 0) {
          this.avatarUrl = defaultAvatar;
          return;
        }
        const base64data = await blobToBase64(response.data);
        this.avatarUrl = base64data;
      } catch (error) {
        console.error("Errore nel recupero della foto del destinatario:", error);
        this.avatarUrl = defaultAvatar;
      }
    },
    async fetchGroupPhoto() {
      if (!this.selectedChat || !this.selectedChat.group_conversation_id || this.selectedChatType !== "group") return;
      try {
        const response = await axios.get(
          `${__API_URL__}/groups/${this.selectedChat.group_conversation_id}/photo`,
          { responseType: "blob" }
        );
        if (response.data.size === 0) {
          this.groupImage = defaultAvatar;
          const idx = this.groupChats.findIndex(
            group => group.group_conversation_id === this.selectedChat.group_conversation_id
          );
          if (idx !== -1) {
            this.groupChats[idx].group_avatarUrl = defaultAvatar;
            this.groupChats = [...this.groupChats];
          }
          return;
        }
        const base64data = await blobToBase64(response.data);
        this.groupImage = base64data;
        const idx = this.groupChats.findIndex(
          group => group.group_conversation_id === this.selectedChat.group_conversation_id
        );
        if (idx !== -1) {
          this.groupChats[idx].group_avatarUrl = base64data;
          this.groupChats = [...this.groupChats]; 
        }
        eventBus.emit("groupPhotoUpdated", {
          groupId: this.selectedChat.group_conversation_id,
          image: base64data,
        });
      } catch (error) {
        console.error("Errore nel recupero della foto del gruppo:", error);
        this.groupImage = defaultAvatar;
      }
    },
    async uploadGroupPhoto(event) {
      const file = event.target.files[0];
      if (!file) return;
      const token = localStorage.getItem("token");
      const formData = new FormData();
      formData.append("photo", file);
      try {
        await axios.put(`${__API_URL__}/groups/${this.selectedChat.group_conversation_id}/photo`, formData, {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "multipart/form-data",
          },
        });
        await this.fetchGroupPhoto();
      } catch (error) {
        console.error("Errore nell'upload della foto del gruppo:", error);
      }
    },
    async uploadProfilePicture(event) {
      const file = event.target.files[0];
      if (!file) return;
      const token = localStorage.getItem("token");
      const formData = new FormData();
      formData.append("photo", file);
      try {
        await axios.put(`${__API_URL__}/users/${token}/photo`, formData, {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "multipart/form-data",
          },
        });
        await this.fetchProfilePhoto();
      } catch (error) {
        console.error("Errore nel caricamento dell'immagine:", error);
      }
    },
    async fetchProfilePhoto() {
      const token = localStorage.getItem("token");
      if (!token) return;
      try {
        const response = await axios.get(`${__API_URL__}/users/${token}/photo`, {
          responseType: "blob",
        });
        if (response.data.size === 0) {
          this.userImage = defaultAvatar;
          return;
        }
        const base64data = await blobToBase64(response.data);
        this.userImage = base64data;
      } catch (error) {
        console.error("Errore nel recupero della foto:", error);
        this.userImage = defaultAvatar;
      }
    },
    // --- Funzioni del profilo utente ---
    editNickname() {
      this.editableNickname = this.nickname;
      this.isEditing = true;
      this.$nextTick(() => {
        this.$refs.editableNickname.focus();
      });
    },
    async saveNickname() {
      if (this.editableNickname.trim() && this.editableNickname !== this.nickname) {
        const token = localStorage.getItem("token");
        try {
          await axios.put(`${__API_URL__}/users/${token}`, { nickname: this.editableNickname.trim() }, {
            headers: {
              Authorization: `Bearer ${token}`,
              "Content-Type": "application/json",
            },
          });
          this.nickname = this.editableNickname.trim();
          localStorage.setItem("nickname", this.nickname);
        } catch (error) {
          console.error("Errore nel salvataggio del nickname:", error);
        }
      }
      this.isEditing = false;
    },
    logout() {
      this.$router.replace("/login");
    },
    // --- Funzioni per la creazione di nuove conversazioni ---
    async fetchUsersForChat() {
      this.showChatOptions = false;
      const token = localStorage.getItem("token");
      try {
        const response = await axios.get(`${__API_URL__}/users`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        if (Array.isArray(response.data)) {
          this.users = response.data.filter(user => user.User_id.toString() !== token);
        } else {
          this.users = [];
        }
        this.showUserList = true;
      } catch (error) {
        console.error("Errore nel recupero degli utenti:", error);
        alert("Errore nel recupero degli utenti.");
      }
    },
    async fetchUsersForGroup() {
      this.showChatOptions = false;
      const token = localStorage.getItem("token");
      try {
        const response = await axios.get(`${__API_URL__}/users`, {
          headers: { Authorization: `Bearer ${token}` },
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
    async startChat(user) {
      const existingChat = this.chats.find(
        chat => chat.recipient_id === user.User_id || chat.sender_id === user.User_id
      );
      if (existingChat) {
        this.selectedChat = existingChat;
        this.selectedChatType = "private";
        this.showChatOptions = false;
        this.showUserList = false;
        this.fetchUserPhoto();
        this.fetchMessages();
        return;
      }
      const token = localStorage.getItem("token");
      try {
        const response = await axios.post(
          `${__API_URL__}/conversations/conversation`,
          { recipient_id: user.User_id },
          { headers: { Authorization: `Bearer ${token}` } }
        );
        const conversationId = response.data.conversation_id;
        const newChat = {
          conversation_id: conversationId,
          sender_id: parseInt(token),
          recipient_id: user.User_id,
          name: user.Nickname,
          // Imposta direttamente l'URL backend per la foto
          avatarUrl: user.Avatar ? `${__API_URL__}/users/${user.User_id}/photo` : defaultAvatar,
          lastMessage: "",
        };
        this.chats.push(newChat);
        this.selectedChat = newChat;
        this.selectedChatType = "private";
        this.showChatOptions = false;
        this.showUserList = false;
        this.fetchUserPhoto();
        await this.fetchChats();
        this.fetchMessages();
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
        const response = await axios.post(`${__API_URL__}/groups`, { group_name: this.groupName }, { headers: { Authorization: `Bearer ${token}` } });
        const groupId = response.data.group_id;
        for (const userId of this.selectedUsers) {
          await axios.post(`${__API_URL__}/groups/${groupId}/users/${userId}?state=add`, { role: "member" }, { headers: { Authorization: `Bearer ${token}` } });
        }
        const newGroup = {
          group_conversation_id: groupId,
          group_name: this.groupName,
          group_avatarUrl: defaultAvatar,
          group_last_message_id: "",
        };
        this.groupChats.push(newGroup);
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
  },
  created() {
    this.fetchChats();
    this.fetchGroupChats();
    this.fetchProfilePhoto();
    this.fetchUserPhoto();
    this.fetchGroupPhoto();
    eventBus.on("conversationDeleted", conversationId => {
      this.chats = this.chats.filter(chat => chat.conversation_id !== conversationId);
      this.groupChats = this.groupChats.filter(group => group.group_conversation_id !== conversationId);
    });
  },
  beforeUnmount() {
    eventBus.off("conversationDeleted");
  },
};
</script>

<style scoped>
/* Container */
.container-fluid {
  height: 100vh;
  width: 100%;
  display: flex;
  padding: 0;
  margin: 0;
  overflow: hidden;
}

/* Pannello sinistro */
.left-panel {
  width: 30%;
  background: #ffffff;
  border-right: 1px solid #ddd;
  display: flex;
  flex-direction: column;
  height: 100vh;
}

/* Pannello destro */
.right-panel {
  width: 70%;
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
  box-sizing: border-box;
}

/* Profilo utente */
.user-profile {
  display: flex;
  align-items: center;
  padding: 15px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 10px;
}
.profile-picture {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 50px;
  height: 50px;
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
  width: 50px;
  height: 50px;
  object-fit: cover;
  border-radius: 50%;
  display: block;
}
.user-name-container {
  margin-left: 15px;
  display: flex;
  align-items: center;
  cursor: pointer;
  min-width: 150px;
  max-width: 200px;
}
.user-name {
  font-size: 18px;
  font-weight: bold;
  color: #333;
}
.user-name-input {
  min-width: 120px;
  font-size: 18px;
  font-weight: bold;
  color: #333;
  border: none;
  background: transparent;
  outline: none;
  text-align: left;
}
.logout-btn {
  margin-left: auto;
  padding: 10px 15px;
  background-color: #069327;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

/* Lista Chat */
.chat-list-container {
  flex-grow: 1;
  overflow-y: auto;
  background-color: #f0f0f0;
}
.search-bar-container {
  display: flex;
  align-items: center;
  padding: 8px 10px;
  background-color: #ffffff;
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
  overflow-y: auto;
}
.chat-item {
  display: flex;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #ddd;
  cursor: pointer;
}
.chat-details {
  display: flex;
  flex-direction: column;
  margin-left: 15px;
}
.chat-name {
  font-weight: bold;
  font-size: 16px;
}
.chat-last-message {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.group-chat-last-message {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.error-message {
  color: red;
  margin-top: 10px;
  font-weight: bold;
}
.chat-window {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #e5ddd5;
}
.chat-header {
  position: fixed;
  top: 0;
  left: 30%;
  width: 70%;
  height: 60px;
  background-color: #2f814e;
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
.chat-header-name {
  margin-left: 15px;
  font-size: 1.2rem;
  font-weight: bold;
}
.delete-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 20px;
  color: white;
  margin-left: auto;
}
.delete-btn:hover {
  color: red;
}
.chat-messages {
  margin-top: 60px;
  flex-grow: 1;
  overflow-y: auto;
  padding: 1rem;
}
.message {
  margin-bottom: 10px;
  display: flex;
}
.message-me {
  justify-content: flex-end;
}
.message-other {
  justify-content: flex-start;
}
.message-content {
  padding: 10px;
  border-radius: 10px;
  background-color: #dcf8c6;
  position: relative;
}
.message-time {
  font-size: 10px;
  color: gray;
  white-space: nowrap;
  margin-top: 5px;
  cursor: pointer;
}
.message-input-container {
  padding: 1rem;
  border-top: 1px solid #ddd;
  background-color: #f0f0f0;
  display: flex;
}
.message-input-container input[type="text"] {
  flex-grow: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 20px;
  outline: none;
}
.message-input-container button {
  margin-left: 10px;
  padding: 10px;
  background-color: #069327;
  color: white;
  border: none;
  border-radius: 50%;
  cursor: pointer;
}

/* Chat vuota */
.empty-chat {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-grow: 1;
  background-color: #e5ddd5;
}

/* Modali */
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
  width: 90%;
  max-width: 400px;
  text-align: center;
}
.modal-content p {
  margin-top: 10px;
  font-size: 16px;
  color: #000000; 
}
.modal-error-title {
  color: #ff0800;
}
.modal-error-box {
  margin: 20px 0;
  padding: 15px;
  border-top: 1px solid #ccc;
  border-bottom: 1px solid #ccc;
  font-size: 16px;
  color: #333;
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
  background-color: grey;
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
.selected {
  background-color: lightgreen;
}
.loading-text {
  text-align: center;
  margin-top: 25%;
}
</style>