<template>
  <div id="app">
    <!-- System Alert Banner -->
    <SystemAlert />
    
    <!-- Menu nur anzeigen wenn eingeloggt -->
    <Menu v-if="isLoggedIn" />
     <!-- Main Content Area -->
     <main class="main-content" :class="{ 'full-width': !isLoggedIn }">
      <router-view />
    </main>
  </div>
</template>

<script>
import Menu from "./components/Menu.vue";
import SystemAlert from "./components/SystemAlert.vue";

export default {
  components: {
    Menu,
    SystemAlert,
  },
  data() {
    return {
      loggedIn: false
    }
  },
  computed: {
    isLoggedIn() {
      return this.loggedIn;
    }
  },
  mounted() {
    this.checkAuthStatus();
    // Listen for storage changes
    window.addEventListener('storage', this.handleStorageChange);
    // Listen for custom auth events
    window.addEventListener('auth-changed', this.handleAuthChange);
  },
  beforeUnmount() {
    window.removeEventListener('storage', this.handleStorageChange);
    window.removeEventListener('auth-changed', this.handleAuthChange);
  },
  methods: {
    checkAuthStatus() {
      const token = localStorage.getItem('token');
      this.loggedIn = !!token;
    },
    handleStorageChange() {
      this.checkAuthStatus();
    },
    handleAuthChange() {
      this.checkAuthStatus();
    }
  }
};
</script>

<style src="./assets/main.css"></style>
<style>
/* Global styles can go here */

/* Full-width Layout für Login/Register Seiten */
.main-content.full-width {
  margin: 0;
  padding: 2rem;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}

/* Zentrierte Login-Form */
.main-content.full-width form {
  max-width: 400px;
  width: 100%;
  padding: 2rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: white;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}
</style>
