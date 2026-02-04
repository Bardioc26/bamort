<template>
  <div class="fullwidth-page reset-wrapper">
    <div class="card reset-card">
      <div class="page-header">
        <h2>{{ $t('forgotPassword.title') }}</h2>
        <p class="reset-description">
          {{ $t('forgotPassword.description') }}
        </p>
      </div>

      <form @submit.prevent="requestReset" v-if="!submitted">
        <div class="form-group">
          <label for="email">{{ $t('forgotPassword.emailLabel') }}</label>
          <input
            v-model="email"
            type="email"
            id="email"
            name="email"
            class="form-control"
            :placeholder="$t('forgotPassword.emailPlaceholder')"
            required
          />
        </div>

        <button 
          type="submit" 
          class="btn btn-primary reset-button"
          :disabled="isLoading"
        >
          <span v-if="isLoading">{{ $t('forgotPassword.submitting') }}</span>
          <span v-else>{{ $t('forgotPassword.submit') }}</span>
        </button>
      </form>

      <div v-if="submitted" class="badge badge-success reset-badge">
        <p class="reset-success-title">
          <strong>{{ $t('forgotPassword.successTitle') }}</strong>
        </p>
        <p class="reset-success-text">
          {{ $t('forgotPassword.successInfo') }}
        </p>
        <p class="reset-success-hint">
          {{ $t('forgotPassword.successHint') }}
        </p>
      </div>

      <div v-if="error" class="badge badge-danger reset-badge">
        {{ error }}
      </div>

      <div class="reset-footer">
        <router-link to="/" class="btn btn-secondary">
          {{ $t('forgotPassword.backToLogin') }}
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
        this.error = err.response?.data?.error || this.$t('forgotPassword.error')
      } finally {
        this.isLoading = false
      }
    },
  },
}
</script>

<style scoped>
.reset-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}

.reset-card {
  max-width: 400px;
  width: 100%;
  margin: 20px;
}

.reset-description {
  color: #666;
  font-size: 0.9em;
  margin-top: 10px;
}

.reset-button {
  width: 100%;
  margin-top: 10px;
}

.reset-badge {
  width: 100%;
  margin-top: 15px;
  text-align: center;
  display: block;
}

.reset-success-title {
  margin: 10px 0;
}

.reset-success-text {
  font-size: 0.9em;
  margin: 5px 0;
}

.reset-success-hint {
  font-size: 0.8em;
  margin: 5px 0;
  opacity: 0.8;
}

.reset-footer {
  text-align: center;
  margin-top: 20px;
  padding-top: 15px;
  border-top: 1px solid #dee2e6;
}
</style>
