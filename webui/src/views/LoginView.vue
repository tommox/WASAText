<template>
	<div class="login-background">
	  <div class="login-container">
		<div class="logo-container">
		  <!-- Esempio di immagine logo -->
		  <img src="../assets/images/sfondo.jpeg" alt="Logo" class="logo" />
		</div>
  
		<div class="form-container">
		  <input 
			v-model="nickname"
			type="text"
			placeholder="Inserisci il tuo Nickname"
		  />
		  <button @click="login">Accedi</button>
		</div>
	  </div>
	</div>
  </template>
  
  <script>
  export default {
	name: 'LoginView',
	data() {
	  return {
		nickname: ''
	  }
	},
	methods: {
	  login() {
		// Esempio: per ora basta portarti a /home
		// (in seguito potrai fare la chiamata al backend per validare il nickname).
		this.$router.push('/home');
	  }
	}
  }
  </script>
  
  <style scoped>
  /* Copriamo tutta la finestra con un'immagine di sfondo */
  .login-background {
	background: url('../assets/images/sfondo.jpg') no-repeat center center fixed;
	background-size: cover; /* Adattiamo la grandezza dell'immagine allo schermo */
	height: 100vh;          /* Occupiamo l'intera altezza del viewport */
	width: 100%;
	display: flex;
	justify-content: center; /* Centra orizzontalmente il contenitore login */
	align-items: center;     /* Centra verticalmente il contenitore login */
	margin: 0;
	padding: 0;
  }
  
  /* Contenitore della form di login (fondo chiaro, leggermente trasparente ecc.) */
  .login-container {
	background-color: rgba(255, 255, 255, 0.9);
	border-radius: 8px;
	padding: 30px 20px;
	width: 300px; /* Larghezza fissa per la form */
	display: flex;
	flex-direction: column;
	align-items: center;
  }
  
  /* Sezione per il logo */
  .logo-container {
	margin-bottom: 20px;
  }
  
  .logo {
	max-width: 120px; /* Riduci/ingrandisci se serve */
  }
  
  /* Contenitore del form (input + bottone) */
  .form-container {
	width: 100%;
	display: flex;
	flex-direction: column;
	align-items: center;
  }
  
  .form-container input {
	width: 100%;
	padding: 10px;
	margin-bottom: 15px;
	border: 1px solid #ccc;
	border-radius: 4px;
  }
  
  .form-container button {
	width: 100%;
	padding: 10px;
	background-color: #00bfa5;
	color: white;
	border: none;
	border-radius: 4px;
	font-size: 16px;
	cursor: pointer;
  }
  
  .form-container button:hover {
	background-color: #009d87;
  }
  </style>