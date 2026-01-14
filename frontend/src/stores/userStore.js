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
        this.currentUser = response.data
        
        // Set user's preferred language
        if (response.data.preferred_language) {
          i18n.global.locale.value = response.data.preferred_language
          localStorage.setItem('language', response.data.preferred_language)
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
