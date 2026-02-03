<template>
  <div class="basic-info-form character-creation-container">
    <h2>{{ $t('characters.basicInfo.title') }}</h2>
    
    <form @submit.prevent="handleSubmit">
      <div class="form-row">
        <!-- 1. Name -->
        <div class="form-group">
          <label for="name">{{ $t('characters.basicInfo.characterName') }} {{ $t('characters.basicInfo.required') }}
            <span
            class="help-icon"
            :title="$t('characters.basicInfo.characterNameHelp')"
            role="img"
            :aria-label="$t('characters.basicInfo.characterNameHelp')"
          >
            ?
          </span>
          </label>
          <input 
            id="name"
            v-model="formData.name"
            type="text"
            required
            minlength="2"
            maxlength="50"
            :placeholder="$t('characters.basicInfo.characterNamePlaceholder')"
          />
        </div>

        <!-- 2. Herkunft -->
        <div class="form-group">
          <label for="herkunft">{{ $t('characters.basicInfo.origin') }} {{ $t('characters.basicInfo.required') }}</label>
          <select id="herkunft" v-model="formData.herkunft" required>
            <option value="">{{ $t('characters.basicInfo.selectOrigin') }}</option>
            <option v-for="origin in origins" :key="origin" :value="origin">{{ origin }}</option>
          </select>
        </div>
      </div>

      <!-- 3. Glaube -->
      <div class="form-group">
        <label for="glaube">
          {{ $t('characters.basicInfo.religion') }}
          <span
            class="help-icon"
            :title="$t('characters.basicInfo.religionHelp')"
            role="img"
            :aria-label="$t('characters.basicInfo.religionHelp')"
          >
            ?
          </span>
        </label>
        <div class="belief-search">
          <input 
            id="glaube"
            v-model="beliefSearch"
            type="text"
            :placeholder="$t('characters.basicInfo.religionPlaceholder')"
            @input="searchBeliefs"
          />
          <div v-if="beliefResults.length > 0" class="belief-dropdown">
            <div 
              v-for="belief in beliefResults" 
              :key="belief"
              class="belief-option"
              @click="selectBelief(belief)"
            >
              {{ belief }}
            </div>
          </div>
        </div>
        <div v-if="formData.glaube" class="selected-belief">
          {{ $t('characters.basicInfo.selected') }}: {{ formData.glaube }}
          <button type="button" @click="clearBelief" class="clear-btn">×</button>
        </div>
      </div>

      <div class="form-row">
        <!-- 4. Geschlecht -->
        <div class="form-group">
          <label for="geschlecht">{{ $t('characters.basicInfo.gender') }} {{ $t('characters.basicInfo.required') }}</label>
          <select id="geschlecht" v-model="formData.geschlecht" required>
            <option value="">{{ $t('characters.basicInfo.selectGender') }}</option>
            <option value="Männlich">{{ $t('characters.basicInfo.male') }}</option>
            <option value="Weiblich">{{ $t('characters.basicInfo.female') }}</option>
          </select>
        </div>

        <!-- 5. Rasse -->
        <div class="form-group">
          <label for="rasse">{{ $t('characters.basicInfo.race') }} {{ $t('characters.basicInfo.required') }}</label>
          <select id="rasse" v-model="formData.rasse" required>
            <option value="">{{ $t('characters.basicInfo.selectRace') }}</option>
            <option v-for="race in races" :key="race" :value="race">{{ race }}</option>
          </select>
        </div>
      </div>

      <div class="form-row">
        <!-- 6. Charakterklasse -->
        <div class="form-group">
          <label for="typ">{{ $t('characters.basicInfo.characterClass') }} {{ $t('characters.basicInfo.required') }}</label>
          <select id="typ" v-model="formData.typ" required>
            <option value="">{{ $t('characters.basicInfo.selectClass') }}</option>
            <option v-for="cls in classes" :key="cls" :value="cls">{{ cls }}</option>
          </select>
        </div>

        <!-- 7. Sozialschicht -->
        <div class="form-group">
          <label for="stand">{{ $t('characters.basicInfo.socialClass') }} {{ $t('characters.basicInfo.required') }}</label>
          <div class="input-with-dice">
            <select id="stand" v-model="formData.stand" required>
              <option value="">{{ $t('characters.basicInfo.selectSocialClass') }}</option>
              <option value="Adel">{{ $t('characters.basicInfo.nobility') }}</option>
              <option value="Mittelschicht">{{ $t('characters.basicInfo.middleClass') }}</option>
              <option value="Volk">{{ $t('characters.basicInfo.commonFolk') }}</option>
              <option value="Unfrei">{{ $t('characters.basicInfo.unfree') }}</option>
            </select>
            <button 
              type="button" 
              class="dice-btn" 
              @click="rollSocialClass"
              :disabled="!formData.typ"
              :title="formData.typ ? $t('characters.basicInfo.rollSocialClass') : $t('characters.basicInfo.selectClassFirst')"
            >
              🎲
            </button>
          </div>
          <div v-if="lastSocialClassRoll" class="roll-result">
            {{ $t('characters.basicInfo.rollResult') }}: {{ lastSocialClassRoll.roll }} 
            <span v-if="lastSocialClassRoll.modifier !== 0">
              ({{ lastSocialClassRoll.baseRoll }}{{ lastSocialClassRoll.modifier >= 0 ? '+' : '' }}{{ lastSocialClassRoll.modifier }})
            </span>
            → {{ lastSocialClassRoll.result }}
          </div>
        </div>
      </div>

      <div class="form-actions">
        <button type="submit" class="next-btn" :disabled="!isValid">
          {{ $t('characters.basicInfo.nextAttributes') }}
        </button>
      </div>
    </form>
    
    <!-- Roll Result Overlay -->
    <div v-if="showOverlay && lastSocialClassRoll" class="roll-overlay" @click="hideOverlay">
      <div class="roll-overlay-content">
        <button class="overlay-close" @click="hideOverlay">×</button>
        <div class="overlay-title">🎲 {{ $t('characters.basicInfo.rollResult') }}</div>
        <div class="overlay-roll">
          {{ lastSocialClassRoll.roll }}
          <span v-if="lastSocialClassRoll.modifier !== 0" class="roll-breakdown">
            ({{ lastSocialClassRoll.baseRoll }}{{ lastSocialClassRoll.modifier >= 0 ? '+' : '' }}{{ lastSocialClassRoll.modifier }})
          </span>
        </div>
        <div class="overlay-result">
          → {{ $t('characters.basicInfo.' + lastSocialClassRoll.result.toLowerCase()) || lastSocialClassRoll.result }}
        </div>
        <div class="overlay-hint">{{ $t('characters.basicInfo.clickToClose') }}</div>
      </div>
    </div>
  </div>
</template>

<script>
import API from '../../utils/api'
import { rollNotation } from '../../utils/randomUtils'

export default {
  name: 'CharacterBasicInfo',
  props: {
    sessionData: {
      type: Object,
      required: true,
    }
  },
  emits: ['next', 'save'],
  data() {
    return {
      formData: {
        name: '',
        geschlecht: '',
        rasse: '',
        typ: '',
        herkunft: '',
        stand: '',
        glaube: '',
      },
      races: [],
      classes: [],
      origins: [],
      beliefSearch: '',
      beliefResults: [],
      searchTimeout: null,
      lastSocialClassRoll: null,
      showOverlay: false,
      overlayTimeout: null,
      isInitialized: false, // Flag to prevent early watcher triggers
    }
  },
  computed: {
    isValid() {
      return this.formData.name.length >= 2 && 
             this.formData.geschlecht &&
             this.formData.rasse && 
             this.formData.typ && 
             this.formData.herkunft &&
             this.formData.stand
    }
  },
  watch: {
    'formData.typ'() {
      // Clear social class roll result when character class changes
      this.lastSocialClassRoll = null
    },
    formData: {
      handler(newValue) {
        // Only save if component is fully initialized
        if (this.isInitialized) {
          this.$emit('save', { basic_info: newValue })
        }
      },
      deep: true
    }
  },
  async created() {
    // Initialize form with session data - check both old format and new basic_info format
    const basicInfo = this.sessionData.basic_info || {}
    this.formData = {
      name: basicInfo.name || this.sessionData.name || '',
      geschlecht: basicInfo.geschlecht || this.sessionData.geschlecht || '',
      rasse: basicInfo.rasse || this.sessionData.rasse || '',
      typ: basicInfo.typ || this.sessionData.typ || '',
      herkunft: basicInfo.herkunft || this.sessionData.herkunft || '',
      stand: basicInfo.stand || this.sessionData.stand || '',
      glaube: basicInfo.glaube || this.sessionData.glaube || '',
    }
    
    if (this.formData.glaube) {
      this.beliefSearch = this.formData.glaube
    }
    
    // Save initial state to ensure all fields are captured
    this.$emit('save', { basic_info: this.formData })
    
    await this.loadReferenceData()
    
    // Mark as initialized to enable watcher
    this.isInitialized = true
  },
  beforeUnmount() {
    // Clean up timeouts
    if (this.searchTimeout) {
      clearTimeout(this.searchTimeout)
    }
    if (this.overlayTimeout) {
      clearTimeout(this.overlayTimeout)
    }
  },
  methods: {
    async loadReferenceData() {
      try {
        const token = localStorage.getItem('token')
        const headers = { Authorization: `Bearer ${token}` }
        
        // Load all reference data in parallel
        const [racesRes, classesRes, originsRes] = await Promise.all([
          API.get('/api/characters/races', { headers }),
          API.get('/api/characters/classes', { headers }),
          API.get('/api/characters/origins', { headers }),
        ])
        
        this.races = racesRes.data.races
        this.classes = classesRes.data.classes
        this.origins = originsRes.data.origins
      } catch (error) {
        console.error('Error loading reference data:', error)
      }
    },
    
    searchBeliefs() {
      if (this.searchTimeout) {
        clearTimeout(this.searchTimeout)
      }
      
      this.searchTimeout = setTimeout(async () => {
        if (this.beliefSearch.length >= 2) {
          try {
            const token = localStorage.getItem('token')
            const response = await API.get(`/api/characters/beliefs?q=${this.beliefSearch}`, {
              headers: { Authorization: `Bearer ${token}` },
            })
            this.beliefResults = response.data.beliefs
          } catch (error) {
            console.error('Error searching beliefs:', error)
            this.beliefResults = []
          }
        } else {
          this.beliefResults = []
        }
      }, 300)
    },
    
    selectBelief(belief) {
      this.formData.glaube = belief
      this.beliefSearch = belief
      this.beliefResults = []
    },
    
    clearBelief() {
      this.formData.glaube = ''
      this.beliefSearch = ''
      this.beliefResults = []
    },
    
    rollSocialClass() {
      if (!this.formData.typ) {
        return
      }
      
      // Base 1d100 roll
      const baseRoll = rollNotation('1d100')
      let modifier = 0
      
      // Apply class modifiers
      switch (this.formData.typ) {
        case 'Barde':
        case 'Priester':
          modifier = 20
          break
        case 'Druide':
        case 'Magier':
          modifier = 10
          break
        case 'Assassine':
        case 'Händler':
        case 'Waldläufer':
          modifier = -10
          break
        case 'Spitzbube':
          modifier = -20
          break
      }
      
      const finalRoll = baseRoll.sum + modifier
      
      // Determine social class based on final roll
      let socialClass = ''
      if (finalRoll <= 10) {
        socialClass = 'Unfrei'
      } else if (finalRoll <= 50) {
        socialClass = 'Volk'
      } else if (finalRoll <= 90) {
        socialClass = 'Mittelschicht'
      } else {
        socialClass = 'Adel'
      }
      
      // Set the form data
      this.formData.stand = socialClass
      
      // Store roll information for display
      this.lastSocialClassRoll = {
        baseRoll: baseRoll.sum,
        modifier: modifier,
        roll: finalRoll,
        result: socialClass
      }
      
      // Save the updated form data
      this.$emit('save', { basic_info: this.formData })
      
      // Show overlay notification
      this.showRollOverlay()
    },
    
    showRollOverlay() {
      this.showOverlay = true
      
      // Clear existing timeout if any
      if (this.overlayTimeout) {
        clearTimeout(this.overlayTimeout)
      }
      
      // Hide overlay after 20 seconds
      this.overlayTimeout = setTimeout(() => {
        this.showOverlay = false
      }, 20000)
    },
    
    hideOverlay() {
      this.showOverlay = false
      if (this.overlayTimeout) {
        clearTimeout(this.overlayTimeout)
        this.overlayTimeout = null
      }
    },
    
    handleSubmit() {
      if (this.isValid) {
        // Save the current state before proceeding
        this.$emit('save', { basic_info: this.formData })
        this.$emit('next', { basic_info: this.formData })
      }
    },
  }
}
</script>

<style>
/* All common styles moved to main.css */

.belief-search {
  position: relative;
}

.belief-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: white;
  border: 1px solid #ddd;
  border-top: none;
  border-radius: 0 0 4px 4px;
  max-height: 200px;
  overflow-y: auto;
  z-index: 1000;
}

.belief-option {
  padding: 10px;
  cursor: pointer;
  border-bottom: 1px solid #eee;
}

.belief-option:hover {
  background-color: #f5f5f5;
}

.belief-option:last-child {
  border-bottom: none;
}

.selected-belief {
  margin-top: 10px;
  padding: 8px;
  background-color: #e3f2fd;
  border-radius: 4px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.clear-btn {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  color: #666;
  padding: 0;
  width: 20px;
  height: 20px;
}

.clear-btn:hover {
  color: #f44336;
}

.roll-result {
  margin-top: 8px;
  padding: 8px;
  background-color: #e8f5e8;
  border-radius: 4px;
  font-size: 14px;
  color: #2e7d32;
}

.roll-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 10000;
  cursor: pointer;
}

.roll-overlay-content {
  background: white;
  padding: 30px;
  border-radius: 10px;
  text-align: center;
  position: relative;
  cursor: default;
  min-width: 300px;
}

.overlay-close {
  position: absolute;
  top: 10px;
  right: 15px;
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
}

.overlay-close:hover {
  color: #000;
}

.overlay-title {
  font-size: 20px;
  font-weight: bold;
  margin-bottom: 15px;
  color: #333;
}

.overlay-roll {
  font-size: 36px;
  font-weight: bold;
  color: #4caf50;
  margin: 20px 0;
}

.roll-breakdown {
  font-size: 18px;
  color: #666;
  margin-left: 10px;
}

.overlay-result {
  font-size: 24px;
  font-weight: bold;
  color: #2196f3;
  margin: 20px 0;
}

.overlay-hint {
  font-size: 14px;
  color: #666;
  margin-top: 15px;
}

.help-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  margin-left: 6px;
  border: 1px solid #999;
  border-radius: 50%;
  font-size: 12px;
  line-height: 1;
  /*cursor: help;*/
  background: #f5f5f5;
  color: #555;
}

.help-icon:hover {
  background: #e0e0e0;
  color: #222;
}
</style>
