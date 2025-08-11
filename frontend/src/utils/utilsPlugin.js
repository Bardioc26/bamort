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
 * this.$rollDice(3, 6) // 3 Würfel mit 6 Seiten
 * this.$rollNotation('2d10+3') // RPG-Würfelnotation
 */

import { formatDate, formatDateTime, formatRelativeDate, safeValue, capitalize } from './dateUtils'
import { 
  rollDie, 
  rollDice, 
  rollDiceWithSum, 
  rollNotation, 
  randomBetween, 
  randomChoice, 
  shuffleArray
} from './randomUtils.js'

export default {
  install(app) {
    // Globale Properties für Vue 3 - Date Utils
    app.config.globalProperties.$formatDate = formatDate
    app.config.globalProperties.$formatDateTime = formatDateTime
    app.config.globalProperties.$formatRelativeDate = formatRelativeDate
    app.config.globalProperties.$safeValue = safeValue
    app.config.globalProperties.$capitalize = capitalize
    
    // Globale Properties für Vue 3 - Random Utils
    app.config.globalProperties.$rollDie = rollDie
    app.config.globalProperties.$rollDice = rollDice
    app.config.globalProperties.$rollDiceWithSum = rollDiceWithSum
    app.config.globalProperties.$rollNotation = rollNotation
    app.config.globalProperties.$randomBetween = randomBetween
    app.config.globalProperties.$randomChoice = randomChoice
    app.config.globalProperties.$shuffleArray = shuffleArray
    
    // Provide/Inject für Composition API
    app.provide('utils', {
      // Date Utils
      formatDate,
      formatDateTime,
      formatRelativeDate,
      safeValue,
      capitalize,
      // Random Utils
      rollDie,
      rollDice,
      rollDiceWithSum,
      rollNotation,
      randomBetween,
      randomChoice,
      shuffleArray
    })
  }
}
