<template>
  <div class="header-section">
    <h2>{{ $t('maintenance') }} - {{ $t('litsource.title') }}</h2>
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
            <th class="cd-table-header">{{ $t('litsource.id') }}</th>
            <th class="cd-table-header">{{ $t('litsource.code') }}</th>
            <th class="cd-table-header">{{ $t('litsource.name') }}</th>
            <th class="cd-table-header">{{ $t('litsource.fullName') }}</th>
            <th class="cd-table-header">{{ $t('litsource.edition') }}</th>
            <th class="cd-table-header">{{ $t('litsource.publisher') }}</th>
            <th class="cd-table-header">{{ $t('litsource.year') }}</th>
            <th class="cd-table-header">{{ $t('litsource.active') }}</th>
            <th class="cd-table-header">{{ $t('litsource.core') }}</th>
            <th class="cd-table-header"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="isLoading">
            <td colspan="10">{{ $t('common.loading') }}</td>
          </tr>

          <template v-for="src in filteredSources" :key="src.id">
            <tr v-if="editingId !== src.id">
              <td>{{ src.id }}</td>
              <td>{{ src.code }}</td>
              <td>{{ src.name }}</td>
              <td>{{ src.full_name }}</td>
              <td>{{ src.edition }}</td>
              <td>{{ src.publisher }}</td>
              <td>{{ src.publish_year }}</td>
              <td><input type="checkbox" :checked="src.is_active" disabled /></td>
              <td><input type="checkbox" :checked="src.is_core" disabled /></td>
              <td><button @click="startEdit(src)">{{ $t('litsource.edit') }}</button></td>
            </tr>
            <tr v-else>
              <td>{{ src.id }}</td>
              <td>{{ src.code }}</td>
              <td colspan="8">
                <div class="edit-form">
                  <div class="edit-row">
                    <label>{{ $t('litsource.name') }}</label>
                    <input v-model="editedItem.name" />
                  </div>
                  <div class="edit-row">
                    <label>{{ $t('litsource.fullName') }}</label>
                    <input v-model="editedItem.full_name" />
                  </div>
                  <div class="edit-row">
                    <label>{{ $t('litsource.edition') }}</label>
                    <input v-model="editedItem.edition" />
                    <label class="inline-label">{{ $t('litsource.publisher') }}</label>
                    <input v-model="editedItem.publisher" />
                    <label class="inline-label">{{ $t('litsource.year') }}</label>
                    <input v-model.number="editedItem.publish_year" type="number" />
                  </div>
                  <div class="edit-row">
                    <label>{{ $t('litsource.description') }}</label>
                    <input v-model="editedItem.description" />
                  </div>
                  <div class="edit-row">
                    <label>{{ $t('litsource.active') }}</label>
                    <input type="checkbox" v-model="editedItem.is_active" />
                    <label class="inline-label">{{ $t('litsource.core') }}</label>
                    <input type="checkbox" v-model="editedItem.is_core" />
                  </div>
                  <div class="edit-actions">
                    <button class="btn-primary" :disabled="isSaving" @click="saveEdit">
                      <span v-if="!isSaving">{{ $t('litsource.save') }}</span>
                      <span v-else>{{ $t('litsource.saving') }}</span>
                    </button>
                    <button class="btn-cancel" :disabled="isSaving" @click="cancelEdit">
                      {{ $t('litsource.cancel') }}
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
  flex-wrap: wrap;
  align-items: center;
}
.edit-actions {
  display: flex;
  gap: 10px;
}
.inline-label {
  margin-left: 10px;
}
</style>

<script>
import API from '../../utils/api'

export default {
  name: 'LitSourceView',
  data() {
    return {
      gameSystems: [],
      currentGameSystem: null,
      sources: [],
      editingId: null,
      editedItem: null,
      isLoading: false,
      isSaving: false,
      error: '',
      searchTerm: '',
    }
  },
  async created() {
    await this.initialize()
  },
  computed: {
    filteredSources() {
      const term = this.searchTerm.trim().toLowerCase()
      const list = term
        ? this.sources.filter(src =>
            (src.name || '').toLowerCase().includes(term) ||
            (src.code || '').toLowerCase().includes(term)
          )
        : this.sources
      return [...list].sort((a, b) => (a.code || '').localeCompare(b.code || ''))
    },
  },
  methods: {
    async initialize() {
      this.error = ''
      await this.loadGameSystems()
      if (this.currentGameSystem) {
        await this.loadSources()
      }
    },
    async loadGameSystems() {
      try {
        const resp = await API.get('/api/maintenance/game-systems')
        const systems = (resp.data?.game_systems || []).map(this.normalizeSystem)
        this.gameSystems = systems
        const active = systems.find(s => s.is_active)
        this.currentGameSystem = active || systems[0] || null
      } catch (err) {
        console.error('Failed to load game systems:', err)
        this.error = err.response?.data?.error || err.message
      }
    },
    async loadSources() {
      this.isLoading = true
      this.error = ''
      try {
        const params = this.buildGameSystemParams()
        const resp = await API.get('/api/maintenance/gsm-lit-sources', { params })
        this.sources = resp.data?.sources || []
      } catch (err) {
        console.error('Failed to load sources:', err)
        this.error = err.response?.data?.error || err.message
      } finally {
        this.isLoading = false
      }
    },
    buildGameSystemParams() {
      if (!this.currentGameSystem) return {}
      return {
        game_system_id: this.currentGameSystem.id,
        game_system: this.currentGameSystem.name,
      }
    },
    normalizeSystem(gs) {
      return {
        id: gs.id ?? gs.ID,
        code: gs.code ?? gs.Code,
        name: gs.name ?? gs.Name,
        description: gs.description ?? gs.Description,
        is_active: gs.is_active ?? gs.IsActive,
      }
    },
    startEdit(src) {
      this.editingId = src.id
      this.editedItem = { ...src }
    },
    cancelEdit() {
      this.editingId = null
      this.editedItem = null
    },
    async saveEdit() {
      if (!this.editedItem) return
      const payload = {
        name: this.editedItem.name || '',
        full_name: this.editedItem.full_name || '',
        edition: this.editedItem.edition || '',
        publisher: this.editedItem.publisher || '',
        publish_year: this.editedItem.publish_year || 0,
        description: this.editedItem.description || '',
        is_active: !!this.editedItem.is_active,
        is_core: !!this.editedItem.is_core,
      }
      this.isSaving = true
      try {
        const params = this.buildGameSystemParams()
        const resp = await API.put(`/api/maintenance/gsm-lit-sources/${this.editingId}`, payload, { params })
        const idx = this.sources.findIndex(s => s.id === this.editingId)
        if (idx !== -1) this.sources.splice(idx, 1, resp.data)
        this.cancelEdit()
      } catch (err) {
        console.error('Failed to save source:', err)
        this.error = err.response?.data?.error || err.message
      } finally {
        this.isSaving = false
      }
    },
  },
}
</script>
