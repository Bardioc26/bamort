<template>
  <div v-if="isVisible" class="modal-overlay" @click.self="closeDialog">
    <div class="modal-content skill-learn-dialog">
      <div class="dialog-header">
        <h3>Neue Fertigkeit lernen</h3>
        <button @click="closeDialog" class="btn-close">&times;</button>
      </div>

      <!-- Ressourcen-Anzeige -->
      <div class="resources-section">
        <h4>Verfügbare Ressourcen</h4>
        <div class="current-resources">
          <div class="resource-display-card">
            <span class="resource-icon">⚡</span>
            <div class="resource-info">
              <div class="resource-label">Erfahrungspunkte</div>
              <div class="resource-amount">{{ character.erfahrungsschatz?.ep || 0 }} EP</div>
              <div class="resource-remaining" v-if="totalSelectedEP > 0">
                <small>
                  Verwendet: {{ totalSelectedEP }} EP | 
                  Verbleibend: {{ remainingEP }} EP
                </small>
              </div>
              <div class="resource-remaining" v-if="totalSelectedEP > 0">
                <small :class="{ 'text-warning': remainingEP < 50, 'text-danger': remainingEP <= 0 }">
                  Nach Lernen: {{ remainingEP }} EP
                </small>
              </div>
            </div>
          </div>
          <div class="resource-display-card">
            <span class="resource-icon">💰</span>
            <div class="resource-info">
              <div class="resource-label">Gold</div>
              <div class="resource-amount">{{ character.vermoegen?.goldstücke || 0 }} GS</div>
              <div class="resource-remaining" v-if="totalSelectedGold > 0 && rewardType === 'default'">
                <small>
                  Verwendet: {{ totalSelectedGold }} GS | 
                  Verbleibend: {{ remainingGold }} GS
                </small>
              </div>
              <div class="resource-remaining" v-if="totalSelectedGold > 0 && rewardType === 'default'">
                <small :class="{ 'text-warning': remainingGold < 20, 'text-danger': remainingGold <= 0 }">
                  Nach Lernen: {{ remainingGold }} GS
                </small>
              </div>
              <div class="resource-remaining" v-if="rewardType === 'noGold' && totalSelectedEP > 0">
                <small class="text-info">
                  Kein Gold benötigt (Belohnungslernen)
                </small>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Lernmethode direkt unter den Ressourcen -->
        <div class="reward-method-section">
          <label for="rewardType">Lernen als Belohnung:</label>
          <select id="rewardType" v-model="rewardType" class="form-select">
            <option value="default">Standard (EP + Gold)</option>
            <option value="noGold">Nur EP (kein Gold)</option>
          </select>
          <small class="form-hint">Wählen Sie die Art des Lernens</small>
        </div>
      </div>

      <!-- Formular -->
      <div class="form-section">
        <!-- Fertigkeiten-Auswahl mit Drag & Drop -->
        <div class="skills-selection-container">
          <div class="skills-available">
            <h4>Verfügbare Fertigkeiten</h4>
            
            <!-- Kategorie-Filter -->
            <div class="category-filters">
              <button 
                @click="setCategoryFilter(null)"
                class="category-filter-btn"
                :class="{ 'active': selectedCategoryFilter === null }"
                title="Alle Kategorien anzeigen"
              >
                Alle
              </button>
              <button 
                v-for="category in availableCategories" 
                :key="category.key || category.name"
                @click="setCategoryFilter(category.name)"
                class="category-filter-btn"
                :class="{ 'active': selectedCategoryFilter === category.name }"
                :title="category.description"
              >
                {{ category.name }}
              </button>
            </div>
            
            <!-- Sortier- und Suchbereich -->
            <div class="sort-and-search-controls">
              <div class="sort-controls">
                <span class="sort-label">Sortieren nach:</span>
                <button 
                  @click="setSortBy('name')"
                  class="sort-btn"
                  :class="{ 'active': sortBy === 'name' }"
                  title="Nach Name sortieren"
                >
                  Name 
                  <span v-if="sortBy === 'name'" class="sort-icon">
                    {{ sortOrder === 'asc' ? '↑' : '↓' }}
                  </span>
                </button>
                <button 
                  @click="setSortBy('epCost')"
                  class="sort-btn"
                  :class="{ 'active': sortBy === 'epCost' }"
                  title="Nach EP-Kosten sortieren"
                >
                  EP-Kosten 
                  <span v-if="sortBy === 'epCost'" class="sort-icon">
                    {{ sortOrder === 'asc' ? '↑' : '↓' }}
                  </span>
                </button>
              </div>
              
              <div class="skills-search">
                <input 
                  v-model="skillSearchFilter"
                  type="text" 
                  placeholder="Fertigkeiten filtern..." 
                  class="form-input search-input"
                />
              </div>
            </div>
            <div class="skills-list" v-if="availableSkillsByCategory">
              <div 
                v-for="skill in sortedFilteredSkills" 
                :key="skill.name"
                class="skill-item"
                :class="{ 'skill-affordable': skill.canAfford }"
                draggable="true"
                @dragstart="onDragStart($event, skill)"
                @click="selectSkill(skill)"
              >
                <div class="skill-info">
                  <div class="skill-main-line">
                    <span class="skill-name">{{ skill.name }}</span>
                    <span class="skill-category">({{ skill.category }})</span>
                    <span class="skill-costs">
                      <span v-if="rewardType === 'default'" class="cost-ep">{{ skill.epCost }} EP</span>
                      <span v-if="rewardType === 'default'" class="cost-gold">{{ skill.goldCost }} GS</span>
                      <span v-if="rewardType === 'noGold'" class="cost-ep">{{ skill.epCost }} EP</span>
                    </span>
                  </div>
                </div>
                <div class="skill-actions">
                  <button 
                    @click.stop="selectSkill(skill)"
                    class="btn-select"
                    :disabled="!skill.canAfford"
                  >
                    →
                  </button>
                </div>
              </div>
            </div>
            <div v-if="isLoadingSkills" class="loading-skills">
              <span class="loading-spinner">⏳</span> Lade Fertigkeiten...
            </div>
          </div>
          
          <div class="skills-selected">
            <h4>Zu lernende Fertigkeiten</h4>
            <div 
              class="skills-drop-zone"
              :class="{ 'drag-over': isDragOver }"
              @dragover.prevent="isDragOver = true"
              @dragleave.prevent="isDragOver = false"
              @drop.prevent="onDrop"
            >
              <div v-if="selectedSkills.length === 0" class="drop-zone-placeholder">
                <div class="placeholder-icon">📚</div>
                <div class="placeholder-text">
                  Ziehen Sie Fertigkeiten hierher oder klicken Sie auf → um sie auszuwählen
                </div>
              </div>
              <div v-else class="selected-skills-list">
                <div 
                  v-for="(skill, index) in selectedSkills" 
                  :key="skill.name + index"
                  class="selected-skill-item"
                >
                  <div class="selected-skill-info">
                    <div class="selected-skill-name">{{ skill.name }}</div>
                    <div class="selected-skill-costs">
                      <span v-if="rewardType === 'default'" class="cost-ep">{{ skill.epCost }} EP</span>
                      <span v-if="rewardType === 'default'" class="cost-gold">{{ skill.goldCost }} GS</span>
                      <span v-if="rewardType === 'noGold'" class="cost-ep">{{ skill.epCost }} EP</span>
                    </div>
                  </div>
                  <button 
                    @click="removeSelectedSkill(index)"
                    class="btn-remove"
                    title="Entfernen"
                  >
                    ×
                  </button>
                </div>
              </div>
            </div>
            
            <!-- Gesamtkosten -->
            <div v-if="selectedSkills.length > 0" class="total-costs">
              <div class="total-costs-header">Gesamtkosten:</div>
              <div class="total-costs-breakdown">
                <span v-if="rewardType === 'default'" class="total-ep">{{ totalSelectedEP }} EP</span>
                <span v-if="rewardType === 'default'" class="total-gold">{{ totalSelectedGold }} GS</span>
                <span v-if="rewardType === 'noGold'" class="total-ep">{{ totalSelectedEP }} EP</span>
              </div>
              <div class="affordability-check">
                <span 
                  :class="{ 
                    'text-success': canAffordSelected, 
                    'text-danger': !canAffordSelected 
                  }"
                >
                  {{ canAffordSelected ? '✓ Kann gelernt werden' : '✗ Nicht genügend Ressourcen' }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Einfache Eingabe als Fallback -->
        <div class="simple-input-section" v-if="!availableSkillsByCategory">
          <div class="form-group">
            <label for="skillName">Fertigkeitsname:</label>
            <input 
              id="skillName"
              v-model="skillName" 
              type="text" 
              placeholder="Name der neuen Fertigkeit eingeben..." 
              class="form-input"
              @keyup.enter="learnSkill"
            />
          </div>
        </div>

        <div class="form-group">
          <label for="notes">Notizen (optional):</label>
          <textarea 
            id="notes"
            v-model="notes" 
            placeholder="Zusätzliche Notizen zum Lernen der Fertigkeit..."
            class="form-textarea"
            rows="3"
          ></textarea>
        </div>
      </div>

      <!-- Kosten-Vorschau (falls implementiert) -->
      <div v-if="estimatedCosts" class="costs-preview">
        <h4>Geschätzte Kosten</h4>
        <div class="cost-breakdown">
          <div class="cost-item">
            <span class="cost-label">EP-Kosten:</span>
            <span class="cost-value">{{ estimatedCosts.ep || 0 }} EP</span>
          </div>
          <div v-if="rewardType === 'default'" class="cost-item">
            <span class="cost-label">Gold-Kosten:</span>
            <span class="cost-value">{{ estimatedCosts.gold || 0 }} GS</span>
          </div>
        </div>
      </div>

      <!-- Aktionen -->
      <div class="modal-actions">
        <div class="action-info">
          <span v-if="selectedSkills.length > 0" class="selection-count">
            {{ selectedSkills.length }} Fertigkeit{{ selectedSkills.length !== 1 ? 'en' : '' }} ausgewählt
          </span>
        </div>
        
        <div class="action-buttons">
          <button 
            @click="learnSelectedSkills" 
            class="btn-confirm"
            :disabled="selectedSkills.length === 0 || !canAffordSelected || isLoading"
          >
            <span v-if="isLoading" class="loading-spinner">⏳</span>
            {{ isLoading ? 'Lerne...' : (selectedSkills.length > 1 ? 'Fertigkeiten lernen' : 'Fertigkeit lernen') }}
          </button>
          <button @click="closeDialog" class="btn-cancel" :disabled="isLoading">
            Abbrechen
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import API from '@/utils/api'

export default {
  name: 'SkillLearnDialog',
  props: {
    character: {
      type: Object,
      required: true
    },
    isVisible: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      skillName: '',
      rewardType: 'default',
      notes: '',
      isLoading: false,
      estimatedCosts: null,
      
      // Neue Eigenschaften für Fertigkeiten-Auswahl
      availableSkillsByCategory: null,
      selectedSkills: [],
      skillSearchFilter: '',
      isDragOver: false,
      isLoadingSkills: false,
      
      // Kategorie-Filter
      availableCategories: [],
      selectedCategoryFilter: null,
      
      // Sortierung
      sortBy: 'name', // 'name', 'epCost'
      sortOrder: 'asc' // 'asc', 'desc'
    }
  },
  computed: {
    filteredSkillsByCategory() {
      let filtered = this.availableSkillsByCategory
      
      if (!filtered) {
        return filtered
      }
      
      // Kategorie-Filter anwenden
      if (this.selectedCategoryFilter) {
        filtered = {
          [this.selectedCategoryFilter]: filtered[this.selectedCategoryFilter] || []
        }
      }
      
      // Text-Suchfilter anwenden
      if (this.skillSearchFilter) {
        const result = {}
        const searchTerm = this.skillSearchFilter.toLowerCase()
        
        Object.keys(filtered).forEach(category => {
          const filteredSkills = filtered[category].filter(skill =>
            skill.name.toLowerCase().includes(searchTerm)
          )
          if (filteredSkills.length > 0) {
            result[category] = filteredSkills
          }
        })
        
        return result
      }
      
      return filtered
    },
    
    sortedFilteredSkills() {
      if (!this.availableSkillsByCategory) {
        return []
      }
      
      let allSkills = []
      
      // Sammle alle Fertigkeiten aus allen Kategorien
      Object.keys(this.availableSkillsByCategory).forEach(category => {
        this.availableSkillsByCategory[category].forEach(skill => {
          allSkills.push({
            ...skill,
            category: category,
            canAfford: this.canAffordSkill(skill) // Berechne canAfford im Frontend
          })
        })
      })
      
      // Entferne bereits ausgewählte Fertigkeiten
      const selectedSkillNames = this.selectedSkills.map(s => s.name)
      allSkills = allSkills.filter(skill => 
        !selectedSkillNames.includes(skill.name)
      )
      
      // Anwenden der Filter
      if (this.selectedCategoryFilter) {
        allSkills = allSkills.filter(skill => skill.category === this.selectedCategoryFilter)
      }
      
      if (this.skillSearchFilter) {
        const searchTerm = this.skillSearchFilter.toLowerCase()
        allSkills = allSkills.filter(skill => 
          skill.name.toLowerCase().includes(searchTerm)
        )
      }
      
      // Sortiere nach dem gewählten Kriterium
      if (this.sortBy === 'name') {
        allSkills.sort((a, b) => {
          const comparison = a.name.localeCompare(b.name)
          return this.sortOrder === 'asc' ? comparison : -comparison
        })
      } else if (this.sortBy === 'epCost') {
        allSkills.sort((a, b) => {
          const comparison = (a.epCost || 0) - (b.epCost || 0)
          return this.sortOrder === 'asc' ? comparison : -comparison
        })
      }
      
      return allSkills
    },
    
    totalSelectedEP() {
      return this.selectedSkills.reduce((total, skill) => total + (skill.epCost || 0), 0)
    },
    
    totalSelectedGold() {
      return this.selectedSkills.reduce((total, skill) => total + (skill.goldCost || 0), 0)
    },
    
    canAffordSelected() {
      const availableEP = this.character.erfahrungsschatz?.ep || 0
      const availableGold = this.character.vermoegen?.goldstücke || 0
      
      if (this.rewardType === 'default') {
        return availableEP >= this.totalSelectedEP && availableGold >= this.totalSelectedGold
      } else if (this.rewardType === 'noGold') {
        return availableEP >= this.totalSelectedEP
      }
      return false
    },
    
    remainingEP() {
      const current = this.character.erfahrungsschatz?.ep || 0
      return Math.max(0, current - this.totalSelectedEP)
    },
    
    remainingGold() {
      const current = this.character.vermoegen?.goldstücke || 0
      return Math.max(0, current - this.totalSelectedGold)
    }
  },
  watch: {
    isVisible(newVal) {
      if (newVal) {
        this.loadAvailableSkills()
      } else {
        this.resetForm()
      }
    },
    
    rewardType() {
      // Lade Fertigkeiten neu bei Änderung der Lernmethode
      if (this.isVisible) {
        this.loadAvailableSkills()
      }
    }
  },
  
  mounted() {
    this.$api = API
    if (this.isVisible) {
      this.loadAvailableSkills()
    }
  },
  methods: {
    closeDialog() {
      this.$emit('close')
    },

    canAffordSkill(skill) {
      const availableEP = this.character.erfahrungsschatz?.ep || 0
      const availableGold = this.character.vermoegen?.goldstücke || 0
      
      if (this.rewardType === 'default' || this.rewardType === '') {
        return availableEP >= (skill.epCost || 0) && availableGold >= (skill.goldCost || 0)
      } else if (this.rewardType === 'noGold') {
        return availableEP >= (skill.epCost || 0)
      }
      return false
    },

    resetForm() {
      this.skillName = ''
      this.rewardType = 'default'
      this.notes = ''
      this.isLoading = false
      this.estimatedCosts = null
      this.selectedSkills = []
      this.skillSearchFilter = ''
      this.isDragOver = false
      this.selectedCategoryFilter = null
      this.availableSkillsByCategory = null
      this.sortBy = 'name'
      this.sortOrder = 'asc'
    },

    async loadAvailableSkills() {
      this.isLoadingSkills = true
      try {
        // Lade verfügbare Kategorien
        await this.loadAvailableCategories()
        
        // Lade alle verfügbaren Fertigkeiten mit Kosten (bereits ohne gelernte gefiltert)
        const requestData = {
          char_id: this.character.id,
          name: '', // Wird für jede Fertigkeit einzeln gesetzt
          current_level: 0,
          target_level: 1,
          type: 'skill',
          action: 'learn',
          use_pp: 0,
          use_gold: 0,
          reward: this.rewardType || 'default'
        }
        
        const response = await this.$api.post('/api/characters/available-skills-new', requestData)
        
        if (response.data && response.data.skills_by_category) {
          this.availableSkillsByCategory = response.data.skills_by_category
          console.log('Loaded skills by category:', this.availableSkillsByCategory)
        } else {
          // Fallback: Generiere Beispiel-Fertigkeiten
          this.generateSampleSkills()
        }
        
      } catch (error) {
        console.error('Fehler beim Laden der Fertigkeiten:', error)
        // Fallback: Generiere Beispiel-Fertigkeiten
        this.generateSampleSkills()
      } finally {
        this.isLoadingSkills = false
      }
    },

    async loadAvailableCategories() {
      try {
        const response = await this.$api.get('/api/characters/skill-categories')
        if (response.data && response.data.skill_categories) {
          // Extrahiere die Namen und Beschreibungen der Kategorien
          const categories = response.data.skill_categories
          this.availableCategories = Object.keys(categories).map(key => ({
            name: categories[key].name,
            description: categories[key].description,
            key: key
          }))
          console.log('Loaded categories:', this.availableCategories)
        } else {
          // Fallback: Standard-Kategorien
          this.availableCategories = [
            { name: 'Körperliche Fertigkeiten', description: 'Körperliche Fertigkeiten', key: 'körper' },
            { name: 'Geistige Fertigkeiten', description: 'Geistige Fertigkeiten', key: 'geist' }, 
            { name: 'Handwerkliche Fertigkeiten', description: 'Handwerkliche Fertigkeiten', key: 'handwerk' },
            { name: 'Magische Fertigkeiten', description: 'Magische Fertigkeiten', key: 'magie' },
            { name: 'Soziale Fertigkeiten', description: 'Soziale Fertigkeiten', key: 'sozial' }
          ]
        }
      } catch (error) {
        console.error('Fehler beim Laden der Kategorien:', error)
        // Fallback: Standard-Kategorien
        this.availableCategories = [
          { name: 'Körperliche Fertigkeiten', description: 'Körperliche Fertigkeiten', key: 'körper' },
          { name: 'Geistige Fertigkeiten', description: 'Geistige Fertigkeiten', key: 'geist' }, 
          { name: 'Handwerkliche Fertigkeiten', description: 'Handwerkliche Fertigkeiten', key: 'handwerk' },
          { name: 'Magische Fertigkeiten', description: 'Magische Fertigkeiten', key: 'magie' },
          { name: 'Soziale Fertigkeiten', description: 'Soziale Fertigkeiten', key: 'sozial' }
        ]
      }
    },

    generateSampleSkills() {
      // Fallback für Testzwecke
      const availableEP = this.character.erfahrungsschatz?.ep || 0
      const availableGold = this.character.vermoegen?.goldstücke || 0
      
      this.availableSkillsByCategory = {
        'Körperliche Fertigkeiten': [
          { name: 'Klettern', epCost: 100, goldCost: 50, canAfford: availableEP >= 100 && availableGold >= 50 },
          { name: 'Schwimmen', epCost: 80, goldCost: 40, canAfford: availableEP >= 80 && availableGold >= 40 },
          { name: 'Springen', epCost: 60, goldCost: 30, canAfford: availableEP >= 60 && availableGold >= 30 }
        ],
        'Geistige Fertigkeiten': [
          { name: 'Erste Hilfe', epCost: 120, goldCost: 60, canAfford: availableEP >= 120 && availableGold >= 60 },
          { name: 'Naturkunde', epCost: 150, goldCost: 75, canAfford: availableEP >= 150 && availableGold >= 75 },
          { name: 'Menschenkenntnis', epCost: 130, goldCost: 65, canAfford: availableEP >= 130 && availableGold >= 65 }
        ],
        'Handwerkliche Fertigkeiten': [
          { name: 'Bogenbau', epCost: 200, goldCost: 100, canAfford: availableEP >= 200 && availableGold >= 100 },
          { name: 'Schmieden', epCost: 250, goldCost: 125, canAfford: availableEP >= 250 && availableGold >= 125 }
        ]
      }
      
      // Setze verfügbare Kategorien aus den Beispieldaten
      this.availableCategories = Object.keys(this.availableSkillsByCategory).map(key => ({
        name: key,
        description: key,
        key: key.toLowerCase().replace(/\s+/g, '_')
      }))
      
      console.log('Generated sample skills:', this.availableSkillsByCategory)
    },

    setCategoryFilter(categoryName) {
      this.selectedCategoryFilter = categoryName
      console.log('Category filter set to:', categoryName)
    },

    setSortBy(sortBy) {
      if (this.sortBy === sortBy) {
        // Gleicher Sortiertyp - Reihenfolge umkehren
        this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc'
      } else {
        // Neuer Sortiertyp - standardmäßig aufsteigend
        this.sortBy = sortBy
        this.sortOrder = 'asc'
      }
      console.log('Sort set to:', this.sortBy, this.sortOrder)
    },

    onDragStart(event, skill) {
      event.dataTransfer.setData('application/json', JSON.stringify(skill))
      event.dataTransfer.effectAllowed = 'copy'
    },

    onDrop(event) {
      this.isDragOver = false
      try {
        const skillData = JSON.parse(event.dataTransfer.getData('application/json'))
        this.selectSkill(skillData)
      } catch (error) {
        console.error('Fehler beim Drag & Drop:', error)
      }
    },

    selectSkill(skill) {
      if (!skill.canAfford) {
        alert('Sie haben nicht genügend Ressourcen für diese Fertigkeit.')
        return
      }
      
      // Prüfe ob Fertigkeit bereits ausgewählt (sollte eigentlich nicht passieren, da gefiltert)
      const alreadySelected = this.selectedSkills.some(s => s.name === skill.name)
      if (alreadySelected) {
        alert('Diese Fertigkeit ist bereits ausgewählt.')
        return
      }
      
      // Füge Fertigkeit zu den ausgewählten hinzu
      this.selectedSkills.push({ ...skill })
      
      console.log('Skill selected:', skill.name)
      console.log('Currently selected skills:', this.selectedSkills.map(s => s.name))
    },

    removeSelectedSkill(index) {
      const removedSkill = this.selectedSkills[index]
      this.selectedSkills.splice(index, 1)
      
      console.log('Skill removed from selection:', removedSkill.name)
      console.log('Currently selected skills:', this.selectedSkills.map(s => s.name))
    },

    async learnSelectedSkills() {
      if (this.selectedSkills.length === 0) {
        alert('Bitte wählen Sie mindestens eine Fertigkeit aus.')
        return
      }

      if (!this.canAffordSelected) {
        alert('Sie haben nicht genügend Ressourcen für die ausgewählten Fertigkeiten.')
        return
      }

      this.isLoading = true
      try {
        // Lerne alle ausgewählten Fertigkeiten
        const learnPromises = this.selectedSkills.map(skill => {
          const requestData = {
            char_id: this.character.id,
            name: skill.name,
            current_level: 0,
            target_level: 1,
            type: 'skill',
            action: 'learn',
            reward: this.rewardType || 'default'  // Immer das reward-Feld setzen
          }

          return this.$api.post(`/api/characters/${this.character.id}/learn-skill`, requestData)
        })

        const responses = await Promise.all(learnPromises)
        
        console.log('Fertigkeiten erfolgreich gelernt:', responses.map(r => r.data))
        
        const skillNames = this.selectedSkills.map(s => s.name).join(', ')
        alert(`Fertigkeiten erfolgreich gelernt: ${skillNames}`)
        
        this.closeDialog()
        this.$emit('skill-learned', {
          skills: this.selectedSkills,
          responses: responses.map(r => r.data)
        })
        
      } catch (error) {
        console.error('Fehler beim Lernen der Fertigkeiten:', error)
        const errorMessage = error.response?.data?.error || error.message || 'Unbekannter Fehler'
        alert('Fehler beim Lernen der Fertigkeiten: ' + errorMessage)
      } finally {
        this.isLoading = false
      }
    },

    // Utility: Debounce function
    debounce(func, wait) {
      let timeout
      return function executedFunction(...args) {
        const later = () => {
          clearTimeout(timeout)
          func(...args)
        }
        clearTimeout(timeout)
        timeout = setTimeout(later, wait)
      }
    }
  }
}
</script>

<style scoped>
/* Component-specific styles - common styles are in main.css */

/* Fullscreen modal override */
.modal-content {
  background: white;
  border-radius: 0;
  width: 100vw;
  height: 100vh;
  max-width: none;
  max-height: none;
  overflow: hidden;
  box-shadow: none;
  animation: modalSlideIn 0.3s ease;
  display: flex;
  flex-direction: column;
  position: absolute;
  top: 0;
  left: 0;
}

@keyframes modalSlideIn {
  from {
    opacity: 0;
    transform: scale(0.9) translateY(-20px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

/* Component-specific overrides */
.resource-display-card {
  background: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.resource-amount {
  color: #1da766;
}

.modal-actions {
  border-radius: 0;
  flex-shrink: 0;
}

/* Button overrides for this specific dialog */
.btn-confirm {
  padding: 12px 24px;
  background: #1da766;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-confirm:hover:not(:disabled) {
  background: #16a085;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(29, 167, 102, 0.3);
}

.btn-confirm:disabled {
  background: #6c757d;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.btn-cancel {
  padding: 12px 24px;
  background: #6c757d;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.2s ease;
}

.btn-cancel:hover:not(:disabled) {
  background: #5a6268;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(108, 117, 125, 0.3);
}

.btn-cancel:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.loading-spinner {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Responsive Design */
@media (max-width: 768px) {
  .modal-content {
    width: 95%;
    margin: 10px;
    max-width: none;
  }
  
  .cost-breakdown {
    flex-direction: column;
    gap: 10px;
  }
  
  .modal-actions {
    flex-direction: column;
    align-items: stretch;
  }
  
  .action-buttons {
    width: 100%;
  }
  
  .btn-confirm,
  .btn-cancel {
    width: 100%;
  }
  
  .total-costs-breakdown {
    flex-direction: column;
    gap: 4px;
  }
  
  .category-filters {
    flex-direction: column;
    gap: 6px;
  }
  
  .sort-and-search-controls {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
  
  .sort-controls {
    justify-content: center;
  }
  
  .skills-search {
    min-width: auto;
  }
  
  .category-filter-btn {
    width: 100%;
    text-align: center;
  }
}
</style>
