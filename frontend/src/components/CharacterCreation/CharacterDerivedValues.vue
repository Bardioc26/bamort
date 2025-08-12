<template>
  <div class="derived-values-form">
    <h2>{{ $t('characters.derivedValues.title') }}</h2>
    <p class="instruction">{{ $t('characters.derivedValues.instruction') }}</p>
    
    <form @submit.prevent="handleSubmit">
      <div class="values-grid">
        <div class="value-group" v-for="value in derivedValues" :key="value.key">
          <label :for="value.key">{{ $t(value.name) }}</label>
          <div class="value-input-group">
            <div class="input-with-dice">
              <input 
                :id="value.key"
                v-model.number="formData[value.key]"
                type="number"
                :min="value.min"
                :max="value.max"
                required
              />
              <button 
                v-if="value.key === 'pa' || value.key === 'wk' || value.key === 'lp_max'"
                type="button" 
                class="dice-btn" 
                @click="value.key === 'pa' ? rollPA() : 
                       value.key === 'wk' ? rollWK() : 
                       value.key === 'lp_max' ? rollLP() : null"
                :title="value.key === 'pa' ? 'Roll PA: 1d100 + 4×(In/10) - 20' : 
                       value.key === 'wk' ? 'Roll WK: 1d100 + 2×(Ko/10 + In/10) - 20' : 
                       value.key === 'lp_max' ? 'Roll LP: 1d3 + 7 + (Ko/10)' : ''"
              >
                🎲
              </button>
            </div>
            <div class="value-info">
              <span class="calculated-value">{{ $t('characters.derivedValues.calculated') }}: {{ calculatedValues[value.key] }}</span>
              <span class="value-description">{{ $t(value.description) }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="calculation-info">
        <h3>{{ $t('characters.derivedValues.calculationRules') }}</h3>
        <div class="calculation-rules">
          <div class="rule">
            <strong>{{ $t('characters.derivedValues.lpFormula') }}:</strong> {{ $t('characters.derivedValues.lpDescription') }}
          </div>
          <div class="rule">
            <strong>{{ $t('characters.derivedValues.apFormula') }}:</strong> {{ $t('characters.derivedValues.apDescription') }}
          </div>
          <div class="rule">
            <strong>{{ $t('characters.derivedValues.bFormula') }}:</strong> {{ $t('characters.derivedValues.bDescription') }}
          </div>
          <div class="rule">
            <strong>{{ $t('characters.derivedValues.benniesFormula') }}:</strong> {{ $t('characters.derivedValues.benniesDescription') }}
          </div>
        </div>
      </div>

      <div class="form-actions">
        <button type="button" @click="handlePrevious" class="prev-btn">
          ← {{ $t('characters.derivedValues.previousAttributes') }}
        </button>
        <button type="button" @click="recalculate" class="calc-btn">
          {{ $t('characters.derivedValues.recalculate') }}
        </button>
        <button type="submit" class="next-btn" :disabled="!isValid">
          {{ $t('characters.derivedValues.nextSkills') }} →
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
        pa: 50, // Persönliche Ausstrahlung
        wk: 50, // Willenskraft
        lp_max: 20,
        ap_max: 20,
        b_max: 50,
        resistenz_koerper: 50, // Resistenz Körper
        resistenz_geist: 50, // Resistenz Geist
        resistenz_bonus_koerper: 0, // Resistenz Bonus Körper
        resistenz_bonus_geist: 0, // Resistenz Bonus Geist
        abwehr: 50, // Abwehr
        abwehr_bonus: 0, // Abwehr Bonus
        ausdauer_bonus: 0, // Ausdauer Bonus
        angriffs_bonus: 0, // Angriffs Bonus
        zaubern: 50, // Zaubern
        zauber_bonus: 0, // Zauber Bonus
        raufen: 50, // Raufen
        schadens_bonus: 0, // Schadens Bonus
        sg: 1, // Schicksalsgunst
        gg: 1, // Göttliche Gnade
        gp: 1, // Glückspunkte
      },
      derivedValues: [
        {
          key: 'pa',
          name: 'characters.derivedValues.pa',
          description: 'characters.derivedValues.paDescription',
          min: 1,
          max: 100
        },
        {
          key: 'wk',
          name: 'characters.derivedValues.wk',
          description: 'characters.derivedValues.wkDescription',
          min: 1,
          max: 100
        },
        {
          key: 'lp_max',
          name: 'characters.derivedValues.lpMax',
          description: 'characters.derivedValues.lpMaxDescription',
          min: 1,
          max: 200
        },
        {
          key: 'ap_max',
          name: 'characters.derivedValues.apMax',
          description: 'characters.derivedValues.apMaxDescription',
          min: 1,
          max: 200
        },
        {
          key: 'b_max',
          name: 'characters.derivedValues.bMax',
          description: 'characters.derivedValues.bMaxDescription',
          min: 1,
          max: 500
        },
        {
          key: 'resistenz_koerper',
          name: 'characters.derivedValues.resistenzKoerper',
          description: 'characters.derivedValues.resistenzKoerperDescription',
          min: 1,
          max: 100
        },
        {
          key: 'resistenz_geist',
          name: 'characters.derivedValues.resistenzGeist',
          description: 'characters.derivedValues.resistenzGeistDescription',
          min: 1,
          max: 100
        },
        {
          key: 'resistenz_bonus_koerper',
          name: 'characters.derivedValues.resistenzBonusKoerper',
          description: 'characters.derivedValues.resistenzBonusKoerperDescription',
          min: -50,
          max: 50
        },
        {
          key: 'resistenz_bonus_geist',
          name: 'characters.derivedValues.resistenzBonusGeist',
          description: 'characters.derivedValues.resistenzBonusGeistDescription',
          min: -50,
          max: 50
        },
        {
          key: 'abwehr',
          name: 'characters.derivedValues.abwehr',
          description: 'characters.derivedValues.abwehrDescription',
          min: 1,
          max: 100
        },
        {
          key: 'abwehr_bonus',
          name: 'characters.derivedValues.abwehrBonus',
          description: 'characters.derivedValues.abwehrBonusDescription',
          min: -50,
          max: 50
        },
        {
          key: 'ausdauer_bonus',
          name: 'characters.derivedValues.ausdauerBonus',
          description: 'characters.derivedValues.ausdauerBonusDescription',
          min: -50,
          max: 50
        },
        {
          key: 'angriffs_bonus',
          name: 'characters.derivedValues.angriffsBonus',
          description: 'characters.derivedValues.angriffsBonusDescription',
          min: -50,
          max: 50
        },
        {
          key: 'zaubern',
          name: 'characters.derivedValues.zaubern',
          description: 'characters.derivedValues.zaubernDescription',
          min: 1,
          max: 100
        },
        {
          key: 'zauber_bonus',
          name: 'characters.derivedValues.zauberBonus',
          description: 'characters.derivedValues.zauberBonusDescription',
          min: -50,
          max: 50
        },
        {
          key: 'raufen',
          name: 'characters.derivedValues.raufen',
          description: 'characters.derivedValues.raufenDescription',
          min: 1,
          max: 100
        },
        {
          key: 'schadens_bonus',
          name: 'characters.derivedValues.schadensBonus',
          description: 'characters.derivedValues.schadensBonusDescription',
          min: -50,
          max: 50
        },
        {
          key: 'sg',
          name: 'characters.derivedValues.sg',
          description: 'characters.derivedValues.sgDescription',
          min: 0,
          max: 10
        },
        {
          key: 'gg',
          name: 'characters.derivedValues.gg',
          description: 'characters.derivedValues.ggDescription',
          min: 0,
          max: 10
        },
        {
          key: 'gp',
          name: 'characters.derivedValues.gp',
          description: 'characters.derivedValues.gpDescription',
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
      const wkValue = this.formData?.wk || this.calculateWK(attrs.ko || 50, attrs.in || 50)
      
      return {
        pa: this.calculatePA(attrs.in || 50),
        wk: this.calculateWK(attrs.ko || 50, attrs.in || 50),
        lp_max: this.calculateLP(attrs.ko || 50),
        ap_max: Math.floor(((attrs.au || 50) + wkValue) / 2) + 5,
        b_max: (attrs.st || 50) + 10,
        resistenz_koerper: attrs.ko || 50, // Resistenz Körper = Konstitution
        resistenz_geist: wkValue, // Resistenz Geist = Willenskraft
        resistenz_bonus_koerper: Math.floor((attrs.ko || 50) / 20) - 2, // Bonus basierend auf Ko
        resistenz_bonus_geist: Math.floor(wkValue / 20) - 2, // Bonus basierend auf WK
        abwehr: Math.floor(((attrs.gw || 50) + (attrs.gs || 50)) / 2), // Abwehr = (Gewandtheit + Geschicklichkeit) / 2
        abwehr_bonus: Math.floor(((attrs.gw || 50) + (attrs.gs || 50)) / 40) - 2, // Abwehr Bonus
        ausdauer_bonus: Math.floor((attrs.ko || 50) / 20) - 2, // Ausdauer Bonus basierend auf Ko
        angriffs_bonus: Math.floor((attrs.gs || 50) / 20) - 2, // Angriffs Bonus basierend auf Gs
        zaubern: attrs.zt || 50, // Zaubern = Zaubertalent
        zauber_bonus: Math.floor((attrs.zt || 50) / 20) - 2, // Zauber Bonus
        raufen: Math.floor(((attrs.st || 50) + (attrs.gw || 50)) / 2), // Raufen = (Stärke + Gewandtheit) / 2
        schadens_bonus: Math.floor((attrs.st || 50) / 20) - 2, // Schadens Bonus basierend auf St
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
    calculateLP(constitution) {
      // LP = 1d3 + 7 + (Ko/10)
      const diceRoll = Math.floor(Math.random() * 3) + 1 // 1d3
      const constitutionBonus = Math.floor(constitution / 10)
      const result = diceRoll + 7 + constitutionBonus
      return result
    },
    
    calculatePA(intelligence) {
      // PA = 1d100 + 4×(In/10) - 20
      const baseRoll = Math.floor(Math.random() * 100) + 1
      const intelligenceBonus = Math.floor(intelligence / 10) * 4
      const result = baseRoll + intelligenceBonus - 20
      return Math.max(1, Math.min(100, result))
    },
    
    calculateWK(constitution, intelligence) {
      // WK = 1d100 + 2×(Ko/10 + In/10) - 20
      const baseRoll = Math.floor(Math.random() * 100) + 1
      const constitutionBonus = Math.floor(constitution / 10)
      const intelligenceBonus = Math.floor(intelligence / 10)
      const combinedBonus = (constitutionBonus + intelligenceBonus) * 2
      const result = baseRoll + combinedBonus - 20
      return Math.max(1, Math.min(100, result))
    },
    
    rollLP() {
      const attrs = this.sessionData.attributes || {}
      this.formData.lp_max = this.calculateLP(attrs.ko || 50)
    },
    
    rollPA() {
      const attrs = this.sessionData.attributes || {}
      this.formData.pa = this.calculatePA(attrs.in || 50)
    },
    
    rollWK() {
      const attrs = this.sessionData.attributes || {}
      this.formData.wk = this.calculateWK(attrs.ko || 50, attrs.in || 50)
    },
    
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
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 15px;
  margin-bottom: 30px;
  max-height: 60vh;
  overflow-y: auto;
  padding: 10px;
  border: 1px solid #eee;
  border-radius: 8px;
  background-color: #fefefe;
}

.value-group {
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background-color: #fafafa;
  min-width: 0;
}

.value-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: bold;
  color: #333;
}

.input-with-dice {
  display: flex;
  gap: 8px;
  align-items: center;
}

.value-input-group input {
  flex: 1;
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
