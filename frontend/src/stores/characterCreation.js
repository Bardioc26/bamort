import { defineStore } from 'pinia'

export const useCharacterCreationStore = defineStore('characterCreation', {
  state: () => ({
    sessionData: {
      basic_info: null,
      attributes: null,
      derived_values: null,
      skills: [],
      skills_meta: {
        totalUsedPoints: 0,
        selectedCategory: null
      },
      spells: [],
      equipment: null
    },
    currentStep: 1,
    sessionId: null,
    isLoading: false,
    error: null
  }),

  getters: {
    characterClass: (state) => state.sessionData.basic_info?.typ || 'Barbar',
    characterStand: (state) => state.sessionData.basic_info?.stand || 'Buerger',
    characterRace: (state) => state.sessionData.basic_info?.rasse || 'Menschen',
    
    hasSelectedSkills: (state) => state.sessionData.skills && state.sessionData.skills.length > 0,
    totalSkillPoints: (state) => state.sessionData.skills_meta?.totalUsedPoints || 0,
    
    isValid: (state) => {
      // Basic validation - can be extended as needed
      return !!(state.sessionData.basic_info && 
                state.sessionData.attributes && 
                state.sessionData.derived_values)
    }
  },

  actions: {
    // Initialize or load existing session
    async initializeSession(sessionId = null) {
      this.isLoading = true
      this.error = null
      
      try {
        if (sessionId) {
          // Load existing session from backend
          await this.loadSession(sessionId)
        } else {
          // Create new session
          await this.createNewSession()
        }
      } catch (error) {
        this.error = error.message
        console.error('Error initializing session:', error)
      } finally {
        this.isLoading = false
      }
    },

    // Create new character creation session
    async createNewSession() {
      // For now, just initialize empty session
      // In future, this could call backend to create session
      this.sessionData = {
        basic_info: null,
        attributes: null,
        derived_values: null,
        skills: [],
        skills_meta: {
          totalUsedPoints: 0,
          selectedCategory: null
        },
        spells: [],
        equipment: null
      }
      this.currentStep = 1
      console.log('Created new character creation session')
    },

    // Load existing session from backend
    async loadSession(sessionId) {
      // TODO: Implement backend call to load session
      this.sessionId = sessionId
      console.log('Loading session:', sessionId)
    },

    // Save session data to backend
    async saveSession() {
      if (!this.sessionId) {
        console.warn('No session ID available for saving')
        return
      }
      
      try {
        // TODO: Implement backend call to save session
        console.log('Saving session data:', this.sessionData)
      } catch (error) {
        console.error('Error saving session:', error)
        throw error
      }
    },

    // Update specific step data
    updateStepData(stepData) {
      Object.assign(this.sessionData, stepData)
      console.log('Updated session data:', stepData)
      
      // Auto-save after update
      this.saveSession()
    },

    // Navigate between steps
    goToStep(step) {
      this.currentStep = step
    },

    nextStep() {
      this.currentStep += 1
    },

    previousStep() {
      this.currentStep -= 1
    },

    // Clear session (for starting over)
    clearSession() {
      this.sessionData = {
        basic_info: null,
        attributes: null,
        derived_values: null,
        skills: [],
        skills_meta: {
          totalUsedPoints: 0,
          selectedCategory: null
        },
        spells: [],
        equipment: null
      }
      this.currentStep = 1
      this.sessionId = null
      this.error = null
    }
  }
})
