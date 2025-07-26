<template>
  <div v-if="isVisible" class="modal-overlay" @click.self="closeDialog">
    <div class="modal-content modal-wide">
      <h3>{{ skill?.name }} verbessern</h3>
      
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
                <small>
                  Verwendet: {{ totalGoldCost || 0 }} GS | 
                  Verbleibend: {{ Math.max(0, (character.vermoegen?.goldstücke || 0) - (totalGoldCost || 0)) }} GS
                </small>
              </div>
              <div class="resource-remaining">
                <small :class="{ 'text-warning': remainingGold < 20, 'text-danger': remainingGold <= 0 }">
                  Nach Lernen: {{ remainingGold }} GS
                </small>
              </div>
            </div>
          </div>
          <div class="resource-display-card">
            <span class="resource-icon">📝</span>
            <div class="resource-info">
              <div class="resource-label">Praxispunkte</div>
              <div class="resource-amount">{{ skill?.pp || 0 }} PP</div>
              <div class="resource-remaining">
                <small>
                  Verwendet: {{ totalPPCost || 0 }} PP | 
                  Verbleibend: {{ Math.max(0, (skill?.pp || 0) - (totalPPCost || 0)) }} PP
                </small>
              </div>
              <div class="resource-remaining">
                <small :class="{ 'text-warning': remainingPP < 5, 'text-danger': remainingPP <= 0 }">
                  Nach Lernen: {{ remainingPP }} PP
                </small>
              </div>
            </div>
          </div>
        </div>      <!-- Belohnungsart, PP-Eingabe und Gold-Eingabe nebeneinander -->
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
          <label>Praxispunkte verwenden:</label>
          <input 
            v-model.number="ppUsed" 
            type="number" 
            min="0" 
            :max="skill?.pp || 0"
            placeholder="PP verwenden"
            @input="updatePPUsage"
          />
          <small class="help-text">
            {{ ppUsed || 0 }} / {{ skill?.pp || 0 }} PP
          </small>
        </div>
        <div class="form-col form-col-input">
          <label>Goldstücke verwenden:</label>
          <input 
            v-model.number="goldUsed" 
            type="number" 
            min="0" 
            :max="character.vermoegen?.goldstücke || 0"
            placeholder="GS verwenden"
            @input="updateGoldUsage"
          />
          <small class="help-text">
            {{ goldUsed || 0 }} / {{ character.vermoegen?.goldstücke || 0 }} GS
          </small>
        </div>
      </div>

      <!-- Lernbare Stufen -->
      <div class="form-group">
        <label>Lernbare Stufen (mehrere auswählbar):</label>
        <div v-if="selectedLevels.length > 0" class="selection-summary">
          <strong>Ausgewählt:</strong> {{ skill?.fertigkeitswert || 0 }} → {{ finalTargetLevel }} 
          ({{ selectedLevels.length }} Level{{ selectedLevels.length !== 1 ? 's' : '' }})
          <br>
          <span class="cost-summary">
            Gesamtkosten: 
            <span v-if="selectedRewardType === 'default'">{{ totalCost }} EP + {{ totalGoldCost }} GS</span>
            <span v-else-if="selectedRewardType === 'noGold'">{{ totalCost }} EP</span>
            <span v-else-if="selectedRewardType === 'halveepnoGold'">{{ Math.floor(totalCost / 2) }} EP</span>
            <span v-else-if="selectedRewardType === 'pp'">{{ totalPPCost }} PP</span>
            <span v-else-if="selectedRewardType === 'mixed'">{{ totalCost }} EP + {{ totalPPCost }} PP</span>
            <span v-else>{{ totalCost }} EP + {{ totalGoldCost }} GS</span>
          </span>
        </div>
        <div class="learning-levels">
          <div 
            v-for="level in availableLevels" 
            :key="level.targetLevel"
            class="level-option"
            :class="{ 
              selected: selectedLevels.includes(level.targetLevel), 
              disabled: !level.canAfford,
              'in-sequence': isInSelectedSequence(level.targetLevel)
            }"
            @click="selectLevel(level)"
          >
            <div class="level-header">
              <span class="level-target">{{ skill?.fertigkeitswert || 0 }} → {{ level.targetLevel }}</span>
              <span class="level-cost" v-if="selectedRewardType === 'default'">{{ level.epCost }} EP + {{ level.goldCost }} GS</span>
              <span class="level-cost" v-else-if="selectedRewardType === 'noGold'">{{ level.epCost }} EP</span>
              <span class="level-cost" v-else-if="selectedRewardType === 'halveepnoGold'">{{ Math.floor(level.epCost / 2) }} EP</span>
              <span class="level-cost" v-else-if="selectedRewardType === 'pp'">{{ level.ppCost }} PP</span>
              <span class="level-cost" v-else-if="selectedRewardType === 'mixed'">{{ level.epCost }} EP + {{ level.ppUsed || 0 }} PP</span>
              <span class="level-cost" v-else>{{ level.epCost }} EP + {{ level.goldCost }} GS</span>
            </div>
            <div class="level-details" v-if="selectedRewardType === 'mixed'">
              <small>PP verwenden: {{ level.ppUsed }} / {{ skill?.pp || 0 }}</small>
            </div>
          </div>
        </div>
      </div>

      <!-- Notizen -->
      <div class="form-group">
        <label>Notizen (optional):</label>
        <textarea v-model="notes" placeholder="Zusätzliche Notizen zum Lernvorgang..."></textarea>
      </div>

      <div class="modal-actions">
        <button 
          @click="executeDetailedLearning" 
          class="btn-confirm" 
          :disabled="!selectedLevels.length || !canAffordSelectedLevels || isLoading"
        >
          {{ isLoading ? 'Wird gelernt...' : `${selectedLevels.length > 1 ? selectedLevels.length + ' Level' : '1 Level'} jetzt lernen` }}
        </button>
        <button @click="closeDialog" class="btn-cancel" :disabled="isLoading">Abbrechen</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  padding: 24px;
  max-width: 500px;
  width: 90%;
  max-height: 80vh;
  overflow-y: auto;
  animation: modalSlideIn 0.3s ease;
}

.modal-wide {
  max-width: 700px;
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

/* Lernbare Stufen */
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
  max-height: 200px;
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

.level-option.in-sequence:not(.selected) {
  background: #f0f8ff;
  border-left: 2px solid #87ceeb;
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

.level-details {
  margin-top: 4px;
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

<script>
import API from '@/utils/api'

export default {
  name: 'SkillImproveDialog',
  props: {
    character: {
      type: Object,
      required: true
    },
    skill: {
      type: Object,
      default: null,
      required: true
    },
    isVisible: {
      type: Boolean,
      default: false
    },
    learningType: {
      type: String,
      default: 'improve', // 'improve', 'learn', 'spell'
      validator: value => ['improve', 'learn', 'spell'].includes(value)
    }
  },
  emits: ['close', 'skill-updated', 'auth-error'],
  data() {
    return {
      selectedRewardType: '',
      selectedLevels: [], // Array von ausgewählten Leveln
      ppUsed: 0,
      goldUsed: 0,
      notes: '',
      availableLevels: [],
      availableRewardTypes: [],
      isLoading: false,
      isLoadingRewardTypes: false
    };
  },
  computed: {
    canAffordSelectedLevels() {
      if (!this.selectedLevels || this.selectedLevels.length === 0) return false;
      
      // Prüfe ob alle ausgewählten Level bezahlbar sind (kumulativ)
      return this.selectedLevels.every(levelNum => {
        const level = this.availableLevels.find(l => l.targetLevel === levelNum);
        return level && level.canAfford;
      }) && this.totalCost <= (this.character.erfahrungsschatz?.ep || 0) &&
             this.totalGoldCost <= (this.character.vermoegen?.goldstücke || 0) &&
             this.totalPPCost <= (this.skill?.pp || 0);
    },
    
    totalCost() {
      return this.selectedLevels.reduce((sum, levelNum) => {
        const level = this.availableLevels.find(l => l.targetLevel === levelNum);
        return sum + (level ? level.epCost : 0);
      }, 0);
    },
    
    totalGoldCost() {
      return this.selectedLevels.reduce((sum, levelNum) => {
        const level = this.availableLevels.find(l => l.targetLevel === levelNum);
        return sum + (level ? level.goldCost : 0);
      }, 0);
    },
    
    totalPPCost() {
      return this.selectedLevels.reduce((sum, levelNum) => {
        const level = this.availableLevels.find(l => l.targetLevel === levelNum);
        if (!level) return sum;
        
        // PP werden nur bei bestimmten Belohnungstypen tatsächlich verwendet
        switch (this.selectedRewardType) {
          case 'pp':
            return sum + (level.ppCost || 0);
          case 'mixed':
            return sum + (level.ppUsed || 0);
          default:
            return sum; // Bei anderen Belohnungstypen werden keine PP verwendet
        }
      }, 0);
    },
    
    finalTargetLevel() {
      if (this.selectedLevels.length === 0) return this.skill?.fertigkeitswert || 0;
      return Math.max(...this.selectedLevels);
    },
    
    selectedCost() {
      return this.totalCost;
    },
    
    selectedGoldCost() {
      return this.totalGoldCost;
    },
    
    selectedPPCost() {
      return this.totalPPCost;
    },
    
    remainingEP() {
      const current = this.character.erfahrungsschatz?.ep || 0;
      return Math.max(0, current - this.totalCost);
    },
    
    remainingGold() {
      const current = this.character.vermoegen?.goldstücke || 0;
      return Math.max(0, current - this.totalGoldCost);
    },
    
    remainingPP() {
      const current = this.skill?.pp || 0;
      return Math.max(0, current - this.totalPPCost);
    }
  },
  watch: {
    skill: {
      handler(newSkill) {
        if (newSkill) {
          this.loadRewardTypes();
          // loadLearningCosts wird durch selectedRewardType watcher ausgelöst
        }
      },
      immediate: true
    },
    learningType: {
      handler() {
        if (this.skill) {
          this.loadRewardTypes();
          // loadLearningCosts wird durch selectedRewardType watcher ausgelöst
        }
      },
      immediate: true
    },
    selectedRewardType() {
      // Nur hier die Lernkosten laden, wenn Belohnungsart geändert wird
      if (this.selectedRewardType && this.skill) {
        this.loadLearningCosts();
      }
      this.selectedLevels = []; // Reset selection when reward type changes
    },
    isVisible(newValue) {
      if (newValue) {
        this.resetDialog();
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
    
    resetDialog() {
      this.selectedRewardType = '';
      this.selectedLevels = [];
      this.ppUsed = 0;
      this.goldUsed = 0;
      this.notes = '';
      this.availableLevels = [];
      this.availableRewardTypes = [];
      if (this.skill) {
        this.loadRewardTypes();
        // loadLearningCosts wird automatisch durch selectedRewardType watcher ausgelöst
      }
    },
    
    async loadRewardTypes() {
      if (!this.skill) return;
      
      // Prüfe ob Token vorhanden ist
      const token = localStorage.getItem('token');
      if (!token) {
        console.error('No authentication token available - cannot load reward types');
        this.availableRewardTypes = [{
          value: 'error',
          label: 'Anmeldung erforderlich'
        }];
        return;
      }
      
      this.isLoadingRewardTypes = true;
      try {
        console.log('Loading reward types for:', {
          character_id: this.character.id,
          learning_type: this.learningType,
          skill_name: this.skill.name,
          current_level: this.skill.fertigkeitswert || 0,
          skill_type: this.skill.type || 'skill',
          has_token: !!token
        });
        
        // API-Endpunkt für verfügbare Belohnungsarten
        const response = await this.$api.get(`/api/characters/${this.character.id}/reward-types`, {
          params: {
            learning_type: this.learningType,
            skill_name: this.skill.name,
            current_level: this.skill.fertigkeitswert || 0,
            skill_type: this.skill.type || 'skill' // 'skill', 'weapon', 'spell'
          }
        });
        
        console.log('API Response:', response.data);
        console.log('Reward types from API:', response.data.reward_types);
        
        this.availableRewardTypes = response.data.reward_types || [];
        
        // Setze den ersten verfügbaren Belohnungstyp als Standard
        if (this.availableRewardTypes.length > 0 && !this.selectedRewardType) {
          this.selectedRewardType = this.availableRewardTypes[0].value;
          console.log('Set default reward type to:', this.selectedRewardType);
          // loadLearningCosts wird automatisch durch selectedRewardType watcher ausgelöst
        } else if (this.availableRewardTypes.length === 0) {
          console.error('No reward types received from API - cannot proceed without proper reward types');
          // Zeige Fehlermeldung statt Fallback zu verwenden
          this.availableRewardTypes = [{
            value: 'error',
            label: 'Fehler beim Laden der Belohnungsarten'
          }];
        }
        
      } catch (error) {
        console.error('Fehler beim Laden der Belohnungsarten:', error);
        console.error('Error details:', {
          status: error.response?.status,
          statusText: error.response?.statusText,
          data: error.response?.data,
          url: error.config?.url,
          headers: error.config?.headers
        });
        
        // Spezielle Behandlung für Auth-Fehler
        if (error.response?.status === 401) {
          console.error('Authentication failed - token may be invalid or expired');
          console.log('Current token:', token ? 'Present' : 'Missing');
          
          // Optional: Token entfernen wenn ungültig
          // localStorage.removeItem('token');
          
          // Event emittieren um Parent über Auth-Problem zu informieren
          this.$emit('auth-error', 'Authentication required for reward types');
        } else if (error.response?.status === 404) {
          console.warn('Reward types endpoint not found, using fallback');
        } else if (error.response?.status === 403) {
          console.error('Access forbidden - insufficient permissions');
        }
        
        // Bei Fehlern klare Fehlermeldung statt Fallback
        this.availableRewardTypes = [{
          value: 'error',
          label: 'Fehler beim Laden der Belohnungsarten'
        }];
        console.error('Could not load reward types from API, showing error state');
      } finally {
        this.isLoadingRewardTypes = false;
      }
    },
    
    calculateAvailableLevels() {
      if (!this.skill) return;
      
      // Diese Methode ist jetzt redundant, da loadLearningCosts() 
      // automatisch durch den selectedRewardType watcher aufgerufen wird
      console.warn('calculateAvailableLevels() ist veraltet - verwende selectedRewardType watcher');
    },
    
    async loadLearningCosts() {
      if (!this.skill || !this.selectedRewardType) return;
      
      this.isLoading = true;
      try {
        const token = localStorage.getItem('token');
        if (!token) {
          console.warn('No authentication token available for cost calculation');
          this.generateFallbackLevels();
          return;
        }
        
        console.log('Loading learning costs for:', {
          character_id: this.character.id,
          skill_name: this.skill.name,
          skill_type: this.skill.type || 'skill',
          learning_type: this.learningType,
          current_level: this.skill.fertigkeitswert || 0,
          reward_type: this.selectedRewardType
        });
        
        // Verwende den neuen /lerncost Endpunkt mit gsmaster.LernCostRequest Struktur
        const requestData = {
          char_id: parseInt(this.character.id),
          name: this.skill.name,
          current_level: this.skill.fertigkeitswert || 0,
          type: this.skill.type || 'skill',
          action: this.learningType === 'learn' ? 'learn' : 'improve',
          target_level: 0, // Wird vom Backend automatisch bis Level 18 berechnet
          use_pp: this.ppUsed || 0,
          use_gold: this.goldUsed || 0,
          reward: this.selectedRewardType
        };
        
        const response = await this.$api.post(`/api/characters/lerncost`, requestData);
        
        console.log('Learning costs API response:', response.data);
        
        if (response.data && Array.isArray(response.data) && response.data.length > 0) {
          // Konvertiere gsmaster.SkillCostResultNew Array zu unserem internen Format
          const availableEP = this.character.erfahrungsschatz?.ep || 0;
          const availableGold = this.character.vermoegen?.goldstücke || 0;
          const availablePP = this.skill?.pp || 0;
          
          this.availableLevels = response.data.map(cost => {
            // Backend liefert bereits die korrekten Kosten basierend auf dem Belohnungstyp
            let canAfford = false;
            
            switch (this.selectedRewardType) {
              case 'noGold':
              case 'halveepnoGold':
                canAfford = availableEP >= cost.ep;
                break;
              case 'pp':
                canAfford = availablePP >= cost.le;
                break;
              case 'mixed':
                canAfford = availableEP >= cost.ep && availablePP >= (cost.pp_used || 0);
                break;
              case 'default':
              default:
                canAfford = availableEP >= cost.ep && availableGold >= cost.gold_cost;
                break;
            }
            
            return {
              targetLevel: cost.target_level,
              epCost: cost.ep,
              goldCost: cost.gold_cost,
              ppCost: cost.le,
              ppUsed: cost.pp_used || 0,
              canAfford: canAfford
            };
          });
          
          console.log('Processed level costs for reward type', this.selectedRewardType, ':', this.availableLevels);
        } else {
          console.warn('No level costs returned from API, using fallback');
          this.generateFallbackLevels();
        }
        
      } catch (error) {
        console.error('Fehler beim Laden der Lernkosten:', error);
        console.error('Error details:', {
          status: error.response?.status,
          statusText: error.response?.statusText,
          data: error.response?.data
        });
        
        if (error.response?.status === 401) {
          console.error('Authentication failed for learning costs');
          this.$emit('auth-error', 'Authentication required for learning costs');
        }
        
        // Fallback auf berechnete Kosten
        this.generateFallbackLevels();
        
      } finally {
        this.isLoading = false;
      }
    },
    
    generateFallbackLevels() {
      // Einfache Fallback-Methode für Kostenberechnung (nur für Notfälle)
      const currentLevel = this.skill.fertigkeitswert || 0;
      const maxLevel = 20;
      const availableEP = this.character.erfahrungsschatz?.ep || 0;
      const availableGold = this.character.vermoegen?.goldstücke || 0;
      const availablePP = this.skill?.pp || 0;
      
      this.availableLevels = [];
      
      for (let targetLevel = currentLevel + 1; targetLevel <= Math.min(currentLevel + 5, maxLevel); targetLevel++) {
        const levelDiff = targetLevel - currentLevel;
        
        // Sehr einfache Basis-Kosten (nur als Fallback)
        const epCost = levelDiff * 100;
        const goldCost = levelDiff * 50;
        const ppCost = levelDiff * 20;
        
        // Verfügbarkeit basierend auf Belohnungstyp
        let canAfford = false;
        switch (this.selectedRewardType) {
          case 'noGold':
          case 'halveepnoGold':
            canAfford = availableEP >= epCost;
            break;
          case 'pp':
            canAfford = availablePP >= ppCost;
            break;
          case 'mixed':
            canAfford = availableEP >= epCost && availablePP >= Math.min(levelDiff * 5, availablePP);
            break;
          case 'default':
          default:
            canAfford = availableEP >= epCost && availableGold >= goldCost;
            break;
        }
        
        this.availableLevels.push({
          targetLevel,
          epCost,
          goldCost,
          ppCost,
          ppUsed: 0,
          canAfford
        });
      }
    },
    
    selectLevel(level) {
      if (!level.canAfford) return;
      
      const targetLevel = level.targetLevel;
      const currentLevel = this.skill?.fertigkeitswert || 0;
      
      // Toggle-Logik für Multi-Level-Auswahl
      const isSelected = this.selectedLevels.includes(targetLevel);
      
      if (isSelected) {
        // Entferne das Level und alle höheren Level
        this.selectedLevels = this.selectedLevels.filter(l => l < targetLevel);
      } else {
        // Füge das Level hinzu und alle Level zwischen dem aktuellen Level und diesem Level
        const levelsToAdd = [];
        
        // Bestimme den niedrigsten bereits ausgewählten Level oder das aktuelle Level
        const minSelectedLevel = this.selectedLevels.length > 0 
          ? Math.min(...this.selectedLevels) 
          : currentLevel + 1;
        
        // Füge alle Level vom niedrigsten bis zum angeklickten Level hinzu
        for (let i = Math.min(minSelectedLevel, targetLevel); i <= targetLevel; i++) {
          if (i > currentLevel && !this.selectedLevels.includes(i)) {
            // Prüfe ob das Level verfügbar und bezahlbar ist
            const levelData = this.availableLevels.find(l => l.targetLevel === i);
            if (levelData && levelData.canAfford) {
              levelsToAdd.push(i);
            }
          }
        }
        
        // Aktualisiere die Auswahl
        this.selectedLevels = [...new Set([...this.selectedLevels, ...levelsToAdd])].sort((a, b) => a - b);
      }
      
      console.log('Selected levels:', this.selectedLevels);
    },
    
    isInSelectedSequence(targetLevel) {
      // Prüft ob ein Level Teil der ausgewählten Sequenz ist (visueller Hinweis)
      if (this.selectedLevels.length === 0) return false;
      
      const minSelected = Math.min(...this.selectedLevels);
      const maxSelected = Math.max(...this.selectedLevels);
      
      return targetLevel >= minSelected && targetLevel <= maxSelected;
    },
    
    updatePPUsage() {
      // Stelle sicher, dass PP-Verwendung die verfügbaren PP nicht überschreitet
      const maxPP = this.skill?.pp || 0;
      if (this.ppUsed > maxPP) {
        this.ppUsed = maxPP;
      }
      if (this.ppUsed < 0) {
        this.ppUsed = 0;
      }
      
      // Lernkosten immer neu laden, da PP-Verwendung die Kosten beeinflussen kann
      if (this.selectedRewardType && this.skill) {
        this.loadLearningCosts();
      }
    },
    
    updateGoldUsage() {
      // Stelle sicher, dass Gold-Verwendung die verfügbaren GS nicht überschreitet
      const maxGold = this.character.vermoegen?.goldstücke || 0;
      if (this.goldUsed > maxGold) {
        this.goldUsed = maxGold;
      }
      if (this.goldUsed < 0) {
        this.goldUsed = 0;
      }
      
      // Lernkosten immer neu laden, da Gold-Verwendung die Kosten beeinflussen kann
      if (this.selectedRewardType && this.skill) {
        this.loadLearningCosts();
      }
    },
    
    updateMixedCosts() {
      // Diese Methode ist jetzt redundant, da updatePPUsage() alles übernimmt
      this.updatePPUsage();
    },
    
    async executeDetailedLearning() {
      if (!this.skill || !this.selectedLevels.length) {
        alert('Bitte wählen Sie mindestens eine Zielstufe aus.');
        return;
      }
      
      if (!this.selectedRewardType) {
        alert('Bitte wählen Sie eine Belohnungsart aus.');
        return;
      }
      
      if (!this.canAffordSelectedLevels) {
        alert('Sie haben nicht genügend Ressourcen für diese Verbesserung(en).');
        return;
      }
      
      this.isLoading = true;
      try {
        // Für Multi-Level-Learning senden wir das höchste Level als Ziel
        const finalLevel = Math.max(...this.selectedLevels);
        
        const requestData = {
          char_id: this.character.id,
          name: this.skill.name,
          current_level: this.skill.fertigkeitswert,
          target_level: finalLevel, // Höchstes ausgewähltes Level
          type: this.learningType === 'spell' ? 'spell' : 'skill',
          action: this.learningType === 'learn' ? 'learn' : 'improve',
          reward: this.selectedRewardType,
          use_pp: this.ppUsed || 0,
          use_gold: this.goldUsed || 0,
          levels_to_learn: this.selectedLevels, // Alle ausgewählten Level
          notes: this.notes || `${this.learningType === 'spell' ? 'Zauber' : 'Fertigkeit'} ${this.skill.name} von ${this.skill.fertigkeitswert} auf ${finalLevel} ${this.learningType === 'learn' ? 'gelernt' : 'verbessert'} (${this.selectedLevels.length} Level)`
        };
        
        // Wähle den richtigen API-Endpunkt basierend auf Lerntyp
        let endpoint;
        switch (this.learningType) {
          case 'learn':
            endpoint = `/api/characters/${this.character.id}/learn-skill`;
            break;
          case 'spell':
            endpoint = `/api/characters/${this.character.id}/improve-spell`;
            break;
          case 'improve':
          default:
            endpoint = `/api/characters/improve-skill`;
            break;
        }
        
        const response = await this.$api.post(endpoint, requestData);
        
        console.log(`${this.learningType} erfolgreich ausgeführt:`, response.data);
        alert(`${this.learningType === 'spell' ? 'Zauber' : 'Fertigkeit'} "${this.skill.name}" erfolgreich ${this.learningType === 'learn' ? 'gelernt' : 'auf Stufe ' + finalLevel + ' verbessert'} (${this.selectedLevels.length} Level)!`);
        this.$emit('skill-updated');
        this.closeDialog();
        
      } catch (error) {
        console.error(`Fehler beim ${this.learningType}:`, error);
        alert(`Fehler beim ${this.learningType === 'learn' ? 'Lernen' : 'Verbessern'}: ` + (error.response?.data?.error || error.message));
      } finally {
        this.isLoading = false;
      }
    }
  }
};
</script>
