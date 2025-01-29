<script>
export default {
	data() {
		return {
			logged: false, // Stato di login
		};
	},

	methods: {
		updateLogged(newLogged) {
			this.logged = newLogged;
			if (newLogged) {
				this.$router.replace("/dashboard"); // Cambia destinazione se necessario
			}
		},
	},

	mounted() {
		// Controllo token per autenticazione
		if (localStorage.getItem('token')) {
			this.logged = true;
			this.$router.replace("/home"); // Reindirizza se già loggato
		}
	},
};
</script>

<template>
	<div class="container-fluid">
		<div class="row justify-content-center align-items-center vh-100">
			<div class="col-md-4">
				<RouterView @updatedLoggedChild="updateLogged" />
			</div>
		</div>
	</div>
</template>

<style>
/* Stile base per il layout */
.container-fluid {
	background-color: #f8f9fa; /* Grigio chiaro */
}
</style>