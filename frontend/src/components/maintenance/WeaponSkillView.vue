<template>
  <div class="header-section">
    <h2>{{ $t('maintenance') }}</h2>
      <!-- Add search input -->
      <div class="search-box">
        <input
          type="text"
          v-model="searchTerm"
          :placeholder="`${$t('search')} ${$t('WaeponSkill')}...`"
        />
      </div>
    </div>

  <div class="cd-view">
    <div class="cd-list">
      <!-- Filter Row -->
      <div class="filter-row">
        <div class="filter-item">
          <label>{{ $t('weaponskill.difficulty') }}:</label>
          <select v-model="filterDifficulty">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option v-for="diff in availableDifficulties" :key="diff" :value="diff">{{ diff }}</option>
          </select>
        </div>
        <div class="filter-item">
          <label>{{ $t('weaponskill.quelle') }}:</label>
          <select v-model="filterQuelle">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option v-for="quelle in availableQuellen" :key="quelle" :value="quelle">{{ quelle }}</option>
          </select>
        </div>
        <button @click="clearFilters" class="btn-clear-filters">{{ $t('clearFilters') || 'Clear Filters' }}</button>
      </div>

      <div class="tables-container">
          <table class="cd-table">
            <thead>
              <tr>
                <th class="cd-table-header">{{ $t('weaponskill.id') }}</th>
                <th class="cd-table-header">{{ $t('weaponskill.name') }} <button @click="sortBy('name')">-{{ sortField === 'name' ? (sortAsc ? '↑' : '↓') : '' }}</button></th>
                <th class="cd-table-header">{{ $t('weaponskill.difficulty') }}<button @click="sortBy('difficulty')">-{{ sortField === 'difficulty' ? (sortAsc ? '↑' : '↓') : '' }}</button></th>
                <th class="cd-table-header">{{ $t('weaponskill.initialwert') }}</th>
                <th class="cd-table-header">{{ $t('weaponskill.description') }}</th>
                <th class="cd-table-header">{{ $t('weaponskill.quelle') }}</th>
                <th class="cd-table-header">{{ $t('weaponskill.system') }}</th>
                <th class="cd-table-header"> </th>
              </tr>
            </thead>
            <tbody>
              <template v-for="(dtaItem, index) in filteredAndSortedWaeponSkills" :key="dtaItem.id">
                <tr v-if="editingIndex !== index">
                  <td>{{ dtaItem.id || '' }}</td>
                  <td>{{ dtaItem.name || '-' }}</td>
                  <td>{{ dtaItem.difficulty || '-' }}</td>
                  <td>{{ dtaItem.initialwert || '0' }}</td>
                  <td>{{ dtaItem.beschreibung || '-' }}</td>
                  <td>{{ formatQuelle(dtaItem) }}</td>
                  <td>{{ getSystemCodeById(dtaItem.game_system_id, dtaItem.system || 'midgard') }}</td>
                  <td>
                    <button @click="startEdit(index)">Edit</button>
                  </td>
                </tr>
                <!-- Edit Mode -->
                <tr v-else>
                  <td><input v-model="editedItem.id" style="width:20px;" disabled /></td>
                  <td colspan="6">
                    <!-- Expanded edit form -->
                    <div class="edit-form">
                      <div class="edit-row">
                        <div class="edit-field">
                          <label>{{ $t('weaponskill.name') }}:</label>
                          <input v-model="editedItem.name" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('weaponskill.difficulty') }}:</label>
                          <select v-model="editedItem.difficulty" style="width:120px;">
                            <option value="leicht">leicht</option>
                            <option value="normal">normal</option>
                            <option value="schwer">schwer</option>
                            <option value="sehr schwer">sehr schwer</option>
                          </select>
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('weaponskill.initialwert') }}:</label>
                          <input v-model.number="editedItem.initialwert" type="number" style="width:60px;" />
                        </div>
                      </div>

                      <div class="edit-row">
                        <div class="edit-field full-width">
                          <label>{{ $t('weaponskill.description') }}:</label>
                          <input v-model="editedItem.beschreibung" />
                        </div>
                      </div>

                      <div class="edit-row">
                        <div class="edit-field">
                          <label>{{ $t('weaponskill.quelle') }}:</label>
                          <select v-model="editedItem.sourceCode" style="width:100px;">
                            <option value="">-</option>
                            <option v-for="source in availableSources" :key="source.code" :value="source.code">
                              {{ source.code }}
                            </option>
                          </select>
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('weaponskill.page') || 'Page' }}:</label>
                          <input v-model.number="editedItem.page_number" type="number" style="width:60px;" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('weaponskill.system') }}:</label>
                          <select v-model.number="selectedSystemId" style="width:140px;">
                            <option value="">-</option>
                            <option v-for="system in systemOptions" :key="system.id" :value="system.id">
                              {{ system.label }}
                            </option>
                          </select>
                        </div>
                      </div>

                      <div class="edit-actions">
                        <button @click="saveEdit(index)" class="btn-save">Save</button>
                        <button @click="cancelEdit" class="btn-cancel">Cancel</button>
                      </div>
                    </div>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
      </div>
    </div> <!--- end cd-list-->
  </div> <!--- end character -datasheet-->
</template>

<!-- <style scoped> -->
<style>
.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.3rem;
  height: fit-content;
  padding: 0.5rem;
}
.search-box {
  margin-bottom: 1rem;
}
.search-box input {
  padding: 0.2rem;
  width: 200px;
  border: 1px solid #ddd;
  border-radius: 4px;
}
.tables-container {
  display: flex;
  gap: 1rem;
  width: 100%;
}

.table-wrapper-left {
  flex: 6;
  min-width: 0; /* Prevent table from overflowing */
}
.table-wrapper-right {
  flex: 4;
  min-width: 0; /* Prevent table from overflowing */
}

.cd-table {
  width: 100%;
}
.cd-table-header {
  background-color: #1da766;
  font-weight: bold;
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
  name: "WaeponSkillView",
  props: {
    mdata: {
      type: Object,
      required: true,
      default: () => ({
        skills: [],
        skillcategories: []
      })
    }
  },
  data() {
    return {
      searchTerm: '',
      sortField: 'name',
      sortAsc: true,
      editingIndex: -1,
      editedItem: null,
      filterDifficulty: '',
      filterQuelle: '',
      enhancedWeaponSkills: [],
      availableSources: [],
      availableDifficultiesData: [],
      gameSystems: [],
      selectedSystemId: null
    }
  },
  async created() {
    await Promise.all([
      this.loadGameSystems(),
      this.loadEnhancedWeaponSkills()
    ])
  },
  computed: {
    availableDifficulties() {
      const difficulties = new Set()
      this.enhancedWeaponSkills.forEach(ws => {
        if (ws.difficulty) difficulties.add(ws.difficulty)
      })
      return Array.from(difficulties).sort()
    },
    availableQuellen() {
      const quellen = new Set()
      this.enhancedWeaponSkills.forEach(ws => {
        if (ws.source_id && this.availableSources.length > 0) {
          const source = this.availableSources.find(s => s.id === ws.source_id)
          if (source) {
            quellen.add(source.code)
          }
        }
      })
      return Array.from(quellen).sort()
    },
    filteredAndSortedWaeponSkills() {
      let filtered = [...this.enhancedWeaponSkills]

      // Apply search filter
      if (this.searchTerm) {
        const searchLower = this.searchTerm.toLowerCase()
        filtered = filtered.filter(ws =>
          ws.name?.toLowerCase().includes(searchLower) ||
          ws.difficulty?.toLowerCase().includes(searchLower)
        )
      }

      // Apply difficulty filter
      if (this.filterDifficulty) {
        filtered = filtered.filter(ws => ws.difficulty === this.filterDifficulty)
      }

      // Apply Quelle filter (only by source code, ignoring page number)
      if (this.filterQuelle) {
        filtered = filtered.filter(ws => {
          if (ws.source_id && this.availableSources.length > 0) {
            const source = this.availableSources.find(s => s.id === ws.source_id)
            return source && source.code === this.filterQuelle
          }
          return false
        })
      }

      // Apply sorting
      filtered.sort((a, b) => {
        const aValue = (a[this.sortField] || '').toString().toLowerCase()
        const bValue = (b[this.sortField] || '').toString().toLowerCase()
        return this.sortAsc ? aValue.localeCompare(bValue) : bValue.localeCompare(aValue)
      })

      return filtered
    },
    systemOptions() {
      return buildSystemOptions(this.gameSystems)
    }
  },
  methods: {
    async loadGameSystems() {
      try {
        this.gameSystems = await fetchGameSystems()
      } catch (error) {
        console.error('Failed to load game systems:', error)
      }
    },
    async loadEnhancedWeaponSkills() {
      try {
        const response = await API.get('/api/maintenance/weaponskills-enhanced')
        this.enhancedWeaponSkills = response.data.weaponskills || []
        this.availableSources = response.data.sources || []
        this.availableDifficultiesData = response.data.difficulties || []
      } catch (error) {
        console.error('Failed to load enhanced weapon skills:', error)
      }
    },
    startEdit(index) {
      const weaponSkill = this.filteredAndSortedWaeponSkills[index]
      this.editedItem = {
        ...weaponSkill,
        sourceCode: this.getSourceCode(weaponSkill.source_id)
      }
      this.selectedSystemId = weaponSkill.game_system_id ?? this.findSystemIdByCode(weaponSkill.system)
      this.editingIndex = index
    },
    async saveEdit(index) {
      try {
        // Find source ID from code
        const source = this.availableSources.find(s => s.code === this.editedItem.sourceCode)
        const selectedSystem = this.gameSystems.find(gs => gs.id === this.selectedSystemId)
        
        const updateData = {
          ...this.editedItem,
          source_id: source ? source.id : null,
          page_number: this.editedItem.page_number || 0,
          difficulty: this.editedItem.difficulty,
          category: 'Waffen', // Weapon skills always use 'Waffen' category
          system: selectedSystem ? selectedSystem.code : (this.editedItem.system || ''),
          game_system_id: selectedSystem ? selectedSystem.id : (this.editedItem.game_system_id ?? null)
        }
        
        const response = await API.put(
          `/api/maintenance/weaponskills-enhanced/${this.editedItem.id}`,
          updateData
        )

        // Reload the list to get updated difficulty from backend
        await this.loadEnhancedWeaponSkills()

        this.editingIndex = -1
        this.editedItem = null
        this.selectedSystemId = null
      } catch (error) {
        console.error('Failed to save weapon skill:', error)
        alert('Failed to save weapon skill: ' + (error.response?.data?.error || error.message))
      }
    },
    cancelEdit() {
      this.editingIndex = -1
      this.editedItem = null
      this.selectedSystemId = null
    },
    findSystemIdByCode(code) {
      return findSystemIdByCode(this.gameSystems, code)
    },
    sortBy(field) {
      if (this.sortField === field) {
        this.sortAsc = !this.sortAsc
      } else {
        this.sortField = field
        this.sortAsc = true
      }
    },
    formatQuelle(weaponSkill) {
      if (weaponSkill.source_id && this.availableSources.length > 0) {
        const source = this.availableSources.find(s => s.id === weaponSkill.source_id)
        if (source) {
          if (weaponSkill.page_number) {
            return `${source.code}:${weaponSkill.page_number}`
          } else {
            // No page number - show code and quelle if available
            const quelle = weaponSkill.quelle ? ` (${weaponSkill.quelle})` : ''
            return `${source.code}${quelle}`
          }
        }
      }
      return weaponSkill.quelle || '-'
    },
    getSourceCode(sourceId) {
      return getSourceCode(this.availableSources, sourceId)
    },
    getSystemCodeById(systemId, fallback = '') {
      return getSystemCodeById(this.gameSystems, systemId, fallback)
    },
    clearFilters() {
      this.searchTerm = ''
      this.filterDifficulty = ''
      this.filterQuelle = ''
    }
  }
};
</script>
