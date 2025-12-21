package pdfrender

import "testing"

func TestGenerateContinuationTemplateName(t *testing.T) {
	tests := []struct {
		name     string
		template string
		pageNum  int
		expected string
	}{
		{
			name:     "First page returns original",
			template: "page_1.html",
			pageNum:  1,
			expected: "page_1.html",
		},
		{
			name:     "Second page gets continuation name",
			template: "page_1.html",
			pageNum:  2,
			expected: "page_1.2.html",
		},
		{
			name:     "Third page uses same continuation template (.2)",
			template: "page_2.html",
			pageNum:  3,
			expected: "page_2.2.html",
		},
		{
			name:     "All continuation pages use .2 template",
			template: "page_1.html",
			pageNum:  10,
			expected: "page_1.2.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateContinuationTemplateName(tt.template, tt.pageNum)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestExtractBaseTemplateName(t *testing.T) {
	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "Base template returns itself",
			template: "page_1.html",
			expected: "page_1.html",
		},
		{
			name:     "Continuation template returns base",
			template: "page_1.2.html",
			expected: "page_1.html",
		},
		{
			name:     "Multiple digit continuation",
			template: "page_1.10.html",
			expected: "page_1.html",
		},
		{
			name:     "Different page type",
			template: "page_2.3.html",
			expected: "page_2.html",
		},
		{
			name:     "Template without underscore",
			template: "simple.html",
			expected: "simple.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractBaseTemplateName(tt.template)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestContinuationRoundTrip(t *testing.T) {
	// Test that generating a continuation name and extracting the base works correctly
	original := "page_1.html"

	for pageNum := 1; pageNum <= 5; pageNum++ {
		continuation := GenerateContinuationTemplateName(original, pageNum)
		extracted := ExtractBaseTemplateName(continuation)

		if extracted != original {
			t.Errorf("Round trip failed for page %d: original '%s', continuation '%s', extracted '%s'",
				pageNum, original, continuation, extracted)
		}
	}
}
