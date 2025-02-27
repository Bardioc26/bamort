package equipment

import (
	"bamort/database"
	"bamort/models"
)

var dbPrefix = "equi"

type Ausruestung struct {
	models.BamortCharTrait
	models.Magisch
	Beschreibung  string  `json:"beschreibung"`
	Anzahl        int     `json:"anzahl"`
	BeinhaltetIn  string  `json:"beinhaltet_in"`
	ContainedIn   uint    `json:"contained_in"`
	ContainerType string  `json:"container_type"`
	Bonus         int     `json:"bonus,omitempty"`
	Gewicht       float64 `json:"gewicht"`
	Wert          float64 `json:"wert"`
}

type Waffe struct {
	models.BamortCharTrait
	models.Magisch
	Beschreibung            string  `json:"beschreibung"`
	Abwb                    int     `json:"abwb"`
	Anb                     int     `json:"anb"`
	Anzahl                  int     `json:"anzahl"`
	BeinhaltetIn            string  `json:"beinhaltet_in"`
	ContainedIn             uint    `json:"contained_in"`
	ContainerType           string  `json:"container_type"`
	Gewicht                 float64 `json:"gewicht"`
	NameFuerSpezialisierung string  `json:"nameFuerSpezialisierung"`
	Schb                    int     `json:"schb"`
	Wert                    float64 `json:"wert"`
}

type Container struct {
	models.BamortCharTrait
	models.Magisch
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

func (object *Ausruestung) TableName() string {
	return dbPrefix + "_" + "equipments"
}
func (object *Waffe) TableName() string {
	return dbPrefix + "_" + "weapons"
}
func (object *Container) TableName() string {
	return dbPrefix + "_" + "containers"
}

func (object *Container) FirstExtId(id string) error {
	err := database.DB.
		First(&object, "ext_id = ?", id).Error
	if err != nil {
		// Container not found
		return err
	}
	return nil
}
func (object *Container) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// Container saved
		return err
	}
	return nil
}

func (object *Ausruestung) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// Ausruestung saved
		return err
	}
	return nil
}
func (object *Waffe) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// Waffe saved
		return err
	}
	return nil
}

func (object *Ausruestung) LinkContainer() error {
	//var err error
	if object.BeinhaltetIn != "" {
		co := Container{}
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

func (object *Waffe) LinkContainer() error {
	//var err error
	if object.BeinhaltetIn != "" {
		co := Container{}
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

func (object *Container) LinkContainer() error {
	//var err error
	if object.BeinhaltetIn != "" {
		co := Container{}
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
