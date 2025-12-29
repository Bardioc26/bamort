<template>
  <nav class="top-nav"><!---<nav class="menu"> --->
    <ul>
      <li>
        <router-link to="/dashboard" active-class="active">{{ $t('menu.Dashboard') }}</router-link>
      </li>
      <li>
        <router-link to="/upload" active-class="active">{{ $t('menu.ImportData') }}</router-link>
      </li>
      <li v-if="isLoggedIn">
        <button @click="logout">{{ $t('menu.Logout') }}</button>
      </li>
      <li v-if="!isLoggedIn">
        <router-link to="/register" active-class="active">{{ $t('menu.Register') }}</router-link>
      </li>
      <li v-if="isLoggedIn">
        <router-link to="/maintenance" active-class="active">{{ $t('menu.Maintenance') }}</router-link>
      </li>
    </ul>
    <div class="menu-right">
      <LanguageSwitcher />
      <router-link v-if="isLoggedIn" to="/profile" active-class="active" class="profile-link">{{ $t('menu.Profile') }}</router-link>
    </div>
  </nav>
</template>

<script>
import { isLoggedIn, logout } from "../utils/auth";
import LanguageSwitcher from "./LanguageSwitcher.vue";


export default {
  name: "Menu",
  components: {
    LanguageSwitcher,
  },
  computed: {
    isLoggedIn() {
      return isLoggedIn();
    },
  },
  methods: {
    logout() {
      logout();
      // Emit auth change event
      window.dispatchEvent(new Event('auth-changed'));
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

.menu-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.profile-link {
  color: white;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: var(--border-radius);
  transition: background-color 0.2s;
}

.profile-link:hover {
  background-color: rgba(255, 255, 255, 0.1);
  text-decoration: none;
}

.profile-link.active {
  background-color: rgba(255, 255, 255, 0.2);
  font-weight: bold;
}
</style>
