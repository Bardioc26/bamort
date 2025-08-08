/**
 * Gemeinsame Utility-Funktionen für die gesamte Anwendung
 */

/**
 * Formatiert ein Datum-String in ein lokales Datumsformat
 * @param {string} dateString - ISO-Datum-String
 * @param {string} locale - Sprach-/Ländercode (optional, default: browser locale)
 * @param {Object} options - Intl.DateTimeFormat Optionen (optional)
 * @returns {string} - Formatiertes Datum
 */
export function formatDate(dateString, locale = undefined, options = {}) {
  if (!dateString) {
    return 'Unknown'
  }
  
  try {
    const date = new Date(dateString)
    
    // Prüfe ob das Datum gültig ist
    if (isNaN(date.getTime())) {
      return 'Invalid Date'
    }
    
    // Default-Optionen für Datumsformatierung
    const defaultOptions = {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      ...options
    }
    
    return date.toLocaleDateString(locale, defaultOptions)
  } catch (error) {
    console.warn('Error formatting date:', error)
    return 'Invalid Date'
  }
}

/**
 * Formatiert ein Datum-String mit Uhrzeit
 * @param {string} dateString - ISO-Datum-String
 * @param {string} locale - Sprach-/Ländercode (optional)
 * @returns {string} - Formatiertes Datum mit Uhrzeit
 */
export function formatDateTime(dateString, locale = undefined) {
  return formatDate(dateString, locale, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

/**
 * Formatiert ein Datum relativ (z.B. "vor 2 Stunden")
 * @param {string} dateString - ISO-Datum-String
 * @param {string} locale - Sprach-/Ländercode (optional)
 * @returns {string} - Relatives Datum
 */
export function formatRelativeDate(dateString, locale = undefined) {
  if (!dateString) {
    return 'Unknown'
  }
  
  try {
    const date = new Date(dateString)
    const now = new Date()
    
    if (isNaN(date.getTime())) {
      return 'Invalid Date'
    }
    
    // Für moderne Browser mit Intl.RelativeTimeFormat
    if (typeof Intl !== 'undefined' && Intl.RelativeTimeFormat) {
      const rtf = new Intl.RelativeTimeFormat(locale, { numeric: 'auto' })
      const diffTime = date.getTime() - now.getTime()
      const diffDays = Math.round(diffTime / (1000 * 60 * 60 * 24))
      
      if (Math.abs(diffDays) < 1) {
        const diffHours = Math.round(diffTime / (1000 * 60 * 60))
        if (Math.abs(diffHours) < 1) {
          const diffMinutes = Math.round(diffTime / (1000 * 60))
          return rtf.format(diffMinutes, 'minute')
        }
        return rtf.format(diffHours, 'hour')
      }
      
      return rtf.format(diffDays, 'day')
    }
    
    // Fallback für ältere Browser
    return formatDate(dateString, locale)
  } catch (error) {
    console.warn('Error formatting relative date:', error)
    return formatDate(dateString, locale)
  }
}

/**
 * Weitere gemeinsame Utility-Funktionen können hier hinzugefügt werden
 */

/**
 * Sicherheits-Check für leere Werte
 * @param {*} value - Zu prüfender Wert
 * @param {*} fallback - Fallback-Wert
 * @returns {*} - Wert oder Fallback
 */
export function safeValue(value, fallback = '-') {
  return value != null && value !== '' ? value : fallback
}

/**
 * Kapitalisiert den ersten Buchstaben eines Strings
 * @param {string} str - Input String
 * @returns {string} - String mit großem ersten Buchstaben
 */
export function capitalize(str) {
  if (!str || typeof str !== 'string') return str
  return str.charAt(0).toUpperCase() + str.slice(1)
}
