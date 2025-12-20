<template>
  <div class="character-details">
    <!-- Character Header -->
    <div class="character-header">
      <div class="header-content">
        <button @click="showExportDialog = true" class="export-button-small" :title="$t('export.exportPDF')">
          📄
        </button>
        <h2>{{ $t('char') }}: {{ character.name }} ({{ $t(currentView) }})</h2>
      </div>
    </div>

    <!-- Export Dialog -->
    <div v-if="showExportDialog" class="modal-overlay" @click.self="showExportDialog = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ $t('export.exportPDF') }}</h3>
          <button @click="showExportDialog = false" class="close-button">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>{{ $t('export.selectTemplate') }}:</label>
            <select v-model="selectedTemplate" class="template-select">
              <option value="">{{ $t('export.pleaseSelectTemplate') }}</option>
              <option v-for="template in templates" :key="template.id" :value="template.id">
                {{ template.name }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="showUserName">
              {{ $t('export.showUserName') }}
            </label>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="showExportDialog = false" class="btn-cancel">
            {{ $t('export.cancel') }}
          </button>
          <button @click="exportToPDF" class="btn-export" :disabled="!selectedTemplate || isExporting">
            <span v-if="!isExporting">{{ $t('export.export') }}</span>
            <span v-else>{{ $t('export.exporting') }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Submenu Content -->
    <!-- <div class="character-aspect"> -->
      <component :is="currentView" :character="character" @character-updated="refreshCharacter"/>
    <!-- </div> -->

    <!-- Submenu -->
    <div class="submenu">
      <button
        v-for="menu in menus"
        :key="menu.id"
        :class="{ active: currentView === menu.component }"
        @click="changeView(menu.component)"
      >
      {{ $t( 'menu.'+ menu.name ) }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.character-details {
  width: 100%;
  height: 100%;
  padding: 20px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

.character-header {
  margin-bottom: 20px;
}

.header-content {
  display: flex;
  align-items: center;
  gap: 15px;
}

.export-button-small {
  width: 40px;
  height: 40px;
  padding: 0;
  border: 1px solid #007bff;
  border-radius: 8px;
  background: #007bff;
  color: white;
  font-size: 1.2rem;
  cursor: pointer;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.export-button-small:hover {
  background: #0056b3;
  border-color: #0056b3;
  transform: scale(1.05);
}

.character-header h2 {
  margin: 0;
  color: #333;
  font-size: 1.5rem;
  border-bottom: 2px solid #007bff;
  padding-bottom: 10px;
  flex: 1;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #dee2e6;
}

.modal-header h3 {
  margin: 0;
  color: #333;
  font-size: 1.25rem;
}

.close-button {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #999;
  cursor: pointer;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-button:hover {
  color: #333;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #495057;
}

.template-select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  background: white;
  color: #495057;
  font-size: 0.95rem;
  cursor: pointer;
}

.template-select:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  user-select: none;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 20px;
  border-top: 1px solid #dee2e6;
}

.btn-cancel {
  padding: 10px 20px;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  background: #f8f9fa;
  color: #495057;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s ease;
}

.btn-cancel:hover {
  background: #e9ecef;
  border-color: #adb5bd;
}

.btn-export {
  padding: 10px 20px;
  border: 1px solid #007bff;
  border-radius: 6px;
  background: #007bff;
  color: white;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s ease;
}

.btn-export:hover:not(:disabled) {
  background: #0056b3;
  border-color: #0056b3;
}

.btn-export:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.submenu {
  display: flex;
  gap: 10px;
  margin: 20px 0;
  flex-wrap: wrap;
}

.submenu button {
  padding: 10px 16px;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  background: #f8f9fa;
  color: #495057;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s ease;
}

.submenu button:hover {
  background: #e9ecef;
  border-color: #007bff;
}

.submenu button.active {
  background: #007bff;
  color: white;
  border-color: #007bff;
}
</style>


<script>
import API from '../utils/api'
import DatasheetView from "./DatasheetView.vue"; // Component for character stats
import SkillView from "./SkillView.vue"; // Component for character history
import WeaponView from "./WeaponView.vue"; // Component for character history
import SpellView from "./SpellView.vue"; // Component for character history
import EquipmentView from "./EquipmentView.vue"; // Component for character equipment
import ExperianceView from "./ExperianceView.vue"; // Component for character history
import DeleteCharView from "./DeleteCharView.vue"; // Component for character history


export default {
  name: "CharacterDetails",
  props: ["id"], // Receive the route parameter as a prop
  components: {
    DatasheetView,
    SkillView,
    WeaponView,
    SpellView,
    EquipmentView,
    ExperianceView,
    DeleteCharView,
  },
  data() {
    return {
      character: {},
      currentView: "DatasheetView", // Default view
      lastView: "DatasheetView",
      templates: [],
      selectedTemplate: "",
      showUserName: false,
      showExportDialog: false,
      isExporting: false,
      menus: [
        { id: 1, name: "Datasheet", component: "DatasheetView" },
        { id: 2, name: "Skill", component: "SkillView" },
        { id: 3, name: "Weapon", component: "WeaponView" },
        { id: 4, name: "Spell", component: "SpellView" },
        { id: 5, name: "Equipment", component: "EquipmentView" },
        { id: 6, name: "Experiance", component: "ExperianceView" },
        { id: 6, name: "DeleteChar", component: "DeleteCharView" },

        //{ id: 3, name: "History", component: "HistoryView" },
        //{ id: 2, name: "Notes", component: "NotesView" },
        //{ id: 2, name: "Campagne", component: "CampagneView" },
      ],
    };
  },
  async created() {
    const token = localStorage.getItem('token')
    const response = await API.get(`/api/characters/${this.id}`, {
      headers: { Authorization: `Bearer ${token}` },
    })
    this.character = response.data
    
    // Load available templates
    await this.loadTemplates()
  },
  methods: {
    async loadTemplates() {
      try {
        const response = await API.get('/api/pdf/templates')
        this.templates = response.data
        // Auto-select first template if available
        if (this.templates.length > 0) {
          this.selectedTemplate = this.templates[0].id
        }
      } catch (error) {
        console.error('Failed to load templates:', error)
      }
    },
    
    async exportToPDF() {
      if (!this.selectedTemplate) {
        alert(this.$t('export.pleaseSelectTemplate'))
        return
      }
      
      this.isExporting = true
      try {
        const params = { template: this.selectedTemplate }
        if (this.showUserName) {
          params.showUserName = true
        }
        
        const response = await API.get(`/api/pdf/export/${this.id}`, {
          params,
          responseType: 'blob'
        })
        
        // Create blob URL and open in new tab
        const blob = new Blob([response.data], { type: 'application/pdf' })
        const url = window.URL.createObjectURL(blob)
        window.open(url, '_blank')
        
        // Clean up blob URL after a delay
        setTimeout(() => window.URL.revokeObjectURL(url), 100)
        
        // Close dialog on success
        this.showExportDialog = false
      } catch (error) {
        console.error('Failed to export PDF:', error)
        alert(this.$t('export.exportFailed') + ': ' + (error.response?.data?.error || error.message))
      } finally {
        this.isExporting = false
      }
    },
    
    changeView(view) {
      this.lastView = this.currentView;
      this.currentView = view;
    },
    
    async refreshCharacter() {
      // Lade die Charakterdaten neu nach einer Aktualisierung
      try {
        const token = localStorage.getItem('token');
        const response = await API.get(`/api/characters/${this.id}`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        this.character = response.data;
        console.log('Character data refreshed after skill update');
      } catch (error) {
        console.error('Failed to refresh character data:', error);
        // Optional: Zeige eine Fehlermeldung an
        alert('Fehler beim Aktualisieren der Charakterdaten: ' + (error.response?.data?.error || error.message));
      }
    },
  },
};
</script>
