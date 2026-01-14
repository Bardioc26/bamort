<template>
  <div class="fullwidth-page">
    <div class="page-header" style="flex-direction: row;">
      <button @click="$router.back()" class="btn btn-secondary back-button">
        ← {{ $t('common.back') }}
      </button>
      <h2 style="padding-top: 8px; padding-left: 10px;">{{ $t('help.title') }}</h2>
    </div>

    <div class="card">
      <p>{{ $t('help.introduction') }}</p>
    </div>

    <div class="section-header">
      <h3>{{ $t('help.gettingStarted') }}</h3>
    </div>

    <div class="card">
      <h4>{{ $t('help.step1Title') }}</h4>
      <p>{{ $t('help.step1Text') }}</p>
      
      <h4 style="margin-top: 20px;">{{ $t('help.step2Title') }}</h4>
      <p>{{ $t('help.step2Text') }}</p>
      
      <h4 style="margin-top: 20px;">{{ $t('help.step3Title') }}</h4>
      <p>{{ $t('help.step3Text') }}</p>
    </div>

    <div class="section-header">
      <h3>{{ $t('help.features') }}</h3>
    </div>

    <div class="grid-container grid-2-columns">
      <div class="card">
        <h4>{{ $t('help.featureCharacterManagement') }}</h4>
        <p>{{ $t('help.featureCharacterManagementText') }}</p>
      </div>

      <div class="card">
        <h4>{{ $t('help.featureSkills') }}</h4>
        <p>{{ $t('help.featureSkillsText') }}</p>
      </div>

      <div class="card">
        <h4>{{ $t('help.featureSpells') }}</h4>
        <p>{{ $t('help.featureSpellsText') }}</p>
      </div>

      <div class="card">
        <h4>{{ $t('help.featureEquipment') }}</h4>
        <p>{{ $t('help.featureEquipmentText') }}</p>
      </div>

      <div class="card">
        <h4>{{ $t('help.featureExport') }}</h4>
        <p>{{ $t('help.featureExportText') }}</p>
      </div>

      <div class="card">
        <h4>{{ $t('help.featureImport') }}</h4>
        <p>{{ $t('help.featureImportText') }}</p>
      </div>
    </div>

    <div class="section-header">
      <h3>{{ $t('help.faq') }}</h3>
    </div>

    <div class="card" v-for="(faq, index) in faqs" :key="index">
      <h4>{{ faq.question }}</h4>
      <p>{{ faq.answer }}</p>
    </div>

    <div class="section-header">
      <h3>{{ $t('help.support') }}</h3>
    </div>

    <div class="card">
      <p>{{ $t('help.supportText') }}</p>
      <div style="display: flex; gap: 10px; margin-top: 15px;">
        <a :href="githubUrl" target="_blank" rel="noopener noreferrer" class="btn btn-primary">
          {{ $t('help.github') }}
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
  name: "HelpView",
  data() {
    return {
      githubUrl: "https://github.com/Bardioc26/bamort",
      koFiUrl: "https://ko-fi.com/bardioc26"
    }
  },
  computed: {
    faqs() {
      const faqs = []
      let i = 1
      // Dynamically load FAQs until we can't find any more
      while (this.$te(`help.faq${i}.question`)) {
        faqs.push({
          question: this.$t(`help.faq${i}.question`),
          answer: this.$t(`help.faq${i}.answer`)
        })
        i++
      }
      return faqs
    }
  }
}
</script>
