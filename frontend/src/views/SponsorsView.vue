<template>
  <div class="fullwidth-page">
    <div class="page-header" style="flex-direction: row;">
      <button @click="$router.back()" class="btn btn-secondary back-button">
        ← {{ $t('common.back') }}
      </button>
      <h2 style="padding-top: 8px; padding-left: 10px;">{{ $t('sponsors.title') }}</h2>
    </div>

    <div class="card">
      <p>{{ $t('sponsors.introduction') }}</p>
    </div>

    <div class="section-header">
      <h3>{{ $t('sponsors.contributors') }}</h3>
    </div>

    <div class="grid-container grid-2-columns">
      <div class="card" v-for="contributor in contributors" :key="contributor.name">
        <h4>{{ contributor.name }}</h4>
        <p>{{ contributor.role }}</p>
        <p v-if="contributor.contribution">{{ contributor.contribution }}</p>
      </div>
    </div>

    <div class="section-header">
      <h3>{{ $t('sponsors.specialThanks') }}</h3>
    </div>

    <div class="card">
      <ul>
        <li v-for="thanks in specialThanks" :key="thanks">{{ thanks }}</li>
      </ul>
    </div>

    <div class="section-header">
      <h3>{{ $t('sponsors.support') }}</h3>
    </div>

    <div class="card">
      <p>{{ $t('sponsors.supportText') }}</p>
      <div style="display: flex; gap: 10px; margin-top: 15px;">
        <a :href="githubUrl" target="_blank" rel="noopener noreferrer" class="btn btn-primary">
          {{ $t('sponsors.github') }}
        </a>
        <a :href="koFiUrl" target="_blank" rel="noopener noreferrer" class="btn btn-primary">
          {{ $t('sponsors.koFi') }}
        </a>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "SponsorsView",
  data() {
    return {
      githubUrl: "https://github.com/Bardioc26/bamort",
      koFiUrl: "https://ko-fi.com/bardioc26",
      contributors: [
        {
          name: "Bardioc26",
          role: this.$t('sponsors.roleCreator'),
          contribution: this.$t('sponsors.contributionCreator')
        }
      ],
      specialThanks: []
    }
  },
  created() {
    // Populate special thanks with translations
    this.specialThanks = [
      this.$t('sponsors.thanksOpenSource'),
      this.$t('sponsors.thanksVue'),
      this.$t('sponsors.thanksGo'),
      this.$t('sponsors.thanksCommunity')
    ]
  }
}
</script>
