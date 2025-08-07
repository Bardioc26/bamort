<template>
  <div class="experiance-view">
    <h3>{{ $t('experience.title') }}</h3>
    
    <!-- Erfahrungspunkte -->
    <div class="experience-section">
      <h4>{{ $t('experience.experience_points') }}</h4>
      <div class="stat-box">
        <div class="stat-item">
          <span class="stat-label">{{ $t('experience.available_ep') }}:</span>
          <span class="stat-value">{{ character.erfahrungsschatz?.ep || 0 }} EP</span>
        </div>
        <div class="control-row">
          <div class="input-group">
            <input 
              v-model.number="experienceAmount" 
              type="number" 
              class="amount-input"
              placeholder="Anzahl EP"
              min="1"
            />
            <div class="button-group">
              <button @click="addExperience" class="btn-add" :disabled="!experienceAmount || experienceAmount <= 0 || isLoading">
                <span v-if="isLoading">⏳</span>
                <span v-else>+ Hinzufügen</span>
              </button>
              <button @click="removeExperience" class="btn-remove" :disabled="!experienceAmount || experienceAmount <= 0 || isLoading">
                <span v-if="isLoading">⏳</span>
                <span v-else>- Entfernen</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Vermögen -->
    <div class="wealth-section">
      <h4>{{ $t('experience.wealth') }}</h4>
      <div class="stat-box">
        <div class="wealth-item">
          <span class="wealth-label">{{ $t('experience.gold_coins') }}:</span>
          <span class="wealth-value">{{ character.vermoegen?.goldstücke || 0 }} GS</span>
        </div>
        <div class="control-row">
          <div class="input-group">
            <input 
              v-model.number="goldAmount" 
              type="number" 
              class="amount-input"
              placeholder="Anzahl GS"
              min="1"
            />
            <div class="button-group">
              <button @click="addGold" class="btn-add" :disabled="!goldAmount || goldAmount <= 0 || isLoading">
                <span v-if="isLoading">⏳</span>
                <span v-else>+ Hinzufügen</span>
              </button>
              <button @click="removeGold" class="btn-remove" :disabled="!goldAmount || goldAmount <= 0 || isLoading">
                <span v-if="isLoading">⏳</span>
                <span v-else>- Entfernen</span>
              </button>
            </div>
          </div>
        </div>
        <!-- Silberstücke und Kupferstücke temporär ausgeblendet -->
        <!-- 
        <div class="wealth-item">
          <span class="wealth-label">{{ $t('experience.silver_coins') }}:</span>
          <span class="wealth-value">{{ character.vermoegen?.silberstücke || 0 }} SS</span>
        </div>
        <div class="wealth-item">
          <span class="wealth-label">{{ $t('experience.copper_coins') }}:</span>
          <span class="wealth-value">{{ character.vermoegen?.kupferstücke || 0 }} KS</span>
        </div>
        -->
        <div class="wealth-item total">
          <span class="wealth-label">{{ $t('experience.total_in_gs') }}:</span>
          <span class="wealth-value">{{ totalWealthInGS }} GS</span>
        </div>
      </div>
    </div>
    
    <!-- Audit Log -->
    <AuditLogView 
      v-if="character" 
      :character="character"
    />
  </div>
</template>

<style scoped>
.experiance-view {
  padding: 20px;
  max-width: 800px;
}

.experiance-view h3 {
  color: #333;
  border-bottom: 2px solid #007bff;
  padding-bottom: 10px;
  margin-bottom: 20px;
}

.experiance-view h4 {
  color: #555;
  margin-bottom: 15px;
  margin-top: 25px;
}

.stat-box {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 20px;
}

.stat-item, .wealth-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #e9ecef;
}

.stat-item:last-child, .wealth-item:last-child {
  border-bottom: none;
}

.wealth-item.total {
  border-top: 2px solid #007bff;
  margin-top: 10px;
  padding-top: 15px;
  font-weight: bold;
  color: #007bff;
}

.stat-label, .wealth-label {
  font-weight: 500;
  color: #555;
}

.stat-value, .wealth-value {
  font-weight: bold;
  color: #333;
  background: #fff;
  padding: 5px 10px;
  border-radius: 4px;
  border: 1px solid #ddd;
}

.wealth-item.total .wealth-value {
  background: #007bff;
  color: white;
  border-color: #007bff;
}

.value-controls {
  display: flex;
  align-items: center;
  gap: 5px;
}

.control-row {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #e9ecef;
}

.input-group {
  display: flex;
  gap: 10px;
  align-items: center;
}

.amount-input {
  font-weight: bold;
  color: #333;
  background: #fff;
  padding: 8px 12px;
  border-radius: 4px;
  border: 1px solid #ddd;
  width: 120px;
  text-align: right;
  transition: border-color 0.2s;
}

.amount-input:focus {
  border-color: #007bff;
  outline: none;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
}

.button-group {
  display: flex;
  gap: 8px;
}

.btn-add, .btn-remove {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
}

.btn-add {
  background: #28a745;
  color: white;
}

.btn-add:hover:not(:disabled) {
  background: #218838;
}

.btn-remove {
  background: #dc3545;
  color: white;
}

.btn-remove:hover:not(:disabled) {
  background: #c82333;
}

.btn-add:disabled, .btn-remove:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.unit {
  font-weight: bold;
  color: #666;
}

.experience-section, .wealth-section {
  margin-bottom: 30px;
}
</style>


<script>
import API from '@/utils/api'
import AuditLogView from './AuditLogView.vue'

export default {
  name: "ExperianceView",
  components: {
    AuditLogView
  },
  props: {
    character: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      experienceAmount: null,
      goldAmount: null,
      isLoading: false,
      testMode: false // Für Debugging - setzt auf true um Backend zu umgehen
    };
  },
  created() {
    this.$api = API;
  },
  computed: {
    totalWealthInGS() {
      const vermoegen = this.character.vermoegen || {};
      const goldstücke = vermoegen.goldstücke || 0;    // GS
      const silberstücke = vermoegen.silberstücke || 0; // SS
      const kupferstücke = vermoegen.kupferstücke || 0; // KS
      
      // Midgard Währungsumrechnung: 1 GS = 10 SS = 10 KS
      // Alles in Goldstücke umrechnen
      return goldstücke + Math.floor(silberstücke / 10) + Math.floor(kupferstücke / 10);
    }
  },
  methods: {
    async addExperience() {
      if (!this.experienceAmount || this.experienceAmount <= 0 || this.isLoading) return;
      
      this.isLoading = true;
      try {
        const currentEP = this.character.erfahrungsschatz?.ep || 0;
        const newEP = currentEP + this.experienceAmount;
        
        await this.updateExperience(newEP);
        this.experienceAmount = null;
      } catch (error) {
        console.error('Fehler beim Hinzufügen der Erfahrungspunkte:', error);
      } finally {
        this.isLoading = false;
      }
    },

    async removeExperience() {
      if (!this.experienceAmount || this.experienceAmount <= 0 || this.isLoading) return;
      
      this.isLoading = true;
      try {
        const currentEP = this.character.erfahrungsschatz?.ep || 0;
        const newEP = Math.max(0, currentEP - this.experienceAmount);
        
        await this.updateExperience(newEP);
        this.experienceAmount = null;
      } catch (error) {
        console.error('Fehler beim Entfernen der Erfahrungspunkte:', error);
      } finally {
        this.isLoading = false;
      }
    },

    async addGold() {
      if (!this.goldAmount || this.goldAmount <= 0 || this.isLoading) return;
      
      this.isLoading = true;
      try {
        const currentGold = this.character.vermoegen?.goldstücke || 0;
        const newGold = currentGold + this.goldAmount;
        
        await this.updateGold(newGold);
        this.goldAmount = null;
      } catch (error) {
        console.error('Fehler beim Hinzufügen der Goldstücke:', error);
      } finally {
        this.isLoading = false;
      }
    },

    async removeGold() {
      if (!this.goldAmount || this.goldAmount <= 0 || this.isLoading) return;
      
      this.isLoading = true;
      try {
        const currentGold = this.character.vermoegen?.goldstücke || 0;
        const newGold = Math.max(0, currentGold - this.goldAmount);
        
        await this.updateGold(newGold);
        this.goldAmount = null;
      } catch (error) {
        console.error('Fehler beim Entfernen der Goldstücke:', error);
      } finally {
        this.isLoading = false;
      }
    },

    async updateExperience(newValue) {
      try {
        console.log('Updating experience to:', newValue);
        
        if (!this.testMode) {
          // API-Call zum Speichern der Erfahrungspunkte
          const response = await this.$api.put(`/api/characters/${this.character.id}/experience`, {
            experience_points: newValue,
            reason: "manual",
            notes: `Manual adjustment: ${this.experienceAmount > 0 ? 'Added' : 'Removed'} ${Math.abs(this.experienceAmount || 0)} EP`
          });
          
          console.log('Experience update response:', response.data);
        } else {
          console.log('Test mode - skipping API call');
          // Simuliere einen kurzen delay
          await new Promise(resolve => setTimeout(resolve, 500));
        }

        // Direkte Aktualisierung der lokalen Daten
        if (this.character.erfahrungsschatz) {
          this.character.erfahrungsschatz.value = newValue;
        } else {
          this.$set(this.character, 'erfahrungsschatz', { value: newValue });
        }

        // Emit event to parent component to refresh character data
        this.$emit('character-updated');
        
      } catch (error) {
        console.error('Fehler beim Speichern der Erfahrungspunkte:', error);
        console.error('Error details:', error.response);
        
        if (error.response?.status === 401) {
          alert('Authentifizierung fehlgeschlagen. Bitte loggen Sie sich erneut ein.');
        } else {
          alert('Fehler beim Speichern der Erfahrungspunkte: ' + (error.response?.data?.error || error.message));
        }
        throw error; // Re-throw to be caught by calling function
      }
    },

    async updateGold(newValue) {
      try {
        console.log('Updating gold to:', newValue);
        
        if (!this.testMode) {
          // API-Call zum Speichern der Goldstücke
          const response = await this.$api.put(`/api/characters/${this.character.id}/wealth`, {
            goldstücke: newValue,
            reason: "manual",
            notes: `Manual adjustment: ${this.goldAmount > 0 ? 'Added' : 'Removed'} ${Math.abs(this.goldAmount || 0)} GS`
          });
          
          console.log('Gold update response:', response.data);
        } else {
          console.log('Test mode - skipping API call');
          // Simuliere einen kurzen delay
          await new Promise(resolve => setTimeout(resolve, 500));
        }

        // Direkte Aktualisierung der lokalen Daten
        if (this.character.vermoegen) {
          this.character.vermoegen.goldstücke = newValue;
        } else {
          this.$set(this.character, 'vermoegen', { goldstücke: newValue, silberstücke: 0, kupferstücke: 0 });
        }

        // Emit event to parent component to refresh character data
        this.$emit('character-updated');
        
      } catch (error) {
        console.error('Fehler beim Speichern der Goldstücke:', error);
        console.error('Error details:', error.response);
        
        if (error.response?.status === 401) {
          alert('Authentifizierung fehlgeschlagen. Bitte loggen Sie sich erneut ein.');
        } else {
          alert('Fehler beim Speichern der Goldstücke: ' + (error.response?.data?.error || error.message));
        }
        throw error; // Re-throw to be caught by calling function
      }
    }
  }
};
</script>
