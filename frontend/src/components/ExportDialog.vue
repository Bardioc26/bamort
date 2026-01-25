<template>
  <div v-if="showDialog" class="modal-overlay" @click.self="closeDialog">
    <div class="modal-content">
      <div class="modal-header">
        <h3>{{ $t('export.title') }}</h3>
        <button @click="closeDialog" class="close-button">&times;</button>
      </div>
      <div class="modal-body">
        <div v-if="isExporting" class="loading-overlay">
          <div class="spinner"></div>
          <p>{{ $t('export.generating') }}</p>
        </div>
        <div class="form-group">
          <label>{{ $t('export.selectFormat') }}:</label>
          <select v-model="selectedFormat" class="template-select" :disabled="isExporting">
            <option value="">{{ $t('export.pleaseSelectFormat') }}</option>
            <option value="pdf">{{ $t('export.formatPDF') }}</option>
            <option value="vtt">{{ $t('export.formatVTT') }}</option>
            <option value="bamort">{{ $t('export.formatBaMoRT') }}</option>
          </select>
        </div>
        <div v-if="selectedFormat === 'pdf'" class="form-group">
          <label>{{ $t('export.selectTemplate') }}:</label>
          <select v-model="selectedTemplate" class="template-select" :disabled="isExporting">
            <option value="">{{ $t('export.pleaseSelectTemplate') }}</option>
            <option v-for="template in templates" :key="template.id" :value="template.id">
              {{ template.name }}
            </option>
          </select>
        </div>
        <div v-if="selectedFormat === 'pdf'" class="form-group">
          <label class="checkbox-label">
            <input type="checkbox" v-model="showUserName" :disabled="isExporting">
            {{ $t('export.showUserName') }}
          </label>
        </div>
      </div>
      <div class="modal-footer">
        <button @click="closeDialog" class="btn-cancel" :disabled="isExporting">
          {{ $t('export.cancel') }}
        </button>
        <button @click="performExport" class="btn-export" :disabled="!canExport || isExporting">
          <span v-if="!isExporting">{{ $t('export.export') }}</span>
          <span v-else>{{ $t('export.exporting') }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<style>
/* All common styles moved to main.css - no component-specific styles needed */
</style>

<script>
import API from '../utils/api'

export default {
  name: "ExportDialog",
  props: {
    characterId: {
      type: [String, Number],
      required: true
    },
    showDialog: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      templates: [],
      selectedFormat: "",
      selectedTemplate: "",
      showUserName: false,
      isExporting: false
    }
  },
  computed: {
    canExport() {
      if (!this.selectedFormat) return false
      if (this.selectedFormat === 'pdf' && !this.selectedTemplate) return false
      return true
    }
  },
  async created() {
    await this.loadTemplates()
  },
  methods: {
    async loadTemplates() {
      try {
        const response = await API.get('/api/pdf/templates')
        this.templates = response.data
        // Auto-select first template if available
        if (this.templates.length > 0) {
          this.selectedTemplate = this.templates[0].id
        }
      } catch (error) {
        console.error('Failed to load templates:', error)
      }
    },
    
    async performExport() {
      if (!this.selectedFormat) {
        alert(this.$t('export.pleaseSelectFormat'))
        return
      }

      if (this.selectedFormat === 'pdf') {
        await this.exportToPDF()
      } else if (this.selectedFormat === 'vtt') {
        await this.exportToVTT()
      } else if (this.selectedFormat === 'bamort') {
        await this.exportToBaMoRT()
      }
    },
    
    async exportToPDF() {
      if (!this.selectedTemplate) {
        alert(this.$t('export.pleaseSelectTemplate'))
        return
      }
      
      this.isExporting = true
      
      try {
        // Build URL parameters
        const params = new URLSearchParams({
          template: this.selectedTemplate
        })
        if (this.showUserName) {
          params.append('showUserName', 'true')
        }
        
        // Get filename from export API (saves PDF to file)
        const response = await API.get(`/api/pdf/export/${this.characterId}`, {
          params: Object.fromEntries(params)
        })
        
        const filename = response.data.filename
        if (!filename) {
          throw new Error('No filename returned from export')
        }
        
        // Open PDF in new window using file endpoint
        const pdfUrl = `${API.defaults.baseURL}/api/pdf/file/${filename}`
        window.open(pdfUrl, '_blank')
        
        // Emit success event and close dialog
        this.$emit('export-success')
        this.closeDialog()
      } catch (error) {
        console.error('Failed to export PDF:', error)
        alert(this.$t('export.exportFailed') + ': ' + (error.response?.data?.error || error.message))
      } finally {
        this.isExporting = false
      }
    },

    async exportToVTT() {
      this.isExporting = true
      
      try {
        // Get VTT data and trigger download
        const response = await API.get(`/api/importer/export/vtt/${this.characterId}/file`, {
          responseType: 'blob'
        })
        
        // Create download link
        const blob = new Blob([response.data], { type: 'application/json' })
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.download = `character_${this.characterId}_vtt.json`
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        window.URL.revokeObjectURL(url)
        
        this.$emit('export-success')
        this.closeDialog()
      } catch (error) {
        console.error('Failed to export VTT:', error)
        alert(this.$t('export.exportFailed') + ': ' + (error.response?.data?.error || error.message))
      } finally {
        this.isExporting = false
      }
    },

    async exportToBaMoRT() {
      this.isExporting = true
      
      try {
        // Get BaMoRT JSON data and trigger download
        const response = await API.get(`/api/transfer/download/${this.characterId}`, {
          responseType: 'blob'
        })
        
        // Create download link
        const blob = new Blob([response.data], { type: 'application/json' })
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.download = `character_${this.characterId}_bamort.json`
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        window.URL.revokeObjectURL(url)
        
        this.$emit('export-success')
        this.closeDialog()
      } catch (error) {
        console.error('Failed to export BaMoRT format:', error)
        alert(this.$t('export.exportFailed') + ': ' + (error.response?.data?.error || error.message))
      } finally {
        this.isExporting = false
      }
    },
    
    closeDialog() {
      this.$emit('update:showDialog', false)
    }
  }
}
</script>
