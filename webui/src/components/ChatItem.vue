<template>
  <div class="chat-item" @click="$emit('selectChat', chat)">
    <img :src="avatarUrl" alt="Avatar" class="profile-img" />
    <div class="chat-details">
      <div class="chat-name">{{ chatName }}</div>
      <div class="chat-last-message">{{ lastMessage || 'Nessun messaggio' }}</div>
    </div>
  </div>
</template>

<script>
import defaultAvatar from "@/assets/images/user.png";
import axios from "axios";
import eventBus from "@/eventBus";

export default {
  props: { chat: Object, type: String }, 
  data() {
    return {
      avatarUrl: defaultAvatar,
      lastMessage: "Nessun messaggio"
    };
  },

  created() {
  eventBus.on("newMessage", (data) => {
    if (data.type === "private" && data.conversation_id === this.chat.conversation_id) {
      this.lastMessage = data.lastMessage;
    } else if (data.type === "group" && data.conversation_id === this.chat.group_conversation_id) {
      this.lastMessage = data.lastMessage;
    }
    });

    if (this.type === "private") {
      this.fetchLastMessage();
    } else if (this.type === "group") {
      this.fetchLastMessageGroup();
    }
  },

  beforeUnmount() {
    eventBus.off("newMessage");
  },

  computed: {
    chatName() {
      return this.type === "private" ? this.chat.name : this.chat.group_name;
    }
  },

  watch: {
    chat: {
      immediate: true,
      deep: true,
      handler() {
        this.fetchUserPhoto();
        if (this.type === "private") {
          this.fetchLastMessage();
        } else if (this.type === "group") {
          this.fetchLastMessageGroup();
        }
      }
    }
  },

  methods: {
    async fetchUserPhoto() {
      if (!this.chat || this.type !== "private") return;
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
    },

    async fetchGroupPhoto() {
      if (!this.chat || this.type !== "group") return;
      this.avatarUrl = defaultAvatar; 

      try {
        const response = await axios.get(`${__API_URL__}/groups/${this.chat.group_id}/photo`, {
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
        console.error("Errore nel recupero della foto del gruppo:", error);
        this.avatarUrl = defaultAvatar;
      }
    },


    async fetchLastMessage() {
      if (!this.chat || !this.chat.group_last_message_id || this.chat.group_last_message_id === "Nessun messaggio") {
        this.lastMessage = "Nessun messaggio";
        return;
      }
      try {
        const response = await axios.get(`${__API_URL__}/messages/${this.chat.lastMessage}?type=${this.type}`, {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        });
        if (response.data && response.data.message_content) {
          this.lastMessage = response.data.message_content;
        }
      } catch (error) {
        console.error("Errore nel recupero dell'ultimo messaggio:", error);
        this.lastMessage = "Nessun messaggio";
      }
    },

    async fetchLastMessageGroup() {
      if (!this.chat || !this.chat.group_last_message_id || this.chat.group_last_message_id === "Nessun messaggio") {
        this.lastMessage = "Nessun messaggio";
        return;
      }
      try {
        const response = await axios.get(`${__API_URL__}/messages/${this.chat.group_last_message_id}?type=${this.type}`, {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        });
        console.log("resp: ",response);
        if (response.data && response.data.message_content) {
          this.lastMessage = response.data.message_content;
        }
      } catch (error) {
        console.error("Errore nel recupero dell'ultimo messaggio del gruppo:", error);
        this.lastMessage = "Nessun messaggio";
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

.chat-details {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  min-width: 0; 
}

.chat-last-message {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex-grow: 1;
  min-width: 0;
}

.chat-item:hover {
  background-color: #f0f0f0;
}
</style>
