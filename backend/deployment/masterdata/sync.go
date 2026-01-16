package masterdata

import (
	"bamort/gsmaster"
	"bamort/logger"
	"fmt"

	"gorm.io/gorm"
)

// MasterDataSync orchestrates master data synchronization
type MasterDataSync struct {
	ImportDir string
	DB        *gorm.DB
	DryRun    bool
	Verbose   bool
}

// NewMasterDataSync creates a new master data sync instance
func NewMasterDataSync(db *gorm.DB, importDir string) *MasterDataSync {
	return &MasterDataSync{
		ImportDir: importDir,
		DB:        db,
		DryRun:    false,
		Verbose:   false,
	}
}

// SyncAll synchronizes all master data in dependency order
func (s *MasterDataSync) SyncAll() error {
	logger.Info("Starting master data synchronization from %s", s.ImportDir)

	if s.DryRun {
		logger.Info("[DRY RUN] No changes will be made")
	}

	// Import in dependency order (no dependencies → dependencies)
	steps := []struct {
		Name     string
		ImportFn func() error
	}{
		{"Sources", s.importSources},
		{"Character Classes", s.importCharacterClasses},
		{"Skill Categories", s.importSkillCategories},
		{"Skill Difficulties", s.importSkillDifficulties},
		{"Spell Schools", s.importSpellSchools},
		{"Skills", s.importSkills},
		{"Weapon Skills", s.importWeaponSkills},
		{"Spells", s.importSpells},
		{"Equipment", s.importEquipment},
		{"Learning Costs", s.importLearningCosts},
	}

	for _, step := range steps {
		if s.Verbose {
			logger.Info("Importing %s...", step.Name)
		}

		if s.DryRun {
			logger.Info("[DRY RUN] Would import %s", step.Name)
			continue
		}

		if err := step.ImportFn(); err != nil {
			return fmt.Errorf("failed to import %s: %w", step.Name, err)
		}

		if s.Verbose {
			logger.Info("✓ %s imported successfully", step.Name)
		}
	}

	logger.Info("Master data synchronization completed successfully")
	return nil
}

// Import functions delegate to existing gsmaster package
func (s *MasterDataSync) importSources() error {
	return gsmaster.ImportSources(s.ImportDir)
}

func (s *MasterDataSync) importCharacterClasses() error {
	return gsmaster.ImportCharacterClasses(s.ImportDir)
}

func (s *MasterDataSync) importSkillCategories() error {
	return gsmaster.ImportSkillCategories(s.ImportDir)
}

func (s *MasterDataSync) importSkillDifficulties() error {
	return gsmaster.ImportSkillDifficulties(s.ImportDir)
}

func (s *MasterDataSync) importSpellSchools() error {
	return gsmaster.ImportSpellSchools(s.ImportDir)
}

func (s *MasterDataSync) importSkills() error {
	return gsmaster.ImportSkills(s.ImportDir)
}

func (s *MasterDataSync) importWeaponSkills() error {
	return gsmaster.ImportWeaponSkills(s.ImportDir)
}

func (s *MasterDataSync) importSpells() error {
	return gsmaster.ImportSpells(s.ImportDir)
}

func (s *MasterDataSync) importEquipment() error {
	return gsmaster.ImportEquipment(s.ImportDir)
}

func (s *MasterDataSync) importLearningCosts() error {
	return gsmaster.ImportSkillImprovementCosts(s.ImportDir)
}
