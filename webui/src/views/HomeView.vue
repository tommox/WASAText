<template>
	<div class="flex h-screen">
	  <!-- Sezione Sinistra: Lista Chat -->
	  <div class="w-1/3 h-full">
		<ChatList 
		  :chats="chats" 
		  @chatSelected="selectedChat = $event" 
		/>
	  </div>
  
	  <!-- Sezione Destra: Chat Attiva o Placeholder -->
	  <div class="w-2/3 h-full flex items-center justify-center bg-white">
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
  import axios from "axios";
  
  export default {
	components: { ChatList, ChatWindow, EmptyChat },
	data() {
	  return {
		selectedChat: null,
		chats: [],
	  };
	},
	methods: {
	  async fetchChats() {
		try {
		  const response = await axios.get("/conversations");
		  this.chats = response.data;
		} catch (error) {
		  console.error("Errore nel caricamento delle conversazioni:", error);
		}
	  }
	},
	created() {
	  this.fetchChats();
	}
  };
  </script>
  
  <style scoped>
  .h-screen {
	height: 100vh;
  }
  .h-full {
	height: 100%;
  }
  </style>
  