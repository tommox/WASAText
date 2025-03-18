<template>
  <div class="container-fluid h-100 w-100 p-0">
    <div class="left-panel">
      <div class="user-profile">
        <label class="profile-picture">
          <img :src="userImage || defaultAvatar" alt="Foto Profilo" class="profile-img" />
          <input type="file" @change="uploadProfilePicture" accept="image/*" class="upload-input" />
        </label>
        <div class="user-name-container" @click="editNickname">
          <span v-if="!isEditing" class="user-name">{{ nickname }}</span>
          <input v-else ref="editableNickname" v-model="editableNickname" class="user-name-input" @blur="saveNickname" @keyup.enter="saveNickname" />
        </div>
        <button @click="logout" class="logout-btn">Logout</button>
      </div>
      <ChatList @chatSelected="openChat" />
    </div>
    <div class="right-panel">
      <template v-if="selectedChat">
        <ChatWindow 
          :chat="selectedChat"
          :type="selectedChatType"
          @conversationDeleted="handleConversationDeleted"
          @closeChat="selectedChat = null"/>
      </template>
      <template v-else>
        <EmptyChat />
      </template>
    </div>
  </div>
</template>
  
  <script>
  import axios from "axios";
  import ChatList from "@/components/ChatList.vue";
  import ChatWindow from "@/components/ChatWindow.vue";
  import EmptyChat from "@/components/EmptyChat.vue";
  import defaultAvatar from "@/assets/images/user.png";
  
  export default {
    components: { ChatList, ChatWindow, EmptyChat },
    data() {
      return {
        selectedChat: null,
        selectedChatType: "",
        nickname: localStorage.getItem("nickname") || "Utente",
        userImage: localStorage.getItem("profileImage") || defaultAvatar,
        editableNickname: "",
        isEditing: false,
      };
    },
    methods: {

      openChat(chat) {
      if (chat.conversation_id) {
        this.selectedChatType = "private"; 
      } else if (chat.group_conversation_id) {
        this.selectedChatType = "group"; 
      }
      this.selectedChat = chat;
    },

      handleConversationDeleted(conversationId) {
        if (this.selectedChat?.conversation_id === conversationId) {
            this.selectedChat = null;
        }
        this.$emit("conversationDeleted", conversationId);
      },

      logout() {
        this.$router.replace("/login");
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
          await axios.put(`${__API_URL__}/users/${token}`, {
            nickname: this.editableNickname.trim()
          }, {
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
          this.fetchUserPhoto();
        } catch (error) {
          console.error("Errore nel caricamento dell'immagine:", error);
        }
      },
      async fetchUserPhoto() {
        const token = localStorage.getItem("token");
        if (!token) return;
        try {
          const response = await axios.get(`${__API_URL__}/users/${token}/photo`, {
            responseType: "blob",
          });
          if (response.data.size === 0) {
            this.userImage = defaultAvatar;
            localStorage.setItem("profileImage", defaultAvatar);
            return;
          }
  
          const imageUrl = URL.createObjectURL(response.data);
          this.userImage = imageUrl;
          localStorage.setItem("profileImage", imageUrl);
          this.$forceUpdate();
        } catch (error) {
          console.error("Errore nel recupero della foto:", error);
          this.userImage = defaultAvatar;
          localStorage.setItem("profileImage", defaultAvatar);
        }
      },
    },
    created() {
      this.fetchUserPhoto();
    },
  };
  </script>
  
  <style scoped>
  .container-fluid {
    height: 100%;
    width: 100%;
    display: flex;
    padding: 0;
    margin: 0;
	  overflow: hidden;
  }
  
  .left-panel {
    width: 30%;
    background: #ffffff;
    border-right: 1px solid #ddd;
    display: flex;
    flex-direction: column; 
    height: 100vh;
  }

  .chat-list-container {
    flex-grow: 1; 
    overflow-y: auto;
    background-color: #f0f0f0;
  }
  
  .right-panel {
    width: 70%;
    display: flex;
    flex-direction: column;
    height: 100vh;
    overflow: hidden;
    box-sizing: border-box;
    padding: 0;
  }
  
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
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  .user-name {
    font-size: 18px;
    font-weight: bold;
    color: #333;
  }
  
  .logout-btn {
    background-color: #069327;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    padding: 10px 15px;
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
    width: auto;
    font-size: 18px;
    font-weight: bold;
    color: #333;
    border: none;
    background: transparent;
    outline: none;
    text-align: left;
  }
  
  </style>
  