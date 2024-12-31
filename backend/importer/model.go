package importer

type BamortBase struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//gsmaster

type Magisch struct {
	Abw         int  `json:"abw"`
	Ausgebrannt bool `json:"ausgebrannt"`
	IstMagisch  bool `json:"ist_magisch"`
}

type Ausruestung struct {
	BamortBase
	Beschreibung string  `json:"beschreibung"`
	Anzahl       int     `json:"anzahl"`
	BeinhaltetIn *string `json:"beinhaltet_in"`
	Bonus        int     `json:"bonus,omitempty"`
	Gewicht      float64 `json:"gewicht"`
	Magisch      Magisch `json:"magisch"`
	Wert         float64 `json:"wert"`
}

type Waffe struct {
	BamortBase
	Beschreibung            string  `json:"beschreibung"`
	Abwb                    int     `json:"abwb"`
	Anb                     int     `json:"anb"`
	Anzahl                  int     `json:"anzahl"`
	BeinhaltetIn            *string `json:"beinhaltet_in"`
	Gewicht                 float64 `json:"gewicht"`
	Magisch                 Magisch `json:"magisch"`
	NameFuerSpezialisierung string  `json:"nameFuerSpezialisierung"`
	Schb                    int     `json:"schb"`
	Wert                    float64 `json:"wert"`
}

type Behaeltniss struct {
	BamortBase
	Beschreibung string  `json:"beschreibung"`
	BeinhaltetIn any     `json:"beinhaltet_in"`
	Gewicht      float64 `json:"gewicht"`
	Magisch      Magisch `json:"magisch"`
	Tragkraft    float64 `json:"tragkraft"`
	Volumen      float64 `json:"volumen"`
	Wert         float64 `json:"wert"`
}

type Transportation struct {
	BamortBase
	Beschreibung string  `json:"beschreibung"`
	BeinhaltetIn any     `json:"beinhaltet_in"`
	Gewicht      int     `json:"gewicht"`
	Tragkraft    float64 `json:"tragkraft"`
	Wert         float64 `json:"wert"`
	Magisch      Magisch `json:"magisch"`
	//Magisch   Magisch `gorm:"polymorphic:Item;polymorphicValue:Transportmittel" json:"magisch"`
}

type Fertigkeit struct {
	BamortBase
	Beschreibung    string `json:"beschreibung"`
	Fertigkeitswert int    `json:"fertigkeitswert"`
	Bonus           int    `json:"bonus,omitempty"`
	Pp              int    `json:"pp,omitempty"`
	Quelle          string `json:"quelle"`
}
type Zauber struct {
	BamortBase
	Beschreibung string `json:"beschreibung"`
	Bonus        int    `json:"bonus"`
	Quelle       string `json:"quelle"`
}

type Waffenfertigkeit struct {
	BamortBase
	Beschreibung    string `json:"beschreibung"`
	Bonus           int    `json:"bonus"`
	Fertigkeitswert int    `json:"fertigkeitswert"`
	Pp              int    `json:"pp"`
	Quelle          string `json:"quelle"`
}

// Character data
type Eigenschaft struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Merkmale struct {
	Augenfarbe string `json:"augenfarbe"`
	Haarfarbe  string `json:"haarfarbe"`
	Sonstige   string `json:"sonstige"`
}

type Lp struct {
	Max   int `json:"max"`
	Value int `json:"value"`
}

type Gestalt struct {
	Breite  string `json:"breite"`
	Groesse string `json:"groesse"`
}

type Erfahrungsschatz struct {
	Value int `json:"value"`
}

type Eigenschaften struct {
	Au int `json:"au"`
	Gs int `json:"gs"`
	Gw int `json:"gw"`
	In int `json:"in"`
	Ko int `json:"ko"`
	Pa int `json:"pa"`
	St int `json:"st"`
	Wk int `json:"wk"`
	Zt int `json:"zt"`
}

type Bennies struct {
	Gg int `json:"gg"`
	Gp int `json:"gp"`
	Sg int `json:"sg"`
}

type Ap struct {
	Max   int `json:"max"`
	Value int `json:"value"`
}

type B struct {
	Max   int `json:"max"`
	Value int `json:"value"`
}

/*
Define models for each table
Add other models for Ausruestung, Fertigkeiten, etc., following the same pattern.
*/
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
