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
      deep: true,
      handler: "fetchUserPhoto"
    }
  },
  methods: {
    async fetchUserPhoto() {
      console.log("chat, sono in chatitem: ",this.chat);
      if (!this.chat || !this.chat.recipient_id) return;
      this.avatarUrl = defaultAvatar;
      try {
        const response = await axios.get(`${__API_URL__}/users/${this.chat.recipient_id}/photo`, {
          responseType: "blob"
        });
        if (response.data.size === 0) {
          this.avatarUrl = defaultAvatar;
          return;
        }
        const imageUrl = URL.createObjectURL(response.data);
        this.avatarUrl = ""; 
        this.$nextTick(() => {
          this.avatarUrl = imageUrl;
        });
      } catch (error) {
        console.error("Errore nel recupero della foto profilo:", error);
        this.avatarUrl = defaultAvatar;
      }
    }
  }
};
</script>


<style scoped>
.chat-item {
  display: flex; 
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #ddd;
  cursor: pointer;
}

.profile-img {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 15px; 
}

.chat-name {
  font-weight: bold;
  font-size: 16px;
}

.chat-last-message {
  color: #666;
  font-size: 14px;
}

.chat-item:hover {
  background-color: #f0f0f0;
}

.no-hover:hover {
  background-color: transparent; 
}
</style>
