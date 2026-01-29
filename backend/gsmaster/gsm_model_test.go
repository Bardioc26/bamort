package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkill_Create(t *testing.T) {
	database.SetupTestDB()
	testDefinition := []struct {
		name    string
		skill   *models.Skill
		wantErr bool
	}{
		{
			name: "valid skill",
			skill: &models.Skill{
				Name:             "Test Skill",
				Beschreibung:     "Test Description",
				Quelle:           "Test Source",
				Initialwert:      5,
				Bonuseigenschaft: "st",
				Improvable:       true,
				InnateSkill:      false,
				Category:         "Test Category",
			},
			wantErr: false,
		},
		/*{
			name:    "nil skill",
			skill:   nil,
			wantErr: true,
		},*/
	}

	for _, tt := range testDefinition {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.skill.Create()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, "midgard", tt.skill.GameSystem)
			assert.NotZero(t, tt.skill.GameSystemId)
			assert.NotZero(t, tt.skill.ID)
		})
	}
	database.ResetTestDB()
}

func TestWeaponSkill_Create(t *testing.T) {
	database.SetupTestDB()
	testDefinition := []struct {
		name        string
		weaponSkill *models.WeaponSkill
		wantErr     bool
	}{
		{
			name: "valid weapon skill",
			weaponSkill: &models.WeaponSkill{
				Skill: models.Skill{
					Name:             "Test Weapon Skill",
					Beschreibung:     "Test Description",
					Quelle:           "Test Source",
					Initialwert:      5,
					Bonuseigenschaft: "st",
					Improvable:       true,
					InnateSkill:      false,
					Category:         "Test Category",
				},
			},
			wantErr: false,
		},
		/*{
			name:        "nil weapon skill",
			weaponSkill: nil,
			wantErr:     true,
		},*/
	}

	for _, tt := range testDefinition {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.weaponSkill.Create()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, "midgard", tt.weaponSkill.GameSystem)
			assert.NotZero(t, tt.weaponSkill.GameSystemId)
			assert.NotZero(t, tt.weaponSkill.ID)
		})
	}
	database.ResetTestDB()
}

func TestSkill_TableName(t *testing.T) {
	database.SetupTestDB()
	skill := &models.Skill{}
	expected := "gsm_skills"
	assert.Equal(t, expected, skill.TableName())
	database.ResetTestDB()
}

func TestSkill_First(t *testing.T) {
	database.SetupTestDB()
	testDefinition := []struct {
		name     string
		skill    *models.Skill
		findName string
		wantErr  bool
	}{
		{
			name: "existing skill",
			skill: &models.Skill{
				Name:         "Test Skill",
				GameSystemId: 1,
			},
			findName: "Test Skill",
			wantErr:  false,
		},
		{
			name:     "non-existing skill",
			skill:    &models.Skill{},
			findName: "NonExistent",
			wantErr:  true,
		},
		{
			name:     "empty name",
			skill:    &models.Skill{},
			findName: "",
			wantErr:  true,
		},
	}

	for _, tt := range testDefinition {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				// Create test data first
				err := tt.skill.Create()
				assert.NoError(t, err)
			}

			s := &models.Skill{}
			err := s.First(tt.findName)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.skill.Name, s.Name)
			assert.Equal(t, "midgard", s.GameSystem)
			assert.NotZero(t, s.GameSystemId)
		})
	}
	database.ResetTestDB()
}

func TestSkill_FirstId(t *testing.T) {
	database.SetupTestDB()
	testDefinition := []struct {
		name    string
		skill   *models.Skill
		findId  uint
		wantErr bool
	}{
		{
			name: "existing skill",
			skill: &models.Skill{
				Name:         "Test Skill",
				GameSystemId: 1,
			},
			findId:  1,
			wantErr: false,
		},
		{
			name:    "non-existing id",
			skill:   &models.Skill{},
			findId:  9999,
			wantErr: true,
		},
		{
			name:    "zero id",
			skill:   &models.Skill{},
			findId:  0,
			wantErr: true,
		},
	}

	for _, tt := range testDefinition {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				// Create test data first
				err := tt.skill.Create()
				assert.NoError(t, err)
				tt.findId = tt.skill.ID
			}

			s := &models.Skill{}
			err := s.FirstId(tt.findId)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.skill.Name, s.Name)
			assert.Equal(t, "midgard", s.GameSystem)
			assert.NotZero(t, s.GameSystemId)
			assert.Equal(t, tt.findId, s.ID)
		})
	}
	database.ResetTestDB()
}

func TestSkill_Save(t *testing.T) {
	database.SetupTestDB()
	testDefinition := []struct {
		name    string
		skill   *models.Skill
		wantErr bool
	}{
		{
			name: "update existing skill",
			skill: &models.Skill{
				Name:         "Test Skill",
				Beschreibung: "Original Description",
				GameSystem:   "midgard",
			},
			wantErr: false,
		},
		/*{
			name:    "nil skill",
			skill:   nil,
			wantErr: true,
		},*/
	}

	for _, tt := range testDefinition {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				// Create initial record
				err := tt.skill.Create()
				assert.NoError(t, err)

				// Modify the skill
				tt.skill.Beschreibung = "Updated Description"
			}

			if tt.skill != nil {
				err := tt.skill.Save()
				if tt.wantErr {
					assert.Error(t, err)
					return
				}

				assert.NoError(t, err)

				// Verify the update
				saved := &models.Skill{}
				err = saved.FirstId(tt.skill.ID)
				assert.NoError(t, err)
				assert.Equal(t, "Updated Description", saved.Beschreibung)
				assert.NotZero(t, saved.GameSystemId)
			}
		})
	}
	database.ResetTestDB()
}

func TestSkill_Delete(t *testing.T) {
	database.SetupTestDB()
	testDefinition := []struct {
		name    string
		skill   *models.Skill
		wantErr bool
	}{
		{
			name: "delete existing skill",
			skill: &models.Skill{
				Name:         "Test Skill",
				GameSystemId: 1,
			},
			wantErr: false,
		},
		{
			name: "delete non-existing skill",
			skill: &models.Skill{
				ID: 9999,
			},
			wantErr: true,
		},
		/*{
			name:    "delete nil skill",
			skill:   nil,
			wantErr: true,
		},*/
	}

	for _, tt := range testDefinition {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr && tt.skill != nil {
				// Create test data first
				err := tt.skill.Create()
				assert.NoError(t, err)
			}

			var err error
			if tt.skill != nil {
				err = tt.skill.Delete()
			}

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Verify deletion
			s := &models.Skill{}
			err = s.FirstId(tt.skill.ID)
			assert.Error(t, err) // Should error since record is deleted
		})
	}
	database.ResetTestDB()
}

func TestSkill_GetSkillCategories(t *testing.T) {
	database.SetupTestDB()
	// Create test skill categories in the lookup table
	// Note: GetSkillCategories() reads from gsm_skill_categories table, not from skills
	testCategories := []*models.SkillCategory{
		{
			Name:         "Category1",
			GameSystemId: 1,
		},
		{
			Name:         "Category2",
			GameSystemId: 1,
		},
	}

	// Create test categories in the lookup table
	for _, cat := range testCategories {
		err := cat.Create()
		assert.NoError(t, err)
	}

	testDefinition := []struct {
		name                string
		expectedMinCount    int
		expectedMustContain []string
		wantErr             bool
	}{
		{
			name:                "get categories",
			expectedMinCount:    2,
			expectedMustContain: []string{"Category1", "Category2"},
			wantErr:             false,
		},
	}

	for _, tt := range testDefinition {
		t.Run(tt.name, func(t *testing.T) {
			s := &models.Skill{}
			categories, err := s.GetSkillCategories()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.GreaterOrEqual(t, len(categories), tt.expectedMinCount, "Should have at least %d categories", tt.expectedMinCount)

			// Check that all expected categories are present
			for _, expected := range tt.expectedMustContain {
				assert.Contains(t, categories, expected, "Categories should contain %q", expected)
			}
		})
	}
	database.ResetTestDB()
}
