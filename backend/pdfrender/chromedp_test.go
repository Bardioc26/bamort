package pdfrender

import (
	"strings"
	"testing"
)

func TestRenderHTMLToPDF_BasicHTML(t *testing.T) {
	// Arrange
	html := `<!DOCTYPE html>
<html>
<head><title>Test PDF</title></head>
<body><h1>Hello PDF</h1></body>
</html>`

	renderer := NewPDFRenderer()

	// Act
	pdfBytes, err := renderer.RenderHTMLToPDF(html)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Error("Expected non-empty PDF bytes")
	}

	// Check PDF magic number
	if len(pdfBytes) < 4 || string(pdfBytes[0:4]) != "%PDF" {
		t.Error("Output does not appear to be a PDF (missing %PDF header)")
	}
}

func TestRenderHTMLToPDF_WithStyles(t *testing.T) {
	// Arrange
	html := `<!DOCTYPE html>
<html>
<head>
<style>
body { font-family: Arial; margin: 20px; }
h1 { color: blue; }
</style>
</head>
<body>
<h1>Styled Content</h1>
<p>This is a test paragraph with styling.</p>
</body>
</html>`

	renderer := NewPDFRenderer()

	// Act
	pdfBytes, err := renderer.RenderHTMLToPDF(html)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Error("Expected non-empty PDF bytes")
	}

	// Check PDF magic number
	if string(pdfBytes[0:4]) != "%PDF" {
		t.Error("Output does not appear to be a PDF")
	}
}

func TestRenderHTMLToPDF_EmptyHTML(t *testing.T) {
	// Arrange
	html := ""
	renderer := NewPDFRenderer()

	// Act
	pdfBytes, err := renderer.RenderHTMLToPDF(html)

	// Assert - should still generate a PDF even if empty
	if err != nil {
		t.Fatalf("Expected no error for empty HTML, got %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Error("Expected non-empty PDF bytes even for empty HTML")
	}
}

func TestImageToBase64DataURI_PNG(t *testing.T) {
	// Arrange - simple 1x1 red PNG
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
	}

	// Act
	dataURI := ImageToBase64DataURI(pngData, "image/png")

	// Assert
	if !strings.HasPrefix(dataURI, "data:image/png;base64,") {
		t.Errorf("Expected data URI to start with 'data:image/png;base64,', got %s", dataURI[:30])
	}

	if len(dataURI) < 30 {
		t.Error("Expected base64 encoded data URI to be longer")
	}
}

func TestImageToBase64DataURI_JPEG(t *testing.T) {
	// Arrange
	jpegData := []byte{0xFF, 0xD8, 0xFF} // JPEG magic number

	// Act
	dataURI := ImageToBase64DataURI(jpegData, "image/jpeg")

	// Assert
	if !strings.HasPrefix(dataURI, "data:image/jpeg;base64,") {
		t.Error("Expected data URI to start with 'data:image/jpeg;base64,'")
	}
}
