<template>
  <div class="basic-info-form">
    <h2>Basic Character Information</h2>
    
    <form @submit.prevent="handleSubmit">
      <div class="form-group">
        <label for="name">Character Name *</label>
        <input 
          id="name"
          v-model="formData.name"
          type="text"
          required
          minlength="2"
          maxlength="50"
          placeholder="Enter character name"
        />
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="rasse">Race *</label>
          <select id="rasse" v-model="formData.rasse" required>
            <option value="">Select Race</option>
            <option v-for="race in races" :key="race" :value="race">{{ race }}</option>
          </select>
        </div>

        <div class="form-group">
          <label for="typ">Character Class *</label>
          <select id="typ" v-model="formData.typ" required>
            <option value="">Select Class</option>
            <option v-for="cls in classes" :key="cls" :value="cls">{{ cls }}</option>
          </select>
        </div>
      </div>

      <div class="form-group">
        <label for="herkunft">Origin *</label>
        <select id="herkunft" v-model="formData.herkunft" required>
          <option value="">Select Origin</option>
          <option v-for="origin in origins" :key="origin" :value="origin">{{ origin }}</option>
        </select>
      </div>

      <div class="form-group">
        <label for="glaube">Religion/Belief</label>
        <div class="belief-search">
          <input 
            id="glaube"
            v-model="beliefSearch"
            type="text"
            placeholder="Type at least 2 characters to search beliefs..."
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
          Selected: {{ formData.glaube }}
          <button type="button" @click="clearBelief" class="clear-btn">×</button>
        </div>
      </div>

      <div class="form-actions">
        <button type="submit" class="next-btn" :disabled="!isValid">
          Next: Attributes →
        </button>
      </div>
    </form>
  </div>
</template>

<script>
import API from '../../utils/api'

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
        rasse: '',
        typ: '',
        herkunft: '',
        glaube: '',
      },
      races: [],
      classes: [],
      origins: [],
      beliefSearch: '',
      beliefResults: [],
      searchTimeout: null,
    }
  },
  computed: {
    isValid() {
      return this.formData.name.length >= 2 && 
             this.formData.rasse && 
             this.formData.typ && 
             this.formData.herkunft
    }
  },
  async created() {
    // Initialize form with session data
    this.formData = {
      name: this.sessionData.name || '',
      rasse: this.sessionData.rasse || '',
      typ: this.sessionData.typ || '',
      herkunft: this.sessionData.herkunft || '',
      glaube: this.sessionData.glaube || '',
    }
    
    if (this.formData.glaube) {
      this.beliefSearch = this.formData.glaube
    }
    
    await this.loadReferenceData()
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
    
    handleSubmit() {
      if (this.isValid) {
        this.$emit('next', this.formData)
      }
    },
  }
}
</script>

<style scoped>
.basic-info-form {
  max-width: 600px;
  margin: 0 auto;
}

.basic-info-form h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
}

.form-group {
  margin-bottom: 20px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
  color: #555;
}

input, select {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
  box-sizing: border-box;
}

input:focus, select:focus {
  outline: none;
  border-color: #2196f3;
  box-shadow: 0 0 5px rgba(33, 150, 243, 0.3);
}

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

.form-actions {
  text-align: center;
  margin-top: 30px;
}

.next-btn {
  background-color: #2196f3;
  color: white;
  padding: 12px 30px;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.next-btn:hover:not(:disabled) {
  background-color: #1976d2;
}

.next-btn:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}
</style>
