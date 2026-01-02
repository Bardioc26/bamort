package models

import (
	"bamort/database"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Source repräsentiert ein Regelwerk/Quellenbuch
type Source struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Code        string `gorm:"uniqueIndex;size:10;not null" json:"code"` // z.B. "KOD", "ARK", "MYS"
	Name        string `gorm:"unique;not null" json:"name"`              // z.B. "Kodex", "Arkanum", "Mysterium"
	FullName    string `json:"full_name,omitempty"`                      // z.B. "Midgard Regelwerk - Kodex"
	Edition     string `json:"edition,omitempty"`                        // z.B. "5. Edition"
	Publisher   string `json:"publisher,omitempty"`                      // z.B. "Pegasus Spiele"
	PublishYear int    `json:"publish_year,omitempty"`                   // Erscheinungsjahr
	Description string `json:"description,omitempty"`                    // Beschreibung des Werks
	IsCore      bool   `gorm:"default:false" json:"is_core"`             // Ist es ein Grundregelwerk?
	IsActive    bool   `gorm:"default:true" json:"is_active"`            // Ist das Werk aktiv/verfügbar?
	GameSystem  string `gorm:"index;default:midgard" json:"game_system"`
}

// CharacterClass repräsentiert eine Charakterklasse mit Code und vollständigem Namen
type CharacterClass struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Code        string `gorm:"uniqueIndex;size:3" json:"code"`  // z.B. "Hx", "Ma", "Kr"
	Name        string `gorm:"unique;not null" json:"name"`     // z.B. "Hexer", "Magier", "Krieger"
	Description string `json:"description,omitempty"`           // Optional: Beschreibung der Klasse
	SourceID    uint   `gorm:"not null;index" json:"source_id"` // Verweis auf das Quellenbuch
	GameSystem  string `gorm:"index;default:midgard" json:"game_system"`
	Source      Source `gorm:"foreignKey:SourceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"source"`
	Quelle      string `gorm:"index;size:3;default:KOD" json:"quelle,omitempty"` // Optional: Quelle der Kategorie, falls abweichend
}

// SkillCategory repräsentiert eine Fertigkeitskategorie
type SkillCategory struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `gorm:"unique;not null" json:"name"`     // z.B. "Alltag", "Kampf", "Waffen"
	SourceID   uint   `gorm:"not null;index" json:"source_id"` // Verweis auf das Quellenbuch
	GameSystem string `gorm:"index;default:midgard" json:"game_system"`
	Source     Source `gorm:"foreignKey:SourceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"source"`
	Quelle     string `gorm:"index;size:3;default:KOD" json:"quelle,omitempty"` // Optional: Quelle der Kategorie, falls abweichend
}

// SkillDifficulty repräsentiert Schwierigkeitsgrade
type SkillDifficulty struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"unique;not null" json:"name"` // z.B. "leicht", "normal", "schwer", "sehr schwer"
	Description string `json:"description,omitempty"`       // Optional: Beschreibung
	GameSystem  string `gorm:"index;default:midgard" json:"game_system"`
}

// SpellSchool repräsentiert Zauberschulen
type SpellSchool struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"unique;not null" json:"name"`     // z.B. "Beherrschen", "Bewegen", "Erkennen"
	Description string `json:"description,omitempty"`           // Optional: Beschreibung
	SourceID    uint   `gorm:"not null;index" json:"source_id"` // Verweis auf das Quellenbuch
	GameSystem  string `gorm:"index;default:midgard" json:"game_system"`
	Source      Source `gorm:"foreignKey:SourceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"source"`
	Quelle      string `gorm:"index;size:3;default:KOD" json:"quelle,omitempty"` // Optional: Quelle der Kategorie, falls abweichend
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
	ID         uint   `gorm:"primaryKey" json:"id"`
	Level      int    `gorm:"uniqueIndex;not null" json:"level"` // Zauber-Stufe (1-12)
	LERequired int    `gorm:"not null" json:"le_required"`       // Benötigte Lerneinheiten
	GameSystem string `gorm:"index;default:midgard" json:"game_system"`
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
	LearnCost        int    `json:"learn_cost"`
	CharacterClassID uint   `json:"character_class_id"`
	ClassCode        string `json:"class_code"`
	ClassName        string `json:"class_name"`
	EPPerTE          int    `json:"ep_per_te"`
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
	return "learning_sources"
}

func (CharacterClass) TableName() string {
	return "learning_character_classes"
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

func (SkillImprovementCost) TableName() string {
	return "learning_skill_improvement_costs"
}

// CRUD-Methoden

func (s *Source) Create() error {
	return database.DB.Create(s).Error
}

func (s *Source) FirstByCode(code string) error {
	return database.DB.Where("code = ?", code).First(s).Error
}

func (s *Source) FirstByName(name string) error {
	return database.DB.Where("name = ?", name).First(s).Error
}

func (cc *CharacterClass) Create() error {
	return database.DB.Create(cc).Error
}

func (cc *CharacterClass) FirstByCode(code string) error {
	return database.DB.Where("code = ?", code).First(cc).Error
}
func (cc *CharacterClass) FirstByName(code string) error {
	return database.DB.Where("name = ?", code).First(cc).Error
}
func (cc *CharacterClass) FirstByNameOrCode(value string) error {
	return database.DB.Where("name = ? OR code = ?", value, value).First(cc).Error
}

func (sc *SkillCategory) Create() error {
	return database.DB.Create(sc).Error
}

func (sc *SkillCategory) FirstByName(name string) error {
	return database.DB.Where("name = ?", name).First(sc).Error
}

func (sd *SkillDifficulty) Create() error {
	return database.DB.Create(sd).Error
}

func (sd *SkillDifficulty) FirstByName(name string) error {
	return database.DB.Where("name = ?", name).First(sd).Error
}

func (ss *SpellSchool) Create() error {
	return database.DB.Create(ss).Error
}

func (ss *SpellSchool) FirstByName(name string) error {
	return database.DB.Where("name = ?", name).First(ss).Error
}

func (ccec *ClassCategoryEPCost) Create() error {
	return database.DB.Create(ccec).Error
}

func (cssec *ClassSpellSchoolEPCost) Create() error {
	return database.DB.Create(cssec).Error
}

func (sllc *SpellLevelLECost) Create() error {
	return database.DB.Create(sllc).Error
}

func (scd *SkillCategoryDifficulty) Create() error {
	return database.DB.Create(scd).Error
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

	err := database.DB.Raw(`
		SELECT 
			scd.skill_id,
			s.name as skill_name,
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
		WHERE s.name = ? AND ccec.character_class = ?
		ORDER BY total_cost ASC
	`, skillName, classCode).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &results[0], nil
}

// GetSkillInfoCategoryAndDifficultyNewSystem holt die Informationen für eine spezifische Kategorie/Schwierigkeit
func GetSkillInfoCategoryAndDifficultyNewSystem(skillName, category, difficulty, classCode string) (*SkillLearningInfo, error) {
	var result SkillLearningInfo

	err := database.DB.Raw(`
		SELECT 
			scd.skill_id,
			s.name as skill_name,
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
		WHERE s.name = ? AND scd.skill_category = ? AND scd.skill_difficulty = ? AND ccec.character_class = ?
	`, skillName, category, difficulty, classCode).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSpellLearningInfoNewSystem holt alle Informationen für das Erlernen eines Zaubers
func GetSpellLearningInfoNewSystem(spellName string, classCode string) (*SpellLearningInfo, error) {
	var result SpellLearningInfo

	err := database.DB.Raw(`
		SELECT 
			s.id as spell_id,
			s.name as spell_name,
			s.stufe as spell_level,
			s.learning_category as school_name,
			cssec.character_class as class_code,
			cssec.character_class as class_name,
			cssec.ep_per_le,
			COALESCE(sllc.le_required, 0) as le_required
		FROM gsm_spells s
		JOIN learning_class_spell_school_ep_costs cssec ON s.learning_category = cssec.spell_school
		LEFT JOIN learning_spell_level_le_costs sllc ON s.stufe = sllc.level
		WHERE s.name = ? AND cssec.character_class = ?
	`, spellName, classCode).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	// Validate spell level - level 0 is not a valid spell level
	if result.SpellLevel <= 0 {
		return nil, fmt.Errorf("invalid spell level %d for spell '%s' - spell levels must be 1 or higher", result.SpellLevel, spellName)
	}

	// Validate that we found a spell (spell_id should be > 0)
	if result.SpellID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

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
		Joins("LEFT JOIN learning_sources ON gsm_skills.source_id = learning_sources.id").
		Where("gsm_skills.game_system = ? AND (learning_sources.is_active = ? OR gsm_skills.source_id IS NULL)", gameSystem, true).
		Find(&skills).Error
	return skills, err
}

// GetSpellsByActiveSources gibt alle Zauber zurück, die in aktiven Quellen definiert sind
func GetSpellsByActiveSources(gameSystem string) ([]Spell, error) {
	var spells []Spell
	err := database.DB.
		Joins("LEFT JOIN learning_sources ON gsm_spells.source_id = learning_sources.id").
		Where("gsm_spells.game_system = ? AND (learning_sources.is_active = ? OR gsm_spells.source_id IS NULL)", gameSystem, true).
		Find(&spells).Error
	return spells, err
}

// GetCharacterClassesByActiveSources gibt alle Charakterklassen zurück, die in aktiven Quellen definiert sind
func GetCharacterClassesByActiveSources(gameSystem string) ([]CharacterClass, error) {
	var classes []CharacterClass
	err := database.DB.
		Joins("LEFT JOIN learning_sources ON learning_character_classes.source_id = learning_sources.id").
		Where("learning_character_classes.game_system = ? AND (learning_sources.is_active = ? OR learning_character_classes.source_id IS NULL)", gameSystem, true).
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
