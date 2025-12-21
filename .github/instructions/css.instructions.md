---
description: 'Instructions for writing CSS following project conventions and best practices'
applyTo: '**/*.css,**/*.vue'
---

# CSS Development Instructions

Follow project-specific CSS conventions and modern best practices.

## Scoped Styles in Vue Components

### Always Use Scoped Styles
```vue
<style scoped>
.component-class {
  /* Styles only apply to this component */
}
</style>
```

**Critical**: Use `scoped` attribute to prevent style conflicts between components.

## Layout Patterns

### Flexbox for Component Layouts
Standard pattern for headers, modals, and lists:

```css
.header-content {
  display: flex;
  align-items: center;
  gap: 15px; /* Use gap instead of margin */
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
```

### Common Flex Patterns
```css
/* Horizontal layout with spacing */
.horizontal-layout {
  display: flex;
  gap: 10px;
  align-items: center;
}

/* Vertical centering */
.centered {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
}

/* Space between items */
.space-between {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
```

## Modal Dialog Styling

### Standard Modal Pattern
```css
/* Overlay - covers entire viewport */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

/* Modal content container */
.modal-content {
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

/* Modal sections */
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #dee2e6;
}

.modal-body {
  padding: 20px;
  position: relative; /* For loading overlays */
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 20px;
  border-top: 1px solid #dee2e6;
}
```

## Button Styling

### Standard Button Styles
```css
/* Primary action button */
.btn-primary,
.btn-export {
  padding: 10px 20px;
  border: 1px solid #007bff;
  border-radius: 6px;
  background: #007bff;
  color: white;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s ease;
}

.btn-primary:hover:not(:disabled) {
  background: #0056b3;
  border-color: #0056b3;
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Cancel/secondary button */
.btn-cancel {
  padding: 10px 20px;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  background: #f8f9fa;
  color: #495057;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s ease;
}

.btn-cancel:hover {
  background: #e9ecef;
  border-color: #adb5bd;
}

/* Icon-only button */
.export-button-small {
  width: 40px;
  height: 40px;
  padding: 0;
  border: 1px solid #007bff;
  border-radius: 8px;
  background: #007bff;
  color: white;
  font-size: 1.2rem;
  cursor: pointer;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.export-button-small:hover {
  background: #0056b3;
  transform: scale(1.05);
}
```

## Form Elements

### Input Styling
```css
.template-select,
input[type="text"],
input[type="email"] {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  background: white;
  color: #495057;
  font-size: 0.95rem;
}

.template-select:focus,
input:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}

input:disabled,
select:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  background: #e9ecef;
}
```

### Checkbox Styling
```css
.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  user-select: none;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}
```

## Loading Animations

### Spinner Animation
```css
.spinner {
  border: 4px solid #f3f3f3;
  border-top: 4px solid #007bff;
  border-radius: 50%;
  width: 50px;
  height: 50px;
  animation: spin 1s linear infinite;
  margin-bottom: 15px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
```

### Loading Overlay
```css
.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.95);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  z-index: 10;
}

.loading-overlay p {
  color: #007bff;
  font-weight: 500;
  margin: 0;
}
```

## Color Scheme

### Standard Colors
```css
/* Primary */
--primary: #007bff;
--primary-hover: #0056b3;

/* Text */
--text-primary: #333;
--text-secondary: #495057;
--text-muted: #666;

/* Backgrounds */
--bg-light: #f8f9fa;
--bg-gray: #e9ecef;

/* Borders */
--border-light: #dee2e6;
--border-dark: #adb5bd;

/* Semantic colors */
--success: #28a745;
--danger: #dc3545;
--warning: #ffc107;
```

Use these consistently across components for visual coherence.

## Spacing System

### Use Consistent Spacing
```css
/* Prefer these spacing values */
gap: 8px;   /* Tight spacing */
gap: 10px;  /* Default spacing */
gap: 15px;  /* Medium spacing */
gap: 20px;  /* Large spacing */

padding: 10px 12px; /* Inputs */
padding: 20px;      /* Modal sections */
```

## Typography

### Font Sizing
```css
h2 {
  font-size: 1.5rem;
  margin: 0;
  color: #333;
}

h3 {
  font-size: 1.25rem;
  margin: 0;
  color: #333;
}

p, span {
  font-size: 0.95rem;
  color: #495057;
}

small {
  font-size: 0.875rem;
  color: #666;
}
```

### Font Weight
```css
font-weight: 400; /* Normal */
font-weight: 500; /* Medium (buttons, labels) */
font-weight: 600; /* Semibold (headings) */
```

## Transitions and Animations

### Standard Transitions
```css
/* Buttons, interactive elements */
transition: all 0.2s ease;

/* Background changes */
transition: background 0.3s ease;

/* Transform animations */
transition: transform 0.2s ease;
```

### Hover Effects
```css
button:hover {
  transform: scale(1.02); /* Subtle scale */
}

.card:hover {
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15);
}
```

## Responsive Design

### Mobile-First Approach
```css
/* Base styles (mobile) */
.modal-content {
  width: 90%;
  max-width: 500px;
}

/* Tablet and up */
@media (min-width: 768px) {
  .modal-content {
    width: 600px;
  }
}

/* Desktop */
@media (min-width: 1024px) {
  .modal-content {
    width: 700px;
  }
}
```

## Best Practices

1. **Always use `scoped`** on Vue component styles
2. **Use flexbox for layouts** instead of floats or positioning
3. **Use `gap` property** instead of margins for spacing
4. **Keep selectors simple** - avoid deep nesting
5. **Use relative units** (`rem`, `em`) for font sizes
6. **Add transitions** for interactive elements
7. **Use CSS variables** for repeated values
8. **Keep z-index organized** (modals: 1000, dropdowns: 100, etc.)
9. **Test hover states** for all interactive elements
10. **Include disabled states** for buttons and inputs

## Anti-Patterns to Avoid

❌ Don't use `!important` unless absolutely necessary
❌ Don't use inline styles in templates - use classes
❌ Don't use fixed pixel widths for responsive layouts
❌ Don't nest selectors more than 3 levels deep
❌ Don't use IDs for styling - use classes
❌ Don't forget `:hover`, `:focus`, `:disabled` states
❌ Don't use `position: absolute` unless necessary
❌ Don't forget to test in different viewport sizes
❌ Don't use vendor prefixes manually - use autoprefixer

## Common Component Patterns

### Close Button
```css
.close-button {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #999;
  cursor: pointer;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-button:hover {
  color: #333;
}
```

### Full-Height Container
```css
.character-details {
  width: 100%;
  height: 100%;
  padding: 20px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}
```

### Submenu/Tabs
```css
.submenu {
  display: flex;
  gap: 10px;
  margin: 20px 0;
  flex-wrap: wrap;
}

.submenu button {
  padding: 10px 16px;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  background: #f8f9fa;
  color: #495057;
  cursor: pointer;
  transition: all 0.2s ease;
}

.submenu button.active {
  background: #007bff;
  color: white;
  border-color: #007bff;
}
```
