<template>
  <div class="fullwidth-page" style="display: flex; justify-content: center; align-items: center; min-height: 100vh;">
    <div class="card" style="max-width: 400px; width: 100%; margin: 20px;">
      <div class="page-header">
        <h2>Passwort zurücksetzen</h2>
        <p style="color: #666; font-size: 0.9em; margin-top: 10px;">
          Geben Sie Ihre E-Mail-Adresse ein, um einen Reset-Link zu erhalten.
        </p>
      </div>
      
      <form @submit.prevent="requestReset" v-if="!submitted">
        <div class="form-group">
          <label for="email">E-Mail-Adresse</label>
          <input
            v-model="email"
            type="email"
            id="email"
            name="email"
            class="form-control"
            placeholder="ihre@email.de"
            required
          />
        </div>
        
        <button 
          type="submit" 
          class="btn btn-primary" 
          style="width: 100%; margin-top: 10px;"
          :disabled="isLoading"
        >
          <span v-if="isLoading">Wird gesendet...</span>
          <span v-else>Reset-Link senden</span>
        </button>
      </form>
      
      <div v-if="submitted" class="badge badge-success" style="width: 100%; margin-top: 15px; text-align: center; display: block;">
        <p style="margin: 10px 0;">
          <strong>E-Mail gesendet!</strong>
        </p>
        <p style="font-size: 0.9em; margin: 5px 0;">
          Falls ein Account mit dieser E-Mail-Adresse existiert, wurde ein Reset-Link gesendet.
        </p>
        <p style="font-size: 0.8em; margin: 5px 0; opacity: 0.8;">
          Prüfen Sie Ihre E-Mails und folgen Sie dem Link.
        </p>
      </div>
      
      <div v-if="error" class="badge badge-danger" style="width: 100%; margin-top: 15px; text-align: center; display: block;">
        {{ error }}
      </div>
      
      <div style="text-align: center; margin-top: 20px; padding-top: 15px; border-top: 1px solid #dee2e6;">
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
  name: 'ForgotPasswordForm',
  data() {
    return {
      email: '',
      error: '',
      submitted: false,
      isLoading: false,
    }
  },
  methods: {
    async requestReset() {
      this.error = ''
      this.isLoading = true

      try {
        await API.post('/password-reset/request', {
          email: this.email,
          redirect_url: window.location.origin, // Aktuelle Frontend-URL
        })
        
        this.submitted = true
        console.log('Password reset email requested for:', this.email)
      } catch (err) {
        console.error('Password reset request error:', err)
        this.error = err.response?.data?.error || 'Fehler beim Senden der E-Mail. Versuchen Sie es später erneut.'
      } finally {
        this.isLoading = false
      }
    },
  },
}
</script>

<style>
/* All common styles moved to main.css */
</style>
