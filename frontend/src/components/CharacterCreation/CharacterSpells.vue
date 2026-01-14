<template>
  <div class="fullwidth-page">
    <div class="page-header">
      <h2>{{ $t('characters.create.spells.title') }}</h2>
    </div>
    <p style="color: #666; margin-bottom: 30px; font-size: 16px; line-height: 1.5;">
      {{ $t('characters.create.spells.description') }}
    </p>
    
    <!-- Spell Points Display -->
    <div v-if="spellPoints.total > 0" class="spell-points-display">
      <div class="points-header">
        <h4>Zauber-Lernpunkte</h4>
      </div>
      <div class="points-info">
        <span class="points-remaining">{{ spellPoints.remaining }}</span>
        <span class="points-separator">/</span>
        <span class="points-total">{{ spellPoints.total }}</span>
        <span class="points-label">verfügbare LE</span>
      </div>
      <div class="points-usage" v-if="totalSelectedLE > 0">
        <span class="usage-label">Verwendet:</span>
        <span class="usage-value">{{ totalSelectedLE }} LE</span>
      </div>
    </div>
    
    <!-- Loading State -->
    <div v-if="isLoading" class="loading-message">
      <div class="loading-spinner"></div>
      <p>Lade Zauber für {{ characterClass }}...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-message">
      {{ error }}
      <br>
      <button @click="retry" class="btn btn-primary" style="margin-top: 10px;">Erneut versuchen</button>
    </div>

    <!-- Main Content -->
    <div v-else class="spells-content">
      <!-- Three Column Layout -->
      <div class="three-column-grid">
        
        <!-- Left Column: Available Spells for Selected Category -->
        <div class="card">
          <div class="section-header">
            <h4>Verfügbare Zauber</h4>
            <!--<span v-if="selectedCategory" class="category-badge">{{ getSelectedCategoryName() }}</span>-->
          </div>
          
          <div v-if="!selectedCategory" class="no-selection-message">
            <p>Wählen Sie eine Kategorie aus den Zauberarten rechts, um verfügbare Zauber zu sehen.</p>
          </div>
          
          <div v-else-if="isLoadingSpells" class="loading-message">
            <div class="loading-spinner"></div>
            <p>Lade Zauber...</p>
          </div>
          
          <div v-else class="spells-list">
            <div v-if="availableSpellsForSelectedCategory.length === 0" class="no-spells-message">
              <p>Keine Zauber in dieser Kategorie verfügbar.</p>
            </div>
            
            <div v-else>
              <div
                v-for="spell in availableSpellsForSelectedCategory"
                :key="spell.name"
                class="spell-item"
                :class="{ 
                  'can-select': canAffordSpellInCategory(spell), 
                  'cannot-select': !canAffordSpellInCategory(spell) 
                }"
                @click="selectSpellForLearning(spell)"
              >
                <div class="spell-header">
                  <span class="spell-name">{{ spell.name }}</span>
                  <span class="spell-level">Stufe {{ spell.level }}</span>
                </div>
                <div class="spell-details">
                  <span class="spell-school">{{ spell.school }}</span>
                  <span class="spell-cost">{{ getSpellCost(spell) }} LE</span>
                </div>
                <div v-if="spell.description" class="spell-description">
                  {{ spell.description }}
                </div>
                <div v-if="canAffordSpellInCategory(spell)" class="select-icon">
                  <i class="fas fa-plus"></i>
                </div>
                <div v-else class="cannot-select-reason">
                  <span v-if="isSpellSelected(spell)">Bereits gewählt</span>
                  <span v-else>Nicht genug LE</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Center Column: Selected Spells -->
        <div class="card">
          <div class="section-header">
            <h4>Gewählte Zauber</h4>
            <!--<span class="count-badge">{{ selectedSpells.length }}</span>-->
          </div>
          
          <div v-if="selectedSpells.length === 0" class="no-selection-message">
            <p>Noch keine Zauber gewählt.</p>
          </div>
          
          <div v-else class="selected-spells-list">
            <div
              v-for="spell in selectedSpells"
              :key="spell.name"
              class="selected-spell-item"
            >
              <div class="spell-header-with-remove">
                <div class="spell-header">
                  <span class="spell-name">{{ spell.name }}</span>
                  <span class="spell-level">Grad {{ spell.level }}</span>
                </div>
                <button
                  @click="removeSpellFromSelection(spell)"
                  class="remove-button"
                  :title="'Zauber entfernen'"
                >
                  <i class="fas fa-times"></i>
                </button>
              </div>
              <div class="spell-details">
                <span class="spell-school">{{ spell.school }}</span>
                <span class="spell-cost">{{ getSpellCost(spell) }} LE</span>
                <span class="spell-category">{{ spell.category }}</span>
              </div>
            </div>
          </div>
          
          <!-- Total costs display -->
          <div v-if="selectedSpells.length > 0" class="total-costs">
            <div class="cost-item">
              <span class="cost-label">Gesamt LE:</span>
              <span class="cost-value">{{ totalSelectedLE }}</span>
            </div>
          </div>
        </div>

        <!-- Right Column: Spell Categories -->
        <div class="card">
          <div class="section-header">
            <h4>Zauberarten</h4>
            <!--<span class="info-badge">Verfügbare Kategorien</span>-->
          </div>
          
          <div v-if="spellCategories.length === 0" class="no-categories-message">
            <p>Keine Zauberarten für diese Charakterklasse verfügbar.</p>
          </div>
          
          <div v-else class="categories-list">
            <div
              v-for="category in spellCategories"
              :key="category.name"
              class="category-item"
              :class="{ 
                'selected': selectedCategory === category.name,
                'has-spells': category.spellCount > 0,
                'no-spells': category.spellCount === 0
              }"
              @click="selectCategory(category.name)"
            >
              <div class="category-header">
                <span class="category-name">{{ category.displayName }}</span>
                <span class="spell-count">{{ category.spellCount }} Zauber</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Navigation -->
    <div class="form-row" style="justify-content: space-between; padding-top: 20px; border-top: 1px solid #dee2e6; margin-top: 30px;">
      <button type="button" @click="handlePrevious" class="btn btn-secondary">
        ← Zurück: Fertigkeiten
      </button>
      <button 
        type="button" 
        @click="handleFinalize" 
        class="btn btn-primary"
      >
        Charakter erstellen
        <i class="fas fa-check"></i>
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
import API from '@/utils/api'

export default {
  name: 'CharacterSpells',
  
  props: {
    sessionData: {
      type: Object,
      required: true
    }
  },
  
  emits: ['previous', 'finalize', 'save'],
  
  data() {
    return {
      // Loading states
      isLoading: false,
      isLoadingSpells: false,
      error: null,
      
      // Learning points data
      learningPointsData: null,
      
      // Spells data
      availableSpellsByCategory: null,
      spellCategories: [],
      
      // Spells selection
      selectedCategory: null,
      selectedSpells: [],
      
      // Spell points tracking
      spellPoints: {
        total: 0,
        remaining: 0
      },
      
      // Debug
      showDebug: false // Set to true for debugging
    }
  },
  
  computed: {
    // Get character class from session data
    characterClass() {
      return this.sessionData?.typ || ''
    },
    
    // Available spells for the selected category
    availableSpellsForSelectedCategory() {
      if (!this.selectedCategory || !this.availableSpellsByCategory) {
        return []
      }
      
      // Try to find the matching category key
      const categoryKey = this.findCategoryKey(this.selectedCategory)
      if (!categoryKey) {
        return []
      }
      
      // Get spells from the selected category
      const categorySpells = this.availableSpellsByCategory[categoryKey] || []
      
      const filteredSpells = categorySpells.map(spell => ({
        ...spell,
        cost: this.getSpellCost(spell),
        category: categoryKey // Use the actual category key from availableSpellsByCategory
      }))
      .filter(spell => {
        // Remove already selected spells
        const selectedSpellNames = this.selectedSpells.map(s => s.name)
        return !selectedSpellNames.includes(spell.name)
      })
      .sort((a, b) => {
        const aLeCost = this.getSpellCost(a)
        const bLeCost = this.getSpellCost(b)
        
        // First sort by LE cost (ascending)
        if (aLeCost !== bLeCost) {
          return aLeCost - bLeCost
        }
        
        // If costs are equal, sort by level, then alphabetically
        if (a.level !== b.level) {
          return a.level - b.level
        }
        
        return a.name.localeCompare(b.name)
      })
      
      return filteredSpells
    },

    totalSelectedLE() {
      return this.selectedSpells.reduce((total, spell) => total + this.getSpellCost(spell), 0)
    },

    debugInfo() {
      return {
        // Session Data
        hasSessionData: !!this.sessionData,
        sessionDataKeys: this.sessionData ? Object.keys(this.sessionData) : null,
        
        // Character Info
        characterClass: this.characterClass,
        
        // Loading States
        isLoading: this.isLoading,
        isLoadingSpells: this.isLoadingSpells,
        error: this.error,
        
        // Data States
        hasAvailableSpells: !!this.availableSpellsByCategory,
        spellCategoriesCount: this.spellCategories.length,
        selectedSpellsCount: this.selectedSpells.length,
        selectedCategory: this.selectedCategory,
        totalSelectedLE: this.totalSelectedLE,
        
        // Raw Data (for debugging)
        availableSpellsByCategory: this.availableSpellsByCategory,
        spellCategories: this.spellCategories
      }
    }
  },
  
  watch: {
    selectedSpells: {
      handler(newSpells) {
        // Save spells automatically when they change
        this.$emit('save', { 
          spells: newSpells,
          spell_points: {
            total: this.spellPoints.total,
            remaining: this.spellPoints.remaining
          },
          spells_meta: {
            selectedCategory: this.selectedCategory,
            totalLE: this.totalSelectedLE
          }
        })
      },
      deep: true
    }
  },
  
  created() {
    // Initialize with existing session data if available
    if (this.sessionData.spells && Array.isArray(this.sessionData.spells)) {
      this.selectedSpells = [...this.sessionData.spells]
    }
    
    // Restore selected category if available
    if (this.sessionData.spells_meta?.selectedCategory) {
      this.selectedCategory = this.sessionData.spells_meta.selectedCategory
    }
    
    // Initialize component
    this.initializeComponent()
  },
  
  methods: {
    async initializeComponent() {
      try {
        this.isLoading = true
        this.error = null
        
        // Load learning points data first
        await this.loadLearningPoints()
        
        // Load available spells
        await this.loadAvailableSpells()
        
        // Restore previously selected spells if any
        this.restoreSelectedSpells()
        
        // Update spell points based on selected spells
        this.updateSpellPoints()
        
      } catch (error) {
        console.error('Error initializing component:', error)
        this.error = 'Fehler beim Laden der Zauber. Bitte versuchen Sie es erneut.'
      } finally {
        this.isLoading = false
      }
    },
    
    async loadLearningPoints() {
      if (!this.characterClass) {
        throw new Error('Charakterklasse nicht verfügbar')
      }
      
      try {
        const params = {
          class: this.characterClass
        }
        
        // Add stand if available
        if (this.sessionData.stand) {
          params.stand = this.sessionData.stand
        }
        
        const response = await API.get('/api/characters/classes/learning-points', { params })
        this.learningPointsData = response.data
        
        // Initialize spell points from learning points data
        const spellPoints = this.learningPointsData.spell_points || 0
        this.spellPoints = {
          total: spellPoints,
          remaining: spellPoints
        }
        
      } catch (error) {
        console.error('Error loading learning points:', error)
        // Set default spell points if API fails
        this.spellPoints = {
          total: 10, // Default fallback
          remaining: 10
        }
      }
    },
    
    async loadAvailableSpells() {
      if (!this.characterClass) {
        throw new Error('Charakterklasse nicht verfügbar')
      }

      this.isLoadingSpells = true
      try {
        const token = localStorage.getItem('token')
        
        const request = {
          characterClass: this.characterClass
        }
        
        const response = await API.post('/api/characters/available-spells-creation', request, {
          headers: { Authorization: `Bearer ${token}` },
        })
        
        this.availableSpellsByCategory = response.data.spells_by_category || {}
        this.processSpellCategories()
        
      } catch (error) {
        console.error('Error loading spells:', error)
        this.generateSampleSpells()
      } finally {
        this.isLoadingSpells = false
      }
    },

    processSpellCategories() {
      // Convert spell categories from the backend response
      this.spellCategories = Object.entries(this.availableSpellsByCategory || {}).map(([categoryKey, spells]) => ({
        name: categoryKey.toLowerCase(),
        displayName: categoryKey,
        spellCount: spells.length
      }))
      
      // Sort categories by name
      this.spellCategories.sort((a, b) => a.displayName.localeCompare(b.displayName))
    },

    generateSampleSpells() {
      // Fallback for testing
      this.availableSpellsByCategory = {
        'Dweomer': [
          { name: 'Licht', level: 1, school: 'Dweomer', le_cost: 1, description: 'Erzeugt ein helles Licht' },
          { name: 'Zauberschutz', level: 1, school: 'Dweomer', le_cost: 1, description: 'Schutz vor Magie' }
        ],
        'Erkennen': [
          { name: 'Zaubersicht', level: 1, school: 'Erkennen', le_cost: 1, description: 'Erkennt magische Auren' },
          { name: 'Gedankenlesen', level: 2, school: 'Erkennen', le_cost: 2, description: 'Liest Oberflächengedanken' }
        ],
        'Verändern': [
          { name: 'Heilen von Wunden', level: 1, school: 'Verändern', le_cost: 1, description: 'Heilt leichte Verletzungen' },
          { name: 'Stärke', level: 2, school: 'Verändern', le_cost: 2, description: 'Erhöht die körperliche Stärke' }
        ]
      }
      
      this.processSpellCategories()
      console.log('Generated sample spells:', this.availableSpellsByCategory)
    },
    
    restoreSelectedSpells() {
      if (this.sessionData.spells && Array.isArray(this.sessionData.spells)) {
        this.selectedSpells = [...this.sessionData.spells]
      }
    },
    
    async selectCategory(categoryName) {
      this.selectedCategory = categoryName
      // Spells are now provided by computed property availableSpellsForSelectedCategory
    },
    
    getSelectedCategoryName() {
      const category = this.spellCategories.find(cat => cat.name === this.selectedCategory)
      return category ? category.displayName : this.selectedCategory
    },

    findCategoryKey(selectedCategoryName) {
      if (!this.availableSpellsByCategory) return null
      
      // Try to find by spell category mapping first (most likely scenario)
      const spellCategory = this.spellCategories?.find(cat => cat.name === selectedCategoryName)
      if (spellCategory && this.availableSpellsByCategory[spellCategory.displayName]) {
        return spellCategory.displayName
      }
      
      // Try direct match
      if (this.availableSpellsByCategory[selectedCategoryName]) {
        return selectedCategoryName
      }
      
      // Try case-insensitive search
      const availableKeys = Object.keys(this.availableSpellsByCategory)
      const foundKey = availableKeys.find(key => 
        key.toLowerCase() === selectedCategoryName.toLowerCase()
      )
      
      return foundKey || null
    },

    getSpellCost(spell) {
      // Unified method to get spell cost from various possible properties
      return spell.cost || spell.le_cost || spell.leCost || 0
    },

    canAffordSpellInCategory(spell) {
      // Check if already selected
      if (this.isSpellSelected(spell)) {
        return false
      }
      
      // Check if enough spell points remaining
      const spellCost = this.getSpellCost(spell)
      return this.spellPoints.remaining >= spellCost
    },

    isSpellSelected(spell) {
      return this.selectedSpells.some(s => s.name === spell.name)
    },

    updateSpellPoints() {
      // Reset to total points
      this.spellPoints.remaining = this.spellPoints.total
      
      // Deduct points for selected spells
      this.selectedSpells.forEach(spell => {
        this.spellPoints.remaining -= this.getSpellCost(spell)
      })
      
      // Ensure remaining points don't go below 0
      this.spellPoints.remaining = Math.max(0, this.spellPoints.remaining)
    },

    restoreSelectedSpells() {
      if (this.sessionData.spells && Array.isArray(this.sessionData.spells)) {
        this.selectedSpells = [...this.sessionData.spells]
      }
      
      // Also restore spell points if saved
      if (this.sessionData.spell_points) {
        this.spellPoints.remaining = this.sessionData.spell_points.remaining || this.spellPoints.total
      }
    },

    selectSpellForLearning(spell) {
      if (this.isSpellSelected(spell)) {
        return
      }
      
      // Check if the spell can be afforded
      if (!this.canAffordSpellInCategory(spell)) {
        return
      }
      
      // Add spell to selected list with proper cost and category
      const spellToAdd = {
        ...spell,
        cost: this.getSpellCost(spell)
      }
      
      this.selectedSpells.push(spellToAdd)
      
      // Update spell points
      this.updateSpellPoints()
      
      console.log('Spell selected for learning:', spell.name, 'Cost:', spellToAdd.cost)
    },

    removeSpellFromSelection(spell) {
      const index = this.selectedSpells.findIndex(s => s.name === spell.name)
      if (index !== -1) {
        this.selectedSpells.splice(index, 1)
        
        // Update spell points
        this.updateSpellPoints()
        
        console.log('Spell removed from selection:', spell.name)
      }
    },
    
    async retry() {
      await this.initializeComponent()
    },
    
    handlePrevious() {
      this.$emit('previous')
    },
    
    handleFinalize() {
      const data = {
        spells: this.selectedSpells,
        spell_points: {
          total: this.spellPoints.total,
          remaining: this.spellPoints.remaining
        },
        spells_meta: {
          selectedCategory: this.selectedCategory,
          totalLE: this.totalSelectedLE
        }
      }
      
      console.log('Finalizing with data:', data)
      this.$emit('finalize', data)
    }
  }
}
</script>

<style>
/* All common styles moved to main.css */

.spell-points-display {
  margin: 0 20px 30px 20px;
  padding: 16px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border: 1px solid #dee2e6;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}

.points-header h4 {
  margin: 0 0 12px 0;
  color: #495057;
  font-size: 1.1rem;
  font-weight: 600;
}

.points-info {
  display: flex;
  align-items: baseline;
  gap: 6px;
  margin-bottom: 8px;
}

.points-remaining {
  font-size: 1.8rem;
  font-weight: 700;
  color: #007bff;
}

.points-separator {
  font-size: 1.4rem;
  color: #6c757d;
  margin: 0 2px;
}

.points-total {
  font-size: 1.4rem;
  font-weight: 600;
  color: #6c757d;
}

.points-label {
  font-size: 0.9rem;
  color: #6c757d;
  margin-left: 4px;
}

.points-usage {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9rem;
}

.usage-label {
  color: #6c757d;
}

.usage-value {
  color: #dc3545;
  font-weight: 600;
}

.fullwidth-page {
  /* 
  padding: 0 !important;
  margin: 0 !important;
  width: 100vw !important;
  max-width: 100vw !important;
  box-sizing: border-box !important;
  */
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

.spells-content {
  width: 100%;
  max-width: 100%;
}

.spell-item {
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.spell-item.can-select {
  border-color: #28a745;
  background-color: #f8fff9;
}

.spell-item.can-select:hover {
  border-color: #20c997;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.spell-item.cannot-select {
  border-color: #dc3545;
  background-color: #fff5f5;
  cursor: not-allowed;
  opacity: 0.7;
}

.spell-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.spell-name {
  font-weight: 600;
  color: #333;
}

.spell-level {
  font-size: 0.85em;
  color: #666;
  background: #f8f9fa;
  padding: 2px 6px;
  border-radius: 4px;
}

.spell-details {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.9em;
  color: #666;
  margin-bottom: 6px;
}

.spell-school {
  color: #007bff;
  font-weight: 500;
}

.spell-cost {
  color: #28a745;
  font-weight: 600;
}

.spell-description {
  font-size: 0.85em;
  color: #666;
  font-style: italic;
  margin-top: 6px;
  line-height: 1.3;
}

.select-icon {
  position: absolute;
  top: 8px;
  right: 8px;
  color: #28a745;
  font-size: 1.1em;
}

.cannot-select-reason {
  position: absolute;
  top: 8px;
  right: 8px;
  font-size: 0.8em;
  color: #dc3545;
  text-align: right;
}

.selected-spell-item {
  padding: 12px;
  border: 1px solid #007bff;
  border-radius: 6px;
  margin-bottom: 8px;
  background: #f8fcff;
}

.spell-header-with-remove {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 6px;
}

.selected-spell-item .spell-category {
  color: #6c757d;
  font-size: 0.8em;
  background: #e9ecef;
  padding: 1px 4px;
  border-radius: 3px;
}

.remove-button {
  background: #dc3545;
  border: none;
  color: white;
  cursor: pointer;
  padding: 6px;
  border-radius: 4px;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s ease;
  flex-shrink: 0;
}

.remove-button:hover {
  background-color: #c82333;
}

.category-item {
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.category-item:hover {
  border-color: #007bff;
  background-color: #f8fcff;
}

.category-item.selected {
  border-color: #007bff;
  background-color: #e3f2fd;
}

.category-item.no-spells {
  opacity: 0.6;
  cursor: not-allowed;
}

.category-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.category-name {
  font-weight: 600;
  color: #333;
}

.spell-count {
  font-size: 0.85em;
  color: #666;
  background: #f8f9fa;
  padding: 2px 6px;
  border-radius: 4px;
}

.total-costs {
  border-top: 1px solid #dee2e6;
  padding-top: 12px;
  margin-top: 12px;
}

.cost-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.cost-label {
  font-weight: 500;
  color: #495057;
}

.cost-value {
  font-weight: 600;
  color: #007bff;
}

.opacity-50 {
  opacity: 0.6;
}

.no-selection-message,
.no-spells-message,
.no-categories-message {
  text-align: center;
  color: #6c757d;
  font-style: italic;
  padding: 20px;
}

.category-badge,
.count-badge,
.info-badge {
  font-size: 0.8em;
  padding: 2px 6px;
  border-radius: 4px;
  background: #e9ecef;
  color: #495057;
}

.count-badge {
  background: #007bff;
  color: white;
}

.info-badge {
  background: #17a2b8;
  color: white;
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
