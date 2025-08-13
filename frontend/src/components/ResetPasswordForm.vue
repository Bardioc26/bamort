<template>
  <div class="fullwidth-page" style="display: flex; justify-content: center; align-items: center; min-height: 100vh;">
    <div class="card" style="max-width: 400px; width: 100%; margin: 20px;">
      
      <!-- Loading State -->
      <div v-if="isValidating" class="page-header" style="text-align: center;">
        <h2>Validiere Reset-Link...</h2>
        <div style="margin-top: 20px;">
          <div class="spinner" style="margin: 0 auto;"></div>
        </div>
      </div>
      
      <!-- Invalid Token -->
      <div v-else-if="!isValidToken" class="page-header">
        <h2>Ungültiger Reset-Link</h2>
        <div class="badge badge-danger" style="width: 100%; margin-top: 15px; text-align: center; display: block;">
          <p style="margin: 10px 0;">
            Dieser Reset-Link ist ungültig oder abgelaufen.
          </p>
          <p style="font-size: 0.9em; margin: 5px 0;">
            Bitte fordern Sie einen neuen Reset-Link an.
          </p>
        </div>
        <div style="text-align: center; margin-top: 20px;">
          <router-link to="/forgot-password" class="btn btn-primary">
            Neuen Reset-Link anfordern
          </router-link>
        </div>
      </div>
      
      <!-- Valid Token - Reset Form -->
      <div v-else>
        <div class="page-header">
          <h2>Neues Passwort setzen</h2>
          <p style="color: #666; font-size: 0.9em; margin-top: 10px;" v-if="userInfo.username">
            Für Benutzer: <strong>{{ userInfo.username }}</strong>
          </p>
        </div>
        
        <form @submit.prevent="resetPassword" v-if="!resetSuccess">
          <div class="form-group">
            <label for="newPassword">Neues Passwort</label>
            <input
              v-model="newPassword"
              type="password"
              id="newPassword"
              name="newPassword"
              class="form-control"
              placeholder="Mindestens 6 Zeichen"
              required
              minlength="6"
            />
          </div>
          
          <div class="form-group">
            <label for="confirmPassword">Passwort bestätigen</label>
            <input
              v-model="confirmPassword"
              type="password"
              id="confirmPassword"
              name="confirmPassword"
              class="form-control"
              placeholder="Passwort wiederholen"
              required
              minlength="6"
            />
          </div>
          
          <button 
            type="submit" 
            class="btn btn-primary" 
            style="width: 100%; margin-top: 10px;"
            :disabled="isResetting || !passwordsMatch"
          >
            <span v-if="isResetting">Passwort wird gesetzt...</span>
            <span v-else>Passwort zurücksetzen</span>
          </button>
          
          <div v-if="!passwordsMatch && confirmPassword" class="badge badge-warning" style="width: 100%; margin-top: 10px; text-align: center; display: block; font-size: 0.8em;">
            Die Passwörter stimmen nicht überein
          </div>
        </form>
        
        <!-- Success Message -->
        <div v-if="resetSuccess" class="badge badge-success" style="width: 100%; margin-top: 15px; text-align: center; display: block;">
          <p style="margin: 10px 0;">
            <strong>Passwort erfolgreich zurückgesetzt!</strong>
          </p>
          <p style="font-size: 0.9em; margin: 5px 0;">
            Sie können sich jetzt mit Ihrem neuen Passwort anmelden.
          </p>
          <div style="margin-top: 15px;">
            <router-link to="/" class="btn btn-primary">
              Zum Login
            </router-link>
          </div>
        </div>
        
        <div v-if="error" class="badge badge-danger" style="width: 100%; margin-top: 15px; text-align: center; display: block;">
          {{ error }}
        </div>
      </div>
      
      <!-- Back to Login (always visible) -->
      <div style="text-align: center; margin-top: 20px; padding-top: 15px; border-top: 1px solid #dee2e6;" v-if="!isValidating && !resetSuccess">
        <router-link to="/" class="btn btn-secondary">
          Zurück zum Login
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
import API from '../utils/api'

export default {
  name: 'ResetPasswordForm',
  data() {
    return {
      token: '',
      newPassword: '',
      confirmPassword: '',
      error: '',
      isValidating: true,
      isValidToken: false,
      isResetting: false,
      resetSuccess: false,
      userInfo: {},
    }
  },
  computed: {
    passwordsMatch() {
      return this.newPassword === this.confirmPassword
    }
  },
  async mounted() {
    // Get token from URL query parameter
    this.token = this.$route.query.token
    
    if (!this.token) {
      this.isValidating = false
      this.isValidToken = false
      this.error = 'Kein Reset-Token gefunden'
      return
    }
    
    await this.validateToken()
  },
  methods: {
    async validateToken() {
      try {
        const response = await API.get(`/password-reset/validate/${this.token}`)
        
        this.isValidToken = response.data.valid
        this.userInfo = {
          username: response.data.username,
          expires: response.data.expires
        }
        
        console.log('Token validation successful:', response.data)
      } catch (err) {
        console.error('Token validation error:', err)
        this.isValidToken = false
        this.error = 'Reset-Link ist ungültig oder abgelaufen'
      } finally {
        this.isValidating = false
      }
    },
    
    async resetPassword() {
      if (!this.passwordsMatch) {
        this.error = 'Die Passwörter stimmen nicht überein'
        return
      }
      
      if (this.newPassword.length < 6) {
        this.error = 'Das Passwort muss mindestens 6 Zeichen lang sein'
        return
      }
      
      this.error = ''
      this.isResetting = true
      
      try {
        await API.post('/password-reset/reset', {
          token: this.token,
          new_password: this.newPassword,
        })
        
        this.resetSuccess = true
        console.log('Password reset successful')
      } catch (err) {
        console.error('Password reset error:', err)
        this.error = err.response?.data?.error || 'Fehler beim Zurücksetzen des Passworts. Versuchen Sie es erneut.'
      } finally {
        this.isResetting = false
      }
    },
  },
}
</script>

<style scoped>
.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid #f3f3f3;
  border-top: 2px solid #007bff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>
