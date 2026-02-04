<template>
  <div class="header-section">
    <h2>{{ $t('maintenance') }}</h2>
    <!-- Add search input -->
    <div class="search-box">
      <input
        type="text"
        v-model="searchTerm"
        :placeholder="`${$t('search')} ${$t('Spell')}...`"
      />
      <button @click="startCreate" class="btn-primary">{{ $t('newEntry') }}</button>
    </div>
  </div>
  
  <!-- Import CSV Section -->
  <div class="import-section">
    <h3>{{ $t('spell.import') || 'Import Spells' }}</h3>
    <div class="file-upload">
      <input
        ref="fileInput"
        type="file"
        accept=".csv"
        @change="handleFileSelect"
        style="display: none;"
      />
      <button 
        @click="$refs.fileInput.click()"
        :disabled="importing"
        class="upload-btn"
      >
        {{ importing ? ($t('spell.importing') || 'Importing...') : ($t('spell.selectCsv') || 'Select CSV File') }}
      </button>
      <span v-if="selectedFile" class="file-name">{{ selectedFile.name }}</span>
    </div>
    <button 
      v-if="selectedFile && !importing"
      @click="importSpells"
      class="import-btn"
    >
      {{ $t('spell.import') || 'Import Spells' }}
    </button>
    <div v-if="importResult" class="import-result" :class="importResult.success ? 'success' : 'error'">
      async saveCreate() {
        if (!this.newItem) return
        try {
          const source = this.availableSources.find(s => s.code === this.newItem.sourceCode)
          const selectedSystem = this.gameSystems.find(gs => gs.id === this.createSelectedSystemId)

          const createData = {
            ...this.newItem,
            source_id: source ? source.id : null,
            page_number: this.newItem.page_number || 0,
            system: selectedSystem ? selectedSystem.code : (this.newItem.system || ''),
            game_system_id: selectedSystem ? selectedSystem.id : null
          }

          const response = await API.post(
            '/api/maintenance/spells-enhanced',
            createData
          )

          this.enhancedSpells.push(response.data)
          this.cancelCreate()
        } catch (error) {
          console.error('Failed to create spell:', error)
          alert('Failed to create spell: ' + (error.response?.data?.error || error.message))
        }
      },
      {{ importResult.message }}
      <span v-if="importResult.total_spells"> ({{ importResult.total_spells }} spells total)</span>
    </div>
  </div>

  <div class="cd-view">
    <div class="cd-list">
      <!-- Filter Row -->
      <div class="filter-row">
        <div class="filter-item">
          <label>{{ $t('spell.category') }}:</label>
          <select v-model="filterCategory">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option v-for="cat in availableCategories" :key="cat" :value="cat">{{ cat }}</option>
          </select>
        </div>
        <div class="filter-item">
          <label>{{ $t('spell.level') }}:</label>
          <select v-model="filterLevel">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option v-for="level in availableLevels" :key="level" :value="level">{{ level }}</option>
          </select>
        </div>
        <div class="filter-item">
          <label>{{ $t('spell.ursprung') }}:</label>
          <select v-model="filterUrsprung">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option v-for="ursprung in availableUrsprungs" :key="ursprung" :value="ursprung">{{ ursprung }}</option>
          </select>
        </div>
        <div class="filter-item">
          <label>{{ $t('spell.reichweite') }}:</label>
          <select v-model="filterReichweite">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option v-for="reichweite in availableReichweiten" :key="reichweite" :value="reichweite">{{ reichweite }}</option>
          </select>
        </div>
        <div class="filter-item">
          <label>{{ $t('spell.wirkungsziel') }}:</label>
          <select v-model="filterWirkungsziel">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option v-for="wirkungsziel in availableWirkungsziele" :key="wirkungsziel" :value="wirkungsziel">{{ wirkungsziel }}</option>
          </select>
        </div>
        <div class="filter-item">
          <label>{{ $t('spell.quelle') }}:</label>
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
                <th class="cd-table-header">{{ $t('spell.id') }}</th>
                <th class="cd-table-header">
                  {{ $t('spell.category') }}
                  <button @click="sortBy('category')">{{ sortField === 'category' ? (sortAsc ? '↑' : '↓') : '-' }}</button>
                </th>
                <th class="cd-table-header">
                  {{ $t('spell.name') }}
                  <button @click="sortBy('name')">{{ sortField === 'name' ? (sortAsc ? '↑' : '↓') : '-' }}</button>
                </th>
                <th class="cd-table-header">{{ $t('spell.level') }}</th>
                <th class="cd-table-header">{{ $t('spell.apverbrauch') }}</th>
                <th class="cd-table-header">{{ $t('spell.zauberdauer') }}</th>
                <th class="cd-table-header">{{ $t('spell.reichweite') }}</th>
                <th class="cd-table-header">{{ $t('spell.wirkungsziel') }}</th>
                <th class="cd-table-header">{{ $t('spell.wirkungsbereich') }}</th>
                <th class="cd-table-header">{{ $t('spell.wirkungsdauer') }}</th>
                <th class="cd-table-header">{{ $t('spell.ursprung') }}</th>
                <th class="cd-table-header">{{ $t('spell.description') }}</th>
                <th class="cd-table-header">{{ $t('spell.quelle') }}</th>
                <th class="cd-table-header">{{ $t('spell.system') }}</th>
                <th class="cd-table-header"> </th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="creatingNew">
                <td>New</td>
                <td colspan="14">
                  <div class="edit-form">
                    <div class="edit-row">
                      <div class="edit-field">
                        <label>{{ $t('spell.name') }}:</label>
                        <input v-model="newItem.name" />
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('spell.category') }}:</label>
                        <select v-model="newItem.category" style="width:120px;">
                          <option v-for="category in mdata['spellcategories']" :key="category" :value="category">
                            {{ category }}
                          </option>
                        </select>
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('spell.level') }}:</label>
                        <input v-model.number="newItem.level" type="number" style="width:60px;" />
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('spell.apverbrauch') }}:</label>
                        <input v-model="newItem.ap" style="width:60px;" />
                      </div>
                    </div>

                    <div class="edit-row">
                      <div class="edit-field">
                        <label>{{ $t('spell.zauberdauer') }}:</label>
                        <input v-model="newItem.zauberdauer" style="width:120px;" />
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('spell.reichweite') }}:</label>
                        <input v-model="newItem.reichweite" style="width:120px;" />
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('spell.wirkungsdauer') }}:</label>
                        <input v-model="newItem.wirkungsdauer" style="width:120px;" />
                      </div>
                    </div>

                    <div class="edit-row">
                      <div class="edit-field">
                        <label>{{ $t('spell.wirkungsziel') }}:</label>
                        <input v-model="newItem.wirkungsziel" style="width:150px;" />
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('spell.wirkungsbereich') }}:</label>
                        <input v-model="newItem.wirkungsbereich" style="width:150px;" />
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('spell.ursprung') }}:</label>
                        <input v-model="newItem.ursprung" style="width:120px;" />
                      </div>
                    </div>

                    <div class="edit-row">
                      <div class="edit-field full-width">
                        <label>{{ $t('spell.description') }}:</label>
                        <input v-model="newItem.beschreibung" />
                      </div>
                    </div>

                    <div class="edit-row">
                      <div class="edit-field">
                        <label>{{ $t('spell.quelle') }}:</label>
                        <select v-model="newItem.sourceCode" style="width:100px;">
                          <option value="">-</option>
                          <option v-for="source in availableSources" :key="source.code" :value="source.code">
                            {{ source.code }}
                          </option>
                        </select>
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('spell.page') || 'Page' }}:</label>
                        <input v-model.number="newItem.page_number" type="number" style="width:60px;" />
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('spell.system') }}:</label>
                        <select v-model.number="createSelectedSystemId" style="width:140px;">
                          <option value="">-</option>
                          <option v-for="system in systemOptions" :key="system.id" :value="system.id">
                            {{ system.label }}
                          </option>
                        </select>
                      </div>
                    </div>

                    <div class="edit-actions">
                      <button @click="saveCreate" class="btn-save">{{ $t('common.save') }}</button>
                      <button @click="cancelCreate" class="btn-cancel">{{ $t('common.cancel') }}</button>
                    </div>
                  </div>
                </td>
              </tr>
              <template v-for="(dtaItem, index) in filteredAndSortedSpells" :key="dtaItem.id">
                <tr v-if="editingIndex !== index">
                  <td>{{ dtaItem.id || '' }}</td>
                  <td>{{ dtaItem.category|| '-' }}</td>
                  <td>{{ dtaItem.name || '-' }}</td>
                  <td>{{ dtaItem.level || '0' }}</td>
                  <td>{{ dtaItem.ap || '0' }}</td>
                  <td>{{ dtaItem.zauberdauer || '-' }}</td>
                  <td>{{ dtaItem.reichweite || '0' }}</td>
                  <td>{{ dtaItem.wirkungsziel || '-' }}</td>
                  <td>{{ dtaItem.wirkungsbereich || '-' }}</td>
                  <td>{{ dtaItem.wirkungsdauer || '-' }}</td>
                  <td>{{ dtaItem.ursprung || '-' }}</td>
                  <td>{{ dtaItem.beschreibung || '-' }}</td>
                  <td>{{ formatQuelle(dtaItem) }}</td>
                  <td>{{ getSystemCodeById(dtaItem.game_system_id, dtaItem.system || 'midgard') }}</td>
                  <td>
                    <button @click="startEdit(index)">{{ $t('common.edit') }}</button>
                  </td>
                </tr>
                <!-- Edit Mode -->
                <tr v-else>
                  <td><input v-model="editedItem.id" style="width:20px;" disabled /></td>
                  <td colspan="14">
                    <!-- Expanded edit form -->
                    <div class="edit-form">
                      <div class="edit-row">
                        <div class="edit-field">
                          <label>{{ $t('spell.name') }}:</label>
                          <input v-model="editedItem.name" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('spell.category') }}:</label>
                          <select v-model="editedItem.category" style="width:120px;">
                            <option v-for="category in mdata['spellcategories']" :key="category" :value="category">
                              {{ category }}
                            </option>
                          </select>
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('spell.level') }}:</label>
                          <input v-model.number="editedItem.level" type="number" style="width:60px;" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('spell.apverbrauch') }}:</label>
                          <input v-model="editedItem.ap" style="width:60px;" />
                        </div>
                      </div>

                      <div class="edit-row">
                        <div class="edit-field">
                          <label>{{ $t('spell.zauberdauer') }}:</label>
                          <input v-model="editedItem.zauberdauer" style="width:120px;" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('spell.reichweite') }}:</label>
                          <input v-model="editedItem.reichweite" style="width:120px;" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('spell.wirkungsdauer') }}:</label>
                          <input v-model="editedItem.wirkungsdauer" style="width:120px;" />
                        </div>
                      </div>

                      <div class="edit-row">
                        <div class="edit-field">
                          <label>{{ $t('spell.wirkungsziel') }}:</label>
                          <input v-model="editedItem.wirkungsziel" style="width:150px;" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('spell.wirkungsbereich') }}:</label>
                          <input v-model="editedItem.wirkungsbereich" style="width:150px;" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('spell.ursprung') }}:</label>
                          <input v-model="editedItem.ursprung" style="width:120px;" />
                        </div>
                      </div>

                      <div class="edit-row">
                        <div class="edit-field full-width">
                          <label>{{ $t('spell.description') }}:</label>
                          <input v-model="editedItem.beschreibung" />
                        </div>
                      </div>

                      <div class="edit-row">
                        <div class="edit-field">
                          <label>{{ $t('spell.quelle') }}:</label>
                          <select v-model="editedItem.sourceCode" style="width:100px;">
                            <option value="">-</option>
                            <option v-for="source in availableSources" :key="source.code" :value="source.code">
                              {{ source.code }}
                            </option>
                          </select>
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('spell.page') || 'Page' }}:</label>
                          <input v-model.number="editedItem.page_number" type="number" style="width:60px;" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('spell.system') }}:</label>
                          <select v-model.number="selectedSystemId" style="width:140px;">
                            <option value="">-</option>
                            <option v-for="system in systemOptions" :key="system.id" :value="system.id">
                              {{ system.label }}
                            </option>
                          </select>
                        </div>
                      </div>

                      <div class="edit-actions">
                        <button @click="saveEdit(index)" class="btn-save">{{ $t('common.save') }}</button>
                        <button @click="cancelEdit" class="btn-cancel">{{ $t('common.cancel') }}</button>
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
  margin-bottom: 1rem;
  height: fit-content;
  padding: 0.5rem;
  border-bottom: 1px solid #ddd;
}

.search-box {
  margin-left: auto;
}

.search-box input {
  padding: 0.5rem;
  width: 250px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.import-section {
  margin-bottom: 1.5rem;
  padding: 1.5rem;
  border: 2px solid #1da766;
  border-radius: 8px;
  background-color: #f8fcfa;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.import-section h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.file-upload {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

.upload-btn, .import-btn {
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  background-color: #1da766;
  color: white;
  cursor: pointer;
  transition: background-color 0.2s;
}

.upload-btn:hover, .import-btn:hover {
  background-color: #166d4a;
}

.upload-btn:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.file-name {
  font-style: italic;
  color: #666;
}

.import-result {
  padding: 0.5rem;
  border-radius: 4px;
  margin-top: 1rem;
}

.import-result.success {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.import-result.error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
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
  name: "SpellView",
  props: {
    mdata: {
      type: Object,
      required: true,
      default: () => ({
        spells: [],
        spellcategories: []
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
      selectedFile: null,
      importing: false,
      importResult: null,
      filterCategory: '',
      filterLevel: '',
      filterUrsprung: '',
      filterReichweite: '',
      filterWirkungsziel: '',
      filterQuelle: '',
      enhancedSpells: [],
      availableSources: [],
      gameSystems: [],
      selectedSystemId: null,
      creatingNew: false,
      newItem: null,
      createSelectedSystemId: null
    }
  },
  async created() {
    await Promise.all([
      this.loadGameSystems(),
      this.loadEnhancedSpells()
    ])
  },
  computed: {
    availableCategories() {
      const categories = new Set()
      this.enhancedSpells.forEach(spell => {
        if (spell.category) categories.add(spell.category)
      })
      return Array.from(categories).sort()
    },
    availableLevels() {
      const levels = new Set()
      this.enhancedSpells.forEach(spell => {
        if (spell.level !== null && spell.level !== undefined) levels.add(spell.level)
      })
      return Array.from(levels).sort((a, b) => a - b)
    },
    availableUrsprungs() {
      const ursprungs = new Set()
      this.enhancedSpells.forEach(spell => {
        if (spell.ursprung) ursprungs.add(spell.ursprung)
      })
      return Array.from(ursprungs).sort()
    },
    availableReichweiten() {
      const reichweiten = new Set()
      this.enhancedSpells.forEach(spell => {
        if (spell.reichweite) reichweiten.add(spell.reichweite)
      })
      return Array.from(reichweiten).sort()
    },
    availableWirkungsziele() {
      const wirkungsziele = new Set()
      this.enhancedSpells.forEach(spell => {
        if (spell.wirkungsziel) wirkungsziele.add(spell.wirkungsziel)
      })
      return Array.from(wirkungsziele).sort()
    },
    availableQuellen() {
      const quellen = new Set()
      this.enhancedSpells.forEach(spell => {
        if (spell.source_id && this.availableSources.length > 0) {
          const source = this.availableSources.find(s => s.id === spell.source_id)
          if (source) {
            quellen.add(source.code)
          }
        }
      })
      return Array.from(quellen).sort()
    },
    filteredAndSortedSpells() {
      let filtered = [...this.enhancedSpells]

      // Apply search filter
      if (this.searchTerm) {
        const searchLower = this.searchTerm.toLowerCase()
        filtered = filtered.filter(spell =>
          spell.name?.toLowerCase().includes(searchLower) ||
          spell.category?.toLowerCase().includes(searchLower)
        )
      }

      // Apply category filter
      if (this.filterCategory) {
        filtered = filtered.filter(spell => spell.category === this.filterCategory)
      }

      // Apply level filter
      if (this.filterLevel !== '') {
        filtered = filtered.filter(spell => spell.level === this.filterLevel)
      }

      // Apply ursprung filter
      if (this.filterUrsprung) {
        filtered = filtered.filter(spell => spell.ursprung === this.filterUrsprung)
      }

      // Apply Reichweite filter
      if (this.filterReichweite) {
        filtered = filtered.filter(spell => spell.reichweite === this.filterReichweite)
      }

      // Apply Wirkungsziel filter
      if (this.filterWirkungsziel) {
        filtered = filtered.filter(spell => spell.wirkungsziel === this.filterWirkungsziel)
      }

      // Apply Quelle filter (only by source code, ignoring page number)
      if (this.filterQuelle) {
        filtered = filtered.filter(spell => {
          if (spell.source_id && this.availableSources.length > 0) {
            const source = this.availableSources.find(s => s.id === spell.source_id)
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
    sortedSpells() {
      return [...this.mdata.spells].sort((a, b) => {
        const aValue = (a[this.sortField] || '').toLowerCase();
        const bValue = (b[this.sortField] || '').toLowerCase();
        return this.sortAsc
          ? aValue.localeCompare(bValue)
          : bValue.localeCompare(aValue);
      });
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
    async loadEnhancedSpells() {
      try {
        const response = await API.get('/api/maintenance/spells-enhanced')
        this.enhancedSpells = response.data.spells || []
        this.availableSources = response.data.sources || []
        // Also update mdata for compatibility
        if (response.data.categories) {
          this.mdata.spellcategories = response.data.categories
        }
      } catch (error) {
        console.error('Failed to load enhanced spells:', error)
      }
    },
    startEdit(index) {
      const spell = this.filteredAndSortedSpells[index]
      this.editedItem = {
        ...spell,
        sourceCode: this.getSourceCode(spell.source_id)
      }
      this.selectedSystemId = spell.game_system_id ?? this.findSystemIdByCode(spell.system)
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
          system: selectedSystem ? selectedSystem.code : (this.editedItem.system || ''),
          game_system_id: selectedSystem ? selectedSystem.id : (this.editedItem.game_system_id ?? null)
        }
        
        const response = await API.put(
          `/api/maintenance/spells-enhanced/${this.editedItem.id}`,
          updateData
        )

        // Update the spell in the list using splice for proper reactivity
        const spellIndex = this.enhancedSpells.findIndex(s => s.id === this.editedItem.id)
        if (spellIndex !== -1) {
          this.enhancedSpells.splice(spellIndex, 1, response.data)
        }

        this.editingIndex = -1
        this.editedItem = null
        this.selectedSystemId = null
      } catch (error) {
        console.error('Failed to save spell:', error)
        alert('Failed to save spell: ' + (error.response?.data?.error || error.message))
      }
    },
    cancelEdit() {
      this.editingIndex = -1;
      this.editedItem = null;
      this.selectedSystemId = null;
    },
    startCreate() {
      this.cancelEdit()
      const defaultSystem = this.gameSystems.find(gs => gs.is_active) || this.gameSystems[0] || null
      this.createSelectedSystemId = defaultSystem ? defaultSystem.id : null
      this.newItem = {
        name: '',
        category: this.mdata.spellcategories?.[0] || '',
        level: 0,
        ap: '',
        zauberdauer: '',
        reichweite: '',
        wirkungsziel: '',
        wirkungsbereich: '',
        wirkungsdauer: '',
        ursprung: '',
        beschreibung: '',
        sourceCode: '',
        page_number: 0,
        system: defaultSystem ? defaultSystem.code : ''
      }
      this.creatingNew = true
    },
    cancelCreate() {
      this.creatingNew = false
      this.newItem = null
      this.createSelectedSystemId = null
    },
    findSystemIdByCode(code) {
      return findSystemIdByCode(this.gameSystems, code)
    },
    sortBy(field) {
      if (this.sortField === field) {
        this.sortAsc = !this.sortAsc;
      } else {
        this.sortField = field;
        this.sortAsc = true;
      }
    },
    async handleSpellUpdate({ index, spell }) {
      try {
        const response = await API.put(
          `/api/maintenance/spells/${spell.id}`, spell,
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('token')}` ,
              'Content-Type': 'application/json'
            }
          }
        )
        if (!response.statusText== "OK") throw new Error('Update failed');
        const updatedSkill = response.data;
        // Update the spell in mdata
        this.mdata.spells = this.mdata.spells.map(s =>
          s.id === updatedSkill.id ? updatedSkill : s
        );
      } catch (error) {
        console.error('Failed to update spell:', error);
      }
    },
    async saveCreate() {
      if (!this.newItem) return
      try {
        const source = this.availableSources.find(s => s.code === this.newItem.sourceCode)
        const selectedSystem = this.gameSystems.find(gs => gs.id === this.createSelectedSystemId)

        const createData = {
          ...this.newItem,
          source_id: source ? source.id : null,
          page_number: this.newItem.page_number || 0,
          system: selectedSystem ? selectedSystem.code : (this.newItem.system || ''),
          game_system_id: selectedSystem ? selectedSystem.id : null
        }

        const response = await API.post(
          '/api/maintenance/spells-enhanced',
          createData
        )

        this.enhancedSpells.push(response.data)
        this.cancelCreate()
      } catch (error) {
        console.error('Failed to create spell:', error)
        alert('Failed to create spell: ' + (error.response?.data?.error || error.message))
      }
    },
    handleFileSelect(event) {
      const file = event.target.files[0];
      if (file && file.type === 'text/csv') {
        this.selectedFile = file;
        this.importResult = null;
      } else {
        this.selectedFile = null;
        this.importResult = {
          success: false,
          message: 'Please select a valid CSV file.'
        };
      }
    },
    async importSpells() {
      if (!this.selectedFile) {
        this.importResult = {
          success: false,
          message: 'Please select a CSV file first.'
        };
        return;
      }

      this.importing = true;
      this.importResult = null;

      try {
        const formData = new FormData();
        formData.append('file', this.selectedFile);

        const response = await API.post('/api/importer/spells/csv', formData, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
            'Content-Type': 'multipart/form-data'
          }
        });

        this.importResult = {
          success: true,
          message: response.data.message,
          total_spells: response.data.total_spells
        };

        // Refresh the spells data after successful import
        await this.refreshSpellsData();

      } catch (error) {
        console.error('Failed to import spells:', error);
        this.importResult = {
          success: false,
          message: error.response?.data?.message || 'Import failed. Please try again.'
        };
      } finally {
        this.importing = false;
      }
    },
    async refreshSpellsData() {
      try {
        const token = localStorage.getItem('token');
        const response = await API.get('/api/maintenance', {
          headers: { Authorization: `Bearer ${token}` }
        });
        
        // Update the spells data
        if (response.data.spells) {
          this.mdata.spells = response.data.spells;
        }
        if (response.data.sources) {
          this.availableSources = response.data.sources;
        }
      } catch (error) {
        console.error('Failed to refresh spells data:', error);
      }
    },
    formatQuelle(spell) {
      if (spell.source_id && this.availableSources.length > 0) {
        const source = this.availableSources.find(s => s.id === spell.source_id)
        if (source) {
          if (spell.page_number) {
            return `${source.code}:${spell.page_number}`
          } else {
            // No page number - show code and quelle if available
            const quelle = spell.quelle ? ` (${spell.quelle})` : ''
            return `${source.code}${quelle}`
          }
        }
      }
      return spell.quelle || '-'
    },
    getSourceCode(sourceId) {
      return getSourceCode(this.availableSources, sourceId)
    },
    getSystemCodeById(systemId, fallback = '') {
      return getSystemCodeById(this.gameSystems, systemId, fallback)
    },
    clearFilters() {
      this.searchTerm = ''
      this.filterCategory = ''
      this.filterLevel = ''
      this.filterUrsprung = ''
      this.filterReichweite = ''
      this.filterWirkungsziel = ''
      this.filterQuelle = ''
    }
  }
};
</script>
