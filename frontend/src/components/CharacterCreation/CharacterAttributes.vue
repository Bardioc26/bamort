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
            <div class="input-with-dice">
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
              <button 
                type="button" 
                class="dice-btn" 
                @click="rollAttribute(attr.key)"
                :title="attr.key === 'pa' ? 'Roll PA: 1d100 + 4×(In/10) - 20' : 
                       attr.key === 'wk' ? 'Roll WK: 1d100 + 2×(Ko/10 + In/10) - 20' : 
                       attr.key === 'au' ? 'Roll AU with race restrictions (Elfen ≥81, Gnome/Zwerge ≤80)' :
                       'Roll max(2d100) for ' + attr.name"
              >
                🎲
              </button>
            </div>
          </div>
          <span class="attribute-description">{{ attr.description }}</span>
          <!-- Race restriction warning for AU -->
          <div v-if="attr.key === 'au' && auRaceRestriction" class="race-restriction-warning">
            ⚠️ {{ auRaceRestriction.message }}
          </div>
          <div v-if="lastAttributeRoll && lastAttributeRoll.attribute === attr.key" class="roll-result">
            {{ attr.name }}: {{ lastAttributeRoll.roll }} 
            <span class="roll-breakdown">
              <span v-if="lastAttributeRoll.isSpecialCalculation">
                ({{ lastAttributeRoll.description }})
              </span>
              <span v-else>
                (max of {{ lastAttributeRoll.rolls.join(', ') }})
              </span>
            </span>
            → {{ lastAttributeRoll.result }}
          </div>
        </div>
      </div>

      <div class="attribute-summary">
        <div class="total-points">
          <strong>Total Points: {{ totalPoints }}</strong>
        </div>
        <div class="average-points">
          <strong>Average: {{ averagePoints.toFixed(1) }}</strong>
        </div>
        <button type="button" @click="rollAllAttributes" class="roll-all-btn">
          🎲 Roll All Attributes
        </button>
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
    
    <!-- Roll Result Overlay -->
    <div v-if="showOverlay && lastAttributeRoll" class="roll-overlay" @click="hideOverlay">
      <div class="roll-overlay-content">
        <button class="overlay-close" @click="hideOverlay">×</button>
        <div class="overlay-title">🎲 {{ lastAttributeRoll.attributeName }}</div>
        <div class="overlay-roll">
          {{ lastAttributeRoll.result }}
          <span class="roll-breakdown">
            <span v-if="lastAttributeRoll.isSpecialCalculation">
              ({{ lastAttributeRoll.description }})
            </span>
            <span v-else>
              (max of {{ lastAttributeRoll.rolls.join(', ') }})
            </span>
          </span>
        </div>
        <div class="overlay-result">
          → {{ lastAttributeRoll.attributeName }}: {{ lastAttributeRoll.result }}
        </div>
        <div class="overlay-hint">Click to close</div>
      </div>
    </div>
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
        au: 50, // Ausehen
        pa: 50, // Persönliche Ausstrahlung
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
          name: 'Aussehen',
          description: 'Beauty and appearance (Race restrictions: Elfen ≥81, Gnome/Zwerge ≤80)'
        },
        {
          key: 'pa',
          name: 'Persönliche Ausstrahlung',
          description: 'Charisma and leadership (Special: 1d100 + 4×(In/10) - 20)'
        },
        {
          key: 'wk',
          name: 'Willenskraft',
          description: 'Mental fortitude and resistance (Special: 1d100 + 2×(Ko/10 + In/10) - 20)'
        },
      ],
      totalPoints: 0,
      lastAttributeRoll: null,
      showOverlay: false,
      overlayTimeout: null,
    }
  },
  computed: {
    isValid() {
      return Object.values(this.formData).every(val => val >= 1 && val <= 100)
    },
    averagePoints() {
      return this.totalPoints / Object.keys(this.formData).length
    },
    auRaceRestriction() {
      const race = this.sessionData.rasse || ''
      if (race === 'Elfen') {
        return { type: 'minimum', value: 81, message: 'Elfen must have AU ≥ 81' }
      } else if (race === 'Gnome') {
        return { type: 'maximum', value: 80, message: 'Gnome must have AU ≤ 80' }
      } else if (race === 'Zwerge') {
        return { type: 'maximum', value: 80, message: 'Zwerge must have AU ≤ 80' }
      }
      return null
    }
  },
  created() {
    // Initialize form with session data
    if (this.sessionData.attributes && Object.keys(this.sessionData.attributes).length > 0) {
      this.formData = { ...this.formData, ...this.sessionData.attributes }
    }
    this.updateTotal()
  },
  beforeUnmount() {
    // Clean up timeout
    if (this.overlayTimeout) {
      clearTimeout(this.overlayTimeout)
    }
  },
  methods: {
    updateTotal() {
      this.totalPoints = Object.values(this.formData).reduce((sum, val) => sum + (val || 0), 0)
    },
    
    rollAttribute(attributeKey) {
      let roll, rollValue, modifier = 0, rollDescription = ''
      
      if (attributeKey === 'pa') {
        // Special calculation for PA: PA = 1d100 + 4×(In/10) - 20
        const baseRoll = this.$rollNotation('1d100')
        const intelligenceBonus = Math.floor(this.formData.in / 10) * 4
        modifier = intelligenceBonus - 20
        rollValue = baseRoll.sum + modifier
        
        // Ensure PA stays within valid range (1-100)
        rollValue = Math.max(1, Math.min(100, rollValue))
        
        roll = {
          ...baseRoll,
          selectedValue: rollValue
        }
        
        rollDescription = `1d100 + 4×(${this.formData.in}/10) - 20 = ${baseRoll.sum} + ${intelligenceBonus} - 20`
      } else if (attributeKey === 'wk') {
        // Special calculation for WK: WK = 1d100 + 2×(Ko/10 + In/10) - 20
        const baseRoll = this.$rollNotation('1d100')
        const constitutionBonus = Math.floor(this.formData.ko / 10)
        const intelligenceBonus = Math.floor(this.formData.in / 10)
        const combinedBonus = (constitutionBonus + intelligenceBonus) * 2
        modifier = combinedBonus - 20
        rollValue = baseRoll.sum + modifier
        
        // Ensure WK stays within valid range (1-100)
        rollValue = Math.max(1, Math.min(100, rollValue))
        
        roll = {
          ...baseRoll,
          selectedValue: rollValue
        }
        
        rollDescription = `1d100 + 2×(${this.formData.ko}/10 + ${this.formData.in}/10) - 20 = ${baseRoll.sum} + 2×(${constitutionBonus} + ${intelligenceBonus}) - 20`
      } else if (attributeKey === 'au') {
        // Standard max(2d100) roll for AU with race-based restrictions
        roll = this.$rollNotation('max(2d100)')
        rollValue = roll.selectedValue
        
        // Apply race-based restrictions for AU (Aussehen)
        const race = this.sessionData.rasse || ''
        let minValue = 1, maxValue = 100, raceRestriction = ''
        
        if (race === 'Elf') {
          minValue = 81
          raceRestriction = ' (Elfen minimum: 81)'
        } else if (race === 'Gnom' || race === 'Zwerg') {
          maxValue = 80
          raceRestriction = ` (${race} maximum: 80)`
        }
        //console.log('rollvalue, minval, maxval, race:', rollValue ,minValue, maxValue, race)
        
        // Store original roll value for comparison
        const originalRollValue = rollValue
        
        // Apply race restrictions
        if (rollValue < minValue) {
          rollValue = minValue
        } else if (rollValue > maxValue) {
          rollValue = maxValue
        }
        //console.log('2 rollvalue, minval, maxval, race:', rollValue ,minValue, maxValue, race)
        roll = {
          ...roll,
          selectedValue: rollValue
        }
        
        rollDescription = `max of ${roll.rolls.join(', ')}${raceRestriction}`
        if (rollValue !== originalRollValue) {
          rollDescription += ` → adjusted to ${rollValue}`
        }
      } else {
        // Standard max(2d100) roll for other attributes
        roll = this.$rollNotation('max(2d100)')
        rollValue = roll.selectedValue
        rollDescription = `max of ${roll.rolls.join(', ')}`
        //console.log('Standard max(2d100) roll for other attributes:', rollValue)
      }
      
      const attributeName = this.attributes.find(attr => attr.key === attributeKey)?.name || attributeKey
      
      this.formData[attributeKey] = rollValue
      this.updateTotal()
      
      // Store roll information for display
      this.lastAttributeRoll = {
        attribute: attributeKey,
        attributeName: attributeName,
        rolls: roll.rolls || [roll.sum],
        roll: rollValue,
        result: rollValue,
        description: rollDescription,
        modifier: modifier,
        isSpecialCalculation: attributeKey === 'pa' || attributeKey === 'wk' || attributeKey === 'au'
      }
      
      // Show overlay notification
      this.showRollOverlay()
    },
    
    rollAllAttributes() {
      // Roll all attributes at once
      Object.keys(this.formData).forEach(key => {
        //const roll = this.$rollNotation('max(2d100)')
        //this.formData[key] = roll.selectedValue
        this.rollAttribute(key)
      })
      this.updateTotal()
      
      // Clear individual roll display when rolling all
      this.lastAttributeRoll = null
    },
    
    showRollOverlay() {
      this.showOverlay = true
      
      // Clear existing timeout if any
      if (this.overlayTimeout) {
        clearTimeout(this.overlayTimeout)
      }
      
      // Hide overlay after 10 seconds (shorter for attributes)
      this.overlayTimeout = setTimeout(() => {
        this.showOverlay = false
      }, 10000)
    },
    
    hideOverlay() {
      this.showOverlay = false
      if (this.overlayTimeout) {
        clearTimeout(this.overlayTimeout)
        this.overlayTimeout = null
      }
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

.input-with-dice {
  display: flex;
  gap: 8px;
  align-items: center;
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

.race-restriction-warning {
  font-size: 10px;
  color: #ff5722;
  font-weight: bold;
  margin-top: 2px;
  display: block;
}

.attribute-summary {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 30px;
  margin-bottom: 20px;
  padding: 15px;
  background-color: #e3f2fd;
  border-radius: 8px;
  flex-shrink: 0;
  flex-wrap: wrap;
}

.total-points, .average-points {
  font-size: 18px;
  color: #1976d2;
}

.roll-all-btn {
  background-color: #ff9800;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 10px 20px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.roll-all-btn:hover {
  background-color: #f57c00;
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
