package main

type CharType struct {
	Alter  int    `json:"alter"`
	Anrede string `json:"anrede"`
	Ap     struct {
		Max   int `json:"max"`
		Value int `json:"value"`
	} `json:"ap"`
	Ausruestung []struct {
		Anzahl       int     `json:"anzahl"`
		BeinhaltetIn *string `json:"beinhaltet_in"`
		Beschreibung *string `json:"beschreibung"`
		Bonus        int     `json:"bonus,omitempty"`
		Gewicht      float64 `json:"gewicht"`
		ID           string  `json:"id"`
		Magisch      struct {
			Abw         int  `json:"abw"`
			Ausgebrannt bool `json:"ausgebrannt"`
			IstMagisch  bool `json:"ist_magisch"`
		} `json:"magisch"`
		Name string  `json:"name"`
		Wert float64 `json:"wert"`
	} `json:"ausruestung"`
	B struct {
		Max int `json:"max"`
	} `json:"b"`
	Behaeltnisse []struct {
		BeinhaltetIn any     `json:"beinhaltet_in"`
		Beschreibung *string `json:"beschreibung"`
		Gewicht      float64 `json:"gewicht"`
		ID           string  `json:"id"`
		Magisch      struct {
			Abw         int  `json:"abw"`
			Ausgebrannt bool `json:"ausgebrannt"`
			IstMagisch  bool `json:"ist_magisch"`
		} `json:"magisch"`
		Name      string  `json:"name"`
		Tragkraft float64 `json:"tragkraft"`
		Volumen   float64 `json:"volumen"`
		Wert      float64 `json:"wert"`
	} `json:"behaeltnisse"`
	Bennies struct {
		Gg int `json:"gg"`
		Gp int `json:"gp"`
		Sg int `json:"sg"`
	} `json:"bennies"`
	Eigenschaften struct {
		Au int `json:"au"`
		Gs int `json:"gs"`
		Gw int `json:"gw"`
		In int `json:"in"`
		Ko int `json:"ko"`
		Pa int `json:"pa"`
		St int `json:"st"`
		Wk int `json:"wk"`
		Zt int `json:"zt"`
	} `json:"eigenschaften"`
	Erfahrungsschatz struct {
		Value int `json:"value"`
	} `json:"erfahrungsschatz"`
	Fertigkeiten []struct {
		Beschreibung    *string `json:"beschreibung"`
		Bonus           int     `json:"bonus,omitempty"`
		Fertigkeitswert int     `json:"fertigkeitswert"`
		ID              string  `json:"id"`
		Name            string  `json:"name"`
		Pp              int     `json:"pp,omitempty"`
		Quelle          string  `json:"quelle"`
	} `json:"fertigkeiten"`
	Gestalt struct {
		Breite  string `json:"breite"`
		Groesse string `json:"groesse"`
	} `json:"gestalt"`
	Gewicht int    `json:"gewicht"`
	Glaube  string `json:"glaube"`
	Grad    int    `json:"grad"`
	Groesse int    `json:"groesse"`
	Hand    string `json:"hand"`
	ID      string `json:"id"`
	Image   string `json:"image,omitempty"`
	Lp      struct {
		Max   int `json:"max"`
		Value int `json:"value"`
	} `json:"lp"`
	Merkmale struct {
		Augenfarbe string `json:"augenfarbe"`
		Haarfarbe  string `json:"haarfarbe"`
		Sonstige   string `json:"sonstige"`
	} `json:"merkmale"`
	Name            string   `json:"name"`
	Rasse           string   `json:"rasse"`
	Spezialisierung []string `json:"spezialisierung"`
	Stand           string   `json:"stand"`
	Transportmittel []struct {
		BeinhaltetIn any     `json:"beinhaltet_in"`
		Beschreibung *string `json:"beschreibung"`
		Gewicht      int     `json:"gewicht"`
		ID           string  `json:"id"`
		Magisch      struct {
			Abw         int  `json:"abw"`
			Ausgebrannt bool `json:"ausgebrannt"`
			IstMagisch  bool `json:"ist_magisch"`
		} `json:"magisch"`
		Name      string  `json:"name"`
		Tragkraft float64 `json:"tragkraft"`
		Wert      float64 `json:"wert"`
	} `json:"transportmittel"`
	Typ    string `json:"typ"`
	Waffen []struct {
		Abwb         int     `json:"abwb"`
		Anb          int     `json:"anb"`
		Anzahl       int     `json:"anzahl"`
		BeinhaltetIn *string `json:"beinhaltet_in"`
		Beschreibung *string `json:"beschreibung"`
		Gewicht      float64 `json:"gewicht"`
		ID           string  `json:"id"`
		Magisch      struct {
			Abw         int  `json:"abw"`
			Ausgebrannt bool `json:"ausgebrannt"`
			IstMagisch  bool `json:"ist_magisch"`
		} `json:"magisch"`
		Name                    string  `json:"name"`
		NameFuerSpezialisierung string  `json:"nameFuerSpezialisierung"`
		Schb                    int     `json:"schb"`
		Wert                    float64 `json:"wert"`
	} `json:"waffen"`
	Waffenfertigkeiten []struct {
		Beschreibung    *string `json:"beschreibung"`
		Bonus           int     `json:"bonus"`
		Fertigkeitswert int     `json:"fertigkeitswert"`
		ID              string  `json:"id"`
		Name            string  `json:"name"`
		Pp              int     `json:"pp"`
		Quelle          string  `json:"quelle"`
	} `json:"waffenfertigkeiten"`
	Zauber []struct {
		Beschreibung any    `json:"beschreibung"`
		Bonus        int    `json:"bonus"`
		ID           string `json:"id"`
		Name         string `json:"name"`
		Quelle       string `json:"quelle"`
	} `json:"zauber"`
}
