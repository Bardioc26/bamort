<template>
  <div class="character-details">
    <!-- Character Header -->
    <div class="character-header">
      <div class="header-content">
        <button @click="showExportDialog = true" class="export-button-small" :title="$t('export.title')">
          📄
        </button>
        <h2>{{ $t('char') }}: {{ character.name }} ({{ $t(currentView) }})</h2>
      </div>
    </div>

    <!-- Export Dialog -->
    <ExportDialog 
      :characterId="id" 
      :showDialog="showExportDialog"
      @update:showDialog="showExportDialog = $event"
      @export-success="handleExportSuccess"
    />

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
import ExportDialog from "./ExportDialog.vue";
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
    ExportDialog,
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
      showExportDialog: false,
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
  },
  methods: {
    handleExportSuccess() {
      console.log('PDF exported successfully')
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