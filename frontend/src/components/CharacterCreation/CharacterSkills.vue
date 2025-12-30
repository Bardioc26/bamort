<template>
  <div class="fullwidth-page">
    <div class="page-header">
      <h2>Fertigkeiten & Lernpunkte</h2>
    </div>
    <p style="color: #666; margin-bottom: 30px; font-size: 16px; line-height: 1.5;">
      Wählen Sie Fertigkeiten für Ihren {{ characterClass }}-Charakter aus. 
      Jede Kategorie hat eine begrenzte Anzahl von Lernpunkten.
    </p>

    <!-- Loading State -->
    <div v-if="isLoading" class="loading-message">
      <div class="loading-spinner"></div>
      <p>Lade Lernpunkte für {{ characterClass }}...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-message">
      {{ error }}
      <br>
      <button @click="retry" class="btn btn-primary" style="margin-top: 10px;">Erneut versuchen</button>
    </div>

    <!-- Main Content -->
    <div v-else class="skills-content">
      <!-- Three Column Layout -->
      <div class="three-column-grid">
        
        <!-- Left Column: Available Skills for Selected Category -->
        <div>
          <div v-if="selectedCategory" class="skills-content">
            <div class="list-container">
              <div class="section-header" style="padding: 20px 20px 10px;">
                <h3>{{ getSelectedCategoryName() }} - Verfügbare Fertigkeiten</h3>
              </div>
              
              <!-- Loading skills -->
              <div v-if="isLoadingSkills" class="loading-message">
                <div class="loading-spinner"></div>
                <p>Lade Fertigkeiten...</p>
              </div>

              <!-- Available skills -->
              <div v-else-if="availableSkillsForSelectedCategory.length > 0">
                <div 
                  v-for="skill in availableSkillsForSelectedCategory"
                  :key="skill.name"
                  class="list-item"
                  :class="{ 'opacity-50': !canAffordSkillInCategory(skill) }"
                >
                  <div class="list-item-content">
                    <div class="list-item-title">{{ skill.name }}</div>
                    <div class="list-item-details">
                      <span class="badge badge-primary">{{ skill.cost }} LE</span>
                      <span v-if="skill.attribute" class="badge badge-secondary">({{ skill.attribute }})</span>
                    </div>
                  </div>
                  <div class="list-item-actions">
                    <button 
                      @click="selectSkillForLearning(skill)"
                      :disabled="!canAffordSkillInCategory(skill)"
                      class="btn btn-primary"
                      style="font-size: 12px;"
                    >
                      Hinzufügen
                    </button>
                  </div>
                </div>
              </div>

              <!-- No skills found -->
              <div v-else class="empty-state">
                <p>Keine Fertigkeiten für diese Kategorie gefunden.</p>
              </div>
            </div>
          </div>
          
          <!-- No category selected -->
          <div v-else-if="learningCategories.length > 0" class="empty-state">
            <p>Wählen Sie eine Lernpunkte-Kategorie aus, um verfügbare Fertigkeiten zu sehen.</p>
          </div>
        </div>

        <!-- Center Column: Selected Skills -->
        <div>
          <div class="list-container">
            <div class="section-header" style="padding: 20px 20px 10px;">
              <h3>Gewählte Fertigkeiten</h3>
            </div>
            
            <div v-if="selectedSkills.length > 0">
              <div 
                v-for="skill in selectedSkills"
                :key="skill.name"
                class="list-item"
              >
                <div class="list-item-content">
                  <div class="list-item-title">{{ skill.name }}</div>
                  <div class="list-item-details">
                    <span class="badge badge-primary">{{ skill.cost }} LE</span>
                    <span class="badge badge-info">{{ skill.category }}</span>
                  </div>
                </div>
                <div class="list-item-actions">
                  <button 
                    @click="removeSkill(skill)"
                    class="btn btn-danger"
                    style="width: 30px; height: 30px; border-radius: 50%; padding: 0;"
                    title="Fertigkeit entfernen"
                  >
                    ×
                  </button>
                </div>
              </div>
              
              <!-- Summary -->
              <div class="card" style="margin: 15px 20px;">
                <div class="resource-card" style="justify-content: center;">
                  <div class="resource-info" style="text-align: center;">
                    <div class="resource-label">Gesamt verbrauchte LE:</div>
                    <div class="resource-amount">{{ totalUsedPoints }}</div>
                  </div>
                </div>
              </div>
            </div>

            <div v-else class="empty-state">
              <p>Noch keine Fertigkeiten gewählt.</p>
              <p style="font-style: italic; margin-top: 10px; font-size: 14px; color: #6c757d;">Wählen Sie eine Kategorie und fügen Sie Fertigkeiten hinzu.</p>
            </div>
          </div>
        </div>

        <!-- Right Column: Learning Points + Typical Skills -->
        <div>
          <!-- Learning Points Overview -->
          <div v-if="learningCategories.length > 0" style="margin-bottom: 30px;">
            <div class="section-header">
              <h3>Verfügbare Lernpunkte</h3>
            </div>
            <div style="display: flex; flex-direction: column; gap: 15px;">
              <div 
                v-for="category in learningCategories"
                :key="category.name"
                class="card"
                :class="{
                  'border-primary': selectedCategory === category.name,
                  'border-danger': category.remainingPoints === 0,
                  'border-warning': category.remainingPoints < category.totalPoints && category.remainingPoints > 0
                }"
                @click="selectCategory(category.name)"
                style="cursor: pointer;"
              >
                <div class="list-item-title">{{ category.displayName }}</div>
                <div class="resource-display">
                  <div class="resource-card" style="justify-content: center; text-align: center;">
                    <div class="resource-info">
                      <div class="resource-amount">
                        <span style="color: #28a745; font-size: 16px;">{{ category.remainingPoints }}</span>
                        <span style="color: #6c757d;">/</span>
                        <span style="color: #6c757d;">{{ category.totalPoints }}</span>
                        <span style="color: #6c757d; font-size: 12px;">LE</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Typical Skills Info -->
          <div v-if="typicalSkills.length > 0" class="card">
            <div class="section-header">
              <h3>Empfohlene Fertigkeiten für {{ characterClass }}</h3>
            </div>
            <p style="color: #6c757d; margin-bottom: 15px;">
              Die folgenden Fertigkeiten werden häufig von Charakteren Ihrer Klasse erlernt:
            </p>
            <div style="display: flex; flex-wrap: wrap; gap: 8px;">
              <span 
                v-for="skill in typicalSkills" 
                :key="skill.name"
                class="badge badge-info"
                :title="`${skill.name} (${skill.attribute}) - Bonus: +${skill.bonus}`"
              >
                {{ skill.name }} ({{ skill.attribute }})
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Navigation -->
    <div class="form-row" style="justify-content: space-between; padding-top: 20px; border-top: 1px solid #dee2e6; margin-top: 30px;">
      <button type="button" @click="handlePrevious" class="btn btn-secondary">
        ← Zurück: Abgeleitete Werte
      </button>
      <button 
        type="button" 
        @click="handleNext" 
        class="btn btn-primary"
        :disabled="!canProceed"
      >
        Weiter: Zauber →
      </button>
    </div>

    <!-- Debug Info (removable) -->
    <div v-if="showDebug || error" class="card" style="margin-top: 20px; font-family: monospace; font-size: 12px;">
      <div class="section-header">
        <h4>Debug Information</h4>
        <button @click="showDebug = !showDebug" class="btn btn-sm">{{ showDebug ? 'Hide' : 'Show' }} Debug</button>
      </div>
      <div v-if="showDebug || error">
        <pre style="margin: 0; white-space: pre-wrap; word-break: break-all;">{{ debugInfo }}</pre>
      </div>
    </div>
  </div>
</template>

<script>
import API from '../../utils/api'

export default {
  name: 'CharacterSkills',
  props: {
    sessionData: {
      type: Object,
      required: true,
    }
  },
  emits: ['next', 'previous', 'save'],
  data() {
    return {
      // Loading states
      isLoading: false,
      isLoadingSkills: false,
      error: null,
      
      // Learning data
      learningPointsData: null,
      learningCategories: [],
      typicalSkills: [],
      
      // Skills selection
      selectedCategory: null,
      availableSkills: [],
      selectedSkills: [],
      
      // Available skills fetching
      availableSkillsByCategory: null,
      
      // Debug
      showDebug: false // Set to true for debugging
    }
  },
  computed: {
    // Get character class from session data (direct access like in working version)
    characterClass() {
      return this.sessionData?.typ || ''
    },
    
    // Get character stand from session data (direct access like in working version)
    characterStand() {
      return this.sessionData?.stand || ''
    },
    
    totalUsedPoints() {
      return this.selectedSkills.reduce((sum, skill) => sum + (skill.cost || 0), 0)
    },
    
    canProceed() {
      // Allow proceeding even without skills selected - this is optional
      return true
    },
    
    debugInfo() {
      return {
        // Session Data
        hasSessionData: !!this.sessionData,
        sessionDataKeys: this.sessionData ? Object.keys(this.sessionData) : null,
        basicInfo: this.sessionData?.basic_info || null,
        
        // Character Info
        characterClass: this.characterClass,
        characterStand: this.characterStand,
        
        // Loading States
        isLoading: this.isLoading,
        isLoadingSkills: this.isLoadingSkills,
        error: this.error,
        
        // Data States
        hasLearningPointsData: !!this.learningPointsData,
        learningCategories: this.learningCategories.length,
        hasAvailableSkills: !!this.availableSkillsByCategory,
        selectedSkillsCount: this.selectedSkills.length,
        selectedCategory: this.selectedCategory,
        totalUsedPoints: this.totalUsedPoints,
        
        // Raw Data (for debugging)
        learningPointsData: this.learningPointsData,
        availableSkillsByCategory: this.availableSkillsByCategory
      }
    },

    // Skills for the selected category (from Learning Points section)
    availableSkillsForSelectedCategory() {
      if (!this.selectedCategory || !this.availableSkillsByCategory) {
        return []
      }
      
      // Try to find the matching category key
      const categoryKey = this.findCategoryKey(this.selectedCategory)
      if (!categoryKey) {
        return []
      }
      
      // Get skills from the selected category
      const categorySkills = this.availableSkillsByCategory[categoryKey] || []
      
      const filteredSkills = categorySkills.map(skill => ({
        ...skill,
        cost: this.getSkillCost(skill),
        category: categoryKey // Use the actual category key from availableSkillsByCategory
      }))
      .filter(skill => {
        // Remove already selected skills
        const selectedSkillNames = this.selectedSkills.map(s => s.name)
        return !selectedSkillNames.includes(skill.name)
      })
      .filter(skill => this.canAffordSkillInCategory(skill))
      .sort((a, b) => {
        const aLeCost = this.getSkillCost(a)
        const bLeCost = this.getSkillCost(b)
        
        // First sort by LE cost (ascending)
        if (aLeCost !== bLeCost) {
          return aLeCost - bLeCost
        }
        
        // If costs are equal, sort alphabetically
        return a.name.localeCompare(b.name)
      })
      
      return filteredSkills
    },

    totalSelectedEP() {
      return this.selectedSkills.reduce((total, skill) => total + this.getSkillCost(skill), 0)
    },

    totalSelectedGold() {
      // For character creation, we only track learning costs, not gold costs
      return 0
    }
  },
  watch: {
    selectedSkills: {
      handler(newSkills) {
        // Save skills automatically when they change
        this.saveSkillsToSession()
      },
      deep: true
    }
  },
  created() {
    // Initialize with existing session data if available
    if (this.sessionData.skills && Array.isArray(this.sessionData.skills)) {
      this.selectedSkills = [...this.sessionData.skills]
    }
    
    // Restore selected category if available
    if (this.sessionData.skills_meta?.selectedCategory) {
      this.selectedCategory = this.sessionData.skills_meta.selectedCategory
    }
    
    // Initialize component
    this.initializeComponent()
  },
  beforeUnmount() {
    // Ensure skills are saved when component is about to be unmounted
    this.saveSkillsToSession()
  },
  methods: {
    saveSkillsToSession() {
      // Save skills to session data via emit
      this.$emit('save', { 
        skills: this.selectedSkills,
        skills_meta: {
          totalUsedPoints: this.totalUsedPoints,
          selectedCategory: this.selectedCategory
        }
      })
    },
    
    async initializeComponent() {
      try {
        this.isLoading = true
        this.error = null
        
        // Load learning points from backend
        await this.loadLearningPoints()
        
        // Load available skills
        await this.loadAvailableSkills()
        
        // Restore previously selected skills if any
        this.restoreSelectedSkills()
        
        // Update points based on selected skills
        this.updateRemainingPoints()
        
      } catch (error) {
        console.error('Error initializing component:', error)
        this.error = 'Fehler beim Laden der Lernpunkte. Bitte versuchen Sie es erneut.'
      } finally {
        this.isLoading = false
      }
    },
    
    async loadLearningPoints() {
      if (!this.characterClass) {
        throw new Error('Charakterklasse nicht verfügbar')
      }
      
      const params = {
        class: this.characterClass
      }
      
      if (this.characterStand) {
        params.stand = this.characterStand
      }
      
      const response = await API.get('/api/characters/classes/learning-points', { params })
      this.learningPointsData = response.data
      
      // Process learning points into categories
      this.processLearningPoints()
      
      // Store typical skills
      this.typicalSkills = this.learningPointsData.typical_skills || []
    },
    
    processLearningPoints() {
      const learningPoints = this.learningPointsData.learning_points || {}
      
      this.learningCategories = Object.entries(learningPoints).map(([categoryKey, points]) => ({
        name: categoryKey.toLowerCase(),
        displayName: categoryKey,
        totalPoints: points,
        remainingPoints: points
      }))
      
      // Add weapon points as separate category if available
      if (this.learningPointsData.weapon_points && this.learningPointsData.weapon_points > 0) {
        console.log('Adding weapon category with points:', this.learningPointsData.weapon_points)
        this.learningCategories.push({
          name: 'waffen',
          displayName: 'Waffenfertigkeiten',
          totalPoints: this.learningPointsData.weapon_points,
          remainingPoints: this.learningPointsData.weapon_points
        })
      } else {
        console.log('No weapon points found:', this.learningPointsData.weapon_points)
      }
    },
    
    restoreSelectedSkills() {
      if (this.sessionData.skills && Array.isArray(this.sessionData.skills)) {
        // Simply restore skills as they are
        this.selectedSkills = this.sessionData.skills.map(skill => {
          return { ...skill }
        })
        console.log('Restored selected skills:', this.selectedSkills)
      }
    },
    
    updateRemainingPoints() {
      // Reset all categories to full points
      this.learningCategories.forEach(category => {
        category.remainingPoints = category.totalPoints
      })
      
      // Deduct points for selected skills
      this.selectedSkills.forEach(skill => {
        const category = this.learningCategories.find(cat => 
          cat.name === skill.category?.toLowerCase()
        )
        if (category && skill.cost) {
          category.remainingPoints = Math.max(0, category.remainingPoints - this.getSkillCost(skill))
        }
      })
    },
    
    async selectCategory(categoryName) {
      this.selectedCategory = categoryName
      // Skills are now provided by computed property availableSkillsForSelectedCategory
    },
    
    getSelectedCategoryName() {
      const category = this.learningCategories.find(cat => cat.name === this.selectedCategory)
      return category ? category.displayName : this.selectedCategory
    },

    findCategoryKey(selectedCategoryName) {
      if (!this.availableSkillsByCategory) return null
      
      // Try to find by learning category mapping first (most likely scenario)
      const learningCategory = this.learningCategories?.find(cat => cat.name === selectedCategoryName)
      if (learningCategory && this.availableSkillsByCategory[learningCategory.displayName]) {
        return learningCategory.displayName
      }
      
      // Try direct match
      if (this.availableSkillsByCategory[selectedCategoryName]) {
        return selectedCategoryName
      }
      
      // Try case-insensitive search
      const availableKeys = Object.keys(this.availableSkillsByCategory)
      const foundKey = availableKeys.find(key => 
        key.toLowerCase() === selectedCategoryName.toLowerCase()
      )
      
      return foundKey || null
    },

    getSkillCost(skill) {
      // Unified method to get skill cost from various possible properties
      return skill.cost || skill.learnCost || skill.leCost || 0
    },

    canAffordSkillInCategory(skill) {
      // Check if character has enough learning points in the skill's category
      // Handle both displayName format (e.g., "Sozial") and internal name format (e.g., "sozial")
      let category = this.learningCategories.find(cat => 
        cat.displayName === skill.category
      )
      
      // If not found by displayName, try by internal name
      if (!category) {
        category = this.learningCategories.find(cat => 
          cat.name === skill.category?.toLowerCase()
        )
      }
      
      if (!category) {
        console.warn('No learning category found for skill:', skill.name, 'category:', skill.category)
        return false
      }
      
      const skillCost = this.getSkillCost(skill)
      return category.remainingPoints >= skillCost
    },
    
    removeSkill(skill) {
      const index = this.selectedSkills.findIndex(s => s.name === skill.name)
      if (index >= 0) {
        this.selectedSkills.splice(index, 1)
        this.updateRemainingPoints()
        
        console.log('Removed skill:', skill.name, 'Selected skills:', this.selectedSkills)
      }
    },
    
    async retry() {
      await this.initializeComponent()
    },
    
    handlePrevious() {
      this.$emit('previous')
    },
    
    handleNext() {
      const data = {
        skills: this.selectedSkills,
        skill_points: this.learningCategories.reduce((acc, cat) => {
          acc[cat.name] = cat.remainingPoints
          return acc
        }, {}),
        learning_points_data: this.learningPointsData
      }
      
      console.log('Proceeding with data:', data)
      this.$emit('next', data)
    },

    // Skills methods
    async loadAvailableSkills() {
      if (!this.characterClass) {
        return
      }

      this.isLoadingSkills = true
      try {
        const requestData = {
          characterClass: this.characterClass
        }
        
        const response = await API.post('/api/characters/available-skills-creation', requestData)
        
        if (response.data && response.data.skills_by_category) {
          this.availableSkillsByCategory = response.data.skills_by_category
        } else {
          this.generateSampleSkills()
        }
        
      } catch (error) {
        console.error('Error loading skills:', error)
        this.generateSampleSkills()
      } finally {
        this.isLoadingSkills = false
      }
    },

    generateSampleSkills() {
      // Fallback for testing
      this.availableSkillsByCategory = {
        'Körperliche Fertigkeiten': [
          { name: 'Klettern', learnCost: 50 },
          { name: 'Schwimmen', learnCost: 40 },
          { name: 'Springen', learnCost: 30 }
        ],
        'Geistige Fertigkeiten': [
          { name: 'Erste Hilfe', learnCost: 60 },
          { name: 'Naturkunde', learnCost: 75 },
          { name: 'Menschenkenntnis', learnCost: 65 }
        ],
        'Handwerkliche Fertigkeiten': [
          { name: 'Bogenbau', learnCost: 100 },
          { name: 'Schmieden', learnCost: 125 }
        ]
      }
      
      console.log('Generated sample skills:', this.availableSkillsByCategory)
    },

    isSkillSelected(skill) {
      return this.selectedSkills.some(s => s.name === skill.name)
    },

    selectSkillForLearning(skill) {
      if (this.isSkillSelected(skill)) {
        return // Already selected
      }
      
      // Check if the skill can be afforded
      if (!this.canAffordSkillInCategory(skill)) {
        console.log('Cannot afford skill:', skill.name)
        return
      }
      
      // Add skill to selected list with proper cost and category
      const skillToAdd = {
        ...skill,
        cost: this.getSkillCost(skill), // Ensure cost is properly set
        category: skill.category // Use the category as-is
      }
      
      this.selectedSkills.push(skillToAdd)
      
      // Update remaining points for all categories
      this.updateRemainingPoints()
      
      console.log('Skill selected for learning:', skill.name, 'Cost:', skillToAdd.cost)
      console.log('Updated remaining points')
    },

    removeSkillFromSelection(skill) {
      const index = this.selectedSkills.findIndex(s => s.name === skill.name)
      if (index !== -1) {
        this.selectedSkills.splice(index, 1)
        
        // Update remaining points for all categories after removal
        this.updateRemainingPoints()
        
        console.log('Skill removed from selection:', skill.name)
        console.log('Updated remaining points after removal')
      }
    },

    // Navigation methods
    handlePrevious() {
      this.$emit('previous')
    },

    handleNext() {
      // Save current skills before navigating
      this.$emit('next', { 
        skills: this.selectedSkills,
        skills_meta: {
          totalUsedPoints: this.totalUsedPoints,
          selectedCategory: this.selectedCategory
        }
      })
    },

    retry() {
      this.error = null
      this.initializeComponent()
    }
  }
}
</script>

<style>
/* All common styles moved to main.css */

.fullwidth-page {
  padding: 0 !important;
  margin: 0 !important;
  width: 100vw !important;
  max-width: 100vw !important;
  box-sizing: border-box !important;
}

.page-header {
  padding: 15px 20px;
  margin-bottom: 20px;
}

.three-column-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 30px;
  margin: 0 20px 30px 20px;
  width: calc(100vw - 40px);
  max-width: calc(100vw - 40px);
  box-sizing: border-box;
}

.skills-content {
  width: 100%;
  max-width: 100%;
}

.opacity-50 {
  opacity: 0.6;
}

.border-primary {
  border-color: #007bff !important;
  background-color: #f8fcff;
}

.border-warning {
  border-color: #ffc107 !important;
  background-color: #fff8e1;
}

.border-danger {
  border-color: #dc3545 !important;
  background-color: #ffebee;
}

@media (max-width: 1200px) {
  .three-column-grid {
    grid-template-columns: 1fr 1fr;
    gap: 20px;
  }
}

@media (max-width: 768px) {
  .three-column-grid {
    grid-template-columns: 1fr;
    gap: 15px;
  }
}
</style>
