package gsmaster

/*
Fertigkeiten werden in Fertigkeitsgruppen eingeteilt zum Beispiel Alltag, Freiland und Kampf
innerhalb der Gruppe saind die Fertigkeiten eingeordnet nach Schwierigkeit. Also leicht, normal, schwer und sehr schwer.
Je schwerer eine Fertigkeit zu lernen ist um so mehr Lernpunkte müssen führ ihre Verbesserung investiert werden.
Auch pro Verbesserungsstufe sind je nach Schwierigkeit unterschiedliche viele Lernpunkte zu investieren.

Beispiel
Tabelle Lern und Verbesserungskosten:
Alltag
Lernen:
leicht: 1LE, {Bootfahren, Kochen, Reiten, Schwimmen, Tanzen, ...}
normal: 1LE, {schreiben, Sprache, Lesen, ...}
schwer: 2LE, {Erste Hilfe, Etikette, Fälschen, ...}
sehr schwer: 10LE, {Gerätekunde, Geschäftssinn, ...}
Verbessern:

	+9;		+10;	+11;	+12;	+13;...

leicht: 	-;		-;		-;		-;		1;...
normal: 	1;		1;		1;		1;		2;...
schwer: 	2;		2;		5;		5;		10;...
sehr schwer: 5;		5;		10;		10;		20;...

Freiland
Lernen:
leicht: 1LE, {Überleben, ...}
normal: 1LE, {Naturkunde, ...}
schwer: 2LE, {Tarnen, ...}
Verbessern:

	+9;		+10;	+11;	+12;	+13;...

leicht: 	1;		1;		1;		2;		2;...
normal: 	2;		5;		5;		10;		15;...
schwer: 	5;		5;		10;		10;		20;...

Tabelle EP Kosten pro Lerneeinheit

	Alltag;	Freiland;	Kampf;	...

Krieger		20;		20;			10;	...
Magier		30;		30;			-;	...
Schurke		10;		30;			30;	...
...

Die pro Lerneinheit aufzuwendenden Erfahrungspunkte sind für die einzelnen Fertigkeitsgruppen pro charakterklasse unterschiedlich. Für einige Charakterklassen sind einzelne Fertigkeitsgruppen nicht erlernbar.

Erstelle aus diesen Informationen eine Datenstruktur, die es ermöglicht die Lernpunkte für eine Fertigkeit zu berechnen.
Ziehe dazu die Dateien model.go in den Verzeichnissen backend/gsmaster, backend/skills, und backend/character zu Rate.
*/
type SkillGroup string

const (
	GroupAlltag   SkillGroup = "Alltag"
	GroupFreiland SkillGroup = "Freiland"
	GroupKampf    SkillGroup = "Kampf"
	// ... add others as needed
)

type Difficulty string

const (
	DifficultyLight     Difficulty = "leicht"
	DifficultyNormal    Difficulty = "normal"
	DifficultyHeavy     Difficulty = "schwer"
	DifficultyVeryHeavy Difficulty = "sehr_schwer"
)

// Base cost per difficulty (example values)
var DifficultyBaseCost = map[Difficulty]int{
	DifficultyLight:     1,
	DifficultyNormal:    2,
	DifficultyHeavy:     3,
	DifficultyVeryHeavy: 5,
}

// Character class multipliers by skill group (example values)
var SkillGroupMultiplierByClass = map[string]map[SkillGroup]float64{
	"Krieger": {
		GroupAlltag:   1.0,
		GroupFreiland: 1.5,
		GroupKampf:    1.0,
	},
	"Magier": {
		GroupAlltag:   1.2,
		GroupFreiland: 2.0,
		GroupKampf:    99.0, // not learnable or extremely expensive
	},
	// ... add more classes here
}

// SkillDefinition captures group & difficulty
type SkillDefinition struct {
	Name       string
	Group      SkillGroup
	Difficulty Difficulty
}

// CalculateSkillCost calculates the cost for one improvement step
func CalculateSkillCost(skill SkillDefinition, charClass string) int {
	base := DifficultyBaseCost[skill.Difficulty]
	classMultiplier := 1.0

	if multipliers, ok := SkillGroupMultiplierByClass[charClass]; ok {
		if mult, found := multipliers[skill.Group]; found {
			classMultiplier = mult
		} else {
			// If not found, treat as unlearnable or very high cost
			classMultiplier = 99.0
		}
	}

	return int(float64(base) * classMultiplier)
}
