<template>
  <div class="experiance-view">
    <h3>{{ $t('experience.title') }}</h3>
    
    <!-- Erfahrungspunkte -->
    <div class="experience-section">
      <h4>{{ $t('experience.experience_points') }}</h4>
      <div class="stat-box">
        <div class="stat-item">
          <span class="stat-label">{{ $t('experience.available_ep') }}:</span>
          <span class="stat-value">{{ character.erfahrungsschatz?.value || 0 }} EP</span>
        </div>
      </div>
    </div>

    <!-- Vermögen -->
    <div class="wealth-section">
      <h4>{{ $t('experience.wealth') }}</h4>
      <div class="stat-box">
        <div class="wealth-item">
          <span class="wealth-label">{{ $t('experience.gold_coins') }}:</span>
          <span class="wealth-value">{{ character.vermoegen?.goldstücke || 0 }} GS</span>
        </div>
        <div class="wealth-item">
          <span class="wealth-label">{{ $t('experience.silver_coins') }}:</span>
          <span class="wealth-value">{{ character.vermoegen?.silberstücke || 0 }} SS</span>
        </div>
        <div class="wealth-item">
          <span class="wealth-label">{{ $t('experience.copper_coins') }}:</span>
          <span class="wealth-value">{{ character.vermoegen?.kupferstücke || 0 }} KS</span>
        </div>
        <div class="wealth-item total">
          <span class="wealth-label">{{ $t('experience.total_in_gs') }}:</span>
          <span class="wealth-value">{{ totalWealthInGS }} GS</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.experiance-view {
  padding: 20px;
  max-width: 800px;
}

.experiance-view h3 {
  color: #333;
  border-bottom: 2px solid #007bff;
  padding-bottom: 10px;
  margin-bottom: 20px;
}

.experiance-view h4 {
  color: #555;
  margin-bottom: 15px;
  margin-top: 25px;
}

.stat-box {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 20px;
}

.stat-item, .wealth-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #e9ecef;
}

.stat-item:last-child, .wealth-item:last-child {
  border-bottom: none;
}

.wealth-item.total {
  border-top: 2px solid #007bff;
  margin-top: 10px;
  padding-top: 15px;
  font-weight: bold;
  color: #007bff;
}

.stat-label, .wealth-label {
  font-weight: 500;
  color: #555;
}

.stat-value, .wealth-value {
  font-weight: bold;
  color: #333;
  background: #fff;
  padding: 5px 10px;
  border-radius: 4px;
  border: 1px solid #ddd;
}

.wealth-item.total .wealth-value {
  background: #007bff;
  color: white;
  border-color: #007bff;
}

.experience-section, .wealth-section {
  margin-bottom: 30px;
}
</style>


<script>
export default {
  name: "ExperianceView",
  props: {
    character: {
      type: Object,
      required: true
    }
  },
  computed: {
    totalWealthInGS() {
      const vermoegen = this.character.vermoegen || {};
      const goldstücke = vermoegen.goldstücke || 0;    // GS
      const silberstücke = vermoegen.silberstücke || 0; // SS
      const kupferstücke = vermoegen.kupferstücke || 0; // KS
      
      // Midgard Währungsumrechnung: 1 GS = 10 SS = 10 KS
      // Alles in Goldstücke umrechnen
      return goldstücke + Math.floor(silberstücke / 10) + Math.floor(kupferstücke / 10);
    }
  }
};
</script>
