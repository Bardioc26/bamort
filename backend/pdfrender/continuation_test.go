package pdfrender

import (
	"testing"
)

// TestContinuationPages_WhenSkillsExceedCapacity tests that overflow items
// are distributed to continuation pages when template capacity is exceeded
func TestContinuationPages_WhenSkillsExceedCapacity(t *testing.T) {
	// Arrange - Create more skills than fit on a single page
	// Assuming page1_stats has capacity of 64 (2 columns x 32 each)
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Create 100 skills to force overflow
	skills := make([]SkillViewModel, 100)
	for i := 0; i < 100; i++ {
		skills[i] = SkillViewModel{
			Name:  "Skill " + string(rune('A'+i%26)),
			Value: 10 + i%20,
		}
	}

	// Act - Paginate skills for page1_stats
	pages, err := paginator.PaginateSkills(skills, "page_1.html", "")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should have more than 1 page (100 skills / 64 capacity = 2 pages)
	if len(pages) < 2 {
		t.Errorf("Expected at least 2 pages for 100 skills, got %d", len(pages))
	}

	// First page should be "page_1.html"
	if pages[0].TemplateName != "page_1.html" {
		t.Errorf("Expected first page template 'page_1.html', got '%s'", pages[0].TemplateName)
	}

	// Second page should be continuation page with name pattern "page_1.2.html"
	expectedContinuation := "page_1.2.html"
	if pages[1].TemplateName != expectedContinuation {
		t.Errorf("Expected continuation page template '%s', got '%s'",
			expectedContinuation, pages[1].TemplateName)
	}

	// Verify page numbers are sequential
	for i, page := range pages {
		expectedPageNum := i + 1
		if page.PageNumber != expectedPageNum {
			t.Errorf("Page %d: expected page number %d, got %d",
				i, expectedPageNum, page.PageNumber)
		}
	}

	// Verify all 100 skills are distributed
	totalSkills := 0
	for _, page := range pages {
		for _, data := range page.Data {
			if skillSlice, ok := data.([]SkillViewModel); ok {
				totalSkills += len(skillSlice)
			}
		}
	}

	if totalSkills != 100 {
		t.Errorf("Expected 100 total skills across all pages, got %d", totalSkills)
	}
}

// TestContinuationPages_WhenWeaponsExceedCapacity tests overflow handling for weapons
func TestContinuationPages_WhenWeaponsExceedCapacity(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Get capacity for page2_play weapons (should be 12)
	var weaponsCapacity int
	for _, tmpl := range templateSet.Templates {
		if tmpl.Metadata.Name == "page_2.html" {
			for _, block := range tmpl.Metadata.Blocks {
				if block.ListType == "weapons" {
					weaponsCapacity = block.MaxItems
					t.Logf("Found weapons block: %s with capacity %d", block.Name, block.MaxItems)
				}
			}
			break
		}
	}

	t.Logf("Total weapons capacity per page: %d", weaponsCapacity)

	// Create more weapons than capacity (enough for exactly 2 pages)
	// If capacity is 3, create 4 weapons to get 2 pages (3 + 1)
	numWeapons := weaponsCapacity + 1
	weapons := make([]WeaponViewModel, numWeapons)
	for i := 0; i < len(weapons); i++ {
		weapons[i] = WeaponViewModel{
			Name:  "Weapon " + string(rune('A'+i%26)),
			Value: 10 + i,
		}
	}

	t.Logf("Created %d weapons (capacity %d)", numWeapons, weaponsCapacity)

	// Act
	pages, err := paginator.PaginateWeapons(weapons, "page_2.html")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should have 2 pages (base + continuation)
	if len(pages) != 2 {
		t.Errorf("Expected 2 pages (main + continuation), got %d", len(pages))
	}

	// First page should be original template
	if pages[0].TemplateName != "page_2.html" {
		t.Errorf("Expected first page 'page_2.html', got '%s'", pages[0].TemplateName)
	}

	// Second page should be continuation
	if pages[1].TemplateName != "page_2.2.html" {
		t.Errorf("Expected continuation 'page_2.2.html', got '%s'", pages[1].TemplateName)
	}
}

// TestContinuationPages_MultipleOverflows tests handling of multiple continuation pages
func TestContinuationPages_MultipleOverflows(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Get actual capacity from template
	var skillsCapacity int
	for _, tmpl := range templateSet.Templates {
		if tmpl.Metadata.Name == "page_1.html" {
			for _, block := range tmpl.Metadata.Blocks {
				if block.ListType == "skills" {
					skillsCapacity += block.MaxItems
				}
			}
			break
		}
	}

	t.Logf("Skills capacity per page: %d", skillsCapacity)

	// Create 200 skills to force multiple continuations
	skills := make([]SkillViewModel, 200)
	for i := 0; i < 200; i++ {
		skills[i] = SkillViewModel{
			Name:  "Skill " + string(rune('A'+i%26)),
			Value: 10 + i%20,
		}
	}

	// Calculate expected pages
	expectedPages := (200 + skillsCapacity - 1) / skillsCapacity

	// Act
	pages, err := paginator.PaginateSkills(skills, "page_1.html", "")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should have calculated number of pages
	if len(pages) != expectedPages {
		t.Errorf("Expected %d pages for 200 skills (capacity %d), got %d", expectedPages, skillsCapacity, len(pages))
	}

	t.Logf("Created %d pages for 200 skills", len(pages))

	// Verify template names follow pattern: page_1.html, then all use page_1.2.html
	for i, page := range pages {
		var expectedTemplate string
		if i == 0 {
			expectedTemplate = "page_1.html"
		} else {
			// All continuation pages use the same .2 template
			expectedTemplate = "page_1.2.html"
		}

		if page.TemplateName != expectedTemplate {
			t.Errorf("Page %d: expected template '%s', got '%s'",
				i+1, expectedTemplate, page.TemplateName)
		}
	}
}

// TestContinuationPages_NoOverflow tests that no continuation pages are created
// when items fit within the original page capacity
func TestContinuationPages_NoOverflow(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Create only 10 skills (well below capacity)
	skills := make([]SkillViewModel, 10)
	for i := 0; i < 10; i++ {
		skills[i] = SkillViewModel{
			Name:  "Skill " + string(rune('A'+i)),
			Value: 10 + i,
		}
	}

	// Act
	pages, err := paginator.PaginateSkills(skills, "page_1.html", "")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should have exactly 1 page (no overflow)
	if len(pages) != 1 {
		t.Errorf("Expected 1 page for 10 skills, got %d", len(pages))
	}

	// Should use original template, not continuation
	if pages[0].TemplateName != "page_1.html" {
		t.Errorf("Expected original template 'page_1.html', got '%s'", pages[0].TemplateName)
	}
}
