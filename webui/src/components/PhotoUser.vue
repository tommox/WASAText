<template>
    <img
      :src="imageSrc"
      :alt="altText"
      class="w-12 h-12 rounded-full object-cover"
    />
  </template>
  
  <script>
  export default {
    props: {
      src: {
        type: String,
        default: null, // L'utente potrebbe non aver caricato l'immagine
      },
      altText: {
        type: String,
        default: "User Photo",
      },
    },
    computed: {
      imageSrc() {
        return this.src || require('../assets/images/user.png'); // Placeholder se manca la foto
      },
    },
  };
  </script>