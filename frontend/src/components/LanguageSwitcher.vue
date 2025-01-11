<template>
  <select v-model="selectedLanguage">
    <option value="de">Deutsch</option>
    <option value="en">English</option>
  </select>
</template>

<script>
import { useLanguageStore } from '../stores/languageStore'

export default {
  name: 'LanguageSwitcher',
  computed: {
    selectedLanguage: {
      get() {
        return useLanguageStore().currentLanguage
      },
      set(value) {
        useLanguageStore().setLanguage(value)
      }
    }
  }
}
</script>
