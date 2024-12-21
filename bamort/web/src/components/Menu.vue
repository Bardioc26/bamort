<template>
  <nav class="menu">
    <ul>
      <li>
        <router-link to="/dashboard">Dashboard</router-link>
      </li>
      <li v-if="!isLoggedIn">
        <router-link to="/register">Register</router-link>
      </li>
      <li v-if="isLoggedIn">
        <button @click="logout">Logout</button>
      </li>
    </ul>
  </nav>
</template>

<script>
import { isLoggedIn, logout } from "../utils/auth";

export default {
  computed: {
    isLoggedIn() {
      return isLoggedIn();
    },
  },
  methods: {
    logout() {
      logout();
      this.$router.push("/");
    },
  },
};
</script>
