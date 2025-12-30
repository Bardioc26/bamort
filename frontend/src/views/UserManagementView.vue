<template>
  <div class="page">
    <div class="page-header">
      <h2>{{ $t('userManagement.title') }}</h2>
    </div>

    <div v-if="isLoading" class="loading">{{ $t('userManagement.loading') }}</div>

    <div v-else-if="error" class="badge badge-danger">{{ error }}</div>

    <div v-else class="card">
      <table class="data-table">
        <thead>
          <tr>
            <th>{{ $t('userManagement.id') }}</th>
            <th>{{ $t('userManagement.username') }}</th>
            <th>{{ $t('userManagement.email') }}</th>
            <th>{{ $t('userManagement.role') }}</th>
            <th>{{ $t('userManagement.createdAt') }}</th>
            <th>{{ $t('userManagement.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.id }}</td>
            <td>{{ user.username }}</td>
            <td>{{ user.email }}</td>
            <td>
              <span :class="getRoleBadgeClass(user.role)">
                {{ $t(`userManagement.roles.${user.role}`) }}
              </span>
            </td>
            <td>{{ formatDate(user.created_at) }}</td>
            <td>
              <button 
                @click="openRoleDialog(user)" 
                class="btn btn-secondary btn-sm"
                :disabled="user.id === currentUser.id"
              >
                {{ $t('userManagement.changeRole') }}
              </button>
              <button 
                @click="openPasswordDialog(user)" 
                class="btn btn-sm"
              >
                {{ $t('userManagement.changePassword') }}
              </button>
              <button 
                @click="confirmDeleteUser(user)" 
                class="btn btn-sm"
                :disabled="user.id === currentUser.id"
              >
                {{ $t('userManagement.delete') }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Role Change Dialog -->
    <div v-if="showRoleDialog" class="modal-overlay" @click.self="showRoleDialog = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ $t('userManagement.changeRoleTitle') }}</h3>
        </div>
        <div class="modal-body">
          <p>{{ $t('userManagement.changeRoleFor') }}: <strong>{{ selectedUser.username }}</strong></p>
          <div class="form-group">
            <label>{{ $t('userManagement.selectRole') }}</label>
            <select v-model="newRole" class="form-control">
              <option value="standard">{{ $t('userManagement.roles.standard') }}</option>
              <option value="maintainer">{{ $t('userManagement.roles.maintainer') }}</option>
              <option value="admin">{{ $t('userManagement.roles.admin') }}</option>
            </select>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="updateUserRole" class="btn btn-primary">
            {{ $t('userManagement.save') }}
          </button>
          <button @click="showRoleDialog = false" class="btn btn-secondary">
            {{ $t('userManagement.cancel') }}
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Dialog -->
    <div v-if="showDeleteDialog" class="modal-overlay" @click.self="showDeleteDialog = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ $t('userManagement.deleteUserTitle') }}</h3>
        </div>
        <div class="modal-body">
          <p>{{ $t('userManagement.deleteConfirm') }}: <strong>{{ selectedUser.username }}</strong>?</p>
          <p class="badge badge-warning">{{ $t('userManagement.deleteWarning') }}</p>
        </div>
        <div class="modal-footer">
          <button @click="deleteUser" class="btn btn-danger">
            {{ $t('userManagement.delete') }}
          </button>
          <button @click="showDeleteDialog = false" class="btn btn-secondary">
            {{ $t('userManagement.cancel') }}
          </button>
        </div>
      </div>
    </div>

    <!-- Change Password Dialog -->
    <div v-if="showPasswordDialog" class="modal-overlay" @click.self="showPasswordDialog = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ $t('userManagement.changePasswordTitle') }}</h3>
        </div>
        <div class="modal-body">
          <p>{{ $t('userManagement.changePasswordFor') }}: <strong>{{ selectedUser.username }}</strong></p>
          <div class="form-group">
            <label>{{ $t('userManagement.newPassword') }}</label>
            <input 
              v-model="newPassword" 
              type="password" 
              class="form-control" 
              :placeholder="$t('userManagement.newPasswordPlaceholder')"
              minlength="6"
            />
          </div>
          <div class="form-group">
            <label>{{ $t('userManagement.confirmPassword') }}</label>
            <input 
              v-model="confirmPassword" 
              type="password" 
              class="form-control" 
              :placeholder="$t('userManagement.confirmPasswordPlaceholder')"
              minlength="6"
            />
          </div>
        </div>
        <div class="modal-footer">
          <button @click="changeUserPassword" class="btn btn-primary">
            {{ $t('userManagement.save') }}
          </button>
          <button @click="showPasswordDialog = false" class="btn btn-secondary">
            {{ $t('userManagement.cancel') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #dee2e6;
}

.data-table th {
  background-color: #f8f9fa;
  font-weight: 600;
}

.data-table tr:hover {
  background-color: #f8f9fa;
}

.btn-sm {
  padding: 5px 10px;
  font-size: 0.875rem;
  margin-right: 5px;
}

.badge-role-standard {
  background-color: #6c757d;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.875rem;
}

.badge-role-maintainer {
  background-color: #0dcaf0;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.875rem;
}

.badge-role-admin {
  background-color: #dc3545;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.875rem;
}

.loading {
  text-align: center;
  padding: 20px;
  color: #6c757d;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 0;
  border-radius: 8px;
  max-width: 500px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  padding: 20px;
  border-bottom: 1px solid #dee2e6;
}

.modal-header h3 {
  margin: 0;
}

.modal-body {
  padding: 20px;
}

.modal-footer {
  padding: 15px 20px;
  border-top: 1px solid #dee2e6;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>

<script>
import API from '../utils/api'
import { useUserStore } from '../stores/userStore'

export default {
  name: 'UserManagementView',
  data() {
    return {
      users: [],
      isLoading: false,
      error: null,
      showRoleDialog: false,
      showDeleteDialog: false,
      showPasswordDialog: false,
      selectedUser: null,
      newRole: '',
      newPassword: '',
      confirmPassword: ''
    }
  },
  computed: {
    currentUser() {
      const userStore = useUserStore()
      return userStore.currentUser
    }
  },
  async created() {
    await this.loadUsers()
  },
  methods: {
    async loadUsers() {
      this.isLoading = true
      this.error = null
      try {
        const response = await API.get('/api/users')
        this.users = response.data
      } catch (error) {
        console.error('Failed to load users:', error)
        this.error = this.$t('userManagement.loadError')
      } finally {
        this.isLoading = false
      }
    },
    openRoleDialog(user) {
      this.selectedUser = user
      this.newRole = user.role
      this.showRoleDialog = true
    },
    async updateUserRole() {
      try {
        await API.put(`/api/users/${this.selectedUser.id}/role`, {
          role: this.newRole
        })
        this.showRoleDialog = false
        await this.loadUsers()
      } catch (error) {
        console.error('Failed to update user role:', error)
        this.error = this.$t('userManagement.updateError')
      }
    },
    confirmDeleteUser(user) {
      this.selectedUser = user
      this.showDeleteDialog = true
    },
    async deleteUser() {
      try {
        await API.delete(`/api/users/${this.selectedUser.id}`)
        this.showDeleteDialog = false
        await this.loadUsers()
      } catch (error) {
        console.error('Failed to delete user:', error)
        this.error = this.$t('userManagement.deleteError')
      }
    },
    openPasswordDialog(user) {
      this.selectedUser = user
      this.newPassword = ''
      this.confirmPassword = ''
      this.showPasswordDialog = true
    },
    async changeUserPassword() {
      if (!this.newPassword || !this.confirmPassword) {
        this.error = this.$t('userManagement.passwordRequired')
        return
      }
      if (this.newPassword.length < 6) {
        this.error = this.$t('userManagement.passwordTooShort')
        return
      }
      if (this.newPassword !== this.confirmPassword) {
        this.error = this.$t('userManagement.passwordMismatch')
        return
      }
      try {
        await API.put(`/api/users/${this.selectedUser.id}/password`, {
          new_password: this.newPassword
        })
        this.showPasswordDialog = false
        this.error = null
      } catch (error) {
        console.error('Failed to change password:', error)
        this.error = this.$t('userManagement.passwordChangeError')
      }
    },
    getRoleBadgeClass(role) {
      return `badge-role-${role}`
    },
    formatDate(dateString) {
      const date = new Date(dateString)
      return date.toLocaleDateString() + ' ' + date.toLocaleTimeString()
    }
  }
}
</script>
