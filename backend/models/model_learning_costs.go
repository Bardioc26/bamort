package models

import (
	"bamort/database"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Source repräsentiert ein Regelwerk/Quellenbuch
type Source struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Code         string `gorm:"uniqueIndex;size:10;not null" json:"code"` // z.B. "KOD", "ARK", "MYS"
	Name         string `gorm:"unique;not null" json:"name"`              // z.B. "Kodex", "Arkanum", "Mysterium"
	FullName     string `json:"full_name,omitempty"`                      // z.B. "Midgard Regelwerk - Kodex"
	Edition      string `json:"edition,omitempty"`                        // z.B. "5. Edition"
	Publisher    string `json:"publisher,omitempty"`                      // z.B. "Pegasus Spiele"
	PublishYear  int    `json:"publish_year,omitempty"`                   // Erscheinungsjahr
	Description  string `json:"description,omitempty"`                    // Beschreibung des Werks
	IsCore       bool   `json:"is_core"`                                  // Ist es ein Grundregelwerk?
	IsActive     bool   `json:"is_active"`                                // Ist das Werk aktiv/verfügbar?
	GameSystem   string `gorm:"index;default:midgard" json:"game_system"`
	GameSystemId uint   `json:"game_system_id,omitempty"`
}

// CharacterClass repräsentiert eine Charakterklasse mit Code und vollständigem Namen
type CharacterClass struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Code         string `gorm:"uniqueIndex;size:3" json:"code"`  // z.B. "Hx", "Ma", "Kr"
	Name         string `gorm:"unique;not null" json:"name"`     // z.B. "Hexer", "Magier", "Krieger"
	Description  string `json:"description,omitempty"`           // Optional: Beschreibung der Klasse
	SourceID     uint   `gorm:"not null;index" json:"source_id"` // Verweis auf das Quellenbuch
	GameSystem   string `gorm:"index;default:midgard" json:"game_system"`
	GameSystemId uint   `json:"game_system_id,omitempty"`
	Source       Source `gorm:"foreignKey:SourceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"source"`
	Quelle       string `gorm:"index;size:3;default:KOD" json:"quelle,omitempty"` // Optional: Quelle der Kategorie, falls abweichend
}

// SkillCategory repräsentiert eine Fertigkeitskategorie
type SkillCategory struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"unique;not null" json:"name"`     // z.B. "Alltag", "Kampf", "Waffen"
	SourceID     uint   `gorm:"not null;index" json:"source_id"` // Verweis auf das Quellenbuch
	GameSystem   string `gorm:"index;default:midgard" json:"game_system"`
	GameSystemId uint   `json:"game_system_id,omitempty"`
	Source       Source `gorm:"foreignKey:SourceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"source"`
	Quelle       string `gorm:"index;size:3;default:KOD" json:"quelle,omitempty"` // Optional: Quelle der Kategorie, falls abweichend
}

// SkillDifficulty repräsentiert Schwierigkeitsgrade
type SkillDifficulty struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"unique;not null" json:"name"` // z.B. "leicht", "normal", "schwer", "sehr schwer"
	Description  string `json:"description,omitempty"`       // Optional: Beschreibung
	GameSystem   string `gorm:"index;default:midgard" json:"game_system"`
	GameSystemId uint   `json:"game_system_id,omitempty"`
}

// SpellSchool repräsentiert Zauberschulen
type SpellSchool struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"unique;not null" json:"name"`     // z.B. "Beherrschen", "Bewegen", "Erkennen"
	Description  string `json:"description,omitempty"`           // Optional: Beschreibung
	SourceID     uint   `gorm:"not null;index" json:"source_id"` // Verweis auf das Quellenbuch
	GameSystem   string `gorm:"index;default:midgard" json:"game_system"`
	GameSystemId uint   `json:"game_system_id,omitempty"`
	Source       Source `gorm:"foreignKey:SourceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"source"`
	Quelle       string `gorm:"index;size:3;default:KOD" json:"quelle,omitempty"` // Optional: Quelle der Kategorie, falls abweichend
}

// ClassCategoryEPCost definiert EP-Kosten für 1 Trainingseinheit (TE) pro Charakterklasse und Fertigkeitskategorie
type ClassCategoryEPCost struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	CharacterClassID uint           `gorm:"not null;index" json:"character_class_id"`
	SkillCategoryID  uint           `gorm:"not null;index" json:"skill_category_id"`
	EPPerTE          int            `gorm:"not null" json:"ep_per_te"` // EP-Kosten für 1 Trainingseinheit
	CharacterClass   CharacterClass `gorm:"foreignKey:CharacterClassID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"character_class"`
	SkillCategory    SkillCategory  `gorm:"foreignKey:SkillCategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"skill_category"`
	CCLass           string         `gorm:"column:character_class;size:3;not null;index;default:Or" json:"characterClass"`    // Abkürzung der Charakterklasse
	SCategory        string         `gorm:"column:skill_category;size:15;not null;index;default:Alltag" json:"skillCategory"` // SkillCategory
}

// ClassSpellSchoolEPCost definiert EP-Kosten für 1 Lerneinheit (LE) für Zauber pro Charakterklasse und Zauberschule
type ClassSpellSchoolEPCost struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	CharacterClassID uint           `gorm:"not null;index" json:"character_class_id"`
	SpellSchoolID    uint           `gorm:"not null;index" json:"spell_school_id"`
	EPPerLE          int            `gorm:"not null" json:"ep_per_le"` // EP-Kosten für 1 Lerneinheit
	CharacterClass   CharacterClass `gorm:"foreignKey:CharacterClassID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"character_class"`
	SpellSchool      SpellSchool    `gorm:"foreignKey:SpellSchoolID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"spell_school"`
	CCLass           string         `gorm:"column:character_class;size:3;not null;index;default:Or" json:"characterClass"` // Abkürzung der Charakterklasse
	SCategory        string         `gorm:"column:spell_school;size:15;not null;index;default:Dweomer" json:"spellSchool"` // Spell School Name
}

// SpellLevelLECost definiert LE-Kosten pro Zauber-Stufe
type SpellLevelLECost struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Level        int    `gorm:"uniqueIndex;not null" json:"level"` // Zauber-Stufe (1-12)
	LERequired   int    `gorm:"not null" json:"le_required"`       // Benötigte Lerneinheiten
	GameSystem   string `gorm:"index;default:midgard" json:"game_system"`
	GameSystemId uint   `json:"game_system_id,omitempty"`
}

// SkillCategoryDifficulty definiert die Schwierigkeit einer Fertigkeit in einer bestimmten Kategorie
// Eine Fertigkeit kann in mehreren Kategorien mit unterschiedlichen Schwierigkeiten existieren
type SkillCategoryDifficulty struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	SkillID           uint            `gorm:"not null;index" json:"skill_id"`
	SkillCategoryID   uint            `gorm:"not null;index" json:"skill_category_id"`
	SkillDifficultyID uint            `gorm:"not null;index" json:"skill_difficulty_id"`
	LearnCost         int             `gorm:"not null" json:"learn_cost"` // LE-Kosten für das Erlernen
	Skill             Skill           `gorm:"foreignKey:SkillID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"skill"`
	SkillCategory     SkillCategory   `gorm:"foreignKey:SkillCategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"skill_category"`
	SkillDifficulty   SkillDifficulty `gorm:"foreignKey:SkillDifficultyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"skill_difficulty"`
	SDifficulty       string          `gorm:"column:skill_difficulty;size:15;not null;index;default:normal" json:"skillDifficulty"` // Abkürzung der Charakterklasse
	SCategory         string          `gorm:"column:skill_category;size:25;not null;index;default:Alltag" json:"skillCategory"`     // SkillCategory
}

// WeaponSkillCategoryDifficulty definiert die Schwierigkeit einer Waffenfertigkeit in einer bestimmten Kategorie
// Separate Tabelle für Waffenfertigkeiten, da diese in gsm_weaponskills gespeichert sind
type WeaponSkillCategoryDifficulty struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	WeaponSkillID     uint            `gorm:"not null;index" json:"weapon_skill_id"`
	SkillCategoryID   uint            `gorm:"not null;index" json:"skill_category_id"`
	SkillDifficultyID uint            `gorm:"not null;index" json:"skill_difficulty_id"`
	LearnCost         int             `gorm:"not null" json:"learn_cost"` // LE-Kosten für das Erlernen
	WeaponSkill       WeaponSkill     `gorm:"foreignKey:WeaponSkillID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"weapon_skill"`
	SkillCategory     SkillCategory   `gorm:"foreignKey:SkillCategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"skill_category"`
	SkillDifficulty   SkillDifficulty `gorm:"foreignKey:SkillDifficultyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"skill_difficulty"`
	SDifficulty       string          `gorm:"column:skill_difficulty;size:15;not null;index;default:normal" json:"skillDifficulty"`
	SCategory         string          `gorm:"column:skill_category;size:25;not null;index;default:Waffen" json:"skillCategory"`
}

// SkillImprovementCost definiert TE-Kosten für Verbesserungen basierend auf Kategorie, Schwierigkeit und aktuellem Wert
type SkillImprovementCost struct {
	ID                        uint                    `gorm:"primaryKey" json:"id"`
	SkillCategoryDifficultyID uint                    `gorm:"not null;index" json:"skill_category_difficulty_id"`
	CurrentLevel              int                     `gorm:"not null;index" json:"current_level"` // Aktueller Fertigkeitswert
	TERequired                int                     `gorm:"not null" json:"te_required"`         // Benötigte Trainingseinheiten
	SkillCategoryDifficulty   SkillCategoryDifficulty `gorm:"foreignKey:SkillCategoryDifficultyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"skill_category_difficulty"`
}

// Für komplexere Queries: View oder Helper-Strukturen

// SkillLearningInfo kombiniert alle Informationen für eine Fertigkeit in einer bestimmten Kategorie
type SkillLearningInfo struct {
	SkillID          uint   `json:"skill_id"`
	SkillName        string `json:"skill_name"`
	CategoryID       uint   `json:"category_id"`
	CategoryName     string `json:"category_name"`
	DifficultyID     uint   `json:"difficulty_id"`
	DifficultyName   string `json:"difficulty_name"`
	LearnCost        int    `json:"learn_cost"` // LE-Kosten für das Erlernen
	CharacterClassID uint   `json:"character_class_id"`
	ClassCode        string `json:"class_code"`
	ClassName        string `json:"class_name"`
	EPPerTE          int    `json:"ep_per_te"`
	GameSystem       string `json:"game_system"`
	GameSystemId     uint   `json:"game_system_id,omitempty"`
}

// SpellLearningInfo kombiniert alle Informationen für einen Zauber
type SpellLearningInfo struct {
	SpellID          uint   `json:"spell_id"`
	SpellName        string `json:"spell_name"`
	SpellLevel       int    `json:"spell_level"`
	SchoolID         uint   `json:"school_id"`
	SchoolName       string `json:"school_name"`
	CharacterClassID uint   `json:"character_class_id"`
	ClassCode        string `json:"class_code"`
	ClassName        string `json:"class_name"`
	EPPerLE          int    `json:"ep_per_le"`
	LERequired       int    `json:"le_required"`
	GameSystem       string `json:"game_system"`
	GameSystemId     uint   `json:"game_system_id,omitempty"`
}

// AuditLogEntry repräsentiert einen Eintrag im Audit-Log
type AuditLogEntry struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CharacterID uint      `json:"character_id" gorm:"not null"`
	FieldName   string    `json:"field_name" gorm:"not null"` // "experience_points", "gold", "silver", "copper"
	OldValue    int       `json:"old_value"`
	NewValue    int       `json:"new_value"`
	Difference  int       `json:"difference"` // NewValue - OldValue
	Reason      string    `json:"reason"`     // "manual", "skill_learning", "skill_improvement", "spell_learning", etc.
	UserID      uint      `json:"user_id,omitempty"`
	Timestamp   time.Time `json:"timestamp" gorm:"autoCreateTime"`
	Notes       string    `json:"notes,omitempty"` // Zusätzliche Informationen
}

// TableName-Methoden für GORM
func (Source) TableName() string {
	//learning_sources
	return "gsm_lit_sources"
}

func (CharacterClass) TableName() string {
	//learning_character_classes
	return "gsm_character_classes"
}

func (SkillCategory) TableName() string {
	return "learning_skill_categories"
}

func (SkillDifficulty) TableName() string {
	return "learning_skill_difficulties"
}

func (SpellSchool) TableName() string {
	return "learning_spell_schools"
}

func (ClassCategoryEPCost) TableName() string {
	return "learning_class_category_ep_costs"
}

func (ClassSpellSchoolEPCost) TableName() string {
	return "learning_class_spell_school_ep_costs"
}

func (SpellLevelLECost) TableName() string {
	return "learning_spell_level_le_costs"
}

func (SkillCategoryDifficulty) TableName() string {
	return "learning_skill_category_difficulties"
}

func (WeaponSkillCategoryDifficulty) TableName() string {
	return "learning_weaponskill_category_difficulties"
}

func (SkillImprovementCost) TableName() string {
	return "learning_skill_improvement_costs"
}

func (sllc *SpellLevelLECost) ensureGameSystem() {
	gs := GetGameSystem(sllc.GameSystemId, sllc.GameSystem)
	sllc.GameSystemId = gs.ID
	sllc.GameSystem = gs.Name
}

func (sllc *SpellLevelLECost) BeforeCreate(tx *gorm.DB) error {
	sllc.ensureGameSystem()
	return nil
}

func (sllc *SpellLevelLECost) BeforeSave(tx *gorm.DB) error {
	sllc.ensureGameSystem()
	return nil
}

func (sllc *SpellLevelLECost) Save() error {
	sllc.ensureGameSystem()
	return database.DB.Save(sllc).Error
}

func (s *Source) ensureGameSystem() {
	gs := GetGameSystem(s.GameSystemId, s.GameSystem)
	s.GameSystemId = gs.ID
	s.GameSystem = gs.Name
}

func (s *Source) BeforeCreate(tx *gorm.DB) error {
	s.ensureGameSystem()
	return nil
}

func (s *Source) BeforeSave(tx *gorm.DB) error {
	s.ensureGameSystem()
	return nil
}

func (s *Source) Save() error {
	s.ensureGameSystem()
	return database.DB.Save(s).Error
}

func (cc *CharacterClass) ensureGameSystem() {
	gs := GetGameSystem(cc.GameSystemId, cc.GameSystem)
	cc.GameSystemId = gs.ID
	cc.GameSystem = gs.Name
}

func (cc *CharacterClass) BeforeCreate(tx *gorm.DB) error {
	cc.ensureGameSystem()
	return nil
}

func (cc *CharacterClass) BeforeSave(tx *gorm.DB) error {
	cc.ensureGameSystem()
	return nil
}

func (cc *CharacterClass) Save() error {
	cc.ensureGameSystem()
	return database.DB.Save(cc).Error
}

func (sc *SkillCategory) ensureGameSystem() {
	gs := GetGameSystem(sc.GameSystemId, sc.GameSystem)
	sc.GameSystemId = gs.ID
	sc.GameSystem = gs.Name
}

func (sc *SkillCategory) BeforeCreate(tx *gorm.DB) error {
	sc.ensureGameSystem()
	return nil
}

func (sc *SkillCategory) BeforeSave(tx *gorm.DB) error {
	sc.ensureGameSystem()
	return nil
}

func (sc *SkillCategory) Save() error {
	sc.ensureGameSystem()
	return database.DB.Save(sc).Error
}

func (sd *SkillDifficulty) ensureGameSystem() {
	gs := GetGameSystem(sd.GameSystemId, sd.GameSystem)
	sd.GameSystemId = gs.ID
	sd.GameSystem = gs.Name
}

func (sd *SkillDifficulty) BeforeCreate(tx *gorm.DB) error {
	sd.ensureGameSystem()
	return nil
}

func (sd *SkillDifficulty) BeforeSave(tx *gorm.DB) error {
	sd.ensureGameSystem()
	return nil
}

func (sd *SkillDifficulty) Save() error {
	sd.ensureGameSystem()
	return database.DB.Save(sd).Error
}

func (ss *SpellSchool) ensureGameSystem() {
	gs := GetGameSystem(ss.GameSystemId, ss.GameSystem)
	ss.GameSystemId = gs.ID
	ss.GameSystem = gs.Name
}

func (ss *SpellSchool) BeforeCreate(tx *gorm.DB) error {
	ss.ensureGameSystem()
	return nil
}

func (ss *SpellSchool) BeforeSave(tx *gorm.DB) error {
	ss.ensureGameSystem()
	return nil
}

func (ss *SpellSchool) Save() error {
	ss.ensureGameSystem()
	return database.DB.Save(ss).Error
}

// CRUD-Methoden
func (object *Source) CreatewGS(gamesystemId uint) error {
	gs := GetGameSystem(gamesystemId, "")
	if gs == nil {
		return fmt.Errorf("invalid game system")
	}
	object.GameSystemId = gs.ID
	return database.DB.Create(object).Error
}

func (object *Source) Create() error {
	object.ensureGameSystem()
	return object.CreatewGS(object.GameSystemId)
}

func (object *Source) FirstByCodewGS(code string, gameSystemId uint) error {
	if code == "" {
		return fmt.Errorf("name cannot be empty")
	}
	gs := GetGameSystem(gameSystemId, "")
	if gs == nil {
		return fmt.Errorf("invalid game system")
	}
	return database.DB.Where("(game_system = ? OR game_system_id = ?) AND code = ?", gs.Name, gs.ID, code).First(object).Error
}

func (object *Source) FirstByCode(code string) error {
	gs := GetGameSystem(object.GameSystemId, object.GameSystem)
	return object.FirstByCodewGS(code, gs.ID)
}

func (object *Source) FirstByNamewGS(name string, gameSystemId uint) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	gs := GetGameSystem(gameSystemId, "")
	return database.DB.Where("(game_system = ? OR game_system_id = ?) AND name = ?", gs.Name, gs.ID, name).First(object).Error
}

func (object *Source) FirstByName(name string) error {
	gs := GetGameSystem(object.GameSystemId, object.GameSystem)
	return object.FirstByNamewGS(name, gs.ID)
}

func (cc *CharacterClass) Create() error {
	cc.ensureGameSystem()
	return database.DB.Create(cc).Error
}

func (cc *CharacterClass) FirstByCode(code string) error {
	gs := GetGameSystem(cc.GameSystemId, cc.GameSystem)
	return database.DB.Where("(game_system = ? OR game_system_id = ?) AND code = ?", gs.Name, gs.ID, code).First(cc).Error
}
func (cc *CharacterClass) FirstByName(code string) error {
	gs := GetGameSystem(cc.GameSystemId, cc.GameSystem)
	return database.DB.Where("(game_system = ? OR game_system_id = ?) AND name = ?", gs.Name, gs.ID, code).First(cc).Error
}
func (cc *CharacterClass) FirstByNameOrCode(value string) error {
	gs := GetGameSystem(cc.GameSystemId, cc.GameSystem)
	return database.DB.Where("(game_system = ? OR game_system_id = ?) AND (name = ? OR code = ?)", gs.Name, gs.ID, value, value).First(cc).Error
}

func (sc *SkillCategory) Create() error {
	sc.ensureGameSystem()
	return database.DB.Create(sc).Error
}

func (sc *SkillCategory) FirstByName(name string) error {
	gs := GetGameSystem(sc.GameSystemId, sc.GameSystem)
	return database.DB.Where("(game_system = ? OR game_system_id = ?) AND name = ?", gs.Name, gs.ID, name).First(sc).Error
}

func (sd *SkillDifficulty) Create() error {
	sd.ensureGameSystem()
	return database.DB.Create(sd).Error
}

func (sd *SkillDifficulty) FirstByName(name string) error {
	gs := GetGameSystem(sd.GameSystemId, sd.GameSystem)
	return database.DB.Where("(game_system = ? OR game_system_id = ?) AND name = ?", gs.Name, gs.ID, name).First(sd).Error
}

func (ss *SpellSchool) Create() error {
	ss.ensureGameSystem()
	return database.DB.Create(ss).Error
}

func (ss *SpellSchool) FirstByName(name string) error {
	gs := GetGameSystem(ss.GameSystemId, ss.GameSystem)
	return database.DB.Where("(game_system = ? OR game_system_id = ?) AND name = ?", gs.Name, gs.ID, name).First(ss).Error
}

func (ccec *ClassCategoryEPCost) Create() error {
	return database.DB.Create(ccec).Error
}

func (cssec *ClassSpellSchoolEPCost) Create() error {
	return database.DB.Create(cssec).Error
}

func (sllc *SpellLevelLECost) Create() error {
	sllc.ensureGameSystem()
	return database.DB.Create(sllc).Error
}

func (scd *SkillCategoryDifficulty) Create() error {
	return database.DB.Create(scd).Error
}

func (wscd *WeaponSkillCategoryDifficulty) Create() error {
	return database.DB.Create(wscd).Error
}

func (sic *SkillImprovementCost) Create() error {
	return database.DB.Create(sic).Error
}

// Komplexere Query-Methoden

// GetEPPerTEForClassAndCategory holt die EP-Kosten für eine Klasse und Kategorie
func GetEPPerTEForClassAndCategory(classCode string, categoryName string) (int, error) {
	var result ClassCategoryEPCost
	err := database.DB.
		Where("character_class = ? AND skill_category = ?", classCode, categoryName).
		First(&result).Error

	if err != nil {
		return 0, err
	}
	return result.EPPerTE, nil
}

// GetEPPerLEForClassAndSpellSchool holt die EP-Kosten für eine Klasse und Zauberschule
func GetEPPerLEForClassAndSpellSchool(classCode string, schoolName string) (int, error) {
	var result ClassSpellSchoolEPCost
	err := database.DB.
		Where("character_class = ? AND spell_school = ?", classCode, schoolName).
		First(&result).Error

	if err != nil {
		return 0, err
	}
	return result.EPPerLE, nil
}

// GetSkillCategoryAndDifficultyNewSystem findet die beste Kategorie für eine Fertigkeit basierend auf niedrigsten EP-Kosten
func GetSkillCategoryAndDifficultyNewSystem(skillName string, classCode string) (*SkillLearningInfo, error) {
	var results []SkillLearningInfo
	gs := GetGameSystem(0, "midgard")

	err := database.DB.Raw(`
		SELECT 
			scd.skill_id,
			s.name as skill_name,
			s.game_system as game_system,
			s.game_system_id as game_system_id,
			scd.skill_category as category_name,
			scd.skill_difficulty as difficulty_name,
			scd.learn_cost,
			ccec.character_class as class_code,
			ccec.character_class as class_name,
			ccec.ep_per_te,
			(scd.learn_cost * ccec.ep_per_te) as total_cost
		FROM learning_skill_category_difficulties scd
		JOIN learning_class_category_ep_costs ccec ON scd.skill_category = ccec.skill_category
		JOIN gsm_skills s ON scd.skill_id = s.id
		WHERE s.name = ? AND ccec.character_class = ? AND (s.game_system = ? OR s.game_system_id = ?)
		ORDER BY total_cost ASC
	`, skillName, classCode, gs.Name, gs.ID).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	for i := range results {
		gs := GetGameSystem(results[i].GameSystemId, results[i].GameSystem)
		results[i].GameSystemId = gs.ID
		results[i].GameSystem = gs.Name
	}

	if len(results) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &results[0], nil
}

// GetSkillInfoCategoryAndDifficultyNewSystem holt die Informationen für eine spezifische Kategorie/Schwierigkeit
func GetSkillInfoCategoryAndDifficultyNewSystem(skillName, category, difficulty, classCode string) (*SkillLearningInfo, error) {
	var result SkillLearningInfo
	gs := GetGameSystem(0, "midgard")

	err := database.DB.Raw(`
		SELECT 
			scd.skill_id,
			s.name as skill_name,
			s.game_system as game_system,
			s.game_system_id as game_system_id,
			scd.skill_category as category_name,
			scd.skill_difficulty as difficulty_name,
			scd.learn_cost,
			ccec.character_class as class_code,
			ccec.character_class as class_name,
			ccec.ep_per_te,
			(scd.learn_cost * ccec.ep_per_te) as total_cost
		FROM learning_skill_category_difficulties scd
		JOIN learning_class_category_ep_costs ccec ON scd.skill_category = ccec.skill_category
		JOIN gsm_skills s ON scd.skill_id = s.id
		WHERE s.name = ? AND scd.skill_category = ? AND scd.skill_difficulty = ? AND ccec.character_class = ? AND (s.game_system = ? OR s.game_system_id = ?)
	`, skillName, category, difficulty, classCode, gs.Name, gs.ID).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	gs = GetGameSystem(result.GameSystemId, result.GameSystem)
	result.GameSystemId = gs.ID
	result.GameSystem = gs.Name

	return &result, nil
}

// GetSpellLearningInfoNewSystem holt alle Informationen für das Erlernen eines Zaubers
func GetSpellLearningInfoNewSystem(spellName string, classCode string) (*SpellLearningInfo, error) {
	var result SpellLearningInfo

	err := database.DB.Raw(`
		SELECT 
			s.id as spell_id,
			s.name as spell_name,
			s.game_system as game_system,
			s.game_system_id as game_system_id,
			s.stufe as spell_level,
			COALESCE(NULLIF(s.category, ''), s.learning_category) as school_name,
			cssec.character_class as class_code,
			cssec.character_class as class_name,
			cssec.ep_per_le,
			COALESCE(sllc.le_required, 0) as le_required
		FROM gsm_spells s
		JOIN learning_class_spell_school_ep_costs cssec ON COALESCE(NULLIF(s.category, ''), s.learning_category) = cssec.spell_school
		LEFT JOIN learning_spell_level_le_costs sllc ON s.stufe = sllc.level AND (sllc.game_system = s.game_system OR sllc.game_system_id = s.game_system_id OR sllc.game_system_id IS NULL)
		WHERE s.name = ? AND cssec.character_class = ?
	`, spellName, classCode).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	// Validate that we found a spell (spell_id should be > 0)
	if result.SpellID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// Validate spell level - level 0 is not a valid spell level
	if result.SpellLevel <= 0 {
		return nil, fmt.Errorf("invalid spell level %d for spell '%s' - spell levels must be 1 or higher", result.SpellLevel, spellName)
	}

	gs := GetGameSystem(result.GameSystemId, result.GameSystem)
	result.GameSystemId = gs.ID
	result.GameSystem = gs.Name

	return &result, nil
}

// GetImprovementCost holt die Verbesserungskosten für eine Fertigkeit
func GetImprovementCost(skillName string, categoryName string, difficultyName string, currentLevel int) (int, error) {
	var result SkillImprovementCost

	err := database.DB.Raw(`
        SELECT sic.te_required
        FROM learning_skill_improvement_costs sic
        JOIN learning_skill_category_difficulties scd ON sic.skill_category_difficulty_id = scd.id
        JOIN gsm_skills s ON scd.skill_id = s.id
        WHERE s.name = ? AND scd.skill_category = ? AND scd.skill_difficulty = ? AND sic.current_level = ?
    `, skillName, categoryName, difficultyName, currentLevel).Scan(&result).Error

	if err != nil {
		return 0, err
	}

	return result.TERequired, nil
}

// GetSkillLearnCost holt die Lernkosten (LE) für eine Fertigkeit basierend auf Kategorie und Schwierigkeit
func GetSkillLearnCost(categoryName string, difficultyName string) (int, error) {
	var result SkillCategoryDifficulty
	err := database.DB.
		Where("skill_category = ? AND skill_difficulty = ?", categoryName, difficultyName).
		First(&result).Error

	if err != nil {
		return 0, err
	}
	return result.LearnCost, nil
}

// Quellenmanagement-Funktionen

// GetActiveSourceCodes gibt alle aktiven Quellencodes zurück
func GetActiveSourceCodes() ([]string, error) {
	var codes []string
	err := database.DB.Model(&Source{}).Where("is_active = ?", true).Pluck("code", &codes).Error
	return codes, err
}

// GetSourcesByGameSystem gibt alle Quellen für ein bestimmtes Spielsystem zurück
func GetSourcesByGameSystem(gameSystem string) ([]Source, error) {
	var sources []Source
	err := database.DB.Where("game_system = ?", gameSystem).Order("is_core DESC, code ASC").Find(&sources).Error
	return sources, err
}

// GetSkillsByActiveSources gibt alle Fertigkeiten zurück, die in aktiven Quellen definiert sind
func GetSkillsByActiveSources(gameSystem string) ([]Skill, error) {
	var skills []Skill
	err := database.DB.
		Joins("LEFT JOIN gsm_lit_sources ON gsm_skills.source_id = gsm_lit_sources.id").
		Where("gsm_skills.game_system = ? AND (gsm_lit_sourcesis_active = ? OR gsm_skills.source_id IS NULL)", gameSystem, true).
		Find(&skills).Error
	return skills, err
}

// GetSpellsByActiveSources gibt alle Zauber zurück, die in aktiven Quellen definiert sind
func GetSpellsByActiveSources(gameSystem string) ([]Spell, error) {
	var spells []Spell
	err := database.DB.
		Joins("LEFT JOIN gsm_lit_sources ON gsm_spells.source_id = gsm_lit_sources.id").
		Where("gsm_spells.game_system = ? AND (gsm_lit_sources.is_active = ? OR gsm_spells.source_id IS NULL)", gameSystem, true).
		Find(&spells).Error
	return spells, err
}

// GetCharacterClassesByActiveSources gibt alle Charakterklassen zurück, die in aktiven Quellen definiert sind
func GetCharacterClassesByActiveSources(gameSystem string) ([]CharacterClass, error) {
	var classes []CharacterClass
	err := database.DB.
		Joins("LEFT JOIN gsm_lit_sources ON gsm_character_classes.source_id = gsm_lit_sources.id").
		Where("gsm_character_classes.game_system = ? AND (gsm_lit_sources.is_active = ? OR gsm_character_classes.source_id IS NULL)", gameSystem, true).
		Find(&classes).Error
	return classes, err
}

// SetSourceActive aktiviert oder deaktiviert eine Quelle
func SetSourceActive(sourceCode string, active bool) error {
	return database.DB.Model(&Source{}).Where("code = ?", sourceCode).Update("is_active", active).Error
}

// GetDefaultSourceForContentType gibt die Standardquelle für einen Inhaltstyp zurück
func GetDefaultSourceForContentType(contentType string) (string, error) {
	switch contentType {
	case "skill", "character_class", "skill_category":
		return "KOD", nil
	case "spell", "spell_school":
		return "ARK", nil
	case "unknown":
		return "UNB", nil
	default:
		return "KOD", nil // Fallback auf Kodex
	}
}

// GetContentTypeDefaultSources gibt eine Übersicht aller Standardquellen-Zuordnungen zurück
func GetContentTypeDefaultSources() map[string]string {
	return map[string]string{
		"character_class": "KOD", // Charakterklassen -> Kodex
		"skill":           "KOD", // Fertigkeiten -> Kodex
		"skill_category":  "KOD", // Fertigkeitskategorien -> Kodex
		"spell":           "ARK", // Zauber -> Arkanum
		"spell_school":    "ARK", // Zauberschulen -> Arkanum (erweiterte), KOD (Basis)
		"unknown":         "UNB", // Unbekannte Inhalte -> Unbekannt
	}
}
