package importer

import (
	"testing"
)

func TestValidation_RequiredFieldsRule(t *testing.T) {
	rule := &RequiredFieldsRule{}

	// Test valid character
	char := CharacterImport{
		Name: "Test Character",
	}
	bmrt := NewBMRTCharacter(char, "test-adapter", "test-format")

	result := rule.Validate(bmrt)
	if !result.Valid {
		t.Errorf("Character with name should be valid, got errors: %v", result.Errors)
	}

	// Test invalid character (no name)
	charInvalid := CharacterImport{}
	bmrtInvalid := NewBMRTCharacter(charInvalid, "test-adapter", "test-format")

	result = rule.Validate(bmrtInvalid)
	if result.Valid {
		t.Error("Character without name should be invalid")
	}
	if len(result.Errors) == 0 {
		t.Error("Should have validation errors")
	}
}

func TestValidation_BmrtVersionRule(t *testing.T) {
	rule := &BmrtVersionRule{}

	// Test valid version
	char := CharacterImport{Name: "Test"}
	bmrt := NewBMRTCharacter(char, "test-adapter", "test-format")

	result := rule.Validate(bmrt)
	if !result.Valid {
		t.Errorf("Valid BMRT version should pass, got errors: %v", result.Errors)
	}

	// Test invalid version
	bmrt.BmrtVersion = "99.9"
	result = rule.Validate(bmrt)
	if result.Valid {
		t.Error("Invalid BMRT version should fail validation")
	}
}

func TestValidation_StatsRangeRule(t *testing.T) {
	rule := &StatsRangeRule{}

	// Test valid stats
	char := CharacterImport{
		Name: "Test",
		Eigenschaften: Eigenschaften{
			St: 50,
			Gw: 60,
			Ko: 70,
			In: 80,
			Zt: 90,
		},
	}
	bmrt := NewBMRTCharacter(char, "test-adapter", "test-format")

	result := rule.Validate(bmrt)
	if !result.Valid {
		t.Errorf("Valid stats should pass, got errors: %v", result.Errors)
	}

	// Test invalid stats (out of range)
	charInvalid := CharacterImport{
		Name: "Test",
		Eigenschaften: Eigenschaften{
			St: 150, // Out of range - should only warn, not fail
			Gw: 60,
		},
	}
	bmrtInvalid := NewBMRTCharacter(charInvalid, "test-adapter", "test-format")

	result = rule.Validate(bmrtInvalid)
	// Stats > 100 are warnings, not errors
	if !result.Valid {
		t.Error("High stats should be valid (warning only)")
	}
	if len(result.Warnings) == 0 {
		t.Error("High stats should generate warnings")
	}

	// Test negative stats (should fail)
	charNegative := CharacterImport{
		Name: "Test",
		Eigenschaften: Eigenschaften{
			St: -10, // Negative - should fail
		},
	}
	bmrtNegative := NewBMRTCharacter(charNegative, "test-adapter", "test-format")

	result = rule.Validate(bmrtNegative)
	if result.Valid {
		t.Error("Negative stats should fail validation")
	}
}

func TestValidator_ValidateCharacter(t *testing.T) {
	validator := NewValidator()

	// Add rules
	validator.AddRule(&RequiredFieldsRule{})
	validator.AddRule(&BmrtVersionRule{})

	// Test valid character
	char := CharacterImport{
		Name: "Test Character",
	}
	bmrt := NewBMRTCharacter(char, "test-adapter", "test-format")

	result := validator.ValidateCharacter(bmrt)
	if !result.Valid {
		t.Errorf("Valid character should pass validation, errors: %v", result.Errors)
	}

	// Test invalid character
	charInvalid := CharacterImport{} // Missing name
	bmrtInvalid := NewBMRTCharacter(charInvalid, "test-adapter", "test-format")
	bmrtInvalid.BmrtVersion = "99.9" // Invalid version

	result = validator.ValidateCharacter(bmrtInvalid)
	if result.Valid {
		t.Error("Invalid character should fail validation")
	}
	if len(result.Errors) < 2 {
		t.Errorf("Expected at least 2 errors, got %d", len(result.Errors))
	}
}

func TestValidator_Warnings(t *testing.T) {
	validator := NewValidator()

	validator.AddRule(&testWarningRule{})

	char := CharacterImport{Name: "Test"}
	bmrt := NewBMRTCharacter(char, "test-adapter", "test-format")

	result := validator.ValidateCharacter(bmrt)
	if !result.Valid {
		t.Error("Character should still be valid with warnings")
	}
	if len(result.Warnings) == 0 {
		t.Error("Should have warnings")
	}
}

// testWarningRule is a helper rule for testing warnings
type testWarningRule struct{}

func (r *testWarningRule) Validate(char *BMRTCharacter) ValidationResult {
	return ValidationResult{
		Valid: true,
		Warnings: []ValidationWarning{
			{
				Field:   "test_field",
				Message: "This is a warning",
				Source:  "test",
			},
		},
	}
}

func TestValidationResult_Combine(t *testing.T) {
	result1 := ValidationResult{
		Valid: true,
		Errors: []ValidationError{
			{Field: "field1", Message: "error1", Source: "source1"},
		},
	}

	result2 := ValidationResult{
		Valid: true,
		Warnings: []ValidationWarning{
			{Field: "field2", Message: "warning1", Source: "source2"},
		},
	}

	combined := CombineValidationResults(result1, result2)

	if len(combined.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(combined.Errors))
	}
	if len(combined.Warnings) != 1 {
		t.Errorf("Expected 1 warning, got %d", len(combined.Warnings))
	}
	if !combined.Valid {
		t.Error("Combined result should be valid when both are valid")
	}

	// Test with one invalid
	result3 := ValidationResult{
		Valid: false,
		Errors: []ValidationError{
			{Field: "field3", Message: "error3", Source: "source3"},
		},
	}

	combined = CombineValidationResults(result1, result3)
	if combined.Valid {
		t.Error("Combined result should be invalid when one is invalid")
	}
}
