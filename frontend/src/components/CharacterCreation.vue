<template>
  <div class="character-creation">
    <div class="creation-header">
      <h1>Create New Character</h1>
      <div class="progress-indicator">
        <div 
          v-for="step in steps" 
          :key="step.number"
          :class="['step', { 
            active: currentStep === step.number, 
            completed: currentStep > step.number,
            clickable: currentStep > step.number || currentStep === step.number
          }]"
          @click="navigateToStep(step.number)"
        >
          <span class="step-number">{{ step.number }}</span>
          <span class="step-title">{{ step.title }}</span>
        </div>
      </div>
    </div>

    <div class="creation-content">
      <!-- Step 1: Basic Information -->
      <CharacterBasicInfo 
        v-if="currentStep === 1"
        :session-data="sessionData"
        @next="handleNext"
        @save="saveProgress"
      />

      <!-- Step 2: Attributes -->
      <CharacterAttributes 
        v-if="currentStep === 2"
        :session-data="sessionData"
        @next="handleNext"
        @previous="handlePrevious"
        @save="saveProgress"
      />

      <!-- Step 3: Derived Values -->
      <CharacterDerivedValues 
        v-if="currentStep === 3"
        :session-data="sessionData"
        @next="handleNext"
        @previous="handlePrevious"
        @save="saveProgress"
      />

      <!-- Step 4: Skills -->
      <CharacterSkills 
        v-if="currentStep === 4"
        :session-data="sessionData"
        :skill-categories="skillCategories"
        @previous="handlePrevious"
        @next="handleNext"
        @save="saveProgress"
      />

      <!-- Step 5: Spells -->
      <CharacterSpells 
        v-if="currentStep === 5"
        :session-data="sessionData"
        :skill-categories="skillCategories"
        @previous="handlePrevious"
        @finalize="handleFinalize"
        @save="saveProgress"
      />
    </div>

    <!-- Session Info -->
    <div class="session-info">
      <p>Session expires: {{ formatDate(sessionData.expires_at) }}</p>
      <button @click="deleteDraft" class="delete-btn">Delete Draft</button>
    </div>
  </div>
</template>

<script>
import API from '../utils/api'
import CharacterBasicInfo from './CharacterCreation/CharacterBasicInfo.vue'
import CharacterAttributes from './CharacterCreation/CharacterAttributes.vue'
import CharacterDerivedValues from './CharacterCreation/CharacterDerivedValues.vue'
import CharacterSkills from './CharacterCreation/CharacterSkills.vue'
import CharacterSpells from './CharacterCreation/CharacterSpells.vue'

export default {
  name: 'CharacterCreation',
  components: {
    CharacterBasicInfo,
    CharacterAttributes,
    CharacterDerivedValues,
    CharacterSkills,
    CharacterSpells,
  },
  props: {
    sessionId: {
      type: String,
      required: true,
    }
  },
  data() {
    return {
      currentStep: 1,
      sessionData: {
        id: '',
        name: '',
        rasse: '',
        typ: '',
        herkunft: '',
        glaube: '',
        geschlecht: '',
        stand: '',
        attributes: {},
        derived_values: {},
        skills: [],
        spells: [],
        skill_points: {},
        spell_points: {},
        expires_at: null,
        current_step: 1,
      },
      steps: [
        { number: 1, title: 'Basic Info' },
        { number: 2, title: 'Attributes' },
        { number: 3, title: 'Derived Values' },
        { number: 4, title: 'Skills' },
        { number: 5, title: 'Spells' },
      ],
      skillCategories: [],
    }
  },
  async created() {
    await this.loadSession()
    await this.loadSkillCategories()
  },
  methods: {
    async loadSession() {
      try {
        const token = localStorage.getItem('token')
        const response = await API.get(`/api/characters/create-session/${this.sessionId}`, {
          headers: { Authorization: `Bearer ${token}` },
        })
        
        this.sessionData = response.data
        this.currentStep = response.data.current_step || 1
      } catch (error) {
        console.error('Error loading session:', error)
        this.$router.push('/dashboard')
      }
    },
    
    async loadSkillCategories() {
      try {
        const token = localStorage.getItem('token')
        const response = await API.get('/api/characters/skill-categories', {
          headers: { Authorization: `Bearer ${token}` },
        })
        
        this.skillCategories = response.data.categories || []
      } catch (error) {
        console.error('Error loading skill categories:', error)
        // Fallback dummy data
        this.skillCategories = [
          { name: 'körperlich', display_name: 'Körperliche Fertigkeiten', max_points: 200, points: 200 },
          { name: 'gesellschaftlich', display_name: 'Gesellschaftliche Fertigkeiten', max_points: 150, points: 150 },
          { name: 'natur', display_name: 'Natur Fertigkeiten', max_points: 100, points: 100 },
          { name: 'wissen', display_name: 'Wissens Fertigkeiten', max_points: 180, points: 180 },
          { name: 'handwerk', display_name: 'Handwerks Fertigkeiten', max_points: 120, points: 120 },
          { name: 'zauber', display_name: 'Zauber', max_points: 300, points: 300 },
        ]
      }
    },
    
    async handleNext(data) {
      try {
        // Merge the new data
        this.sessionData = { ...this.sessionData, ...data }
        
        // Save progress for current step before moving to next
        await this.saveProgressForStep(this.currentStep, data)
        
        // Move to next step
        this.currentStep++
      } catch (error) {
        console.error('Failed to save progress before moving to next step:', error)
        // Don't move to next step if save failed
      }
    },
    
    async saveProgressForStep(step, data) {
      try {
        const token = localStorage.getItem('token')
        
        let endpoint = ''
        let payload = {}
        
        switch (step) {
          case 1:
            endpoint = `/api/characters/create-session/${this.sessionId}/basic`
            // Handle both old format and new basic_info format
            const basicInfo = data.basic_info || data
            payload = {
              name: basicInfo.name || this.sessionData.name || '',
              geschlecht: basicInfo.geschlecht || this.sessionData.geschlecht || '',
              rasse: basicInfo.rasse || this.sessionData.rasse || '',
              typ: basicInfo.typ || this.sessionData.typ || '',
              herkunft: basicInfo.herkunft || this.sessionData.herkunft || '',
              stand: basicInfo.stand || this.sessionData.stand || '',
              glaube: basicInfo.glaube || this.sessionData.glaube || '',
            }
            // Validate that all required fields are present
            if (!payload.name || !payload.geschlecht || !payload.rasse || !payload.typ || !payload.herkunft || !payload.stand) {
              throw new Error(`Missing required fields: name=${payload.name}, geschlecht=${payload.geschlecht}, rasse=${payload.rasse}, typ=${payload.typ}, herkunft=${payload.herkunft}, stand=${payload.stand}`)
            }
            break
          case 2:
            endpoint = `/api/characters/create-session/${this.sessionId}/attributes`
            payload = data.attributes || data
            break
          case 3:
            endpoint = `/api/characters/create-session/${this.sessionId}/derived`
            payload = data.derived_values || data
            break
          case 4:
            endpoint = `/api/characters/create-session/${this.sessionId}/skills`
            payload = {
              skills: data.skills || this.sessionData.skills,
              spells: data.spells || this.sessionData.spells,
              skill_points: data.skill_points || this.sessionData.skill_points,
            }
            break
        }
        
        if (endpoint) {
          const response = await API.put(endpoint, payload, {
            headers: { Authorization: `Bearer ${token}` },
          })
        }
      } catch (error) {
        console.error('Error saving progress for step', step, ':', error)
        
        // Provide more specific error messages
        if (error.response && error.response.status === 401) {
          alert('Your session has expired. Please log in again.')
        } else if (error.response && error.response.status === 400) {
          const errorMsg = error.response.data?.error || 'Invalid data submitted'
          alert(`Error saving character data: ${errorMsg}`)
        } //else {
          //alert('Failed to save character data. Please try again.')
        //}
        throw error // Re-throw to handle in calling function
      }
    },
    
    handlePrevious() {
      this.currentStep--
    },
    
    navigateToStep(stepNumber) {
      // Only allow navigation to current step or previously completed steps
      if (stepNumber <= this.currentStep) {
        this.currentStep = stepNumber
        // Save current progress before switching steps (no data parameter needed here)
        this.saveProgress().catch(error => {
          console.error('Failed to save progress during navigation:', error)
        })
      }
    },
    
    async saveProgress(data = null) {
      try {
        // Use provided data or current sessionData as fallback
        const dataToSave = data || this.sessionData
        
        // Update sessionData with new data if provided
        if (data) {
          this.sessionData = { ...this.sessionData, ...data }
        }
        
        // Save progress for current step
        await this.saveProgressForStep(this.currentStep, dataToSave)
      } catch (error) {
        console.error('Failed to save progress:', error)
        throw error
      }
    },
    
    async handleFinalize() {
      try {
        const token = localStorage.getItem('token')
        const response = await API.post(`/api/characters/create-session/${this.sessionId}/finalize`, {}, {
          headers: { Authorization: `Bearer ${token}` },
        })
        
        const characterId = response.data.character_id
        
        // Success message
        //alert('Character successfully created!')
        
        // Navigate to character view or back to character list
        this.$router.push(`/character/${characterId}`)
      } catch (error) {
        console.error('Error finalizing character:', error)
        if (error.response?.data?.error) {
          alert(`Error: ${error.response.data.error}`)
        } else {
          alert('Fehler beim Abschließen der Charakter-Erstellung')
        }
      }
    },
    
    async deleteDraft() {
      if (confirm('Are you sure you want to delete this character draft?')) {
        try {
          const token = localStorage.getItem('token')
          await API.delete(`/api/characters/create-session/${this.sessionId}`, {
            headers: { Authorization: `Bearer ${token}` },
          })
          
          this.$router.push('/dashboard')
        } catch (error) {
          console.error('Error deleting session:', error)
        }
      }
    },
    
    formatDate(dateString) {
      if (!dateString) return ''
      return new Date(dateString).toLocaleDateString()
    },
  },
}
</script>

<style scoped>
.character-creation {
  width: 100%;
  max-width: none;
  margin: 0;
  padding: 10px;
}

.creation-content {
  width: 100%;
}

.creation-header {
  margin-bottom: 30px;
}

.creation-header h1 {
  text-align: center;
  margin-bottom: 20px;
  color: #333;
}

.progress-indicator {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
  margin-bottom: 20px;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 10px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.step.active {
  background-color: #e3f2fd;
  border: 2px solid #2196f3;
}

.step.completed {
  background-color: #e8f5e8;
  border: 2px solid #4caf50;
}

.step-number {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background-color: #ddd;
  color: #666;
  font-weight: bold;
  margin-bottom: 5px;
}

.step.active .step-number {
  background-color: #2196f3;
  color: white;
}

.step.completed .step-number {
  background-color: #4caf50;
  color: white;
}

.step.clickable {
  cursor: pointer;
  transition: all 0.3s ease;
}

.step.clickable:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}

.step.completed:hover .step-number {
  background-color: #45a049;
}

.step.active:hover .step-number {
  background-color: #1976d2;
}

.step-title {
  font-size: 12px;
  color: #666;
  text-align: center;
}

.creation-content {
  background: white;
  border-radius: 8px;
  padding: 30px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  margin-bottom: 20px;
}

.session-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  background-color: #f5f5f5;
  border-radius: 4px;
  font-size: 14px;
  color: #666;
}

.delete-btn {
  background-color: #f44336;
  color: white;
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.delete-btn:hover {
  background-color: #d32f2f;
}
</style>
