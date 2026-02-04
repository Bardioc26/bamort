<template>
  <div class="header-section">
    <h2>{{ $t('maintenance') }} - {{ $t('misc.title') }}</h2>
    <div class="search-box">
      <input v-model="searchTerm" type="text" :placeholder="$t('search')" />
      <select v-model="filterKey">
        <option value="">{{ $t('all') || 'All' }}</option>
        <option v-for="key in keyOptions" :key="key" :value="key">{{ key }}</option>
      </select>
      <button class="btn-primary" @click="startCreate">{{ $t('newEntry') }}</button>
    </div>
  </div>

  <div v-if="error" class="error-box">{{ error }}</div>

  <div class="cd-view">
    <div class="cd-list">
      <table class="cd-table">
        <thead>
          <tr>
            <th class="cd-table-header">{{ $t('misc.id') }}</th>
            <th class="cd-table-header">{{ $t('misc.key') }}</th>
            <th class="cd-table-header">{{ $t('misc.value') }}</th>
            <th class="cd-table-header">{{ $t('misc.source') }}</th>
            <th class="cd-table-header">{{ $t('misc.page') }}</th>
            <th class="cd-table-header">{{ $t('misc.system') }}</th>
            <th class="cd-table-header"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="isLoading">
            <td colspan="7">{{ $t('common.loading') }}</td>
          </tr>

          <tr v-if="creatingNew">
            <td>New</td>
            <td colspan="6">
              <div class="edit-form">
                <div class="edit-row">
                  <label>{{ $t('misc.key') }}</label>
                  <input v-model="newItem.key" list="misc-key-options" />
                  <datalist id="misc-key-options">
                    <option v-for="key in keyOptions" :key="key" :value="key" />
                  </datalist>
                  <label class="inline-label">{{ $t('misc.value') }}</label>
                  <input v-model="newItem.value" />
                </div>
                <div class="edit-row">
                  <label>{{ $t('misc.source') }}</label>
                  <select v-model.number="newItem.source_id">
                    <option :value="null">-</option>
                    <option v-for="src in sourceOptions" :key="src.id" :value="src.id">
                      {{ src.label }}
                    </option>
                  </select>
                  <label class="inline-label">{{ $t('misc.page') }}</label>
                  <input v-model.number="newItem.page_number" type="number" min="0" />
                </div>
                <div class="edit-row">
                  <label>{{ $t('misc.system') }}</label>
                  <select v-model.number="createSelectedSystemId">
                    <option v-for="sys in systemOptions" :key="sys.id" :value="sys.id">{{ sys.label }}</option>
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

          <template v-for="item in filteredItems" :key="item.id">
            <tr v-if="editingId !== item.id">
              <td>{{ item.id }}</td>
              <td>{{ item.key }}</td>
              <td>{{ item.value }}</td>
              <td>{{ sourceCodeFor(item.source_id) }}</td>
              <td>{{ item.page_number || '-' }}</td>
              <td>{{ systemCodeFor(item.game_system_id, item.game_system) || '-' }}</td>
              <td><button @click="startEdit(item)">{{ $t('common.edit') }}</button></td>
            </tr>
            <tr v-else>
              <td>{{ item.id }}</td>
              <td colspan="6">
                <div class="edit-form">
                  <div class="edit-row">
                    <label>{{ $t('misc.key') }}</label>
                    <select v-model="editedItem.key">
                      <option :value="''">-</option>
                      <option v-for="key in keyOptionsWithCurrent" :key="key" :value="key">{{ key }}</option>
                    </select>
                    <label class="inline-label">{{ $t('misc.value') }}</label>
                    <input v-model="editedItem.value" />
                  </div>
                  <div class="edit-row">
                    <label>{{ $t('misc.source') }}</label>
                    <select v-model.number="editedItem.source_id">
                      <option :value="null">-</option>
                      <option v-for="src in sourceOptions" :key="src.id" :value="src.id">
                        {{ src.label }}
                      </option>
                    </select>
                    <label class="inline-label">{{ $t('misc.page') }}</label>
                    <input v-model.number="editedItem.page_number" type="number" min="0" />
                  </div>
                  <div class="edit-row">
                    <label>{{ $t('misc.system') }}</label>
                    <select v-model.number="selectedSystemId">
                      <option v-for="sys in systemOptions" :key="sys.id" :value="sys.id">{{ sys.label }}</option>
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
  loadGameSystems as fetchGameSystems,
  systemCodeFor as resolveSystemCode,
  buildSystemOptions,
} from '../../utils/maintenanceGameSystems'

export default {
  name: 'MiscLookupView',
  data() {
    return {
      items: [],
      gameSystems: [],
      currentGameSystem: null,
      sources: [],
      editingId: null,
      editedItem: null,
      selectedSystemId: null,
      creatingNew: false,
      newItem: null,
      createSelectedSystemId: null,
      isLoading: false,
      isSaving: false,
      error: '',
      searchTerm: '',
      filterKey: '',
    }
  },
  async created() {
    await this.initialize()
  },
  computed: {
    filteredItems() {
      const term = this.searchTerm.trim().toLowerCase()
      let list = term
        ? this.items.filter(it =>
            (it.key || '').toLowerCase().includes(term) ||
            (it.value || '').toLowerCase().includes(term)
          )
        : this.items

      if (this.filterKey) {
        list = list.filter(it => it.key === this.filterKey)
      }

      return [...list].sort((a, b) => (a.key || '').localeCompare(b.key || ''))
    },
    keyOptions() {
      const set = new Set()
      this.items.forEach(it => {
        if (it.key) set.add(it.key)
      })
      return Array.from(set.values()).sort()
    },
    keyOptionsWithCurrent() {
      const set = new Set(this.keyOptions)
      if (this.editedItem?.key) set.add(String(this.editedItem.key))
      return Array.from(set.values()).sort()
    },
    systemOptions() {
      const labelBuilder = system => {
        const code = system.code || ''
        const name = system.name || ''
        if (code && name) return `${code} - ${name}`.trim()
        return code || name || String(system.id ?? '')
      }
      return buildSystemOptions(this.gameSystems, labelBuilder)
    },
    sourceMap() {
      const map = new Map()
      this.sources.forEach(src => {
        map.set(src.id, src.code || src.name || src.id)
      })
      return map
    },
    sourceOptions() {
      return this.sources.map(src => ({
        id: src.id,
        label: src.code ? `${src.code} - ${src.name || ''}`.trim() : src.name || src.id,
      }))
    },
  },
  methods: {
    async initialize() {
      this.error = ''
      await this.loadGameSystems()
      if (!this.currentGameSystem) return
      await this.loadSources()
      await this.loadItems()
    },
    startCreate() {
      this.cancelEdit()
      const defaultSystem = this.currentGameSystem || this.gameSystems.find(gs => gs.is_active) || this.gameSystems[0] || null
      this.createSelectedSystemId = defaultSystem ? defaultSystem.id : null
      this.newItem = {
        key: '',
        value: '',
        source_id: null,
        page_number: 0,
      }
      this.creatingNew = true
    },
    cancelCreate() {
      this.creatingNew = false
      this.newItem = null
      this.createSelectedSystemId = null
    },
    async loadGameSystems() {
      try {
        const systems = await fetchGameSystems()
        this.gameSystems = systems
        const active = systems.find(s => s.is_active)
        this.currentGameSystem = active || systems[0] || null
      } catch (err) {
        console.error('Failed to load game systems:', err)
        this.error = err.response?.data?.error || err.message
      }
    },
    async loadSources() {
      try {
        const params = buildGameSystemParams(this.currentGameSystem)
        const resp = await API.get('/api/maintenance/gsm-lit-sources', { params })
        this.sources = resp.data?.sources || []
      } catch (err) {
        console.error('Failed to load sources:', err)
        this.error = err.response?.data?.error || err.message
      }
    },
    async loadItems() {
      this.isLoading = true
      this.error = ''
      try {
        const params = buildGameSystemParams(this.currentGameSystem)
        const resp = await API.get('/api/maintenance/gsm-misc', { params })
        this.items = resp.data?.misc || []
      } catch (err) {
        console.error('Failed to load misc:', err)
        this.error = err.response?.data?.error || err.message
      } finally {
        this.isLoading = false
      }
    },
    buildParamsForSystemId(systemId) {
      const sys = findSystemById(this.gameSystems, systemId) || this.currentGameSystem
      if (!sys) return {}
      return buildGameSystemParams(sys)
    },
    sourceCodeFor(id) {
      if (!id) return '-'
      const code = this.sourceMap.get(id)
      return code || id
    },
    systemCodeFor(systemId, fallback = '') {
      return resolveSystemCode(this.gameSystems, systemId, fallback)
    },
    startEdit(item) {
      this.editingId = item.id
      this.editedItem = { ...item }
      this.selectedSystemId = item.game_system_id ?? this.currentGameSystem?.id ?? null
    },
    cancelEdit() {
      this.editingId = null
      this.editedItem = null
      this.selectedSystemId = null
    },
    async saveEdit() {
      if (!this.editedItem) return
      const payload = {
        key: this.editedItem.key || '',
        value: this.editedItem.value || '',
        source_id: this.editedItem.source_id || null,
        page_number: this.editedItem.page_number || 0,
      }
      this.isSaving = true
      try {
        const params = this.buildParamsForSystemId(this.selectedSystemId)
        const resp = await API.put(`/api/maintenance/gsm-misc/${this.editingId}`, payload, { params })
        const idx = this.items.findIndex(i => i.id === this.editingId)
        if (idx !== -1) this.items.splice(idx, 1, resp.data)
        this.cancelEdit()
      } catch (err) {
        console.error('Failed to save misc:', err)
        this.error = err.response?.data?.error || err.message
      } finally {
        this.isSaving = false
      }
    },
    async saveCreate() {
      if (!this.newItem) return
      const payload = {
        key: this.newItem.key || '',
        value: this.newItem.value || '',
        source_id: this.newItem.source_id || null,
        page_number: this.newItem.page_number || 0,
      }
      this.isSaving = true
      try {
        const params = this.buildParamsForSystemId(this.createSelectedSystemId)
        const resp = await API.post('/api/maintenance/gsm-misc', payload, { params })
        this.items.push(resp.data)
        this.cancelCreate()
      } catch (err) {
        console.error('Failed to create misc entry:', err)
        this.error = err.response?.data?.error || err.message
      } finally {
        this.isSaving = false
      }
    },
  },
}
</script>
