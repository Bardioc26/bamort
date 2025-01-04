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
        const response = await API.post("/register", {
          username: this.username,
          email: this.email,
          password: this.password,
        });
        this.success = "Registration successful! You can now log in.";
        this.error = "";
        this.username = "";
        this.email = "";
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
