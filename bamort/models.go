package main

type User struct {
	UserID       uint   `gorm:"primaryKey"`
	Username     string `gorm:"unique"`
	PasswordHash string
	Email        string `gorm:"unique"`
}

type Character struct {
	ID                 uint               `gorm:"primaryKey" json:"id"`
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
	Fertigkeiten       []Fertigkeit       `gorm:"foreignKey:CharacterID" json:"fertigkeiten"`
	Zauber             []Zauber           `gorm:"foreignKey:CharacterID" json:"zauber"`
	Lp                 Lp                 `gorm:"foreignKey:CharacterID" json:"lp"`
	Eigenschaften      []Eigenschaft      `gorm:"foreignKey:CharacterID" json:"eigenschaften"`
	Merkmale           Merkmale           `gorm:"foreignKey:CharacterID" json:"merkmale"`
	Bennies            Bennies            `gorm:"foreignKey:CharacterID" json:"bennies"`
	Gestalt            Gestalt            `gorm:"foreignKey:CharacterID" json:"gestalt"`
	Ap                 Ap                 `gorm:"foreignKey:CharacterID" json:"ap"`
	B                  B                  `gorm:"foreignKey:CharacterID" json:"b"`
	Erfahrungsschatz   Erfahrungsschatz   `gorm:"foreignKey:CharacterID" json:"erfahrungsschatz"`
	Transportmittel    []Transportation   `gorm:"foreignKey:CharacterID" json:"transportmittel"`
	Ausruestung        []Ausruestung      `gorm:"foreignKey:CharacterID" json:"ausruestung"`
	Behaeltnisse       []Behaeltniss      `gorm:"foreignKey:CharacterID" json:"behaeltnisse"`
	Waffen             []Waffe            `gorm:"foreignKey:CharacterID" json:"waffen"`
	Waffenfertigkeiten []Waffenfertigkeit `gorm:"foreignKey:CharacterID" json:"waffenfertigkeiten"`
	Spezialisierung    StringArray        `gorm:"type:TEXT"  json:"spezialisierung"`
	/*
		Image              string             `json:"image,omitempty"`
	*/
}

type Eigenschaft struct {
	ID          uint   `gorm:"primaryKey"`
	CharacterID uint   `gorm:"index" json:"character_id"`
	Name        string `json:"name"`
	Value       int    `json:"value"`
}

type Ausruestung struct {
	ID           uint               `gorm:"primaryKey"`
	CharacterID  uint               `gorm:"index" json:"character_id"`
	Name         string             `json:"name"`
	Beschreibung string             `json:"beschreibung"`
	Anzahl       int                `json:"anzahl"`
	BeinhaltetIn *string            `json:"beinhaltet_in"`
	Bonus        int                `json:"bonus,omitempty"`
	Gewicht      float64            `json:"gewicht"`
	Magisch      MagischAusruestung `gorm:"foreignKey:AusruestungID" json:"magisch"`
	Wert         float64            `json:"wert"`
}

type Fertigkeit struct {
	ID              uint   `gorm:"primaryKey"`
	CharacterID     uint   `gorm:"index" json:"character_id"`
	Name            string `json:"name"`
	Beschreibung    string `json:"beschreibung"`
	Fertigkeitswert int    `json:"fertigkeitswert"`
	Bonus           int    `json:"bonus,omitempty"`
	Pp              int    `json:"pp,omitempty"`
	Quelle          string `json:"quelle"`
}

type Zauber struct {
	ID           uint   `gorm:"primaryKey"`
	CharacterID  uint   `gorm:"index" json:"character_id"`
	Name         string `json:"name"`
	Beschreibung string `json:"beschreibung"`
	Bonus        int    `json:"bonus"`
	Quelle       string `json:"quelle"`
}

type Waffenfertigkeit struct {
	ID              uint   `gorm:"primaryKey"`
	CharacterID     uint   `gorm:"index" json:"character_id"`
	Name            string `json:"name"`
	Beschreibung    string `json:"beschreibung"`
	Bonus           int    `json:"bonus"`
	Fertigkeitswert int    `json:"fertigkeitswert"`
	Pp              int    `json:"pp"`
	Quelle          string `json:"quelle"`
}

type Waffe struct {
	ID                      uint         `gorm:"primaryKey"`
	CharacterID             uint         `gorm:"index" json:"character_id"`
	Name                    string       `json:"name"`
	Beschreibung            string       `json:"beschreibung"`
	Abwb                    int          `json:"abwb"`
	Anb                     int          `json:"anb"`
	Anzahl                  int          `json:"anzahl"`
	BeinhaltetIn            *string      `json:"beinhaltet_in"`
	Gewicht                 float64      `json:"gewicht"`
	Magisch                 MagischWaffe `gorm:"foreignKey:WaffenID" json:"magisch"`
	NameFuerSpezialisierung string       `json:"nameFuerSpezialisierung"`
	Schb                    int          `json:"schb"`
	Wert                    float64      `json:"wert"`
}

type Merkmale struct {
	ID          uint   `gorm:"primaryKey"`
	CharacterID uint   `gorm:"index" json:"character_id"`
	Augenfarbe  string `json:"augenfarbe"`
	Haarfarbe   string `json:"haarfarbe"`
	Sonstige    string `json:"sonstige"`
}

type Lp struct {
	ID          uint `gorm:"primaryKey"`
	CharacterID uint `gorm:"index" json:"character_id"`
	Max         int  `json:"max"`
	Value       int  `json:"value"`
}

type Gestalt struct {
	ID          uint   `gorm:"primaryKey"`
	CharacterID uint   `gorm:"index" json:"character_id"`
	Breite      string `json:"breite"`
	Groesse     string `json:"groesse"`
}

type Erfahrungsschatz struct {
	ID          uint `gorm:"primaryKey"`
	CharacterID uint `gorm:"index" json:"character_id"`
	Value       int  `json:"value"`
}

/*
type Eigenschaften struct {
	ID          uint `gorm:"primaryKey"`
	CharacterID uint `gorm:"index" json:"character_id"`
	Au          int  `json:"au"`
	Gs          int  `json:"gs"`
	Gw          int  `json:"gw"`
	In          int  `json:"in"`
	Ko          int  `json:"ko"`
	Pa          int  `json:"pa"`
	St          int  `json:"st"`
	Wk          int  `json:"wk"`
	Zt          int  `json:"zt"`
}
*/

type Bennies struct {
	ID          uint `gorm:"primaryKey"`
	CharacterID uint `gorm:"index" json:"character_id"`
	Gg          int  `json:"gg"`
	Gp          int  `json:"gp"`
	Sg          int  `json:"sg"`
}

type Behaeltniss struct {
	ID           uint   `gorm:"primaryKey"`
	CharacterID  uint   `gorm:"index" json:"character_id"`
	Name         string `json:"name"`
	Beschreibung string `json:"beschreibung"`
	//BeinhaltetIn any              `json:"beinhaltet_in"`
	Gewicht   float64          `json:"gewicht"`
	Magisch   MagischBehaelter `gorm:"foreignKey:BehaeltnissID" json:"magisch"`
	Tragkraft float64          `json:"tragkraft"`
	Volumen   float64          `json:"volumen"`
	Wert      float64          `json:"wert"`
}

type Ap struct {
	ID          uint `gorm:"primaryKey"`
	CharacterID uint `gorm:"index" json:"character_id"`
	Max         int  `json:"max"`
	Value       int  `json:"value"`
}

type B struct {
	ID          uint `gorm:"primaryKey"`
	CharacterID uint `gorm:"index" json:"character_id"`
	Max         int  `json:"max"`
	Value       int  `json:"value"`
}

type Transportation struct {
	ID           uint   `gorm:"primaryKey"`
	CharacterID  uint   `gorm:"index" json:"character_id"`
	Name         string `json:"name"`
	Beschreibung string `json:"beschreibung"`
	//BeinhaltetIn any     `json:"beinhaltet_in"`
	Gewicht   int              `json:"gewicht"`
	Tragkraft float64          `json:"tragkraft"`
	Wert      float64          `json:"wert"`
	Magisch   MagischTransport `gorm:"foreignKey:TransportationID" json:"magisch"`
	//Magisch   Magisch `gorm:"polymorphic:Item;polymorphicValue:Transportmittel" json:"magisch"`
}

type MagischAusruestung struct {
	ID uint `gorm:"primaryKey" json:"id"`
	//ItemType    string `gorm:"index" json:"item_type"` // Type of the referenced item (e.g., "Ausruestung")
	//ItemID      uint   `gorm:"index" json:"item_id"`   // ID of the referenced item
	AusruestungID int  `gorm:"index" json:"ausruestung_id"`
	Abw           int  `json:"abw"`
	Ausgebrannt   bool `json:"ausgebrannt"`
	IstMagisch    bool `json:"ist_magisch"`
}

type MagischWaffe struct {
	ID uint `gorm:"primaryKey" json:"id"`
	//ItemType    string `gorm:"index" json:"item_type"` // Type of the referenced item (e.g., "Ausruestung")
	//ItemID      uint   `gorm:"index" json:"item_id"`   // ID of the referenced item
	WaffenID    int  `gorm:"index" json:"waffen_id"`
	Abw         int  `json:"abw"`
	Ausgebrannt bool `json:"ausgebrannt"`
	IstMagisch  bool `json:"ist_magisch"`
}

type MagischBehaelter struct {
	ID uint `gorm:"primaryKey" json:"id"`
	//ItemType    string `gorm:"index" json:"item_type"` // Type of the referenced item (e.g., "Ausruestung")
	//ItemID      uint   `gorm:"index" json:"item_id"`   // ID of the referenced item
	BehaeltnissID int  `gorm:"index" json:"behaeltniss_id"`
	Abw           int  `json:"abw"`
	Ausgebrannt   bool `json:"ausgebrannt"`
	IstMagisch    bool `json:"ist_magisch"`
}

type MagischTransport struct {
	ID uint `gorm:"primaryKey" json:"id"`
	//ItemType    string `gorm:"index" json:"item_type"` // Type of the referenced item (e.g., "Ausruestung")
	//ItemID      uint   `gorm:"index" json:"item_id"`   // ID of the referenced item
	TransportationID int  `gorm:"index" json:"transportation_id"`
	Abw              int  `json:"abw"`
	Ausgebrannt      bool `json:"ausgebrannt"`
	IstMagisch       bool `json:"ist_magisch"`
}

/*
Define models for each table
Add other models for Ausruestung, Fertigkeiten, etc., following the same pattern.
*/
