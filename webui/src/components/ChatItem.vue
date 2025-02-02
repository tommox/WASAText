<template>
    <div 
      class="flex items-center p-4 border border-black bg-white cursor-pointer transition duration-200"
      @click="$emit('selectChat', chat)"
    >
      <!-- Immagine profilo -->
      <img :src="chat.avatar || require('@/assets/images/user.png')" class="w-12 h-12 rounded-full border border-black">
  
      <!-- Dettagli chat -->
      <div class="flex-grow ml-3">
        <div class="font-bold text-black">{{ chat.name }}</div>
        <div class="text-sm text-black truncate">{{ chat.lastMessage || "Nessun messaggio" }}</div>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    props: { chat: Object }
  };
  </script>
  
  <style scoped>
  div:hover {
    background-color: #f0f0f0;
  }
  </style>
  