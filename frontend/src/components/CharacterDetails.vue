<template>
  <div class="character-details">
    <!-- Character Header -->
    <div class="character-header">
      <h2>Figur: {{ character.name }}</h2>
    </div>

    <!-- Character Overview -->
    <div class="character-overview">
      <div class="character-image">
        <img :src="imageSrc" alt="Character Image" v-if="this.character.image"/>
      </div>
      <div class="character-stats">
        <div class="stat">
          <span>St</span>
          <strong>{{ character.eigenschaften[6].value }}</strong>
        </div>
        <div class="stat">
          <span>Gs</span>
          <strong>{{ character.eigenschaften[1].value }}</strong>
        </div>
        <div class="stat">
          <span>Gw</span>
          <strong>{{ character.eigenschaften[2].value }}</strong>
        </div>
        <div class="stat">
          <span>Ko</span>
          <strong>{{ character.eigenschaften[4].value }}</strong>
        </div>
        <div class="stat">
          <span>In</span>
          <strong>{{ character.eigenschaften[3].value }}</strong>
        </div>
        <div class="stat">
          <span>Zt</span>
          <strong>{{ character.eigenschaften[8].value }}</strong>
        </div>
        <div class="stat">
          <span>Au</span>
          <strong>{{ character.eigenschaften[0].value }}</strong>
        </div>
        <div class="stat">
          <span>pA</span>
          <strong>{{ character.eigenschaften[5].value }}</strong>
        </div>
        <div class="stat">
          <span>Wk</span>
          <strong>{{ character.eigenschaften[7].value }}</strong>
        </div>
        <div class="stat">
          <span>GiT</span>
          <strong>64</strong>
        </div>
        <div class="stat">
          <span>B</span>
          <strong>{{ character.b.max }}</strong>
        </div>
        <div class="stat">
          <span>LP</span>
          <strong>{{ character.lp.max }}</strong>
        </div>
        <div class="stat">
          <span>AP</span>
          <strong>{{ character.ap.max }}</strong>
        </div>
        <div class="stat">
          <span>GG</span>
          <strong>{{ character.bennies.gg }}</strong>
        </div>
        <div class="stat">
          <span>SG</span>
          <strong>{{ character.bennies.sg }}</strong>
        </div>
      </div>
    </div>

    <!-- Character Information -->
    <div class="character-info">
      <p>
        <strong>Aktive Figur?</strong> ✔
        <strong>Aktuelle Kampagne:</strong> Melzindar
      </p>
      <p>
        {{ character.typ }} (xmännlich), Grad:  {{ character.grad }}, Rasse:  {{ character.rasse }}, Heimat: xAlba, Stand:
        xMittelschicht.
      </p>
      <p>
        Hort für Grad 3: 125 GS, für nächsten Grad: 250 GS.
      </p>
      <p>
        <strong>Spezialisierung:</strong> {{ character.spezialisierung }}.
      </p>
      <p>
        Alter: {{ character.alter }},<strong v-if="character.hand='rechts'"> Rechtshänder</strong><strong v-else-if="character.hand='links'"> Linkshänder</strong><strong   v-else> Beidhändig</strong  >, Größe: {{ character.groesse }}cm, Gewicht: {{ character.gewicht }}kg, Gestalt: {{ character.merkmale.groesse }},
        und {{ character.merkmale.breite }}, Augen: {{ character.merkmale.augenfarbe }}, Haare: {{ character.merkmale.haarfarbe }}, Glaube: {{ character.glaube }}.
      </p>
      <p>
        <strong>Merkmale:</strong> {{ character.merkmale.sonstige }}
      </p>
      <p>
        <em>Persönlicher Bonus für</em> Ausdauer 12, Schaden 5, Angriff 2,
        Abwehr 0, Zauber 0, Resistenz 3 / 4.
      </p>
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
</style>


<script>
import API from '../utils/api'

export default {
  name: "CharacterDetails",
  props: ["id"], // Receive the route parameter as a prop
  data() {
    return {
      character: {},
    };
  },
  async created() {
    const token = localStorage.getItem('token')
    const response = await API.get(`/api/characters/${this.id}`, {
      headers: { Authorization: `Bearer ${token}` },
    })
    this.character = response.data
  },
  computed: {
    imageSrc() {
      return this.character.image
        ? `${this.character.image}`
        : "";
    },
  },
};
</script>
