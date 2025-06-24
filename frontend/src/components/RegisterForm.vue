<template>
  <form @submit.prevent="register">
    <h2>Register</h2>
    <input v-model="username" placeholder="Username" required />
    <input v-model="email" type="email" placeholder="Email" required />
    <input v-model="password" type="password" placeholder="Password" required />
    <button type="submit">Register</button>
    <p v-if="error" class="error">{{ error }}</p>
    <p v-if="success" class="success">{{ success }}</p>
  </form>
</template>

<script>
import API from "../utils/api";

export default {
  data() {
    return {
      username: "",
      email: "",
      password: "",
      error: "",
      success: "",
    };
  },
  methods: {
    async register() {
      try {
        const response = await API.post('/api/register', {
          username: this.username,
          password: this.password,
          email: this.email
        });
        this.success = "Registration successful!";
        this.error = "";
        this.password = "";
      } catch (err) {
        this.error = err.response?.data?.error || "Registration failed.";
        this.success = "";
      }
    },
  },
};
</script>

<style>
.error {
  color: red;
}
.success {
  color: green;
}
</style>
