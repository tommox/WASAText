<template>
  <div 
    class="chat-item flex items-center p-3 border-b cursor-pointer transition duration-200"
    @click="$emit('selectChat', chat)"
  >
    <!-- Immagine profilo -->
    <img :src="chat.avatar || defaultAvatar" alt="Foto Profilo" class="profile-img mr-3" />
    
    <!-- Dettagli chat -->
    <div class="flex-grow">
      <div class="font-bold text-gray-800">{{ chat.name }}</div>
      <div class="text-sm text-gray-600 truncate no-hover">{{ chat.lastMessage || "Nessun messaggio" }}</div>
    </div>
  </div>
</template>

<script>
import defaultAvatar from "@/assets/images/user.png";

export default {
  props: { chat: Object },
  data() {
    return { 
      defaultAvatar,
    };
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
