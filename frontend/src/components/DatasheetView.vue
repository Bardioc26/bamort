<template>
    <div class="fullwidth-container datasheet-container" v-if="character">
      <!-- Character Overview -->
      <div class="character-overview">
        <div class="character-image">
          <img :src="imageSrc" alt="Character Image" v-if="this.character.image"/>
        </div>
        <div class="character-stats">
          <div class="stat" v-for="(stat, index) in characterStats" :key="index">
            <span>{{ $t(stat.label) }}</span>
            <strong>{{ getStat(stat.path) }}</strong>
          </div>
        </div>
      </div>
      
      <!-- Character Information -->
      <div class="character-info">
        <div class="info-section">
          <p>
            <strong>Aktive Figur?</strong> ✔
            <strong>Aktuelle Kampagne:</strong> Melzindar
          </p>
          <p>
            {{ character.typ || 'x' }} ({{ character.geschlecht || 'x' }}nännlich),
            Grad:  {{ character.grad || 'x' }},
            Rasse:  {{ character.rasse || 'x' }},
            Heimat: {{ character.heimat || 'x' }}Alba,
            Stand:  {{ character.heimat || 'x' }}Mittelschicht.
          </p>
          <p v-if="character.rasse==='Zwerg'">
            Hort für Grad {{ character.grad || 'x' }}: 125 GS, für nächsten Grad: 250 GS.
          </p>
          <p>
            <strong>Spezialisierung:</strong> {{ character.spezialisierung || '-'}}.
          </p>
          <p>
            Alter: {{ character.alter || 'xx' }},
            <strong v-if="character.hand=='rechts'"> Rechtshänder</strong>
            <strong v-else-if="character.hand=='links'"> Linkshänder</strong>
            <strong   v-else> Beidhändig</strong>,
            Größe: {{ character.groesse }}cm,
            Gewicht: {{ character.gewicht }}kg,
            Gestalt: {{ character.merkmale?.groesse || '-'}}
            und {{ character.merkmale?.breite  || '-'}},
            Augen: {{ character.merkmale?.augenfarbe || '-' }},
            Haare: {{ character.merkmale?.haarfarbe || '-' }},
            Glaube: {{ character.glaube }}.
          </p>
          <p>
            <strong>Merkmale:</strong> {{ character.merkmale?.sonstige || '-' }}
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
/* DatasheetView spezifische Styles */
.datasheet-container {
  padding-top: 10px; /* Reduziertes oberes Padding */
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
  margin-top: 0; /* Kein zusätzlicher oberer Margin */
}

.character-info {
  margin-top: 20px;
}
</style>


<script>
export default {
  name: "DatasheetView",
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
        : "";
    },
  },
  data() {
    return {
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
    getStat(path) {
      if (path === 'git' ){
        return '64!'
      }
      return path.split('.').reduce((obj, key) => obj?.[key], this.character) ?? '-'
    }
  }
};
</script>
