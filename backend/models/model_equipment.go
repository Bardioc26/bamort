package models

import (
	"bamort/database"
)

type EqAusruestung struct {
	BamortCharTrait
	Magisch
	Beschreibung  string  `json:"beschreibung"`
	Anzahl        int     `json:"anzahl"`
	BeinhaltetIn  string  `json:"beinhaltet_in"`
	ContainedIn   uint    `json:"contained_in"`
	ContainerType string  `json:"container_type"`
	Bonus         int     `json:"bonus,omitempty"`
	Gewicht       float64 `json:"gewicht"`
	Wert          float64 `json:"wert"`
}

type EqWaffe struct {
	BamortCharTrait
	Magisch
	Beschreibung            string  `json:"beschreibung"`
	Abwb                    int     `json:"abwb"` // Abwehrbonus
	Anb                     int     `json:"anb"`  // Angriffsbonus
	Anzahl                  int     `json:"anzahl"`
	BeinhaltetIn            string  `json:"beinhaltet_in"`
	ContainedIn             uint    `json:"contained_in"`
	ContainerType           string  `json:"container_type"`
	Gewicht                 float64 `json:"gewicht"`
	NameFuerSpezialisierung string  `json:"nameFuerSpezialisierung"`
	Schb                    int     `json:"schb"` // Schadensbonus
	Wert                    float64 `json:"wert"`
}

type EqContainer struct {
	BamortCharTrait
	Magisch
	Beschreibung     string  `json:"beschreibung"`
	BeinhaltetIn     string  `json:"beinhaltet_in"`
	ContainedIn      uint    `json:"contained_in"`
	IsTransportation bool    `json:"is_transportation"`
	Gewicht          float64 `json:"gewicht"`
	Wert             float64 `json:"wert"`
	Tragkraft        float64 `json:"tragkraft"`
	Volumen          float64 `json:"volumen"`
	ExtID            string  `json:"ext_id"`
}

func (object *EqAusruestung) TableName() string {
	dbPrefix := "equi"
	return dbPrefix + "_" + "equipments"
}
func (object *EqWaffe) TableName() string {
	dbPrefix := "equi"
	return dbPrefix + "_" + "weapons"
}
func (object *EqContainer) TableName() string {
	dbPrefix := "equi"
	return dbPrefix + "_" + "containers"
}

func (object *EqContainer) FirstExtId(id string) error {
	err := database.DB.
		First(&object, "ext_id = ?", id).Error
	if err != nil {
		// Container not found
		return err
	}
	return nil
}
func (object *EqContainer) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// Container saved
		return err
	}
	return nil
}

func (object *EqAusruestung) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// Ausruestung saved
		return err
	}
	return nil
}
func (object *EqWaffe) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// Waffe saved
		return err
	}
	return nil
}

func (object *EqAusruestung) LinkContainer() error {
	//var err error
	if object.BeinhaltetIn != "" {
		co := EqContainer{}
		co.FirstExtId(object.BeinhaltetIn)
		if co.ID > 0 {
			object.ContainedIn = co.ID
			//err= object.Save()
		}
		if object.ContainedIn > 0 {
			err := object.Save()
			if err != nil {
				// ContainedIn saved
				return err
			}
		}
	}
	return nil
}

func (object *EqWaffe) LinkContainer() error {
	//var err error
	if object.BeinhaltetIn != "" {
		co := EqContainer{}
		co.FirstExtId(object.BeinhaltetIn)
		if co.ID > 0 {
			object.ContainedIn = co.ID
			//err= object.Save()
		}
		if object.ContainedIn > 0 {
			err := object.Save()
			if err != nil {
				// ContainedIn saved
				return err
			}
		}
	}
	return nil
}

func (object *EqContainer) LinkContainer() error {
	//var err error
	if object.BeinhaltetIn != "" {
		co := EqContainer{}
		co.FirstExtId(object.BeinhaltetIn)
		if co.ID > 0 {
			object.ContainedIn = co.ID
			//err= object.Save()
		}
		if object.ContainedIn > 0 {
			err := object.Save()
			if err != nil {
				// ContainedIn saved
				return err
			}
		}
	}
	return nil
}
