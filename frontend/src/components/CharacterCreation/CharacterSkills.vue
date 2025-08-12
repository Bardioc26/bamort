<template>
  <div class="skills-form character-creation-container">
    <h2>Skills & Spells</h2>
    <p class="instruction">Select skills and spells for your character. Each category has a limited number of learning points.</p>
    
    <!-- Debug Info -->
    <div v-if="true" style="background: #f0f0f0; padding: 10px; margin: 10px 0; border-radius: 4px; font-size: 12px;">
      <strong>Debug Info:</strong><br>
      skillCategories.length: {{ skillCategories.length }}<br>
      selectedCategory: {{ selectedCategory }}<br>
      availableSkills.length: {{ availableSkills.length }}<br>
      selectedSkills.length: {{ selectedSkills.length }}
    </div>
    
    <div class="skills-content">
      <!-- Left Column: Categories and Skills -->
      <div class="left-column">
        <!-- Skill Categories -->
        <div class="categories-section">
          <h3>Skill Categories</h3>
          
          <div v-if="skillCategories.length === 0" class="no-categories">
            <p>No skill categories available. Loading...</p>
          </div>
          
          <div v-else class="categories-grid">
            <div 
              v-for="category in skillCategories" 
              :key="category.name"
              :class="['category-card', { active: selectedCategory === category.name }]"
              @click="selectCategory(category.name)"
              v-if="category.name !== 'zauber'"
            >
              <div class="category-header">
                <h4>{{ category.display_name }}</h4>
                <div class="points-info">
                  <span class="remaining">{{ category.points }}</span> / 
                  <span class="total">{{ category.max_points }}</span>
                </div>
              </div>
              <div class="progress-bar">
                <div 
                  class="progress-fill" 
                  :style="{ width: ((category.max_points - category.points) / category.max_points * 100) + '%' }"
                ></div>
              </div>
            </div>
          </div>
        </div>

        <!-- Skills List for Selected Category -->
        <div v-if="selectedCategory" class="skills-section">
          <h3>{{ getCategoryDisplayName(selectedCategory) }} Skills</h3>
          <div class="skills-list">
            <div 
              v-for="skill in availableSkills" 
              :key="skill.name"
              class="skill-item"
            >
              <div class="skill-info">
                <span class="skill-name">{{ skill.name }}</span>
                <span class="skill-cost">Cost: {{ skill.cost }} EP</span>
              </div>
              <button 
                @click="addSkill(skill)"
                :disabled="!canAddSkill(skill)"
                class="add-btn"
              >
                Add
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column: Selected Skills -->
      <div class="right-column">
        <div class="selected-section">
          <h3>Selected Skills</h3>
          
          <div v-if="selectedSkills.length > 0" class="selected-skills">
            <div class="selected-list">
              <div 
                v-for="skill in selectedSkills" 
                :key="skill.name"
                class="selected-item"
              >
                <span class="item-name">{{ skill.name }}</span>
                <span class="item-category">{{ skill.category }}</span>
                <span class="item-cost">{{ skill.cost }} EP</span>
                <button @click="removeSkill(skill)" class="remove-btn">×</button>
              </div>
            </div>
          </div>

          <div v-if="selectedSkills.length === 0" class="no-selection">
            No skills selected yet. Click on a category above to start selecting skills.
          </div>
        </div>
      </div>
    </div>

    <div class="form-actions">
      <button type="button" @click="handlePrevious" class="prev-btn">
        ← Previous: Derived Values
      </button>
      <button type="button" @click="handleNext" class="next-btn">
        Next: Spells →
      </button>
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
    },
    skillCategories: {
      type: Array,
      required: true,
    }
  },
  emits: ['previous', 'next', 'save'],
  data() {
    return {
      selectedCategory: null,
      availableSkills: [],
      selectedSkills: [],
    }
  },
  async created() {
    console.log('CharacterSkills created, sessionData:', this.sessionData)
    console.log('CharacterSkills created, skillCategories:', this.skillCategories)
    
    // Initialize with session data
    if (this.sessionData.skills) {
      this.selectedSkills = [...this.sessionData.skills]
      console.log('Initialized selectedSkills:', this.selectedSkills)
    }
    
    this.updateCategoryPoints()
    console.log('Updated skillCategories after points update:', this.skillCategories)
  },
  methods: {
    async selectCategory(categoryName) {
      this.selectedCategory = categoryName
      await this.loadSkills(categoryName)
    },
    
    async loadSkills(category) {
      try {
        const token = localStorage.getItem('token')
        
        // Create a dummy request for skills in this category
        const request = {
          characterClass: this.sessionData.typ || 'Abenteurer',
          characterId: '0', // Dummy for new character
          category: category,
        }
        
        const response = await API.post('/api/characters/available-skills-new', request, {
          headers: { Authorization: `Bearer ${token}` },
        })
        
        this.availableSkills = response.data.skills || []
      } catch (error) {
        console.error('Error loading skills:', error)
        // Fallback dummy data
        this.availableSkills = [
          { name: 'Sample Skill 1', cost: 30, category: category },
          { name: 'Sample Skill 2', cost: 40, category: category },
          { name: 'Sample Skill 3', cost: 50, category: category },
        ]
      }
    },
    
    getCategoryDisplayName(categoryName) {
      const category = this.skillCategories.find(c => c.name === categoryName)
      return category ? category.display_name : categoryName
    },
    
    canAddSkill(skill) {
      const category = this.skillCategories.find(c => c.name === skill.category)
      const alreadySelected = this.selectedSkills.some(s => s.name === skill.name)
      
      return category && category.points >= skill.cost && !alreadySelected
    },
    
    addSkill(skill) {
      if (this.canAddSkill(skill)) {
        this.selectedSkills.push({ ...skill })
        this.updateCategoryPoints()
      }
    },
    
    removeSkill(skill) {
      const index = this.selectedSkills.findIndex(s => s.name === skill.name)
      if (index >= 0) {
        this.selectedSkills.splice(index, 1)
        this.updateCategoryPoints()
      }
    },
    
    updateCategoryPoints() {
      // Reset all categories to max points
      this.skillCategories.forEach(category => {
        category.points = category.max_points
      })
      
      // Deduct points for selected skills
      this.selectedSkills.forEach(skill => {
        const category = this.skillCategories.find(c => c.name === skill.category)
        if (category) {
          category.points -= skill.cost
        }
      })
    },
    
    handlePrevious() {
      this.$emit('previous')
    },
    
    handleNext() {
      const data = {
        skills: this.selectedSkills,
        skill_points: this.skillCategories.reduce((acc, cat) => {
          acc[cat.name] = cat.points
          return acc
        }, {})
      }
      
      this.$emit('next', data)
    },
  }
}
</script>

<style scoped>
.skills-form {
  width: 100% !important;
  max-width: none !important;
  margin: 0;
}

.skills-form h2 {
  text-align: center;
  margin-bottom: 10px;
  color: #333;
}

.instruction {
  text-align: center;
  margin-bottom: 30px;
  color: #666;
  font-style: italic;
}

.skills-content {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 30px;
  margin-bottom: 30px;
}

.left-column {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.right-column {
  position: sticky;
  top: 20px;
  height: fit-content;
}

.categories-section h3, .skills-section h3, .spells-section h3, .selected-section h3 {
  margin-bottom: 15px;
  color: #333;
}

.categories-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 15px;
}

.no-categories {
  text-align: center;
  padding: 20px;
  color: #666;
  font-style: italic;
}

.category-card {
  padding: 15px;
  border: 2px solid #ddd;
  border-radius: 8px;
  background-color: #fafafa;
  cursor: pointer;
  transition: all 0.3s ease;
}

.category-card:hover {
  border-color: #2196f3;
  background-color: #f0f8ff;
}

.category-card.active {
  border-color: #2196f3;
  background-color: #e3f2fd;
}

.category-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.category-header h4 {
  margin: 0;
  font-size: 14px;
  color: #333;
}

.points-info {
  font-size: 12px;
  color: #666;
}

.remaining {
  font-weight: bold;
  color: #2196f3;
}

.progress-bar {
  height: 4px;
  background-color: #e0e0e0;
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background-color: #4caf50;
  transition: width 0.3s ease;
}

.skills-list, .spells-list, .selected-list {
  max-height: 400px;
  overflow-y: auto;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.skill-item, .spell-item, .selected-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #eee;
}

.skill-item:last-child, .spell-item:last-child, .selected-item:last-child {
  border-bottom: none;
}

.skill-info, .spell-info {
  flex: 1;
}

.skill-name, .spell-name, .item-name {
  display: block;
  font-weight: bold;
  color: #333;
}

.skill-cost, .spell-cost, .item-cost {
  font-size: 12px;
  color: #666;
}

.item-category {
  font-size: 12px;
  color: #888;
  margin-right: 10px;
}

.add-btn, .remove-btn {
  padding: 5px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.add-btn {
  background-color: #4caf50;
  color: white;
}

.add-btn:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.add-btn:hover:not(:disabled) {
  background-color: #45a049;
}

.remove-btn {
  background-color: #f44336;
  color: white;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.remove-btn:hover {
  background-color: #d32f2f;
}

.selected-section {
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 8px;
  height: fit-content;
}

.selected-skills, .selected-spells {
  margin-bottom: 20px;
}

.selected-skills h4, .selected-spells h4 {
  margin-bottom: 10px;
  color: #555;
}

.no-selection {
  text-align: center;
  color: #999;
  font-style: italic;
  padding: 20px;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.prev-btn, .finalize-btn {
  padding: 12px 30px;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.prev-btn {
  background-color: #6c757d;
  color: white;
}

.prev-btn:hover {
  background-color: #5a6268;
}

.finalize-btn {
  background-color: #4caf50;
  color: white;
}

.finalize-btn:hover {
  background-color: #45a049;
}
</style>
