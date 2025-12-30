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
          <p>{{ $t('landing.version') }}: {{ version }}</p>
          <p>{{ $t('landing.commit') }}: {{ commit }}</p>
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
export default {
  name: "LandingView",
  data() {
    return {
      version: "0.1.0",
      commit: "9775290",
      githubUrl: "https://github.com/Bardioc26/bamort"
    }
  },
  mounted() {
    // Fetch version and commit from API if available
    this.fetchVersionInfo()
  },
  methods: {
    async fetchVersionInfo() {
      try {
        // For now, using static values
        // In future, could fetch from backend API endpoint
        this.version = import.meta.env.VITE_APP_VERSION || "0.1.0"
        this.commit = import.meta.env.VITE_GIT_COMMIT || "9775290"
      } catch (error) {
        console.warn("Could not fetch version info:", error)
      }
    }
  }
}
</script>
