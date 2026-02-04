<template>
  <div v-if="showDialog" class="modal-overlay" @click.self="closeDialog">
    <div class="modal-content modal-large">
      <div class="modal-header">
        <h3>{{ $t('visibility.title') }}</h3>
        <button @click="closeDialog" class="close-button">&times;</button>
      </div>
      <div class="modal-body">
        <div v-if="isUpdating" class="loading-overlay">
          <div class="spinner"></div>
          <p>{{ $t('visibility.updating') }}</p>
        </div>
        
        <!-- Visibility Options -->
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

        <!-- Share with Specific Users Section -->
        <div class="form-group">
          <h4>{{ $t('visibility.shareWithUsers') }}</h4>
          <p class="section-description">{{ $t('visibility.shareDescription') }}</p>
          
          <div class="share-sections-container">
            
            <!-- Add Users Section -->
            <div class="add-users-section">
              <h5>{{ $t('visibility.addUsers') }}</h5>
              <div class="user-search">
                <input 
                  v-model="searchQuery" 
                  type="text" 
                  :placeholder="$t('visibility.searchUsers')"
                  class="form-control"
                  :disabled="isUpdating"
                />
              </div>
              
              <div v-if="isLoadingUsers" class="loading">{{ $t('visibility.loadingUsers') }}</div>
              
              <div v-else-if="filteredAvailableUsers.length > 0" class="available-users-list">
                <div 
                  v-for="user in filteredAvailableUsers" 
                  :key="user.user_id"
                  class="user-item"
                  @click="toggleUser(user.user_id)"
                >
                  <div class="user-info">
                    <span class="user-name">{{ user.display_name || user.username }}</span>
                    <span class="user-email">{{ user.email }}</span>
                  </div>
                </div>
              </div>
              
              <div v-else-if="!isLoadingUsers && availableUsers.length === 0" class="no-users">
                {{ $t('visibility.noOtherUsers') }}
              </div>
              
              <div v-else class="no-users">
                {{ $t('visibility.noMatchingUsers') }}
              </div>
            </div>

            <!-- Currently Shared Users -->
            <div class="shared-users-list">
              <h5>{{ $t('visibility.currentlySharedWith') }}</h5>
              <div v-if="sharedUserIds.length > 0" class="shared-users-items">
                <div 
                  v-for="userId in sharedUserIds" 
                  :key="userId" 
                  class="user-item shared-user"
                >
                  <div class="user-info">
                    <span class="user-name">{{ getUserName(userId) }}</span>
                    <span class="user-email">{{ getUserEmail(userId) }}</span>
                  </div>
                  <button @click="removeUser(userId)" class="remove-btn" :disabled="isUpdating">&times;</button>
                </div>
              </div>
              <div v-else class="no-users">
                {{ $t('visibility.noSharedUsers') }}
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="modal-footer">
        <button @click="updateVisibilityAndShares" class="btn-primary btn-save" :disabled="isUpdating">
          <span v-if="!isUpdating">{{ $t('common.save') }}</span>
          <span v-else>{{ $t('common.saving') }}</span>
        </button>
        <button @click="closeDialog" class="btn-cancel" :disabled="isUpdating">
          {{ $t('common.cancel') }}
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
      isUpdating: false,
      isLoadingUsers: false,
      availableUsers: [],
      sharedUserIds: [],
      searchQuery: ''
    }
  },
  computed: {
    filteredAvailableUsers() {
      let users = this.availableUsers.filter(user => !this.sharedUserIds.includes(user.user_id))
      
      if (!this.searchQuery) {
        return users
      }
      const query = this.searchQuery.toLowerCase()
      return users.filter(user => 
        (user.display_name && user.display_name.toLowerCase().includes(query)) ||
        user.username.toLowerCase().includes(query) ||
        user.email.toLowerCase().includes(query)
      )
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
        this.loadAvailableUsers()
        this.loadCurrentShares()
      }
    }
  },
  methods: {
    async loadAvailableUsers() {
      this.isLoadingUsers = true
      try {
        const token = localStorage.getItem('token')
        const response = await API.get(`/api/characters/${this.characterId}/available-users`, {
          headers: { Authorization: `Bearer ${token}` }
        })
        this.availableUsers = response.data || []
      } catch (error) {
        console.error('Failed to load available users:', error)
      } finally {
        this.isLoadingUsers = false
      }
    },
    
    async loadCurrentShares() {
      try {
        const token = localStorage.getItem('token')
        const response = await API.get(`/api/characters/${this.characterId}/shares`, {
          headers: { Authorization: `Bearer ${token}` }
        })
        this.sharedUserIds = (response.data || []).map(share => share.user_id)
      } catch (error) {
        console.error('Failed to load current shares:', error)
        this.sharedUserIds = []
      }
    },
    
    toggleUser(userId) {
      const index = this.sharedUserIds.indexOf(userId)
      if (index > -1) {
        this.sharedUserIds.splice(index, 1)
      } else {
        this.sharedUserIds.push(userId)
      }
    },
    
    removeUser(userId) {
      const index = this.sharedUserIds.indexOf(userId)
      if (index > -1) {
        this.sharedUserIds.splice(index, 1)
      }
    },
    
    getUserName(userId) {
      const user = this.availableUsers.find(u => u.user_id === userId)
      return user ? (user.display_name || user.username) : 'Unknown'
    },
    
    getUserEmail(userId) {
      const user = this.availableUsers.find(u => u.user_id === userId)
      return user ? user.email : ''
    },
    
    closeDialog() {
      if (!this.isUpdating) {
        this.searchQuery = ''
        this.$emit('update:showDialog', false)
      }
    },
    
    async updateVisibilityAndShares() {
      this.isUpdating = true
      
      try {
        const token = localStorage.getItem('token')
        
        // Update visibility
        await API.patch(`/api/characters/${this.characterId}`, 
          { public: this.isPublic },
          {
            headers: { Authorization: `Bearer ${token}` }
          }
        )
        
        // Update shares
        await API.put(`/api/characters/${this.characterId}/shares`,
          { user_ids: this.sharedUserIds },
          {
            headers: { Authorization: `Bearer ${token}` }
          }
        )
        
        this.$emit('visibility-updated', this.isPublic)
        this.closeDialog()
      } catch (error) {
        console.error('Failed to update character visibility/shares:', error)
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

.modal-large {
  max-width: 700px;
  max-height: 90vh;
  overflow-y: auto;
}

.section-description {
  color: #666;
  font-size: 0.9rem;
  margin-bottom: 15px;
}

.share-sections-container {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

.shared-users-list {
  flex: 1;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 8px;
  min-height: 400px;
}

.shared-users-list h5 {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 1rem;
}

.shared-users-items {
  border: 1px solid #dee2e6;
  border-radius: 8px;
  max-height: 300px;
  overflow-y: auto;
  background: white;
}

.user-item.shared-user {
  background: white;
  cursor: default;
}

.user-item.shared-user:hover {
  background: #fff5f5;
}

.user-item.shared-user .remove-btn {
  background: none;
  border: none;
  color: #dc3545;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 0 8px;
  line-height: 1;
  transition: color 0.2s ease;
  font-weight: bold;
  flex-shrink: 0;
}

.user-item.shared-user .remove-btn:hover {
  color: #a71d2a;
}

.user-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.user-chip {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  background: white;
  border: 1px solid #dee2e6;
  border-radius: 20px;
  font-size: 0.9rem;
}

.user-chip .remove-btn {
  background: none;
  border: none;
  color: #dc3545;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 0;
  line-height: 1;
  transition: color 0.2s ease;
}

.user-chip .remove-btn:hover {
  color: #a71d2a;
}

.add-users-section {
  flex: 1;
}

.add-users-section h5 {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 1rem;
}

.user-search {
  margin-bottom: 15px;
}

.available-users-list {
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid #dee2e6;
  border-radius: 8px;
}

.user-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 15px;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: background 0.2s ease;
}

.user-item:last-child {
  border-bottom: none;
}

.user-item:hover {
  background: #f8f9fa;
}

.user-item.selected {
  background: #e7f3ff;
  border-left: 3px solid #007bff;
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.user-name {
  font-weight: 600;
  color: #333;
}

.user-email {
  font-size: 0.85rem;
  color: #666;
}

.check-icon {
  color: #28a745;
  font-size: 1.2rem;
  font-weight: bold;
}

.no-users {
  text-align: center;
  padding: 40px 20px;
  color: #999;
  font-style: italic;
}

.loading {
  text-align: center;
  padding: 20px;
  color: #666;
}
</style>
