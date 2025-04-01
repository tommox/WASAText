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
            <div class="chat-last-time">
              <template v-if="chat.lastMessage && chat.lastMessage !== 'Nessun messaggio' && chat.lastMessageSenderId === parseInt(token)">
                <img
                  class="checkmark-icon"
                  :src="chat.lastMessageIsRead ? doubleCheckmark : checkmark"
                  alt="Spunta"
                />
              </template>
              {{ formatTime(chat.lastMessageTimestamp) }}
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
            <div class="group-chat-last-time">
              <template v-if="group.group_lastMessage && group.group_lastMessage !== 'Nessun messaggio' && group.group_lastMessageSenderId === parseInt(token)">
                <img
                  class="checkmark-icon"
                  :src="group.group_lastMessageIsRead ? doubleCheckmark : checkmark"
                  alt="Spunta"
                />
              </template>
              {{ formatTime(group.group_lastMessageTimestamp) }}
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
          <div class="chat-header-name" v-if="selectedChatType === 'group'" @click="toggleGroupMenu">
            {{ selectedChat.group_name }}
            <div v-if="showGroupMenu" class="modal-overlay" @click="showGroupMenu = false">
              <div class="modal-content" @click.stop>
                <h2>Azioni gruppo</h2>
                <div class="user-list">
                  <div class="user-item" @click="openRenameGroup"> Cambia nome al gruppo</div>
                  <div class="user-item" @click="openManageMembers"> Gestisci membri del gruppo</div>
                </div>
                <button @click="showGroupMenu = false" class="cancel-btn">Chiudi</button>
              </div>
            </div>
          </div>
          <span class="chat-header-name" v-else>
            {{ selectedChat.name }}
          </span>
          <!-- Modale Cambia Nome -->
          <div v-if="showRenameGroupModal" class="modal-overlay">
            <div class="modal-content">
              <h2>Modifica nome gruppo</h2>
              <input v-model="editedGroupName" class="modal-input" placeholder="Nuovo nome" />
              <button class="confirm-btn" @click="renameGroup">Salva</button>
              <button class="cancel-btn" @click="showRenameGroupModal = false">Annulla</button>
            </div>
          </div>

          <!-- Modale Aggiunta Membri -->
          <div v-if="showAddMembersModal" class="modal-overlay">
            <div class="modal-content">
              <h2>Aggiungi membri</h2>
              <input
                v-model="newMemberSearch"
                type="text"
                placeholder="Cerca utente..."
                class="modal-input"
              />
              <div class="user-list">
                <div
                  v-for="user in filteredAvailableUsers"
                  :key="user.User_id"
                  class="user-item"
                  @click="addUserToGroup(user.User_id)"
                >
                  {{ user.Nickname }}
                </div>
              </div>
              <button class="cancel-btn" @click="showAddMembersModal = false">Chiudi</button>
            </div>
          </div>

          <!-- Modale Gestione Membri -->
          <div v-if="showManageMembersModal" class="modal-overlay">
            <div class="modal-content">
              <div class="modal-header-with-button">
                <h2>Gestisci membri del gruppo</h2>
                <button @click="openAddMembersModal" class="new-chat-btn" title="Aggiungi membro">➕</button>
              </div>
              <div v-if="groupMembers.length === 0" class="loading-text">Caricamento membri...</div>
              <div class="user-list">
                <div
                  v-for="user in groupMembers"
                  :key="user.user_id"
                  class="user-item"
                >
                  {{ user.nickname }}
                  <span v-if="user.role === 'admin'" class="badge">Admin</span>
                  <button
                    v-if="user.role !== 'admin'"
                    class="small-btn"
                    @click="promoteToAdmin(user.user_id)">
                    Promuovi
                  </button>
                  <button class="small-btn" @click="removeFromGroup(user.user_id)">Rimuovi</button>
                </div>
              </div>
              <button class="cancel-btn" @click="showManageMembersModal = false">Chiudi</button>
            </div>
          </div>
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
              <template v-if="message.isReply">
                <div class="reply-indicator">
                  <strong>Risposta a:</strong> <span>{{ message.replyMessageText }}</span>
                </div>
              </template>
              <template v-if="selectedChatType === 'group'">
                <strong>{{ message.sender === 'me' ? 'Tu' : userMap[message.rawSenderId] || 'Utente' }}</strong><br />
              </template>
              <template v-if="message.imageData">
                <div class="message-image">
                  <img :src="message.imageData" alt="Messaggio immagine" class="message-img" />
                </div>
              </template>
              
              <template v-else>
                <span>{{ message.text }}</span>
              </template>
              <div class="message-footer">
                <span class="message-time" @click="openMessageMenu(message.id, message.sender)">
                  {{ formatTime(message.timestamp) }}
                </span>
                <span v-if="message.sender === 'me'">
                  <img
                    :src="message.isRead ? doubleCheckmark : checkmark"
                    alt="Spunta"
                    class="checkmark-icon"
                  />
                </span>
              </div>
              <div class="message-reactions">
                <span 
                  v-for="reaction in message.reactions" 
                  :key="reaction.user_id" 
                  class="reaction"
                  @click="removeReaction(message.id, reaction)"
                  :title="userMap[reaction.user_id] || 'Utente Sconosciuto'">
                  {{ reaction.emoji }}
                </span>
              </div>
              <div v-if="showOptions && selectedMessageId === message.id" class="modal-overlay">
                <div class="modal-content">
                  <h2>Seleziona un'opzione</h2>
                  <div class="user-list">
                    <div
                      v-if="selectedMessageSender === 'me'"
                      class="user-item"
                      @click="deleteMessage(selectedMessageId)"
                    >
                      Elimina
                    </div>
                    <div class="user-item" @click="forwardMessage(selectedMessageId)">
                      Inoltra
                    </div>
                    <div class="user-item" @click="replyToMessage(selectedMessageId)">
                      Rispondi
                    </div>
                    <div class="emoji-container">
                      <span class="emoji" @click="reactToMessage(selectedMessageId, '👍')">👍</span>
                      <span class="emoji" @click="reactToMessage(selectedMessageId, '❤️')">❤️</span>
                      <span class="emoji" @click="reactToMessage(selectedMessageId, '😂')">😂</span>
                      <span class="emoji" @click="reactToMessage(selectedMessageId, '😢')">😢</span>
                      <span class="emoji" @click="reactToMessage(selectedMessageId, '🔥')">🔥</span>
                    </div>
                  </div>
                  <button @click="showOptions = false" class="cancel-btn">Chiudi</button>
                </div>
              </div>
              <div v-if="forwardMode" class="modal-overlay">
                <div class="modal-content">
                  <h2>Inoltra a...</h2>
                  <div class="user-list">
                    <div
                      v-for="user in users"
                      :key="'fwd-user-' + user.User_id"
                      class="user-item"
                      @click="confirmForward('private', user.User_id)"
                    >
                      {{ user.Nickname }}
                    </div>
                    <div
                      v-for="group in groupChats"
                      :key="'fwd-group-' + group.group_conversation_id"
                      class="user-item"
                      @click="confirmForward('group', group.group_conversation_id)"
                    >
                      {{ group.group_name }}
                    </div>
                  </div>
                  <button @click="forwardMode = false" class="cancel-btn">Annulla</button>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div v-if="isReplying" class="reply-indicator">
          Risposta a: <strong>{{ replyMessageText }}</strong>
        </div>
        <!-- Input per inviare messaggio -->
        <div class="message-input-container">
          <input
            v-model="newMessage"
            type="text"
            placeholder="Scrivi un messaggio..."
            @keyup.enter="sendCurrentMessage"
          />
          <label class="image-upload-btn">
            📷
            <input type="file" accept="image/*" @change="sendPhotoMessage" style="display: none" ref="imageInput" />
          </label>
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
import checkmark from "@/assets/images/checkmark.png";
import doubleCheckmark from "@/assets/images/doublecheckmark.png";


function blobToBase64(blob) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onloadend = () => resolve(reader.result);
    reader.onerror = reject;
    reader.readAsDataURL(blob);
  });
}
export default {
  name: "HomeView",
  data() {
    return {
      // Profilo utente
      nickname: localStorage.getItem("nickname") || "Utente",
      userImage: defaultAvatar, 
      editableNickname: "",
      isEditing: false,
      showGroupErrorModal: false,
      groupErrorMessage: "",
      // Lista chat
      search: "",
      chats: [],
      groupChats: [],
      userMap: {},
      checkmark,
      doubleCheckmark,
      // Selezione chat
      selectedChat: null,
      selectedChatType: "",
      // Modali per creare chat/gruppo
      showChatOptions: false,
      showUserList: false,
      showGroupUserList: false,
      userSearch: "",
      selectedImage: null,
      users: [],
      groupName: "",
      selectedUsers: [],
      showGroupMenu: false,
      showRenameGroupModal: false,
      showManageMembersModal: false,
      editedGroupName: "",
      groupMembers: [],
      newMemberSearch: "",
      availableUsers: [],
      showAddMembersModal: false,
      // Finestra chat
      messages: [],
      newMessage: "",
      forwardMode: false,
      forwardMessageId: null,
      forwardMessageType: "",
      loading: false,
      avatarUrl: defaultAvatar, 
      groupImage: defaultAvatar,
      isReplying: false, 
      replyMessageText: "", 
      replyMessageId: null, 
      // Opzioni messaggio
      showOptions: false,
      selectedMessageId: null,
      selectedMessageSender: null,
    };
  },
  computed: {
    token() {
      return parseInt(localStorage.getItem("token"));
    },
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
    filteredAvailableUsers() {
      return this.availableUsers.filter(user =>
        user.Nickname.toLowerCase().includes(this.newMemberSearch.toLowerCase())
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
    async fetchChats() {
      const token = localStorage.getItem("token");
      if (!token) return;
      try {
        const response = await this.$axios.get(`/conversations`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        const userResponse = await this.$axios.get(`/users`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        const allUsers = userResponse.data;
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
                const photoResponse = await this.$axios.get(`/users/${recipientId}/photo`, {
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
            let lastMessageTimestamp = null;
            let lastMessageSenderId = null; 
            let lastMessageIsRead = false;  
            if (chat.last_message_id) {
              try {
                const msgResponse = await this.$axios.get(
                  `/messages/${chat.last_message_id}?type=private`,
                  { headers: { Authorization: `Bearer ${token}` } }
                );

                if (msgResponse.data.message_content && msgResponse.data.message_content.trim() !== "") {
                  lastMessage = msgResponse.data.message_content;
                } else {
                  lastMessage = "📷 Foto";
                }
                lastMessageTimestamp = msgResponse.data.timestamp;
                lastMessageSenderId = msgResponse.data.sender_id;
                lastMessageIsRead = msgResponse.data.isRead === true;
              } catch (error) {
                console.error("Errore nel recupero dell'ultimo messaggio", error);
              }
            }
            return {
              ...chat,
              recipient_id: recipientId,
              name: recipient ? recipient.Nickname : "Utente Sconosciuto",
              avatarUrl,
              lastMessage,
              lastMessageTimestamp,
              lastMessageSenderId,
              lastMessageIsRead, 
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
        const response = await this.$axios.get(`/conversations`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        let groupChats = [];
        if (Array.isArray(response.data.group_conversations)) {
          groupChats = await Promise.all(response.data.group_conversations.map(async group => {
            const groupConversationId = group.group_conversation_id;
            let lastMessage = "Nessun messaggio";
            let lastMessageTimestamp = null;
            let avatarUrl = defaultAvatar;
            let lastMessageSenderId = null;
            let lastMessageIsRead = false;
            if (group.last_message_id) {
              try {
                const msgResponse = await this.$axios.get(
                  `/messages/${group.last_message_id}?type=group`,
                  { headers: { Authorization: `Bearer ${token}` } }
                );
                if (msgResponse.data.message_content && msgResponse.data.message_content.trim() !== "") {
                  lastMessage = msgResponse.data.message_content;
                } else {
                  lastMessage = "📷 Foto";
                }
                lastMessageTimestamp = msgResponse.data.timestamp;
                lastMessageSenderId = msgResponse.data.sender_id;
                lastMessageIsRead = msgResponse.data.isRead === true;
              } catch (error) {
                console.error("Errore nel recupero dell'ultimo messaggio per il gruppo", error);
              }
            }
            try {
              const photoResponse = await this.$axios.get(
                `/groups/${groupConversationId}/photo`,
                {
                  headers: { Authorization: `Bearer ${token}` },
                  responseType: "blob"
                }
              );
              if (photoResponse.data && photoResponse.data.size > 0) {
                avatarUrl = await blobToBase64(photoResponse.data);
              }
            } catch (e) {
              console.warn("Errore nel recupero foto gruppo:", e);
              avatarUrl = defaultAvatar;
            }
            return {
              group_conversation_id: groupConversationId,
              group_name: group.group_name || "Gruppo Sconosciuto",
              group_avatarUrl: avatarUrl,
              group_lastMessage: lastMessage,
              group_lastMessageTimestamp: lastMessageTimestamp,
              group_lastMessageSenderId: lastMessageSenderId,
              group_lastMessageIsRead: lastMessageIsRead,
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
        const response = await this.$axios.get(`/groups/${groupId}/users`, {
          headers: { Authorization: `Bearer ${token}` }
        });

        const members = response.data;
        this.groupMembers = members;
        this.userMap = Object.fromEntries(members.map(user => [user.user_id, user.nickname]));
      } catch (error) {
        console.error("Errore nel recupero dei membri del gruppo:", error);
      }
    },
    toggleGroupMenu() {
      this.showGroupMenu = !this.showGroupMenu;
    },
    async renameGroup() {
      const token = localStorage.getItem("token");
      try {
        await this.$axios.patch(`/groups/${this.selectedChat.group_conversation_id}`, {
          group_name: this.editedGroupName,
        }, {
          headers: { Authorization: `Bearer ${token}` },
        });

        this.selectedChat.group_name = this.editedGroupName;
        const idx = this.groupChats.findIndex(group => group.group_conversation_id === this.selectedChat.group_conversation_id);
        if (idx !== -1) {
          this.groupChats[idx].group_name = this.editedGroupName;
        }

        this.showRenameGroupModal = false;
        this.showGroupMenu = false;
      } catch (err) {
        console.error("Errore nel cambiare nome al gruppo:", err);
        alert("Errore nella modifica del nome.");
      }
    },
    async promoteToAdmin(userId) {
      const token = localStorage.getItem("token");
      try {
        await this.$axios.post(`/groups/${this.selectedChat.group_conversation_id}/users/${userId}?state=promote`, null, {
          headers: { Authorization: `Bearer ${token}` },
        });
        this.fetchGroupMembers(this.selectedChat.group_conversation_id);
      } catch (err) {
        console.error("Errore nella promozione:", err);
      }
    },

    async removeFromGroup(userId) {
      const token = localStorage.getItem("token");
      try {
        await this.$axios.post(`/groups/${this.selectedChat.group_conversation_id}/users/${userId}?state=remove`, null, {
          headers: { Authorization: `Bearer ${token}` },
        });
        this.fetchGroupMembers(this.selectedChat.group_conversation_id);
        eventBus.emit("userRemovedFromGroup", {
          groupId: this.selectedChat.group_conversation_id,
          userId: userId,
        });
      } catch (err) {
        console.error("Errore nella rimozione dal gruppo:", err);
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
        this.editedGroupName = chat.group_name;
      }
    },
    async fetchMessages() {
      if (!this.selectedChat) return;
      this.loading = true;
      const token = localStorage.getItem("token");
      try {
        let response;
        if (this.selectedChatType === "private" && this.selectedChat.conversation_id) {
          response = await this.$axios.get(
            `/conversations/${this.selectedChat.conversation_id}?type=private`,
            { headers: { Authorization: `Bearer ${token}` } }
          );
          if (Array.isArray(response.data)) {
            this.messages = await Promise.all(response.data.map(async (msg) => {
              let replyMessageText = null;
              if (msg.isReply) {
                try {
                  const replyResponse = await this.$axios.get(`/messages/${msg.isReply}?type=private`, {
                    headers: { Authorization: `Bearer ${token}` }
                  });
                  replyMessageText = replyResponse.data.message_content;
                } catch (error) {
                  console.error("Errore nel recupero del messaggio di risposta", error);
                }
              }
              return {
                id: msg.message_id,
                text: msg.image_data ? "" : (msg.message_content || ""),
                imageData: msg.image_data ? `data:image/jpeg;base64,${msg.image_data}` : null,
                sender: msg.sender_id === Number(token) ? "me" : "other",
                rawSenderId: msg.sender_id,
                timestamp: new Date(msg.timestamp),
                isRead: msg.isRead || false,
                isReply: msg.isReply,
                replyMessageText: replyMessageText,
                reactions: await this.fetchReactionsForMessage(msg.message_id),
              };
            }));
            console.log("fetch_priv",this.messages);
          } else {
            this.messages = [];
          }
        } else if (this.selectedChatType === "group" && this.selectedChat.group_conversation_id) {
          response = await this.$axios.get(
            `/conversations/${this.selectedChat.group_conversation_id}?type=group`,
            { headers: { Authorization: `Bearer ${token}` } }
          );
          if (Array.isArray(response.data)) {
            this.messages = await Promise.all(response.data.map(async (msg) => {
              let replyMessageText = null;
              if (msg.isReply) {
                try {
                  const replyResponse = await this.$axios.get(`/messages/${msg.isReply}?type=group`, {
                    headers: { Authorization: `Bearer ${token}` }
                  });
                  replyMessageText = replyResponse.data.message_content;
                } catch (error) {
                  console.error("Errore nel recupero del messaggio di risposta", error);
                }
              }
              return {
                id: msg.message_id,
                text: msg.image_data ? "" : (msg.message_content || ""),
                imageData: msg.image_data ? `data:image/jpeg;base64,${msg.image_data}` : null,
                sender: msg.sender_id === Number(token) ? "me" : "other",
                rawSenderId: msg.sender_id,
                timestamp: new Date(msg.timestamp),
                isRead: msg.isRead || false,
                isReply: msg.isReply,
                replyMessageText: replyMessageText,
                reactions: await this.fetchReactionsForMessage(msg.message_id),
              };
            }));
            console.log("fetch_group",this.messages);
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
    async fetchReactionsForMessage(messageId) {
      const token = localStorage.getItem("token");
      const isGroup = this.selectedChatType === 'group';
      try {
        const response = await this.$axios.get(
          `/messages/${messageId}/reactions?isGroup=${isGroup}`,
          {
            headers: { Authorization: `Bearer ${token}` }
          }
        );
        const reactions = response.data || [];
        for (const reaction of reactions) {
          if (!this.userMap[reaction.user_id]) {
            this.userMap[reaction.user_id] = await this.getUserNickname(reaction.user_id);
          }
        }
        return reactions;
      } catch (error) {
        console.error("Errore nel recupero delle reazioni:", error);
        return [];
      }
    },
    async sendCurrentMessage() {
      if ((this.newMessage.trim() === "" && !this.selectedImage) || !this.selectedChat) return;
      const token = localStorage.getItem("token");
      if (!token) {
        alert("Sessione scaduta. Effettua nuovamente il login.");
        this.$router.push("/login");
        return;
      }
      try {
        let response;
        if (this.selectedChatType === "private") {
          if (this.selectedImage) {
            const formData = new FormData();
            formData.append("photo", this.selectedImage);
            formData.append("conversation_id", this.selectedChat.conversation_id);
            if (this.replyMessageId) {
              formData.append("isReply", this.replyMessageId);
            }
            response = await this.$axios.post(`/messages`, formData, {
              headers: {
                Authorization: `Bearer ${token}`,
                "Content-Type": "multipart/form-data",
              },
            });
            this.selectedImage = null; 
          } else {
            response = await this.$axios.post(
              `/messages`,
              {
                conversation_id: this.selectedChat.conversation_id,
                message_content: this.newMessage,
                isReply: this.replyMessageId,
              },
              { headers: { Authorization: `Bearer ${token}` } }
            );
          }
        } else if (this.selectedChatType === "group") {
          if (this.selectedImage) {
            const formData = new FormData();
            formData.append("photo", this.selectedImage);
            formData.append("group_conversation_id", this.selectedChat.group_conversation_id);
            if (this.replyMessageId) {
              formData.append("isReply", this.replyMessageId);
            }
            response = await this.$axios.post(`/groups/${this.selectedChat.group_conversation_id}/messages`, formData, {
              headers: {
                Authorization: `Bearer ${token}`,
                "Content-Type": "multipart/form-data",
              },
            });
            this.selectedImage = null; 
          } else {
            response = await this.$axios.post(
              `/groups/${this.selectedChat.group_conversation_id}/messages`,
              { message_content: this.newMessage, isReply: this.replyMessageId },
              { headers: { Authorization: `Bearer ${token}` } }
            );
          }
        }
        const getMessageResponse = await this.$axios.get(`/messages/${response.data.message_id}`, {
          headers: { Authorization: `Bearer ${token}` },
          params: { type: this.selectedChatType }
        });
        const data = getMessageResponse.data;
        this.messages.push({
          id: data.message_id,
          text: data.message_content || "",
          sender: "me",
          timestamp: new Date(data.timestamp),
          imageData: data.image_data ? `data:image/jpeg;base64,${data.image_data}` : null,
          isRead: data.isRead || false,
          isReply: data.isReply || null,
          replyMessageText: this.replyMessageText || null,
          reactions: [],
        });
        console.log("sendMess",this.messages);
        if (this.selectedChatType === "private") {
          this.selectedChat.lastMessage = this.newMessage || "📷 Foto";
          const idx = this.chats.findIndex(chat => chat.conversation_id === this.selectedChat.conversation_id);
          if (idx !== -1) {
            this.chats[idx] = {
              ...this.chats[idx],
              lastMessage: this.newMessage || "📷 Foto",
            };
          }
        } else if (this.selectedChatType === "group") {
          this.selectedChat.lastMessage = this.newMessage || "📷 Foto";
          const idx = this.groupChats.findIndex(group => group.group_conversation_id === this.selectedChat.group_conversation_id);
          if (idx !== -1) {
            this.groupChats[idx] = {
              ...this.groupChats[idx],
              group_lastMessage: this.newMessage || "📷 Foto",
            };
          }
        }
        this.scrollToBottom();
        this.fetchChats();
        this.fetchGroupChats();
      } catch (error) {
        console.error("Errore nell'invio del messaggio:", error);
      }
      this.newMessage = "";
      this.isReplying = false;
      this.replyMessageText = "";
      this.replyMessageId = null;
    },

    sendPhotoMessage(event) {
      const file = event.target.files[0];
      if (file) {
        this.selectedImage = file;
        this.sendCurrentMessage();
      }
    },
    async deleteMessage(messageId) {
      if (!messageId) return;
      const token = localStorage.getItem("token");
      try {
        const type = this.selectedChatType;
        await this.$axios.delete(`/messages/${messageId}?type=${type}`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        this.messages = this.messages.filter(msg => msg.id !== messageId);
        this.showOptions = false;
      } catch (error) {
        console.error("Errore nell'eliminazione del messaggio:", error);
      }
      this.fetchChats();
      this.fetchGroupChats();
    },
    replyToMessage(messageId) {
      this.selectedMessageId = messageId;
      const message = this.messages.find(msg => msg.id === messageId);
      this.replyMessageText = message.text;
      this.isReplying = true;
      this.replyMessageId = messageId;
      this.showOptions = false; 
    },
    openMessageMenu(messageId, messageSender) {
      this.selectedMessageId = messageId;
      this.selectedMessageSender = messageSender;
      this.showOptions = true;
    },
    openAddMembersModal() {
      this.showManageMembersModal = false; 
      this.showAddMembersModal = true; 
    },
    openRenameGroup() {
      this.showRenameGroupModal = true;
      this.showGroupMenu = false;
    },
    openManageMembers() {
      this.showManageMembersModal = true;
      this.showGroupMenu = false;
      this.fetchAvailableUsers();
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
          await this.$axios.delete(`/conversations/${this.selectedChat.conversation_id}`, {
            headers: { Authorization: `Bearer ${token}` },
          });
          this.chats = this.chats.filter(chat => chat.conversation_id !== this.selectedChat.conversation_id);
          this.selectedChat = null;
        } catch (error) {
          console.error("Errore nell'eliminazione della conversazione:", error);
        }
      } else if (this.selectedChatType === "group" && this.selectedChat.group_conversation_id) {
        try {
          await this.$axios.delete(`/conversations/${this.selectedChat.group_conversation_id}`, {
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
        await this.$axios.delete(`/groups/${groupId}`, {
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
        const response = await this.$axios.get(`/users/${recipientId}/photo`, {
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
        const response = await this.$axios.get(
          `/groups/${this.selectedChat.group_conversation_id}/photo`,
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
        await this.$axios.put(`/groups/${this.selectedChat.group_conversation_id}/photo`, formData, {
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
        await this.$axios.put(`/users/${token}/photo`, formData, {
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
        const response = await this.$axios.get(`/users/${token}/photo`, {
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
          await this.$axios.put(`/users/${token}`, { nickname: this.editableNickname.trim() }, {
            headers: {
              Authorization: `Bearer ${token}`,
              "Content-Type": "application/json",
            },
          });
          this.nickname = this.editableNickname.trim();
          localStorage.setItem("nickname", this.nickname);
        } catch (error) {
          if (error.response.status === 409) {
            alert("Questo nickname è già stato utilizzato. Scegli un altro nome.");
          } else {
            console.error("Errore nel salvataggio del nickname:", error);
          }
        }
      }
      this.isEditing = false;
    },
    logout() {
      this.$router.replace("/login");
    },
    async fetchUsersForChat() {
      this.showChatOptions = false;
      const token = localStorage.getItem("token");
      try {
        const response = await this.$axios.get(`/users`, {
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
        const response = await this.$axios.get(`/users`, {
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
        const response = await this.$axios.post(
          `/conversations/conversation`,
          { recipient_id: user.User_id },
          { headers: { Authorization: `Bearer ${token}` } }
        );
        const conversationId = response.data.conversation_id;
        const newChat = {
          conversation_id: conversationId,
          sender_id: parseInt(token),
          recipient_id: user.User_id,
          name: user.Nickname,
          avatarUrl: user.Avatar ? `/users/${user.User_id}/photo` : defaultAvatar,
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
    async forwardMessage(messageId) {
      this.forwardMessageId = messageId;
      this.forwardMessageType = this.selectedChatType;
      this.showOptions = false;
      this.forwardMode = true;
      const token = localStorage.getItem("token");
      try {
        const response = await this.$axios.get(`/users`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        this.users = response.data.filter(user => user.User_id.toString() !== token);
      } catch (error) {
        console.error("Errore nel recupero degli utenti per il forward:", error);
      }
    },
    async confirmForward(type, destinationId) {
      const token = localStorage.getItem("token");
      let payload = {};
      let finalConversationId = destinationId;
      try {
        if (type === "private") {
          const existingChat = this.chats.find(chat =>
            chat.recipient_id === destinationId || chat.sender_id === destinationId
          );
          if (existingChat) {
            finalConversationId = existingChat.conversation_id;
            payload = { conversation_id: finalConversationId };
          } else {
            const createResponse = await this.$axios.post(
              `/conversations/conversation`,
              { recipient_id: destinationId },
              { headers: { Authorization: `Bearer ${token}` } }
            );
            finalConversationId = createResponse.data.conversation_id;
            payload = { conversation_id: finalConversationId };
          }
        } else if (type === "group") {
          payload = { group_id: destinationId };
        }
        await this.$axios.post(
          `/messages/${this.forwardMessageId}/forwards?type=${this.forwardMessageType}`,
          payload,
          {
            headers: {
              Authorization: `Bearer ${token}`,
              "Content-Type": "application/json",
            },
          }
        );
        eventBus.emit("messageForwarded", {
              type,
              destinationId,
            });
        await this.fetchChats();
        await this.fetchGroupChats();
        this.forwardMode = false;
        this.forwardMessageId = null;
        this.forwardMessageType = "";
      } catch (error) {
        console.error("Errore durante l'inoltro del messaggio:", error);
        alert("Errore durante l'inoltro del messaggio.");
        this.forwardMode = false;
      }
    },
    async handleMessageForwarded({type, destinationId}) {
      const isActiveChat =
        (type === "private" &&
          this.selectedChatType === "private" &&
          this.selectedChat?.conversation_id === destinationId) ||
        (type === "group" &&
          this.selectedChatType === "group" &&
          this.selectedChat?.group_conversation_id === destinationId);

      if (isActiveChat) {
        await this.fetchMessages();
        this.scrollToBottom();
      }
      await this.fetchChats();
      await this.fetchGroupChats();
    },
    handleUserRemovedFromGroup({ groupId, userId }) {
      if (
        this.selectedChatType === "group" &&
        this.selectedChat &&
        this.selectedChat.group_conversation_id === groupId
      ) {
        this.messages = this.messages.filter(msg => msg.rawSenderId !== userId);
      }
    },
    async getUserNickname(userId) {
      const token = localStorage.getItem("token");
      try {
        const userResponse = await this.$axios.get(`/users`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        const user = userResponse.data.find(user => user.User_id === userId);
        if (user) {
          this.userMap[userId] = user.Nickname;
          return user.Nickname;
        } else {
          return "Utente Sconosciuto";
        }
      } catch (error) {
        console.error("Errore nel recupero degli utenti:", error);
        return "Utente Sconosciuto";
      }
    },
    async reactToMessage(messageId, reaction) {
      const token = localStorage.getItem("token");
      try {
        const body = {
          emoji: reaction,
          add: true,  
          isGroup: this.selectedChatType === "group",
        };
        await this.$axios.post(`/messages/${messageId}/reactions`, body, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        const message = this.messages.find(msg => msg.id === messageId);
        if (!message.reactions) {
          message.reactions = [];
        }
        message.reactions.push({reaction});
        this.fetchMessages();  
        this.showOptions = false;
      } catch (error) {
        console.error("Errore nell'aggiungere la reazione:", error);
      }
    },
    async removeReaction(messageId, reaction) {
      const token = localStorage.getItem("token");
      try {
        const body = {
          emoji: reaction.emoji,  
          add: false,  
          isGroup: this.selectedChatType === "group", 
        };
        await this.$axios.post(`/messages/${messageId}/reactions`, body, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        const message = this.messages.find(msg => msg.id === messageId);
        if (message) {
          message.reactions = message.reactions.filter(r => r.reaction !== reaction.reaction);
        }
        this.fetchMessages(); 
      } catch (error) {
        console.error("Errore nella rimozione della reazione:", error);
      }
    },
    async createGroup() {
      if (!this.canCreateGroup) return;
      try {
        const token = localStorage.getItem("token");
        const response = await this.$axios.post(`/groups`, { group_name: this.groupName }, { headers: { Authorization: `Bearer ${token}` } });
        const groupId = response.data.group_id;
        for (const userId of this.selectedUsers) {
          await this.$axios.post(`/groups/${groupId}/users/${userId}?state=add`, { role: "member" }, { headers: { Authorization: `Bearer ${token}` } });
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
    async fetchAvailableUsers() {
    const token = localStorage.getItem("token");
    try {
      const response = await this.$axios.get(`/users`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      const allUsers = response.data;
      const currentMemberIds = this.groupMembers.map(u => u.user_id);
      this.availableUsers = allUsers.filter(user => !currentMemberIds.includes(user.User_id));
    } catch (err) {
      console.error("Errore nel recuperare utenti disponibili:", err);
    }
  },
  async addUserToGroup(userId) {
    const token = localStorage.getItem("token");
    try {
      await this.$axios.post(
      `/groups/${this.selectedChat.group_conversation_id}/users/${userId}?state=add`,
      { role: "member" }, 
      {
        headers: { Authorization: `Bearer ${token}` }
      }
    );
      await this.fetchGroupMembers(this.selectedChat.group_conversation_id);
      await this.fetchAvailableUsers();
      this.showAddMembersModal = false;
      this.showManageMembersModal = false;
      this.showGroupMenu = false;
    } catch (err) {
      console.error("Errore nell'aggiunta del membro:", err);
    }
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
    eventBus.on("messageForwarded", this.handleMessageForwarded);
    eventBus.on("userRemovedFromGroup", this.handleUserRemovedFromGroup);
  },
  beforeUnmount() {
    eventBus.off("conversationDeleted");
    eventBus.off("messageForwarded", this.handleMessageForwarded);
    eventBus.off("userRemovedFromGroup", this.handleUserRemovedFromGroup);
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
.chat-last-time, .group-chat-last-time {
  font-size: 12px; 
  color: #888;
  margin-left: auto; 
  text-align: right; 
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
  cursor: pointer;
  position: relative;
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
.emoji-container {
  display: flex;
  justify-content: space-around;
  margin-top: 10px;
}
.emoji {
  font-size: 20px;
  cursor: pointer;
  margin: 0 5px;
}
.message-reactions {
  display: flex;
  justify-content: flex-start;
  margin-top: 5px;
}
.reaction {
  font-size: 24px; 
  margin-right: 10px; 
  cursor: pointer;
}
.reaction:hover {
  opacity: 0.7; 
}
.image-upload-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #ffffff00;
  color: white;
  font-size: 20px;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  margin-left: 10px;
  cursor: pointer;
}
.hidden-image-input {
  display: none;
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
  padding: 10px 15px;
  background-color: #069327;
  color: white;
  border: none;
  border-radius: 50%;
  cursor: pointer;
}
.message-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 5px;
  margin-top: 5px;
  font-size: 10px;
  color: gray;
}
.checkmark-icon {
  width: 20px;
  height: 20px;
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
.modal-header-with-button {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.modal-content h2 {
  color: #000; 
  font-size: 20px; 
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
  color: #000;
  background: #fff;
  max-height: 200px;
  overflow-y: auto;
  border-top: 1px solid #ddd;
  margin-top: 10px;
}
.user-item {
  color: #000;
  font-size: 14px;
  font-weight: normal !important;
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
.group-dropdown {
  position: absolute;
  background: white;
  border: 1px solid #ccc;
  border-radius: 8px;
  top: 60px;
  left: 10px;
  z-index: 999;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}
.dropdown-item {
  padding: 10px 15px;
  cursor: pointer;
  border-bottom: 1px solid #eee;
}
.dropdown-item:last-child {
  border-bottom: none;
}
.dropdown-item:hover {
  background-color: #f5f5f5;
}
.reply-indicator {
  font-size: 14px;
  color: #777;
  margin-bottom: 8px;
  padding: 5px;
  background-color: #f9f9f9;
  border-radius: 5px;
}
.badge {
  background-color: #069327;
  color: white;
  font-size: 12px;
  padding: 2px 5px;
  border-radius: 5px;
  margin-left: 10px;
}
.small-btn {
  background: none;
  border: none;
  color: #069327;
  cursor: pointer;
  margin-left: 10px;
}
.small-btn:hover {
  text-decoration: underline;
}
</style>