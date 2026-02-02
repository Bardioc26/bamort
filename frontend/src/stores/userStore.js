import { defineStore } from 'pinia'
import API from '../utils/api'
import { i18n } from './languageStore'

export const useUserStore = defineStore('user', {
  state: () => ({
    currentUser: null,
    isLoading: false
  }),

  getters: {
    isAuthenticated: (state) => !!state.currentUser,
    userRole: (state) => state.currentUser?.role || 'standard',
    isAdmin: (state) => state.currentUser?.role === 'admin',
    isMaintainer: (state) => state.currentUser?.role === 'maintainer' || state.currentUser?.role === 'admin',
    isStandardUser: (state) => !!state.currentUser
  },

  actions: {
    async fetchCurrentUser() {
      this.isLoading = true
      try {
        const response = await API.get('/api/user/profile')
        const profile = { ...response.data }
        profile.display_name = profile.display_name || profile.username
        this.currentUser = profile
        
        // Set user's preferred language
        if (profile.preferred_language) {
          i18n.global.locale.value = profile.preferred_language
          localStorage.setItem('language', profile.preferred_language)
        }
      } catch (error) {
        console.error('Failed to fetch user profile:', error)
        this.currentUser = null
      } finally {
        this.isLoading = false
      }
    },

    clearUser() {
      this.currentUser = null
    }
  }
})
