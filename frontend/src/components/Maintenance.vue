<template>
  <div class="character-details">
    Maintenance
    <div class="character-header">
      <h2>{{ $t('maintenance') }}</h2>
    </div>
    <!-- Submenu Content -->
      <component :is="currentView"  :mdata="mdata"/>
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
/*import WeaponSkillView from "./maintenance/WeaponSkillView.vue"; // Component for character equipment
import SpellView from "./maintenance/SpellView.vue"; // Component for character history
import EquipmentView from "./maintenance/EquipmentView.vue"; // Component for character equipment
import WeaponView from "./maintenance/WeaponView.vue"; // Component for character history
*/


export default {
  name: "Maintenance",
  //props: ["id"], // Receive the route parameter as a prop
  components: {
    SkillView,
    /*WeaponSkillView,
    SpellView,
    WeaponView,
    EquipmentView,*/
  },
  data() {
    return {
      mdata: {},
      /*
      skills: {},
      weaponskills: {},
      spells: {},
      equipment: {},
      weapons: {},*/
      currentView: "SkillView", // Default view
      lastView: "SkillView",
      menus: [
        { id: 0, name: "skill", component: "SkillView" },
        /*{ id: 1, name: "weaponskill", component: "WeaponSkillView" },
        { id: 2, name: "spell", component: "SpellView" },
        { id: 3, name: "equipment", component: "EquipmentView" },
        { id: 1, name: "weapon", component: "WeaponView" },
         */
      ],
    };
  },
  async created() {
    const token = localStorage.getItem('token')
    const response = await API.get(`/api/maintenance`, {
      headers: { Authorization: `Bearer ${token}` },
    })
    this.mdata= response.data
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
    },
  },
};
</script>
