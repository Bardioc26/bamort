---
description: 'Instructions for writing Vue 3 components following project conventions and best practices'
applyTo: '**/*.vue, **/*.ts, **/*.js, **/*.scss'
---

# Vue 3 Development Instructions

Follow Vue 3 best practices and project-specific conventions when writing components.

## Component Structure

### Standard Component Layout
```vue
<template>
  <!-- HTML template -->
</template>

<style scoped>
/* Component-specific styles */
</style>

<script>
// Component logic
</script>
```

**Order matters**: Template, Style, Script (as seen throughout the codebase)

### Options API Pattern (Primary)
Use Options API for consistency with existing codebase:

```vue
<script>
export default {
  name: "ComponentName",
  props: ["id"],
  data() {
    return {
      items: [],
      isLoading: false
    }
  },
  async created() {
    // Initialization logic
  },
  methods: {
    async methodName() {
      // Method implementation
    }
  }
}
</script>
```

## API Communication

### Using the API Utility
Always use `API` from `utils/api.js` - it handles authentication automatically:

```js
import API from '../utils/api'

// In methods:
const response = await API.get(`/api/characters/${this.id}`)
const data = await API.post('/api/characters', character)
```

**Never** manually add Authorization headers - the interceptor handles this.

### Error Handling Pattern
```js
try {
  const response = await API.get(`/api/endpoint`)
  this.data = response.data
} catch (error) {
  console.error('Failed to load data:', error)
  alert(this.$t('errors.loadFailed') + ': ' + (error.response?.data?.error || error.message))
}
```

## Internationalization (i18n)

### Using Translations
```vue
<template>
  <h2>{{ $t('char') }}: {{ character.name }}</h2>
  <button>{{ $t('export.exportPDF') }}</button>
</template>
```

### Adding New Translations
**ALWAYS** add to both `src/locales/de` and `src/locales/en`:

```js
// src/locales/de
export default {
  export: {
    selectTemplate: 'Vorlage wählen',
    exportPDF: 'PDF Export'
  }
}

// src/locales/en
export default {
  export: {
    selectTemplate: 'Select Template',
    exportPDF: 'Export PDF'
  }
}
```

**Note**: Locale files use `.js` extension and export objects, not JSON.

## Modal Dialog Pattern

### Standard Modal Structure
```vue
<template>
  <!-- Trigger -->
  <button @click="showDialog = true">Open</button>

  <!-- Modal -->
  <div v-if="showDialog" class="modal-overlay" @click.self="showDialog = false">
    <div class="modal-content">
      <div class="modal-header">
        <h3>{{ $t('modal.title') }}</h3>
        <button @click="showDialog = false" class="close-button">&times;</button>
      </div>
      <div class="modal-body">
        <!-- Content -->
      </div>
      <div class="modal-footer">
        <button @click="showDialog = false" class="btn-cancel">{{ $t('cancel') }}</button>
        <button @click="handleSubmit" class="btn-primary">{{ $t('submit') }}</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      showDialog: false
    }
  }
}
</script>
```

**Key conventions:**
- Use `@click.self` on overlay to close on outside click
- Include close button (×) in header
- Separate header, body, footer sections

## Component Communication

### Props (Parent → Child)
```vue
<script>
export default {
  props: {
    character: Object,
    id: [String, Number]
  }
}
</script>
```

### Events (Child → Parent)
```vue
<template>
  <button @click="notifyParent">Update</button>
</template>

<script>
export default {
  methods: {
    notifyParent() {
      this.$emit('character-updated', this.character)
    }
  }
}
</script>

<!-- Parent component -->
<template>
  <ChildComponent @character-updated="refreshCharacter" />
</template>
```

## State Management

### Loading States
Always show feedback for async operations:

```vue
<template>
  <button @click="submit" :disabled="isLoading">
    <span v-if="!isLoading">{{ $t('submit') }}</span>
    <span v-else>{{ $t('loading') }}</span>
  </button>
</template>

<script>
export default {
  data() {
    return { isLoading: false }
  },
  methods: {
    async submit() {
      this.isLoading = true
      try {
        await API.post('/api/endpoint', this.data)
      } finally {
        this.isLoading = false
      }
    }
  }
}
</script>
```

### Disabling Form Elements During Loading
```vue
<select v-model="selected" :disabled="isLoading">
<input type="checkbox" v-model="option" :disabled="isLoading">
```

## Browser Compatibility

### Popup Blocker Workaround
Open windows **synchronously** before async operations:

```js
async exportToPDF() {
  // Open window FIRST (synchronously)
  const pdfWindow = window.open('', '_blank')
  if (!pdfWindow) {
    alert(this.$t('export.popupBlocked'))
    return
  }
  
  // Show loading page
  pdfWindow.document.write('<html>...</html>')
  
  // Then do async work
  const response = await API.get('/api/pdf/export')
  
  // Update window with result
  pdfWindow.location.href = url
}
```

**Critical**: `window.open()` must be called synchronously in the click handler, not after `await`.

## Common Patterns

### Dynamic Component Loading
```vue
<template>
  <component :is="currentView" :character="character" @character-updated="refresh"/>
</template>

<script>
import ViewA from './ViewA.vue'
import ViewB from './ViewB.vue'

export default {
  components: { ViewA, ViewB },
  data() {
    return { currentView: 'ViewA' }
  }
}
</script>
```

### Conditional Rendering
- Use `v-if` for elements that toggle rarely
- Use `v-show` for frequent toggles
- Use `v-for` with `:key` attribute always

### Event Modifiers
- `@click.self` - only trigger if clicked element itself
- `@submit.prevent` - prevent form submission
- `@keyup.enter` - keyboard event handling

## Best Practices

1. **use global CSS definition to ensure consistent style 
2. **Always use scoped styles** to avoid CSS conflicts
3. **Name components with PascalCase** (e.g., `CharacterDetails.vue`)
4. **Use meaningful prop names** and validate when possible
5. **Handle errors gracefully** with user-friendly messages
6. **Keep template logic simple** - move complex logic to methods
7. **Clean up resources** in `beforeUnmount` if needed
8. **Test with actual API** running in Docker container
9. **Check HMR reload** - view logs: `docker logs bamort-frontend-dev`

## Anti-Patterns to Avoid

❌ Don't use `v-if` and `v-for` on the same element
❌ Don't mutate props directly
❌ Don't forget to handle loading and error states
❌ Don't use inline styles - use scoped CSS
❌ Don't call API methods in template expressions
❌ Don't forget translations - add to both DE and EN
