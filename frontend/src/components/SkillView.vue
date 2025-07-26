<template>
  <div class="cd-view">
    <div class="cd-list">
      <div class="tables-container">
        <div class="table-wrapper-left">
          <div class="header-section">
            <h2 style="line-height: 1.5; margin-top: 5px;">Fertigkeiten</h2>
            
            <!-- Lernmodus Toggle Button -->
            <div class="learning-mode-controls">
              <!-- Ressourcen-Anzeige (nur sichtbar wenn Lernmodus aktiv) -->
              <div v-if="learningMode" class="resources-display">
                <div class="resource-item">
                  <span class="resource-icon">⚡</span>
                  <span class="resource-value">{{ character.erfahrungsschatz?.ep || 0 }} EP</span>
                </div>
                <div class="resource-item">
                  <span class="resource-icon">💰</span>
                  <span class="resource-value">{{ character.vermoegen?.goldstücke || 0 }} Gold</span>
                </div>
              </div>
              
              <button 
                @click="toggleLearningMode" 
                class="btn-learning-mode"
                :class="{ active: learningMode }"
                :title="learningMode ? 'Lernmodus beenden' : 'Lernmodus aktivieren'"
              >
                <span class="icon">🎓</span>
              </button>
              
              <!-- Lernmodus-Buttons (nur sichtbar wenn Lernmodus aktiv) -->
              <div v-if="learningMode" class="learning-actions">
                <button 
                  @click="showLearnNewDialog" 
                  class="btn-learn-new"
                  title="Neue Fertigkeit lernen"
                >
                  <span class="icon">📚</span>
                </button>
                <button 
                  @click="showAddDialog" 
                  class="btn-add"
                  title="Fertigkeit hinzufügen"
                >
                  <span class="icon">➕</span>
                </button>
              </div>
            </div>
          </div>
          <table class="cd-table">
      <thead>
        <tr>
          <th class="cd-table-header" width="60%">{{ $t('skill.name') }}</th>
          <th class="cd-table-header" width="35">{{ $t('skill.value') }}</th>
          <th class="cd-table-header" width="35">{{ $t('skill.bonus') }}</th>
          <th class="cd-table-header" width="35">{{ $t('skill.pp') }}</th>
          <th class="cd-table-header" width="30%">{{ $t('skill.note') }}</th>
          <th v-if="learningMode" class="cd-table-header" width="80">Aktionen</th>
        </tr>
      </thead>
      <tbody>
        <template v-for="skills,categorie in character.categorizedskills">
        <tr>
          <td :colspan="learningMode ? 6 : 5">{{ categorie || '-' }}</td>
        </tr>
        <template v-for="skill in skills">
          <tr>
            <td>{{ skill.name || '-' }}</td>
            <td>{{ skill.fertigkeitswert || '-' }}</td>
            <td>{{ skill.bonus || '0' }}</td>
            <td class="pp-cell">
              <div class="pp-container">
                <button 
                  @click="decreasePP(skill)" 
                  class="pp-btn pp-btn-minus"
                  :disabled="(skill.pp || 0) <= 0"
                  title="Praxispunkt entfernen"
                >
                  −
                </button>
                <span class="pp-value">{{ skill.pp || '0' }}</span>
                <button 
                  @click="increasePP(skill)" 
                  class="pp-btn pp-btn-plus"
                  title="Praxispunkt hinzufügen"
                >
                  +
                </button>
              </div>
            </td>
            <td>{{ skill.bemerkung || '-' }}</td>
            <td v-if="learningMode" class="action-cell">
              <button 
                @click="improveSkill(skill)" 
                class="btn-action btn-improve-small"
                title="Fertigkeit verbessern"
              >
                ⬆️
              </button>
            </td>
          </tr>
        </template>
      </template>
      <tr>
        <td class="cd-table-header" :colspan="learningMode ? 6 : 5">Waffenfertigkeiten</td>
      </tr>
      <template v-for="skill in character.waffenfertigkeiten">
        <tr>
          <td>{{ skill.name || '-' }}</td>
          <td>{{ skill.fertigkeitswert || '-' }}</td>
          <td>{{ skill.bonus || '0' }}</td>
          <td class="pp-cell">
            <div class="pp-container">
              <button 
                @click="decreaseWeaponPP(skill)" 
                class="pp-btn pp-btn-minus"
                :disabled="(skill.pp || 0) <= 0"
                title="Praxispunkt entfernen"
              >
                −
              </button>
              <span class="pp-value">{{ skill.pp || '0' }}</span>
              <button 
                @click="increaseWeaponPP(skill)" 
                class="pp-btn pp-btn-plus"
                title="Praxispunkt hinzufügen"
              >
                +
              </button>
            </div>
          </td>
          <td>{{ skill.bemerkung || '-' }}</td>
          <td v-if="learningMode" class="action-cell">
            <button 
              @click="improveWeaponSkill(skill)" 
              class="btn-action btn-improve-small"
              title="Waffenfertigkeit verbessern"
            >
              ⬆️
            </button>
          </td>
        </tr>
      </template>
      </tbody>
          </table>
        </div>
        <div class="table-wrapper-right">
          <h2 style="line-height: 1.5; margin-top: 5px;"> Angeborene Fertigkeiten</h2>
          <table class="cd-table">
      <thead>
        <tr>
          <th class="cd-table-header" width="80%">{{ $t('skill.name') }}</th>
          <th class="cd-table-header" width="35">{{ $t('skill.value') }}</th>
        </tr>
      </thead>
      <tbody>
      <template v-for="skill in character.InnateSkills">
        <tr>
          <td>{{ skill.name || '-' }}</td>
          <td>{{ skill.fertigkeitswert || '-' }}</td>
        </tr>
      </template>
      </tbody>
          </table>
        </div>
      </div>
    </div> <!--- end cd-list-->
    
    <!-- Dialog für neue Fertigkeit lernen -->
    <div v-if="showLearnDialog" class="modal-overlay" @click.self="closeDialogs">
      <div class="modal-content">
        <h3>Neue Fertigkeit lernen</h3>
        <div class="form-group">
          <label>Fertigkeitsname:</label>
          <input v-model="newSkillName" type="text" placeholder="Name der Fertigkeit" />
        </div>
        <div class="form-group">
          <label>Notizen (optional):</label>
          <textarea v-model="learnNotes" placeholder="Zusätzliche Notizen..."></textarea>
        </div>
        <div class="modal-actions">
          <button @click="learnNewSkill" class="btn-confirm">Lernen</button>
          <button @click="closeDialogs" class="btn-cancel">Abbrechen</button>
        </div>
      </div>
    </div>

    <!-- Dialog für Fertigkeit verbessern -->
    <div v-if="showImproveSelectionDialog" class="modal-overlay" @click.self="closeDialogs">
      <div class="modal-content">
        <h3>Fertigkeit verbessern</h3>
        <div class="form-group">
          <label>Fertigkeit auswählen:</label>
          <select v-model="selectedSkillToImprove">
            <option value="">-- Fertigkeit wählen --</option>
            <optgroup label="Fertigkeiten">
              <template v-for="(skills, category) in character.categorizedskills" :key="category">
                <option v-for="skill in skills" 
                        :key="skill.name" 
                        :value="skill">
                  {{ skill.name }} ({{ skill.fertigkeitswert }})
                </option>
              </template>
            </optgroup>
            <optgroup label="Waffenfertigkeiten">
              <option v-for="skill in character.waffenfertigkeiten" 
                      :key="skill.name" 
                      :value="skill">
                {{ skill.name }} ({{ skill.fertigkeitswert }})
              </option>
            </optgroup>
          </select>
        </div>
        <div class="form-group">
          <label>Praxispunkte verwenden:</label>
          <input v-model.number="usePP" type="number" min="0" placeholder="Anzahl PP" />
        </div>
        <div class="form-group">
          <label>Notizen (optional):</label>
          <textarea v-model="improveNotes" placeholder="Zusätzliche Notizen..."></textarea>
        </div>
        <div class="modal-actions">
          <button @click="improveSelectedSkill" class="btn-confirm" :disabled="!selectedSkillToImprove">Verbessern</button>
          <button @click="closeDialogs" class="btn-cancel">Abbrechen</button>
        </div>
      </div>
    </div>

    <!-- Dialog für Fertigkeit hinzufügen -->
    <div v-if="showAddDialog" class="modal-overlay" @click.self="closeDialogs">
      <div class="modal-content">
        <h3>Fertigkeit hinzufügen</h3>
        <div class="form-group">
          <label>Fertigkeitsname:</label>
          <input v-model="addSkillName" type="text" placeholder="Name der Fertigkeit" />
        </div>
        <div class="form-group">
          <label>Fertigkeitswert:</label>
          <input v-model.number="addSkillValue" type="number" min="0" max="20" placeholder="Wert (0-20)" />
        </div>
        <div class="form-group">
          <label>Notizen (optional):</label>
          <textarea v-model="addNotes" placeholder="Zusätzliche Notizen..."></textarea>
        </div>
        <div class="modal-actions">
          <button @click="addNewSkill" class="btn-confirm">Hinzufügen</button>
          <button @click="closeDialogs" class="btn-cancel">Abbrechen</button>
        </div>
      </div>
    </div>

    <!-- Neue Dialog-Komponente für detailliertes Fertigkeiten-Lernen -->
    <SkillImproveDialog 
      :character="character"
      :skill="selectedSkillToLearn"
      :isVisible="showDetailedLearnDialog"
      :learningType="selectedLearningType"
      @close="closeDialogs"
      @skill-updated="handleSkillUpdated"
    />
  </div> <!--- end character -datasheet-->

</template>

<style>
.tables-container {
  display: flex;
  gap: 1rem;
  width: 100%;
}

.table-wrapper-left {
  flex: 6;
  min-width: 0; /* Prevent table from overflowing */
}
.table-wrapper-right {
  flex: 4;
  min-width: 0; /* Prevent table from overflowing */
}

.cd-table {
  width: 100%;
}
.cd-table-header {
  background-color: #1da766;
  font-weight: bold;
}

/* Header mit Lernmodus-Kontrollen */
.header-section {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 15px;
}

.learning-mode-controls {
  display: flex;
  align-items: center;
  gap: 15px;
}

/* Ressourcen-Anzeige */
.resources-display {
  display: flex;
  gap: 15px;
  animation: slideIn 0.3s ease;
}

.resource-item {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 6px 12px;
  background: #f8f9fa;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  font-weight: bold;
  color: #495057;
}

.resource-icon {
  font-size: 16px;
}

.resource-value {
  font-size: 14px;
  white-space: nowrap;
}

/* Lernmodus Toggle Button */
.btn-learning-mode {
  padding: 8px 16px;
  border: 2px solid #1da766;
  background: white;
  color: #1da766;
  border-radius: 6px;
  cursor: pointer;
  font-weight: bold;
  display: flex;
  align-items: center;
  gap: 5px;
  transition: all 0.3s ease;
  position: relative;
}

.btn-learning-mode:hover {
  background: #1da766;
  color: white;
}

.btn-learning-mode.active {
  background: #1da766;
  color: white;
}

/* Lernmodus Action Buttons */
.learning-actions {
  display: flex;
  gap: 5px;
  animation: slideIn 0.3s ease;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.btn-learn-new,
.btn-improve,
.btn-add {
  width: 40px;
  height: 40px;
  border: 2px solid #007bff;
  background: white;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  transition: all 0.2s ease;
  position: relative;
}

.btn-learn-new:hover {
  background: #007bff;
  color: white;
}

.btn-improve {
  border-color: #28a745;
}

.btn-improve:hover {
  background: #28a745;
  color: white;
}

.btn-add {
  border-color: #17a2b8;
}

.btn-add:hover {
  background: #17a2b8;
  color: white;
}

/* Aktions-Buttons in der Tabelle */
.action-cell {
  text-align: center;
  padding: 4px;
}

.btn-action {
  padding: 4px 8px;
  border: 1px solid #28a745;
  background: white;
  color: #28a745;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.2s ease;
  position: relative;
}

.btn-action:hover {
  background: #28a745;
  color: white;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  padding: 24px;
  max-width: 500px;
  width: 90%;
  max-height: 80vh;
  overflow-y: auto;
  animation: modalSlideIn 0.3s ease;
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

.modal-content h3 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #333;
  border-bottom: 2px solid #1da766;
  padding-bottom: 10px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
  color: #555;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}

.form-group textarea {
  height: 80px;
  resize: vertical;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
  padding-top: 15px;
  border-top: 1px solid #eee;
}

.btn-confirm {
  padding: 8px 20px;
  background: #1da766;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  transition: background 0.2s ease;
}

.btn-confirm:hover:not(:disabled) {
  background: #16a085;
}

.btn-confirm:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.btn-cancel {
  padding: 8px 20px;
  background: #6c757d;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background 0.2s ease;
}

.btn-cancel:hover {
  background: #5a6268;
}

.icon {
  font-size: 14px;
}

/* PP-Button Styles */
.pp-cell {
  padding: 4px 8px;
}

.pp-container {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 2px;
}

.pp-btn {
  width: 20px;
  height: 20px;
  border: 1px solid #007bff;
  background: white;
  color: #007bff;
  border-radius: 3px;
  cursor: pointer;
  font-size: 14px;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  line-height: 1;
  padding: 0;
}

.pp-btn:hover:not(:disabled) {
  background: #007bff;
  color: white;
}

.pp-btn:disabled {
  border-color: #ccc;
  color: #ccc;
  cursor: not-allowed;
  opacity: 0.5;
}

.pp-btn-plus {
  border-color: #28a745;
  color: #28a745;
}

.pp-btn-plus:hover:not(:disabled) {
  background: #28a745;
  color: white;
}

.pp-btn-minus {
  border-color: #dc3545;
  color: #dc3545;
}

.pp-btn-minus:hover:not(:disabled) {
  background: #dc3545;
  color: white;
}

.pp-value {
  min-width: 20px;
  text-align: center;
  font-weight: bold;
  color: #495057;
  font-size: 13px;
}


</style>


<script>
import API from '@/utils/api'
import SkillImproveDialog from './SkillImproveDialog.vue'

export default {
  name: "SkillView",
  components: {
    SkillImproveDialog
  },
  props: {
    character: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      learningMode: false,
      showLearnDialog: false,
      showImproveSelectionDialog: false,
      showAddDialog: false,
      showDetailedLearnDialog: false,
      
      // Formulardaten
      newSkillName: '',
      learnNotes: '',
      selectedSkillToImprove: null,
      usePP: 0,
      improveNotes: '',
      addSkillName: '',
      addSkillValue: 0,
      addNotes: '',
      
      // Detailliertes Lernen
      selectedSkillToLearn: null,
      selectedLearningType: 'improve', // 'improve', 'learn', 'spell'
      
      isLoading: false
    };
  },
  created() {
    this.$api = API;
  },
  methods: {
    toggleLearningMode() {
      this.learningMode = !this.learningMode;
      if (!this.learningMode) {
        this.closeDialogs();
      }
    },
    
    showLearnNewDialog() {
      this.closeDialogs();
      this.showLearnDialog = true;
    },
    
    showImproveDialog() {
      this.closeDialogs();
      this.showImproveSelectionDialog = true;
    },
    
    showAddDialog() {
      this.closeDialogs();
      this.showAddDialog = true;
    },
    
    closeDialogs() {
      this.showLearnDialog = false;
      this.showImproveSelectionDialog = false;
      this.showAddDialog = false;
      this.showDetailedLearnDialog = false;
      this.clearFormData();
    },
    
    clearFormData() {
      this.newSkillName = '';
      this.learnNotes = '';
      this.selectedSkillToImprove = null;
      this.usePP = 0;
      this.improveNotes = '';
      this.addSkillName = '';
      this.addSkillValue = 0;
      this.addNotes = '';
      
      // Detailliertes Lernen zurücksetzen
      this.selectedSkillToLearn = null;
      this.selectedLearningType = 'improve';
    },
    
    async learnNewSkill() {
      if (!this.newSkillName.trim()) {
        alert('Bitte geben Sie einen Fertigkeitsnamen ein.');
        return;
      }
      
      this.isLoading = true;
      try {
        const response = await this.$api.post(`/api/characters/${this.character.id}/learn-skill`, {
          name: this.newSkillName.trim(),
          current_level: 0,
          target_level: 1,
          type: 'fertigkeit',
          action: 'lernen',
          reward: 'ep',
          notes: this.learnNotes || `Fertigkeit ${this.newSkillName} über Frontend gelernt`
        });
        
        console.log('Fertigkeit erfolgreich gelernt:', response.data);
        alert(`Fertigkeit "${this.newSkillName}" erfolgreich gelernt!`);
        this.closeDialogs();
        this.$emit('character-updated');
        
      } catch (error) {
        console.error('Fehler beim Lernen der Fertigkeit:', error);
        alert('Fehler beim Lernen der Fertigkeit: ' + (error.response?.data?.error || error.message));
      } finally {
        this.isLoading = false;
      }
    },
    
    async improveSelectedSkill() {
      if (!this.selectedSkillToImprove) {
        alert('Bitte wählen Sie eine Fertigkeit aus.');
        return;
      }
      
      this.isLoading = true;
      try {
        const response = await this.$api.post(`/api/characters/improve-skill`, {
          char_id: this.character.id,
          name: this.selectedSkillToImprove.name,
          current_level: this.selectedSkillToImprove.fertigkeitswert,
          target_level: this.selectedSkillToImprove.fertigkeitswert + 1,
          type: 'skill',
          action: 'improve',
          reward: 'default',
          use_pp: this.usePP || 0,
          notes: this.improveNotes || `Fertigkeit ${this.selectedSkillToImprove.name} über Frontend verbessert`
        });
        
        console.log('Fertigkeit erfolgreich verbessert:', response.data);
        alert(`Fertigkeit "${this.selectedSkillToImprove.name}" erfolgreich verbessert!`);
        this.closeDialogs();
        this.$emit('character-updated');
        
      } catch (error) {
        console.error('Fehler beim Verbessern der Fertigkeit:', error);
        alert('Fehler beim Verbessern der Fertigkeit: ' + (error.response?.data?.error || error.message));
      } finally {
        this.isLoading = false;
      }
    },
    
    async improveSkill(skill) {
      // Detailliertes Lernformular über neue Komponente öffnen
      this.selectedSkillToLearn = skill;
      this.selectedLearningType = 'improve';
      this.showDetailedLearnDialog = true;
    },
    
    async improveWeaponSkill(skill) {
      // Waffenfertigkeit verbessern
      this.isLoading = true;
      try {
        const response = await this.$api.post(`/api/characters/improve-skill`, {
          char_id: this.character.id,
          name: skill.name,
          current_level: skill.fertigkeitswert,
          target_level: skill.fertigkeitswert + 1,
          type: 'skill',
          action: 'improve',
          reward: 'default',
          use_pp: 0,
          notes: `Waffenfertigkeit ${skill.name} direkt aus Tabelle verbessert`
        });
        
        console.log('Waffenfertigkeit erfolgreich verbessert:', response.data);
        alert(`Waffenfertigkeit "${skill.name}" erfolgreich verbessert!`);
        this.$emit('character-updated');
        
      } catch (error) {
        console.error('Fehler beim Verbessern der Waffenfertigkeit:', error);
        alert('Fehler beim Verbessern der Waffenfertigkeit: ' + (error.response?.data?.error || error.message));
      } finally {
        this.isLoading = false;
      }
    },
    
    async addNewSkill() {
      if (!this.addSkillName.trim()) {
        alert('Bitte geben Sie einen Fertigkeitsnamen ein.');
        return;
      }
      
      if (this.addSkillValue < 0 || this.addSkillValue > 20) {
        alert('Der Fertigkeitswert muss zwischen 0 und 20 liegen.');
        return;
      }
      
      // TODO: Hier würde ein API-Endpunkt zum direkten Hinzufügen von Fertigkeiten benötigt
      // Da dieser noch nicht existiert, verwenden wir eine Platzhalter-Implementierung
      
      alert(`Fertigkeit "${this.addSkillName}" mit Wert ${this.addSkillValue} hinzugefügt! (Noch nicht implementiert)`);
      this.closeDialogs();
    },
    
    handleSkillUpdated() {
      // Event-Handler für die neue Dialog-Komponente
      this.$emit('character-updated');
    },
    
    async increasePP(skill) {
      try {
        const response = await this.$api.post(`/api/characters/${this.character.id}/practice-points/add`, {
          skill_name: skill.name,
          amount: 1
        });
        
        // Verwende die Enhanced Response Daten
        const data = response.data;
        if (data.success) {
          // Aktualisiere die lokalen Daten mit den Server-Daten
          if (data.practice_points) {
            // Aktualisiere alle Praxispunkte basierend auf der Server-Antwort
            this.updateLocalPracticePoints(data.practice_points);
          }
          
          // Zeige informative Nachricht über Zauber-Detection
          if (data.is_spell && data.requested_skill !== data.target_skill) {
            console.log(`Zauber erkannt: PP für "${data.requested_skill}" wurde zu "${data.target_skill}" hinzugefügt`);
          }
          
          console.log('Praxispunkt hinzugefügt:', data.message);
          
          // Charakter-Ansicht aktualisieren
          this.$emit('character-updated');
        }
        
      } catch (error) {
        console.error('Fehler beim Hinzufügen von Praxispunkten:', error);
        alert('Fehler beim Hinzufügen von Praxispunkten: ' + (error.response?.data?.error || error.message));
      }
    },
    
    async decreasePP(skill) {
      if ((skill.pp || 0) <= 0) return;
      
      try {
        const response = await this.$api.post(`/api/characters/${this.character.id}/practice-points/use`, {
          skill_name: skill.name,
          amount: 1
        });
        
        // Verwende die Enhanced Response Daten
        const data = response.data;
        if (data.success) {
          // Aktualisiere die lokalen Daten mit den Server-Daten
          if (data.practice_points) {
            // Aktualisiere alle Praxispunkte basierend auf der Server-Antwort
            this.updateLocalPracticePoints(data.practice_points);
          }
          
          // Zeige informative Nachricht über Zauber-Detection
          if (data.is_spell && data.requested_skill !== data.target_skill) {
            console.log(`Zauber erkannt: PP für "${data.requested_skill}" wurde von "${data.target_skill}" verwendet`);
          }
          
          console.log('Praxispunkt entfernt:', data.message);
          
          // Charakter-Ansicht aktualisieren
          this.$emit('character-updated');
        }
        
      } catch (error) {
        console.error('Fehler beim Entfernen von Praxispunkten:', error);
        alert('Fehler beim Entfernen von Praxispunkten: ' + (error.response?.data?.error || error.message));
      }
    },
    
    async increaseWeaponPP(skill) {
      try {
        const response = await this.$api.post(`/api/characters/${this.character.id}/practice-points/add`, {
          skill_name: skill.name,
          amount: 1
        });
        
        // Verwende die Enhanced Response Daten
        const data = response.data;
        if (data.success) {
          // Aktualisiere die lokalen Daten mit den Server-Daten
          if (data.practice_points) {
            // Aktualisiere alle Praxispunkte basierend auf der Server-Antwort
            this.updateLocalPracticePoints(data.practice_points);
          }
          
          console.log('Praxispunkt für Waffenfertigkeit hinzugefügt:', data.message);
          
          // Charakter-Ansicht aktualisieren
          this.$emit('character-updated');
        }
        
      } catch (error) {
        console.error('Fehler beim Hinzufügen von Praxispunkten für Waffenfertigkeit:', error);
        alert('Fehler beim Hinzufügen von Praxispunkten: ' + (error.response?.data?.error || error.message));
      }
    },
    
    async decreaseWeaponPP(skill) {
      if ((skill.pp || 0) <= 0) return;
      
      try {
        const response = await this.$api.post(`/api/characters/${this.character.id}/practice-points/use`, {
          skill_name: skill.name,
          amount: 1
        });
        
        // Verwende die Enhanced Response Daten
        const data = response.data;
        if (data.success) {
          // Aktualisiere die lokalen Daten mit den Server-Daten
          if (data.practice_points) {
            // Aktualisiere alle Praxispunkte basierend auf der Server-Antwort
            this.updateLocalPracticePoints(data.practice_points);
          }
          
          console.log('Praxispunkt für Waffenfertigkeit entfernt:', data.message);
          
          // Charakter-Ansicht aktualisieren
          this.$emit('character-updated');
        }
        
      } catch (error) {
        console.error('Fehler beim Entfernen von Praxispunkten für Waffenfertigkeit:', error);
        alert('Fehler beim Entfernen von Praxispunkten: ' + (error.response?.data?.error || error.message));
      }
    },

    // Helper-Methode zum Aktualisieren der lokalen Praxispunkte basierend auf Server-Response
    updateLocalPracticePoints(practicePointsFromServer) {
      // Erstelle ein Map für schnellen Zugriff
      const ppMap = {};
      practicePointsFromServer.forEach(pp => {
        ppMap[pp.skill_name] = pp.amount;
      });

      // Aktualisiere categorizedskills (Fertigkeiten nach Kategorien)
      if (this.character.categorizedskills) {
        Object.values(this.character.categorizedskills).forEach(skillCategory => {
          if (Array.isArray(skillCategory)) {
            skillCategory.forEach(skill => {
              skill.pp = ppMap[skill.name] || 0;
            });
          }
        });
      }

      // Aktualisiere Waffen-Fertigkeiten
      if (this.character.waffenfertigkeiten) {
        this.character.waffenfertigkeiten.forEach(skill => {
          skill.pp = ppMap[skill.name] || 0;
        });
      }

      // Aktualisiere auch flache fertigkeiten Liste falls vorhanden
      if (this.character.fertigkeiten) {
        this.character.fertigkeiten.forEach(skill => {
          skill.pp = ppMap[skill.name] || 0;
        });
      }
    }
  }
};
</script>
