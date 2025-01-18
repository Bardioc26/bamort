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

type Behaeltniss struct {
	models.BamortCharTrait
	models.Magisch
	Beschreibung  string  `json:"beschreibung"`
	BeinhaltetIn  string  `json:"beinhaltet_in"`
	ContainedIn   uint    `json:"contained_in"`
	ContainerType string  `json:"container_type"`
	Gewicht       float64 `json:"gewicht"`
	Wert          float64 `json:"wert"`
	Tragkraft     float64 `json:"tragkraft"`
	Volumen       float64 `json:"volumen"`
	ExtID         string  `json:"ext_id"`
}

type Transportation struct {
	models.BamortCharTrait
	models.Magisch
	Beschreibung  string  `json:"beschreibung"`
	BeinhaltetIn  string  `json:"beinhaltet_in"`
	ContainedIn   uint    `json:"contained_in"`
	ContainerType string  `json:"container_type"`
	Gewicht       float64 `json:"gewicht"`
	Wert          float64 `json:"wert"`
	Tragkraft     float64 `json:"tragkraft"`
	Volumen       float64 `json:"volumen"`
	ExtID         string  `json:"ext_id"`
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

func (object *Behaeltniss) FirstExtId(id string) error {
	err := database.DB.
		First(&object, "ext_id = ?", id).Error
	if err != nil {
		// Behaeltniss not found
		return err
	}
	return nil
}
func (object *Behaeltniss) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// Behaeltniss saved
		return err
	}
	return nil
}
func (object *Transportation) FirstExtId(id string) error {
	err := database.DB.
		First(&object, "ext_id = ?", id).Error
	if err != nil {
		// Transportation not found
		return err
	}
	return nil
}
func (object *Transportation) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// Transportation saved
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
		co := Behaeltniss{}
		co.FirstExtId(object.BeinhaltetIn)
		if co.ID > 0 {
			object.ContainedIn = co.ID
			//err= object.Save()
		} else {
			obj := Transportation{}
			if co.ID > 0 {
				object.ContainedIn = obj.ID
				//err= object.Save()
			}
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
		co := Behaeltniss{}
		co.FirstExtId(object.BeinhaltetIn)
		if co.ID > 0 {
			object.ContainedIn = co.ID
			object.ContainerType = "Behaeltniss"
			//err= object.Save()
		} else {
			obj := Transportation{}
			obj.FirstExtId(object.BeinhaltetIn)
			if obj.ID > 0 {
				object.ContainedIn = obj.ID
				object.ContainerType = "Transportation"
				//err= object.Save()
			}
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

func (object *Behaeltniss) LinkContainer() error {
	//var err error
	if object.BeinhaltetIn != "" {
		co := Behaeltniss{}
		co.FirstExtId(object.BeinhaltetIn)
		if co.ID > 0 {
			object.ContainedIn = co.ID
			//err= object.Save()
		} else {
			obj := Transportation{}
			if co.ID > 0 {
				object.ContainedIn = obj.ID
				//err= object.Save()
			}
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

func (object *Transportation) LinkContainer() error {
	//var err error
	if object.BeinhaltetIn != "" {
		co := Behaeltniss{}
		co.FirstExtId(object.BeinhaltetIn)
		if co.ID > 0 {
			object.ContainedIn = co.ID
			//err= object.Save()
		} else {
			obj := Transportation{}
			if co.ID > 0 {
				object.ContainedIn = obj.ID
				//err= object.Save()
			}
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
