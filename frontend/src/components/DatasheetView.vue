<template>
    <div class="datasheet-view">
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
    </div> <!--- end character -datasheet-->

</template>

<style>
.datasheet-view {
  color: #fff;
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
  }
};
</script>
