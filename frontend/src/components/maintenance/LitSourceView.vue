<template>
  <div class="header-section">
    <h2>{{ $t('maintenance') }} - {{ $t('litsource.title') }}</h2>
    <div class="search-box">
      <input v-model="searchTerm" type="text" :placeholder="$t('search')" />
    </div>
    <div class="search-box">
      <select v-model.number="selectedSystemId" @change="handleGameSystemChange">
        <option value="">{{ $t('gamesystem.title') }}</option>
        <option v-for="system in systemOptions" :key="system.id" :value="system.id">
          {{ system.label }}
        </option>
      </select>
    </div>
    <button class="btn-primary" @click="startCreate">{{ $t('newEntry') }}</button>
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

          <tr v-if="creatingNew">
            <td>New</td>
            <td colspan="9">
              <div class="edit-form">
                <div class="edit-row">
                  <label>{{ $t('litsource.code') }}</label>
                  <input v-model="newItem.code" />
                  <label class="inline-label">{{ $t('litsource.name') }}</label>
                  <input v-model="newItem.name" />
                </div>
                <div class="edit-row">
                  <label>{{ $t('litsource.fullName') }}</label>
                  <input v-model="newItem.full_name" />
                </div>
                <div class="edit-row">
                  <label>{{ $t('litsource.edition') }}</label>
                  <input v-model="newItem.edition" />
                  <label class="inline-label">{{ $t('litsource.publisher') }}</label>
                  <input v-model="newItem.publisher" />
                  <label class="inline-label">{{ $t('litsource.year') }}</label>
                  <input v-model.number="newItem.publish_year" type="number" />
                </div>
                <div class="edit-row">
                  <label>{{ $t('litsource.description') }}</label>
                  <input v-model="newItem.description" />
                </div>
                <div class="edit-row">
                  <label>{{ $t('litsource.active') }}</label>
                  <input type="checkbox" v-model="newItem.is_active" />
                  <label class="inline-label">{{ $t('litsource.core') }}</label>
                  <input type="checkbox" v-model="newItem.is_core" />
                </div>
                <div class="edit-row">
                  <label>{{ $t('gamesystem.title') }}</label>
                  <select v-model.number="createSelectedSystemId">
                    <option value="">-</option>
                    <option v-for="system in systemOptions" :key="system.id" :value="system.id">
                      {{ system.label }}
                    </option>
                  </select>
                </div>
                <div class="edit-actions">
                  <button class="btn-primary btn-save" :disabled="isSaving" @click="saveCreate">
                    <span v-if="!isSaving">{{ $t('common.save') }}</span>
                    <span v-else>{{ $t('common.saving') }}</span>
                  </button>
                  <button class="btn-cancel" :disabled="isSaving" @click="cancelCreate">
                    {{ $t('common.cancel') }}
                  </button>
                </div>
              </div>
            </td>
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
              <td><button @click="startEdit(src)">{{ $t('common.edit') }}</button></td>
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
                  <div class="edit-row">
                    <label>{{ $t('gamesystem.title') }}</label>
                    <select v-model.number="selectedSystemId">
                      <option value="">-</option>
                      <option v-for="system in systemOptions" :key="system.id" :value="system.id">
                        {{ system.label }}
                      </option>
                    </select>
                  </div>
                  <div class="edit-actions">
                    <button class="btn-primary btn-save" :disabled="isSaving" @click="saveEdit">
                      <span v-if="!isSaving">{{ $t('common.save') }}</span>
                      <span v-else>{{ $t('common.saving') }}</span>
                    </button>
                    <button class="btn-cancel" :disabled="isSaving" @click="cancelEdit">
                      {{ $t('common.cancel') }}
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
import {
  buildGameSystemParams,
  findSystemById,
  findSystemIdByCode,
  loadGameSystems as fetchGameSystems,
  buildSystemOptions,
} from '../../utils/maintenanceGameSystems'

export default {
  name: 'LitSourceView',
  data() {
    return {
      gameSystems: [],
      currentGameSystem: null,
      selectedSystemId: null,
      sources: [],
      editingId: null,
      editedItem: null,
      creatingNew: false,
      newItem: null,
      createSelectedSystemId: null,
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
    systemOptions() {
      return buildSystemOptions(this.gameSystems)
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
        const systems = await fetchGameSystems()
        this.gameSystems = systems
        const active = systems.find(s => s.is_active)
        this.currentGameSystem = active || systems[0] || null
        this.selectedSystemId = this.currentGameSystem ? this.currentGameSystem.id : null
      } catch (err) {
        console.error('Failed to load game systems:', err)
        this.error = err.response?.data?.error || err.message
      }
    },
    async loadSources() {
      this.isLoading = true
      this.error = ''
      try {
        const params = buildGameSystemParams(this.currentGameSystem)
        const resp = await API.get('/api/maintenance/gsm-lit-sources', { params })
        this.sources = resp.data?.sources || []
      } catch (err) {
        console.error('Failed to load sources:', err)
        this.error = err.response?.data?.error || err.message
      } finally {
        this.isLoading = false
      }
    },
    startEdit(src) {
      this.editingId = src.id
      this.editedItem = { ...src }
      this.selectedSystemId = src.game_system_id
        || this.findSystemIdByCode(src.game_system)
        || this.currentGameSystem?.id
        || null
    },
    cancelEdit() {
      this.editingId = null
      this.editedItem = null
      this.selectedSystemId = this.currentGameSystem ? this.currentGameSystem.id : null
    },
    handleGameSystemChange() {
      const target = this.findSystemById(this.selectedSystemId)
      this.currentGameSystem = target || this.currentGameSystem
      if (this.currentGameSystem) {
        this.loadSources()
      }
    },
    startCreate() {
      this.cancelEdit()
      const defaultSystem = this.currentGameSystem || this.gameSystems.find(s => s.is_active) || this.gameSystems[0] || null
      this.createSelectedSystemId = defaultSystem ? defaultSystem.id : null
      this.newItem = {
        code: '',
        name: '',
        full_name: '',
        edition: '',
        publisher: '',
        publish_year: 0,
        description: '',
        is_active: true,
        is_core: false,
      }
      this.creatingNew = true
    },
    cancelCreate() {
      this.creatingNew = false
      this.newItem = null
      this.createSelectedSystemId = this.currentGameSystem ? this.currentGameSystem.id : null
    },
    findSystemById(id) {
      return findSystemById(this.gameSystems, id)
    },
    findSystemIdByCode(code) {
      return findSystemIdByCode(this.gameSystems, code)
    },
    async saveEdit() {
      if (!this.editedItem) return
      const targetSystem = this.findSystemById(this.selectedSystemId) || this.currentGameSystem
      const payload = {
        name: this.editedItem.name || '',
        full_name: this.editedItem.full_name || '',
        edition: this.editedItem.edition || '',
        publisher: this.editedItem.publisher || '',
        publish_year: this.editedItem.publish_year || 0,
        description: this.editedItem.description || '',
        is_active: !!this.editedItem.is_active,
        is_core: !!this.editedItem.is_core,
        game_system_id: targetSystem ? targetSystem.id : null,
        game_system: targetSystem ? targetSystem.code : '',
      }
      this.isSaving = true
      try {
        const params = buildGameSystemParams(targetSystem)
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
    async saveCreate() {
      if (!this.newItem) return
      const targetSystem = this.findSystemById(this.createSelectedSystemId) || this.currentGameSystem
      const payload = {
        code: this.newItem.code || '',
        name: this.newItem.name || '',
        full_name: this.newItem.full_name || '',
        edition: this.newItem.edition || '',
        publisher: this.newItem.publisher || '',
        publish_year: this.newItem.publish_year || 0,
        description: this.newItem.description || '',
        is_active: !!this.newItem.is_active,
        is_core: !!this.newItem.is_core,
        game_system_id: targetSystem ? targetSystem.id : null,
        game_system: targetSystem ? targetSystem.code : '',
      }
      this.isSaving = true
      try {
        const params = buildGameSystemParams(targetSystem)
        const resp = await API.post('/api/maintenance/gsm-lit-sources', payload, { params })
        this.sources.push(resp.data)
        this.cancelCreate()
      } catch (err) {
        console.error('Failed to create source:', err)
        this.error = err.response?.data?.error || err.message
      } finally {
        this.isSaving = false
      }
    },
  },
}
</script>
