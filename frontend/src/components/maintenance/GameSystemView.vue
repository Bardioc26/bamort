<template>
  <div class="header-section">
    <h2>{{ $t('maintenance') }} - {{ $t('gamesystem.title') }}</h2>
    <div class="search-box">
      <input v-model="searchTerm" type="text" :placeholder="$t('search')" />
    </div>
  </div>

  <div v-if="error" class="error-box">{{ error }}</div>

  <div class="cd-view">
    <div class="cd-list">
      <table class="cd-table">
        <thead>
          <tr>
            <th class="cd-table-header">{{ $t('gamesystem.id') }}</th>
            <th class="cd-table-header">{{ $t('gamesystem.code') }}</th>
            <th class="cd-table-header">{{ $t('gamesystem.name') }}</th>
            <th class="cd-table-header">{{ $t('gamesystem.description') }}</th>
            <th class="cd-table-header">{{ $t('gamesystem.active') }}</th>
            <th class="cd-table-header"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="isLoading">
            <td colspan="6">{{ $t('common.loading') }}</td>
          </tr>

          <template v-for="gs in filteredSystems" :key="gs.id">
            <tr v-if="editingId !== gs.id">
              <td>{{ gs.id }}</td>
              <td>{{ gs.code }}</td>
              <td>{{ gs.name }}</td>
              <td>{{ gs.description || '-' }}</td>
              <td><input type="checkbox" :checked="gs.is_active" disabled /></td>
              <td><button @click="startEdit(gs)">{{ $t('gamesystem.edit') }}</button></td>
            </tr>
            <tr v-else>
              <td>{{ gs.id }}</td>
              <td>{{ gs.code }}</td>
              <td colspan="4">
                <div class="edit-form">
                  <div class="edit-row">
                    <label>{{ $t('gamesystem.name') }}</label>
                    <input v-model="editedItem.name" />
                  </div>
                  <div class="edit-row">
                    <label>{{ $t('gamesystem.description') }}</label>
                    <input v-model="editedItem.description" />
                  </div>
                  <div class="edit-row">
                    <label>{{ $t('gamesystem.active') }}</label>
                    <input type="checkbox" v-model="editedItem.is_active" />
                  </div>
                  <div class="edit-actions">
                    <button class="btn-primary" :disabled="isSaving" @click="saveEdit">
                      <span v-if="!isSaving">{{ $t('gamesystem.save') }}</span>
                      <span v-else>{{ $t('gamesystem.saving') }}</span>
                    </button>
                    <button class="btn-cancel" :disabled="isSaving" @click="cancelEdit">
                      {{ $t('gamesystem.cancel') }}
                    </button>
                  </div>
                </div>
              </td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.error-box {
  margin: 10px 0;
  padding: 10px 12px;
  background: #ffe3e3;
  color: #8a1c1c;
  border: 1px solid #f5c2c2;
  border-radius: 6px;
}
.edit-form {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.edit-row {
  display: flex;
  gap: 10px;
  align-items: center;
}
.edit-actions {
  display: flex;
  gap: 10px;
}
</style>

<script>
import API from '../../utils/api'
import { loadGameSystems as fetchGameSystems } from '../../utils/maintenanceGameSystems'

export default {
  name: 'GameSystemView',
  data() {
    return {
      systems: [],
      editingId: null,
      editedItem: null,
      isLoading: false,
      isSaving: false,
      error: '',
      searchTerm: '',
    }
  },
  async created() {
    await this.loadSystems()
  },
  computed: {
    filteredSystems() {
      const term = this.searchTerm.trim().toLowerCase()
      const list = term
        ? this.systems.filter(gs =>
            (gs.name || '').toLowerCase().includes(term) ||
            (gs.code || '').toLowerCase().includes(term)
          )
        : this.systems
      return [...list].sort((a, b) => (a.code || '').localeCompare(b.code || ''))
    },
  },
  methods: {
    async loadSystems() {
      this.isLoading = true
      this.error = ''
      try {
        this.systems = await fetchGameSystems()
      } catch (err) {
        console.error('Failed to load game systems:', err)
        this.error = err.response?.data?.error || err.message
      } finally {
        this.isLoading = false
      }
    },
    startEdit(gs) {
      this.editingId = gs.id
      this.editedItem = { ...gs }
    },
    cancelEdit() {
      this.editingId = null
      this.editedItem = null
    },
    async saveEdit() {
      if (!this.editedItem) return
      const payload = {
        name: this.editedItem.name || '',
        description: this.editedItem.description || '',
        is_active: !!this.editedItem.is_active,
      }
      this.isSaving = true
      try {
        const resp = await API.put(`/api/maintenance/game-systems/${this.editingId}`, payload)
        const idx = this.systems.findIndex(s => s.id === this.editingId)
        if (idx !== -1) this.systems.splice(idx, 1, resp.data)
        this.cancelEdit()
      } catch (err) {
        console.error('Failed to save game system:', err)
        this.error = err.response?.data?.error || err.message
      } finally {
        this.isSaving = false
      }
    },
  },
}
</script>
