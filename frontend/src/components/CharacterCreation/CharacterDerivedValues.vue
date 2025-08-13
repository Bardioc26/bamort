<template>
  <div class="derived-values-form character-creation-container">
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
                v-if="value.key === 'pa' || value.key === 'wk' || value.key === 'lp_max' || value.key === 'ap_max' || value.key === 'b_max'"
                type="button" 
                class="dice-btn" 
                @click="rollField(value.key)"
                :title="getDiceTooltip(value.key)"
                :disabled="isCalculating"
              >
                {{ isCalculating ? '⏳' : '🎲' }}
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
        <button type="button" @click="calculateAllStatic" class="calc-btn" :disabled="isCalculating">
          {{ isCalculating ? $t('characters.derivedValues.calculating') : $t('characters.derivedValues.recalculate') }}
        </button>
        <button type="submit" class="next-btn" :disabled="!isValid">
          {{ $t('characters.derivedValues.nextSkills') }} →
        </button>
      </div>
    </form>
  </div>
</template>

<script>
import API from '../../utils/api'
import { rollDie, rollDice } from '../../utils/randomUtils'

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
        pa: 0, // Persönliche Ausstrahlung
        wk: 0, // Willenskraft
        lp_max: 0,
        ap_max: 0,
        b_max: 0,
        resistenz_koerper: 0, // Resistenz Körper
        resistenz_geist: 0, // Resistenz Geist
        resistenz_bonus_koerper: 0, // Resistenz Bonus Körper
        resistenz_bonus_geist: 0, // Resistenz Bonus Geist
        abwehr: 0, // Abwehr
        abwehr_bonus: 0, // Abwehr Bonus
        ausdauer_bonus: 0, // Ausdauer Bonus
        angriffs_bonus: 0, // Angriffs Bonus
        zaubern: 0, // Zaubern
        zauber_bonus: 0, // Zauber Bonus
        raufen: 0, // Raufen
        schadens_bonus: 0, // Schadens Bonus
        sg: 0, // Schicksalsgunst
        gg: 0, // Göttliche Gnade
        gp: 0, // Glückspunkte
      },
      isCalculating: false,
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
          max: 50
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
          max: 50
        },
        {
          key: 'resistenz_koerper',
          name: 'characters.derivedValues.resistenzKoerper',
          description: 'characters.derivedValues.resistenzKoerperDescription',
          min: 1,
          max: 20
        },
        {
          key: 'resistenz_geist',
          name: 'characters.derivedValues.resistenzGeist',
          description: 'characters.derivedValues.resistenzGeistDescription',
          min: 1,
          max: 20
        },
        {
          key: 'resistenz_bonus_koerper',
          name: 'characters.derivedValues.resistenzBonusKoerper',
          description: 'characters.derivedValues.resistenzBonusKoerperDescription',
          min: -5,
          max: 5
        },
        {
          key: 'resistenz_bonus_geist',
          name: 'characters.derivedValues.resistenzBonusGeist',
          description: 'characters.derivedValues.resistenzBonusGeistDescription',
          min: -5,
          max: 5
        },
        {
          key: 'abwehr',
          name: 'characters.derivedValues.abwehr',
          description: 'characters.derivedValues.abwehrDescription',
          min: 1,
          max: 20
        },
        {
          key: 'abwehr_bonus',
          name: 'characters.derivedValues.abwehrBonus',
          description: 'characters.derivedValues.abwehrBonusDescription',
          min: -5,
          max: 5
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
          min: -5,
          max: 5
        },
        {
          key: 'zaubern',
          name: 'characters.derivedValues.zaubern',
          description: 'characters.derivedValues.zaubernDescription',
          min: 1,
          max: 20
        },
        {
          key: 'zauber_bonus',
          name: 'characters.derivedValues.zauberBonus',
          description: 'characters.derivedValues.zauberBonusDescription',
          min: -5,
          max: 5
        },
        {
          key: 'raufen',
          name: 'characters.derivedValues.raufen',
          description: 'characters.derivedValues.raufenDescription',
          min: 1,
          max: 20
        },
        {
          key: 'schadens_bonus',
          name: 'characters.derivedValues.schadensBonus',
          description: 'characters.derivedValues.schadensBonusDescription',
          min: -10,
          max: 10
        },
        {
          key: 'sg',
          name: 'characters.derivedValues.sg',
          description: 'characters.derivedValues.sgDescription',
          min: 0,
          max: 50
        },
        {
          key: 'gg',
          name: 'characters.derivedValues.gg',
          description: 'characters.derivedValues.ggDescription',
          min: 0,
          max: 50
        },
        {
          key: 'gp',
          name: 'characters.derivedValues.gp',
          description: 'characters.derivedValues.gpDescription',
          min: 0,
          max: 50
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
      // Return currently loaded values or defaults
      // The actual calculation happens via API calls
      return this.formData
    }
  },
  watch: {
    formData: {
      handler(newValue) {
        // Save changes automatically when form data changes
        this.$emit('save', { derived_values: newValue })
      },
      deep: true
    }
  },
  created() {
    // Initialize with existing session data if available
    if (this.sessionData.derived_values && Object.keys(this.sessionData.derived_values).length > 0) {
      this.formData = { ...this.formData, ...this.sessionData.derived_values }
    } else {
      // Calculate initial values using new API
      this.calculateAllStatic()
    }
  },
  methods: {
    async calculateAllStatic() {
      if (this.isCalculating) return
      
      this.isCalculating = true
      try {
        const attrs = this.sessionData.attributes || {}
        const basic = this.sessionData.basic_info || {}
        const token = localStorage.getItem('token')
        
        const response = await API.post('/api/characters/calculate-static-fields', {
          st: attrs.st || 0,
          gs: attrs.gs || 0, 
          gw: attrs.gw || 0,
          ko: attrs.ko || 0,
          in: attrs.in || 0,
          zt: attrs.zt || 0,
          au: attrs.au || 0,
          rasse: basic.rasse || 'Menschen',
          typ: basic.typ || 'Barbar'
        }, {
          headers: { Authorization: `Bearer ${token}` }
        })
        
        const staticValues = response.data
        
        // Update form data with calculated static values
        this.formData = {
          ...this.formData,
          ausdauer_bonus: staticValues.ausdauer_bonus,
          schadens_bonus: staticValues.schadens_bonus,
          angriffs_bonus: staticValues.angriffs_bonus,
          abwehr_bonus: staticValues.abwehr_bonus,
          zauber_bonus: staticValues.zauber_bonus,
          resistenz_bonus_koerper: staticValues.resistenz_bonus_koerper,
          resistenz_bonus_geist: staticValues.resistenz_bonus_geist,
          resistenz_koerper: staticValues.resistenz_koerper,
          resistenz_geist: staticValues.resistenz_geist,
          abwehr: staticValues.abwehr,
          zaubern: staticValues.zaubern,
          raufen: staticValues.raufen
        }
        
        // Save the updated values to session
        this.$emit('save', { derived_values: this.formData })
      } catch (error) {
        console.error('Error calculating static values:', error)
      } finally {
        this.isCalculating = false
      }
    },
    
    async rollField(fieldName) {
      if (this.isCalculating) return
      
      this.isCalculating = true
      try {
        const attrs = this.sessionData.attributes || {}
        const basic = this.sessionData.basic_info || {}
        const token = localStorage.getItem('token')
        
        // Generate dice roll based on field type
        let roll
        switch (fieldName) {
          case 'pa':
          case 'wk':
            roll = rollDie(100) // 1d100
            break
          case 'lp_max':
            roll = rollDie(3) // 1d3 - single number
            break
          case 'ap_max':
            roll = rollDie(3) // 1d3 - array of 3 values
            break
          case 'b_max':
            // B Max depends on race: Gnome/Halblinge=2d3, Zwerge=3d3, others=4d3
            let diceCount = 4 // default for most races
            if (basic.rasse === 'Gnome' || basic.rasse === 'Halblinge') {
              diceCount = 2
            } else if (basic.rasse === 'Zwerge') {
              diceCount = 3
            }
            roll = rollDice(diceCount, 3) // XdY where X depends on race, Y=3
            break
        }
        
        const response = await API.post('/api/characters/calculate-rolled-field', {
          st: attrs.st || 0,
          gs: attrs.gs || 0,
          gw: attrs.gw || 0,
          ko: attrs.ko || 0,
          in: attrs.in || 0,
          zt: attrs.zt || 0,
          au: attrs.au || 0,
          rasse: basic.rasse || 'Menschen',
          typ: basic.typ || 'Barbar',
          field: fieldName,
          roll: roll
        }, {
          headers: { Authorization: `Bearer ${token}` }
        })
        
        const result = response.data
        this.formData[fieldName] = result.value
        
        // Save the updated values to session
        this.$emit('save', { derived_values: this.formData })
      } catch (error) {
        console.error('Error calculating rolled field:', error)
      } finally {
        this.isCalculating = false
      }
    },
    
    getDiceTooltip(fieldName) {
      switch (fieldName) {
        case 'pa': 
          return this.$t('characters.derivedValues.paRollTooltip')
        case 'wk':
          return this.$t('characters.derivedValues.wkRollTooltip')
        case 'lp_max':
          return this.$t('characters.derivedValues.lpRollTooltip')
        case 'ap_max':
          return this.$t('characters.derivedValues.apRollTooltip')
        case 'b_max':
          return this.$t('characters.derivedValues.bRollTooltip')
        default:
          return ''
      }
    },
    
    // Legacy methods for backward compatibility
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
      this.rollField('lp_max')
    },
    
    rollPA() {
      this.rollField('pa')
    },
    
    rollWK() {
      this.rollField('wk')
    },
    
    getClassBonnie(type) {
      // TODO: Implement class-specific bonnie calculations
      // For now, return base values
      const bonnieMap = {
        'sg': 0,
        'gg': 0,  
        'gp': 0,
      }
      return bonnieMap[type] || 1
    },
    
    recalculate() {
      this.calculateAllStatic()
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
