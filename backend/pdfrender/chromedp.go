package pdfrender

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// PDFRenderer handles HTML to PDF conversion using chromedp
type PDFRenderer struct {
	// Configuration options can be added here later
}

// NewPDFRenderer creates a new PDF renderer
func NewPDFRenderer() *PDFRenderer {
	return &PDFRenderer{}
}

// RenderHTMLToPDF converts HTML string to PDF bytes using chromedp
func (r *PDFRenderer) RenderHTMLToPDF(html string) ([]byte, error) {
	// Create context with timeout
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set timeout for PDF generation
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var pdfBytes []byte

	// Configure PDF printing options for A4 landscape
	printParams := page.PrintToPDF().
		WithPaperWidth(11.69). // A4 landscape width in inches
		WithPaperHeight(8.27). // A4 landscape height in inches
		WithMarginTop(0).      // No margins (template handles spacing)
		WithMarginBottom(0).
		WithMarginLeft(0).
		WithMarginRight(0).
		WithPrintBackground(true).  // Include background colors/images
		WithPreferCSSPageSize(true) // Use CSS page size if specified

	// Execute chromedp tasks
	err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Set HTML content
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Wait a bit for rendering
			time.Sleep(500 * time.Millisecond)
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Print to PDF
			var err error
			pdfBytes, _, err = printParams.Do(ctx)
			return err
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("chromedp failed to render PDF: %w", err)
	}

	return pdfBytes, nil
}

// ImageToBase64DataURI converts image bytes to a data URI for embedding in HTML
func ImageToBase64DataURI(imageData []byte, mimeType string) string {
	encoded := base64.StdEncoding.EncodeToString(imageData)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)
}
