<template>
  <div class="skill-view">
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
      </div>
    </div>

  </div> <!--- end character -datasheet-->

</template>

<style>
.skill-view {
color: #055;
}
</style>


<script>
export default {
  name: "SkillView",
  props: {
    character: {
      type: Object,
      required: true
    }
  }
};
</script>
