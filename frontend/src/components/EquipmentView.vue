<template>
  <div class="fullwidth-container">
    <div class="header-section">
      <h2>{{ $t('EquipmentView') }}</h2>
      <button @click="openAddEquipmentDialog" class="btn-add-equipment">
        {{ $t('equipment.add') }}
      </button>
    </div>

    <div class="cd-list">
      <table class="cd-table">
      <thead>
        <tr>
          <th>{{ $t('equipment.name') }}</th>
          <th>{{ $t('equipment.description') }}</th>
          <th>{{ $t('equipment.weight') }}</th>
          <th>{{ $t('equipment.value') }}</th>
          <th>{{ $t('equipment.amount') }}</th>
          <th>{{ $t('equipment.contained_in') }}</th>
          <th>{{ $t('equipment.bonus') }}</th>
          <th>{{ $t('equipment.actions') }}</th>
        </tr>
      </thead>
      <tbody>
      <template v-if="character.ausruestung && character.ausruestung.length > 0">
        <tr v-for="equipment in character.ausruestung" :key="equipment.id">
          <td>{{ equipment.name || '-' }}</td>
          <td>{{ equipment.beschreibung || '-' }}</td>
          <td>{{ equipment.gewicht || '-' }}</td>
          <td>{{ equipment.wert || '-' }}</td>
          <td>{{ equipment.anzahl || '-' }}</td>
          <td>{{ equipment.beinhaltet_in || '-' }}</td>
          <td>{{ equipment.bonus || '-' }}</td>
          <td class="action-cell">
            <button @click="deleteEquipment(equipment)" class="btn-delete" title="Löschen">
              🗑️
            </button>
          </td>
        </tr>
      </template>
      <template v-else>
        <tr>
          <td colspan="8" class="empty-state">{{ $t('equipment.noEquipment') }}</td>
        </tr>
      </template>
      </tbody>
      </table>
    </div>

    <!-- Dialog für Ausrüstung hinzufügen -->
    <div v-if="showAddDialog" class="modal-overlay" @click.self="closeDialog">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ $t('equipment.addEquipment') }}</h3>
          <button @click="closeDialog" class="close-button">&times;</button>
        </div>
        
        <div class="modal-body">
          <div class="form-group">
            <label>{{ $t('equipment.search') }}:</label>
            <input 
              v-model="searchQuery" 
              type="text" 
              :placeholder="$t('equipment.searchPlaceholder')"
              @input="filterEquipment"
            />
          </div>

          <div v-if="isLoading" class="loading">{{ $t('equipment.loading') }}</div>

          <div v-else-if="filteredMasterEquipment.length > 0" class="equipment-list">
            <div 
              v-for="equipment in filteredMasterEquipment" 
              :key="equipment.id"
              class="equipment-item"
              :class="{ selected: selectedEquipment && selectedEquipment.id === equipment.id }"
              @click="selectEquipment(equipment)"
            >
              <div class="equipment-name">{{ equipment.name }}</div>
              <div class="equipment-details">
                {{ equipment.beschreibung || '-' }} | 
                {{ equipment.weight }}kg | 
                {{ equipment.value || '-' }}
              </div>
            </div>
          </div>

          <div v-else class="no-results">
            {{ searchQuery ? $t('equipment.noResults') : $t('equipment.noMasterData') }}
          </div>

          <div v-if="selectedEquipment" class="selected-equipment-details">
            <h4>{{ $t('equipment.selectedEquipment') }}</h4>
            <p><strong>{{ $t('equipment.name') }}:</strong> {{ selectedEquipment.name }}</p>
            <p><strong>{{ $t('equipment.description') }}:</strong> {{ selectedEquipment.beschreibung || '-' }}</p>
            <p><strong>{{ $t('equipment.weight') }}:</strong> {{ selectedEquipment.weight }}kg</p>
            <p><strong>{{ $t('equipment.value') }}:</strong> {{ selectedEquipment.value || '-' }}</p>
            
            <div class="form-group">
              <label>{{ $t('equipment.amount') }}:</label>
              <input v-model.number="equipmentAmount" type="number" min="1" value="1" />
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button @click="closeDialog" class="btn-cancel">{{ $t('equipment.cancel') }}</button>
          <button 
            @click="addEquipment" 
            class="btn-confirm" 
            :disabled="!selectedEquipment || isSubmitting"
          >
            {{ isSubmitting ? $t('equipment.adding') : $t('equipment.add') }}
          </button>
        </div>
      </div>
    </div>
  </div>

</template>

<style>
/* All styles moved to main.css */

/* Equipment-specific table header color override */
.cd-table th {
  background-color: #1da766;
}
</style>

<script>
import API from '@/utils/api'

export default {
name: "EquipmentView",
  props: {
    character: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      showAddDialog: false,
      isLoading: false,
      isSubmitting: false,
      masterEquipment: [],
      filteredMasterEquipment: [],
      selectedEquipment: null,
      searchQuery: '',
      equipmentAmount: 1
    }
  },
  created() {
    this.$api = API
  },
  methods: {
    async openAddEquipmentDialog() {
      this.showAddDialog = true
      this.isLoading = true
      
      try {
        const response = await this.$api.get('/api/maintenance/equipment')
        this.masterEquipment = response.data || []
        this.filteredMasterEquipment = this.masterEquipment
      } catch (error) {
        console.error('Fehler beim Laden der Ausrüstungs-Stammdaten:', error)
        alert(this.$t('equipment.loadError') + ': ' + (error.response?.data?.error || error.message))
      } finally {
        this.isLoading = false
      }
    },
    
    filterEquipment() {
      const query = this.searchQuery.toLowerCase().trim()
      if (!query) {
        this.filteredMasterEquipment = this.masterEquipment
      } else {
        this.filteredMasterEquipment = this.masterEquipment.filter(equipment => 
          equipment.name.toLowerCase().includes(query) ||
          (equipment.beschreibung && equipment.beschreibung.toLowerCase().includes(query))
        )
      }
    },
    
    selectEquipment(equipment) {
      this.selectedEquipment = equipment
    },
    
    async addEquipment() {
      if (!this.selectedEquipment) {
        alert(this.$t('equipment.pleaseSelect'))
        return
      }
      
      this.isSubmitting = true
      
      try {
        const equipmentData = {
          character_id: this.character.id,
          name: this.selectedEquipment.name,
          beschreibung: this.selectedEquipment.beschreibung || '',
          gewicht: this.selectedEquipment.weight || 0,
          wert: this.selectedEquipment.value || 0,
          anzahl: this.equipmentAmount,
          bonus: this.selectedEquipment.bonus || 0,
          beinhaltet_in: '',
          contained_in: 0,
          container_type: ''
        }
        
        await this.$api.post('/api/equipment', equipmentData)
        
        alert(this.$t('equipment.addSuccess'))
        this.closeDialog()
        this.$emit('character-updated')
        
      } catch (error) {
        console.error('Fehler beim Hinzufügen der Ausrüstung:', error)
        alert(this.$t('equipment.addError') + ': ' + (error.response?.data?.error || error.message))
      } finally {
        this.isSubmitting = false
      }
    },
    
    async deleteEquipment(equipment) {
      if (!confirm(this.$t('equipment.confirmDelete').replace('{name}', equipment.name))) {
        return
      }
      
      try {
        await this.$api.delete(`/api/equipment/${equipment.id}`)
        alert(this.$t('equipment.deleteSuccess'))
        this.$emit('character-updated')
      } catch (error) {
        console.error('Fehler beim Löschen der Ausrüstung:', error)
        alert(this.$t('equipment.deleteError') + ': ' + (error.response?.data?.error || error.message))
      }
    },
    
    closeDialog() {
      this.showAddDialog = false
      this.selectedEquipment = null
      this.searchQuery = ''
      this.equipmentAmount = 1
      this.filteredMasterEquipment = []
      this.masterEquipment = []
    }
  }
}
</script>
