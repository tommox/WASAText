<template>
	<div class="container-fluid h-100 d-flex">
	  <!-- Sezione Sinistra: Lista Chat -->
	  <div class="left-panel">
		<!-- Logo (solo in HomeView) -->
		<div class="logo-container">
		  <img src="@/assets/images/WASAText.png" alt="Logo" class="logo" />
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
  
  export default {
	components: { ChatList, ChatWindow, EmptyChat },
	data() {
	  return {
		selectedChat: null
	  };
	}
  };
  </script>
  
  <style scoped>
  /* Contenitore per il logo */
  .logo-container {
	text-align: center;
	padding: 0px 0;
  }
  
  /* Stile per il logo */
  .logo {
    width: 300px; 
    height: auto;
    display: inline-block;
  }
  
  /* Stile per il layout */
  .container-fluid {
	height: 100vh;
  }
  
  /* Sezione sinistra (chat list) */
  .left-panel {
	width: 30%;
	background: #ffffff;
	padding: 20px;
	border-right: 1px solid #ddd;
	height: 100vh;
	overflow-y: auto;
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
  