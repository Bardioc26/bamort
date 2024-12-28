<template>
  <form @submit.prevent="login">
    <h2>Login</h2>
    <input v-model="username" placeholder="Username" required />
    <input v-model="password" type="password" placeholder="Password" required />
    <button type="submit">Login</button>
    <p v-if="error" class="error">{{ error }}</p>
    <p>Don't have an account? <router-link to="/register">Register here</router-link>.</p>
  </form>
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
        this.$router.push('/dashboard')
      } catch (err) {
        this.error = 'Invalid credentials'
      }
    },
  },
}
</script>

<style>
.error {
  color: red;
}
</style>
