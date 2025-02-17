<template>
    <div class="container-fluid h-100 w-100 p-0">
      <div class="left-panel">
        <div class="user-profile">
          <label class="profile-picture">
            <img :src="userImage || defaultAvatar" alt="Foto Profilo" class="profile-img" />
            <input type="file" @change="uploadProfilePicture" accept="image/*" class="upload-input" />
          </label>
          <span class="user-name">{{ nickname }}</span>
          <button @click="logout" class="logout-btn">Logout</button>
        </div>
        <ChatList @chatSelected="openChat" />
      </div>
      <div class="right-panel">
        <template v-if="selectedChat">
          <ChatWindow :chat="selectedChat" @closeChat="selectedChat = null" />
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
        nickname: localStorage.getItem("nickname") || "Utente",
        userImage: localStorage.getItem("profileImage") || defaultAvatar,
      };
    },
    methods: {
      openChat(chat) {
        this.selectedChat = chat;
      },
      logout() {
        this.$router.replace("/login");
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
	position: relative;
	width: 50px;
	height: 50px;
	border-radius: 50%;
	overflow: hidden;
	margin-right: 10px;
	border: 2px solid transparent; 
	transition: border-color 0.3s ease-in-out;
  }

  .profile-picture:hover {
	border-color: #069327; 
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
  </style>
  