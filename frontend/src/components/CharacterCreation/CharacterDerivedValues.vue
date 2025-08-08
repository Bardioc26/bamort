<template>
  <div class="derived-values-form">
    <h2>Derived Values</h2>
    <p class="instruction">These values are calculated from your attributes. You can adjust them as needed.</p>
    
    <form @submit.prevent="handleSubmit">
      <div class="values-grid">
        <div class="value-group" v-for="value in derivedValues" :key="value.key">
          <label :for="value.key">{{ value.name }}</label>
          <div class="value-input-group">
            <input 
              :id="value.key"
              v-model.number="formData[value.key]"
              type="number"
              :min="value.min"
              :max="value.max"
              required
            />
            <div class="value-info">
              <span class="calculated-value">Calculated: {{ calculatedValues[value.key] }}</span>
              <span class="value-description">{{ value.description }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="calculation-info">
        <h3>Calculation Rules</h3>
        <div class="calculation-rules">
          <div class="rule">
            <strong>LP (Life Points):</strong> Base formula: (KO + ST) / 2 + modifier
          </div>
          <div class="rule">
            <strong>AP (Adventure Points):</strong> Base formula: (AU + WK) / 2 + modifier
          </div>
          <div class="rule">
            <strong>B (Burden):</strong> Base formula: ST + modifier
          </div>
          <div class="rule">
            <strong>Bennies:</strong> Base values based on character class
          </div>
        </div>
      </div>

      <div class="form-actions">
        <button type="button" @click="handlePrevious" class="prev-btn">
          ← Previous: Attributes
        </button>
        <button type="button" @click="recalculate" class="calc-btn">
          Recalculate from Attributes
        </button>
        <button type="submit" class="next-btn" :disabled="!isValid">
          Next: Skills & Spells →
        </button>
      </div>
    </form>
  </div>
</template>

<script>
export default {
  name: 'CharacterDerivedValues',
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
        lp_max: 20,
        ap_max: 20,
        b_max: 50,
        sg: 1, // Schicksalsgunst
        gg: 1, // Göttliche Gnade
        gp: 1, // Glückspunkte
      },
      derivedValues: [
        {
          key: 'lp_max',
          name: 'Life Points (LP) Maximum',
          description: 'Maximum life/health points',
          min: 1,
          max: 200
        },
        {
          key: 'ap_max',
          name: 'Adventure Points (AP) Maximum',
          description: 'Maximum adventure points for special actions',
          min: 1,
          max: 200
        },
        {
          key: 'b_max',
          name: 'Burden (B) Maximum',
          description: 'Maximum carrying capacity',
          min: 1,
          max: 500
        },
        {
          key: 'sg',
          name: 'Schicksalsgunst (SG)',
          description: 'Fate points for rerolls',
          min: 0,
          max: 10
        },
        {
          key: 'gg',
          name: 'Göttliche Gnade (GG)',
          description: 'Divine grace points',
          min: 0,
          max: 10
        },
        {
          key: 'gp',
          name: 'Glückspunkte (GP)',
          description: 'Luck points',
          min: 0,
          max: 10
        },
      ],
    }
  },
  computed: {
    isValid() {
      return Object.entries(this.formData).every(([key, val]) => {
        const valueConfig = this.derivedValues.find(v => v.key === key)
        return val >= valueConfig.min && val <= valueConfig.max
      })
    },
    
    calculatedValues() {
      const attrs = this.sessionData.attributes || {}
      
      return {
        lp_max: Math.floor(((attrs.ko || 50) + (attrs.st || 50)) / 2) + 5,
        ap_max: Math.floor(((attrs.au || 50) + (attrs.wk || 50)) / 2) + 5,
        b_max: (attrs.st || 50) + 10,
        sg: this.getClassBonnie('sg'),
        gg: this.getClassBonnie('gg'),
        gp: this.getClassBonnie('gp'),
      }
    }
  },
  created() {
    // Initialize with calculated values first
    this.formData = { ...this.calculatedValues }
    
    // Then override with session data if available
    if (this.sessionData.derived_values && Object.keys(this.sessionData.derived_values).length > 0) {
      this.formData = { ...this.formData, ...this.sessionData.derived_values }
    }
  },
  methods: {
    getClassBonnie(type) {
      // TODO: Implement class-specific bonnie calculations
      // For now, return base values
      const bonnieMap = {
        'sg': 1,
        'gg': 1,  
        'gp': 1,
      }
      return bonnieMap[type] || 1
    },
    
    recalculate() {
      this.formData = { ...this.calculatedValues }
    },
    
    handlePrevious() {
      this.$emit('previous')
    },
    
    handleSubmit() {
      if (this.isValid) {
        this.$emit('next', { derived_values: this.formData })
      }
    },
  }
}
</script>

<style scoped>
.derived-values-form {
  max-width: 800px;
  margin: 0 auto;
}

.derived-values-form h2 {
  text-align: center;
  margin-bottom: 10px;
  color: #333;
}

.instruction {
  text-align: center;
  margin-bottom: 30px;
  color: #666;
  font-style: italic;
}

.values-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.value-group {
  padding: 15px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background-color: #fafafa;
}

.value-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: bold;
  color: #333;
}

.value-input-group input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
  margin-bottom: 8px;
}

.value-input-group input:focus {
  outline: none;
  border-color: #2196f3;
  box-shadow: 0 0 5px rgba(33, 150, 243, 0.3);
}

.value-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.calculated-value {
  font-size: 12px;
  color: #4caf50;
  font-weight: bold;
}

.value-description {
  font-size: 12px;
  color: #666;
  font-style: italic;
}

.calculation-info {
  margin-bottom: 30px;
  padding: 20px;
  background-color: #e8f5e8;
  border-radius: 8px;
}

.calculation-info h3 {
  margin-bottom: 15px;
  color: #2e7d32;
}

.calculation-rules {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.rule {
  font-size: 14px;
  color: #555;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 15px;
}

.prev-btn, .calc-btn, .next-btn {
  padding: 12px 20px;
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

.calc-btn {
  background-color: #ff9800;
  color: white;
}

.calc-btn:hover {
  background-color: #f57c00;
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
</style>
