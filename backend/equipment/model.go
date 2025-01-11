package equipment

import "bamort/models"

var dbPrefix = "equi"

type Ausruestung struct {
	models.BamortCharTrait
	models.Magisch
	Beschreibung string  `json:"beschreibung"`
	Anzahl       int     `json:"anzahl"`
	BeinhaltetIn string  `json:"beinhaltet_in"`
	Bonus        int     `json:"bonus,omitempty"`
	Gewicht      float64 `json:"gewicht"`
	Wert         float64 `json:"wert"`
}

type Waffe struct {
	models.BamortCharTrait
	models.Magisch
	Beschreibung            string  `json:"beschreibung"`
	Abwb                    int     `json:"abwb"`
	Anb                     int     `json:"anb"`
	Anzahl                  int     `json:"anzahl"`
	BeinhaltetIn            string  `json:"beinhaltet_in"`
	Gewicht                 float64 `json:"gewicht"`
	NameFuerSpezialisierung string  `json:"nameFuerSpezialisierung"`
	Schb                    int     `json:"schb"`
	Wert                    float64 `json:"wert"`
}

type Behaeltniss struct {
	models.BamortCharTrait
	models.Magisch
	Beschreibung string  `json:"beschreibung"`
	BeinhaltetIn string  `json:"beinhaltet_in"`
	Gewicht      float64 `json:"gewicht"`
	Wert         float64 `json:"wert"`
	Tragkraft    float64 `json:"tragkraft"`
	Volumen      float64 `json:"volumen"`
}

type Transportation struct {
	models.BamortCharTrait
	models.Magisch
	Beschreibung string  `json:"beschreibung"`
	BeinhaltetIn string  `json:"beinhaltet_in"`
	Gewicht      float64 `json:"gewicht"`
	Wert         float64 `json:"wert"`
	Tragkraft    float64 `json:"tragkraft"`
	Volumen      float64 `json:"volumen"`
}

func (object *Ausruestung) TableName() string {
	return dbPrefix + "_" + "equipments"
}
func (object *Waffe) TableName() string {
	return dbPrefix + "_" + "weapons"
}
func (object *Behaeltniss) TableName() string {
	return dbPrefix + "_" + "containers"
}
func (object *Transportation) TableName() string {
	return dbPrefix + "_" + "transportationss"
}
