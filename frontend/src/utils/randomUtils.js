/**
 * Utility-Funktionen für Zufallszahlen
 */

/**
 * Erzeugt eine einzelne Zufallszahl im Bereich von 1 bis max
 * @param {number} max - Maximaler Wert (inklusive)
 * @returns {number} - Zufallszahl zwischen 1 und max
 */
export function rollDie(max = 6) {
  if (typeof max !== 'number' || max < 1) {
    throw new Error('Max value must be a positive number')
  }
  return Math.floor(Math.random() * max) + 1
}

/**
 * Erzeugt mehrere Zufallszahlen im Bereich von 1 bis max
 * @param {number} count - Anzahl der zu erzeugenden Zufallszahlen
 * @param {number} max - Maximaler Wert (inklusive, default: 6)
 * @returns {Array<number>} - Array mit Zufallszahlen
 */
export function rollDice(count = 1, max = 6) {
  if (typeof count !== 'number' || count < 1) {
    throw new Error('Count must be a positive number')
  }
  if (typeof max !== 'number' || max < 1) {
    throw new Error('Max value must be a positive number')
  }
  
  const results = []
  for (let i = 0; i < count; i++) {
    results.push(rollDie(max))
  }
  return results
}

/**
 * Erzeugt Zufallszahlen und gibt sowohl die einzelnen Ergebnisse als auch die Summe zurück
 * @param {number} count - Anzahl der Würfel
 * @param {number} max - Maximaler Wert pro Würfel (default: 6)
 * @returns {Object} - Objekt mit rolls (Array) und sum (Summe)
 */
export function rollDiceWithSum(count = 1, max = 6) {
  const rolls = rollDice(count, max)
  const sum = rolls.reduce((total, roll) => total + roll, 0)
  
  return {
    rolls,
    sum,
    count,
    max
  }
}

/**
 * Simuliert RPG-Würfelnotation (z.B. "3d6", "1d20", "2d10+5", "max(2d20)")
 * @param {string} notation - Würfelnotation (z.B. "3d6", "1d20+5", "2d8-1", "max(2d20)", "min(3d6)")
 * @returns {Object} - Objekt mit rolls, sum, modifier und der ursprünglichen notation
 */
export function rollNotation(notation = '1d6') {
  if (typeof notation !== 'string') {
    throw new Error('Notation must be a string')
  }
  
  // Entferne Leerzeichen
  const cleanNotation = notation.replace(/\s+/g, '')
  
  // Check für max/min Funktionen (z.B. "max(2d20)", "min(3d6)")
  const functionMatch = cleanNotation.match(/^(max|min)\((\d+)d(\d+)([+-]\d+)?\)$/i)
  if (functionMatch) {
    const func = functionMatch[1].toLowerCase()
    const count = parseInt(functionMatch[2])
    const sides = parseInt(functionMatch[3])
    const modifier = functionMatch[4] ? parseInt(functionMatch[4]) : 0
    
    const rolls = rollDice(count, sides)
    let selectedValue
    
    if (func === 'max') {
      selectedValue = Math.max(...rolls)
    } else if (func === 'min') {
      selectedValue = Math.min(...rolls)
    }
    
    const finalSum = selectedValue + modifier
    
    return {
      notation,
      rolls,
      selectedValue,
      selectedFunction: func,
      baseSum: selectedValue,
      modifier,
      sum: finalSum,
      count,
      sides
    }
  }
  
  // Standard-Notation (z.B. "3d6+2" oder "1d20-1")
  const standardMatch = cleanNotation.match(/^(\d+)d(\d+)([+-]\d+)?$/i)
  if (!standardMatch) {
    throw new Error('Invalid dice notation. Use format like "3d6", "1d20+5", "2d8-2", "max(2d20)", or "min(3d6)"')
  }
  
  const count = parseInt(standardMatch[1])
  const sides = parseInt(standardMatch[2])
  const modifier = standardMatch[3] ? parseInt(standardMatch[3]) : 0
  
  const rolls = rollDice(count, sides)
  const baseSum = rolls.reduce((total, roll) => total + roll, 0)
  const finalSum = baseSum + modifier
  
  return {
    notation,
    rolls,
    baseSum,
    modifier,
    sum: finalSum,
    count,
    sides
  }
}

/**
 * Erzeugt eine Zufallszahl in einem bestimmten Bereich (min bis max, inklusive)
 * @param {number} min - Minimaler Wert (inklusive)
 * @param {number} max - Maximaler Wert (inklusive)
 * @returns {number} - Zufallszahl zwischen min und max
 */
export function randomBetween(min = 1, max = 6) {
  if (typeof min !== 'number' || typeof max !== 'number') {
    throw new Error('Min and max must be numbers')
  }
  if (min > max) {
    throw new Error('Min value cannot be greater than max value')
  }
  
  return Math.floor(Math.random() * (max - min + 1)) + min
}

/**
 * Wählt ein zufälliges Element aus einem Array
 * @param {Array} array - Array mit Elementen zur Auswahl
 * @returns {*} - Zufällig gewähltes Element
 */
export function randomChoice(array) {
  if (!Array.isArray(array) || array.length === 0) {
    throw new Error('Array must be a non-empty array')
  }
  
  const randomIndex = Math.floor(Math.random() * array.length)
  return array[randomIndex]
}

/**
 * Mischt ein Array zufällig (Fisher-Yates Shuffle)
 * @param {Array} array - Array zum Mischen
 * @returns {Array} - Neues gemischtes Array (Original bleibt unverändert)
 */
export function shuffleArray(array) {
  if (!Array.isArray(array)) {
    throw new Error('Input must be an array')
  }
  
  const shuffled = [...array]
  for (let i = shuffled.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]]
  }
  return shuffled
}

