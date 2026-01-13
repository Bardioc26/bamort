<template>
  <div class="landing-page">
    <div class="landing-content">
      <div class="dragon-container">
        <img src="/bamorty.png" alt="Bamort Dragon" class="dragon-image" />
      </div>
      
      <div class="info-container">
        <h1>{{ $t('landing.title') }}</h1>
        <p class="description">{{ $t('landing.description') }}</p>
        
        <div class="version-info">
          <p>{{ $t('landing.frontendVersion') }}: {{ frontendVersion }}<!-- ({{ frontendCommit }})--> </p>
          <p>{{ $t('landing.backendVersion') }}: {{ backendVersion }}<!-- ({{ backendCommit }})-->  </p>
        </div>
        
        <div class="action-links">
          <router-link to="/login" class="btn btn-primary" :class="{ disabled: !isBackendAvailable }" :event="isBackendAvailable ? 'click' : ''">
            {{ $t('landing.login') }}
          </router-link>
          <a :href="githubUrl" target="_blank" rel="noopener noreferrer" class="btn btn-secondary">
            {{ $t('landing.github') }}
          </a>
        </div>

        <div class="quick-links">
          <router-link to="/help" class="quick-link">
            {{ $t('landing.help') }}
          </router-link>
          <router-link to="/sponsors" class="quick-link">
            {{ $t('landing.sponsors') }}
          </router-link>
          <router-link to="/system-info" class="quick-link">
            {{ $t('landing.systemInfo') }}
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Component-specific styles moved to main.css as per project conventions */
</style>

<script>
import axios from 'axios'
import { getVersion, getGitCommit } from '../version'

export default {
  name: "LandingView",
  data() {
    return {
      frontendVersion: getVersion(),
      frontendCommit: getGitCommit(),
      backendVersion: "Loading...",
      backendCommit: "Loading...",
      githubUrl: "https://github.com/Bardioc26/bamort",
      retryCount: 0,
      maxRetries: 24,
      retryInterval: null
    }
  },
  mounted() {
    this.fetchBackendVersion()
  },
  beforeUnmount() {
    if (this.retryInterval) {
      clearInterval(this.retryInterval)
    }
  },
  computed: {
    isBackendAvailable() {
      return this.backendVersion !== "Loading..." && 
             this.backendVersion !== "Unavailable" && 
             this.backendVersion !== "Unreachable" &&
             this.backendVersion !== "Unknown"
    }
  },
  methods: {
    async fetchBackendVersion() {
      try {
        const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8180'
        const response = await axios.get(`${apiUrl}/api/public/version`)
        
        if (response.data) {
          this.backendVersion = response.data.version || "Unknown"
          this.backendCommit = response.data.gitCommit || "Unknown"
          if (this.retryInterval) {
            clearInterval(this.retryInterval)
            this.retryInterval = null
          }
        }
      } catch (error) {
        console.warn("Could not fetch backend version:", error)
        this.backendVersion = "Unavailable"
        this.backendCommit = "N/A"
        
        if (this.retryCount < this.maxRetries && !this.retryInterval) {
          this.retryInterval = setInterval(() => {
            this.retryCount++
            if (this.retryCount >= this.maxRetries) {
              clearInterval(this.retryInterval)
              this.retryInterval = null
              console.warn("Max retries reached for backend version")
              this.backendVersion = "Unreachable"
              return
            }
            this.fetchBackendVersion()
          }, 5000)
        }
      }
    }
  }
}
</script>
