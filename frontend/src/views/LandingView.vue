<template>
  <div class="landing-page">
    <div class="landing-content">
      <div class="dragon-container">
        <img src="/Drache.png" alt="Bamort Dragon" class="dragon-image" />
      </div>
      
      <div class="info-container">
        <h1>{{ $t('landing.title') }}</h1>
        <p class="description">{{ $t('landing.description') }}</p>
        
        <div class="version-info">
          <p>{{ $t('landing.frontendVersion') }}: {{ frontendVersion }}<!-- ({{ frontendCommit }})--> </p>
          <p>{{ $t('landing.backendVersion') }}: {{ backendVersion }}<!-- ({{ backendCommit }})-->  </p>
        </div>
        
        <div class="action-links">
          <router-link to="/login" class="btn btn-primary">
            {{ $t('landing.login') }}
          </router-link>
          <a :href="githubUrl" target="_blank" rel="noopener noreferrer" class="btn btn-secondary">
            {{ $t('landing.github') }}
          </a>
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
      githubUrl: "https://github.com/Bardioc26/bamort"
    }
  },
  mounted() {
    this.fetchBackendVersion()
  },
  methods: {
    async fetchBackendVersion() {
      try {
        const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8180'
        const response = await axios.get(`${apiUrl}/api/public/version`)
        
        if (response.data) {
          this.backendVersion = response.data.version || "Unknown"
          this.backendCommit = response.data.gitCommit || "Unknown"
        }
      } catch (error) {
        console.warn("Could not fetch backend version:", error)
        this.backendVersion = "Unavailable"
        this.backendCommit = "N/A"
      }
    }
  }
}
</script>
