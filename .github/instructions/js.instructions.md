---
description: 'Instructions for writing JavaScript following project conventions and ES6+ best practices'
applyTo: '**/*.js,**/*.mjs'
---

# JavaScript Development Instructions

Follow ES6+ best practices and project-specific patterns for JavaScript code.

## Module System

### ES6 Modules
Use ES6 import/export syntax:

```js
// Named exports
export const API = axios.create({ ... })
export function helper() { ... }

// Default export
export default {
  messages: { ... }
}

// Imports
import API from '../utils/api'
import { createI18n } from 'vue-i18n'
```

## API Configuration (`utils/api.js`)

### Standard Axios Instance Pattern
```js
import axios from 'axios'

const API = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8180'
})

// Request interceptor - adds auth token
API.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor - handles 401
API.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      // Optional: redirect to login
    }
    return Promise.reject(error)
  }
)

export default API
```

**Key points:**
- Single Axios instance for the entire app
- Auto-adds Authorization header from localStorage
- Auto-handles 401 responses
- Uses Vite environment variables

## Pinia Store Pattern (`stores/`)

### Standard Store Structure
```js
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
```

## Locale Files (`locales/de`, `locales/en`)

### Translation Object Structure
**Important**: Locale files use `.js` extension and export objects (not `.json`):

```js
// locales/de
export default {
  char: 'Figur',
  menu: {
    Datasheet: 'Datenblatt',
    Skill: 'Fertigkeiten'
  },
  export: {
    selectTemplate: 'Vorlage wählen',
    exportPDF: 'PDF Export',
    exporting: 'Exportiere...',
    pleaseSelectTemplate: 'Bitte Vorlage auswählen'
  }
}
```

**Conventions:**
- Nested objects for grouping related translations
- camelCase for keys
- Always add to both `de` and `en` files simultaneously
- Keep structure identical between languages

## Async/Await Patterns

### Error Handling
```js
try {
  const response = await API.get('/api/endpoint')
  return response.data
} catch (error) {
  console.error('Operation failed:', error)
  throw error // or handle gracefully
}
```

### Multiple Parallel Requests
```js
const [characters, templates] = await Promise.all([
  API.get('/api/characters'),
  API.get('/api/pdf/templates')
])
```

## Browser APIs

### LocalStorage Usage
```js
// Save
localStorage.setItem('token', response.data.token)
localStorage.setItem('language', 'de')

// Retrieve
const token = localStorage.getItem('token')
const lang = localStorage.getItem('language') || 'de'

// Remove
localStorage.removeItem('token')
```

### Blob/File Handling
```js
// Create blob from response
const blob = new Blob([response.data], { type: 'application/pdf' })
const url = window.URL.createObjectURL(blob)

// Open in new window
const pdfWindow = window.open(url, '_blank')

// Clean up after use
setTimeout(() => window.URL.revokeObjectURL(url), 10000)
```

### URL Parameters
```js
// Build query string
const params = new URLSearchParams({
  template: templateId,
  showUserName: 'true'
})
const url = `/api/export?${params.toString()}`

// Extract from params
const queryParams = Object.fromEntries(params)
```

## Event Handling

### Debouncing/Throttling
For search inputs or resize events:

```js
let debounceTimer
function debounce(func, delay = 300) {
  return (...args) => {
    clearTimeout(debounceTimer)
    debounceTimer = setTimeout(() => func(...args), delay)
  }
}

// Usage
const search = debounce(async (query) => {
  const results = await API.get(`/api/search?q=${query}`)
}, 300)
```

### Cleanup
```js
// Save timer ID for cleanup
this.timer = setTimeout(() => { ... }, 5000)

// Clean up in component lifecycle
beforeUnmount() {
  if (this.timer) clearTimeout(this.timer)
}
```

## Array/Object Operations

### Array Methods
```js
// Filter
const learned = skills.filter(s => s.Fertigkeitswert > 0)

// Map
const names = characters.map(c => c.name)

// Find
const char = characters.find(c => c.id === 18)

// Some/Every
const hasSkills = character.fertigkeiten.some(f => f.Fertigkeitswert > 0)
```

### Object Destructuring
```js
// Response destructuring
const { data, headers, status } = response

// Props destructuring
const { character, template, showUserName = false } = options
```

### Spread Operator
```js
// Merge objects
const merged = { ...defaults, ...userOptions }

// Copy array
const copy = [...originalArray]
```

## Common Patterns

### Loading State Management
```js
export default {
  data() {
    return {
      isLoading: false,
      data: null,
      error: null
    }
  },
  async created() {
    await this.loadData()
  },
  methods: {
    async loadData() {
      this.isLoading = true
      this.error = null
      try {
        const response = await API.get('/api/data')
        this.data = response.data
      } catch (error) {
        this.error = error.message
      } finally {
        this.isLoading = false
      }
    }
  }
}
```

### Form Validation
```js
validateForm() {
  if (!this.selectedTemplate) {
    alert(this.$t('export.pleaseSelectTemplate'))
    return false
  }
  return true
}

async submit() {
  if (!this.validateForm()) return
  
  // Proceed with submission
}
```

## Best Practices

1. **Use `const` by default**, `let` when reassignment needed, never `var`
2. **Prefer arrow functions** for callbacks and short functions
3. **Use template literals** for string interpolation
4. **Handle promise rejections** with try/catch or .catch()
5. **Clean up timers and intervals** in component lifecycle
6. **Use optional chaining** `?.` for nested properties
7. **Use nullish coalescing** `??` instead of `||` for default values
8. **Keep functions small** and single-purpose
9. **Document complex logic** with comments
10. **Use meaningful variable names** - avoid single letters except loops

## Anti-Patterns to Avoid

❌ Don't use `var` - use `const` or `let`
❌ Don't ignore promise rejections
❌ Don't mutate function parameters
❌ Don't create memory leaks (clean up listeners, timers)
❌ Don't use `eval()` or `new Function()`
❌ Don't mix callbacks and promises
❌ Don't forget to handle edge cases (null, undefined, empty arrays)
❌ Don't use `==` - always use `===` for comparisons

## Environment Variables (Vite)

### Accessing Variables
```js
const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8180'
const isDev = import.meta.env.DEV
const isProd = import.meta.env.PROD
```

**Convention**: Prefix all custom variables with `VITE_`
