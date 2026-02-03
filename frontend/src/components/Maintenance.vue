<template>
  <div class="character-details">
    <div v-if="loading">Loading...</div>
    <div v-else>

    <!-- Submenu Content -->
      <component
        :is="currentView"
        :mdata="mdata"
      />
    </div>
    <!-- Submenu -->
    <div class="submenu">
      <button
        v-for="menu in menus"
        :key="menu.id"
        :class="{ active: currentView === menu.component }"
        @click="changeView(menu.component)"
      >
      {{ $t( 'maintmenu.'+ menu.name ) }}
      </button>
    </div>
  </div>
</template>

<style>
</style>


<script>
import API from '../utils/api'
import SkillView from "./maintenance/SkillView.vue"; // Component for character history
import SpellView from "./maintenance/SpellView.vue"; // Component for character history
import EquipmentView from "./maintenance/EquipmentView.vue"; // Component for character equipment
import WeaponView from "./maintenance/WeaponView.vue"; // Component for character history
import WeaponSkillView from "./maintenance/WeaponSkillView.vue"; // Component for character equipment
import BelieveView from "./maintenance/BelieveView.vue"; // Component for believes maintenance
import GameSystemView from "./maintenance/GameSystemView.vue";
import LitSourceView from "./maintenance/LitSourceView.vue";
import MiscLookupView from "./maintenance/MiscLookupView.vue";
import SkillImprovementCostView from "./maintenance/SkillImprovementCostView.vue";


export default {
  name: "Maintenance",
  //props: ["id"], // Receive the route parameter as a prop
  components: {
    SkillView,
    SpellView,
    EquipmentView,
    WeaponView,
    WeaponSkillView,
    BelieveView,
    GameSystemView,
    LitSourceView,
    MiscLookupView,
    SkillImprovementCostView,
  },
  data() {
    return {
      mdata: {
        skills: [],
        skillcategories: [],
        weaponskills: [],
        spells: [],
        spellcategories: [],
        equipment:[],
        weapons: [],
        weaponskills: [],
      },
      loading: true,
      currentView: "SkillView", // Default view
      lastView: "SkillView",
      menus: [
        { id: 0, name: "skill", component: "SkillView" },
        { id: 1, name: "spell", component: "SpellView" },
        { id: 2, name: "equipment", component: "EquipmentView" },
        { id: 3, name: "weapon", component: "WeaponView" },
        { id: 4, name: "weaponskill", component: "WeaponSkillView" },
        { id: 5, name: "believe", component: "BelieveView" },
        { id: 6, name: "gamesystem", component: "GameSystemView" },
        { id: 7, name: "litsource", component: "LitSourceView" },
        { id: 8, name: "misc", component: "MiscLookupView" },
        { id: 9, name: "skillimprovement", component: "SkillImprovementCostView" },

      ],
    };
  },
  async created() {
    try {
      const token = localStorage.getItem('token')
      const response = await API.get(`/api/maintenance`, {
        headers: { Authorization: `Bearer ${token}` },
      })
      this.mdata= response.data
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      this.loading = false
    }
    /*
    this.skills = response.data['skills']
    this.weaponskills = response.data["weaponskills"]
    this.spells = response.data["spells"]
    this.equipment = response.data["equipment"]
    this.weapons = response.data["weapons"]
    */
  },
  methods: {
    changeView(view) {
      this.lastView = this.currentView;
      this.currentView = view;
    }
  },
};
</script>
