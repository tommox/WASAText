<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/");
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div class="h-screen flex flex-col bg-gray-100">
	  <!-- Barra superiore -->
	  <div class="bg-green-600 text-white p-4 flex justify-between items-center shadow-md">
		<h1 class="text-xl font-bold">WhatsApp</h1>
		<div class="flex gap-4">
		  <BellIcon class="w-6 h-6" />
		  <EllipsisVerticalIcon class="w-6 h-6" />
		</div>
	  </div>
  
	  <!-- Lista Chat -->
	  <div class="flex-1 overflow-y-auto">
		<div v-for="chat in chats" :key="chat.id" class="flex items-center p-4 border-b border-gray-300 hover:bg-gray-200 cursor-pointer">
		  <img :src="chat.avatar" alt="Avatar" class="w-12 h-12 rounded-full mr-4" />
		  <div class="flex-1">
			<h2 class="font-semibold">{{ chat.name }}</h2>
			<p class="text-gray-500 text-sm truncate">{{ chat.lastMessage }}</p>
		  </div>
		  <span class="text-gray-400 text-xs">{{ chat.time }}</span>
		</div>
	  </div>
  
	  <!-- Navigazione inferiore -->
	  <div class="bg-white shadow-md p-3 flex justify-around border-t border-gray-300">
		<HomeIcon class="w-6 h-6 text-gray-500" />
		<UsersIcon class="w-6 h-6 text-gray-500" />
		<PhoneIcon class="w-6 h-6 text-gray-500" />
		<CogIcon class="w-6 h-6 text-gray-500" />
	  </div>
	</div>
  </template>
  
  <script>
  import { HomeIcon, UsersIcon, PhoneIcon, CogIcon, BellIcon, EllipsisVerticalIcon } from '@heroicons/vue/24/outline';
  
  export default {
	components: { HomeIcon, UsersIcon, PhoneIcon, CogIcon, BellIcon, EllipsisVerticalIcon },
	data() {
	  return {
		chats: [
		  { id: 1, name: "Mario Rossi", avatar: "https://i.pravatar.cc/100?img=1", lastMessage: "Ciao, come stai?", time: "12:30" },
		  { id: 2, name: "Anna Bianchi", avatar: "https://i.pravatar.cc/100?img=2", lastMessage: "Ci vediamo domani!", time: "11:15" },
		  { id: 3, name: "Luca Verdi", avatar: "https://i.pravatar.cc/100?img=3", lastMessage: "Ok, a dopo!", time: "10:45" },
		],
	  };
	},
  };
  </script>
  
  <style scoped>
	.truncate {
	  overflow: hidden;
	  text-overflow: ellipsis;
	  white-space: nowrap;
	}
  </style>
  
