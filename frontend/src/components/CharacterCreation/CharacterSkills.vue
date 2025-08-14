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
      <!-- Typical Skills Info -->
      <div v-if="typicalSkills.length > 0" class="card" style="margin-bottom: 30px;">
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

      <!-- Learning Points Overview -->
      <div v-if="learningCategories.length > 0" style="margin-bottom: 30px;">
        <div class="section-header">
          <h3>Verfügbare Lernpunkte</h3>
        </div>
        <div class="grid-container grid-3-columns">
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
                    <span style="color: #6c757d; font-size: 12px;">LP</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Skills Selection -->
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
          <div v-else-if="availableSkills.length > 0">
            <div 
              v-for="skill in availableSkills"
              :key="skill.name"
              class="list-item"
              :class="{ 'opacity-50': !canAddSkill(skill) }"
            >
              <div class="list-item-content">
                <div class="list-item-title">{{ skill.name }}</div>
                <div class="list-item-details">
                  <span class="badge badge-primary">{{ skill.cost }} LP</span>
                  <span v-if="skill.attribute" class="badge badge-secondary">({{ skill.attribute }})</span>
                </div>
              </div>
              <div class="list-item-actions">
                <button 
                  @click="addSkill(skill)"
                  :disabled="!canAddSkill(skill)"
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

        <!-- Selected Skills -->
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
                  <span class="badge badge-primary">{{ skill.cost }} LP</span>
                  <span class="badge badge-info">{{ skill.categoryDisplay }}</span>
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
                  <div class="resource-label">Gesamt verbrauchte LP:</div>
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

      <!-- Available Skills Section -->
      <div v-if="!isLoadingSkills && learningCategories.length > 0" style="margin-top: 30px;">
        <div class="section-header">
          <h3>Verfügbare Fertigkeiten</h3>
        </div>
        <p style="color: #6c757d; margin-bottom: 20px;">
          Wählen Sie aus den verfügbaren Fertigkeiten für jede Kategorie:
        </p>

        <!-- Category Filter -->
        <div class="form-row" style="margin-bottom: 20px; gap: 10px; flex-wrap: wrap;">
          <button 
            @click="setSkillCategoryFilter(null)"
            class="btn"
            :class="selectedSkillCategoryFilter === null ? 'btn-primary' : 'btn-secondary'"
          >
            Alle Kategorien
          </button>
          <button 
            v-for="category in availableSkillCategories" 
            :key="category"
            @click="setSkillCategoryFilter(category)"
            class="btn"
            :class="selectedSkillCategoryFilter === category ? 'btn-primary' : 'btn-secondary'"
          >
            {{ category }}
          </button>
        </div>

        <!-- Skills Loading State -->
        <div v-if="isLoadingSkills" class="loading-message">
          <div class="loading-spinner"></div>
          <p>Lade verfügbare Fertigkeiten...</p>
        </div>

        <!-- Skills List -->
        <div v-else-if="filteredAvailableSkills.length > 0" class="list-container">
          <div 
            v-for="skill in filteredAvailableSkills" 
            :key="skill.name"
            class="list-item"
            style="display: flex; justify-content: space-between; align-items: center;"
            :class="{ 'opacity-50': !canAffordSkill(skill) }"
          >
            <div>
              <div style="font-weight: 600; color: #2c3e50;">{{ skill.name }}</div>
              <div style="display: flex; gap: 15px; font-size: 14px; color: #6c757d; margin-top: 4px;">
                <span>{{ skill.category }}</span>
                <span style="color: #007acc; font-weight: 500;">{{ skill.epCost }} EP</span>
                <span style="color: #28a745; font-weight: 500;">{{ skill.goldCost }} GS</span>
              </div>
            </div>
            <button 
              @click="selectSkillForLearning(skill)"
              class="btn btn-primary"
              :disabled="!canAffordSkill(skill) || isSkillSelected(skill)"
            >
              {{ isSkillSelected(skill) ? '✓' : '→' }}
            </button>
          </div>
        </div>

        <!-- No Skills Available -->
        <div v-else-if="!isLoadingSkills" class="empty-state">
          <p>Keine Fertigkeiten verfügbar für die gewählte Kategorie.</p>
        </div>
      </div>

      <!-- No category selected -->
      <div v-else-if="learningCategories.length > 0" class="empty-state">
        <p>Wählen Sie eine Lernpunkte-Kategorie aus, um verfügbare Fertigkeiten zu sehen.</p>
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
    <div v-if="showDebug" class="card" style="margin-top: 20px; font-family: monospace; font-size: 12px;">
      <div class="section-header">
        <h4>Debug Information</h4>
      </div>
      <pre style="margin: 0; white-space: pre-wrap; word-break: break-all;">{{ debugInfo }}</pre>
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
      required: true
    },
    skillCategories: {
      type: Array,
      default: () => []
    }
  },
  emits: ['previous', 'next', 'save'],
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
      selectedSkillCategoryFilter: null,
      
      // Debug
      showDebug: false // Set to true for debugging
    }
  },
  computed: {
    characterClass() {
      return this.sessionData?.typ || 'Unbekannt'
    },
    
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
        characterClass: this.characterClass,
        characterStand: this.characterStand,
        learningPointsData: this.learningPointsData,
        learningCategories: this.learningCategories,
        selectedCategory: this.selectedCategory,
        selectedSkills: this.selectedSkills,
        totalUsedPoints: this.totalUsedPoints,
        availableSkillsByCategory: this.availableSkillsByCategory
      }
    },

    availableSkillCategories() {
      if (!this.availableSkillsByCategory) return []
      return Object.keys(this.availableSkillsByCategory)
    },

    filteredAvailableSkills() {
      if (!this.availableSkillsByCategory) return []
      
      let allSkills = []
      
      // Collect all skills from all categories
      Object.keys(this.availableSkillsByCategory).forEach(category => {
        this.availableSkillsByCategory[category].forEach(skill => {
          allSkills.push({
            ...skill,
            category: category
          })
        })
      })
      
      // Apply category filter
      if (this.selectedSkillCategoryFilter) {
        allSkills = allSkills.filter(skill => skill.category === this.selectedSkillCategoryFilter)
      }
      
      // Remove already selected skills
      const selectedSkillNames = this.selectedSkills.map(s => s.name)
      allSkills = allSkills.filter(skill => !selectedSkillNames.includes(skill.name))
      
      return allSkills
    },

    totalSelectedEP() {
      return this.selectedSkills.reduce((total, skill) => total + (skill.epCost || 0), 0)
    },

    totalSelectedGold() {
      return this.selectedSkills.reduce((total, skill) => total + (skill.goldCost || 0), 0)
    }
  },
  async mounted() {
    console.log('CharacterSkills mounted with sessionData:', this.sessionData)
    await this.initializeComponent()
  },
  methods: {
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
      const params = {
        class: this.characterClass
      }
      
      if (this.characterStand) {
        params.stand = this.characterStand
      }
      
      console.log('Loading learning points with params:', params)
      
      const response = await API.get('/api/characters/classes/learning-points', { params })
      this.learningPointsData = response.data
      
      console.log('Learning points response:', this.learningPointsData)
      
      // Process learning points into categories
      this.processLearningPoints()
      
      // Store typical skills
      this.typicalSkills = this.learningPointsData.typical_skills || []
      
      console.log('Processed categories:', this.learningCategories)
      console.log('Typical skills:', this.typicalSkills)
    },
    
    processLearningPoints() {
      const learningPoints = this.learningPointsData.learning_points || {}
      
      this.learningCategories = Object.entries(learningPoints).map(([categoryKey, points]) => ({
        name: categoryKey.toLowerCase(),
        displayName: this.getCategoryDisplayName(categoryKey),
        totalPoints: points,
        remainingPoints: points
      }))
    },
    
    getCategoryDisplayName(categoryKey) {
      const displayNames = {
        'Alltag': 'Alltag',
        'Kampf': 'Kampf',
        'Körper': 'Körper',
        'Gesellschaft': 'Gesellschaft',
        'Sozial': 'Sozial',
        'Natur': 'Natur',
        'Wissen': 'Wissen',
        'Handwerk': 'Handwerk',
        'Gaben': 'Gaben',
        'Halbwelt': 'Halbwelt',
        'Unterwelt': 'Unterwelt',
        'Freiland': 'Freiland'
      }
      return displayNames[categoryKey] || categoryKey
    },
    
    restoreSelectedSkills() {
      if (this.sessionData.skills && Array.isArray(this.sessionData.skills)) {
        this.selectedSkills = [...this.sessionData.skills]
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
          category.remainingPoints = Math.max(0, category.remainingPoints - skill.cost)
        }
      })
    },
    
    async selectCategory(categoryName) {
      this.selectedCategory = categoryName
      await this.loadSkillsForCategory(categoryName)
    },
    
    async loadSkillsForCategory(categoryName) {
      try {
        this.isLoadingSkills = true
        
        // Create request for skills in this category
        const request = {
          char_id: 0, // New character - send as number, not string
          name: '', // Will be set for each skill individually
          current_level: 0,
          target_level: 1,
          type: 'skill',
          action: 'learn',
          use_pp: 0,
          use_gold: 0,
          reward: 'default',
          characterClass: this.characterClass,
          category: categoryName
        }
        
        console.log('Loading skills for category:', categoryName, 'with request:', request)
        
        const response = await API.post('/api/characters/available-skills-new', request, {
          headers: { 
            Authorization: `Bearer ${localStorage.getItem('token')}` 
          }
        })
        
        this.availableSkills = (response.data.skills || []).map(skill => ({
          ...skill,
          category: categoryName,
          categoryDisplay: this.getCategoryDisplayName(categoryName)
        }))
        
        console.log('Loaded skills for category:', this.availableSkills)
        
      } catch (error) {
        console.error('Error loading skills:', error)
        // Provide fallback dummy data for development
        this.availableSkills = [
          { 
            name: 'Beispiel Fertigkeit 1', 
            cost: 30, 
            category: categoryName,
            categoryDisplay: this.getCategoryDisplayName(categoryName),
            attribute: 'In'
          },
          { 
            name: 'Beispiel Fertigkeit 2', 
            cost: 40, 
            category: categoryName,
            categoryDisplay: this.getCategoryDisplayName(categoryName),
            attribute: 'Gs'
          }
        ]
      } finally {
        this.isLoadingSkills = false
      }
    },
    
    getSelectedCategoryName() {
      const category = this.learningCategories.find(cat => cat.name === this.selectedCategory)
      return category ? category.displayName : this.selectedCategory
    },
    
    canAddSkill(skill) {
      // Check if skill is already selected
      const alreadySelected = this.selectedSkills.some(s => s.name === skill.name)
      if (alreadySelected) return false
      
      // Check if category has enough points
      const category = this.learningCategories.find(cat => 
        cat.name === skill.category?.toLowerCase()
      )
      if (!category) return false
      
      return category.remainingPoints >= (skill.cost || 0)
    },
    
    addSkill(skill) {
      if (!this.canAddSkill(skill)) return
      
      this.selectedSkills.push({ ...skill })
      this.updateRemainingPoints()
      
      console.log('Added skill:', skill.name, 'Selected skills:', this.selectedSkills)
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
        console.log('No character class available for loading skills')
        return
      }

      this.isLoadingSkills = true
      try {
        // Create request similar to SkillLearnDialog
        const requestData = {
          char_id: 0, // New character - send as number, not string
          name: '', // Will be set for each skill individually
          current_level: 0,
          target_level: 1,
          type: 'skill',
          action: 'learn',
          use_pp: 0,
          use_gold: 0,
          reward: 'default'
        }
        
        console.log('Loading skills with request:', requestData)
        
        const response = await API.post('/api/characters/available-skills-new', requestData)
        
        if (response.data && response.data.skills_by_category) {
          this.availableSkillsByCategory = response.data.skills_by_category
          console.log('Loaded skills by category:', this.availableSkillsByCategory)
        } else {
          // Fallback: Generate sample skills for development
          this.generateSampleSkills()
        }
        
      } catch (error) {
        console.error('Error loading skills:', error)
        // Fallback: Generate sample skills
        this.generateSampleSkills()
      } finally {
        this.isLoadingSkills = false
      }
    },

    generateSampleSkills() {
      // Fallback for testing
      this.availableSkillsByCategory = {
        'Körperliche Fertigkeiten': [
          { name: 'Klettern', epCost: 100, goldCost: 50 },
          { name: 'Schwimmen', epCost: 80, goldCost: 40 },
          { name: 'Springen', epCost: 60, goldCost: 30 }
        ],
        'Geistige Fertigkeiten': [
          { name: 'Erste Hilfe', epCost: 120, goldCost: 60 },
          { name: 'Naturkunde', epCost: 150, goldCost: 75 },
          { name: 'Menschenkenntnis', epCost: 130, goldCost: 65 }
        ],
        'Handwerkliche Fertigkeiten': [
          { name: 'Bogenbau', epCost: 200, goldCost: 100 },
          { name: 'Schmieden', epCost: 250, goldCost: 125 }
        ]
      }
      
      console.log('Generated sample skills:', this.availableSkillsByCategory)
    },

    setSkillCategoryFilter(categoryName) {
      this.selectedSkillCategoryFilter = categoryName
      console.log('Skill category filter set to:', categoryName)
    },

    canAffordSkill(skill) {
      // For character creation, we don't have actual EP/Gold yet
      // This is more of a placeholder for the UI
      return true
    },

    isSkillSelected(skill) {
      return this.selectedSkills.some(s => s.name === skill.name)
    },

    selectSkillForLearning(skill) {
      if (this.isSkillSelected(skill)) {
        return // Already selected
      }
      
      // Add skill to selected list
      this.selectedSkills.push({ ...skill })
      console.log('Skill selected for learning:', skill.name)
    },

    removeSkillFromSelection(skill) {
      const index = this.selectedSkills.findIndex(s => s.name === skill.name)
      if (index !== -1) {
        this.selectedSkills.splice(index, 1)
        console.log('Skill removed from selection:', skill.name)
      }
    }
  }
}
</script>

<style scoped>
/* Minimal custom styles - most styling comes from main.css */

/* Utility classes for dynamic styling that can't be expressed in main.css */
.opacity-50 {
  opacity: 0.6;
}

/* Border color variants for category states */
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
</style>
