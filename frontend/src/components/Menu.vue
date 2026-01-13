<template>
  <nav class="top-nav">
    <ul class="menu-left">
      <!-- Home -->
      <li>
        <router-link to="/" active-class="active">{{ $t('menu.Home') }}</router-link>
      </li>
      
      <!-- Info Dropdown -->
      <li class="dropdown" @mouseenter="openMenu('info')" @mouseleave="closeMenu('info')">
        <span class="dropdown-trigger">{{ $t('menu.Info') }} ▾</span>
        <ul v-show="showInfoMenu" class="dropdown-menu" @mouseenter="openMenu('info')" @mouseleave="closeMenu('info')">
          <li><router-link to="/help" @click="closeAllMenus">{{ $t('menu.Help') }}</router-link></li>
          <li><router-link to="/sponsors" @click="closeAllMenus">{{ $t('menu.Sponsors') }}</router-link></li>
          <li><router-link to="/system-info" @click="closeAllMenus">{{ $t('menu.SystemInfo') }}</router-link></li>
        </ul>
      </li>
      
      <!-- Characters Dropdown (only when logged in) -->
      <li v-if="isLoggedIn" class="dropdown" @mouseenter="openMenu('char')" @mouseleave="closeMenu('char')">
        <span class="dropdown-trigger">{{ $t('menu.Characters') }} ▾</span>
        <ul v-show="showCharMenu" class="dropdown-menu" @mouseenter="openMenu('char')" @mouseleave="closeMenu('char')">
          <li><router-link to="/dashboard" @click="closeAllMenus">{{ $t('menu.Dashboard') }}</router-link></li>
          <li><router-link to="/upload" @click="closeAllMenus">{{ $t('menu.ImportData') }}</router-link></li>
        </ul>
      </li>
      
      <!-- Admin Dropdown (only for maintainers/admins) -->
      <li v-if="isLoggedIn && (isMaintainer || isAdmin)" class="dropdown" @mouseenter="openMenu('admin')" @mouseleave="closeMenu('admin')">
        <span class="dropdown-trigger">{{ $t('menu.Admin') }} ▾</span>
        <ul v-show="showAdminMenu" class="dropdown-menu" @mouseenter="openMenu('admin')" @mouseleave="closeMenu('admin')">
          <li v-if="isMaintainer"><router-link to="/maintenance" @click="closeAllMenus">{{ $t('menu.Maintenance') }}</router-link></li>
          <li v-if="isAdmin"><router-link to="/users" @click="closeAllMenus">{{ $t('menu.UserManagement') }}</router-link></li>
        </ul>
      </li>
      
      <!-- Register (only when not logged in) -->
      <li v-if="!isLoggedIn">
        <router-link to="/register" active-class="active">{{ $t('menu.Register') }}</router-link>
      </li>
    </ul>
    
    <div class="menu-right">
      <LanguageSwitcher />
      <router-link v-if="isLoggedIn" to="/profile" active-class="active" class="profile-link">
        {{ $t('menu.Profile') }}
      </router-link>
      <button v-if="isLoggedIn" @click="logout" class="logout-btn">{{ $t('menu.Logout') }}</button>
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
      userStore: null,
      showInfoMenu: false,
      showCharMenu: false,
      showAdminMenu: false,
      closeTimeout: null
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
    openMenu(menu) {
      // Clear any pending close timeout
      if (this.closeTimeout) {
        clearTimeout(this.closeTimeout)
        this.closeTimeout = null
      }
      
      // Open the requested menu
      if (menu === 'info') this.showInfoMenu = true
      if (menu === 'char') this.showCharMenu = true
      if (menu === 'admin') this.showAdminMenu = true
    },
    closeMenu(menu) {
      // Delay closing to allow mouse to move to submenu
      this.closeTimeout = setTimeout(() => {
        if (menu === 'info') this.showInfoMenu = false
        if (menu === 'char') this.showCharMenu = false
        if (menu === 'admin') this.showAdminMenu = false
      }, 200)
    },
    closeAllMenus() {
      this.showInfoMenu = false
      this.showCharMenu = false
      this.showAdminMenu = false
    },
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
