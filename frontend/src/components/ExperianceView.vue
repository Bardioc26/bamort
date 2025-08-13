<template>
    <div class="fullwidth-container">
    
    <!-- Erfahrungspunkte -->
    <div class="experience-section">
      <div class="section-header">
        <h4>{{ $t('experience.experience_points') }}</h4>
      </div>
      <div class="resource-display">
        <div class="resource-card">
          <span class="resource-icon">⚡</span>
          <div class="resource-info">
            <div class="resource-label">{{ $t('experience.available_ep') }}</div>
            <div class="resource-amount">{{ character.erfahrungsschatz?.ep || 0 }} EP</div>
          </div>
        </div>
        <div class="form-row control-row">
          <div class="form-group">
            <input 
              v-model.number="experienceAmount" 
              type="number" 
              class="form-control amount-input"
              placeholder="Anzahl EP"
              min="1"
            />
          </div>
          <div class="button-group">
            <button @click="addExperience" class="btn btn-success" :disabled="!experienceAmount || experienceAmount <= 0 || isLoading">
              <span v-if="isLoading">⏳</span>
              <span v-else>+ Hinzufügen</span>
            </button>
            <button @click="removeExperience" class="btn btn-danger" :disabled="!experienceAmount || experienceAmount <= 0 || isLoading">
              <span v-if="isLoading">⏳</span>
              <span v-else>- Entfernen</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Vermögen -->
    <div class="wealth-section">
      <div class="section-header">
        <h4>{{ $t('experience.wealth') }}</h4>
      </div>
      <div class="resource-display">
        <div class="resource-card">
          <span class="resource-icon">💰</span>
          <div class="resource-info">
            <div class="resource-label">{{ $t('experience.gold_coins') }}</div>
            <div class="resource-amount">{{ character.vermoegen?.goldstücke || 0 }} GS</div>
          </div>
        </div>
        <div class="form-row control-row">
          <div class="form-group">
            <input 
              v-model.number="goldAmount" 
              type="number" 
              class="form-control amount-input"
              placeholder="Anzahl GS"
              min="1"
            />
          </div>
          <div class="button-group">
            <button @click="addGold" class="btn btn-success" :disabled="!goldAmount || goldAmount <= 0 || isLoading">
              <span v-if="isLoading">⏳</span>
              <span v-else>+ Hinzufügen</span>
            </button>
            <button @click="removeGold" class="btn btn-danger" :disabled="!goldAmount || goldAmount <= 0 || isLoading">
              <span v-if="isLoading">⏳</span>
              <span v-else>- Entfernen</span>
            </button>
          </div>
        </div>
        
        <!-- Total wealth display -->
        <div class="resource-card wealth-total">
          <span class="resource-icon">💎</span>
          <div class="resource-info">
            <div class="resource-label">{{ $t('experience.total_in_gs') }}</div>
            <div class="resource-amount total-amount">{{ totalWealthInGS }} GS</div>
          </div>
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
/* ExperianceView spezifische Styles */

.experience-section, 
.wealth-section {
  margin-bottom: 30px;
}

.control-row {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #e9ecef;
  align-items: center;
  gap: 15px;
}

.button-group {
  display: flex;
  gap: 8px;
}

.amount-input {
  width: 120px;
  text-align: right;
  font-weight: bold;
}

.wealth-total {
  margin-top: 15px;
  border: 2px solid #007bff;
  background: #f0f8ff;
}

.total-amount {
  color: #007bff;
  font-size: 1.1em;
  font-weight: bold;
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
