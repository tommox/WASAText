<template>
	<div class="container-fluid h-100 d-flex">

	  <div class="left-panel">

		<div class="user-profile">
		  <div class="profile-picture">
			<img :src="userImage" alt="Foto Profilo" class="profile-img" />
		  </div>
		  <span class="user-name">{{ nickname }}</span>
		</div>

		<!-- Lista Chat (con barra di ricerca dentro) -->
		<ChatList @chatSelected="selectedChat = $event" />
	  </div>
  
	  <!-- Sezione Destra: Chat Attiva o Placeholder -->
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
	}
  };
  </script>
  
  <style scoped>
  /* Contenitore principale */
  .container-fluid {
	height: 100vh;
  }
  
  /* Sezione sinistra */
  .left-panel {
	width: 30%;
	background: #ffffff;
	padding: 20px;
	border-right: 1px solid #ddd;
	height: 100vh;
	overflow-y: auto;
  }
  
  /* Profilo utente */
  .user-profile {
	display: flex;
	align-items: center;
	padding: 15px;
	margin-bottom: 10px;
	background: rgba(255, 255, 255, 0.9);
	border-radius: 10px;
  }
  
  /* Quadrato per la foto profilo */
  .profile-picture {
	width: 50px;
	height: 50px;
	border-radius: 50%;
	overflow: hidden;
	margin-right: 10px;
	border: 2px solid #007bff;
  }
  
  /* Immagine del profilo */
  .profile-img {
	width: 100%;
	height: 100%;
	object-fit: cover;
  }
  
  /* Nome utente */
  .user-name {
	font-size: 18px;
	font-weight: bold;
	color: #333;
  }
  
  /* Sezione destra (chat attiva) */
  .right-panel {
	width: 70%;
	padding: 20px;
	background: #eeecec;
	height: 100vh;
	overflow-y: auto;
  }
  </style>
  