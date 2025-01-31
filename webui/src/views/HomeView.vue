<script>

export default {
	data: function () {
		return {
			errormsg: null,
		}
	},

	methods: {
		
		async loadStream() {
			try {
				this.errormsg = null
				let response = await this.$axios.get("/Users/" + localStorage.getItem('token') + "/home")
			} catch (e) {
				console.log(e.toString())
				this.errormsg = e.toString()
			}
		}
	},

	async mounted() {
		await this.loadStream()
	}

}
</script>


<template>
	<div class="h-screen flex flex-col bg-gray-100">
	  <!-- Barra di ricerca -->
	  <div class="bg-white p-4 shadow-md">
		<input 
		  type="text" 
		  placeholder="Cerca chat..." 
		  class="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-green-600"
		/>
	  </div>
  
	  <!-- Lista Chat -->
	  <div class="flex-1 overflow-y-auto p-4">
		<div class="flex items-center p-4 bg-white rounded-lg shadow-md mb-4 cursor-pointer">
		  <img src="@/assets/images/user.png" alt="User Avatar" class="w-12 h-12 rounded-full mr-4" />
		  <div>
			<h2 class="font-semibold">Mario Rossi</h2>
			<p class="text-gray-500 text-sm">Ciao, come stai?</p>
		  </div>
		</div>
	  </div>
	</div>
  </template>
 
  <style scoped>
  </style>
  