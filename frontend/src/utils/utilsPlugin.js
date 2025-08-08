/**
 * Vue Plugin für globale Utility-Funktionen
 * 
 * Usage in main.js:
 * import UtilsPlugin from './utils/utilsPlugin'
 * app.use(UtilsPlugin)
 * 
 * Usage in components:
 * this.$formatDate(dateString)
 * this.$safeValue(value, 'fallback')
 */

import { formatDate, formatDateTime, formatRelativeDate, safeValue, capitalize } from './dateUtils'

export default {
  install(app) {
    // Globale Properties für Vue 3
    app.config.globalProperties.$formatDate = formatDate
    app.config.globalProperties.$formatDateTime = formatDateTime
    app.config.globalProperties.$formatRelativeDate = formatRelativeDate
    app.config.globalProperties.$safeValue = safeValue
    app.config.globalProperties.$capitalize = capitalize
    
    // Provide/Inject für Composition API
    app.provide('utils', {
      formatDate,
      formatDateTime,
      formatRelativeDate,
      safeValue,
      capitalize
    })
  }
}
