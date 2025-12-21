package pdfrender

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// InlineResources inlines external CSS and images into the HTML for PDF rendering
func InlineResources(html string, templateDir string) (string, error) {
	// Inline CSS
	cssRegex := regexp.MustCompile(`<link\s+rel="stylesheet"\s+href="([^"]+)"[^>]*>`)
	html = cssRegex.ReplaceAllStringFunc(html, func(match string) string {
		matches := cssRegex.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}

		cssPath := matches[1]
		fullPath := filepath.Join(templateDir, cssPath)

		cssContent, err := os.ReadFile(fullPath)
		if err != nil {
			// Return original tag if CSS can't be loaded
			return match
		}

		return fmt.Sprintf("<style>\n%s\n</style>", string(cssContent))
	})

	// Inline images (except data URIs which are already inline)
	imgRegex := regexp.MustCompile(`<img\s+src="([^"]+)"([^>]*)>`)
	html = imgRegex.ReplaceAllStringFunc(html, func(match string) string {
		matches := imgRegex.FindStringSubmatch(match)
		if len(matches) < 3 {
			return match
		}

		imgSrc := matches[1]
		imgAttrs := matches[2]

		// Skip if already a data URI or template variable
		if strings.HasPrefix(imgSrc, "data:") || strings.Contains(imgSrc, "{{") {
			return match
		}

		fullPath := filepath.Join(templateDir, imgSrc)

		imgData, err := os.ReadFile(fullPath)
		if err != nil {
			// Return original tag if image can't be loaded
			return match
		}

		// Detect mime type from extension
		mimeType := "image/png"
		ext := strings.ToLower(filepath.Ext(imgSrc))
		switch ext {
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".gif":
			mimeType = "image/gif"
		case ".svg":
			mimeType = "image/svg+xml"
		}

		dataURI := ImageToBase64DataURI(imgData, mimeType)
		return fmt.Sprintf(`<img src="%s"%s>`, dataURI, imgAttrs)
	})

	return html, nil
}

// LoadAndInlineTemplate loads a template, renders it, and inlines all resources
func (tl *TemplateLoader) RenderTemplateWithInlinedResources(templateName string, data interface{}) (string, error) {
	// First render the template
	html, err := tl.RenderTemplate(templateName, data)
	if err != nil {
		return "", err
	}

	// Then inline all external resources
	return InlineResources(html, tl.templateDir)
}
