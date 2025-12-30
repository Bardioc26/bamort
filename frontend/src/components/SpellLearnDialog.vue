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
      
      <!-- Verfügbare Zauber und Zu lernende Zauber -->
      <div class="spells-container">
        <!-- Linke Spalte: Verfügbare Zauber -->
        <div class="available-spells-section">
          <div class="form-group">
            <label>{{ $t('spells.learn.available') }}</label>
            <div v-if="filteredSpells.length > 0" class="learning-levels">
              <div 
                v-for="spell in filteredSpells" 
                :key="spell.name"
                class="level-option"
                :class="{ 
                  'selected': selectedSpell?.name === spell.name,
                  'already-selected': isSpellInLearningList(spell.name)
                }"
                @click="selectSpell(spell)"
              >
                <div class="level-header">
                  <span class="level-target">{{ spell.name }}</span>
                  <div class="spell-actions-inline">
                    <span class="level-cost">
                      <span v-if="spell.epCost">{{ spell.epCost }} EP</span>
                      <span v-if="spell.epCost && spell.goldCost"> + </span>
                      <span v-if="spell.goldCost">{{ spell.goldCost }} GS</span>
                    </span>
                    <button 
                      @click.stop="addSpellToLearningListDirect(spell)"
                      class="btn-add-inline"
                      :disabled="isSpellInLearningList(spell.name) || !canAffordSpell(spell)"
                      :title="isSpellInLearningList(spell.name) ? 'Bereits in Liste' : 'Zur Liste hinzufügen'"
                    >
                      {{ isSpellInLearningList(spell.name) ? '✓' : '+' }}
                    </button>
                  </div>
                </div>
              </div>
            </div>
            <div v-else-if="!isLoading" class="no-spells">
              Keine Zauber verfügbar
            </div>
          </div>
        </div>

        <!-- Rechte Spalte: Zu lernende Zauber -->
        <div class="learning-list-section">
          <div class="form-group">
            <label>Zu lernende Zauber</label>
            <div v-if="spellsToLearn.length > 0" class="learning-levels">
              <div 
                v-for="(spell, index) in spellsToLearn" 
                :key="spell.name"
                class="level-option learning-item"
              >
                <div class="level-header">
                  <span class="level-target">{{ spell.name }}</span>
                  <span class="level-cost">
                    <span v-if="spell.epCost">{{ spell.epCost }} EP</span>
                    <span v-if="spell.epCost && spell.goldCost"> + </span>
                    <span v-if="spell.goldCost">{{ spell.goldCost }} GS</span>
                  </span>
                  <button 
                    @click="removeSpellFromLearningList(index)"
                    class="remove-btn"
                    title="Zauber aus Liste entfernen"
                  >
                    ×
                  </button>
                </div>
              </div>
            </div>
            <div v-else class="no-spells">
              Keine Zauber ausgewählt
            </div>
            
            <!-- Gesamtkosten -->
            <div v-if="spellsToLearn.length > 0" class="total-costs">
              <div class="cost-summary">
                <strong>Gesamtkosten:</strong>
                <span v-if="totalLearningCosts.ep > 0">{{ totalLearningCosts.ep }} EP</span>
                <span v-if="totalLearningCosts.ep > 0 && totalLearningCosts.gold > 0"> + </span>
                <span v-if="totalLearningCosts.gold > 0">{{ totalLearningCosts.gold }} GS</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Ausgewählter Zauber Aktionen und Details -->
      <div v-if="selectedSpell" class="form-group">
        <div class="spell-details-section">
          <!--
          <div class="selection-summary">
            <div class="spell-actions">
              <strong>Ausgewählt:</strong> {{ selectedSpell.name }}
              <span class="cost-info">
                - Kosten: 
                <span v-if="selectedSpell.epCost">{{ selectedSpell.epCost }} EP</span>
                <span v-if="selectedSpell.epCost && selectedSpell.goldCost"> + </span>
                <span v-if="selectedSpell.goldCost">{{ selectedSpell.goldCost }} GS</span>
              </span>
              <button 
                @click="addSpellToLearningList"
                class="btn-add-spell"
                :disabled="isSpellInLearningList(selectedSpell.name) || !canAffordSpell(selectedSpell)"
              >
                {{ isSpellInLearningList(selectedSpell.name) ? 'Bereits in Liste' : 'Zur Liste hinzufügen' }}
              </button>
            </div>
          </div>
          -->
          <!-- Detaillierte Zauber-Informationen -->
          <div v-if="isLoadingSpellDetails" class="loading-spell-details">
            <span>Lade Zauber-Details...</span>
          </div>
          
          <div v-else-if="selectedSpellDetails" class="spell-details-grid">
            <div class="spell-detail-card">
              <h4>Grunddaten</h4>
              <div class="detail-row">
                <span class="detail-label">Stufe:</span>
                <span class="detail-value">{{ selectedSpellDetails.level }}</span>
              </div>
              <div class="detail-row" v-if="selectedSpellDetails.bonus">
                <span class="detail-label">Bonus:</span>
                <span class="detail-value">+{{ selectedSpellDetails.bonus }}</span>
              </div>
              <div class="detail-row" v-if="selectedSpellDetails.category">
                <span class="detail-label">Schule:</span>
                <span class="detail-value">{{ selectedSpellDetails.category }}</span>
              </div>
              <div class="detail-row" v-if="selectedSpellDetails.ursprung">
                <span class="detail-label">Ursprung:</span>
                <span class="detail-value">{{ selectedSpellDetails.ursprung }}</span>
              </div>
            </div>

            <div class="spell-detail-card" v-if="selectedSpellDetails.ap || selectedSpellDetails.art || selectedSpellDetails.zauberdauer">
              <h4>Ausführung</h4>
              <div class="detail-row" v-if="selectedSpellDetails.ap">
                <span class="detail-label">AP:</span>
                <span class="detail-value">{{ selectedSpellDetails.ap }}</span>
              </div>
              <div class="detail-row" v-if="selectedSpellDetails.art">
                <span class="detail-label">Art:</span>
                <span class="detail-value">{{ selectedSpellDetails.art }}</span>
              </div>
              <div class="detail-row" v-if="selectedSpellDetails.zauberdauer">
                <span class="detail-label">Zauberdauer:</span>
                <span class="detail-value">{{ selectedSpellDetails.zauberdauer }}</span>
              </div>
            </div>

            <div class="spell-detail-card" v-if="selectedSpellDetails.reichweite || selectedSpellDetails.wirkungsziel || selectedSpellDetails.wirkungsbereich">
              <h4>Reichweite & Ziel</h4>
              <div class="detail-row" v-if="selectedSpellDetails.reichweite">
                <span class="detail-label">Reichweite:</span>
                <span class="detail-value">{{ selectedSpellDetails.reichweite }}</span>
              </div>
              <div class="detail-row" v-if="selectedSpellDetails.wirkungsziel">
                <span class="detail-label">Wirkungsziel:</span>
                <span class="detail-value">{{ selectedSpellDetails.wirkungsziel }}</span>
              </div>
              <div class="detail-row" v-if="selectedSpellDetails.wirkungsbereich">
                <span class="detail-label">Wirkungsbereich:</span>
                <span class="detail-value">{{ selectedSpellDetails.wirkungsbereich }}</span>
              </div>
            </div>

            <div class="spell-detail-card" v-if="selectedSpellDetails.wirkungsdauer">
              <h4>Wirkung</h4>
              <div class="detail-row">
                <span class="detail-label">Wirkungsdauer:</span>
                <span class="detail-value">{{ selectedSpellDetails.wirkungsdauer }}</span>
              </div>
            </div>
          </div>

          <!-- Beschreibung -->
          <div class="spell-description" v-if="selectedSpellDetails && selectedSpellDetails.beschreibung">
            <h4>Beschreibung</h4>
            <p>{{ selectedSpellDetails.beschreibung }}</p>
          </div>
        </div>
      </div>

      <div v-if="isLoading" class="loading-message">
        {{ $t('common.loading') }}
      </div>

      <div class="modal-actions">
        <button 
          @click="learnAllSpells" 
          class="btn-confirm" 
          :disabled="spellsToLearn.length === 0 || isLoading || !canAffordAllSpells"
        >
          {{ isLoading ? 'Wird gelernt...' : `${spellsToLearn.length} Zauber lernen` }}
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
      selectedSpellDetails: null,
      spellsToLearn: [],
      isLoading: false,
      isLoadingSpellDetails: false,
      availableRewardTypes: [],
      isLoadingRewardTypes: false
    };
  },
  computed: {
    remainingEP() {
      const currentEP = this.character.erfahrungsschatz?.ep || 0;
      const usedEP = this.totalLearningCosts.ep;
      return Math.max(0, currentEP - usedEP);
    },
    
    remainingGold() {
      const currentGold = this.character.vermoegen?.goldstücke || 0;
      const usedGold = this.totalLearningCosts.gold;
      return Math.max(0, currentGold - usedGold);
    },
    
    totalLearningCosts() {
      return this.spellsToLearn.reduce((total, spell) => {
        total.ep += spell.epCost || 0;
        total.gold += spell.goldCost || 0;
        return total;
      }, { ep: 0, gold: 0 });
    },
    
    canAffordAllSpells() {
      const currentEP = this.character.erfahrungsschatz?.ep || 0;
      const currentGold = this.character.vermoegen?.goldstücke || 0;
      const costs = this.totalLearningCosts;
      
      return currentEP >= costs.ep && currentGold >= costs.gold;
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
      this.selectedSpellDetails = null;
      this.spellsToLearn = [];
      this.availableRewardTypes = [];
    },
    
    isSpellInLearningList(spellName) {
      return this.spellsToLearn.some(spell => spell.name === spellName);
    },
    
    canAffordSpell(spell) {
      const currentEP = this.character.erfahrungsschatz?.ep || 0;
      const currentGold = this.character.vermoegen?.goldstücke || 0;
      const totalCostsWithSpell = {
        ep: this.totalLearningCosts.ep + (spell.epCost || 0),
        gold: this.totalLearningCosts.gold + (spell.goldCost || 0)
      };
      
      return currentEP >= totalCostsWithSpell.ep && currentGold >= totalCostsWithSpell.gold;
    },
    
    addSpellToLearningList() {
      if (!this.selectedSpell || this.isSpellInLearningList(this.selectedSpell.name)) return;
      if (!this.canAffordSpell(this.selectedSpell)) return;
      
      this.spellsToLearn.push({ ...this.selectedSpell });
      this.selectedSpell = null;
    },
    
    addSpellToLearningListDirect(spell) {
      if (this.isSpellInLearningList(spell.name)) return;
      if (!this.canAffordSpell(spell)) return;
      
      this.spellsToLearn.push({ ...spell });
    },
    
    removeSpellFromLearningList(index) {
      this.spellsToLearn.splice(index, 1);
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
    
    async loadSpellDetails(spellName) {
      if (!spellName) return;
      
      try {
        this.isLoadingSpellDetails = true;
        const response = await this.$api.get('/api/characters/spell-details', {
          params: { name: spellName }
        });
        
        this.selectedSpellDetails = response.data.spell;
      } catch (error) {
        console.error('Fehler beim Laden der Zauber-Details:', error);
        this.selectedSpellDetails = null;
      } finally {
        this.isLoadingSpellDetails = false;
      }
    },
    
    selectSpell(spell) {
      const wasSelected = this.selectedSpell?.name === spell.name;
      this.selectedSpell = wasSelected ? null : spell;
      
      if (this.selectedSpell) {
        // Lade Details nur wenn der Zauber ausgewählt wurde
        this.loadSpellDetails(this.selectedSpell.name);
      } else {
        // Leerere Details wenn kein Zauber ausgewählt ist
        this.selectedSpellDetails = null;
      }
    },
    
    async learnAllSpells() {
      if (this.spellsToLearn.length === 0 || this.isLoading || !this.canAffordAllSpells) return;
      
      try {
        this.isLoading = true;
        const responses = [];
        
        // Lerne jeden Zauber einzeln
        for (const spell of this.spellsToLearn) {
          const response = await this.$api.post(`/api/characters/${this.character.id}/learn-spell-new`, {
            char_id: this.character.id,
            name: spell.name,
            type: 'spell',
            action: 'learn',
            use_pp: 0,
            use_gold: 0,
            reward: this.selectedRewardType
          });
          
          responses.push({
            spell: spell,
            response: response.data
          });
        }
        
        // Erfolgsmeldung mit Details
        const learnedCount = responses.length;
        this.$emit('spell-learned', {
          spells: this.spellsToLearn,
          responses: responses,
          count: learnedCount
        });
        
        // Dialog schließen nach erfolgreichem Lernen
        this.closeDialog();
        
      } catch (error) {
        console.error('Fehler beim Lernen der Zauber:', error);
        alert('Fehler beim Lernen der Zauber: ' + (error.response?.data?.message || error.message));
      } finally {
        this.isLoading = false;
      }
    }
  }
};
</script>
<style>
/* All common styles moved to main.css */

.modal-content {
  width: 100vw;
  height: 100vh;
  max-width: 100vw;
  max-height: 100vh;
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

/* Zweispaltiges Layout */
.spells-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 15px;
}

.available-spells-section,
.learning-list-section {
  min-height: 300px;
}

.learning-item {
  background: #f0f8ff !important;
  border-left: 3px solid #007bff !important;
}

.learning-item .level-header {
  position: relative;
}

.remove-btn {
  position: absolute;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
  background: #dc3545;
  color: white;
  border: none;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  cursor: pointer;
  font-size: 16px;
  line-height: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.remove-btn:hover {
  background: #c82333;
}

.already-selected {
  opacity: 0.5;
  pointer-events: none;
}

.spell-actions-inline {
  display: flex;
  align-items: center;
  gap: 10px;
}

.btn-add-inline {
  background: #28a745;
  color: white;
  border: none;
  border-radius: 50%;
  width: 28px;
  height: 28px;
  cursor: pointer;
  font-size: 14px;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.btn-add-inline:hover:not(:disabled) {
  background: #218838;
  transform: scale(1.1);
}

.btn-add-inline:disabled {
  background: #6c757d;
  cursor: not-allowed;
  transform: none;
}

.already-selected .btn-add-inline {
  background: #17a2b8;
}

.total-costs {
  margin-top: 10px;
  padding: 10px;
  background: #e7f3ff;
  border-radius: 4px;
  border-left: 4px solid #007bff;
}

.spell-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.cost-info {
  color: #28a745;
  font-weight: bold;
}

.btn-add-spell {
  padding: 4px 12px;
  background: #28a745;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  transition: background 0.2s ease;
}

.btn-add-spell:hover:not(:disabled) {
  background: #218838;
}

.btn-add-spell:disabled {
  background: #6c757d;
  cursor: not-allowed;
}

@media (max-width: 1024px) {
  .spells-container {
    grid-template-columns: 1fr;
  }
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

/* Zauber-Auswahl und Details */
.spell-details-section {
  background: #e7f3ff;
  padding: 16px;
  border-radius: 6px;
  margin-bottom: 10px;
  border-left: 4px solid #007bff;
}

.loading-spell-details {
  text-align: center;
  padding: 20px;
  color: #6c757d;
  font-style: italic;
}

.spell-details-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 15px;
  margin-top: 15px;
}

.spell-detail-card {
  background: white;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  padding: 12px;
}

.spell-detail-card h4 {
  margin: 0 0 10px 0;
  color: #495057;
  font-size: 14px;
  font-weight: bold;
  border-bottom: 1px solid #e9ecef;
  padding-bottom: 5px;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
  font-size: 13px;
}

.detail-row:last-child {
  margin-bottom: 0;
}

.detail-label {
  color: #6c757d;
  font-weight: 500;
  flex: 0 0 auto;
  margin-right: 10px;
}

.detail-value {
  color: #495057;
  text-align: right;
  flex: 1 1 auto;
}

.spell-description {
  background: white;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  padding: 12px;
  margin-top: 15px;
}

.spell-description h4 {
  margin: 0 0 8px 0;
  color: #495057;
  font-size: 14px;
  font-weight: bold;
}

.spell-description p {
  margin: 0;
  color: #495057;
  font-size: 13px;
  line-height: 1.4;
}

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

.modal-content h3 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #333;
  border-bottom: 2px solid #1da766;
  padding-bottom: 10px;
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
}
</style>
