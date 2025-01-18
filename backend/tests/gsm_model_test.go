package tests

import (
	"bamort/gsmaster"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkill_Create(t *testing.T) {
	SetupTestDB()
	tests := []struct {
		name    string
		skill   *gsmaster.Skill
		wantErr bool
	}{
		{
			name: "valid skill",
			skill: &gsmaster.Skill{
				LookupList: gsmaster.LookupList{
					Name:         "Test Skill",
					Beschreibung: "Test Description",
					Quelle:       "Test Source",
				},
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.skill.Create()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, "midgard", tt.skill.GameSystem)
			assert.NotZero(t, tt.skill.ID)
		})
	}
	resetDB()
}

func TestWeaponSkill_Create(t *testing.T) {
	SetupTestDB()
	tests := []struct {
		name        string
		weaponSkill *gsmaster.WeaponSkill
		wantErr     bool
	}{
		{
			name: "valid weapon skill",
			weaponSkill: &gsmaster.WeaponSkill{
				Skill: gsmaster.Skill{
					LookupList: gsmaster.LookupList{
						Name:         "Test Weapon Skill",
						Beschreibung: "Test Description",
						Quelle:       "Test Source",
					},
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.weaponSkill.Create()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, "midgard", tt.weaponSkill.GameSystem)
			assert.NotZero(t, tt.weaponSkill.ID)
		})
	}
	resetDB()
}

func TestSkill_TableName(t *testing.T) {
	SetupTestDB()
	skill := &gsmaster.Skill{}
	expected := "gsm_skills"
	assert.Equal(t, expected, skill.TableName())
	resetDB()
}

func TestSkill_First(t *testing.T) {
	SetupTestDB()
	tests := []struct {
		name     string
		skill    *gsmaster.Skill
		findName string
		wantErr  bool
	}{
		{
			name: "existing skill",
			skill: &gsmaster.Skill{
				LookupList: gsmaster.LookupList{
					Name:       "Test Skill",
					GameSystem: "midgard",
				},
			},
			findName: "Test Skill",
			wantErr:  false,
		},
		{
			name:     "non-existing skill",
			skill:    &gsmaster.Skill{},
			findName: "NonExistent",
			wantErr:  true,
		},
		{
			name:     "empty name",
			skill:    &gsmaster.Skill{},
			findName: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				// Create test data first
				err := tt.skill.Create()
				assert.NoError(t, err)
			}

			s := &gsmaster.Skill{}
			err := s.First(tt.findName)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.skill.Name, s.Name)
			assert.Equal(t, "midgard", s.GameSystem)
		})
	}
	resetDB()
}

func TestSkill_FirstId(t *testing.T) {
	SetupTestDB()
	tests := []struct {
		name    string
		skill   *gsmaster.Skill
		findId  uint
		wantErr bool
	}{
		{
			name: "existing skill",
			skill: &gsmaster.Skill{
				LookupList: gsmaster.LookupList{
					Name:       "Test Skill",
					GameSystem: "midgard",
				},
			},
			findId:  1,
			wantErr: false,
		},
		{
			name:    "non-existing id",
			skill:   &gsmaster.Skill{},
			findId:  9999,
			wantErr: true,
		},
		{
			name:    "zero id",
			skill:   &gsmaster.Skill{},
			findId:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				// Create test data first
				err := tt.skill.Create()
				assert.NoError(t, err)
				tt.findId = tt.skill.ID
			}

			s := &gsmaster.Skill{}
			err := s.FirstId(tt.findId)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.skill.Name, s.Name)
			assert.Equal(t, "midgard", s.GameSystem)
			assert.Equal(t, tt.findId, s.ID)
		})
	}
	resetDB()
}

func TestSkill_Save(t *testing.T) {
	SetupTestDB()
	tests := []struct {
		name    string
		skill   *gsmaster.Skill
		wantErr bool
	}{
		{
			name: "update existing skill",
			skill: &gsmaster.Skill{
				LookupList: gsmaster.LookupList{
					Name:         "Test Skill",
					Beschreibung: "Original Description",
					GameSystem:   "midgard",
				},
			},
			wantErr: false,
		},
		/*{
			name:    "nil skill",
			skill:   nil,
			wantErr: true,
		},*/
	}

	for _, tt := range tests {
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
				saved := &gsmaster.Skill{}
				err = saved.FirstId(tt.skill.ID)
				assert.NoError(t, err)
				assert.Equal(t, "Updated Description", saved.Beschreibung)
			}
		})
	}
	resetDB()
}

func TestSkill_Delete(t *testing.T) {
	SetupTestDB()
	tests := []struct {
		name    string
		skill   *gsmaster.Skill
		wantErr bool
	}{
		{
			name: "delete existing skill",
			skill: &gsmaster.Skill{
				LookupList: gsmaster.LookupList{
					Name:       "Test Skill",
					GameSystem: "midgard",
				},
			},
			wantErr: false,
		},
		{
			name: "delete non-existing skill",
			skill: &gsmaster.Skill{
				LookupList: gsmaster.LookupList{
					ID: 9999,
				},
			},
			wantErr: true,
		},
		/*{
			name:    "delete nil skill",
			skill:   nil,
			wantErr: true,
		},*/
	}

	for _, tt := range tests {
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
			s := &gsmaster.Skill{}
			err = s.FirstId(tt.skill.ID)
			assert.Error(t, err) // Should error since record is deleted
		})
	}
	resetDB()
}

func TestSkill_GetSkillCategories(t *testing.T) {
	SetupTestDB()
	// Create test skills with different categories
	testSkills := []*gsmaster.Skill{
		{
			LookupList: gsmaster.LookupList{
				Name:       "Skill1",
				GameSystem: "midgard",
			},
			Category: "Category1",
		},
		{
			LookupList: gsmaster.LookupList{
				Name:       "Skill2",
				GameSystem: "midgard",
			},
			Category: "Category2",
		},
		{
			LookupList: gsmaster.LookupList{
				Name:       "Skill3",
				GameSystem: "midgard",
			},
			Category: "Category1", // Duplicate category
		},
		{
			LookupList: gsmaster.LookupList{
				Name:       "Skill4",
				GameSystem: "midgard",
			},
			Category: "", // Empty category
		},
	}

	// Create test data
	for _, skill := range testSkills {
		err := skill.Create()
		assert.NoError(t, err)
	}

	tests := []struct {
		name          string
		expectedCount int
		expectedFound []string
		wantErr       bool
	}{
		{
			name:          "get categories",
			expectedCount: 3,
			expectedFound: []string{"Category1", "Category2", ""},
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &gsmaster.Skill{}
			categories, err := s.GetSkillCategories()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, categories, tt.expectedCount)
			assert.ElementsMatch(t, tt.expectedFound, categories)
		})
	}
	resetDB()
}
