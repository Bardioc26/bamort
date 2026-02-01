<template>
  <div v-if="showDialog" class="modal-overlay" @click.self="closeDialog">
    <div class="modal-content">
      <div class="modal-header">
        <h3>{{ $t('visibility.title') }}</h3>
        <button @click="closeDialog" class="close-button">&times;</button>
      </div>
      <div class="modal-body">
        <div v-if="isUpdating" class="loading-overlay">
          <div class="spinner"></div>
          <p>{{ $t('visibility.updating') }}</p>
        </div>
        <div class="form-group">
          <p>{{ $t('visibility.description') }}</p>
          <div class="visibility-options">
            <label class="radio-option">
              <input 
                type="radio" 
                :value="false" 
                v-model="isPublic"
                :disabled="isUpdating"
              >
              <div class="option-content">
                <span class="option-label">{{ $t('visibility.private') }}</span>
                <span class="option-description">{{ $t('visibility.privateDescription') }}</span>
              </div>
            </label>
            <label class="radio-option">
              <input 
                type="radio" 
                :value="true" 
                v-model="isPublic"
                :disabled="isUpdating"
              >
              <div class="option-content">
                <span class="option-label">{{ $t('visibility.public') }}</span>
                <span class="option-description">{{ $t('visibility.publicDescription') }}</span>
              </div>
            </label>
          </div>
        </div>
      </div>
      <div class="modal-footer">
        <button @click="closeDialog" class="btn-cancel" :disabled="isUpdating">
          {{ $t('visibility.cancel') }}
        </button>
        <button @click="updateVisibility" class="btn-primary" :disabled="isUpdating">
          <span v-if="!isUpdating">{{ $t('visibility.save') }}</span>
          <span v-else>{{ $t('visibility.saving') }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import API from '../utils/api'

export default {
  name: "VisibilityDialog",
  props: {
    characterId: {
      type: [String, Number],
      required: true
    },
    currentVisibility: {
      type: Boolean,
      default: false
    },
    showDialog: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      isPublic: false,
      isUpdating: false
    }
  },
  watch: {
    currentVisibility: {
      immediate: true,
      handler(newValue) {
        this.isPublic = newValue
      }
    },
    showDialog(newValue) {
      if (newValue) {
        this.isPublic = this.currentVisibility
      }
    }
  },
  methods: {
    closeDialog() {
      if (!this.isUpdating) {
        this.$emit('update:showDialog', false)
      }
    },
    
    async updateVisibility() {
      this.isUpdating = true
      
      try {
        const token = localStorage.getItem('token')
        await API.patch(`/api/characters/${this.characterId}`, 
          { public: this.isPublic },
          {
            headers: { Authorization: `Bearer ${token}` }
          }
        )
        
        this.$emit('visibility-updated', this.isPublic)
        this.closeDialog()
      } catch (error) {
        console.error('Failed to update character visibility:', error)
        alert(this.$t('visibility.updateError') + ': ' + (error.response?.data?.error || error.message))
      } finally {
        this.isUpdating = false
      }
    }
  }
}
</script>

<style scoped>
.visibility-options {
  display: flex;
  flex-direction: column;
  gap: 15px;
  margin-top: 15px;
}

.radio-option {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 15px;
  border: 2px solid #dee2e6;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.radio-option:hover {
  border-color: #007bff;
  background-color: #f8f9fa;
}

.radio-option input[type="radio"] {
  margin-top: 3px;
  cursor: pointer;
}

.option-content {
  display: flex;
  flex-direction: column;
  gap: 5px;
  flex: 1;
}

.option-label {
  font-weight: 600;
  font-size: 1.1rem;
  color: #333;
}

.option-description {
  font-size: 0.9rem;
  color: #666;
}

.radio-option input[type="radio"]:checked + .option-content .option-label {
  color: #007bff;
}

.radio-option:has(input[type="radio"]:checked) {
  border-color: #007bff;
  background-color: #e7f3ff;
}
</style>
