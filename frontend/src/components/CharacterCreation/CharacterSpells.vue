<template>
  <div class="character-spells">
    <h2>{{ $t('characters.create.spells.title') }}</h2>
    <p class="subtitle">{{ $t('characters.create.spells.description') }}</p>
    
    <!-- Spell Points Display -->
    <div class="spell-points-display">
      <div class="category-points">
        <span class="category-name">{{ $t('characters.spells.zauber') }}</span>
        <span class="points">{{ getZauberCategory()?.points || 0 }} / {{ getZauberCategory()?.max_points || 0 }}</span>
      </div>
    </div>
    
    <!-- Available Spells -->
    <div class="available-spells">
      <h3>{{ $t('characters.create.spells.available') }}</h3>
      <div class="spells-grid">
        <div
          v-for="spell in availableSpells"
          :key="spell.name"
          class="spell-card"
          :class="{ 'can-add': canAddSpell(spell), 'cannot-add': !canAddSpell(spell) }"
          @click="addSpell(spell)"
        >
          <div class="spell-name">{{ spell.name }}</div>
          <div class="spell-cost">{{ $t('characters.skills.cost') }}: {{ spell.cost }}</div>
          <div v-if="canAddSpell(spell)" class="add-button">
            <i class="fas fa-plus"></i>
          </div>
          <div v-else class="disabled-reason">
            <span v-if="isSpellSelected(spell)">{{ $t('characters.skills.alreadySelected') }}</span>
            <span v-else>{{ $t('characters.skills.notEnoughPoints') }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Selected Spells -->
    <div class="selected-spells" v-if="selectedSpells.length > 0">
      <h3>{{ $t('characters.create.spells.selected') }}</h3>
      <div class="selected-list">
        <div
          v-for="spell in selectedSpells"
          :key="spell.name"
          class="selected-spell"
        >
          <span class="spell-name">{{ spell.name }}</span>
          <span class="spell-cost">({{ spell.cost }})</span>
          <button
            @click="removeSpell(spell)"
            class="remove-button"
            :title="$t('characters.skills.remove')"
          >
            <i class="fas fa-times"></i>
          </button>
        </div>
      </div>
    </div>
    
    <!-- Navigation Buttons -->
    <div class="navigation-buttons">
      <button @click="handlePrevious" class="btn-secondary">
        <i class="fas fa-arrow-left"></i>
        {{ $t('common.previous') }}
      </button>
      
      <button @click="handleFinalize" class="btn-primary">
        {{ $t('characters.create.finalize') }}
        <i class="fas fa-check"></i>
      </button>
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
    },
    skillCategories: {
      type: Array,
      required: true
    }
  },
  
  emits: ['previous', 'finalize', 'save'],
  
  data() {
    return {
      availableSpells: [],
      selectedSpells: [],
    }
  },
  
  async mounted() {
    await this.loadSpells()
    this.loadSelectedSpells()
  },
  
  methods: {
    async loadSpells() {
      try {
        const token = localStorage.getItem('token')
        
        const request = {
          characterClass: this.sessionData.typ || 'Abenteurer',
          characterId: '0', // Dummy for new character
        }
        
        const response = await API.post('/api/characters/available-spells-new', request, {
          headers: { Authorization: `Bearer ${token}` },
        })
        
        this.availableSpells = response.data.spells || []
      } catch (error) {
        console.error('Error loading spells:', error)
        // Fallback dummy data
        this.availableSpells = [
          { name: 'Licht', cost: 50, description: 'Erzeugt ein helles Licht' },
          { name: 'Magisches Geschoss', cost: 60, description: 'Feuert ein magisches Projektil ab' },
          { name: 'Schildzauber', cost: 70, description: 'Erzeugt einen magischen Schutzschild' },
          { name: 'Heilung', cost: 80, description: 'Heilt leichte Wunden' },
          { name: 'Unsichtbarkeit', cost: 120, description: 'Macht den Zauberer unsichtbar' },
        ]
      }
    },
    
    loadSelectedSpells() {
      // Load from session data if available
      if (this.sessionData.spells) {
        this.selectedSpells = [...this.sessionData.spells]
        this.updateSpellPoints()
      }
    },
    
    getZauberCategory() {
      return this.skillCategories.find(c => c.name === 'zauber')
    },
    
    canAddSpell(spell) {
      const category = this.getZauberCategory()
      const alreadySelected = this.isSpellSelected(spell)
      
      return category && category.points >= spell.cost && !alreadySelected
    },
    
    isSpellSelected(spell) {
      return this.selectedSpells.some(s => s.name === spell.name)
    },
    
    addSpell(spell) {
      if (this.canAddSpell(spell)) {
        this.selectedSpells.push({ ...spell })
        this.updateSpellPoints()
        this.saveData()
      }
    },
    
    removeSpell(spell) {
      const index = this.selectedSpells.findIndex(s => s.name === spell.name)
      if (index >= 0) {
        this.selectedSpells.splice(index, 1)
        this.updateSpellPoints()
        this.saveData()
      }
    },
    
    updateSpellPoints() {
      const zauberCategory = this.getZauberCategory()
      if (zauberCategory) {
        // Reset to max points
        zauberCategory.points = zauberCategory.max_points
        
        // Deduct points for selected spells
        this.selectedSpells.forEach(spell => {
          zauberCategory.points -= spell.cost
        })
      }
    },
    
    saveData() {
      const data = {
        spells: this.selectedSpells,
        spell_points: {
          zauber: this.getZauberCategory()?.points || 0
        }
      }
      
      this.$emit('save', data)
    },
    
    handlePrevious() {
      this.saveData()
      this.$emit('previous')
    },
    
    handleFinalize() {
      const data = {
        spells: this.selectedSpells,
        spell_points: {
          zauber: this.getZauberCategory()?.points || 0
        }
      }
      
      this.$emit('finalize', data)
    },
  },
}
</script>

<style scoped>
.character-spells {
  padding: 1rem;
  max-width: 1200px;
  margin: 0 auto;
}

.subtitle {
  color: #666;
  margin-bottom: 1.5rem;
  font-style: italic;
}

.spell-points-display {
  margin-bottom: 2rem;
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.category-points {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  font-size: 1.1rem;
}

.category-name {
  color: #495057;
}

.points {
  color: #007bff;
}

.available-spells {
  margin-bottom: 2rem;
}

.spells-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1rem;
  margin-top: 1rem;
}

.spell-card {
  border: 2px solid #e9ecef;
  border-radius: 8px;
  padding: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
  background: white;
  position: relative;
}

.spell-card.can-add {
  border-color: #28a745;
  background: #f8fff9;
}

.spell-card.can-add:hover {
  border-color: #20c997;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}

.spell-card.cannot-add {
  border-color: #dc3545;
  background: #fff5f5;
  cursor: not-allowed;
  opacity: 0.7;
}

.spell-name {
  font-weight: 600;
  color: #495057;
  margin-bottom: 0.5rem;
}

.spell-cost {
  color: #6c757d;
  font-size: 0.9rem;
  margin-bottom: 0.5rem;
}

.add-button {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  color: #28a745;
  font-size: 1.2rem;
}

.disabled-reason {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  color: #dc3545;
  font-size: 0.8rem;
  text-align: right;
}

.selected-spells {
  margin-bottom: 2rem;
}

.selected-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 1rem;
}

.selected-spell {
  background: #e3f2fd;
  border: 1px solid #bbdefb;
  border-radius: 20px;
  padding: 0.5rem 1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.selected-spell .spell-name {
  font-weight: 500;
  margin-bottom: 0;
}

.selected-spell .spell-cost {
  color: #1976d2;
  font-size: 0.85rem;
}

.remove-button {
  background: none;
  border: none;
  color: #f44336;
  cursor: pointer;
  padding: 0;
  margin-left: 0.5rem;
  font-size: 0.9rem;
}

.remove-button:hover {
  color: #d32f2f;
}

.navigation-buttons {
  display: flex;
  justify-content: space-between;
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid #e9ecef;
}

.btn-primary, .btn-secondary {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: background-color 0.3s ease;
}

.btn-primary {
  background: #007bff;
  color: white;
}

.btn-primary:hover {
  background: #0056b3;
}

.btn-secondary {
  background: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background: #545b62;
}

/* Responsive Design */
@media (max-width: 768px) {
  .spells-grid {
    grid-template-columns: 1fr;
  }
  
  .navigation-buttons {
    flex-direction: column;
    gap: 1rem;
  }
  
  .btn-primary, .btn-secondary {
    justify-content: center;
  }
}
</style>
