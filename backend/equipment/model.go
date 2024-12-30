package equipment

import "bamort/models"

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
