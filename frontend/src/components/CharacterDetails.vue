<template>
  <div class="character-details">
    <!-- Character Header -->
    <div class="character-header">
      <h2>{{ $t('char') }}: {{ character.name }}</h2>
    </div>
    <!-- Submenu Content -->
    <!-- <div class="character-aspect"> -->
      <component :is="currentView"  :character="character"/>
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
.character-details {
  background-color: #444; /* Background color */
  color: #fff; /* Text color */
  padding: 20px;
  border-radius: 8px;
  width: 90%;
  margin: 0 auto;
  font-family: Arial, sans-serif;
}

.character-header h2 {
  font-size: 1.5rem;
  text-align: center;
  color: #ddd;
  margin-bottom: 20px;
}

.character-overview {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
}

.character-image img {
  width: 150px;
  height: auto;
  border-radius: 8px;
  border: 2px solid #333;
}

.character-stats {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 10px;
  width: 100%;
}

.stat {
  background-color: #555;
  border: 1px solid #333;
  text-align: center;
  padding: 10px;
  border-radius: 5px;
  font-size: 0.9rem;
  font-weight: bold;
}

.stat span {
  display: block;
  font-size: 0.8rem;
  color: #aaa;
}

.character-info {
  background-color: #333;
  padding: 15px;
  border-radius: 8px;
  line-height: 1.6;
  white-space: nowrap;
}

.character-info p {
  margin: 10px 0;
}

.character-info strong {
  color: #eee;
}

.character-info em {
  font-style: italic;
  color: #ccc;
}

.character-details {
  position: relative;
  background-color: #444;
  color: #fff;
  padding: 20px;
  border-radius: 8px;
  width: 90%;
  margin: 0 auto;
  font-family: Arial, sans-serif;
  min-height: 400px; /* Ensure there's space for content */
}

.character-main {
  margin-bottom: 20px;
}

.character-aspect {
  padding: 20px;
  background-color: #333;
  border-radius: 8px;
  min-height: 200px; /* Space for content */
}

.submenu {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background-color: #222;
  display: flex;
  justify-content: center;
  padding: 10px;
  border-top: 1px solid #555;
}

.submenu button {
  background-color: #555;
  color: #fff;
  border: none;
  padding: 10px 20px;
  margin: 0 10px;
  border-radius: 4px;
  font-size: 0.9rem;
  cursor: pointer;
  transition: background-color 0.3s;
}

.submenu button:hover {
  background-color: #777;
}

.submenu button.active {
  background-color: #007BFF;
  color: #fff;
  font-weight: bold;
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
  },
  data() {
    return {
      character: {},
      currentView: "DatasheetView", // Default view
      menus: [
        { id: 1, name: "Datasheet", component: "DatasheetView" },
        { id: 2, name: "Skill", component: "SkillView" },
        { id: 2, name: "Weapon", component: "WeaponView" },
        { id: 2, name: "Spell", component: "SpellView" },
        { id: 2, name: "Equipment", component: "EquipmentView" },
        { id: 2, name: "Experiance", component: "ExperianceView" },
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
    changeView(view) {
      this.currentView = view;
    },
  },
};
</script>
