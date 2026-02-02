<template>
  <div class="character-details">
    <!-- Character Header -->
    <div class="character-header">
      <div class="header-content">
        <button @click="showExportDialog = true" class="export-button-small" :title="$t('export.title')">
          📄
        </button>
        <button v-if="isOwner" @click="showVisibilityDialog = true" class="export-button-small" :title="$t('visibility.title')">
          {{ character.public ? '🌐' : '🔒' }}
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

    <!-- Visibility Dialog -->
    <VisibilityDialog 
      :characterId="id" 
      :currentVisibility="character.public"
      :showDialog="showVisibilityDialog"
      @update:showDialog="showVisibilityDialog = $event"
      @visibility-updated="handleVisibilityUpdated"
    />

    <!-- Submenu Content -->
    <!-- <div class="character-aspect"> -->
      <component :is="currentView" :character="character" :isOwner="isOwner" @character-updated="refreshCharacter"/>
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

<style>
/* Component-specific styles only - global styles in main.css */
.character-details {
  width: 100%;
  height: 100%;
  padding: 20px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}
</style>


<script>
import API from '../utils/api'
import { useUserStore } from '../stores/userStore'
import ExportDialog from "./ExportDialog.vue";
import VisibilityDialog from "./VisibilityDialog.vue";
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
    VisibilityDialog,
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
      showVisibilityDialog: false,
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
  computed: {
    isOwner() {
      const userStore = useUserStore()
      return userStore.currentUser && this.character.user_id === userStore.currentUser.id
    }
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
    
    handleVisibilityUpdated(isPublic) {
      this.character.public = isPublic
      this.showVisibilityDialog = false
      console.log('Character visibility updated to:', isPublic ? 'public' : 'private')
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