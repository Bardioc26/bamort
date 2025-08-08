<template>
  <div class="fullwidth-container cd-view">
    <!-- Header mit Lernmodus-Kontrollen -->
    <div class="page-header header-section">
      <h2>{{ character.name }}'s Zauber</h2>
      
      <div class="learning-mode-controls">
        <!-- Lernmodus Toggle Button -->
        <button 
          @click="showLearnNewDialog"
          class="btn btn-primary btn-learning-mode"
          title="Neuen Zauber lernen"
        >
          <span class="icon">🎓</span>
        </button>
      </div>
    </div>

    <div class="cd-list">
        <table class="cd-table">
      <thead>
        <tr class="cd-table-header">
          <th>{{ $t('spell.name') }}</th>
          <th>{{ $t('spell.description') }}</th>
          <th>{{ $t('spell.bonus') }}</th>
          <th>{{ $t('spell.quelle') }}</th>
        </tr>
      </thead>
      <tbody>
      <template v-for="spell in character.zauber" :key="spell.id || spell.name">
        <tr>
          <td>{{ spell.name || '-' }}</td>
          <td>{{ spell.beschreibung || '-' }}</td>
          <td>{{ spell.bonus || '0' }}</td>
          <td>{{ spell.quelle || '-' }}</td>
        </tr>
      </template>
      </tbody>
      </table>
    </div> <!--- end cd-list-->

    <!-- Dialog für neue Zauber lernen -->
    <SpellLearnDialog 
      :character="character"
      :show="showLearnDialog"
      @close="closeDialogs"
      @spell-learned="handleSpellLearned"
    />
  </div> <!--- end character -datasheet-->

</template>

<style>
/* SpellView spezifische Styles */
.cd-table-header {
  background-color: #1da766;
  color: white;
  font-weight: bold;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.learning-mode-controls {
  display: flex;
  align-items: center;
  gap: 15px;
}

.btn-learning-mode {
  display: flex;
  align-items: center;
  gap: 5px;
  border: 2px solid #1da766;
}

.btn-learning-mode:hover {
  background: #1da766;
  color: white;
}

.icon {
  font-size: 1.2em;
}
</style>


<script>
import API from '@/utils/api'
import SpellLearnDialog from './SpellLearnDialog.vue'

export default {
  name: "SpellView",
  components: {
    SpellLearnDialog
  },
  props: {
    character: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      showLearnDialog: false,
      isLoading: false
    };
  },
  created() {
    this.$api = API;
  },
  methods: {
    showLearnNewDialog() {
      this.showLearnDialog = true;
    },
    
    closeDialogs() {
      this.showLearnDialog = false;
    },
    
    handleSpellLearned(eventData) {
      this.$emit('character-updated');
      this.closeDialogs();
    }
  }
};
</script>
