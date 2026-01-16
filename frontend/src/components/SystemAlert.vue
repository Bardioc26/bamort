<template>
  <div v-if="showAlert" class="system-alert" :class="alertType">
    <div class="alert-content">
      <span class="alert-icon">{{ alertIcon }}</span>
      <div class="alert-message">
        <div class="alert-text">{{ $t(messageKey) }}</div>
        <div v-if="versionInfo" class="version-info">
          {{ $t('system.backendVersion') }}: {{ backendVersion }} | 
          {{ $t('system.databaseVersion') }}: {{ dbVersion }}
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.system-alert {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 9999;
  padding: 1rem;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  animation: slideDown 0.3s ease-out;
}

@keyframes slideDown {
  from {
    transform: translateY(-100%);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.system-alert.warning {
  background-color: #fff3cd;
  border-bottom: 3px solid #ffc107;
  color: #856404;
}

.system-alert.success {
  background-color: #d4edda;
  border-bottom: 3px solid #28a745;
  color: #155724;
}

.system-alert.error {
  background-color: #f8d7da;
  border-bottom: 3px solid #dc3545;
  color: #721c24;
}

.alert-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  max-width: 1200px;
  margin: 0 auto;
}

.alert-icon {
  font-size: 1.5rem;
  font-weight: bold;
}

.alert-message {
  text-align: left;
}

.alert-text {
  font-weight: 500;
  margin-bottom: 0.25rem;
}

.version-info {
  font-size: 0.875rem;
  opacity: 0.8;
}
</style>

<script>
import API from '@/utils/api'

export default {
  name: 'SystemAlert',
  data() {
    return {
      showAlert: false,
      alertType: 'warning',
      messageKey: '',
      versionInfo: false,
      backendVersion: '',
      dbVersion: '',
      pollInterval: null,
      lastCheckTime: null
    }
  },
  computed: {
    alertIcon() {
      switch (this.alertType) {
        case 'warning':
          return '⚠️'
        case 'success':
          return '✓'
        case 'error':
          return '✖'
        default:
          return 'ℹ'
      }
    }
  },
  async mounted() {
    this.$api = API
    await this.checkSystemHealth()
    this.startPolling()
  },
  beforeUnmount() {
    this.stopPolling()
  },
  methods: {
    async checkSystemHealth() {
      try {
        //const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8180'
        //const response = await axios.get(`${baseURL}/api/system/health`)
        const response = await this.$api.get('/api/characters/available-skills-new')
        const health = response.data

        this.backendVersion = health.actual_backend_version
        this.dbVersion = health.db_version

        if (health.migrations_pending) {
          this.showWarning(health)
        } else if (health.compatible) {
          this.hideAlert()
        } else {
          this.showError(health)
        }

        this.lastCheckTime = new Date()
      } catch (error) {
        console.error('Failed to check system health:', error)
      }
    },
    showWarning(health) {
      this.showAlert = true
      this.alertType = 'warning'
      this.messageKey = 'system.migrationRequired'
      this.versionInfo = true
    },
    showError(health) {
      this.showAlert = true
      this.alertType = 'error'
      this.messageKey = 'system.incompatibleVersions'
      this.versionInfo = true
    },
    hideAlert() {
      this.showAlert = false
    },
    startPolling() {
      this.pollInterval = setInterval(() => {
        this.checkSystemHealth()
      }, 30000)
    },
    stopPolling() {
      if (this.pollInterval) {
        clearInterval(this.pollInterval)
        this.pollInterval = null
      }
    }
  }
}
</script>
