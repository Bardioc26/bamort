<template>
  <div class="fullwidth-page" style="display: flex; justify-content: center; align-items: center; min-height: 100vh;">
    <div class="card" style="max-width: 400px; width: 100%; margin: 20px;">
      <div class="page-header">
        <h2>{{ $t('auth.register') }}</h2>
      </div>
      
      <form @submit.prevent="register">
        <div class="form-group">
          <label for="username">{{ $t('auth.username') }}</label>
          <input
            v-model="username"
            type="text"
            id="username"
            name="username"
            class="form-control"
            :placeholder="$t('auth.username')"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="email">{{ $t('auth.email') }}</label>
          <input
            v-model="email"
            type="email"
            id="email"
            name="email"
            class="form-control"
            :placeholder="$t('auth.email')"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="password">{{ $t('auth.password') }}</label>
          <input
            v-model="password"
            type="password"
            id="password"
            name="password"
            class="form-control"
            :placeholder="$t('auth.password')"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="confirmPassword">{{ $t('auth.confirmPassword') }}</label>
          <input
            v-model="confirmPassword"
            type="password"
            id="confirmPassword"
            name="confirmPassword"
            class="form-control"
            :placeholder="$t('auth.confirmPassword')"
            required
          />
        </div>
        
        <button type="submit" class="btn btn-primary" style="width: 100%; margin-top: 10px;">
          {{ $t('auth.register') }}
        </button>
      </form>
      
      <div v-if="error" class="badge badge-danger" style="width: 100%; margin-top: 15px; text-align: center; display: block;">
        {{ error }}
      </div>
      
      <div v-if="success" class="badge badge-success" style="width: 100%; margin-top: 15px; text-align: center; display: block;">
        {{ success }}
      </div>
      
      <div style="text-align: center; margin-top: 20px; padding-top: 15px; border-top: 1px solid #dee2e6;">
        <p>{{ $t('auth.alreadyHaveAccount') }} <router-link to="/" class="btn btn-secondary">{{ $t('auth.loginHere') }}</router-link></p>
      </div>
    </div>
  </div>
</template>

<script>
import API from "../utils/api";

export default {
  data() {
    return {
      username: "",
      email: "",
      password: "",
      confirmPassword: "",
      error: "",
      success: "",
    };
  },
  methods: {
    async register() {
      // Validate passwords match
      if (this.password !== this.confirmPassword) {
        this.error = this.$t('auth.passwordsDontMatch');
        this.success = "";
        return;
      }
      
      try {
        const response = await API.post('/register', {
          username: this.username,
          password: this.password,
          email: this.email
        });
        this.success = this.$t('auth.registrationSuccess');
        this.error = "";
        this.password = "";
        this.confirmPassword = "";
      } catch (err) {
        this.error = err.response?.data?.error || this.$t('auth.registrationFailed');
        this.success = "";
      }
    },
  },
};
</script>

<style scoped>
/* All common styles moved to main.css */
</style>
