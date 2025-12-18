# Pagination Implementation Guide

## Overview

The pagination system intelligently splits character data across multiple columns and pages based on template capacity metadata. It ensures proper distribution of skills, spells, weapons, and equipment without exceeding template limits.

## Core Components

### 1. Paginator (`pagination.go`)

The `Paginator` handles all pagination logic:

```go
paginator := NewPaginator(templateSet)
```

#### Main Methods

- **PaginateSkills**: Splits skills across columns and pages
- **PaginateSpells**: Handles spell list pagination with column support
- **PaginateWeapons**: Distributes weapons across multiple pages
- **PaginateEquipment**: Manages equipment pagination
- **CalculatePagesNeeded**: Pre-calculates page count for planning

### 2. PageDistribution

Represents how data is distributed for a single page:

```go
type PageDistribution struct {
    TemplateName string                 // Template to use
    PageNumber   int                    // 1-indexed page number
    Data         map[string]interface{} // Block name -> data slice
}
```

## Usage Examples

### Example 1: Paginate Skills Across Multiple Pages

```go
// Create 100 skills (exceeds 64 capacity of page1_stats)
skills := make([]SkillViewModel, 100)
for i := 0; i < 100; i++ {
    skills[i] = SkillViewModel{
        Name:  "Skill " + strconv.Itoa(i),
        Value: 10 + i%20,
    }
}

// Initialize paginator
templateSet := DefaultA4QuerTemplateSet()
paginator := NewPaginator(templateSet)

// Paginate skills
pages, err := paginator.PaginateSkills(skills, "page1_stats.html", "")
if err != nil {
    log.Fatal(err)
}

// Result: 2 pages
// Page 1: skills_column1 (32) + skills_column2 (32) = 64 skills
// Page 2: skills_column1 (32) + skills_column2 (4) = 36 skills

// Load templates
loader := NewTemplateLoader("templates/Default_A4_Quer")
loader.LoadTemplates()

// Render each page
renderer := NewPDFRenderer()
for _, page := range pages {
    pageData := &PageData{
        Character:     viewModel.Character,
        Attributes:    viewModel.Attributes,
        DerivedValues: viewModel.DerivedValues,
        Meta: PageMeta{
            Date:       time.Now().Format("02.01.2006"),
            PageNumber: page.PageNumber,
        },
    }
    
    // Add skills from this page's distribution
    col1 := page.Data["skills_column1"].([]SkillViewModel)
    col2 := page.Data["skills_column2"].([]SkillViewModel)
    pageData.Skills = append(col1, col2...)
    
    // Render to HTML
    html, _ := loader.RenderTemplate(page.TemplateName, pageData)
    
    // Generate PDF
    pdfBytes, _ := renderer.RenderHTMLToPDF(html)
    
    // Save or return PDF
    os.WriteFile(fmt.Sprintf("character_page%d.pdf", page.PageNumber), pdfBytes, 0644)
}
```

### Example 2: Paginate Spells with Two Columns

```go
// Create 30 spells (exceeds 24 capacity of page3_spell)
spells := make([]SpellViewModel, 30)
for i := 0; i < 30; i++ {
    spells[i] = SpellViewModel{
        Name:     "Zauber " + strconv.Itoa(i),
        AP:       5,
        Duration: "1 Minute",
    }
}

// Paginate
pages, err := paginator.PaginateSpells(spells, "page3_spell.html")

// Result: 2 pages
// Page 1: spells_column1 (12) + spells_column2 (12) = 24 spells
// Page 2: spells_column1 (6) + spells_column2 (0) = 6 spells

for _, page := range pages {
    col1Spells := page.Data["spells_column1"].([]SpellViewModel)
    col2Spells := page.Data["spells_column2"].([]SpellViewModel)
    
    pageData := &PageData{
        Character: viewModel.Character,
        Spells:    append(col1Spells, col2Spells...),
        Meta: PageMeta{
            PageNumber: page.PageNumber,
        },
    }
    
    // Render and generate PDF...
}
```

### Example 3: Pre-Calculate Page Count

```go
// Check how many pages will be needed before pagination
pagesNeeded, err := paginator.CalculatePagesNeeded(
    "page1_stats.html", 
    "skills", 
    len(skills),
)

fmt.Printf("Will need %d pages for %d skills\n", pagesNeeded, len(skills))
```

### Example 4: Paginate Weapons

```go
// Create 50 weapons (exceeds 30 capacity of page2_play)
weapons := make([]WeaponViewModel, 50)
for i := 0; i < 50; i++ {
    weapons[i] = WeaponViewModel{
        Name:   "Waffe " + strconv.Itoa(i),
        Value:  10 + i,
        Damage: "1W6+2",
    }
}

// Paginate weapons
pages, err := paginator.PaginateWeapons(weapons, "page2_play.html")

// Result: 2 pages
// Page 1: weapons_main (30)
// Page 2: weapons_main (20)

for _, page := range pages {
    weaponsData := page.Data["weapons_main"].([]WeaponViewModel)
    
    pageData := &PageData{
        Character: viewModel.Character,
        Weapons:   weaponsData,
        Meta: PageMeta{
            PageNumber: page.PageNumber,
        },
    }
    
    // Render and generate PDF...
}
```

## Template Capacities

### page1_stats.html (Statistics Page)
- **skills_column1**: MAX 32 (left column)
- **skills_column2**: MAX 32 (right column)
- **Total**: 64 skills per page

### page2_play.html (Adventure Page)
- **skills_learned**: MAX 24 (FILTER: learned)
- **skills_unlearned**: MAX 15 (FILTER: unlearned)
- **skills_languages**: MAX 11 (FILTER: languages)
- **weapons_main**: MAX 30

### page3_spell.html (Spell Page)
- **spells_column1**: MAX 12 (left column)
- **spells_column2**: MAX 12 (right column)
- **magic_items**: MAX 5
- **Total**: 24 spells per page

### page4_equip.html (Equipment Page)
- **equipment_sections**: MAX 20
- **game_results**: MAX 10

## How Pagination Works

1. **Capacity Calculation**: Paginator reads template metadata to determine capacity
2. **Distribution**: Items are distributed across blocks according to MaxItems
3. **Page Creation**: Multiple pages are created when capacity is exceeded
4. **Column Handling**: Items fill first column, then overflow to second column
5. **Page Overflow**: Remaining items continue on next page

## Algorithm Details

### Single Column Distribution
```
Items: 50, Capacity per block: 30
Result:
  Page 1: Block 1 (30 items)
  Page 2: Block 1 (20 items)
```

### Multi-Column Distribution
```
Items: 100, Column 1: 32, Column 2: 32 (64 per page)
Result:
  Page 1: Column 1 (32) + Column 2 (32) = 64
  Page 2: Column 1 (32) + Column 2 (4) = 36
```

## Best Practices

1. **Always check errors** when paginating
2. **Pre-calculate page count** for UI/UX planning
3. **Loop through all pages** to generate complete PDF sets
4. **Maintain page numbers** in metadata for proper labeling
5. **Type assert carefully** when extracting data from PageDistribution.Data

## Error Handling

```go
pages, err := paginator.PaginateSkills(skills, "page1_stats.html", "")
if err != nil {
    // Handle errors:
    // - Template not found
    // - No capacity defined for list type
    // - Invalid template configuration
    log.Printf("Pagination failed: %v", err)
    return
}

if len(pages) == 0 {
    log.Println("No pages needed (empty list)")
    return
}
```

## Testing

All pagination logic is thoroughly tested:

- ✓ Single page scenarios
- ✓ Multi-column distribution
- ✓ Multi-page overflow
- ✓ Empty lists
- ✓ Invalid templates
- ✓ Capacity calculations
- ✓ Integration with PDF generation

Run tests:
```bash
go test -v ./pdfrender/... -run TestPaginate
```

## Performance Considerations

- **Memory**: Pagination creates slices, minimal memory overhead
- **Speed**: O(n) complexity, very fast even for large lists
- **Caching**: Template metadata is parsed once and reused
- **Scalability**: Handles thousands of items efficiently

## Future Enhancements

Potential improvements for consideration:

1. **Filter-based pagination**: Separate learned/unlearned skills
2. **Custom sorting**: Order items before pagination
3. **Dynamic capacity**: Adjust based on font size or layout
4. **Partial rendering**: Generate only specific pages on demand
5. **Merge PDFs**: Combine multiple page PDFs into single document
