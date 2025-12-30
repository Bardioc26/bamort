<template>
  <nav class="top-nav"><!---<nav class="menu"> --->
    <ul>
      <li>
        <router-link to="/" active-class="active">{{ $t('menu.Home') }}</router-link>
      </li>
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
      <li v-if="isLoggedIn && isMaintainer">
        <router-link to="/maintenance" active-class="active">{{ $t('menu.Maintenance') }}</router-link>
      </li>
      <li v-if="isLoggedIn && isAdmin">
        <router-link to="/users" active-class="active">{{ $t('menu.UserManagement') }}</router-link>
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
import { useUserStore } from "../stores/userStore";


export default {
  name: "Menu",
  components: {
    LanguageSwitcher,
  },
  data() {
    return {
      userStore: null
    }
  },
  async created() {
    this.userStore = useUserStore()
    if (isLoggedIn() && !this.userStore.currentUser) {
      await this.userStore.fetchCurrentUser()
    }
    // Listen for auth changes to refresh user data
    window.addEventListener('auth-changed', this.handleAuthChange)
  },
  beforeUnmount() {
    window.removeEventListener('auth-changed', this.handleAuthChange)
  },
  computed: {
    isLoggedIn() {
      return isLoggedIn();
    },
    isAdmin() {
      return this.userStore?.isAdmin || false
    },
    isMaintainer() {
      return this.userStore?.isMaintainer || false
    }
  },
  methods: {
    logout() {
      logout();
      this.userStore.clearUser()
      // Emit auth change event
      window.dispatchEvent(new Event('auth-changed'));
      this.$router.push("/");
    },
    async handleAuthChange() {
      if (isLoggedIn()) {
        await this.userStore.fetchCurrentUser()
      } else {
        this.userStore.clearUser()
      }
    }
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
