package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var migrationDone bool
var isTestDb bool

// SetupTestDB creates an in-memory SQLite database for testing
func setupTestDB(opts ...bool) {
	isTestDb = true
	if len(opts) > 0 {
		isTestDb = opts[0]
	}
	if database.DB == nil {
		var db *gorm.DB
		var err error
		if isTestDb {
			//*
			db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
			if err != nil {
				panic("failed to connect to the test database")
			}
			//*/
		} else {
			//* //testing with persistent MariaDB
			dsn := os.Getenv("TEST_DB_DSN")
			if dsn == "" {
				dsn = "bamort:password@tcp(localhost:3306)/bamort_test?charset=utf8mb4&parseTime=True&loc=Local"
			}
			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				panic("failed to connect to the test database")
			}
			//*/
			migrationDone = true
		}
		database.DB = db
	}
	if !migrationDone {
		err := models.MigrateStructure()
		if err != nil {
			panic("failed to MigrateStructure")
		}
		migrationDone = true
	}
}

func TestLoadLevelingConfig(t *testing.T) {
	// Save original Config
	originalConfig := Config

	// Test invalid file path
	os.Setenv("CONFIG_PATH", "/invalid/path/leveldata.json")
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with invalid file path")
		}
		// Restore original Config
		Config = originalConfig
	}()

	// Call init() which should panic
	loadLevelingConfig("/invalid/path/leveldata.json")
}

func TestInitValidConfig(t *testing.T) {
	loadLevelingConfig()
	setupTestDB(false)
	/*
		// Save original Config
		originalConfig := Config
		defer func() {
			Config = originalConfig
		}()
	*/

	// Test with valid config file
	os.Setenv("CONFIG_PATH", "/data/dev/bamort/config/leveldata.json")
	loadLevelingConfig("../testdata/leveldata.json")

	// Verify Config was populated
	assert.LessOrEqual(t, 1, len(Config.BaseLearnCost), "Expected BaseLearnCost to be populated")
	assert.LessOrEqual(t, 1, len(Config.ImprovementCost), "Expected BaseLearImprovementCostnCost to be populated")
	assert.LessOrEqual(t, 1, len(Config.EPPerTE), "Expected EPPerTE to be populated")
}

func TestCalculateSpellLearnCost(t *testing.T) {
	loadLevelingConfig()
	setupTestDB(false)
	/*
		// Save original Config
		originalConfig := Config
		defer func() {
			Config = originalConfig
		}()

		// Set up test config
		Config = LevelConfig{
			SpellLearnCost: map[int]int{
				1: 1,
				2: 2,
			},
			SpellEPPerSchoolByClass: map[CharClass]map[string]int{
				"Magier": {"Beweg": 10},
				"Elfe":   {"Beweg": 15},
			},
			AllowedSchools: map[CharClass]map[string]bool{
				"Magier": {"Beweg": true},
				"Elfe":   {"Beweg": true},
			},
		}
	*/

	tests := []struct {
		name  string
		spell SpellDefinition
		//class       CharClass
		class       string
		wantEP      int
		wantErr     bool
		errContains string
	}{
		{
			name:   "valid spell for magier",
			spell:  SpellDefinition{Name: "Angst", Stufe: 2, School: "Beherrschen"},
			class:  "Magier",
			wantEP: 180, // 1 LE * (10 EP * 3)
		},
		{
			name:    "valid spell for elf",
			spell:   SpellDefinition{Name: "Angst", Stufe: 2, School: "Beherrschen"},
			class:   "Elfe",
			wantEP:  51, // (1 LE * (15 EP * 3)) + 6
			wantErr: true,
		},
		{
			name:        "invalid spell level",
			spell:       SpellDefinition{Name: "Angst", Stufe: 99, School: "Beherrschen"},
			class:       "Magier",
			wantErr:     true,
			errContains: "ungültige Zauberstufe",
		},
		{
			name:        "invalid class",
			spell:       SpellDefinition{Name: "Angst", Stufe: 2, School: "Beherrschen"},
			class:       "InvalidClass",
			wantErr:     true,
			errContains: "keine EP-Tabelle für Klasse",
		},
		{
			name:        "invalid school",
			spell:       SpellDefinition{Name: "Angst", Stufe: 2, School: "Beherrschen"},
			class:       "Magier",
			wantErr:     true,
			errContains: "unbekannte Schule",
		},
		{
			name:        "not allowed school",
			spell:       SpellDefinition{Name: "Angst", Stufe: 2, School: "Beherrschen"},
			class:       "Krieger",
			wantErr:     true,
			errContains: "darf die Schule",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateSpellLearnCost(tt.spell.Name, tt.class)
			if tt.wantErr {
				if err == nil {
					assert.Error(t, err, "CalculateSpellLearnCost() expected error")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					assert.Error(t, err, "CalculateSpellLearnCost() error = %v, should contain %v", err, tt.errContains)
					//t.Errorf("CalculateSpellLearnCost() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}
			if err != nil {
				assert.NoError(t, err, "CalculateSpellLearnCost() unexpected error = %v", err)
				//t.Errorf("CalculateSpellLearnCost() unexpected error = %v", err)
				return
			}
			if got != tt.wantEP {
				assert.Equal(t, tt.wantEP, got, "CalculateSpellLearnCost() = %v, want %v", got, tt.wantEP)
				//t.Errorf("CalculateSpellLearnCost() = %v, want %v", got, tt.wantEP)
			}
		})
	}
}

func TestCalculateLearnCost(t *testing.T) {
	loadLevelingConfig()
	setupTestDB(false)
	// Save original Config
	/*
		originalConfig := Config
		defer func() {
			Config = originalConfig
		}()
	*/

	// Set up test config
	/*
		Config = LevelConfig{
			BaseLearnCost: map[SkillGroup]map[Difficulty]int{
				"Alltag": {
					"leicht":      1,
					"normal":      1,
					"schwer":      2,
					"sehr_schwer": 10,
				},
			},
			EPPerTE: map[CharClass]map[SkillGroup]int{
				"Krieger": {"Alltag": 20},
				"Elf":     {"Alltag": 30},
			},
		}
	*/

	tests := []struct {
		name        string
		skill       SkillDefinition
		class       string
		wantEP      int
		wantErr     bool
		errContains string
	}{
		{
			name: "valid skill for warrior",
			skill: SkillDefinition{
				Name:       "Bootfahren",
				Group:      "Alltag",
				Difficulty: "leicht",
			},
			class:  "Krieger",
			wantEP: 60, // 1 LE * (20 EP * 3)
		},
		{
			name: "valid skill for elf",
			skill: SkillDefinition{
				Name:       "Bootfahren",
				Group:      "Alltag",
				Difficulty: "leicht",
			},
			class:   "Elf",
			wantEP:  96, // (1 LE * (30 EP * 3)) + 6
			wantErr: true,
		},
		{
			name: "invalid group",
			skill: SkillDefinition{
				Name:       "Erste Hilfe",
				Group:      "InvalidGroup",
				Difficulty: "leicht",
			},
			class:       "Krieger",
			wantErr:     true,
			errContains: "unbekannte Gruppe",
		},
		{
			name: "invalid difficulty",
			skill: SkillDefinition{
				Name:       "Geländelauf",
				Group:      "Körper",
				Difficulty: "invalid",
			},
			class:       "Krieger",
			wantErr:     true,
			errContains: "keine LE-Definition für diese Schwierigkeit",
		},
		{
			name: "invalid class",
			skill: SkillDefinition{
				Name:       "Test",
				Group:      "Alltag",
				Difficulty: "leicht",
			},
			class:       "InvalidClass",
			wantErr:     true,
			errContains: "keine EP-Kosten für",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateSkillLearnCost(tt.skill.Name, tt.class)
			if tt.wantErr {
				if err == nil {
					assert.Error(t, err, "CalculateLearnCost() expected error")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					assert.Error(t, err, "CalculateLearnCost() error = %v, should contain %v", err, tt.errContains)
					//t.Errorf("CalculateLearnCost() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}
			if err != nil {
				assert.NoError(t, err, "CalculateLearnCost() unexpected error = %v", err)
				//t.Errorf("CalculateLearnCost() unexpected error = %v", err)
				return
			}
			if got != tt.wantEP {
				assert.Equal(t, tt.wantEP, got, "CalculateLearnCost() = %v, want %v", got, tt.wantEP)
				//t.Errorf("CalculateLearnCost() = %v, want %v", got, tt.wantEP)
			}
		})
	}
}

func TestCalculateImprovementCost(t *testing.T) {
	loadLevelingConfig()
	setupTestDB(false)
	/*
		// Save original Config
		originalConfig := Config
		defer func() {
			Config = originalConfig
		}()


		// Set up test config
		Config = LevelConfig{
			ImprovementCost: map[SkillGroup]map[Difficulty]map[string]int{
				"Alltag": {
					"leicht": {
						"9":  1,
						"10": 1,
					},
					"normal": {
						"9":  2,
						"10": 2,
					},
				},
			},
			EPPerTE: map[CharClass]map[SkillGroup]int{
				"Krieger": {"Alltag": 20},
				"Magier":  {"Alltag": 30},
			},
		}
	*/
	tests := []struct {
		name         string
		skill        SkillDefinition
		class        CharClass
		currentLevel int
		wantEP       int
		wantErr      bool
		errContains  string
	}{
		{
			name: "valid improvement for warrior",
			skill: SkillDefinition{
				Name:       "Bootfahren",
				Group:      "Alltag",
				Difficulty: "leicht",
			},
			class:        "Krieger",
			currentLevel: 13,
			wantEP:       40, // 1 LE * 20 EP
		},
		{
			name: "valid improvement for mage",
			skill: SkillDefinition{
				Name:       "Schreiben",
				Group:      "Wissen",
				Difficulty: "normal",
			},
			class:        "Hexer",
			currentLevel: 9,
			wantEP:       20, // 2 LE * 30 EP
		},
		{
			name: "invalid group",
			skill: SkillDefinition{
				Name:       "Erste Hilfe",
				Group:      "InvalidGroup",
				Difficulty: "leicht",
			},
			class:        "Krieger",
			currentLevel: 8,
			wantErr:      true,
			errContains:  "keine Improvement-Daten für diese Gruppe",
		},
		{
			name: "invalid difficulty",
			skill: SkillDefinition{
				Name:       "Geländelauf",
				Group:      "Körper",
				Difficulty: "invalid",
			},
			class:        "Krieger",
			currentLevel: 8,
			wantErr:      true,
			errContains:  "keine Improvement-Daten für diese Schwierigkeit",
		},
		{
			name: "invalid next level",
			skill: SkillDefinition{
				Name:       "Schreiben",
				Group:      "Wissen",
				Difficulty: "normal",
			},
			class:        "Krieger",
			currentLevel: 99,
			wantErr:      true,
			errContains:  "kein Eintrag für Bonus",
		},
		{
			name: "invalid class",
			skill: SkillDefinition{
				Name:       "Schreiben",
				Group:      "Wissen",
				Difficulty: "normal",
			},
			class:        "InvalidClass",
			currentLevel: 8,
			wantErr:      true,
			errContains:  "keine EP-Kosten für",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateImprovementCost(tt.skill.Name, string(tt.class), tt.currentLevel)
			if tt.wantErr {
				if err == nil {
					assert.Error(t, err, "CalculateImprovementCost() expected error")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					assert.Error(t, err, "CalculateImprovementCost() error = %v, should contain %v", err, tt.errContains)
					//t.Errorf("CalculateImprovementCost() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}
			if err != nil {
				assert.NoError(t, err, "CalculateImprovementCost() unexpected error = %v", err)
				//t.Errorf("CalculateImprovementCost() unexpected error = %v", err)
				return
			}
			if got.Ep != tt.wantEP {
				assert.Equal(t, tt.wantEP, got, "CalculateImprovementCost() = %v, want %v", got, tt.wantEP)
				//t.Errorf("CalculateImprovementCost() = %v, want %v", got, tt.wantEP)
			}
		})
	}
}
