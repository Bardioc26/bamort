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
          <label>{{ $t('weaponskill.category') }}:</label>
          <select v-model="filterCategory">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option v-for="cat in availableCategories" :key="cat" :value="cat">{{ cat }}</option>
          </select>
        </div>
        <div class="filter-item">
          <label>{{ $t('weaponskill.bonusskill') }}:</label>
          <select v-model="filterBonuseigenschaft">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option v-for="bonus in availableBonuseigenschaften" :key="bonus" :value="bonus">{{ bonus }}</option>
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
                <th class="cd-table-header">{{ $t('weaponskill.category') }}<button @click="sortBy('category')">-{{ sortField === 'category' ? (sortAsc ? '↑' : '↓') : '' }}</button></th>
                <th class="cd-table-header">{{ $t('weaponskill.name') }} <button @click="sortBy('name')">-{{ sortField === 'name' ? (sortAsc ? '↑' : '↓') : '' }}</button></th>
                <th class="cd-table-header">{{ $t('weaponskill.initialwert') }}</th>
                <th class="cd-table-header">{{ $t('weaponskill.description') }}</th>
                <th class="cd-table-header">{{ $t('weaponskill.bonusskill') }}</th>
                <th class="cd-table-header">{{ $t('weaponskill.quelle') }}</th>
                <th class="cd-table-header">{{ $t('weaponskill.system') }}</th>
                <th class="cd-table-header"> </th>
              </tr>
            </thead>
            <tbody>
              <template v-for="(dtaItem, index) in filteredAndSortedWaeponSkills" :key="dtaItem.id">
                <tr v-if="editingIndex !== index">
                  <td>{{ dtaItem.id || '' }}</td>
                  <td>{{ dtaItem.category|| '-' }}</td>
                  <td>{{ dtaItem.name || '-' }}</td>
                  <td>{{ dtaItem.initialwert || '0' }}</td>
                  <td>{{ dtaItem.beschreibung || '-' }}</td>
                  <td>{{ dtaItem.bonuseigenschaft || '-' }}</td>
                  <td>{{ formatQuelle(dtaItem) }}</td>
                  <td>{{ dtaItem.system || 'midgard' }}</td>
                  <td>
                    <button @click="startEdit(index)">Edit</button>
                  </td>
                </tr>
                <!-- Edit Mode -->
                <tr v-else>
                  <td><input v-model="editedItem.id" style="width:20px;" disabled /></td>
                  <td colspan="7">
                    <!-- Expanded edit form -->
                    <div class="edit-form">
                      <div class="edit-row">
                        <div class="edit-field">
                          <label>{{ $t('weaponskill.name') }}:</label>
                          <input v-model="editedItem.name" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('weaponskill.category') }}:</label>
                          <select v-model="editedItem.category" style="width:120px;">
                            <option v-for="category in mdata['skillcategories']" :key="category" :value="category">
                              {{ category }}
                            </option>
                          </select>
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('weaponskill.initialwert') }}:</label>
                          <input v-model.number="editedItem.initialwert" type="number" style="width:60px;" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('weaponskill.bonusskill') }}:</label>
                          <select v-model="editedItem.bonuseigenschaft" style="width:80px;">
                            <option value="">-</option>
                            <option value="St">St</option>
                            <option value="Gs">Gs</option>
                            <option value="Gw">Gw</option>
                            <option value="Ko">Ko</option>
                            <option value="In">In</option>
                            <option value="Zt">Zt</option>
                            <option value="Au">Au</option>
                            <option value="pA">pA</option>
                            <option value="Wk">Wk</option>
                            <option value="B">B</option>
                          </select>
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
                          <input v-model="editedItem.system" style="width:100px;" />
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
      filterCategory: '',
      filterBonuseigenschaft: '',
      filterQuelle: '',
      enhancedWeaponSkills: [],
      availableSources: []
    }
  },
  async created() {
    await this.loadEnhancedWeaponSkills()
  },
  computed: {
    availableCategories() {
      const categories = new Set()
      this.enhancedWeaponSkills.forEach(ws => {
        if (ws.category) categories.add(ws.category)
      })
      return Array.from(categories).sort()
    },
    availableBonuseigenschaften() {
      const bonuses = new Set()
      this.enhancedWeaponSkills.forEach(ws => {
        if (ws.bonuseigenschaft) bonuses.add(ws.bonuseigenschaft)
      })
      return Array.from(bonuses).sort()
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
          ws.category?.toLowerCase().includes(searchLower)
        )
      }

      // Apply category filter
      if (this.filterCategory) {
        filtered = filtered.filter(ws => ws.category === this.filterCategory)
      }

      // Apply bonuseigenschaft filter
      if (this.filterBonuseigenschaft) {
        filtered = filtered.filter(ws => ws.bonuseigenschaft === this.filterBonuseigenschaft)
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
    }
  },
  methods: {
    async loadEnhancedWeaponSkills() {
      try {
        const response = await API.get('/api/maintenance/weaponskills-enhanced')
        this.enhancedWeaponSkills = response.data.weaponskills || []
        this.availableSources = response.data.sources || []
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
      this.editingIndex = index
    },
    async saveEdit(index) {
      try {
        // Find source ID from code
        const source = this.availableSources.find(s => s.code === this.editedItem.sourceCode)
        
        const updateData = {
          ...this.editedItem,
          source_id: source ? source.id : null,
          page_number: this.editedItem.page_number || 0
        }
        
        const response = await API.put(
          `/api/maintenance/weaponskills-enhanced/${this.editedItem.id}`,
          updateData
        )

        // Update the weapon skill in the list using splice for proper reactivity
        const wsIndex = this.enhancedWeaponSkills.findIndex(ws => ws.id === this.editedItem.id)
        if (wsIndex !== -1) {
          this.enhancedWeaponSkills.splice(wsIndex, 1, response.data)
        }

        this.editingIndex = -1
        this.editedItem = null
      } catch (error) {
        console.error('Failed to save weapon skill:', error)
        alert('Failed to save weapon skill: ' + (error.response?.data?.error || error.message))
      }
    },
    cancelEdit() {
      this.editingIndex = -1
      this.editedItem = null
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
      if (!sourceId || !this.availableSources.length) return ''
      const source = this.availableSources.find(s => s.id === sourceId)
      return source ? source.code : ''
    },
    clearFilters() {
      this.searchTerm = ''
      this.filterCategory = ''
      this.filterBonuseigenschaft = ''
      this.filterQuelle = ''
    }
  }
};
</script>
