package importer

import (
	"bamort/gsmaster"
	"fmt"
	"regexp"
)

type ImportBase struct {
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
	ImportBase
	Beschreibung string  `json:"beschreibung"`
	Anzahl       int     `json:"anzahl"`
	BeinhaltetIn string  `json:"beinhaltet_in"`
	Bonus        int     `json:"bonus,omitempty"`
	Gewicht      float64 `json:"gewicht"`
	Magisch      Magisch `json:"magisch"`
	Wert         float64 `json:"wert"`
}

type Waffe struct {
	ImportBase
	Beschreibung            string  `json:"beschreibung"`
	Abwb                    int     `json:"abwb"`
	Anb                     int     `json:"anb"`
	Anzahl                  int     `json:"anzahl"`
	BeinhaltetIn            string  `json:"beinhaltet_in"`
	Gewicht                 float64 `json:"gewicht"`
	Magisch                 Magisch `json:"magisch"`
	NameFuerSpezialisierung string  `json:"nameFuerSpezialisierung"`
	Schb                    int     `json:"schb"`
	Wert                    float64 `json:"wert"`
}

type Behaeltniss struct {
	ImportBase
	Beschreibung string  `json:"beschreibung"`
	BeinhaltetIn string  `json:"beinhaltet_in"`
	Gewicht      float64 `json:"gewicht"`
	Magisch      Magisch `json:"magisch"`
	Tragkraft    float64 `json:"tragkraft"`
	Volumen      float64 `json:"volumen"`
	Wert         float64 `json:"wert"`
}

type Transportation struct {
	ImportBase
	Beschreibung string  `json:"beschreibung"`
	BeinhaltetIn string  `json:"beinhaltet_in"`
	Gewicht      int     `json:"gewicht"`
	Tragkraft    float64 `json:"tragkraft"`
	Wert         float64 `json:"wert"`
	Magisch      Magisch `json:"magisch"`
	//Magisch   Magisch `gorm:"polymorphic:Item;polymorphicValue:Transportmittel" json:"magisch"`
}

type Fertigkeit struct {
	ImportBase
	Beschreibung    string `json:"beschreibung"`
	Fertigkeitswert int    `json:"fertigkeitswert"`
	Bonus           int    `json:"bonus,omitempty"`
	Pp              int    `json:"pp,omitempty"`
	Quelle          string `json:"quelle"`
}
type Zauber struct {
	ImportBase
	Beschreibung string `json:"beschreibung"`
	Bonus        int    `json:"bonus"`
	Quelle       string `json:"quelle"`
}

type Waffenfertigkeit struct {
	ImportBase
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

const (
	SKILL_HOEREN       = "Hören"
	SKILL_NACHTSICHT   = "Nachtsicht"
	SKILL_RIECHEN      = "Riechen"
	SKILL_SECHSTERSINN = "Sechster Sinn"
	SKILL_SEHEN        = "Sehen"
	SKILL_WAHRNEHMUNG  = "Wahrnehmung"
)

var nonImprovableSkills = []string{
	SKILL_HOEREN,
	SKILL_NACHTSICHT,
	SKILL_RIECHEN,
	SKILL_SECHSTERSINN,
	SKILL_SEHEN,
	SKILL_WAHRNEHMUNG,
}

func isImprovableSkill(name string) bool {
	for _, skill := range nonImprovableSkills {
		if skill == name {
			return false
		}
	}
	return true
}
func TransformImportFertigkeit2GSDMaster(object *Fertigkeit) (*gsmaster.Skill, error) {
	gsmobj := gsmaster.Skill{}

	err := gsmobj.First(object.Name)
	// if found check if we need to adjust masterdata
	if err == nil {
		return &gsmobj, nil
	}
	// if not found insert to masterdata
	gsmobj.Name = object.Name
	gsmobj.Beschreibung = object.Beschreibung
	gsmobj.Initialwert = object.Fertigkeitswert
	gsmobj.Quelle = object.Quelle
	gsmobj.Bonuseigenschaft = "check"
	gsmobj.Improvable = false
	re := regexp.MustCompile(`moam-ability-\\d+`)
	if re.MatchString(object.ID) {
		gsmobj.Improvable = true
	}
	if !gsmobj.Improvable {
		gsmobj.Improvable = isImprovableSkill(gsmobj.Name)
	}
	gsmobj.System = "midgard"
	err = gsmobj.Create()
	if err != nil {
		return nil, fmt.Errorf("creating gsmaster record failed: %s", err)
	}
	return &gsmobj, nil
}

func TransformImportWaffenFertigkeit2GSDMaster(object *Waffenfertigkeit) (*gsmaster.WeaponSkill, error) {
	gsmobj := gsmaster.WeaponSkill{}

	err := gsmobj.First(object.Name)
	// if found check if we need to adjust masterdata
	if err == nil {
		return &gsmobj, nil
	}
	// if not found insert to masterdata
	gsmobj.Name = object.Name
	gsmobj.Beschreibung = object.Beschreibung
	gsmobj.Initialwert = object.Fertigkeitswert
	gsmobj.Quelle = object.Quelle
	gsmobj.Bonuseigenschaft = "check"
	gsmobj.Improvable = true
	gsmobj.System = "midgard"
	err = gsmobj.Create()
	if err != nil {
		return nil, fmt.Errorf("creating gsmaster record failed: %s", err)
	}
	return &gsmobj, nil
}

func TransformImportSpell2GSDMaster(object *Zauber) (*gsmaster.Spell, error) {
	gsmobj := gsmaster.Spell{}

	err := gsmobj.First(object.Name)
	// if found check if we need to adjust masterdata
	if err == nil {
		return &gsmobj, nil
	}
	// if not found insert to masterdata
	gsmobj.Name = object.Name
	gsmobj.Beschreibung = object.Beschreibung
	gsmobj.AP = 0
	gsmobj.Stufe = 0
	gsmobj.Reichweite = 0
	gsmobj.Wirkungsziel = ""
	gsmobj.Quelle = object.Quelle
	gsmobj.System = "midgard"
	err = gsmobj.Create()
	if err != nil {
		return nil, fmt.Errorf("creating gsmaster record failed: %s", err)
	}
	return &gsmobj, nil
}

func TransformImportWeapon2GSDMaster(object *Waffe) (*gsmaster.Weapon, error) {
	gsmobj := gsmaster.Weapon{}

	err := gsmobj.First(object.Name)
	// if found check if we need to adjust masterdata
	if err == nil {
		return &gsmobj, nil
	}
	// if not found insert to masterdata
	gsmobj.System = "midgard"
	gsmobj.Name = object.Name
	gsmobj.Beschreibung = object.Beschreibung
	//gsmobj.Quelle = object.Quelle
	gsmobj.Gewicht = object.Gewicht
	gsmobj.Wert = object.Wert
	gsmobj.SkillRequired = "check"
	gsmobj.Damage = "check"
	err = gsmobj.Create()
	if err != nil {
		return nil, fmt.Errorf("creating gsmaster record failed: %s", err)
	}
	return &gsmobj, nil
}

func TransformImportContainer2GSDMaster(object *Behaeltniss) (*gsmaster.Container, error) {
	gsmobj := gsmaster.Container{}

	err := gsmobj.First(object.Name)
	// if found check if we need to adjust masterdata
	if err == nil {
		return &gsmobj, nil
	}
	// if not found insert to masterdata
	gsmobj.System = "midgard"
	gsmobj.Name = object.Name
	gsmobj.Beschreibung = object.Beschreibung
	//gsmobj.Quelle = object.Quelle
	gsmobj.Gewicht = object.Gewicht
	gsmobj.Wert = object.Wert
	gsmobj.Volumen = object.Volumen
	gsmobj.Tragkraft = object.Tragkraft
	err = gsmobj.Create()
	if err != nil {
		return nil, fmt.Errorf("creating gsmaster record failed: %s", err)
	}
	return &gsmobj, nil
}

func TransformImportTransportation2GSDMaster(object *Transportation) (*gsmaster.Transportation, error) {
	gsmobj := gsmaster.Transportation{}

	err := gsmobj.First(object.Name)
	// if found check if we need to adjust masterdata
	if err == nil {
		return &gsmobj, nil
	}
	// if not found insert to masterdata
	gsmobj.System = "midgard"
	gsmobj.Name = object.Name
	gsmobj.Beschreibung = object.Beschreibung
	//gsmobj.Quelle = object.Quelle
	gsmobj.Gewicht = float64(object.Gewicht)
	gsmobj.Wert = object.Wert
	gsmobj.Tragkraft = object.Tragkraft
	err = gsmobj.Create()
	if err != nil {
		return nil, fmt.Errorf("creating gsmaster record failed: %s", err)
	}
	return &gsmobj, nil
}

func TransformImportEquipment2GSDMaster(object *Ausruestung) (*gsmaster.Equipment, error) {
	gsmobj := gsmaster.Equipment{}

	err := gsmobj.First(object.Name)
	// if found check if we need to adjust masterdata
	if err == nil {
		return &gsmobj, nil
	}
	// if not found insert to masterdata
	gsmobj.System = "midgard"
	gsmobj.Name = object.Name
	gsmobj.Beschreibung = object.Beschreibung
	//gsmobj.Quelle = object.Quelle
	gsmobj.Gewicht = float64(object.Gewicht)
	gsmobj.Wert = object.Wert
	//gsmobj.Tragkraft = object.Tragkraft
	err = gsmobj.Create()
	if err != nil {
		return nil, fmt.Errorf("creating gsmaster record failed: %s", err)
	}
	return &gsmobj, nil
}

func TransformImportBelieve2GSDMaster(object string) (*gsmaster.Believe, error) {
	gsmobj := gsmaster.Believe{}

	err := gsmobj.First(object)
	// if found check if we need to adjust masterdata
	if err == nil {
		return &gsmobj, nil
	}
	// if not found insert to masterdata
	gsmobj.System = "midgard"
	gsmobj.Name = object
	gsmobj.Beschreibung = ""
	//gsmobj.Quelle = object.Quelle
	err = gsmobj.Create()
	if err != nil {
		return nil, fmt.Errorf("creating gsmaster record failed: %s", err)
	}
	return &gsmobj, nil
}

func CheckFertigkeiten2GSMaster(objects []Fertigkeit) error {
	for i := range objects {
		gsmobj, err := TransformImportFertigkeit2GSDMaster(&objects[i])
		if err != nil {
			return fmt.Errorf("creating gsmaster failed for 1 record: %s, %v", err, gsmobj)
		}
	}
	return nil
}
func CheckWaffenFertigkeiten2GSMaster(objects []Waffenfertigkeit) error {
	for i := range objects {
		gsmobj, err := TransformImportWaffenFertigkeit2GSDMaster(&objects[i])
		if err != nil {
			return fmt.Errorf("creating gsmaster failed for 1 record: %s, %v", err, gsmobj)
		}
	}
	return nil
}
func CheckSpells2GSMaster(objects []Zauber) error {
	for i := range objects {
		gsmobj, err := TransformImportSpell2GSDMaster(&objects[i])
		if err != nil {
			return fmt.Errorf("creating gsmaster failed for 1 record: %s, %v", err, gsmobj)
		}
	}
	return nil
}
func CheckWeapons2GSMaster(objects []Waffe) error {
	for i := range objects {
		gsmobj, err := TransformImportWeapon2GSDMaster(&objects[i])
		if err != nil {
			return fmt.Errorf("creating gsmaster failed for 1 record: %s, %v", err, gsmobj)
		}
	}
	return nil
}
func CheckContainers2GSMaster(objects []Behaeltniss) error {
	for i := range objects {
		gsmobj, err := TransformImportContainer2GSDMaster(&objects[i])
		if err != nil {
			return fmt.Errorf("creating gsmaster failed for 1 record: %s, %v", err, gsmobj)
		}
	}
	return nil
}
func CheckTransportationss2GSMaster(objects []Transportation) error {
	for i := range objects {
		gsmobj, err := TransformImportTransportation2GSDMaster(&objects[i])
		if err != nil {
			return fmt.Errorf("creating gsmaster failed for 1 record: %s, %v", err, gsmobj)
		}
	}
	return nil
}
func CheckEquipments2GSMaster(objects []Ausruestung) error {
	for i := range objects {
		gsmobj, err := TransformImportEquipment2GSDMaster(&objects[i])
		if err != nil {
			return fmt.Errorf("creating gsmaster failed for 1 record: %s, %v", err, gsmobj)
		}
	}
	return nil
}
func CheckBelieve2GSMaster(objects *CharacterImport) error {
	gsmobj, err := TransformImportBelieve2GSDMaster(objects.Glaube)
	if err != nil {
		return fmt.Errorf("creating gsmaster failed for 1 record: %s, %v", err, gsmobj)
	}
	return nil
}
