package importer

// ValidationResult represents the outcome of validating a character
type ValidationResult struct {
	Valid    bool                `json:"valid"`
	Errors   []ValidationError   `json:"errors,omitempty"`
	Warnings []ValidationWarning `json:"warnings,omitempty"`
	Source   string              `json:"source"` // "bmrt", "gamesystem", "adapter"
}

// ValidationError represents a validation error that prevents import
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Source  string `json:"source"` // Which validation phase found this
}

// ValidationWarning represents a non-blocking validation issue
type ValidationWarning struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Source  string `json:"source"`
}

// ValidationRule defines the interface for validation rules
type ValidationRule interface {
	Validate(char *BMRTCharacter) ValidationResult
}

// Validator manages validation rules and performs character validation
type Validator struct {
	rules []ValidationRule
}

// NewValidator creates a new validator with default rules
func NewValidator() *Validator {
	return &Validator{
		rules: make([]ValidationRule, 0),
	}
}

// AddRule adds a validation rule
func (v *Validator) AddRule(rule ValidationRule) {
	v.rules = append(v.rules, rule)
}

// ValidateCharacter runs all validation rules and combines results
func (v *Validator) ValidateCharacter(char *BMRTCharacter) ValidationResult {
	combined := ValidationResult{
		Valid:    true,
		Errors:   make([]ValidationError, 0),
		Warnings: make([]ValidationWarning, 0),
	}

	for _, rule := range v.rules {
		result := rule.Validate(char)
		combined = CombineValidationResults(combined, result)
	}

	return combined
}

// CombineValidationResults combines multiple validation results
func CombineValidationResults(results ...ValidationResult) ValidationResult {
	combined := ValidationResult{
		Valid:    true,
		Errors:   make([]ValidationError, 0),
		Warnings: make([]ValidationWarning, 0),
	}

	for _, result := range results {
		if !result.Valid {
			combined.Valid = false
		}
		combined.Errors = append(combined.Errors, result.Errors...)
		combined.Warnings = append(combined.Warnings, result.Warnings...)
	}

	return combined
}

// ========================================
// Phase 1: BMRT Structural Validation
// ========================================

// RequiredFieldsRule validates that required fields are present
type RequiredFieldsRule struct{}

func (r *RequiredFieldsRule) Validate(char *BMRTCharacter) ValidationResult {
	result := ValidationResult{
		Valid:  true,
		Source: "bmrt",
	}

	if char.Name == "" {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "name",
			Message: "Character name is required",
			Source:  "bmrt",
		})
	}

	return result
}

// BmrtVersionRule validates the BMRT version is supported
type BmrtVersionRule struct{}

func (r *BmrtVersionRule) Validate(char *BMRTCharacter) ValidationResult {
	result := ValidationResult{
		Valid:  true,
		Source: "bmrt",
	}

	if !ValidateBMRTVersion(char.BmrtVersion) {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "bmrt_version",
			Message: "Unsupported BMRT version: " + char.BmrtVersion,
			Source:  "bmrt",
		})
	}

	return result
}

// ========================================
// Phase 2: Game System Semantic Validation
// ========================================

// StatsRangeRule validates game system stats are within valid ranges
type StatsRangeRule struct{}

func (r *StatsRangeRule) Validate(char *BMRTCharacter) ValidationResult {
	result := ValidationResult{
		Valid:  true,
		Source: "gamesystem",
	}

	// For Midgard, stats should be 0-100 (though can exceed in rare cases)
	stats := map[string]int{
		"St": char.Eigenschaften.St,
		"Gw": char.Eigenschaften.Gw,
		"Ko": char.Eigenschaften.Ko,
		"In": char.Eigenschaften.In,
		"Zt": char.Eigenschaften.Zt,
		"Au": char.Eigenschaften.Au,
		"Pa": char.Eigenschaften.Pa,
		"Wk": char.Eigenschaften.Wk,
		"Gs": char.Eigenschaften.Gs,
	}

	for name, value := range stats {
		if value < 0 {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:   "eigenschaften." + name,
				Message: "Stat cannot be negative",
				Source:  "gamesystem",
			})
		} else if value > 100 {
			// Warning only - some characters can have stats > 100
			result.Warnings = append(result.Warnings, ValidationWarning{
				Field:   "eigenschaften." + name,
				Message: "Stat value unusually high (> 100)",
				Source:  "gamesystem",
			})
		}
	}

	return result
}

// ReferentialIntegrityRule validates that referenced items exist
// This is a placeholder - full implementation would check against game system master data
type ReferentialIntegrityRule struct{}

func (r *ReferentialIntegrityRule) Validate(char *BMRTCharacter) ValidationResult {
	result := ValidationResult{
		Valid:  true,
		Source: "gamesystem",
	}

	// Placeholder validation
	// In full implementation, would check:
	// - Skills reference valid skill categories
	// - Spells exist in game system
	// - Equipment types are valid
	// etc.

	return result
}
