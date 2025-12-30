<template>
  <div class="cd-view">
    <div class="header-section">
      <h2>{{ $t('WeaponView') }}</h2>
      <button @click="openAddWeaponDialog" class="btn-add-weapon">
        {{ $t('weapon.add') }}
      </button>
    </div>

    <div class="cd-list">
      <table class="cd-table">
      <thead>
        <tr>
          <th>{{ $t('weapon.name') }}</th>
          <th>{{ $t('weapon.description') }}</th>
          <th>{{ $t('weapon.weight') }}</th>
          <th>{{ $t('weapon.value') }}</th>
          <th>{{ $t('weapon.amount') }}</th>
          <th>{{ $t('weapon.contained_in') }}</th>
          <th>{{ $t('weapon.bonus') }}</th>
          <th>{{ $t('weapon.actions') }}</th>
        </tr>
      </thead>
      <tbody>
      <template v-if="character.waffen && character.waffen.length > 0">
        <tr v-for="weapon in character.waffen" :key="weapon.id">
          <td>{{ weapon.name || '-' }}<span v-if="weapon.ist_magisch" class="magic-indicator">*</span></td>
          <td>{{ weapon.beschreibung || '-' }}</td>
          <td>{{ weapon.gewicht || '-' }}</td>
          <td>{{ weapon.wert || '-' }}</td>
          <td>{{ weapon.anzahl || '-' }}</td>
          <td>{{ weapon.beinhaltet_in || '-' }}</td>
          <td>{{ weapon.anb || '-' }}/{{ weapon.abwb || '-' }}</td>
          <td class="action-cell">
            <button @click="editWeapon(weapon)" class="btn-edit" title="Bearbeiten">
              ✏️
            </button>
            <button @click="deleteWeapon(weapon)" class="btn-delete" title="Löschen">
              🗑️
            </button>
          </td>
        </tr>
      </template>
      <template v-else>
        <tr>
          <td colspan="8" class="empty-state">{{ $t('weapon.noWeapons') }}</td>
        </tr>
      </template>
      </tbody>
      </table>
    </div>

    <!-- Dialog für Waffe hinzufügen -->
    <div v-if="showAddDialog" class="modal-overlay" @click.self="closeDialog">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ $t('weapon.addWeapon') }}</h3>
          <button @click="closeDialog" class="close-button">&times;</button>
        </div>
        
        <div class="modal-body">
          <div class="form-group">
            <label>{{ $t('weapon.search') }}:</label>
            <input 
              v-model="searchQuery" 
              type="text" 
              :placeholder="$t('weapon.searchPlaceholder')"
              @input="filterWeapons"
            />
          </div>

          <div v-if="isLoading" class="loading">{{ $t('weapon.loading') }}</div>

          <div v-else-if="filteredMasterWeapons.length > 0" class="weapon-list">
            <div 
              v-for="weapon in filteredMasterWeapons" 
              :key="weapon.id"
              class="weapon-item"
              :class="{ selected: selectedWeapon && selectedWeapon.id === weapon.id }"
              @click="selectWeapon(weapon)"
            >
              <div class="weapon-name">{{ weapon.name }}</div>
              <div class="weapon-details">
                {{ weapon.damage || '-' }} | 
                {{ weapon.weight }}kg | 
                {{ weapon.value || '-' }}
              </div>
            </div>
          </div>

          <div v-else class="no-results">
            {{ searchQuery ? $t('weapon.noResults') : $t('weapon.noMasterData') }}
          </div>

          <div v-if="selectedWeapon" class="selected-weapon-details">
            <h4>{{ $t('weapon.selectedWeapon') }}</h4>
            <p><strong>{{ $t('weapon.name') }}:</strong> {{ selectedWeapon.name }}</p>
            <p><strong>{{ $t('weapon.damage') }}:</strong> {{ selectedWeapon.damage || '-' }}</p>
            <p><strong>{{ $t('weapon.weight') }}:</strong> {{ selectedWeapon.weight }}kg</p>
            <p><strong>{{ $t('weapon.value') }}:</strong> {{ selectedWeapon.value || '-' }}</p>
            
            <div class="form-group">
              <label>{{ $t('weapon.amount') }}:</label>
              <input v-model.number="weaponAmount" type="number" min="1" value="1" />
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button @click="closeDialog" class="btn-cancel">{{ $t('weapon.cancel') }}</button>
          <button 
            @click="addWeapon" 
            class="btn-confirm" 
            :disabled="!selectedWeapon || isSubmitting"
          >
            {{ isSubmitting ? $t('weapon.adding') : $t('weapon.add') }}
          </button>
        </div>
      </div>
    </div>

    <!-- Dialog für Waffe bearbeiten -->
    <div v-if="showEditDialog" class="modal-overlay" @click.self="closeEditDialog">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ $t('weapon.editWeapon') }}</h3>
          <button @click="closeEditDialog" class="close-button">&times;</button>
        </div>
        
        <div class="modal-body">
          <div class="form-group">
            <label>{{ $t('weapon.name') }}:</label>
            <input v-model="editingWeapon.name" type="text" disabled />
          </div>

          <div class="form-group">
            <label>{{ $t('weapon.description') }}:</label>
            <textarea v-model="editingWeapon.beschreibung" rows="3"></textarea>
          </div>

          <div class="form-group">
            <label>
              <input type="checkbox" v-model="editingWeapon.ist_magisch" />
              {{ $t('weapon.magical') }}
            </label>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>{{ $t('weapon.amount') }}:</label>
              <input v-model.number="editingWeapon.anzahl" type="number" min="1" />
            </div>

            <div class="form-group">
              <label>{{ $t('weapon.value') }}:</label>
              <input v-model.number="editingWeapon.wert" type="number" min="0" step="0.01" />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>{{ $t('weapon.attackBonus') }} (Anb):</label>
              <input v-model.number="editingWeapon.anb" type="number" />
            </div>

            <div class="form-group">
              <label>{{ $t('weapon.defenseBonus') }} (Abwb):</label>
              <input v-model.number="editingWeapon.abwb" type="number" />
            </div>
          </div>

          <div class="form-group">
            <label>{{ $t('weapon.damageBonus') }} (Schb):</label>
            <input v-model.number="editingWeapon.schb" type="number" />
          </div>
        </div>

        <div class="modal-footer">
          <button @click="closeEditDialog" class="btn-cancel">{{ $t('weapon.cancel') }}</button>
          <button 
            @click="saveWeapon" 
            class="btn-confirm" 
            :disabled="isSubmitting"
          >
            {{ isSubmitting ? $t('weapon.saving') : $t('weapon.save') }}
          </button>
        </div>
      </div>
    </div>
  </div>

</template>

<style>
/* All common styles moved to main.css */

.cd-view {
  padding: 1rem;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.btn-add-weapon {
  padding: 8px 16px;
  background: #1da766;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: bold;
  transition: background 0.2s ease;
}

.btn-add-weapon:hover {
  background: #16a085;
}

.cd-table {
  width: 100%;
  border-collapse: collapse;
}

.cd-table th,
.cd-table td {
  padding: 8px;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

.cd-table th {
  background-color: #1da766;
  color: white;
  font-weight: bold;
}

.empty-state {
  text-align: center;
  color: #999;
  font-style: italic;
  padding: 2rem !important;
}

.magic-indicator {
  color: #9c27b0;
  font-weight: bold;
  margin-left: 4px;
}

.action-cell {
  text-align: center;
}

.btn-edit {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1.2rem;
  padding: 4px 8px;
  margin-right: 8px;
  transition: transform 0.2s ease;
}

.btn-edit:hover {
  transform: scale(1.2);
}

.btn-delete {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1.2rem;
  padding: 4px 8px;
  transition: transform 0.2s ease;
}

.btn-delete:hover {
  transform: scale(1.2);
}

.form-row {
  display: flex;
  gap: 15px;
  margin-bottom: 15px;
}

.form-row .form-group {
  flex: 1;
  margin-bottom: 0;
}

.form-group textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
  resize: vertical;
  font-family: inherit;
}

.loading {
  text-align: center;
  padding: 2rem;
  color: #666;
}

.weapon-list {
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid #ddd;
  border-radius: 4px;
  margin-bottom: 15px;
}

.weapon-item {
  padding: 12px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
  transition: background 0.2s ease;
}

.weapon-item:hover {
  background: #f5f5f5;
}

.weapon-item.selected {
  background: #e8f5e9;
  border-left: 3px solid #1da766;
}

.weapon-name {
  font-weight: bold;
  color: #333;
  margin-bottom: 4px;
}

.weapon-details {
  font-size: 0.9em;
  color: #666;
}

.no-results {
  text-align: center;
  padding: 2rem;
  color: #999;
  font-style: italic;
}

.selected-weapon-details {
  background: #f8f9fa;
  padding: 15px;
  border-radius: 4px;
  margin-top: 15px;
}

.selected-weapon-details h4 {
  margin-top: 0;
  color: #1da766;
}

.selected-weapon-details p {
  margin: 8px 0;
}
</style>

<script>
import API from '@/utils/api'

export default {
name: "WeaponView",
  props: {
    character: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      showAddDialog: false,
      showEditDialog: false,
      isLoading: false,
      isSubmitting: false,
      masterWeapons: [],
      filteredMasterWeapons: [],
      selectedWeapon: null,
      editingWeapon: null,
      searchQuery: '',
      weaponAmount: 1
    }
  },
  created() {
    this.$api = API
  },
  methods: {
    async openAddWeaponDialog() {
      this.showAddDialog = true
      this.isLoading = true
      
      try {
        const response = await this.$api.get('/api/maintenance/weapons')
        this.masterWeapons = response.data || []
        this.filteredMasterWeapons = this.masterWeapons
      } catch (error) {
        console.error('Fehler beim Laden der Waffen-Stammdaten:', error)
        alert(this.$t('weapon.loadError') + ': ' + (error.response?.data?.error || error.message))
      } finally {
        this.isLoading = false
      }
    },
    
    filterWeapons() {
      const query = this.searchQuery.toLowerCase().trim()
      if (!query) {
        this.filteredMasterWeapons = this.masterWeapons
      } else {
        this.filteredMasterWeapons = this.masterWeapons.filter(weapon => 
          weapon.name.toLowerCase().includes(query) ||
          (weapon.beschreibung && weapon.beschreibung.toLowerCase().includes(query))
        )
      }
    },
    
    selectWeapon(weapon) {
      this.selectedWeapon = weapon
    },
    
    async addWeapon() {
      if (!this.selectedWeapon) {
        alert(this.$t('weapon.pleaseSelect'))
        return
      }
      
      this.isSubmitting = true
      
      try {
        const weaponData = {
          character_id: this.character.id,
          name: this.selectedWeapon.name,
          beschreibung: this.selectedWeapon.beschreibung || '',
          gewicht: this.selectedWeapon.weight || 0,
          wert: this.selectedWeapon.value || 0,
          anzahl: this.weaponAmount,
          ist_magisch: false,
          anb: this.selectedWeapon.attack_bonus || 0,
          abwb: this.selectedWeapon.defense_bonus || 0,
          schb: this.selectedWeapon.damage_bonus || 0,
          beinhaltet_in: '',
          contained_in: 0,
          container_type: ''
        }
        
        await this.$api.post('/api/weapons', weaponData)
        
        alert(this.$t('weapon.addSuccess'))
        this.closeDialog()
        this.$emit('character-updated')
        
      } catch (error) {
        console.error('Fehler beim Hinzufügen der Waffe:', error)
        alert(this.$t('weapon.addError') + ': ' + (error.response?.data?.error || error.message))
      } finally {
        this.isSubmitting = false
      }
    },
    
    editWeapon(weapon) {
      this.editingWeapon = {
        id: weapon.id,
        name: weapon.name,
        beschreibung: weapon.beschreibung || '',
        ist_magisch: weapon.ist_magisch || false,
        anzahl: weapon.anzahl || 1,
        wert: weapon.wert || 0,
        anb: weapon.anb || 0,
        abwb: weapon.abwb || 0,
        schb: weapon.schb || 0
      }
      this.showEditDialog = true
    },
    
    async saveWeapon() {
      if (!this.editingWeapon) {
        return
      }
      
      this.isSubmitting = true
      
      try {
        const weaponData = {
          beschreibung: this.editingWeapon.beschreibung,
          ist_magisch: this.editingWeapon.ist_magisch,
          anzahl: this.editingWeapon.anzahl,
          wert: this.editingWeapon.wert,
          anb: this.editingWeapon.anb,
          abwb: this.editingWeapon.abwb,
          schb: this.editingWeapon.schb
        }
        
        await this.$api.put(`/api/weapons/${this.editingWeapon.id}`, weaponData)
        
        alert(this.$t('weapon.editSuccess'))
        this.closeEditDialog()
        this.$emit('character-updated')
        
      } catch (error) {
        console.error('Fehler beim Speichern der Waffe:', error)
        alert(this.$t('weapon.editError') + ': ' + (error.response?.data?.error || error.message))
      } finally {
        this.isSubmitting = false
      }
    },
    
    async deleteWeapon(weapon) {
      if (!confirm(this.$t('weapon.confirmDelete').replace('{name}', weapon.name))) {
        return
      }
      
      try {
        await this.$api.delete(`/api/weapons/${weapon.id}`)
        alert(this.$t('weapon.deleteSuccess'))
        this.$emit('character-updated')
      } catch (error) {
        console.error('Fehler beim Löschen der Waffe:', error)
        alert(this.$t('weapon.deleteError') + ': ' + (error.response?.data?.error || error.message))
      }
    },
    
    closeEditDialog() {
      this.showEditDialog = false
      this.editingWeapon = null
    },
    
    closeDialog() {
      this.showAddDialog = false
      this.selectedWeapon = null
      this.searchQuery = ''
      this.weaponAmount = 1
      this.filteredMasterWeapons = []
      this.masterWeapons = []
    }
  }
}
</script>
