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
	Name       string     // z.B. "Überleben"
	Group      SkillGroup // z.B. GroupFreiland
	Difficulty Difficulty // z.B. DiffLeicht
}

// BaseLearnCost speichert die Anzahl der Lerneinheiten (LE),
// die nötig sind, um eine Fertigkeit erstmals zu lernen.
var BaseLearnCost = map[SkillGroup]map[Difficulty]int{
	GroupAlltag: {
		DiffLeicht:     1,
		DiffNormal:     1,
		DiffSchwer:     2,
		DiffSehrSchwer: 10,
	},
	GroupFreiland: {
		// Werte laut deiner Freiland-Tabelle:
		DiffLeicht: 1,
		DiffNormal: 1,
		DiffSchwer: 2,
		// Für sehr_schwer ist nichts vorgegeben, ggf. 0 oder ein eigener Wert:
		DiffSehrSchwer: 10,
	},
	GroupKampf: {
		// Beispielwerte, ggf. anpassen
		DiffLeicht:     2,
		DiffNormal:     2,
		DiffSchwer:     4,
		DiffSehrSchwer: 12,
	},
	// weitere Gruppen ...
}

// ImprovementCost speichert die nötigen Lerneinheiten (LE),
// um eine Fertigkeit von aktuellem Bonus X auf X+1 zu verbessern.
var ImprovementCost = map[Difficulty]map[int]int{
	// Alltag-Werte beispielhaft
	DiffLeicht: {
		// +9, +10, +11, +12 => 0 bedeutet „-“ oder kein Aufwand angegeben
		9:  0,
		10: 0,
		11: 0,
		12: 0,
		13: 1,
		14: 1,
		15: 1,
		// usw.
	},
	DiffNormal: {
		9:  1,
		10: 1,
		11: 1,
		12: 1,
		13: 2,
		14: 2,
		15: 2,
		// usw.
	},
	DiffSchwer: {
		9:  2,
		10: 2,
		11: 5,
		12: 5,
		13: 10,
		// usw.
	},
	DiffSehrSchwer: {
		9:  5,
		10: 5,
		11: 10,
		12: 10,
		13: 20,
		// usw.
	},
}

// Zusätzlich Beispielwerte speziell für Freiland-Verbesserungen
// Falls Freiland sich unterscheidet, könnte man entweder
// a) pro Gruppe und Difficulty separate Tabellen führen
// oder b) im obigen ImprovementCost-Feld differenzieren.
// Hier zeigen wir, wie es aussehen könnte, wenn Freiland abweicht:
var FreilandImprovementCost = map[Difficulty]map[int]int{
	DiffLeicht: {
		9:  1,
		10: 1,
		11: 1,
		12: 2,
		13: 2,
	},
	DiffNormal: {
		9:  2,
		10: 5,
		11: 5,
		12: 10,
		13: 15,
	},
	DiffSchwer: {
		9:  5,
		10: 5,
		11: 10,
		12: 10,
		13: 20,
	},
}

// CharClass kennzeichnet eine Charakterklasse für die EP-Berechnung
type CharClass string

const (
	ClassKrieger CharClass = "Krieger"
	ClassMagier  CharClass = "Magier"
	ClassSchurke CharClass = "Schurke"
	// weitere Klassen ...
)

// Wieviele EP kostet 1 Lerneinheit (LE) je Klasse & SkillGroup?
// "-" aus deiner Beispiel-Tabelle kann man als 0 interpretieren (nicht erlernbar),
// oder man trägt in einer separaten Tabelle AllowedGroups die Erlaubnis ein.
var EPPerLE = map[CharClass]map[SkillGroup]int{
	ClassKrieger: {
		GroupAlltag:   20,
		GroupFreiland: 20,
		GroupKampf:    10,
	},
	ClassMagier: {
		GroupAlltag:   30,
		GroupFreiland: 30,
		// kein Kampf -> 0 oder nicht erlernbar
	},
	ClassSchurke: {
		GroupAlltag:   10,
		GroupFreiland: 30,
		GroupKampf:    30,
	},
}

// Optional: Erlaubnis pro Klasse/Gruppe
var AllowedGroups = map[CharClass]map[SkillGroup]bool{
	ClassKrieger: {
		GroupAlltag:   true,
		GroupFreiland: true,
		GroupKampf:    true,
	},
	ClassMagier: {
		GroupAlltag:   true,
		GroupFreiland: true,
		GroupKampf:    false,
	},
	ClassSchurke: {
		GroupAlltag:   true,
		GroupFreiland: true,
		GroupKampf:    true,
	},
}

// CalculateLearnCost berechnet die EP-Kosten für das erstmalige Lernen einer Fertigkeit.
func CalculateLearnCost(skill SkillDefinition, class CharClass) (int, error) {
	// Darf die Klasse diese Gruppe überhaupt lernen?
	if !AllowedGroups[class][skill.Group] {
		return 0, fmt.Errorf("Gruppe %s nicht erlernbar für %s", skill.Group, class)
	}
	groupCosts, ok := BaseLearnCost[skill.Group]
	if !ok {
		return 0, errors.New("unbekannte Skill-Gruppe")
	}
	baseLE, ok := groupCosts[skill.Difficulty]
	if !ok {
		return 0, errors.New("keine LearnCost-Definition für diese Schwierigkeit")
	}
	groupEP, ok := EPPerLE[class][skill.Group]
	if !ok {
		return 0, errors.New("keine EP-Kosten für diese Gruppe & Klasse definiert")
	}
	// Gesamt-EP = Anzahl Lerneinheiten * EP pro Lerneinheit
	return baseLE * groupEP, nil
}

// CalculateImprovementCost berechnet die EP-Kosten, um eine bestehende Fertigkeit
// von currentBonus auf currentBonus+1 zu steigern.
// Falls für Freiland separate Werte gelten, nutze FreilandImprovementCost.
// Alternativ kann man bei ImprovementCost weitere Gruppendifferenzierung hineinbauen.
func CalculateImprovementCost(skill SkillDefinition, class CharClass, currentBonus int) (int, error) {
	if !AllowedGroups[class][skill.Group] {
		return 0, fmt.Errorf("Gruppe %s nicht erlernbar für %s", skill.Group, class)
	}

	var neededLE int
	if skill.Group == GroupFreiland {
		// Verwende die spezielle Freiland-Tabelle
		diffMap, ok := FreilandImprovementCost[skill.Difficulty]
		if !ok {
			return 0, errors.New("keine Freiland-Kosten für diese Schwierigkeit")
		}
		val, found := diffMap[currentBonus+1]
		if !found {
			return 0, fmt.Errorf("kein Eintrag für Bonus %d→%d", currentBonus, currentBonus+1)
		}
		neededLE = val
	} else {
		// Standard-Tabelle
		diffMap, ok := ImprovementCost[skill.Difficulty]
		if !ok {
			return 0, errors.New("keine Improvement-Kosten für diese Schwierigkeit")
		}
		val, found := diffMap[currentBonus+1]
		if !found {
			return 0, fmt.Errorf("kein Eintrag für Bonus %d→%d", currentBonus, currentBonus+1)
		}
		neededLE = val
	}

	epPerLE, ok := EPPerLE[class][skill.Group]
	if !ok {
		return 0, errors.New("keine EP-Kosten für diese Gruppe & Klasse definiert")
	}
	return neededLE * epPerLE, nil
}
