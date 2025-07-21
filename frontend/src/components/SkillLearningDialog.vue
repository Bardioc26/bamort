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
            <div class="resource-amount">{{ character.erfahrungsschatz?.value || 0 }} EP</div>
          </div>
        </div>
        <div class="resource-display-card">
          <span class="resource-icon">💰</span>
          <div class="resource-info">
            <div class="resource-label">Gold</div>
            <div class="resource-amount">{{ character.vermoegen?.goldstücke || 0 }} GS</div>
          </div>
        </div>
        <div class="resource-display-card">
          <span class="resource-icon">📝</span>
          <div class="resource-info">
            <div class="resource-label">Praxispunkte</div>
            <div class="resource-amount">{{ skill?.pp || 0 }} PP</div>
          </div>
        </div>
      </div>

      <!-- Belohnungsart auswählen -->
      <div class="form-group">
        <label>Belohnungsart:</label>
        <select v-model="selectedRewardType">
          <option value="ep">Erfahrungspunkte verwenden</option>
          <option value="gold">Gold verwenden</option>
          <option value="pp">Praxispunkte verwenden</option>
          <option value="mixed">Gemischt (EP + PP)</option>
        </select>
      </div>

      <!-- Lernbare Stufen -->
      <div class="form-group">
        <label>Lernbare Stufen:</label>
        <div class="learning-levels">
          <div 
            v-for="level in availableLevels" 
            :key="level.targetLevel"
            class="level-option"
            :class="{ selected: selectedLevel === level.targetLevel, disabled: !level.canAfford }"
            @click="selectLevel(level)"
          >
            <div class="level-header">
              <span class="level-target">{{ skill?.fertigkeitswert || 0 }} → {{ level.targetLevel }}</span>
              <span class="level-cost" v-if="selectedRewardType === 'ep'">{{ level.epCost }} EP</span>
              <span class="level-cost" v-else-if="selectedRewardType === 'gold'">{{ level.goldCost }} GS</span>
              <span class="level-cost" v-else-if="selectedRewardType === 'pp'">{{ level.ppCost }} PP</span>
              <span class="level-cost" v-else>{{ level.epCost }} EP + {{ level.ppUsed }} PP</span>
            </div>
            <div class="level-details" v-if="selectedRewardType === 'mixed'">
              <small>PP verwenden: {{ level.ppUsed }} / {{ skill?.pp || 0 }}</small>
            </div>
          </div>
        </div>
      </div>

      <!-- PP Eingabe für gemischte Belohnung -->
      <div v-if="selectedRewardType === 'mixed'" class="form-group">
        <label>Praxispunkte verwenden (optional):</label>
        <input 
          v-model.number="ppUsed" 
          type="number" 
          min="0" 
          :max="skill?.pp || 0"
          placeholder="Anzahl PP verwenden"
          @input="updateMixedCosts"
        />
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
          :disabled="!selectedLevel || !canAffordSelectedLevel || isLoading"
        >
          {{ isLoading ? 'Wird gelernt...' : 'Jetzt lernen' }}
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

/* Lernbare Stufen */
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
  name: 'SkillLearningDialog',
  props: {
    character: {
      type: Object,
      required: true
    },
    skill: {
      type: Object,
      default: null
    },
    isVisible: {
      type: Boolean,
      default: false
    }
  },
  emits: ['close', 'skill-updated'],
  data() {
    return {
      selectedRewardType: 'ep',
      selectedLevel: null,
      ppUsed: 0,
      notes: '',
      availableLevels: [],
      isLoading: false
    };
  },
  computed: {
    canAffordSelectedLevel() {
      if (!this.selectedLevel || !this.skill) return false;
      
      const level = this.availableLevels.find(l => l.targetLevel === this.selectedLevel);
      return level ? level.canAfford : false;
    }
  },
  watch: {
    skill: {
      handler(newSkill) {
        if (newSkill) {
          this.calculateAvailableLevels();
        }
      },
      immediate: true
    },
    selectedRewardType() {
      this.updateAffordability();
      this.selectedLevel = null; // Reset selection when reward type changes
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
      this.selectedRewardType = 'ep';
      this.selectedLevel = null;
      this.ppUsed = 0;
      this.notes = '';
      this.availableLevels = [];
      if (this.skill) {
        this.calculateAvailableLevels();
      }
    },
    
    calculateAvailableLevels() {
      if (!this.skill) return;
      
      const currentLevel = this.skill.fertigkeitswert || 0;
      const maxLevel = 20; // Maximaler Fertigkeitswert
      const availableEP = this.character.erfahrungsschatz?.value || 0;
      const availableGold = this.character.vermoegen?.goldstücke || 0;
      const availablePP = this.skill.pp || 0;
      
      this.availableLevels = [];
      
      for (let targetLevel = currentLevel + 1; targetLevel <= Math.min(currentLevel + 5, maxLevel); targetLevel++) {
        const levelDiff = targetLevel - currentLevel;
        
        // Basis-Kosten (diese sollten später aus der API kommen)
        const baseEPCost = levelDiff * 100; // Beispiel: 100 EP pro Stufe
        const baseGoldCost = levelDiff * 50; // Beispiel: 50 Gold pro Stufe
        const ppReduction = Math.floor(levelDiff * 10); // 10 EP Reduktion pro PP
        
        // Kosten berechnen
        const epCost = Math.max(baseEPCost - (availablePP * 10), baseEPCost * 0.5);
        const goldCost = baseGoldCost;
        const ppCost = Math.min(levelDiff * 2, availablePP); // 2 PP pro Stufe
        
        this.availableLevels.push({
          targetLevel,
          epCost: Math.ceil(epCost),
          goldCost,
          ppCost,
          ppUsed: 0,
          canAfford: {
            ep: availableEP >= epCost,
            gold: availableGold >= goldCost,
            pp: availablePP >= ppCost
          }
        });
      }
      
      // Aktualisiere canAfford basierend auf dem gewählten Belohnungstyp
      this.updateAffordability();
    },
    
    updateAffordability() {
      const availableEP = this.character.erfahrungsschatz?.value || 0;
      const availableGold = this.character.vermoegen?.goldstücke || 0;
      const availablePP = this.skill?.pp || 0;
      
      this.availableLevels.forEach(level => {
        switch (this.selectedRewardType) {
          case 'ep':
            level.canAfford = availableEP >= level.epCost;
            break;
          case 'gold':
            level.canAfford = availableGold >= level.goldCost;
            break;
          case 'pp':
            level.canAfford = availablePP >= level.ppCost;
            break;
          case 'mixed':
            const adjustedEPCost = Math.max(level.epCost - (this.ppUsed * 10), level.epCost * 0.5);
            level.canAfford = availableEP >= adjustedEPCost && availablePP >= this.ppUsed;
            level.ppUsed = this.ppUsed;
            break;
        }
      });
    },
    
    selectLevel(level) {
      if (level.canAfford) {
        this.selectedLevel = level.targetLevel;
      }
    },
    
    updateMixedCosts() {
      this.updateAffordability();
    },
    
    async executeDetailedLearning() {
      if (!this.skill || !this.selectedLevel) {
        alert('Bitte wählen Sie eine Zielstufe aus.');
        return;
      }
      
      const selectedLevelData = this.availableLevels.find(l => l.targetLevel === this.selectedLevel);
      if (!selectedLevelData || !selectedLevelData.canAfford) {
        alert('Sie haben nicht genügend Ressourcen für diese Verbesserung.');
        return;
      }
      
      this.isLoading = true;
      try {
        const requestData = {
          name: this.skill.name,
          current_level: this.skill.fertigkeitswert,
          target_level: this.selectedLevel,
          reward_type: this.selectedRewardType,
          use_pp: this.selectedRewardType === 'mixed' ? this.ppUsed : 
                  this.selectedRewardType === 'pp' ? selectedLevelData.ppCost : 0,
          notes: this.notes || `Fertigkeit ${this.skill.name} von ${this.skill.fertigkeitswert} auf ${this.selectedLevel} verbessert`
        };
        
        const response = await this.$api.post(`/api/characters/${this.character.id}/improve-skill`, requestData);
        
        console.log('Fertigkeit erfolgreich verbessert:', response.data);
        alert(`Fertigkeit "${this.skill.name}" erfolgreich auf Stufe ${this.selectedLevel} verbessert!`);
        this.$emit('skill-updated');
        this.closeDialog();
        
      } catch (error) {
        console.error('Fehler beim Verbessern der Fertigkeit:', error);
        alert('Fehler beim Verbessern der Fertigkeit: ' + (error.response?.data?.error || error.message));
      } finally {
        this.isLoading = false;
      }
    }
  }
};
</script>
