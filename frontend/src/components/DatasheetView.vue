<template>
    <div class="fullwidth-container datasheet-container" v-if="character">
      <!-- Character Overview -->
      <div class="character-overview">
        <div class="character-image">
          <img :src="imageSrc" alt="Character Image"/>
          <ImageUploadCropper 
            :characterId="character.id" 
            @image-updated="handleImageUpdate"
          />
        </div>
        <div class="character-stats">
          <div class="stat" v-for="(stat, index) in characterStats" :key="index">
            <span>{{ $t(stat.label) }}</span>
            <strong 
              v-if="editingIndex !== index"
              @dblclick="startEdit(index, stat.path)"
              class="editable-value"
            >
              {{ getStat(stat.path) }}
            </strong>
            <input 
              v-else
              v-model="editValue"
              @blur="saveEdit(stat.path)"
              @keyup.enter="saveEdit(stat.path)"
              @keyup.esc="cancelEdit"
              ref="editInput"
              type="number"
              class="edit-input"
            />
          </div>
        </div>
      </div>
      
      <!-- Character Information -->
      <div class="character-info">
        <div class="info-section">
          <p>
            <span 
              v-if="editingProp !== 'typ'" 
              @dblclick="startEditProp('typ', character.typ)"
              class="editable-prop"
            >{{ character.typ || 'x' }}</span>
            <input v-else v-model="editPropValue" @blur="saveProp('typ')" @keyup.enter="saveProp('typ')" @keyup.esc="cancelEditProp" ref="propInput" class="prop-input" />
            (
            <span 
              v-if="editingProp !== 'gender'" 
              @dblclick="startEditProp('gender', character.gender)"
              class="editable-prop"
            >{{ character.gender || 'x' }}</span>
            <input v-else v-model="editPropValue" @blur="saveProp('gender')" @keyup.enter="saveProp('gender')" @keyup.esc="cancelEditProp" ref="propInput" class="prop-input" />
            ),
            Grad:  
            <span 
              v-if="editingProp !== 'grad'" 
              @dblclick="startEditProp('grad', character.grad, 'number')"
              class="editable-prop"
            >{{ character.grad || 'x' }}</span>
            <input v-else v-model="editPropValue" type="number" @blur="saveProp('grad')" @keyup.enter="saveProp('grad')" @keyup.esc="cancelEditProp" ref="propInput" class="prop-input" />,
            Rasse:  
            <span 
              v-if="editingProp !== 'rasse'" 
              @dblclick="startEditProp('rasse', character.rasse)"
              class="editable-prop"
            >{{ character.rasse || 'x' }}</span>
            <input v-else v-model="editPropValue" @blur="saveProp('rasse')" @keyup.enter="saveProp('rasse')" @keyup.esc="cancelEditProp" ref="propInput" class="prop-input" />,
            Heimat: 
            <span 
              v-if="editingProp !== 'origin'" 
              @dblclick="startEditProp('origin', character.origin)"
              class="editable-prop"
            >{{ character.origin || '-' }}</span>
            <input v-else v-model="editPropValue" @blur="saveProp('origin')" @keyup.enter="saveProp('origin')" @keyup.esc="cancelEditProp" ref="propInput" class="prop-input" />,
            Stand:  
            <span 
              v-if="editingProp !== 'social_class'" 
              @dblclick="startEditProp('social_class', character.social_class)"
              class="editable-prop"
            >{{ character.social_class || '-' }}</span>
            <input v-else v-model="editPropValue" @blur="saveProp('social_class')" @keyup.enter="saveProp('social_class')" @keyup.esc="cancelEditProp" ref="propInput" class="prop-input" />.
          </p>
          <p v-if="character.rasse==='Zwerg'">
            Hort für Grad {{ character.grad || 'x' }}: 125 GS, für nächsten Grad: 250 GS.
          </p>
          <p>
            <strong>Spezialisierung:</strong> 
            <span 
              v-if="editingProp !== 'spezialisierung'" 
              @dblclick="startEditProp('spezialisierung', character.spezialisierung)"
              class="editable-prop"
            >{{ character.spezialisierung || '-' }}</span>
            <input v-else v-model="editPropValue" @blur="saveProp('spezialisierung')" @keyup.enter="saveProp('spezialisierung')" @keyup.esc="cancelEditProp" ref="propInput" class="prop-input" />.
          </p>
          <p>
            Alter: 
            <span 
              v-if="editingProp !== 'alter'" 
              @dblclick="startEditProp('alter', character.alter, 'number')"
              class="editable-prop"
            >{{ character.alter || 'xx' }}</span>
            <input v-else v-model="editPropValue" type="number" @blur="saveProp('alter')" @keyup.enter="saveProp('alter')" @keyup.esc="cancelEditProp" ref="propInput" class="prop-input" />,
            <strong v-if="editingProp !== 'hand'" @dblclick="startEditProp('hand', character.hand)" class="editable-prop">
              <span v-if="character.hand=='rechts'">Rechtshänder</span>
              <span v-else-if="character.hand=='links'">Linkshänder</span>
              <span v-else>Beidhändig</span>
            </strong>
            <input v-else v-model="editPropValue" @blur="saveProp('hand')" @keyup.enter="saveProp('hand')" @keyup.esc="cancelEditProp" ref="propInput" class="prop-input" />,
            Größe: 
            <span 
              v-if="editingProp !== 'groesse'" 
              @dblclick="startEditProp('groesse', character.groesse, 'number')"
              class="editable-prop"
            >{{ character.groesse }}</span>
            <input v-else v-model="editPropValue" type="number" @blur="saveProp('groesse')" @keyup.enter="saveProp('groesse')" @keyup.esc="cancelEditProp" ref="propInput" class="prop-input" />cm,
            Gewicht: 
            <span 
              v-if="editingProp !== 'gewicht'" 
              @dblclick="startEditProp('gewicht', character.gewicht, 'number')"
              class="editable-prop"
            >{{ character.gewicht }}</span>
            <input v-else v-model="editPropValue" type="number" @blur="saveProp('gewicht')" @keyup.enter="saveProp('gewicht')" @keyup.esc="cancelEditProp" ref="propInput" class="prop-input" />kg,
            Gestalt: 
            <span 
              v-if="editingProp !== 'merkmale.groesse'" 
              @dblclick="startEditProp('merkmale.groesse', character.merkmale?.groesse)"
              class="editable-prop"
            >{{ character.merkmale?.groesse || '-' }}</span>
            <input v-else v-model="editPropValue" @blur="saveProp('merkmale.groesse')" @keyup.enter="saveProp('merkmale.groesse')" @keyup.esc="cancelEditProp" ref="propInput" class="prop-input" />
            und 
            <span 
              v-if="editingProp !== 'merkmale.breite'" 
              @dblclick="startEditProp('merkmale.breite', character.merkmale?.breite)"
              class="editable-prop"
            >{{ character.merkmale?.breite || '-' }}</span>
            <input v-else v-model="editPropValue" @blur="saveProp('merkmale.breite')" @keyup.enter="saveProp('merkmale.breite')" @keyup.esc="cancelEditProp" ref="propInput" class="prop-input" />,
            Augen: 
            <span 
              v-if="editingProp !== 'merkmale.augenfarbe'" 
              @dblclick="startEditProp('merkmale.augenfarbe', character.merkmale?.augenfarbe)"
              class="editable-prop"
            >{{ character.merkmale?.augenfarbe || '-' }}</span>
            <input v-else v-model="editPropValue" @blur="saveProp('merkmale.augenfarbe')" @keyup.enter="saveProp('merkmale.augenfarbe')" @keyup.esc="cancelEditProp" ref="propInput" class="prop-input" />,
            Haare: 
            <span 
              v-if="editingProp !== 'merkmale.haarfarbe'" 
              @dblclick="startEditProp('merkmale.haarfarbe', character.merkmale?.haarfarbe)"
              class="editable-prop"
            >{{ character.merkmale?.haarfarbe || '-' }}</span>
            <input v-else v-model="editPropValue" @blur="saveProp('merkmale.haarfarbe')" @keyup.enter="saveProp('merkmale.haarfarbe')" @keyup.esc="cancelEditProp" ref="propInput" class="prop-input" />,
            Glaube: 
            <span 
              v-if="editingProp !== 'glaube'" 
              @dblclick="startEditProp('glaube', character.glaube)"
              class="editable-prop"
            >{{ character.glaube || '-' }}</span>
            <input v-else v-model="editPropValue" @blur="saveProp('glaube')" @keyup.enter="saveProp('glaube')" @keyup.esc="cancelEditProp" ref="propInput" class="prop-input" />
          </p>
          <p>
            <strong>Merkmale:</strong> 
            <span 
              v-if="editingProp !== 'merkmale.sonstige'" 
              @dblclick="startEditProp('merkmale.sonstige', character.merkmale?.sonstige)"
              class="editable-prop"
            >{{ character.merkmale?.sonstige || '-' }}</span>
            <input v-else v-model="editPropValue" @blur="saveProp('merkmale.sonstige')" @keyup.enter="saveProp('merkmale.sonstige')" @keyup.esc="cancelEditProp" ref="propInput" class="prop-input" style="width: 300px;" />
          </p>
          <p>
            <em>Persönlicher Bonus für</em> Ausdauer 12, Schaden 5, Angriff 2,
            Abwehr 0, Zauber 0, Resistenz 3 / 4.
          </p>
        </div>
      </div>
    </div>
    <div v-else>Loading character data...</div>
</template>

<style>
/* All common styles moved to main.css */

/* DatasheetView specific styles */
.datasheet-container {
  padding-top: 10px;
}

.info-section {
  max-width: none;
  white-space: normal;
  line-height: 1.6;
}

.info-section p {
  margin: 15px 0;
  padding: 0;
}

.character-overview {
  margin-bottom: 30px;
  margin-top: 0;
}

.character-image {
  position: relative;
}

.character-image .image-upload-container {
  position: absolute;
  bottom: 10px;
  right: 10px;
}

.character-info {
  margin-top: 20px;
}

.editable-value {
  cursor: pointer;
  padding: 2px 4px;
  border-radius: 3px;
  transition: background-color 0.2s;
}

.editable-value:hover {
  background-color: rgba(0, 123, 255, 0.1);
}

.edit-input {
  width: 60px;
  padding: 2px 4px;
  font-size: inherit;
  font-weight: bold;
  border: 2px solid var(--primary-color);
  border-radius: 3px;
  text-align: center;
}

.edit-input:focus {
  outline: none;
  border-color: #0056b3;
}

.editable-prop {
  cursor: pointer;
  padding: 1px 3px;
  border-radius: 2px;
  transition: background-color 0.2s;
  display: inline-block;
  min-width: 20px;
}

.editable-prop:hover {
  background-color: rgba(0, 123, 255, 0.1);
}

.prop-input {
  padding: 1px 4px;
  font-size: inherit;
  border: 1px solid var(--primary-color);
  border-radius: 3px;
  min-width: 60px;
}

.prop-input:focus {
  outline: none;
  border-color: #0056b3;
}
</style>

<script>
import ImageUploadCropper from './ImageUploadCropper.vue'
import API from '../utils/api'

export default {
  name: "DatasheetView",
  components: {
    ImageUploadCropper
  },
  props: {
    character: {
      type: Object,
      required: true
    }
  },
  computed: {
    imageSrc() {
      return this.character.image
        ? `${this.character.image}`
        : "/token_default.png";
    },
  },
  data() {
    return {
      editingIndex: null,
      editValue: '',
      editingProp: null,
      editPropValue: '',
      editPropType: 'text',
      characterStats: [
        { label: 'stats.strength', path: 'eigenschaften.6.value' },
        { label: 'stats.dexterity', path: 'eigenschaften.1.value' },
        { label: 'stats.agility', path: 'eigenschaften.2.value' },
        { label: 'stats.constitution', path: 'eigenschaften.4.value' },
        { label: 'stats.intelligence', path: 'eigenschaften.3.value' },
        { label: 'stats.spelltalent', path: 'eigenschaften.8.value' },
        { label: 'stats.beauty', path: 'eigenschaften.0.value' },
        { label: 'stats.charisma', path: 'eigenschaften.5.value' },
        { label: 'stats.willpower', path: 'eigenschaften.7.value' },
        { label: 'stats.poisontolerance', path: 'git' },
        { label: 'stats.movement', path: 'b.max' },
        { label: 'stats.lifepoints', path: 'lp.max'},
        { label: 'stats.staminapoints', path: 'ap.max'},
        { label: 'stats.divinegrace', path: 'bennies.gg'},
        { label: 'stats.fatesfavor', path: 'bennies.sg' }
      ]
    }
  },
  methods: {
    handleImageUpdate(newImage) {
      this.$emit('character-updated')
    },
    getStat(path) {
      if (path === 'git') {
        return '64!'
      }
      return path.split('.').reduce((obj, key) => obj?.[key], this.character) ?? '-'
    },
    startEdit(index, path) {
      if (path === 'git') return
      
      this.editingIndex = index
      this.editValue = this.getStat(path)
      this.$nextTick(() => {
        if (this.$refs.editInput && this.$refs.editInput[0]) {
          this.$refs.editInput[0].focus()
          this.$refs.editInput[0].select()
        }
      })
    },
    async saveEdit(path) {
      if (this.editingIndex === null) return
      
      const newValue = parseInt(this.editValue)
      if (isNaN(newValue)) {
        this.cancelEdit()
        return
      }
      
      try {
        // Update the character object directly
        const pathParts = path.split('.')
        let obj = this.character
        for (let i = 0; i < pathParts.length - 1; i++) {
          obj = obj[pathParts[i]]
        }
        obj[pathParts[pathParts.length - 1]] = newValue
        
        // Save to backend
        await API.put(`/api/characters/${this.character.id}`, this.character)
        
        this.$emit('character-updated')
        this.cancelEdit()
      } catch (error) {
        console.error('Failed to update stat:', error)
        alert('Fehler beim Speichern: ' + (error.response?.data?.error || error.message))
        this.cancelEdit()
      }
    },
    cancelEdit() {
      this.editingIndex = null
      this.editValue = ''
    },
    startEditProp(prop, value, type = 'text') {
      this.editingProp = prop
      this.editPropValue = value || ''
      this.editPropType = type
      this.$nextTick(() => {
        if (this.$refs.propInput) {
          const input = Array.isArray(this.$refs.propInput) ? this.$refs.propInput[0] : this.$refs.propInput
          if (input) {
            input.focus()
            input.select()
          }
        }
      })
    },
    async saveProp(prop) {
      if (this.editingProp === null) return
      // Update the character object directly
      const pathParts = prop.split('.')
      let obj = this.character
      for (let i = 0; i < pathParts.length - 1; i++) {
        if (!obj[pathParts[i]]) {
          obj[pathParts[i]] = {}
        }
        obj = obj[pathParts[i]]
      }
      obj[pathParts[pathParts.length - 1]] = newValue
      
      // Save to backend
      await API.put(`/api/characters/${this.character.id}`, this.character)
      
      this.$emit('character-updated',this.character)
      let newValue = this.editPropValue
      if (this.editPropType === 'number') {
        newValue = parseInt(this.editPropValue)
        if (isNaN(newValue)) {
          this.cancelEditProp()
          return
        }
      }
      
      try {
        this.$emit('update-property', prop, newValue)
        this.cancelEditProp()
      } catch (error) {
        console.error('Failed to update property:', error)
        alert('Fehler beim Speichern: ' + (error.response?.data?.error || error.message))
        this.cancelEditProp()
      }
    },
    cancelEditProp() {
      this.editingProp = null
      this.editPropValue = ''
      this.editPropType = 'text'
    }
  }
};
</script>
