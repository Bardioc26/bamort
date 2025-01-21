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

import (
	"errors"
	"fmt"
)

// SkillGroup kennzeichnet die übergeordnete Gruppe einer Fertigkeit
type SkillGroup string

const (
	GroupAlltag   SkillGroup = "Alltag"
	GroupFreiland SkillGroup = "Freiland"
	GroupKampf    SkillGroup = "Kampf"
	// weitere Gruppen ...
)

// Difficulty gibt an, wie schwer eine Fertigkeit ist
type Difficulty string

const (
	DiffLeicht     Difficulty = "leicht"
	DiffNormal     Difficulty = "normal"
	DiffSchwer     Difficulty = "schwer"
	DiffSehrSchwer Difficulty = "sehr_schwer"
)

// SkillDefinition definiert eine einzelne Fertigkeit
type SkillDefinition struct {
	Name       string     // z.B. "Schwimmen"
	Group      SkillGroup // z.B. GroupAlltag = "Alltag"
	Difficulty Difficulty // z.B. DiffLeicht = "leicht"
	// weitere Felder aus model.go, z.B. Bonus, AP-Kosten etc.
}

// Kosten für das erstmalige Lernen pro Schwierigkeitsgrad in Lerneinheiten (LE)
var BaseLearnCost = map[SkillGroup]map[Difficulty]int{
	GroupAlltag: {
		DiffLeicht:     1,
		DiffNormal:     1,
		DiffSchwer:     2,
		DiffSehrSchwer: 10,
	},
	GroupFreiland: {
		// Beispielwerte
		DiffLeicht:     2,
		DiffNormal:     3,
		DiffSchwer:     5,
		DiffSehrSchwer: 10,
	},
	GroupKampf: {
		// Beispielwerte
		DiffLeicht:     2,
		DiffNormal:     2,
		DiffSchwer:     4,
		DiffSehrSchwer: 12,
	},
}

// Verbesserungs-Kosten in Lerneinheiten (LE) pro Schwierigkeitsgrad.
// Hier ein Beispiel-Layout: Map von "Fertigkeit +X" -> "Wie viele Lerneinheiten nötig"
var ImprovementCost = map[Difficulty]map[int]int{
	// Ein Beispiel: "leicht" => ab +13 kostet 1 LE usw.
	DiffLeicht: {
		9:  0, // "-" in der Tabelle => 0 oder man ignoriert’s
		10: 0,
		11: 0,
		12: 0,
		13: 1,
		14: 1,
		15: 1,
		// ...
	},
	DiffNormal: {
		9:  1,
		10: 1,
		11: 1,
		12: 1,
		13: 2,
		14: 2,
		15: 2,
		// ...
	},
	DiffSchwer: {
		9:  2,
		10: 2,
		11: 5,
		12: 5,
		13: 10,
		// ...
	},
	DiffSehrSchwer: {
		9:  5,
		10: 5,
		11: 10,
		12: 10,
		13: 20,
		// ...
	},
}

// Charakterklassen mit EP-Kosten pro SkillGroup (Lerneinheit)
type CharClass string

const (
	ClassKrieger CharClass = "Krieger"
	ClassMagier  CharClass = "Magier"
	ClassSchurke CharClass = "Schurke"
	// weitere Klassen ...
)

// EP-Kosten pro Lerneinheit je Klasse und Gruppe
// "-" (nicht erlernbar) wird hier als 0 interpretiert + Sperre in AllowedGroups
var EPPerLE = map[CharClass]map[SkillGroup]int{
	ClassKrieger: {
		GroupAlltag:   20,
		GroupFreiland: 20,
		GroupKampf:    10,
	},
	ClassMagier: {
		GroupAlltag:   30,
		GroupFreiland: 30,
		// GroupKampf: 0 => nicht erlernbar s.u.
	},
	ClassSchurke: {
		GroupAlltag:   10,
		GroupFreiland: 30,
		GroupKampf:    30,
	},
}

// AllowedGroups gibt an, ob eine Klasse eine bestimmte Gruppe lernen darf.
// Falls false oder nicht vorhanden, ist diese Gruppe für die Klasse gesperrt.
var AllowedGroups = map[CharClass]map[SkillGroup]bool{
	ClassKrieger: {
		GroupAlltag:   true,
		GroupFreiland: true,
		GroupKampf:    true,
	},
	ClassMagier: {
		GroupAlltag:   true,
		GroupFreiland: true,
		GroupKampf:    false, // "-" in Tabelle
	},
	ClassSchurke: {
		GroupAlltag:   true,
		GroupFreiland: true,
		GroupKampf:    true,
	},
}

// CalculateLearnCost ermittelt, wie viele Lerneinheiten für das erstmalige Lernen notwendig sind
// und multipliziert ihn mit den EP/Kosten pro LE für die gegebene Klasse
func CalculateLearnCost(skill SkillDefinition, class CharClass) (int, error) {
	// Ist die Gruppe für diese Klasse erlaubt?
	if !AllowedGroups[class][skill.Group] {
		return 0, fmt.Errorf("skill group %s cannot be learned by class %s", skill.Group, class)
	}
	grpMap, ok := BaseLearnCost[skill.Group]
	if !ok {
		return 0, errors.New("unknown skill group")
	}
	baseLE, ok := grpMap[skill.Difficulty]
	if !ok {
		return 0, errors.New("unknown difficulty in base cost table")
	}
	epTable, ok := EPPerLE[class]
	if !ok {
		return 0, errors.New("unknown class in EP table")
	}
	epPerLE, found := epTable[skill.Group]
	if !found {
		return 0, errors.New("no EP cost defined for group in class")
	}

	// Gesamt-EP = Lerneinheiten * EP pro LE
	totalEP := baseLE * epPerLE
	return totalEP, nil
}

// CalculateImprovementCost berechnet die EP-Kosten für die nächste Verbesserung (z.B. von +12 auf +13).
// 'currentBonus' ist der aktuelle Bonus (z.B. +12).
func CalculateImprovementCost(skill SkillDefinition, class CharClass, currentBonus int) (int, error) {
	// Gruppe erlaubt?
	if !AllowedGroups[class][skill.Group] {
		return 0, fmt.Errorf("skill group %s cannot be learned by class %s", skill.Group, class)
	}
	// Passende ImprovementCost-Map
	diffMap, ok := ImprovementCost[skill.Difficulty]
	if !ok {
		return 0, errors.New("unknown difficulty in improvement cost table")
	}
	neededLE, ok := diffMap[currentBonus+1] // z.B. +12 -> +13 => currentBonus+1
	if !ok {
		return 0, fmt.Errorf("no improvement cost for skill level %d", currentBonus+1)
	}
	epTable, ok := EPPerLE[class]
	if !ok {
		return 0, errors.New("unknown class in EP table")
	}
	epPerLE, found := epTable[skill.Group]
	if !found {
		return 0, errors.New("no EP cost defined for group in class")
	}

	totalEP := neededLE * epPerLE
	return totalEP, nil
}
