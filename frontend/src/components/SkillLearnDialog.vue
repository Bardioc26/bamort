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
        <div class="resources-display">
          <div class="resource-item">
            <span class="resource-icon">⚡</span>
            <span class="resource-label">Erfahrungspunkte:</span>
            <span class="resource-value">{{ character.erfahrungsschatz?.ep || 0 }} EP</span>
          </div>
          <div class="resource-item">
            <span class="resource-icon">💰</span>
            <span class="resource-label">Gold:</span>
            <span class="resource-value">{{ character.vermoegen?.goldstücke || 0 }} GS</span>
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
              >
                Alle
              </button>
              <button 
                v-for="category in availableCategories" 
                :key="category"
                @click="setCategoryFilter(category)"
                class="category-filter-btn"
                :class="{ 'active': selectedCategoryFilter === category }"
              >
                {{ category }}
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
      learnedSkillNames: [], // Namen der bereits gelernten Fertigkeiten
      
      // Kategorie-Filter
      availableCategories: [],
      selectedCategoryFilter: null
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
            category: category
          })
        })
      })
      
      // Entferne bereits gelernte Fertigkeiten
      allSkills = allSkills.filter(skill => 
        !this.learnedSkillNames.includes(skill.name)
      )
      
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
      
      // Sortiere nach Namen
      return allSkills.sort((a, b) => a.name.localeCompare(b.name))
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
      this.learnedSkillNames = []
    },

    async loadAvailableSkills() {
      this.isLoadingSkills = true
      try {
        // Lade verfügbare Kategorien
        await this.loadAvailableCategories()
        
        // Lade bereits gelernte Fertigkeiten des Charakters
        await this.loadLearnedSkills()
        
        // Lade alle verfügbaren Fertigkeiten mit Kosten
        const response = await this.$api.get(`/api/characters/${this.character.id}/available-skills`, {
          params: {
            reward_type: this.rewardType
          }
        })
        
        if (response.data && response.data.skills_by_category) {
          this.availableSkillsByCategory = response.data.skills_by_category
          console.log('Loaded skills by category:', this.availableSkillsByCategory)
          console.log('Learned skills to exclude:', this.learnedSkillNames)
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

    async loadLearnedSkills() {
      try {
        // Lade bereits gelernte Fertigkeiten des Charakters
        const response = await this.$api.get(`/api/characters/${this.character.id}/skills`)
        if (response.data && response.data.skills) {
          // Extrahiere die Namen aller bereits gelernten Fertigkeiten
          this.learnedSkillNames = response.data.skills.map(skill => skill.name)
          console.log('Already learned skills:', this.learnedSkillNames)
        } else if (this.character.fertigkeiten && Array.isArray(this.character.fertigkeiten)) {
          // Fallback: Verwende Fertigkeiten aus dem Character-Objekt
          this.learnedSkillNames = this.character.fertigkeiten.map(skill => skill.name)
          console.log('Already learned skills (from character):', this.learnedSkillNames)
        }
      } catch (error) {
        console.error('Fehler beim Laden der gelernten Fertigkeiten:', error)
        // Fallback: Verwende Fertigkeiten aus dem Character-Objekt falls verfügbar
        if (this.character.fertigkeiten && Array.isArray(this.character.fertigkeiten)) {
          this.learnedSkillNames = this.character.fertigkeiten.map(skill => skill.name)
        }
      }
    },

    async loadAvailableCategories() {
      try {
        const response = await this.$api.get('/skill-categories')
        if (response.data && response.data.categories) {
          this.availableCategories = response.data.categories
        } else {
          // Fallback: Standard-Kategorien
          this.availableCategories = [
            'Körperliche Fertigkeiten',
            'Geistige Fertigkeiten', 
            'Handwerkliche Fertigkeiten',
            'Magische Fertigkeiten',
            'Soziale Fertigkeiten'
          ]
        }
      } catch (error) {
        console.error('Fehler beim Laden der Kategorien:', error)
        // Fallback: Standard-Kategorien
        this.availableCategories = [
          'Körperliche Fertigkeiten',
          'Geistige Fertigkeiten', 
          'Handwerkliche Fertigkeiten',
          'Magische Fertigkeiten',
          'Soziale Fertigkeiten'
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
      this.availableCategories = Object.keys(this.availableSkillsByCategory)
      
      console.log('Generated sample skills:', this.availableSkillsByCategory)
    },

    setCategoryFilter(category) {
      this.selectedCategoryFilter = category
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
            action: 'learn'
          }

          if (this.rewardType !== 'default') {
            requestData.reward = this.rewardType
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
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: block;
  z-index: 1000;
}

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

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 2px solid #1da766;
  background: #f8f9fa;
  border-radius: 0;
  flex-shrink: 0;
  z-index: 10;
}

.dialog-header h3 {
  margin: 0;
  color: #333;
  font-size: 1.5rem;
}

.btn-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
  padding: 0;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.btn-close:hover {
  background: #e9ecef;
  color: #333;
}

.resources-section {
  padding: 20px 24px;
  background: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
  flex-shrink: 0;
}

.resources-section h4 {
  margin: 0 0 15px 0;
  color: #495057;
  font-size: 1.1rem;
}

.resources-display {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.resource-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: white;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.resource-icon {
  font-size: 18px;
}

.resource-label {
  font-weight: 600;
  color: #495057;
}

.resource-value {
  font-weight: bold;
  color: #1da766;
  font-size: 1.1rem;
}

.reward-method-section {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #dee2e6;
}

.reward-method-section label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #495057;
  font-size: 0.95rem;
}

.reward-method-section .form-select {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #dee2e6;
  border-radius: 8px;
  font-size: 14px;
  font-family: inherit;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
  box-sizing: border-box;
  background: white;
}

.reward-method-section .form-select:focus {
  outline: none;
  border-color: #1da766;
  box-shadow: 0 0 0 3px rgba(29, 167, 102, 0.1);
}

.reward-method-section .form-hint {
  display: block;
  margin-top: 4px;
  font-size: 0.85rem;
  color: #6c757d;
  font-style: italic;
}

.form-section {
  padding: 24px;
  flex: 1;
  overflow-y: auto;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #495057;
  font-size: 0.95rem;
}

.form-input,
.form-select,
.form-textarea {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #dee2e6;
  border-radius: 8px;
  font-size: 14px;
  font-family: inherit;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
  box-sizing: border-box;
}

.form-input:focus,
.form-select:focus,
.form-textarea:focus {
  outline: none;
  border-color: #1da766;
  box-shadow: 0 0 0 3px rgba(29, 167, 102, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
}

.form-hint {
  display: block;
  margin-top: 4px;
  font-size: 0.85rem;
  color: #6c757d;
  font-style: italic;
}

.costs-preview {
  padding: 20px 24px;
  background: #fff3cd;
  border-top: 1px solid #ffeaa7;
  border-bottom: 1px solid #ffeaa7;
  flex-shrink: 0;
}

.costs-preview h4 {
  margin: 0 0 12px 0;
  color: #856404;
  font-size: 1rem;
}

.cost-breakdown {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.cost-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.cost-label {
  color: #856404;
  font-weight: 500;
}

.cost-value {
  font-weight: bold;
  color: #495057;
  background: white;
  padding: 4px 8px;
  border-radius: 4px;
  border: 1px solid #ffeaa7;
}

.modal-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  padding: 20px 24px;
  background: #f8f9fa;
  border-radius: 0;
  border-top: 1px solid #dee2e6;
  flex-shrink: 0;
}

.action-info {
  flex: 1;
}

.selection-count {
  font-size: 0.9rem;
  color: #6c757d;
  font-weight: 500;
}

.action-buttons {
  display: flex;
  gap: 12px;
}

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

/* Skills Selection Styles */
.skills-selection-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}

.skills-available,
.skills-selected {
  border: 2px solid #dee2e6;
  border-radius: 8px;
  overflow: hidden;
}

.skills-available h4,
.skills-selected h4 {
  margin: 0;
  padding: 12px 16px;
  background: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
  font-size: 1rem;
  color: #495057;
}

.skills-search {
  padding: 12px 16px;
  border-bottom: 1px solid #dee2e6;
}

.category-filters {
  padding: 12px 16px;
  border-bottom: 1px solid #dee2e6;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  background: #f8f9fa;
}

.category-filter-btn {
  padding: 6px 12px;
  background: white;
  border: 2px solid #dee2e6;
  border-radius: 20px;
  font-size: 0.85rem;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
  font-weight: 500;
  color: #495057;
}

.category-filter-btn:hover {
  border-color: #1da766;
  color: #1da766;
}

.category-filter-btn.active {
  background: #1da766;
  border-color: #1da766;
  color: white;
  font-weight: 600;
}

.category-filter-btn:first-child {
  font-weight: 600;
  background: #e9ecef;
  border-color: #adb5bd;
}

.category-filter-btn:first-child.active {
  background: #495057;
  border-color: #495057;
  color: white;
}

.search-input {
  margin: 0;
  font-size: 13px;
}

.skills-list {
  max-height: 60vh;
  overflow-y: auto;
  background: white;
}

.skill-item {
  padding: 12px 16px;
  border-bottom: 1px solid #f8f9fa;
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: grab;
  transition: all 0.2s ease;
}

.skill-item:hover {
  background: #f8f9fa;
}

.skill-item:active {
  cursor: grabbing;
}

.skill-item.skill-affordable {
  border-left: 4px solid #1da766;
}

.skill-item:not(.skill-affordable) {
  opacity: 0.6;
  cursor: not-allowed;
}

.skill-info {
  flex: 1;
}

.skill-main-line {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.skill-name {
  font-weight: 600;
  color: #333;
  font-size: 0.95rem;
  min-width: 120px;
}

.skill-category {
  font-size: 0.8rem;
  color: #6c757d;
  font-style: italic;
  min-width: 100px;
}

.skill-costs {
  display: flex;
  gap: 8px;
  font-size: 0.85rem;
  margin-left: auto;
}

.cost-ep {
  color: #1da766;
  font-weight: 600;
}

.cost-gold {
  color: #ffc107;
  font-weight: 600;
}

.skill-actions {
  display: flex;
  gap: 8px;
}

.btn-select {
  width: 32px;
  height: 32px;
  border: 2px solid #1da766;
  background: white;
  color: #1da766;
  border-radius: 50%;
  cursor: pointer;
  font-weight: bold;
  font-size: 14px;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-select:hover:not(:disabled) {
  background: #1da766;
  color: white;
  transform: scale(1.1);
}

.btn-select:disabled {
  border-color: #6c757d;
  color: #6c757d;
  cursor: not-allowed;
  transform: none;
}

.skills-drop-zone {
  min-height: 60vh;
  padding: 16px;
  background: white;
  border: 2px dashed #dee2e6;
  margin: 16px;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.skills-drop-zone.drag-over {
  border-color: #1da766;
  background: rgba(29, 167, 102, 0.05);
}

.drop-zone-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 200px;
  color: #6c757d;
  text-align: center;
}

.placeholder-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.7;
}

.placeholder-text {
  font-size: 0.9rem;
  line-height: 1.4;
  max-width: 200px;
}

.selected-skills-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.selected-skill-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f8f9fa;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.selected-skill-item:hover {
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.selected-skill-info {
  flex: 1;
}

.selected-skill-name {
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.selected-skill-costs {
  display: flex;
  gap: 12px;
  font-size: 0.85rem;
}

.btn-remove {
  width: 24px;
  height: 24px;
  border: none;
  background: #dc3545;
  color: white;
  border-radius: 50%;
  cursor: pointer;
  font-weight: bold;
  font-size: 16px;
  line-height: 1;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-remove:hover {
  background: #c82333;
  transform: scale(1.1);
}

.total-costs {
  margin-top: 16px;
  padding: 12px;
  background: #fff3cd;
  border: 1px solid #ffeaa7;
  border-radius: 6px;
}

.total-costs-header {
  font-weight: 600;
  color: #856404;
  margin-bottom: 8px;
}

.total-costs-breakdown {
  display: flex;
  gap: 16px;
  margin-bottom: 8px;
}

.total-ep,
.total-gold {
  font-weight: 600;
  font-size: 0.9rem;
}

.total-ep {
  color: #1da766;
}

.total-gold {
  color: #ffc107;
}

.affordability-check {
  font-size: 0.85rem;
  font-weight: 600;
}

.text-success {
  color: #28a745;
}

.text-danger {
  color: #dc3545;
}

.loading-skills {
  padding: 20px;
  text-align: center;
  color: #6c757d;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.simple-input-section {
  margin-bottom: 20px;
  padding: 16px;
  background: #fff3cd;
  border: 1px solid #ffeaa7;
  border-radius: 8px;
}

.simple-input-section .form-group {
  margin-bottom: 0;
}

/* Responsive Design */
@media (max-width: 768px) {
  .modal-content {
    width: 95%;
    margin: 10px;
    max-width: none;
  }
  
  .skills-selection-container {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .resources-display {
    flex-direction: column;
    gap: 10px;
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
  
  .category-filter-btn {
    width: 100%;
    text-align: center;
  }
}
</style>
