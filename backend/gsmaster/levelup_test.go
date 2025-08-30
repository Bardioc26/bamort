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

func TestCalculateImprovementCost(t *testing.T) {
	//loadLevelingConfigOld()
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
