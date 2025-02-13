<template>
  <div class="chat-item" @click="$emit('selectChat', chat)">
    <img :src="avatarUrl" alt="Avatar" class="profile-img" />
    <div class="chat-details">
      <div class="chat-name">{{ chat.name }}</div>
      <div class="chat-last-message">{{ chat.lastMessage || 'Nessun messaggio' }}</div>
    </div>
  </div>
</template>

<script>
import defaultAvatar from "@/assets/images/user.png";
import axios from "axios";

export default {
  props: { chat: Object },
  data() {
    return {
      avatarUrl: defaultAvatar
    };
  },
  watch: {
    chat: {
      immediate: true,
      handler: "loadAvatar"
    }
  },
  methods: {
    async loadAvatar() {
      this.avatarUrl = defaultAvatar; 
      try {
        const response = await axios.get(`${__API_URL__}/users/${this.chat.recipient_id}/photo`, {
          responseType: "blob"
        });

        if (response.data.size > 0) {
          this.avatarUrl = URL.createObjectURL(response.data);
        }
      } catch (error) {
        console.error("Errore nel caricamento della foto profilo:", error);
      }
    }
  }
};
</script>

<style scoped>
.chat-item {
  height: 100px; 
  border-bottom: 1px solid #ddd;
}

.profile-img {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
}

.chat-item:hover {
  background-color: #f0f0f0;
}

.no-hover:hover {
  background-color: transparent; 
}
</style>
