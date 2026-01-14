<template>
  <div class="fullwidth-page">
    <div class="page-header" style="flex-direction: row;">
      <button @click="$router.back()" class="btn btn-secondary back-button">
        ← {{ $t('common.back') }}
      </button>
      <h2 style="padding-top: 8px; padding-left: 10px;">{{ $t('systemInfo.title') }}</h2>
    </div>

    <div class="card">
      <p>{{ $t('systemInfo.introduction') }}</p>
    </div>

    <div class="section-header">
      <h3>{{ $t('systemInfo.versions') }}</h3>
    </div>

    <div class="grid-container grid-2-columns">
      <div class="card">
        <h4>{{ $t('systemInfo.frontend') }}</h4>
        <p><strong>{{ $t('systemInfo.version') }}:</strong> {{ frontendVersion }}</p>
        <p><strong>{{ $t('systemInfo.commit') }}:</strong> <code>{{ frontendCommit }}</code></p>
      </div>

      <div class="card">
        <h4>{{ $t('systemInfo.backend') }}</h4>
        <p><strong>{{ $t('systemInfo.version') }}:</strong> {{ backendVersion }}</p>
        <p><strong>{{ $t('systemInfo.commit') }}:</strong> <code>{{ backendCommit }}</code></p>
        <p><strong>{{ $t('systemInfo.status') }}:</strong> 
          <span :class="statusClass">{{ statusText }}</span>
        </p>
      </div>
    </div>

    <div class="section-header">
      <h3>{{ $t('systemInfo.technologies') }}</h3>
    </div>

    <div class="grid-container grid-3-columns">
      <div class="card">
        <h4>Frontend</h4>
        <ul>
          <li>Vue 3</li>
          <li>Vite</li>
          <li>Vue Router</li>
          <li>Pinia</li>
          <li>Axios</li>
          <li>Vue i18n</li>
        </ul>
      </div>

      <div class="card">
        <h4>Backend</h4>
        <ul>
          <li>Go 1.25</li>
          <li>Gin Framework</li>
          <li>GORM</li>
          <li>MariaDB</li>
          <li>JWT Auth</li>
          <li>Chromedp (PDF)</li>
        </ul>
      </div>

      <div class="card">
        <h4>{{ $t('systemInfo.infrastructure') }}</h4>
        <ul>
          <li>Docker</li>
          <li>Docker Compose</li>
          <li>Air (Hot Reload) in Entwicklungsumgebung</li>
          <!-- <li>phpMyAdmin</li> -->
        </ul>
      </div>
    </div>

    <div class="section-header">
      <h3>{{ $t('systemInfo.features') }}</h3>
    </div>

    <div class="card">
      <ul>
        <li>{{ $t('systemInfo.feature1') }}</li>
        <li>{{ $t('systemInfo.feature2') }}</li>
        <li>{{ $t('systemInfo.feature3') }}</li>
        <li>{{ $t('systemInfo.feature4') }}</li>
        <li>{{ $t('systemInfo.feature5') }}</li>
        <li>{{ $t('systemInfo.feature6') }}</li>
        <li>{{ $t('systemInfo.feature7') }}</li>
        <li>{{ $t('systemInfo.feature8') }}</li>
      </ul>
    </div>

    <div class="section-header">
      <h3>{{ $t('systemInfo.license') }}</h3>
    </div>

    <div class="card">
      <p>{{ $t('systemInfo.licenseText') }}</p>
      <div style="display: flex; gap: 10px; margin-top: 15px;">
        <a :href="githubUrl" target="_blank" rel="noopener noreferrer" class="btn btn-primary">
          {{ $t('systemInfo.github') }}
        </a>
        <a :href="koFiUrl" target="_blank" rel="noopener noreferrer" class="btn btn-primary">
          {{ $t('sponsors.koFi') }}
        </a>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import { getVersion, getGitCommit } from '../version'

export default {
  name: "SystemInfoView",
  data() {
    return {
      frontendVersion: getVersion(),
      frontendCommit: getGitCommit(),
      backendVersion: "Loading...",
      backendCommit: "Loading...",
      githubUrl: "https://github.com/Bardioc26/bamort",
      koFiUrl: "https://ko-fi.com/bardioc26",
    }
  },
  computed: {
    statusClass() {
      if (this.backendVersion === "Loading...") return "status-loading"
      if (this.isBackendAvailable) return "status-available"
      return "status-unavailable"
    },
    statusText() {
      if (this.backendVersion === "Loading...") return this.$t('systemInfo.statusLoading')
      if (this.isBackendAvailable) return this.$t('systemInfo.statusAvailable')
      return this.$t('systemInfo.statusUnavailable')
    },
    isBackendAvailable() {
      return this.backendVersion !== "Loading..." && 
             this.backendVersion !== "Unavailable" && 
             this.backendVersion !== "Unreachable" &&
             this.backendVersion !== "Unknown"
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
