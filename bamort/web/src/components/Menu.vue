<template>
  <nav class="menu">
    <ul>
      <li>
        <router-link to="/dashboard" active-class="active">Dashboard</router-link>
      </li>
      <li v-if="!isLoggedIn">
        <router-link to="/register" active-class="active">Register</router-link>
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

<style>
.menu {
  background-color: #333;
  color: white;
  padding: 1rem;
}

.menu ul {
  list-style: none;
  display: flex;
  gap: 1rem;
}

.menu a {
  color: white;
  text-decoration: none;
}

.menu a:hover {
  text-decoration: underline;
}

.menu .active {
  font-weight: bold;
  text-decoration: underline;
}

.menu button {
  background: none;
  border: none;
  color: white;
  cursor: pointer;
}
</style>
