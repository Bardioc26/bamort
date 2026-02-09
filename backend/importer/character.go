package importer

// This file contains the character import data structures.
// These were copied from the deprecated importero package to break dependencies.

// ImportBase provides common fields for imported items
type ImportBase struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Magisch represents magical properties of items
type Magisch struct {
	Abw         int  `json:"abw"`
	Ausgebrannt bool `json:"ausgebrannt"`
	IstMagisch  bool `json:"ist_magisch"`
}

// Ausruestung represents equipment/gear
type Ausruestung struct {
	ImportBase
	Beschreibung string  `json:"beschreibung"`
	Anzahl       int     `json:"anzahl"`
	BeinhaltetIn string  `json:"beinhaltet_in"`
	ContainedIn  uint    `json:"contained_in"`
	Bonus        int     `json:"bonus,omitempty"`
	Gewicht      float64 `json:"gewicht"`
	Magisch      Magisch `json:"magisch"`
	Wert         float64 `json:"wert"`
}

// Waffe represents a weapon
type Waffe struct {
	ImportBase
	Beschreibung            string  `json:"beschreibung"`
	Abwb                    int     `json:"abwb"`
	Anb                     int     `json:"anb"`
	Anzahl                  int     `json:"anzahl"`
	BeinhaltetIn            string  `json:"beinhaltet_in"`
	ContainedIn             uint    `json:"contained_in"`
	Gewicht                 float64 `json:"gewicht"`
	Magisch                 Magisch `json:"magisch"`
	NameFuerSpezialisierung string  `json:"nameFuerSpezialisierung"`
	Schb                    int     `json:"schb"`
	Wert                    float64 `json:"wert"`
}

// Behaeltniss represents a container
type Behaeltniss struct {
	ImportBase
	Beschreibung string  `json:"beschreibung"`
	BeinhaltetIn string  `json:"beinhaltet_in"`
	ContainedIn  uint    `json:"contained_in"`
	Gewicht      float64 `json:"gewicht"`
	Magisch      Magisch `json:"magisch"`
	Tragkraft    float64 `json:"tragkraft"`
	Volumen      float64 `json:"volumen"`
	Wert         float64 `json:"wert"`
}

// Transportation represents a means of transport
type Transportation struct {
	ImportBase
	Beschreibung string  `json:"beschreibung"`
	BeinhaltetIn string  `json:"beinhaltet_in"`
	ContainedIn  uint    `json:"contained_in"`
	Gewicht      int     `json:"gewicht"`
	Tragkraft    float64 `json:"tragkraft"`
	Wert         float64 `json:"wert"`
	Magisch      Magisch `json:"magisch"`
}

// Fertigkeit represents a skill
type Fertigkeit struct {
	ImportBase
	Beschreibung    string `json:"beschreibung"`
	Fertigkeitswert int    `json:"fertigkeitswert"`
	Bonus           int    `json:"bonus,omitempty"`
	Pp              int    `json:"pp,omitempty"`
	Quelle          string `json:"quelle"`
}

// Zauber represents a spell
type Zauber struct {
	ImportBase
	Beschreibung string `json:"beschreibung"`
	Bonus        int    `json:"bonus"`
	Quelle       string `json:"quelle"`
}

// Waffenfertigkeit represents a weapon skill
type Waffenfertigkeit struct {
	ImportBase
	Beschreibung    string `json:"beschreibung"`
	Bonus           int    `json:"bonus"`
	Fertigkeitswert int    `json:"fertigkeitswert"`
	Pp              int    `json:"pp"`
	Quelle          string `json:"quelle"`
}

// Eigenschaft represents a character attribute (deprecated field, kept for compatibility)
type Eigenschaft struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// Merkmale represents character features/traits
type Merkmale struct {
	Augenfarbe string `json:"augenfarbe"`
	Haarfarbe  string `json:"haarfarbe"`
	Sonstige   string `json:"sonstige"`
}

// Lp represents life points/hit points
type Lp struct {
	Max   int `json:"max"`
	Value int `json:"value"`
}

// Gestalt represents character build/physique
type Gestalt struct {
	Breite  string `json:"breite"`
	Groesse string `json:"groesse"`
}

// Erfahrungsschatz represents experience points
type Erfahrungsschatz struct {
	Value int `json:"value"`
}

// Eigenschaften represents the main character attributes
type Eigenschaften struct {
	Au int `json:"au"` // Aussehen
	Gs int `json:"gs"` // Geschicklichkeit
	Gw int `json:"gw"` // Gewandtheit
	In int `json:"in"` // Intelligenz
	Ko int `json:"ko"` // Konstitution
	Pa int `json:"pa"` // Persönliche Ausstrahlung
	St int `json:"st"` // Stärke
	Wk int `json:"wk"` // Willenskraft
	Zt int `json:"zt"` // Zähigkeit
}

// Bennies represents fate points/luck points
type Bennies struct {
	Gg int `json:"gg"` // Göttliche Gnade
	Gp int `json:"gp"` // Glückspunkte
	Sg int `json:"sg"` // Schicksalsgunst
}

// Ap represents action points/stamina
type Ap struct {
	Max   int `json:"max"`
	Value int `json:"value"`
}

// B represents movement points
type B struct {
	Max   int `json:"max"`
	Value int `json:"value"`
}

// CharacterImport represents the complete character data for import
// This is the canonical interchange format for character data
type CharacterImport struct {
	ID                 string             `json:"id"`
	Name               string             `json:"name"`
	Rasse              string             `json:"rasse"`
	Typ                string             `json:"typ"`
	Alter              int                `json:"alter"`
	Anrede             string             `json:"anrede"`
	Grad               int                `json:"grad"`
	Groesse            int                `json:"groesse"`
	Gewicht            int                `json:"gewicht"`
	Glaube             string             `json:"glaube"`
	Hand               string             `json:"hand"`
	Fertigkeiten       []Fertigkeit       `json:"fertigkeiten"`
	Zauber             []Zauber           `json:"zauber"`
	Lp                 Lp                 `json:"lp"`
	Eigenschaften      Eigenschaften      `json:"eigenschaften"`
	Merkmale           Merkmale           `json:"merkmale"`
	Bennies            Bennies            `json:"bennies"`
	Gestalt            Gestalt            `json:"gestalt"`
	Ap                 Ap                 `json:"ap"`
	B                  B                  `json:"b"`
	Erfahrungsschatz   Erfahrungsschatz   `json:"erfahrungsschatz"`
	Transportmittel    []Transportation   `json:"transportmittel"`
	Ausruestung        []Ausruestung      `json:"ausruestung"`
	Behaeltnisse       []Behaeltniss      `json:"behaeltnisse"`
	Waffen             []Waffe            `json:"waffen"`
	Waffenfertigkeiten []Waffenfertigkeit `json:"waffenfertigkeiten"`
	Spezialisierung    []string           `json:"spezialisierung"`
	Image              string             `json:"image,omitempty"`
}
