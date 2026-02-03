<template>
  <div class="header-section">
    <h2>{{ $t('maintenance') }}</h2>
      <!-- Add search input -->
      <div class="search-box">
        <input
          type="text"
          v-model="searchTerm"
          :placeholder="`${$t('search')} ${$t('Equipment')}...`"
        />
      </div>
    </div>

  <div class="cd-view">
    <div class="cd-list">
      <!-- Filter Row -->
      <div class="filter-row">
        <div class="filter-item">
          <label>{{ $t('equipment.personal_item') }}:</label>
          <select v-model="filterPersonalItem">
            <option value="">{{ $t('all') || 'All' }}</option>
            <option value="true">{{ $t('yes') || 'Yes' }}</option>
            <option value="false">{{ $t('no') || 'No' }}</option>
          </select>
        </div>
        <div class="filter-item">
          <label>{{ $t('equipment.quelle') }}:</label>
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
                <th class="cd-table-header">{{ $t('equipment.id') }}</th>
                <th class="cd-table-header">
                  {{ $t('equipment.name') }}
                  <button @click="sortBy('name')">{{ sortField === 'name' ? (sortAsc ? '↑' : '↓') : '-' }}</button>
                </th>
                <th class="cd-table-header">{{ $t('equipment.gewicht') }}</th>
                <th class="cd-table-header">{{ $t('equipment.wert') }}</th>
                <th class="cd-table-header">{{ $t('equipment.description') }}</th>
                <th class="cd-table-header">{{ $t('equipment.quelle') }}</th>
                <th class="cd-table-header">{{ $t('equipment.personal_item') }}</th>
                <th class="cd-table-header">{{ $t('equipment.system') }}</th>
                <th class="cd-table-header"> </th>
              </tr>
            </thead>
            <tbody>
              <template v-for="(dtaItem, index) in filteredAndSortedEquipments" :key="dtaItem.id">
                <tr v-if="editingIndex !== index">
                  <td>{{ dtaItem.id || '' }}</td>
                  <!-- <td>{{ dtaItem.category|| '-' }}</td> -->
                  <td>{{ dtaItem.name || '-' }}</td>
                  <td>{{ dtaItem.gewicht || '-' }}</td>
                  <td>{{ dtaItem.wert || '-' }}</td>
                  <td>{{ dtaItem.beschreibung || '-' }}</td>
                  <td>{{ formatQuelle(dtaItem) }}</td>
                  <td><input type="checkbox" :checked="dtaItem.personal_item" disabled /></td>
                  <td>{{ getSystemCodeById(dtaItem.game_system_id, dtaItem.system || 'midgard') }}</td>
                  <td>
                    <button @click="startEdit(index)">Edit</button>
                  </td>
                </tr>
                <!-- Edit Mode -->
                <tr v-else>
                  <td><input v-model="editedItem.id" style="width:20px;" disabled /></td>
                  <td colspan="8">
                    <!-- Expanded edit form -->
                    <div class="edit-form">
                      <div class="edit-row">
                        <div class="edit-field">
                          <label>{{ $t('equipment.name') }}:</label>
                          <input v-model="editedItem.name" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('equipment.gewicht') }}:</label>
                          <input v-model.number="editedItem.gewicht" type="number" style="width:80px;" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('equipment.wert') }}:</label>
                          <input v-model="editedItem.wert" style="width:100px;" />
                        </div>
                      </div>

                      <div class="edit-row">
                        <div class="edit-field full-width">
                          <label>{{ $t('equipment.description') }}:</label>
                          <input v-model="editedItem.beschreibung" />
                        </div>
                      </div>

                      <div class="edit-row">
                        <div class="edit-field">
                          <label>{{ $t('equipment.quelle') }}:</label>
                          <select v-model="editedItem.sourceCode" style="width:100px;">
                            <option value="">-</option>
                            <option v-for="source in availableSources" :key="source.code" :value="source.code">
                              {{ source.code }}
                            </option>
                          </select>
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('equipment.page') || 'Page' }}:</label>
                          <input v-model.number="editedItem.page_number" type="number" style="width:60px;" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('equipment.personal_item') }}:</label>
                          <input type="checkbox" v-model="editedItem.personal_item" />
                        </div>
                        <div class="edit-field">
                          <label>{{ $t('equipment.system') }}:</label>
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
  systemOptionsFor,
} from '../../utils/maintenanceGameSystems'
export default {
  name: "EquipmentView",
  props: {
    mdata: {
      type: Object,
      required: true,
      default: () => ({
        equipments: [],
        equipmentcategories: []
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
      filterPersonalItem: '',
      filterQuelle: '',
      enhancedEquipment: [],
      availableSources: [],
      gameSystems: [],
      selectedSystemId: null
    }
  },
  async created() {
    await Promise.all([
      this.loadGameSystems(),
      this.loadEnhancedEquipment()
    ])
  },
  computed: {
    availableQuellen() {
      const quellen = new Set()
      this.enhancedEquipment.forEach(equipment => {
        if (equipment.source_id && this.availableSources.length > 0) {
          const source = this.availableSources.find(s => s.id === equipment.source_id)
          if (source) {
            quellen.add(source.code)
          }
        }
      })
      return Array.from(quellen).sort()
    },
    filteredAndSortedEquipments() {
      let filtered = [...this.enhancedEquipment]

      // Apply search filter
      if (this.searchTerm) {
        const searchLower = this.searchTerm.toLowerCase()
        filtered = filtered.filter(equipment =>
          equipment.name?.toLowerCase().includes(searchLower)
        )
      }

      // Apply personal_item filter
      if (this.filterPersonalItem !== '') {
        const personalItemValue = this.filterPersonalItem === 'true'
        filtered = filtered.filter(equipment => equipment.personal_item === personalItemValue)
      }

      // Apply Quelle filter (only by source code, ignoring page number)
      if (this.filterQuelle) {
        filtered = filtered.filter(equipment => {
          if (equipment.source_id && this.availableSources.length > 0) {
            const source = this.availableSources.find(s => s.id === equipment.source_id)
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
    sortedEquipments() {
      return [...this.mdata.equipment].sort((a, b) => {
        const aValue = (a[this.sortField] || '').toLowerCase();
        const bValue = (b[this.sortField] || '').toLowerCase();
        return this.sortAsc
          ? aValue.localeCompare(bValue)
          : bValue.localeCompare(aValue);
      });
    },
    systemOptions() {
      return systemOptionsFor(this.gameSystems)
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
    async loadEnhancedEquipment() {
      try {
        const response = await API.get('/api/maintenance/equipment-enhanced')
        this.enhancedEquipment = response.data.equipment || []
        this.availableSources = response.data.sources || []
      } catch (error) {
        console.error('Failed to load enhanced equipment:', error)
      }
    },
    startEdit(index) {
      const equipment = this.filteredAndSortedEquipments[index]
      this.editedItem = {
        ...equipment,
        sourceCode: this.getSourceCode(equipment.source_id)
      }
      this.selectedSystemId = equipment.game_system_id ?? this.findSystemIdByCode(equipment.system)
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
          `/api/maintenance/equipment-enhanced/${this.editedItem.id}`,
          updateData
        )

        // Update the equipment in the list using splice for proper reactivity
        const equipmentIndex = this.enhancedEquipment.findIndex(e => e.id === this.editedItem.id)
        if (equipmentIndex !== -1) {
          this.enhancedEquipment.splice(equipmentIndex, 1, response.data)
        }

        this.editingIndex = -1
        this.editedItem = null
        this.selectedSystemId = null
      } catch (error) {
        console.error('Failed to save equipment:', error)
        alert('Failed to save equipment: ' + (error.response?.data?.error || error.message))
      }
    },
    cancelEdit() {
      this.editingIndex = -1;
      this.editedItem = null;
      this.selectedSystemId = null;
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
    formatQuelle(equipment) {
      if (equipment.source_id && this.availableSources.length > 0) {
        const source = this.availableSources.find(s => s.id === equipment.source_id)
        if (source) {
          if (equipment.page_number) {
            return `${source.code}:${equipment.page_number}`
          } else {
            // No page number - show code and quelle if available
            const quelle = equipment.quelle ? ` (${equipment.quelle})` : ''
            return `${source.code}${quelle}`
          }
        }
      }
      return equipment.quelle || '-'
    },
    getSourceCode(sourceId) {
      return getSourceCode(this.availableSources, sourceId)
    },
    clearFilters() {
      this.searchTerm = ''
      this.filterPersonalItem = ''
      this.filterQuelle = ''
    },
    getSystemCodeById(systemId, fallback = '') {
      return getSystemCodeById(this.gameSystems, systemId, fallback)
    },
    async handleEquipmentUpdate({ index, equipment }) {
      try {
          const response = await API.put(
            `/api/maintenance/equipment/${equipment.id}`, equipment,
            {
              headers: {
                Authorization: `Bearer ${localStorage.getItem('token')}` ,
                'Content-Type': 'application/json'
              }
            }
          )
          if (response.status !== 200) throw new Error('Update failed');
          const updatedSkill = response.data;
          // Update the equipment in mdata
          this.mdata.equipment = this.mdata.equipment.map(s =>
            s.id === updatedSkill.id ? updatedSkill : s
          );
        } catch (error) {
          console.error('Failed to update equipment:', error);
        }
      }
  }
};
</script>
