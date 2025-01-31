<script>
export default {
  name: 'LoginView',
  data() {
    return {
      nickname: '',
      errorMessage: ''
    }
  },
  methods: {
    async login() {
      // Validazione minima sul nickname
      if (this.nickname.length < 3 || this.nickname.length > 16) {
        this.errorMessage = 'Il nickname deve essere tra 3 e 16 caratteri!';
        return;
      }

      try {
        // Chiamata POST al tuo endpoint /session
        // Il backend richiede un JSON con chiave "Nickname"
        const response = await this.$axios.post('/session', {
          Nickname: this.nickname
        });
        
        // Se lo status è 200 o 201, la creazione o il recupero dell'utente è andato a buon fine
        if (response.status === 200 || response.status === 201) {
          // Salviamo l'User_id in localStorage
          localStorage.setItem('User_id', response.data.User_id);
          // Se vuoi, salva anche il Nickname:
          // localStorage.setItem('Nickname', response.data.Nickname);

          // Reindirizziamo l'utente alla Home
          this.$router.push('/home');
        }
      } catch (error) {
        if (error.response) {
          if (error.response.status === 400) {
            this.errorMessage = 'Nickname non valido (3-16 caratteri).';
          } else if (error.response.status === 409) {
            this.errorMessage = 'Nickname già in uso.';
          } else if (error.response.status === 500) {
            this.errorMessage = 'Errore del server. Riprova più tardi.';
          } else {
            this.errorMessage = 'Errore sconosciuto: ' + error.response.status;
          }
        } else {
          this.errorMessage = 'Connessione al server non riuscita.';
        }
      }
    }
  }
}
</script>

<template>
  <div class="container-fluid h-100 w-100 d-flex justify-content-center m-0 p-0 login">
    <div class="login-card p-5 text-center">
      <!-- Logo (immagine) -->
      <img
        src="../assets/images/WASAText.png"
        alt="WASAText logo"
        style="width: 520px; height: 180px;"
      />

      <!-- Form di login -->
      <form @submit.prevent="login">
        <div class="form-group">
          <input
            type="text"
            class="form-control"
            v-model="nickname"
            placeholder="Nickname"
          />
        </div>
        <button
          class="btn btn-primary"
          :disabled="nickname.length < 3 || nickname.length > 16"
        >
          Accedi
        </button>
      </form>

      <!-- Mostriamo eventuali messaggi di errore -->
      <p v-if="errorMessage" class="text-danger mt-3">
        {{ errorMessage }}
      </p>
    </div>
  </div>
</template>

<style>
/* Stili personalizzati, come da tuo esempio */
.text-primary {
  color: rgba(0, 0, 200, 255) !important;
}

.login {
  display: flex;
  align-items: center;
  justify-content: center;
  background-image: url("../assets/images/sfondo.jpeg");
  background-size: cover;
  background-position: center;
  height: 100vh;
  width: 100%;
  position: absolute;
  top: 0;
  left: 0;
}

.login-card {
  background-color: rgba(241, 241, 241, 0.35);
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.3);
  border-radius: 25px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  text-align: center;
  width: 350px;
  height: 400px;
}

.login-card img {
  width: auto;
  height: auto;
  text-align: center;
  margin-bottom: 10px;
}

.btn-primary {
  margin-top: 30px;
  margin-left: auto;
  margin-right: auto;
  display: block;
  background-color: #069327 !important;
  color: white;
  border-radius: 50px;
  padding: 15px 25px;
  font-weight: bold;
  font-size: 18px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 10px rgba(0, 0, 0, 0.1);
}

.btn-primary:active {
  transform: translateY(0);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.form-control {
  width: 300px !important;
}
</style>