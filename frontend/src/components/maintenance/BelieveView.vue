<template>
  <div class="header-section">
    <h2>{{ $t('maintenance') }} - {{ $t('believe.title') }}</h2>
    <div class="search-box">
      <input
        v-model="searchTerm"
        type="text"
        :placeholder="$t('search')"
      />
    </div>
  </div>

  <div v-if="error" class="error-box">{{ error }}</div>

  <div class="cd-view">
    <div class="cd-list">
      <table class="cd-table">
        <thead>
          <tr>
            <th class="cd-table-header">{{ $t('believe.id') }}</th>
            <th class="cd-table-header">{{ $t('believe.name') }}</th>
            <th class="cd-table-header">{{ $t('believe.description') }}</th>
            <th class="cd-table-header">{{ $t('believe.source') }}</th>
            <th class="cd-table-header">{{ $t('believe.page') }}</th>
            <th class="cd-table-header">{{ $t('believe.system') }}</th>
            <th class="cd-table-header"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="isLoading">
            <td colspan="7">{{ $t('common.loading') }}</td>
          </tr>

          <template v-for="believe in filteredBelieves" :key="believe.id">
            <tr v-if="editingId !== believe.id">
              <td>{{ believe.id }}</td>
              <td>{{ believe.name }}</td>
              <td>{{ believe.beschreibung || '-' }}</td>
              <td>{{ getSourceCode(believe.source_id) || '-' }}</td>
              <td>{{ believe.page_number || '-' }}</td>
              <td>{{ getSystemCodeById(believe.game_system_id, believe.game_system) || '-' }}</td>
              <td>
                <button @click="startEdit(believe)">{{ $t('believe.edit') }}</button>
              </td>
            </tr>
            <tr v-else>
              <td>{{ believe.id }}</td>
              <td colspan="6">
                <div class="edit-form">
                  <div class="edit-row">
                    <label>{{ $t('believe.name') }}</label>
                    <input v-model="editedItem.name" />
                  </div>

                  <div class="edit-row">
                    <label>{{ $t('believe.description') }}</label>
                    <input v-model="editedItem.beschreibung" />
                  </div>

                  <div class="edit-row">
                    <label>{{ $t('believe.source') }}</label>
                    <select v-model="editedItem.sourceCode">
                      <option value="">-</option>
                      <option v-for="source in sources" :key="source.id" :value="source.code">
                        {{ source.code }}
                      </option>
                    </select>

                    <label class="inline-label">{{ $t('believe.page') }}</label>
                    <input v-model.number="editedItem.page_number" type="number" min="0" />
                  </div>

                  <div class="edit-row">
                    <label>{{ $t('believe.system') }}</label>
                    <select v-model.number="selectedSystemId">
                      <option value="">-</option>
                      <option v-for="system in systemOptions" :key="system.id" :value="system.id">
                        {{ system.label }}
                      </option>
                    </select>
                  </div>

                  <div class="edit-actions">
                    <button
                      class="btn-primary"
                      :disabled="isSaving"
                      @click="saveEdit"
                    >
                      <span v-if="!isSaving">{{ $t('believe.save') }}</span>
                      <span v-else>{{ $t('believe.saving') }}</span>
                    </button>
                    <button
                      class="btn-cancel"
                      :disabled="isSaving"
                      @click="cancelEdit"
                    >
                      {{ $t('believe.cancel') }}
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
/* Uses shared maintenance styles */
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

.edit-row label {
  font-weight: 600;
}

.edit-row input,
.edit-row select {
  padding: 6px 10px;
  border: 1px solid #dee2e6;
  border-radius: 6px;
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
  findSystemIdByCode,
  getSourceCode,
  getSystemCodeById,
  loadGameSystems as fetchGameSystems,
  buildSystemOptions,
} from '../../utils/maintenanceGameSystems'

export default {
  name: "BelieveView",
  props: {
    mdata: {
      type: Object,
      required: false,
      default: () => ({})
    }
  },
  data() {
    return {
      believes: [],
      sources: [],
      editingId: null,
      editedItem: null,
      gameSystems: [],
      selectedSystemId: null,
      isLoading: false,
      isSaving: false,
      error: '',
      searchTerm: ''
    }
  },
  async created() {
    await this.loadGameSystems()
    await this.loadBelieves()
  },
  computed: {
    filteredBelieves() {
      const term = this.searchTerm.trim().toLowerCase()
      const filtered = term
        ? this.believes.filter(believe => {
            const name = (believe.name || '').toLowerCase()
            const desc = (believe.beschreibung || '').toLowerCase()
            return name.includes(term) || desc.includes(term)
          })
        : this.believes

      return [...filtered].sort((a, b) => (a.name || '').localeCompare(b.name || ''))
    },
    systemOptions() {
      return buildSystemOptions(this.gameSystems)
    }
  },
  methods: {
    async loadGameSystems() {
      try {
        this.gameSystems = await fetchGameSystems()
      } catch (err) {
        console.error('Failed to load game systems:', err)
        this.error = err.response?.data?.error || err.message
      }
    },
    async loadBelieves() {
      this.isLoading = true
      this.error = ''
      try {
        const response = await API.get('/api/maintenance/gsm-believes')
        this.believes = response.data?.believes || []
        this.sources = response.data?.sources || []
      } catch (err) {
        console.error('Failed to load believes:', err)
        this.error = err.response?.data?.error || err.message
      } finally {
        this.isLoading = false
      }
    },
    getSourceCode(sourceId) {
      return getSourceCode(this.sources, sourceId)
    },
    getSystemCodeById(systemId, fallback = '') {
      return getSystemCodeById(this.gameSystems, systemId, fallback)
    },
    startEdit(believe) {
      this.editingId = believe.id
      this.editedItem = {
        ...believe,
        sourceCode: this.getSourceCode(believe.source_id)
      }
      this.selectedSystemId = believe.game_system_id ?? this.findSystemIdByCode(believe.game_system)
    },
    cancelEdit() {
      this.editingId = null
      this.editedItem = null
      this.selectedSystemId = null
    },
    findSystemIdByCode(code) {
      return findSystemIdByCode(this.gameSystems, code)
    },
    async saveEdit() {
      if (!this.editedItem || !this.editingId) {
        return
      }

      const trimmedName = (this.editedItem.name || '').trim()
      if (!trimmedName) {
        alert(this.$t('believe.nameRequired'))
        return
      }

      const selectedSource = this.sources.find(src => src.code === this.editedItem.sourceCode)
      const selectedSystem = this.gameSystems.find(gs => gs.id === this.selectedSystemId)

      const payload = {
        name: trimmedName,
        beschreibung: this.editedItem.beschreibung || '',
        source_id: selectedSource ? selectedSource.id : null,
        page_number: this.editedItem.page_number || 0,
        game_system_id: selectedSystem ? selectedSystem.id : null,
        game_system: selectedSystem ? selectedSystem.code : '',
      }

      this.isSaving = true
      try {
        const response = await API.put(`/api/maintenance/gsm-believes/${this.editingId}`, payload)
        const updated = response.data
        const sourceCode = this.getSourceCode(updated.source_id)
        const gameSystemCode = selectedSystem ? selectedSystem.code : updated.game_system
        const gameSystemId = selectedSystem ? selectedSystem.id : (updated.game_system_id ?? null)

        const idx = this.believes.findIndex(b => b.id === this.editingId)
        if (idx !== -1) {
          this.believes.splice(idx, 1, { ...updated, source_code: sourceCode, game_system: gameSystemCode, game_system_id: gameSystemId })
        }

        this.cancelEdit()
      } catch (err) {
        console.error('Failed to save believe:', err)
        this.error = err.response?.data?.error || err.message
      } finally {
        this.isSaving = false
      }
    }
  }
}
</script>
