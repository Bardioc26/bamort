package validator

import (
	"bamort/logger"
	"fmt"

	"gorm.io/gorm"
)

// SchemaValidator validates database schema integrity
type SchemaValidator struct {
	DB *gorm.DB
}

// ValidationReport contains validation results
type ValidationReport struct {
	Success        bool
	TablesChecked  int
	TablesValid    int
	Errors         []string
	Warnings       []string
	MissingTables  []string
	MissingColumns map[string][]string
}

// NewValidator creates a new schema validator
func NewValidator(db *gorm.DB) *SchemaValidator {
	return &SchemaValidator{
		DB: db,
	}
}

// Validate performs comprehensive schema validation
func (v *SchemaValidator) Validate() (*ValidationReport, error) {
	report := &ValidationReport{
		Success:        true,
		MissingColumns: make(map[string][]string),
	}

	logger.Info("Starting database schema validation...")

	// Check ALL tables must exist for the application to work properly
	// If char_*, equi_*, audit_* tables are missing, /api/maintenance/setupcheck must be called
	// So it's best to ensure all tables are present
	criticalTables := []string{
		// System tables
		"schema_version",
		"migration_history",
		"users",

		// Audit tables
		"audit_log_entries",

		// Character tables
		"char_bennies",
		"char_characteristics",
		"char_char_creation_session",
		"char_chars",
		"char_eigenschaften",
		"char_endurances",
		"char_experiances",
		"char_health",
		"char_motionranges",
		"char_skills",
		"char_spells",
		"char_wealth",
		"char_weaponskills",

		// Equipment tables
		"equi_containers",
		"equi_equipments",
		"equi_weapons",

		// GSM Master Data tables
		"gsm_believes",
		"gsm_cc_class_category_points",
		"gsm_cc_class_spell_points",
		"gsm_cc_class_typical_skills",
		"gsm_cc_class_typical_spells",
		"gsm_character_classes",
		"gsm_containers",
		"gsm_equipments",
		"gsm_lit_sources",
		"gsm_misc",
		"gsm_skills",
		"gsm_spells",
		"gsm_transportations",
		"gsm_weapons",
		"gsm_weaponskills",

		// Learning system tables
		"learning_class_category_ep_costs",
		"learning_class_spell_school_ep_costs",
		"learning_skill_categories",
		"learning_skill_category_difficulties",
		"learning_skill_difficulties",
		"learning_skill_improvement_costs",
		"learning_spell_level_le_costs",
		"learning_spell_schools",
		"learning_weaponskill_category_difficulties",
	}

	for _, table := range criticalTables {
		report.TablesChecked++
		if v.tableExists(table) {
			report.TablesValid++
			logger.Debug("✓ Table exists: %s", table)
		} else {
			report.MissingTables = append(report.MissingTables, table)
			report.Errors = append(report.Errors, fmt.Sprintf("Missing table: %s", table))
			report.Success = false
			logger.Error("✗ Missing table: %s", table)
		}
	}

	// Check schema_version table structure
	if v.tableExists("schema_version") {
		requiredColumns := []string{"id", "version", "migration_number", "applied_at"}
		missingCols := v.checkTableColumns("schema_version", requiredColumns)
		if len(missingCols) > 0 {
			report.MissingColumns["schema_version"] = missingCols
			report.Errors = append(report.Errors,
				fmt.Sprintf("schema_version missing columns: %v", missingCols))
			report.Success = false
		}
	}

	// Check migration_history table structure
	if v.tableExists("migration_history") {
		requiredColumns := []string{"id", "migration_number", "description", "applied_at"}
		missingCols := v.checkTableColumns("migration_history", requiredColumns)
		if len(missingCols) > 0 {
			report.MissingColumns["migration_history"] = missingCols
			report.Errors = append(report.Errors,
				fmt.Sprintf("migration_history missing columns: %v", missingCols))
			report.Success = false
		}
	}

	// Validate record counts are reasonable
	if err := v.validateDataIntegrity(report); err != nil {
		report.Warnings = append(report.Warnings, fmt.Sprintf("Data integrity check: %v", err))
	}

	if report.Success {
		logger.Info("✓ Schema validation passed")
	} else {
		logger.Error("✗ Schema validation failed with %d error(s)", len(report.Errors))
	}

	return report, nil
}

// tableExists checks if a table exists in the database
func (v *SchemaValidator) tableExists(tableName string) bool {
	return v.DB.Migrator().HasTable(tableName)
}

// checkTableColumns verifies that required columns exist in a table
func (v *SchemaValidator) checkTableColumns(tableName string, requiredColumns []string) []string {
	var missing []string

	for _, col := range requiredColumns {
		if !v.DB.Migrator().HasColumn(tableName, col) {
			missing = append(missing, col)
		}
	}

	return missing
}

// validateDataIntegrity performs basic sanity checks on data
func (v *SchemaValidator) validateDataIntegrity(report *ValidationReport) error {
	// Check that schema_version has at least one entry
	if v.tableExists("schema_version") {
		var count int64
		if err := v.DB.Table("schema_version").Count(&count).Error; err != nil {
			return fmt.Errorf("failed to count schema_version records: %w", err)
		}
		if count == 0 {
			report.Warnings = append(report.Warnings, "schema_version table is empty")
		}
	}

	// Check for orphaned records (basic check)
	if v.tableExists("chars") && v.tableExists("users") {
		var orphanedChars int64
		if err := v.DB.Raw(`
			SELECT COUNT(*) FROM chars 
			WHERE user_id NOT IN (SELECT id FROM users)
		`).Scan(&orphanedChars).Error; err == nil {
			if orphanedChars > 0 {
				report.Warnings = append(report.Warnings,
					fmt.Sprintf("Found %d orphaned characters (invalid user_id)", orphanedChars))
			}
		}
	}

	return nil
}

// ValidatePostMigration performs post-migration validation
func (v *SchemaValidator) ValidatePostMigration() error {
	logger.Info("Performing post-migration validation...")

	report, err := v.Validate()
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if !report.Success {
		return fmt.Errorf("validation found %d error(s): %v", len(report.Errors), report.Errors)
	}

	if len(report.Warnings) > 0 {
		logger.Warn("Validation passed with %d warning(s):", len(report.Warnings))
		for _, w := range report.Warnings {
			logger.Warn("  - %s", w)
		}
	}

	logger.Info("✓ Post-migration validation successful")
	return nil
}
