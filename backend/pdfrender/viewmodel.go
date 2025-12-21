package pdfrender

import "time"

// CharacterSheetViewModel represents all data needed to render a character sheet PDF
type CharacterSheetViewModel struct {
	Character     CharacterInfo
	Attributes    AttributeValues
	DerivedValues DerivedValueSet
	Skills        []SkillViewModel
	Weapons       []WeaponViewModel
	Spells        []SpellViewModel
	MagicItems    []MagicItemViewModel
	Equipment     []EquipmentViewModel
	GameResults   []GameResultViewModel
	Meta          PageMeta
}

// CharacterInfo contains basic character information
type CharacterInfo struct {
	Name       string
	Player     string
	Type       string // Charaktertyp (z.B. "Krieger", "Magier")
	Grade      int
	Birthdate  string
	Age        int
	Hand       string // "rechts." oder "links." händig
	Height     int    // in cm
	Weight     int    // in kg
	IconBase64 string // base64-kodiertes Charakterbild als Data-URI
	Gender     string
	Homeland   string
	Religion   string
	Stand      string // Sozialer Stand
	Vermoegen  WealthInfo
}

// WealthInfo contains character wealth/money
type WealthInfo struct {
	Goldstuecke   int
	Silberstuecke int
	Kupferstuecke int
}

// AttributeValues contains all character attributes
type AttributeValues struct {
	St int // Stärke
	Gs int // Geschicklichkeit
	Gw int // Gewandtheit
	Ko int // Konstitution
	In int // Intelligenz
	Zt int // Zaubertaltent
	Au int // Aussehen
	PA int // Persönliche Ausstrahlung
	Wk int // Willenskraft
	B  int // Bewegungsweite
}

// DerivedValueSet contains all derived character values
type DerivedValueSet struct {
	// Lebenspunkte & Ausdauer
	LP            int
	AP            int
	LPMax         int
	APMax         int
	LPAktuell     int
	APAktuell     int
	AusdauerBonus int
	// Geschwindigkeit
	GG int // Grundgeschwindigkeit
	SG int // Schrittgeschwindigkeit

	// Kampfwerte
	Abwehr       int // z.B. "Abwehr+12"
	SchadenBonus int
	AngriffBonus int

	// Resistenzen
	ResistenzGift    int
	ResistenzKoerper int
	ResistenzGeist   int

	// Zauberwerte
	Zaubern      int // z.B. "+10/+9"
	ZaubernBonus int // Erster Zauberbonus

	// Sonstige
	Sehen     int // Sehen-Wert
	Horen     int // Hören-Wert
	Riechen   int // Riechen-Wert
	Schmecken int // Schmecken-Wert
	Sechster  int // Sechster Sinn
}

// SkillViewModel represents a skill for display
type SkillViewModel struct {
	Name           string
	Category       string // z.B. "Kampf", "Körper", "Social"
	SkillType      string // z.B. "Fert", "Waff", "Ungelernte Fertigkeit"
	Value          int    // Erfolgswert (EW)
	BaseValue      int    // Grundwert (für Statistikseite)
	Bonus          int    // Bonus/Malus
	PracticePoints int    // Praxispunkte (PP)
	Attribute1     string // Leiteigenschaft Attribut für Bonus (z.B. "St")
	IsLearned      bool   // Ob die Fertigkeit gelernt wurde
}

// WeaponViewModel represents a weapon for display
type WeaponViewModel struct {
	Name       string
	Value      int    // Erfolgswert (EW)
	ParryValue int    // Abwehrwert (falls vorhanden)
	Damage     string // Schaden (z.B. "1W6+2")
	Range      string // Reichweite (für Fernkampfwaffen)
	Notes      string // Besondere Eigenschaften
	IsRanged   bool   // Fernkampfwaffe ja/nein
	IsMagical  bool   // Magische Waffe ja/nein
}

// SpellViewModel represents a spell for display
type SpellViewModel struct {
	Name             string
	Bonus            int    // Character's bonus for this spell
	Stufe            int    // Spell level
	AP               string // Abenteuerpunkte cost
	Art              string // Art des Zaubers (z.B. "Gestenzauber", "Wortzauber")
	Zauberdauer      string // Zauberdauer (z.B. "1 sec", "10 min")
	Reichweite       string // Reichweite (z.B. "0", "30m")
	Wirkungsziel     string // Wirkungsziel (z.B. Körper, Geist, Umgebung)
	Wirkungsbereich  string // Wirkungsbereich (z.B. "1-10 Wesen", "Kegel 5m")
	Wirkungsdauer    string // Wirkungsdauer (z.B. "0", "10 min")
	Ursprung         string // Origin/source of the spell
	Category         string // Spell school/category
	LearningCategory string // Learning category
	Beschreibung     string // Description
}

// MagicItemViewModel represents a magical item
type MagicItemViewModel struct {
	Name        string
	Description string
	Properties  string // Magische Eigenschaften
	Charges     int    // Ladungen (falls zutreffend)
	Notes       string
}

// EquipmentViewModel represents an equipment item
type EquipmentViewModel struct {
	Name        string
	Quantity    int
	Weight      float64 // Gewicht pro Stück in kg
	TotalWeight float64 // Gesamtgewicht
	Location    string  // z.B. "Am Körper", "Container 1"
	Container   string  // Container-Name (z.B. "Becher, Holz")
	Value       int     // Wert in Währungseinheiten
	Notes       string
	IsWorn      bool // Am Körper getragen
	IsContainer bool // Ist selbst ein Container
}

// GameResultViewModel represents a game session result
type GameResultViewModel struct {
	Date        time.Time
	EP          int    // Erfahrungspunkte
	Gold        int    // Gold erhalten
	Description string // Beschreibung der Sitzung
	Location    string // Ort der Handlung
	Notes       string
}

// PageMeta contains metadata about the current page
type PageMeta struct {
	Date           string
	PageNumber     int
	TotalPages     int
	IsContinuation bool   // Ist eine Fortsetzungsseite
	PageType       string // "stats", "play", "spell", "equip"
}

// PageData represents data for a single page (after pagination)
type PageData struct {
	Character     CharacterInfo
	Attributes    AttributeValues
	DerivedValues DerivedValueSet

	// Lists sliced according to template block metadata
	Skills         []SkillViewModel
	SkillsColumn1  []SkillViewModel // For two-column skill layout (page_1)
	SkillsColumn2  []SkillViewModel // For two-column skill layout (page_1)
	SkillsColumn3  []SkillViewModel // For continuation pages (page_1.2)
	SkillsColumn4  []SkillViewModel // For continuation pages (page_1.2)
	SkillsLearned  []SkillViewModel // Filtered learned skills (page_2)
	SkillsLanguage []SkillViewModel // Filtered language skills (page_2)
	Weapons        []WeaponViewModel
	Spells         []SpellViewModel
	SpellsLeft     []SpellViewModel // Left column spells (page_3)
	SpellsRight    []SpellViewModel // Right column spells (page_3)
	MagicItems     []MagicItemViewModel
	Equipment      []EquipmentViewModel
	GameResults    []GameResultViewModel

	Meta PageMeta
}
