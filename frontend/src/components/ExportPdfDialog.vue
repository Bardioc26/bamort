<template>
  <div v-if="showDialog" class="modal-overlay" @click.self="closeDialog">
    <div class="modal-content">
      <div class="modal-header">
        <h3>{{ $t('export.exportPDF') }}</h3>
        <button @click="closeDialog" class="close-button">&times;</button>
      </div>
      <div class="modal-body">
        <div v-if="isExporting" class="loading-overlay">
          <div class="spinner"></div>
          <p>{{ $t('export.generating') }}</p>
        </div>
        <div class="form-group">
          <label>{{ $t('export.selectTemplate') }}:</label>
          <select v-model="selectedTemplate" class="template-select" :disabled="isExporting">
            <option value="">{{ $t('export.pleaseSelectTemplate') }}</option>
            <option v-for="template in templates" :key="template.id" :value="template.id">
              {{ template.name }}
            </option>
          </select>
        </div>
        <div class="form-group">
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
        <button @click="exportToPDF" class="btn-export" :disabled="!selectedTemplate || isExporting">
          <span v-if="!isExporting">{{ $t('export.export') }}</span>
          <span v-else>{{ $t('export.exporting') }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #dee2e6;
}

.modal-header h3 {
  margin: 0;
  color: #333;
  font-size: 1.25rem;
}

.close-button {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #999;
  cursor: pointer;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-button:hover {
  color: #333;
}

.modal-body {
  padding: 20px;
  position: relative;
}

.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.95);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  z-index: 10;
  border-radius: 0 0 8px 8px;
}

.spinner {
  border: 4px solid #f3f3f3;
  border-top: 4px solid #007bff;
  border-radius: 50%;
  width: 50px;
  height: 50px;
  animation: spin 1s linear infinite;
  margin-bottom: 15px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-overlay p {
  color: #007bff;
  font-weight: 500;
  margin: 0;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #495057;
}

.template-select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  background: white;
  color: #495057;
  font-size: 0.95rem;
  cursor: pointer;
}

.template-select:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}

.template-select:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  background: #e9ecef;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  user-select: none;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 20px;
  border-top: 1px solid #dee2e6;
}

.btn-cancel {
  padding: 10px 20px;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  background: #f8f9fa;
  color: #495057;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s ease;
}

.btn-cancel:hover {
  background: #e9ecef;
  border-color: #adb5bd;
}

.btn-export {
  padding: 10px 20px;
  border: 1px solid #007bff;
  border-radius: 6px;
  background: #007bff;
  color: white;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s ease;
}

.btn-export:hover:not(:disabled) {
  background: #0056b3;
  border-color: #0056b3;
}

.btn-export:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>

<script>
import API from '../utils/api'

export default {
  name: "ExportPdfDialog",
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
      selectedTemplate: "",
      showUserName: false,
      isExporting: false
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
    
    closeDialog() {
      this.$emit('update:showDialog', false)
    }
  }
}
</script>
