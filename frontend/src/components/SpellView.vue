<template>
  <div class="cd-view">
    <!-- Header mit Lernmodus-Kontrollen -->
    <div class="header-section">
      <h2>{{ character.name }}'s Zauber</h2>
      
      <div class="learning-mode-controls">
        <!-- Lernmodus Toggle Button -->
        <button 
          @click="showLearnNewDialog"
          class="btn-learning-mode"
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
.cd-table {
  width: 100%;
}

.cd-table-header {
  background-color: #1da766;
  color: white;
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
