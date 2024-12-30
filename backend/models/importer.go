package models

type ImEigenschaft struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type ImAusruestung struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Beschreibung string    `json:"beschreibung"`
	Anzahl       int       `json:"anzahl"`
	BeinhaltetIn *string   `json:"beinhaltet_in"`
	Bonus        int       `json:"bonus,omitempty"`
	Gewicht      float64   `json:"gewicht"`
	Magisch      ImMagisch `json:"magisch"`
	Wert         float64   `json:"wert"`
}

type ImFertigkeit struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Beschreibung    string `json:"beschreibung"`
	Fertigkeitswert int    `json:"fertigkeitswert"`
	Bonus           int    `json:"bonus,omitempty"`
	Pp              int    `json:"pp,omitempty"`
	Quelle          string `json:"quelle"`
}

type ImZauber struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Beschreibung string `json:"beschreibung"`
	Bonus        int    `json:"bonus"`
	Quelle       string `json:"quelle"`
}

type ImWaffenfertigkeit struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Beschreibung    string `json:"beschreibung"`
	Bonus           int    `json:"bonus"`
	Fertigkeitswert int    `json:"fertigkeitswert"`
	Pp              int    `json:"pp"`
	Quelle          string `json:"quelle"`
}

type ImWaffe struct {
	ID                      string    `json:"id"`
	Name                    string    `json:"name"`
	Beschreibung            string    `json:"beschreibung"`
	Abwb                    int       `json:"abwb"`
	Anb                     int       `json:"anb"`
	Anzahl                  int       `json:"anzahl"`
	BeinhaltetIn            *string   `json:"beinhaltet_in"`
	Gewicht                 float64   `json:"gewicht"`
	Magisch                 ImMagisch `json:"magisch"`
	NameFuerSpezialisierung string    `json:"nameFuerSpezialisierung"`
	Schb                    int       `json:"schb"`
	Wert                    float64   `json:"wert"`
}

type ImMerkmale struct {
	Augenfarbe string `json:"augenfarbe"`
	Haarfarbe  string `json:"haarfarbe"`
	Sonstige   string `json:"sonstige"`
}

type ImLp struct {
	Max   int `json:"max"`
	Value int `json:"value"`
}

type ImGestalt struct {
	Breite  string `json:"breite"`
	Groesse string `json:"groesse"`
}

type ImErfahrungsschatz struct {
	Value int `json:"value"`
}

type ImEigenschaften struct {
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

type ImBennies struct {
	Gg int `json:"gg"`
	Gp int `json:"gp"`
	Sg int `json:"sg"`
}

type ImBehaeltniss struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Beschreibung string    `json:"beschreibung"`
	BeinhaltetIn any       `json:"beinhaltet_in"`
	Gewicht      float64   `json:"gewicht"`
	Magisch      ImMagisch `json:"magisch"`
	Tragkraft    float64   `json:"tragkraft"`
	Volumen      float64   `json:"volumen"`
	Wert         float64   `json:"wert"`
}

type ImAp struct {
	Max   int `json:"max"`
	Value int `json:"value"`
}

type ImB struct {
	Max   int `json:"max"`
	Value int `json:"value"`
}

type ImTransportation struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Beschreibung string    `json:"beschreibung"`
	BeinhaltetIn any       `json:"beinhaltet_in"`
	Gewicht      int       `json:"gewicht"`
	Tragkraft    float64   `json:"tragkraft"`
	Wert         float64   `json:"wert"`
	Magisch      ImMagisch `json:"magisch"`
	//Magisch   Magisch `gorm:"polymorphic:Item;polymorphicValue:Transportmittel" json:"magisch"`
}

type ImMagisch struct {
	Abw         int  `json:"abw"`
	Ausgebrannt bool `json:"ausgebrannt"`
	IstMagisch  bool `json:"ist_magisch"`
}

/*
Define models for each table
Add other models for Ausruestung, Fertigkeiten, etc., following the same pattern.
*/
type ImCharacterImport struct {
	ID                 string               `json:"id"`
	Name               string               `json:"name"`
	Rasse              string               `json:"rasse"`
	Typ                string               `json:"typ"`
	Alter              int                  `json:"alter"`
	Anrede             string               `json:"anrede"`
	Grad               int                  `json:"grad"`
	Groesse            int                  `json:"groesse"`
	Gewicht            int                  `json:"gewicht"`
	Glaube             string               `json:"glaube"`
	Hand               string               `json:"hand"`
	Fertigkeiten       []ImFertigkeit       `json:"fertigkeiten"`
	Zauber             []ImZauber           `json:"zauber"`
	Lp                 ImLp                 `json:"lp"`
	Eigenschaften      ImEigenschaften      `json:"eigenschaften"`
	Merkmale           ImMerkmale           `json:"merkmale"`
	Bennies            ImBennies            `json:"bennies"`
	Gestalt            ImGestalt            `json:"gestalt"`
	Ap                 ImAp                 `json:"ap"`
	B                  ImB                  `json:"b"`
	Erfahrungsschatz   ImErfahrungsschatz   `json:"erfahrungsschatz"`
	Transportmittel    []ImTransportation   `json:"transportmittel"`
	Ausruestung        []ImAusruestung      `json:"ausruestung"`
	Behaeltnisse       []ImBehaeltniss      `json:"behaeltnisse"`
	Waffen             []ImWaffe            `json:"waffen"`
	Waffenfertigkeiten []ImWaffenfertigkeit `json:"waffenfertigkeiten"`
	Spezialisierung    []string             `json:"spezialisierung"`
	Image              string               `json:"image,omitempty"`
}
