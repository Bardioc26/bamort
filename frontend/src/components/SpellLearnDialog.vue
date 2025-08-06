<template>
  <div v-if="show" class="modal-overlay" @click.self="closeDialog">
    <div class="modal-content modal-wide">
      <h3>{{ $t('spells.learn.title') }}</h3>
      
      <!-- Aktuelle Ressourcen -->
      <div class="current-resources">
        <div class="resource-display-card">
          <span class="resource-icon">⚡</span>
          <div class="resource-info">
            <div class="resource-label">Erfahrungspunkte</div>
            <div class="resource-amount">{{ character.erfahrungsschatz?.ep || 0 }} EP</div>
            <div class="resource-remaining">
              <small :class="{ 'text-warning': remainingEP < 50, 'text-danger': remainingEP <= 0 }">
                Verbleibend: {{ remainingEP }} EP
              </small>
            </div>
          </div>
        </div>
        <div class="resource-display-card">
          <span class="resource-icon">💰</span>
          <div class="resource-info">
            <div class="resource-label">Gold</div>
            <div class="resource-amount">{{ character.vermoegen?.goldstücke || 0 }} GS</div>
            <div class="resource-remaining">
              <small :class="{ 'text-warning': remainingGold < 20, 'text-danger': remainingGold <= 0 }">
                Nach Lernen: {{ remainingGold }} GS
              </small>
            </div>
          </div>
        </div>
      </div>

      <!-- Belohnungsart, Suche und Sortierung -->
      <div class="form-group form-row">
        <div class="form-col form-col-main">
          <label>Lernen als Belohnung:</label>
          <select v-model="selectedRewardType" :disabled="isLoadingRewardTypes">
            <option value="" disabled>
              {{ isLoadingRewardTypes ? 'Lade Belohnungsarten...' : 'Belohnungsart wählen' }}
            </option>
            <option 
              v-for="rewardType in availableRewardTypes" 
              :key="rewardType.value" 
              :value="rewardType.value"
            >
              {{ rewardType.label }}
            </option>
          </select>
        </div>
        <div class="form-col form-col-input">
          <label>{{ $t('spells.learn.search.label') }}</label>
          <input 
            v-model="searchTerm" 
            type="text" 
            :placeholder="$t('spells.learn.search.placeholder')"
          />
        </div>
        <div class="form-col form-col-input">
          <label>Sortierung:</label>
          <select v-model="sortBy">
            <option value="name">Name</option>
            <option value="epCost">EP Kosten</option>
            <option value="goldCost">Gold Kosten</option>
          </select>
        </div>
      </div>

      <!-- Schule Buttons -->
      <div class="form-group">
        <label>{{ $t('spells.learn.school.label') }}</label>
        <div class="school-buttons">
          <button 
            class="school-btn"
            :class="{ 'active': selectedSchool === '' }"
            @click="selectedSchool = ''"
          >
            {{ $t('spells.learn.school.all') }}
          </button>
          <button 
            v-for="school in availableSchools" 
            :key="school" 
            class="school-btn"
            :class="{ 'active': selectedSchool === school }"
            @click="selectedSchool = school"
          >
            {{ school }}
          </button>
        </div>
      </div>
      
      <!-- Verfügbare Zauber -->
      <div class="form-group">
        <label>{{ $t('spells.learn.available') }}</label>
        <div v-if="filteredSpells.length > 0" class="learning-levels">
          <div 
            v-for="spell in filteredSpells" 
            :key="spell.name"
            class="level-option"
            :class="{ 'selected': selectedSpell?.name === spell.name }"
            @click="selectSpell(spell)"
          >
            <div class="level-header">
              <span class="level-target">{{ spell.name }}</span>
              <span class="level-cost">
                <span v-if="spell.epCost">{{ spell.epCost }} EP</span>
                <span v-if="spell.epCost && spell.goldCost"> + </span>
                <span v-if="spell.goldCost">{{ spell.goldCost }} GS</span>
              </span>
            </div>
          </div>
        </div>
        <div v-else-if="!isLoading" class="no-spells">
          Keine Zauber verfügbar
        </div>
      </div>

      <!-- Ausgewählter Zauber Details -->
      <div v-if="selectedSpell" class="form-group">
        <div class="selection-summary">
          <strong>Ausgewählt:</strong> {{ selectedSpell.name }}
          <br>
          <span class="cost-summary">
            Lernkosten: 
            <span v-if="selectedSpell.epCost">{{ selectedSpell.epCost }} EP</span>
            <span v-if="selectedSpell.epCost && selectedSpell.goldCost"> + </span>
            <span v-if="selectedSpell.goldCost">{{ selectedSpell.goldCost }} GS</span>
          </span>
        </div>
      </div>

      <div v-if="isLoading" class="loading-message">
        {{ $t('common.loading') }}
      </div>

      <div class="modal-actions">
        <button 
          @click="learnSpell" 
          class="btn-confirm" 
          :disabled="!selectedSpell || isLoading"
        >
          {{ isLoading ? 'Wird gelernt...' : $t('spells.learn.action') }}
        </button>
        <button @click="closeDialog" class="btn-cancel" :disabled="isLoading">
          {{ $t('common.cancel') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import API from '@/utils/api'

export default {
  name: "SpellLearnDialog",
  props: {
    show: {
      type: Boolean,
      required: true
    },
    character: {
      type: Object,
      required: true
    }
  },
  emits: ['close', 'spell-learned'],
  data() {
    return {
      searchTerm: '',
      selectedSchool: '',
      selectedRewardType: '',
      sortBy: 'name',
      spellsBySchool: {},
      selectedSpell: null,
      isLoading: false,
      availableRewardTypes: [],
      isLoadingRewardTypes: false
    };
  },
  computed: {
    remainingEP() {
      const currentEP = this.character.erfahrungsschatz?.ep || 0;
      const spellEPCost = this.selectedSpell?.epCost || 0;
      return Math.max(0, currentEP - spellEPCost);
    },
    
    remainingGold() {
      const currentGold = this.character.vermoegen?.goldstücke || 0;
      const spellGoldCost = this.selectedSpell?.goldCost || 0;
      return Math.max(0, currentGold - spellGoldCost);
    },
    
    totalCosts() {
      if (!this.selectedSpell) return 0;
      return (this.selectedSpell.epCost || 0) + (this.selectedSpell.goldCost || 0);
    },
    
    filteredSpells() {
      let allSpells = [];
      
      // Sammle alle Zauber aus allen Schulen
      Object.keys(this.spellsBySchool).forEach(school => {
        if (!this.selectedSchool || school === this.selectedSchool) {
          allSpells = allSpells.concat(this.spellsBySchool[school]);
        }
      });
      
      // Filter nach Suchterm
      let filtered = allSpells.filter(spell => {
        const matchesSearch = !this.searchTerm || 
          spell.name.toLowerCase().includes(this.searchTerm.toLowerCase());
        
        return matchesSearch;
      });
      
      // Sortierung
      return filtered.sort((a, b) => {
        switch (this.sortBy) {
          case 'epCost':
            return (a.epCost || 0) - (b.epCost || 0);
          case 'goldCost':
            return (a.goldCost || 0) - (b.goldCost || 0);
          case 'name':
          default:
            return a.name.localeCompare(b.name);
        }
      });
    },
    
    availableSchools() {
      return Object.keys(this.spellsBySchool).sort();
    }
  },
  watch: {
    show(newVal) {
      if (newVal) {
        this.resetForm();
        this.loadRewardTypes();
      }
    },
    selectedRewardType() {
      if (this.selectedRewardType) {
        this.loadAvailableSpells();
      }
    }
  },
  created() {
    this.$api = API;
  },
  methods: {
    closeDialog() {
      this.$emit('close');
    },
    
    resetForm() {
      this.searchTerm = '';
      this.selectedSchool = '';
      this.selectedRewardType = '';
      this.sortBy = 'name';
      this.spellsBySchool = {};
      this.selectedSpell = null;
      this.availableRewardTypes = [];
    },
    
    async loadRewardTypes() {
      const token = localStorage.getItem('token');
      if (!token) {
        console.error('No authentication token available');
        this.availableRewardTypes = [{
          value: 'error',
          label: 'Anmeldung erforderlich'
        }];
        return;
      }
      
      this.isLoadingRewardTypes = true;
      try {
        const response = await this.$api.get(`/api/characters/${this.character.id}/reward-types`, {
          params: {
            learning_type: 'spell',
            skill_name: 'spell',
            current_level: 0,
            skill_type: 'spell'
          }
        });
        
        this.availableRewardTypes = response.data.reward_types || [];
        
        // Setze Default-Belohnungsart wenn verfügbar
        if (this.availableRewardTypes.length > 0 && !this.selectedRewardType) {
          const defaultReward = this.availableRewardTypes.find(r => r.value === 'default');
          this.selectedRewardType = defaultReward ? defaultReward.value : this.availableRewardTypes[0].value;
        }
      } catch (error) {
        console.error('Fehler beim Laden der Belohnungsarten:', error);
        this.availableRewardTypes = [{
          value: 'default',
          label: 'Standard'
        }];
        this.selectedRewardType = 'default';
      } finally {
        this.isLoadingRewardTypes = false;
      }
    },
    
    performInitialSearch() {
      // Beim Öffnen alle Zauber laden
      this.loadAvailableSpells();
    },
    
    async loadAvailableSpells() {
      if (!this.selectedRewardType) return;
      
      try {
        this.isLoading = true;
        
        // Erstelle LernCostRequest wie vom Backend erwartet
        const request = {
          char_id: this.character.id,
          type: 'spell',
          action: 'learn',
          use_pp: 0,
          use_gold: 0,
          reward: this.selectedRewardType
        };
        
        const response = await this.$api.post('/api/characters/available-spells-new', request);
        this.spellsBySchool = response.data.spells_by_school || {};
      } catch (error) {
        console.error('Fehler beim Laden der verfügbaren Zauber:', error);
        this.spellsBySchool = {};
      } finally {
        this.isLoading = false;
      }
    },
    
    selectSpell(spell) {
      this.selectedSpell = this.selectedSpell?.name === spell.name ? null : spell;
    },
    
    async learnSpell() {
      if (!this.selectedSpell || this.isLoading) return;
      
      try {
        this.isLoading = true;
        
        const response = await this.$api.post(`/characters/${this.character.id}/learn-spell-new`, {
          char_id: this.character.id,
          name: this.selectedSpell.name,
          type: 'spell',
          action: 'learn',
          use_pp: 0,
          use_gold: 0,
          reward: this.selectedRewardType
        });
        
        this.$emit('spell-learned', {
          spell: this.selectedSpell,
          response: response.data
        });
        
      } catch (error) {
        console.error('Fehler beim Lernen des Zaubers:', error);
        alert('Fehler beim Lernen des Zaubers: ' + (error.response?.data?.message || error.message));
      } finally {
        this.isLoading = false;
      }
    }
  }
};
</script>

<style scoped>
/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: flex-start;
  align-items: flex-start;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  padding: 24px;
  width: 100vw;
  height: 100vh;
  max-width: 100vw;
  max-height: 100vh;
  overflow-y: auto;
  animation: modalSlideIn 0.3s ease;
  box-sizing: border-box;
}

.modal-wide {
  max-width: 100vw;
}

/* Ressourcen-Anzeige im Dialog */
.current-resources {
  display: flex;
  gap: 15px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.resource-display-card {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  background: #f8f9fa;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  flex: 1;
  min-width: 160px;
}

.resource-display-card .resource-icon {
  font-size: 20px;
}

.resource-info {
  flex: 1;
}

.resource-label {
  font-size: 12px;
  color: #6c757d;
  font-weight: 500;
}

.resource-amount {
  font-size: 16px;
  font-weight: bold;
  color: #495057;
}

.resource-remaining {
  margin-top: 4px;
}

.resource-remaining small {
  color: #6c757d;
  font-weight: normal;
}

.text-warning {
  color: #f0ad4e !important;
}

.text-danger {
  color: #d9534f !important;
}

/* Filter und Sortierung */
.school-buttons {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-top: 5px;
}

.school-btn {
  padding: 6px 12px;
  border: 1px solid #dee2e6;
  border-radius: 4px;
  background: white;
  color: #495057;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s ease;
}

.school-btn:hover {
  background: #f8f9fa;
  border-color: #007bff;
}

.school-btn.active {
  background: #007bff;
  color: white;
  border-color: #007bff;
}

/* Zauber-Auswahl */
.selection-summary {
  background: #e7f3ff;
  padding: 12px;
  border-radius: 6px;
  margin-bottom: 10px;
  border-left: 4px solid #007bff;
}

.cost-summary {
  color: #28a745;
  font-weight: bold;
}

.learning-levels {
  border: 1px solid #dee2e6;
  border-radius: 6px;
  max-height: 300px;
  overflow-y: auto;
}

.level-option {
  padding: 12px 16px;
  border-bottom: 1px solid #f1f1f1;
  cursor: pointer;
  transition: all 0.2s ease;
}

.level-option:last-child {
  border-bottom: none;
}

.level-option:hover:not(.disabled) {
  background: #f8f9fa;
}

.level-option.selected {
  background: #e7f3ff;
  border-left: 4px solid #007bff;
}

.level-option.disabled {
  background: #f8f9fa;
  color: #6c757d;
  cursor: not-allowed;
  opacity: 0.6;
}

.level-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 500;
}

.level-target {
  color: #495057;
}

.level-cost {
  color: #28a745;
  font-weight: bold;
}

.level-option.disabled .level-cost {
  color: #dc3545;
}

.no-spells {
  text-align: center;
  padding: 20px;
  color: #6c757d;
  font-style: italic;
}

.loading-message {
  text-align: center;
  padding: 20px;
  color: #6c757d;
}

@keyframes modalSlideIn {
  from {
    opacity: 0;
    transform: scale(0.9) translateY(-20px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.modal-content h3 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #333;
  border-bottom: 2px solid #1da766;
  padding-bottom: 10px;
}

.form-group {
  margin-bottom: 15px;
}

.form-row {
  display: flex;
  gap: 15px;
  align-items: flex-start;
}

.form-col {
  flex: 1;
  min-width: 0;
}

.form-col-main {
  flex: 2;
  min-width: 200px;
}

.form-col-input {
  flex: 1;
  min-width: 140px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
  color: #555;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}

.form-group textarea {
  height: 80px;
  resize: vertical;
}

.help-text {
  display: block;
  margin-top: 5px;
  font-size: 12px;
  color: #6c757d;
  font-style: italic;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
  padding-top: 15px;
  border-top: 1px solid #eee;
}

.btn-confirm {
  padding: 8px 20px;
  background: #1da766;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  transition: background 0.2s ease;
}

.btn-confirm:hover:not(:disabled) {
  background: #16a085;
}

.btn-confirm:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.btn-cancel {
  padding: 8px 20px;
  background: #6c757d;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background 0.2s ease;
}

.btn-cancel:hover:not(:disabled) {
  background: #5a6268;
}
</style>
