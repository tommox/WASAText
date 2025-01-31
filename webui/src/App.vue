<template>
<!-- Componente radice dell'applicazione -->
<router-view />
</template>
  
<script>
export default {
	data(){
		return{
			logged: false,
			searchValue: "",
		}
	},
	methods:{
		logout(newValue){
			this.logged = newValue
			this.$router.replace("/login")
		},
		updateLogged(newLogged){
			this.logged = newLogged
		},
		updateView(newRoute){
			this.$router.replace(newRoute)
		},
	},

	created(){
		if (!localStorage.getItem('notFirstStart')){
			localStorage.clear()
			localStorage.setItem('notFirstStart',true)
			// console.log("first start")
		}
		
	},

	mounted(){
		// console.log("Devo modificare ancora lo stile!")
		if (!localStorage.getItem('token')){
			this.$router.replace("/login")
		}else{
			this.logged = true
		}
	},
}
</script>
  
<style>
/* Se vuoi avere degli stili globali minimi, mettili qui */
/* Altrimenti, puoi lasciare anche vuoto. */
</style>