package pdfrender

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

// TemplateLoader manages loading and rendering of HTML templates
type TemplateLoader struct {
	templateDir string
	templates   *template.Template
	metadata    map[string][]BlockMetadata
}

// NewTemplateLoader creates a new template loader for the given directory
func NewTemplateLoader(templateDir string) *TemplateLoader {
	return &TemplateLoader{
		templateDir: templateDir,
		metadata:    make(map[string][]BlockMetadata),
	}
}

// LoadTemplates loads all .html templates from the template directory
func (tl *TemplateLoader) LoadTemplates() error {
	// Check if directory exists
	if _, err := os.Stat(tl.templateDir); os.IsNotExist(err) {
		return fmt.Errorf("template directory does not exist: %s", tl.templateDir)
	}

	// Find all .html files
	pattern := filepath.Join(tl.templateDir, "*.html")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to find template files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no template files found in %s", tl.templateDir)
	}

	// Create template with custom functions
	tmpl := template.New("").Funcs(getTemplateFuncs())

	// Parse all templates
	tmpl, err = tmpl.ParseFiles(files...)
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	tl.templates = tmpl

	// Extract metadata from each template
	for _, file := range files {
		basename := filepath.Base(file)
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", basename, err)
		}

		metadata := ParseTemplateMetadata(string(content))
		tl.metadata[basename] = metadata
	}

	return nil
}

// RenderTemplate renders a specific template with the given data
// For continuation templates (e.g., "page1.2_stats.html"), it automatically
// falls back to the base template (e.g., "page1_stats.html") if the continuation
// template doesn't exist as a physical file
func (tl *TemplateLoader) RenderTemplate(templateName string, data interface{}) (string, error) {
	if tl.templates == nil {
		return "", fmt.Errorf("templates not loaded, call LoadTemplates first")
	}

	tmpl := tl.templates.Lookup(templateName)
	if tmpl == nil {
		// Try to extract base template name for continuation pages
		// e.g., "page1.2_stats.html" -> "page1_stats.html"
		baseTemplateName := ExtractBaseTemplateName(templateName)
		if baseTemplateName != templateName {
			tmpl = tl.templates.Lookup(baseTemplateName)
			if tmpl == nil {
				return "", fmt.Errorf("template not found: %s (and base template %s not found)", templateName, baseTemplateName)
			}
		} else {
			return "", fmt.Errorf("template not found: %s", templateName)
		}
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

// GetTemplateMetadata returns the metadata blocks for a given template
func (tl *TemplateLoader) GetTemplateMetadata(templateName string) []BlockMetadata {
	return tl.metadata[templateName]
}

// getTemplateFuncs returns custom template functions
func getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"iterate": func(n int) []int {
			result := make([]int, n)
			for i := 0; i < n; i++ {
				result[i] = i
			}
			return result
		},
	}
}
