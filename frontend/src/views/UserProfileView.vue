<template>
  <div class="user-profile">
    <div class="container">
      <h1>{{ $t('profile.title') }}</h1>

      <div v-if="loading" class="loading">
        {{ $t('profile.loading') }}
      </div>

      <div v-else class="profile-sections">
        <!-- User Information Section -->
        <div class="profile-section">
          <h2>{{ $t('profile.userInfo') }}</h2>
          <div class="info-row">
            <label>{{ $t('profile.username') }}:</label>
            <span>{{ userProfile.username }}</span>
          </div>
          <div class="info-row">
            <label>{{ $t('profile.currentEmail') }}:</label>
            <span>{{ userProfile.email }}</span>
          </div>
          <div class="info-row">
            <label>{{ $t('profile.role') }}:</label>
            <span :class="getRoleBadgeClass(userProfile.role)">
              {{ $t(`userManagement.roles.${userProfile.role}`) }}
            </span>
          </div>
        </div>

        <!-- Change Email Section -->
        <div class="profile-section">
          <h2>{{ $t('profile.changeEmail') }}</h2>
          <form @submit.prevent="updateEmail" class="profile-form">
            <div class="form-group">
              <label for="newEmail">{{ $t('profile.newEmail') }}:</label>
              <input 
                type="email" 
                id="newEmail" 
                v-model="emailForm.newEmail" 
                :placeholder="$t('profile.emailPlaceholder')"
                required
              />
            </div>
            <button type="submit" :disabled="isUpdating" class="btn-primary">
              <span v-if="!isUpdating">{{ $t('profile.updateEmail') }}</span>
              <span v-else>{{ $t('profile.updating') }}</span>
            </button>
          </form>
        </div>

        <!-- Change Password Section -->
        <div class="profile-section">
          <h2>{{ $t('profile.changePassword') }}</h2>
          <form @submit.prevent="updatePassword" class="profile-form">
            <div class="form-group">
              <label for="currentPassword">{{ $t('profile.currentPassword') }}:</label>
              <input 
                type="password" 
                id="currentPassword" 
                v-model="passwordForm.currentPassword" 
                :placeholder="$t('profile.currentPasswordPlaceholder')"
                required
              />
            </div>
            <div class="form-group">
              <label for="newPassword">{{ $t('profile.newPassword') }}:</label>
              <input 
                type="password" 
                id="newPassword" 
                v-model="passwordForm.newPassword" 
                :placeholder="$t('profile.newPasswordPlaceholder')"
                minlength="6"
                required
              />
            </div>
            <div class="form-group">
              <label for="confirmPassword">{{ $t('profile.confirmPassword') }}:</label>
              <input 
                type="password" 
                id="confirmPassword" 
                v-model="passwordForm.confirmPassword" 
                :placeholder="$t('profile.confirmPasswordPlaceholder')"
                minlength="6"
                required
              />
            </div>
            <button type="submit" :disabled="isUpdating" class="btn-primary">
              <span v-if="!isUpdating">{{ $t('profile.updatePassword') }}</span>
              <span v-else>{{ $t('profile.updating') }}</span>
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
/* All common styles moved to main.css */

.user-profile {
  padding: var(--padding-lg);
  margin-top: 2%;
}

.container {
  max-width: 800px;
  margin: 0 auto;
}

h1 {
  color: var(--color-primary);
  margin-bottom: var(--margin-lg);
  text-align: center;
}

.loading {
  text-align: center;
  padding: var(--padding-lg);
  color: var(--color-text-secondary);
}

.profile-sections {
  display: flex;
  flex-direction: column;
  gap: var(--margin-lg);
}

.profile-section {
  background-color: var(--color-bg-secondary);
  padding: var(--padding-md);
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow-light);
}

.profile-section h2 {
  color: var(--color-primary);
  margin-bottom: var(--margin-md);
  font-size: 1.2em;
}

.info-row {
  display: flex;
  padding: var(--padding-sm) 0;
  border-bottom: 1px solid var(--color-border);
}

.info-row:last-child {
  border-bottom: none;
}

.info-row label {
  font-weight: bold;
  width: 200px;
  color: var(--color-text-secondary);
}

.info-row span {
  flex: 1;
  color: var(--color-text-primary);
}

.profile-form {
  display: flex;
  flex-direction: column;
  gap: var(--margin-md);
}

.badge-role-standard {
  background-color: #6c757d;
  color: white;
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 600;
}

.badge-role-maintainer {
  background-color: #0dcaf0;
  color: white;
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 600;
}

.badge-role-admin {
  background-color: #dc3545;
  color: white;
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 600;
}
</style>

<script>
import API from '../utils/api'

export default {
  name: 'UserProfileView',
  data() {
    return {
      loading: true,
      isUpdating: false,
      userProfile: {
        username: '',
        email: '',
        role: 'standard'
      },
      emailForm: {
        newEmail: ''
      },
      passwordForm: {
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
      }
    }
  },
  async created() {
    await this.loadProfile()
  },
  methods: {
    async loadProfile() {
      this.loading = true
      try {
        const response = await API.get('/api/user/profile')
        this.userProfile = response.data
        this.emailForm.newEmail = this.userProfile.email
      } catch (error) {
        console.error('Failed to load profile:', error)
        alert(this.$t('profile.loadError') + ': ' + (error.response?.data?.error || error.message))
      } finally {
        this.loading = false
      }
    },
    getRoleBadgeClass(role) {
      return `badge-role-${role}`
    },
    async updateEmail() {
      if (!this.emailForm.newEmail) {
        alert(this.$t('profile.emailRequired'))
        return
      }

      if (this.emailForm.newEmail === this.userProfile.email) {
        alert(this.$t('profile.emailUnchanged'))
        return
      }

      this.isUpdating = true
      try {
        const response = await API.put('/api/user/email', {
          email: this.emailForm.newEmail
        })
        
        this.userProfile.email = response.data.email
        alert(this.$t('profile.emailUpdateSuccess'))
      } catch (error) {
        console.error('Failed to update email:', error)
        let errorMsg = this.$t('profile.emailUpdateError')
        if (error.response?.data?.error) {
          if (error.response.data.error.includes('already in use')) {
            errorMsg = this.$t('profile.emailInUse')
          } else {
            errorMsg += ': ' + error.response.data.error
          }
        }
        alert(errorMsg)
      } finally {
        this.isUpdating = false
      }
    },
    async updatePassword() {
      if (!this.passwordForm.currentPassword || !this.passwordForm.newPassword || !this.passwordForm.confirmPassword) {
        alert(this.$t('profile.allFieldsRequired'))
        return
      }

      if (this.passwordForm.newPassword.length < 6) {
        alert(this.$t('profile.passwordTooShort'))
        return
      }

      if (this.passwordForm.newPassword !== this.passwordForm.confirmPassword) {
        alert(this.$t('profile.passwordMismatch'))
        return
      }

      this.isUpdating = true
      try {
        await API.put('/api/user/password', {
          current_password: this.passwordForm.currentPassword,
          new_password: this.passwordForm.newPassword
        })
        
        alert(this.$t('profile.passwordUpdateSuccess'))
        
        // Clear password fields
        this.passwordForm.currentPassword = ''
        this.passwordForm.newPassword = ''
        this.passwordForm.confirmPassword = ''
      } catch (error) {
        console.error('Failed to update password:', error)
        let errorMsg = this.$t('profile.passwordUpdateError')
        if (error.response?.data?.error) {
          if (error.response.data.error.includes('incorrect')) {
            errorMsg = this.$t('profile.currentPasswordIncorrect')
          } else {
            errorMsg += ': ' + error.response.data.error
          }
        }
        alert(errorMsg)
      } finally {
        this.isUpdating = false
      }
    }
  }
}
</script>
