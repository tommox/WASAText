<template>
	<div class="container-fluid h-100 d-flex">
	  <div class="left-panel">
		<div class="user-profile">
			<label class="profile-picture">
				<!-- Immagine del profilo -->
				<img :src="userImage" alt="Foto Profilo" class="profile-img" />

				<input type="file" @change="uploadProfilePicture" accept="image/*" class="upload-input" />
			</label>

		<span class="user-name">{{ nickname }}</span>
	  <button @click="logout" class="logout-btn">Logout</button>
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
	},
	methods: {
    logout() {
        this.$router.replace("/login");
    },

	async uploadProfilePicture(event) {
		const file = event.target.files[0]; 
		if (!file) return;

		const token = localStorage.getItem("token"); 
		const formData = new FormData();
		formData.append("profile_picture", file);

		try {
		const response = await axios.post(`${__API_URL__}/users/${token}/photo`, formData, {
			headers: {
			Authorization: `Bearer ${token}`,
			"Content-Type": "multipart/form-data"
			}
		});

		this.userImage = response.data.profile_picture;
		localStorage.setItem("profileImage", this.userImage);
		} catch (error) {
		console.error("Errore nel caricamento dell'immagine:", error);
		}
	}
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

  .logout-btn {
    position: absolute;
    top: 40px;
    right: 1000px;
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
	position: relative;
	width: 50px;
	height: 50px;
	border-radius: 50%;
	overflow: hidden;
	margin-right: 10px;
	border: 2px solid #069327;
	cursor: pointer; 
}
  
  /* Immagine del profilo */
  .profile-img {
	width: 100%;
	height: 100%;
	object-fit: cover;
}

	/* Input file nascosto */
   .upload-input {
	display: none; 
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
  