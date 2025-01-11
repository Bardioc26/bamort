import { defineStore } from 'pinia'
import { createI18n } from 'vue-i18n'
import de from '@/locales/de'
import en from '@/locales/en'

export const i18n = createI18n({
  legacy: false,
  locale: localStorage.getItem('language') || 'de',
  fallbackLocale: 'en',
  messages: { de, en }
})

export const useLanguageStore = defineStore('language', {
  state: () => ({
    currentLanguage: localStorage.getItem('language') || 'de'
  }),
  actions: {
    setLanguage(lang) {
      this.currentLanguage = lang
      i18n.global.locale.value = lang
      localStorage.setItem('language', lang)
    }
  }
})
