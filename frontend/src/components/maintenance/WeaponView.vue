<template>
  <div class="header-section">
    <h2>{{ $t('maintenance') }}</h2>
      <!-- Add search input -->
      <div class="search-box">
        <input
          type="text"
          v-model="searchTerm"
          :placeholder="`${$t('search')} ${$t('Weapons')}...`"
        />
        <button @click="startCreate" class="btn-primary">{{ $t('newEntry') }}</button>
      </div>
    </div>

  <div class="cd-view">
    <div class="cd-list">
      <!-- Filter Row -->
      <div class="filter-row">
        <div class="filter-item">
          <label>{{ $t('weapon.skillrequired') }}:</label>
          <select v-model="filterSkillRequired">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option v-for="skill in availableSkillsRequired" :key="skill" :value="skill">{{ skill }}</option>
          </select>
        </div>
        <div class="filter-item">
          <label>{{ $t('weapon.damage') }}:</label>
          <select v-model="filterDamage">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option v-for="dmg in availableDamages" :key="dmg" :value="dmg">{{ dmg }}</option>
          </select>
        </div>
        <div class="filter-item">
          <label>{{ $t('weapon.rangenear') }}:</label>
          <select v-model="filterRangeNear">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option v-for="range in availableRangesNear" :key="range" :value="range">{{ range }}</option>
          </select>
        </div>
        <div class="filter-item">
          <label>{{ $t('weapon.rangemiddle') }}:</label>
          <select v-model="filterRangeMiddle">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option v-for="range in availableRangesMiddle" :key="range" :value="range">{{ range }}</option>
          </select>
        </div>
        <div class="filter-item">
          <label>{{ $t('weapon.rangefar') }}:</label>
          <select v-model="filterRangeFar">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option v-for="range in availableRangesFar" :key="range" :value="range">{{ range }}</option>
          </select>
        </div>
        <div class="filter-item">
          <label>{{ $t('weapon.quelle') }}:</label>
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
                <th class="cd-table-header">{{ $t('weapon.id') }}</th>
                <!-- <th class="cd-table-header">{{ $t('weapon.category') }}<button @click="sortBy('category')">-{{ sortField === 'category' ? (sortAsc ? '↑' : '↓') : '' }}</button></th> -->
                <th class="cd-table-header">{{ $t('weapon.name') }} <button @click="sortBy('name')">-{{ sortField === 'name' ? (sortAsc ? '↑' : '↓') : '' }}</button></th>
                <th class="cd-table-header">{{ $t('weapon.skillrequired') || 'Skill Required' }}</th>
                <th class="cd-table-header">{{ $t('weapon.weight') }}</th>
                <th class="cd-table-header">{{ $t('weapon.value') }}</th>
                <th class="cd-table-header">{{ $t('weapon.damage') }}</th>
                <th class="cd-table-header">{{ $t('weapon.rangenear') || 'Range Near' }}</th>
                <th class="cd-table-header">{{ $t('weapon.rangemiddle') || 'Range Middle' }}</th>
                <th class="cd-table-header">{{ $t('weapon.rangefar') || 'Range Far' }}</th>
                <th class="cd-table-header">{{ $t('weapon.description') }}</th>
                <th class="cd-table-header">{{ $t('weapon.quelle') }}</th>
                <th class="cd-table-header">{{ $t('weapon.personal_item') }}</th>
                <th class="cd-table-header">{{ $t('weapon.system') }}</th>
                <th class="cd-table-header"> </th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="creatingNew">
                <td>New</td>
                <td colspan="11">
                  <div class="edit-form">
                    <div class="edit-row">
                      <div class="edit-field">
                        <label>{{ $t('weapon.name') }}:</label>
                        <input v-model="newItem.name" />
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('weapon.skillrequired') || 'Skill Required' }}:</label>
                        <input v-model="newItem.skill_required" style="width:150px;" />
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('weapon.weight') }}:</label>
                        <input v-model.number="newItem.gewicht" type="number" style="width:80px;" />
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('weapon.value') }}:</label>
                        <input v-model="newItem.wert" style="width:100px;" />
                      </div>
                    </div>

                    <div class="edit-row">
                      <div class="edit-field">
                        <label>{{ $t('weapon.damage') }}:</label>
                        <input v-model="newItem.damage" style="width:100px;" />
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('weapon.rangenear') || 'Range Near' }}:</label>
                        <input v-model.number="newItem.range_near" type="number" style="width:80px;" />
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('weapon.rangemiddle') || 'Range Middle' }}:</label>
                        <input v-model.number="newItem.range_middle" type="number" style="width:80px;" />
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('weapon.rangefar') || 'Range Far' }}:</label>
                        <input v-model.number="newItem.range_far" type="number" style="width:80px;" />
                      </div>
                    </div>

                    <div class="edit-row">
                      <div class="edit-field">
                        <label>{{ $t('weapon.bonusskill') || 'Bonus' }}:</label>
                        <select v-model="newItem.bonuseigenschaft" style="width:80px;">
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
                      <div class="edit-field">
                        <label>{{ $t('weapon.personal_item') }}:</label>
                        <input type="checkbox" v-model="newItem.personal_item" />
                      </div>
                    </div>

                    <div class="edit-row">
                      <div class="edit-field full-width">
                        <label>{{ $t('weapon.description') }}:</label>
                        <input v-model="newItem.beschreibung" />
                      </div>
                    </div>

                    <div class="edit-row">
                      <div class="edit-field">
                        <label>{{ $t('weapon.quelle') }}:</label>
                        <select v-model="newItem.sourceCode" style="width:100px;">
                          <option value="">-</option>
                          <option v-for="source in availableSources" :key="source.code" :value="source.code">
                            {{ source.code }}
                          </option>
                        </select>
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('weapon.page') || 'Page' }}:</label>
                        <input v-model.number="newItem.page_number" type="number" style="width:60px;" />
                      </div>
                      <div class="edit-field">
                        <label>{{ $t('weapon.system') }}:</label>
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
              <template v-for="(dtaItem, index) in filteredAndSortedWeaponss" :key="dtaItem.id">
                <tr v-if="editingIndex !== index">
                  <td>{{ dtaItem.id || '' }}</td>
                  <td>{{ dtaItem.name || '-' }}</td>
                  <td>{{ dtaItem.skill_required || '-' }}</td>
                  <td>{{ dtaItem.gewicht || '-' }}</td>
                  <td>{{ dtaItem.wert || '-' }}</td>
                  <td>{{ dtaItem.damage || '-' }}</td>
                  <td>{{ dtaItem.range_near || '-' }}</td>
                  <td>{{ dtaItem.range_middle || '-' }}</td>
                  <td>{{ dtaItem.range_far || '-' }}</td>
                  <td>{{ dtaItem.beschreibung || '-' }}</td>
                  <td>{{ formatQuelle(dtaItem) }}</td>
                  <td><input type="checkbox" :checked="dtaItem.personal_item" disabled /></td>
                  <td>{{ getSystemCodeById(dtaItem.game_system_id, dtaItem.system || 'midgard') }}</td>
                  <td>
                    <button @click="startEdit(index)">{{ $t('common.edit') }}</button>
                  </td>
                </tr>
                <!-- Edit Mode -->
                <tr v-else>
                  <td><input v-model="editedItem.id" style="width:20px;" disabled /></td>
                  <td colspan="11">
                    <!-- Expanded edit form -->
                    <div class="edit-form">
                      <div class="edit-row">
                        <div class="edit-field">
                          <label>{{ $t('weapon.name') }}:</label>
                          <input v-model="editedItem.name" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('weapon.skillrequired') || 'Skill Required' }}:</label>
                          <input v-model="editedItem.skill_required" style="width:150px;" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('weapon.weight') }}:</label>
                          <input v-model.number="editedItem.gewicht" type="number" style="width:80px;" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('weapon.value') }}:</label>
                          <input v-model="editedItem.wert" style="width:100px;" />
                        </div>
                      </div>

                      <div class="edit-row">
                        <div class="edit-field">
                          <label>{{ $t('weapon.damage') }}:</label>
                          <input v-model="editedItem.damage" style="width:100px;" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('weapon.rangenear') || 'Range Near' }}:</label>
                          <input v-model.number="editedItem.range_near" type="number" style="width:80px;" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('weapon.rangemiddle') || 'Range Middle' }}:</label>
                          <input v-model.number="editedItem.range_middle" type="number" style="width:80px;" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('weapon.rangefar') || 'Range Far' }}:</label>
                          <input v-model.number="editedItem.range_far" type="number" style="width:80px;" />
                        </div>
                      </div>

                      <div class="edit-row">
                        <div class="edit-field">
                          <label>{{ $t('weapon.bonusskill') || 'Bonus' }}:</label>
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
                        <div class="edit-field">
                          <label>{{ $t('weapon.personal_item') }}:</label>
                          <input type="checkbox" v-model="editedItem.personal_item" />
                        </div>
                      </div>

                      <div class="edit-row">
                        <div class="edit-field full-width">
                          <label>{{ $t('weapon.description') }}:</label>
                          <input v-model="editedItem.beschreibung" />
                        </div>
                      </div>

                      <div class="edit-row">
                        <div class="edit-field">
                          <label>{{ $t('weapon.quelle') }}:</label>
                          <select v-model="editedItem.sourceCode" style="width:100px;">
                            <option value="">-</option>
                            <option v-for="source in availableSources" :key="source.code" :value="source.code">
                              {{ source.code }}
                            </option>
                          </select>
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('weapon.page') || 'Page' }}:</label>
                          <input v-model.number="editedItem.page_number" type="number" style="width:60px;" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('weapon.system') }}:</label>
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
  name: "WeaponView",
  props: {
    mdata: {
      type: Object,
      required: true,
      default: () => ({
        weaponss: [],
        weaponscategories: []
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
      filterSkillRequired: '',
      filterDamage: '',
      filterRangeNear: '',
      filterRangeMiddle: '',
      filterRangeFar: '',
      filterQuelle: '',
      enhancedWeapons: [],
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
      this.loadEnhancedWeapons()
    ])
  },
  computed: {
    availableSkillsRequired() {
      const skills = new Set()
      this.enhancedWeapons.forEach(w => {
        if (w.skill_required) skills.add(w.skill_required)
      })
      return Array.from(skills).sort()
    },
    availableDamages() {
      const damages = new Set()
      this.enhancedWeapons.forEach(w => {
        if (w.damage) damages.add(w.damage)
      })
      return Array.from(damages).sort()
    },
    availableRangesNear() {
      const ranges = new Set()
      this.enhancedWeapons.forEach(w => {
        if (w.range_near !== null && w.range_near !== undefined) ranges.add(w.range_near)
      })
      return Array.from(ranges).sort((a, b) => a - b)
    },
    availableRangesMiddle() {
      const ranges = new Set()
      this.enhancedWeapons.forEach(w => {
        if (w.range_middle !== null && w.range_middle !== undefined) ranges.add(w.range_middle)
      })
      return Array.from(ranges).sort((a, b) => a - b)
    },
    availableRangesFar() {
      const ranges = new Set()
      this.enhancedWeapons.forEach(w => {
        if (w.range_far !== null && w.range_far !== undefined) ranges.add(w.range_far)
      })
      return Array.from(ranges).sort((a, b) => a - b)
    },
    availableQuellen() {
      const quellen = new Set()
      this.enhancedWeapons.forEach(w => {
        if (w.source_id && this.availableSources.length > 0) {
          const source = this.availableSources.find(s => s.id === w.source_id)
          if (source) {
            quellen.add(source.code)
          }
        }
      })
      return Array.from(quellen).sort()
    },
    filteredAndSortedWeaponss() {
      let filtered = [...this.enhancedWeapons]

      // Apply search filter
      if (this.searchTerm) {
        const searchLower = this.searchTerm.toLowerCase()
        filtered = filtered.filter(w =>
          w.name?.toLowerCase().includes(searchLower) ||
          w.skill_required?.toLowerCase().includes(searchLower)
        )
      }

      // Apply skill_required filter
      if (this.filterSkillRequired) {
        filtered = filtered.filter(w => w.skill_required === this.filterSkillRequired)
      }

      // Apply damage filter
      if (this.filterDamage) {
        filtered = filtered.filter(w => w.damage === this.filterDamage)
      }

      // Apply range_near filter
      if (this.filterRangeNear !== '') {
        filtered = filtered.filter(w => w.range_near === this.filterRangeNear)
      }

      // Apply range_middle filter
      if (this.filterRangeMiddle !== '') {
        filtered = filtered.filter(w => w.range_middle === this.filterRangeMiddle)
      }

      // Apply range_far filter
      if (this.filterRangeFar !== '') {
        filtered = filtered.filter(w => w.range_far === this.filterRangeFar)
      }

      // Apply Quelle filter (only by source code, ignoring page number)
      if (this.filterQuelle) {
        filtered = filtered.filter(w => {
          if (w.source_id && this.availableSources.length > 0) {
            const source = this.availableSources.find(s => s.id === w.source_id)
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
    async loadEnhancedWeapons() {
      try {
        const response = await API.get('/api/maintenance/weapons-enhanced')
        this.enhancedWeapons = response.data.weapons || []
        this.availableSources = response.data.sources || []
      } catch (error) {
        console.error('Failed to load enhanced weapons:', error)
      }
    },
    startEdit(index) {
      const weapon = this.filteredAndSortedWeaponss[index]
      this.editedItem = {
        ...weapon,
        sourceCode: this.getSourceCode(weapon.source_id)
      }
      this.selectedSystemId = weapon.game_system_id ?? this.findSystemIdByCode(weapon.system)
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
          `/api/maintenance/weapons-enhanced/${this.editedItem.id}`,
          updateData
        )

        // Update the weapon in the list using splice for proper reactivity
        const weaponIndex = this.enhancedWeapons.findIndex(w => w.id === this.editedItem.id)
        if (weaponIndex !== -1) {
          this.enhancedWeapons.splice(weaponIndex, 1, response.data)
        }

        this.editingIndex = -1
        this.editedItem = null
        this.selectedSystemId = null
      } catch (error) {
        console.error('Failed to save weapon:', error)
        alert('Failed to save weapon: ' + (error.response?.data?.error || error.message))
      }
    },
    cancelEdit() {
      this.editingIndex = -1
      this.editedItem = null
      this.selectedSystemId = null
    },
    startCreate() {
      this.cancelEdit()
      const defaultSystem = this.gameSystems.find(gs => gs.is_active) || this.gameSystems[0] || null
      this.createSelectedSystemId = defaultSystem ? defaultSystem.id : null
      this.newItem = {
        name: '',
        skill_required: '',
        gewicht: 0,
        wert: '',
        damage: '',
        range_near: 0,
        range_middle: 0,
        range_far: 0,
        beschreibung: '',
        bonuseigenschaft: '',
        personal_item: false,
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
          '/api/maintenance/weapons-enhanced',
          createData
        )

        this.enhancedWeapons.push(response.data)
        this.cancelCreate()
      } catch (error) {
        console.error('Failed to create weapon:', error)
        alert('Failed to create weapon: ' + (error.response?.data?.error || error.message))
      }
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
    formatQuelle(weapon) {
      if (weapon.source_id && this.availableSources.length > 0) {
        const source = this.availableSources.find(s => s.id === weapon.source_id)
        if (source) {
          if (weapon.page_number) {
            return `${source.code}:${weapon.page_number}`
          } else {
            // No page number - show code and quelle if available
            const quelle = weapon.quelle ? ` (${weapon.quelle})` : ''
            return `${source.code}${quelle}`
          }
        }
      }
      return weapon.quelle || '-'
    },
    getSourceCode(sourceId) {
      return getSourceCode(this.availableSources, sourceId)
    },
    getSystemCodeById(systemId, fallback = '') {
      return getSystemCodeById(this.gameSystems, systemId, fallback)
    },
    clearFilters() {
      this.searchTerm = ''
      this.filterSkillRequired = ''
      this.filterDamage = ''
      this.filterRangeNear = ''
      this.filterRangeMiddle = ''
      this.filterRangeFar = ''
      this.filterQuelle = ''
    }
  }
};
</script>
