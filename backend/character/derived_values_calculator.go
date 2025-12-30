package character

import (
	"bamort/logger"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CalculateStaticFieldsRequest für Felder ohne Würfelwürfe
type CalculateStaticFieldsRequest struct {
	// Grundattribute
	St int `json:"st" binding:"required,min=1,max=100"` // Stärke
	Gs int `json:"gs" binding:"required,min=1,max=100"` // Geschicklichkeit
	Gw int `json:"gw" binding:"required,min=1,max=100"` // Gewandtheit
	Ko int `json:"ko" binding:"required,min=1,max=100"` // Konstitution
	In int `json:"in" binding:"required,min=1,max=100"` // Intelligenz
	Zt int `json:"zt" binding:"required,min=1,max=100"` // Zaubertalent
	Au int `json:"au" binding:"required,min=1,max=100"` // Aussehen

	// Charakterdaten
	Rasse string `json:"rasse" binding:"required"`
	Typ   string `json:"typ" binding:"required"`
	Grad  int    `json:"grad" binding:"omitempty,min=1,max=100"` // Charaktergrad (optional, default 1)
}

// StaticFieldsResponse für berechnete Felder ohne Würfelwürfe
type StaticFieldsResponse struct {
	AusdauerBonus         int `json:"ausdauer_bonus"`
	SchadensBonus         int `json:"schadens_bonus"`
	AngriffsBonus         int `json:"angriffs_bonus"`
	AbwehrBonus           int `json:"abwehr_bonus"`
	ZauberBonus           int `json:"zauber_bonus"`
	ResistenzBonusKoerper int `json:"resistenz_bonus_koerper"`
	ResistenzBonusGeist   int `json:"resistenz_bonus_geist"`
	ResistenzKoerper      int `json:"resistenz_koerper"`
	ResistenzGeist        int `json:"resistenz_geist"`
	Abwehr                int `json:"abwehr"`
	Zaubern               int `json:"zaubern"`
	Raufen                int `json:"raufen"`
}

// CalculateRolledFieldRequest für Felder mit Würfelwürfen
type CalculateRolledFieldRequest struct {
	// Grundattribute
	St int `json:"st" binding:"required,min=1,max=100"`
	Gs int `json:"gs" binding:"required,min=1,max=100"`
	Gw int `json:"gw" binding:"required,min=1,max=100"`
	Ko int `json:"ko" binding:"required,min=1,max=100"`
	In int `json:"in" binding:"required,min=1,max=100"`
	Zt int `json:"zt" binding:"required,min=1,max=100"`
	Au int `json:"au" binding:"required,min=1,max=100"`

	// Charakterdaten
	Rasse string `json:"rasse" binding:"required"`
	Typ   string `json:"typ" binding:"required"`
	Field string `json:"field" binding:"required"` // pa, wk, lp_max, ap_max, b_max

	// Würfelwerte vom Frontend
	Roll interface{} `json:"roll" binding:"required"` // Je nach Feld: int für 1d100, []int für mehrere Würfel
}

// RolledFieldResponse für Felder mit Würfelwürfen
type RolledFieldResponse struct {
	Field   string      `json:"field"`
	Value   int         `json:"value"`
	Formula string      `json:"formula"`
	Details interface{} `json:"details"`
}

// CalculateStaticFieldsLogic contains the pure business logic for static field calculation
func CalculateStaticFieldsLogic(req CalculateStaticFieldsRequest) StaticFieldsResponse {
	response := StaticFieldsResponse{}

	// Ausdauer Bonus: Ko/10 + St/20
	response.AusdauerBonus = (req.Ko / 10) + (req.St / 20)

	// Schadens Bonus: St/20 + Gs/30 - 3
	response.SchadensBonus = (req.St / 20) + (req.Gs / 30) - 3

	// Angriffs Bonus basierend auf GS
	response.AngriffsBonus = calculateAttributeBonus(req.Gs)

	// Abwehr Bonus basierend auf GW
	response.AbwehrBonus = calculateAttributeBonus(req.Gw)

	// Zauber Bonus basierend auf Zt
	response.ZauberBonus = calculateAttributeBonus(req.Zt)

	// Resistenz Bonus Körper
	response.ResistenzBonusKoerper = calculateResistenzBonusKoerper(req.Ko, req.Rasse, req.Typ)

	// Resistenz Bonus Geist
	response.ResistenzBonusGeist = calculateResistenzBonusGeist(req.In, req.Rasse, req.Typ)

	// Grad standardmäßig auf 1 setzen wenn nicht angegeben
	grad := req.Grad
	if grad < 1 {
		grad = 1
	}

	// Finale Resistenzwerte (mit Gradbonus)
	response.ResistenzKoerper = getResistenzBaseByGrade(grad) // + response.ResistenzBonusKoerper
	response.ResistenzGeist = getResistenzBaseByGrade(grad)   //+ response.ResistenzBonusGeist

	// Finale Kampfwerte (mit Gradbonus)
	response.Abwehr = getAbwehrBaseByGrade(grad)   //+ response.AbwehrBonus
	response.Zaubern = getZaubernBaseByGrade(grad) //+ response.ZauberBonus

	// Raufen: (St + GW)/20 + angriffs_bonus + Rassenboni
	raceBonus := 0
	if req.Rasse == "Zwerg" {
		raceBonus = 1
	}
	response.Raufen = (req.St+req.Gw)/20 + raceBonus //+ response.AngriffsBonus

	return response
}

// CalculateStaticFields berechnet alle Felder ohne Würfelwürfe (HTTP Handler)
func CalculateStaticFields(c *gin.Context) {
	var req CalculateStaticFieldsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Fehler beim Parsen der Static Fields Anfrage: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültige Anfrage"})
		return
	}

	logger.Info("Berechne statische Felder für %s %s", req.Rasse, req.Typ)

	response := CalculateStaticFieldsLogic(req)

	c.JSON(http.StatusOK, response)
}

// CalculateRolledField berechnet ein Feld mit Würfelwurf
func CalculateRolledField(c *gin.Context) {
	var req CalculateRolledFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Fehler beim Parsen der Rolled Field Anfrage: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültige Anfrage"})
		return
	}

	logger.Info("Berechne Würfelfeld %s für %s %s", req.Field, req.Rasse, req.Typ)

	response, err := calculateRolledField(req)
	if err != nil {
		logger.Error("Fehler beim Berechnen des Würfelfeldes %s: %v", req.Field, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// calculateRolledField führt die Berechnung für Würfelfelder durch
func calculateRolledField(req CalculateRolledFieldRequest) (*RolledFieldResponse, error) {
	field := strings.ToLower(req.Field)

	switch field {
	case "pa":
		roll, ok := req.Roll.(float64) // JSON numbers become float64
		if !ok {
			return nil, fmt.Errorf("ungültiger Würfelwert für PA")
		}
		rollInt := int(roll)
		modifier := (4 * req.In / 10) - 20
		value := rollInt + modifier

		return &RolledFieldResponse{
			Field:   req.Field,
			Value:   value,
			Formula: "1d100 + 4×In/10 - 20",
			Details: map[string]interface{}{
				"roll":          rollInt,
				"in_bonus":      4 * req.In / 10,
				"base_modifier": -20,
				"modifier":      modifier,
			},
		}, nil

	case "wk":
		roll, ok := req.Roll.(float64)
		if !ok {
			return nil, fmt.Errorf("ungültiger Würfelwert für WK")
		}
		rollInt := int(roll)
		modifier := 2*((req.Ko/10)+(req.In/10)) - 20
		value := rollInt + modifier

		return &RolledFieldResponse{
			Field:   req.Field,
			Value:   value,
			Formula: "1d100 + 2×(Ko/10 + In/10) - 20",
			Details: map[string]interface{}{
				"roll":          rollInt,
				"ko_bonus":      req.Ko / 10,
				"in_bonus":      req.In / 10,
				"base_modifier": -20,
				"modifier":      modifier,
			},
		}, nil

	case "lp_max":
		roll, ok := req.Roll.(float64)
		if !ok {
			return nil, fmt.Errorf("ungültiger Würfelwert für LP Max")
		}
		rollInt := int(roll)
		raceModifier := getRaceModifierLP(req.Rasse)
		value := rollInt + 7 + req.Ko/10 + raceModifier

		return &RolledFieldResponse{
			Field:   req.Field,
			Value:   value,
			Formula: "1d3 + 7 + Ko/10 + Rassenmodifikator",
			Details: map[string]interface{}{
				"roll":          rollInt,
				"base":          7,
				"ko_bonus":      req.Ko / 10,
				"race_modifier": raceModifier,
			},
		}, nil

	case "ap_max":
		roll, ok := req.Roll.(float64)
		if !ok {
			return nil, fmt.Errorf("ungültiger Würfelwert für AP Max")
		}
		rollInt := int(roll)
		ausdauerBonus := (req.Ko / 10) + (req.St / 20)
		classModifier := getClassModifierAP(req.Typ)
		value := rollInt + 1 + ausdauerBonus + classModifier

		return &RolledFieldResponse{
			Field:   req.Field,
			Value:   value,
			Formula: "1d3 + 1 + Ausdauerbonus + Klassenmodifikator",
			Details: map[string]interface{}{
				"roll":           rollInt,
				"base":           1,
				"ausdauer_bonus": ausdauerBonus,
				"class_modifier": classModifier,
			},
		}, nil

	case "b_max":
		// Erwarte Array von Würfelwerten
		rollsInterface, ok := req.Roll.([]interface{})
		if !ok {
			return nil, fmt.Errorf("ungültiger Würfelwert für B Max - Array erwartet")
		}

		rolls := make([]int, len(rollsInterface))
		total := 0
		for i, r := range rollsInterface {
			if rollFloat, ok := r.(float64); ok {
				rolls[i] = int(rollFloat)
				total += rolls[i]
			} else {
				return nil, fmt.Errorf("ungültiger Würfelwert in B Max Array")
			}
		}

		baseValue, formula := getMovementBaseAndFormula(req.Rasse)
		value := total + baseValue

		return &RolledFieldResponse{
			Field:   req.Field,
			Value:   value,
			Formula: formula,
			Details: map[string]interface{}{
				"rolls":      rolls,
				"roll_total": total,
				"base":       baseValue,
				"race":       req.Rasse,
			},
		}, nil

	default:
		return nil, fmt.Errorf("unbekanntes Würfelfeld: %s", req.Field)
	}
}

// Hilfsfunktionen

// calculateAttributeBonus berechnet Attributboni basierend auf dem Wert
func calculateAttributeBonus(value int) int {
	if value >= 1 && value <= 5 {
		return -2
	} else if value >= 6 && value <= 20 {
		return -1
	} else if value >= 21 && value <= 80 {
		return 0
	} else if value >= 81 && value <= 95 {
		return 1
	} else if value >= 96 && value <= 100 {
		return 2
	}
	return 0
}

// calculateResistenzBonusKoerper berechnet den Körper-Resistenzbonus
func calculateResistenzBonusKoerper(ko int, rasse string, typ string) int {
	bonus := calculateAttributeBonus(ko)

	switch rasse {
	//case "Mensch":
	// bonus = calculateAttributeBonus(ko)
	case "Elf":
		bonus = 2
	case "Gnom", "Halbling":
		bonus = 4
	case "Zwerg":
		bonus = 3
	}

	// Klassenmodifikator
	if isKaempfer(typ) {
		bonus += 1
	} else if isZauberer(typ) {
		bonus += 2
	}

	return bonus
}

// calculateResistenzBonusGeist berechnet den Geist-Resistenzbonus
func calculateResistenzBonusGeist(in int, rasse string, typ string) int {
	bonus := calculateAttributeBonus(in)

	switch rasse {
	//case "Mensch":
	//bonus = calculateAttributeBonus(in)
	case "Elf":
		bonus = 2
	case "Gnom", "Halbling":
		bonus = 4
	case "Zwerg":
		bonus = 3
	}

	// Klassenmodifikator (nur Zauberer bekommen Geist-Bonus)
	if isZauberer(typ) {
		bonus += 2
	}

	return bonus
}

// isKaempfer prüft ob eine Klasse als Kämpfer gilt
func isKaempfer(typ string) bool {
	kaempferKlassen := []string{"Barbar", "Krieger", "Waldläufer", "Assassine", "Spitzbube"}
	for _, k := range kaempferKlassen {
		if k == typ {
			return true
		}
	}
	return false
}

// isZauberer prüft ob eine Klasse als Zauberer gilt
func isZauberer(typ string) bool {
	zaubererKlassen := []string{"Magier", "Druide", "Priester", "Schamane", "Hexer"}
	for _, z := range zaubererKlassen {
		if z == typ {
			return true
		}
	}
	return false
}

// getRaceModifierLP gibt den LP-Modifikator für eine Rasse zurück
func getRaceModifierLP(rasse string) int {
	switch rasse {
	case "Gnome":
		return -3
	case "Halblinge":
		return -2
	case "Zwerge":
		return 1
	default:
		return 0
	}
}

// getClassModifierAP gibt den AP-Modifikator für eine Klasse zurück
func getClassModifierAP(typ string) int {
	switch typ {
	case "Barbar", "Krieger", "Waldläufer":
		return 2
	case "Assassine", "Spitzbube", "Schamane":
		return 1
	default:
		return 0
	}
}

// getMovementBaseAndFormula gibt den Basiswert und die Formel für Bewegung zurück
func getMovementBaseAndFormula(rasse string) (int, string) {
	switch rasse {
	case "Gnome", "Halblinge":
		return 8, "2d3 + 8"
	case "Zwerge":
		return 12, "3d3 + 12"
	default: // Menschen, Elfen
		return 16, "4d3 + 16"
	}
}

// getAbwehrBaseByGrade gibt den Basis-Fertigkeitswert für Abwehr nach Grad zurück
// Grad: 1->11, 2->12, 5->13, 10->14, 15->15, 20->16, 25->17, 30->18
func getAbwehrBaseByGrade(grad int) int {
	if grad >= 30 {
		return 18
	} else if grad >= 25 {
		return 17
	} else if grad >= 20 {
		return 16
	} else if grad >= 15 {
		return 15
	} else if grad >= 10 {
		return 14
	} else if grad >= 5 {
		return 13
	} else if grad >= 2 {
		return 12
	}
	return 11
}

// getResistenzBaseByGrade gibt den Basis-Fertigkeitswert für Resistenz nach Grad zurück
// Verwendet dieselbe Tabelle wie Abwehr
func getResistenzBaseByGrade(grad int) int {
	return getAbwehrBaseByGrade(grad)
}

// getZaubernBaseByGrade gibt den Basis-Fertigkeitswert für Zaubern nach Grad zurück
// Grad: 1->11, 2->12, 4->13, 6->14, 8->15, 10->16, 15->17, 20->18
func getZaubernBaseByGrade(grad int) int {
	if grad >= 20 {
		return 18
	} else if grad >= 15 {
		return 17
	} else if grad >= 10 {
		return 16
	} else if grad >= 8 {
		return 15
	} else if grad >= 6 {
		return 14
	} else if grad >= 4 {
		return 13
	} else if grad >= 2 {
		return 12
	}
	return 11
}
