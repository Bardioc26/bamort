<template>
  <div class="fullwidth-page" style="display: flex; justify-content: center; align-items: center; min-height: 100vh;">
    <div class="card" style="max-width: 400px; width: 100%; margin: 20px;">
      <div class="page-header">
        <h2>Login</h2>
      </div>
      
      <form @submit.prevent="login">
        <div class="form-group">
          <label for="username">Username</label>
          <input
            v-model="username"
            type="text"
            id="username"
            name="username"
            class="form-control"
            placeholder="Username"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="password">Password</label>
          <input
            v-model="password"
            type="password"
            id="password"
            name="password"
            class="form-control"
            placeholder="Password"
            required
          />
        </div>
        
        <button type="submit" class="btn btn-primary" style="width: 100%; margin-top: 10px;">
          Login
        </button>
      </form>
      
      <div v-if="error" class="badge badge-danger" style="width: 100%; margin-top: 15px; text-align: center; display: block;">
        {{ error }}
      </div>
      
      <div style="text-align: center; margin-top: 20px; padding-top: 15px; border-top: 1px solid #dee2e6;">
        <p>Don't have an account? <router-link to="/register" class="btn btn-secondary">Register here</router-link></p>
      </div>
    </div>
  </div>
</template>

<script>
import API from '../utils/api'

export default {
  data() {
    return {
      username: '',
      password: '',
      error: '',
    }
  },
  methods: {
    async login() {
      try {
        const response = await API.post('/login', {
          username: this.username,
          password: this.password,
        })
        localStorage.setItem('token', response.data.token)
        // Emit auth change event
        window.dispatchEvent(new Event('auth-changed'))
        this.$router.push('/dashboard')
      } catch (err) {
        this.error = 'Invalid credentials'
      }
    },
  },
}
</script>

<style scoped>
/* No custom CSS needed - using main.css classes */
</style>
