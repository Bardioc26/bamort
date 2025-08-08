<template>
  <div class="attributes-form">
    <h2>Character Attributes</h2>
    <p class="instruction">Set the basic attributes for your character (1-100)</p>
    
    <form @submit.prevent="handleSubmit" class="attributes-form-content">
      <div class="attributes-grid">
        <div class="attribute-group" v-for="attr in attributes" :key="attr.key">
          <div class="attribute-row">
            <label :for="attr.key" class="attribute-label">
              {{ attr.name }} ({{ attr.key.toUpperCase() }})
            </label>
            <input 
              :id="attr.key"
              v-model.number="formData[attr.key]"
              type="number"
              min="1"
              max="100"
              required
              class="attribute-input"
              @input="updateTotal"
            />
          </div>
          <span class="attribute-description">{{ attr.description }}</span>
        </div>
      </div>

      <div class="attribute-summary">
        <div class="total-points">
          <strong>Total Points: {{ totalPoints }}</strong>
        </div>
        <div class="average-points">
          <strong>Average: {{ averagePoints.toFixed(1) }}</strong>
        </div>
      </div>

      <div class="form-actions">
        <button type="button" @click="handlePrevious" class="prev-btn">
          ← Previous: Basic Info
        </button>
        <button type="submit" class="next-btn" :disabled="!isValid">
          Next: Derived Values →
        </button>
      </div>
    </form>
  </div>
</template>

<script>
export default {
  name: 'CharacterAttributes',
  props: {
    sessionData: {
      type: Object,
      required: true,
    }
  },
  emits: ['next', 'previous', 'save'],
  data() {
    return {
      formData: {
        st: 50, // Stärke
        gs: 50, // Geschicklichkeit
        gw: 50, // Gewandtheit
        ko: 50, // Konstitution
        in: 50, // Intelligenz
        zt: 50, // Zaubertalent
        au: 50, // Ausstrahlung
        pa: 50, // Psi-Kraft
        wk: 50, // Willenskraft
      },
      attributes: [
        {
          key: 'st',
          name: 'Stärke',
          description: 'Physical strength and power'
        },
        {
          key: 'gs',
          name: 'Geschicklichkeit',
          description: 'Dexterity and manual skill'
        },
        {
          key: 'gw',
          name: 'Gewandtheit',
          description: 'Agility and quick reactions'
        },
        {
          key: 'ko',
          name: 'Konstitution',
          description: 'Health and endurance'
        },
        {
          key: 'in',
          name: 'Intelligenz',
          description: 'Learning ability and logic'
        },
        {
          key: 'zt',
          name: 'Zaubertalent',
          description: 'Magical talent and mana'
        },
        {
          key: 'au',
          name: 'Ausstrahlung',
          description: 'Charisma and leadership'
        },
        {
          key: 'pa',
          name: 'Psi-Kraft',
          description: 'Psychic abilities'
        },
        {
          key: 'wk',
          name: 'Willenskraft',
          description: 'Mental fortitude and resistance'
        },
      ],
      totalPoints: 0,
    }
  },
  computed: {
    isValid() {
      return Object.values(this.formData).every(val => val >= 1 && val <= 100)
    },
    averagePoints() {
      return this.totalPoints / Object.keys(this.formData).length
    }
  },
  created() {
    // Initialize form with session data
    if (this.sessionData.attributes && Object.keys(this.sessionData.attributes).length > 0) {
      this.formData = { ...this.formData, ...this.sessionData.attributes }
    }
    this.updateTotal()
  },
  methods: {
    updateTotal() {
      this.totalPoints = Object.values(this.formData).reduce((sum, val) => sum + (val || 0), 0)
    },
    
    handlePrevious() {
      this.$emit('previous')
    },
    
    handleSubmit() {
      if (this.isValid) {
        this.$emit('next', { attributes: this.formData })
      }
    },
  }
}
</script>

<style scoped>
.attributes-form {
  max-width: 800px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  min-height: 0;
  padding-bottom: 20px;
}

.attributes-form h2 {
  text-align: center;
  margin-bottom: 10px;
  color: #333;
  flex-shrink: 0;
}

.instruction {
  text-align: center;
  margin-bottom: 20px;
  color: #666;
  font-style: italic;
  flex-shrink: 0;
}

.attributes-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
  margin-bottom: 20px;
  max-height: 50vh;
  overflow-y: auto;
  padding: 5px;
  border: 1px solid #eee;
  border-radius: 8px;
  background-color: #fefefe;
}

.attribute-group {
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background-color: #fafafa;
  min-width: 0; /* Prevent overflow */
}

.attribute-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 5px;
}

.attribute-label {
  font-weight: bold;
  color: #333;
  flex: 1;
  margin: 0;
}

.attribute-input {
  width: 60px;
  padding: 6px 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
  text-align: center;
  font-weight: bold;
}

.attribute-input:focus {
  outline: none;
  border-color: #2196f3;
  box-shadow: 0 0 5px rgba(33, 150, 243, 0.3);
}

.attribute-description {
  font-size: 11px;
  color: #666;
  font-style: italic;
  display: block;
  margin-top: 2px;
}

.attribute-summary {
  display: flex;
  justify-content: center;
  gap: 30px;
  margin-bottom: 20px;
  padding: 15px;
  background-color: #e3f2fd;
  border-radius: 8px;
  flex-shrink: 0;
}

.total-points, .average-points {
  font-size: 18px;
  color: #1976d2;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;
  margin-top: 20px;
  padding-top: 15px;
  border-top: 1px solid #eee;
}

.attributes-form-content {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.prev-btn, .next-btn {
  padding: 12px 30px;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.prev-btn {
  background-color: #6c757d;
  color: white;
}

.prev-btn:hover {
  background-color: #5a6268;
}

.next-btn {
  background-color: #2196f3;
  color: white;
}

.next-btn:hover:not(:disabled) {
  background-color: #1976d2;
}

.next-btn:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

/* Responsive Design für sehr kleine Bildschirme */
@media (max-width: 600px) {
  .attributes-grid {
    grid-template-columns: 1fr;
  }
  
  .attribute-group {
    padding: 10px;
  }
}
</style>
