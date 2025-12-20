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
			template: "page1_stats.html",
			pageNum:  1,
			expected: "page1_stats.html",
		},
		{
			name:     "Second page gets continuation name",
			template: "page1_stats.html",
			pageNum:  2,
			expected: "page1.2_stats.html",
		},
		{
			name:     "Third page uses same continuation template (.2)",
			template: "page2_play.html",
			pageNum:  3,
			expected: "page2.2_play.html",
		},
		{
			name:     "All continuation pages use .2 template",
			template: "page1_stats.html",
			pageNum:  10,
			expected: "page1.2_stats.html",
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
			template: "page1_stats.html",
			expected: "page1_stats.html",
		},
		{
			name:     "Continuation template returns base",
			template: "page1.2_stats.html",
			expected: "page1_stats.html",
		},
		{
			name:     "Multiple digit continuation",
			template: "page1.10_stats.html",
			expected: "page1_stats.html",
		},
		{
			name:     "Different page type",
			template: "page2.3_play.html",
			expected: "page2_play.html",
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
	original := "page1_stats.html"

	for pageNum := 1; pageNum <= 5; pageNum++ {
		continuation := GenerateContinuationTemplateName(original, pageNum)
		extracted := ExtractBaseTemplateName(continuation)

		if extracted != original {
			t.Errorf("Round trip failed for page %d: original '%s', continuation '%s', extracted '%s'",
				pageNum, original, continuation, extracted)
		}
	}
}
