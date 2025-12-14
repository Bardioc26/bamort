# Aufteilung der Informationen auf dem Characterbogen

Der Characterbogen besteht aus mehreren Seiten. 
- Seite 1 ist die Statistikseite in der alle Werte in ihrem Grundwert eingetraqgen werden also ohne irgendwelche Boni.
- Seite 2 ist der Spielbogen. sie enthält alle werte so wie sie für das Spiel benötigt werden Boni und Mali jedweder Art sind hier in die Werte der Fertigkeiten eingerechnet. Das Gilt auch für magische Boni für Waffen, Rüstung oder ähnliches.
- Seite 3 enhält die Liste aller gelernten Zauber und aller Zauber die dem Charakter durch Artefakte zur Verfügung stehen
- Seite 4 enthält die Ausrüstung des Abenteurers. Hier ist Platz für "Am Körper getragen", "Auf dem Wagen", "Im Behältnis XYZ"

## Allgemeine Layoutwerte

- Seitengröße A4,quer: 297×210 mm
- Ränder
    - oben: 6 mm
    - unten: 6 mm
    - links: 6 mm
    - rechts: 6 mm
- Spaltenraster
    - Anzahl: 4
    - Breite: alle gleich
    - Gutter: 4 mm
- Font
    - Name: DejaVu Sans
    - size
        - default: 16 
        - small: 14
        - large: 18, 
    - lineHeightFactor
        - default: 1.2
        - small: 1.1
        - large: 1.3
    - spacing
        - vertical: 0.5 mm
        - horizontal: 0.5 mm
- Boxtypen
    - id: 01
        - name: header
        - f-Size: default
        - type: Label: Wert
    - id: 02
        - name: stats
        - f-Size: large
        - type: Label: Wert
        - style: boxed
    - id: 03
        - name: table
        - f-Size: default
        - type: Label: Wert
        - style: nobr, shorten

Wenn anchor für Blöcke nicht angegeben ist wird als default  "top" verwendet.
Wenn heightMode für Blöcke nicht angegeben ist wird als default "fixed" verwendet.

- In jedem Spaltenbereich:
  - zuerst alle anchor=bottom-Blöcke von unten platzieren,
  - dann alle anchor=top-Blöcke von oben,
  - für heightMode=rows-Tabellen:
  - maxRows = floor(availableHeight / rowHeight).

Rowheight wird wie folgt berechnet:
rowHeight_mm = fontSize_pt * lineHeightFactor * ptToMm + spacing.vertical


## Strukturierung der Seiten A4 Quer

Jede seite ist prinzipiell 4 spaltig angelegt. Die Informationen sind in Blöcken organisiert die bei Bedarf mehrere Spalten zusammen fassen können. Wenn nich anders festgelegt laufen die Informationen in den Blöcken von links oben nach rechts unten. „anders festgelegt“ bedeutet:
- anchor=bottom ist eine Ausnahme,
- ansonsten ist die Reihenfolge in der Seitendefinition die Flow-Reihenfolge.

## Seite 1 (Statistikseite)

- Block 1, Spalte 1,2,3,4, Darstellung Label: Wert
    - Name des Charakters 
    - Name des Spielers
- Block 2, Spalte 1, Darstellung Label: Wert
    - Typ, Grad
    - Spezialisierung
- Block 3, Spalte 1, Darstellung Box Label: Wert
    - St, Gs, Gw
    - Ko, In, Zt
    - Au, pA, Wk
    - B, Raufen
- Block 4, Spalte 1, Darstellung Label: Wert
    - "persönliche Boni für:"
    - Ausdauer, Schaden, Angriff
    - Abweht, Resistenz, Zaubern
- Block 5, Spalte 1,2, Darstellung Box Label: Wert
    - LP-Max, AP-Max, GG, SG
- Block 6, Spalte 1,2, Darstellung Tabelle
    - Tabelle Spielergebnisse
- Block 7, Spalte 2, Darstellung Label: Wert
    - Geburtsdatum
    - Alter Händigkeit
    - Größe
    - Gestalt, Gewicht
    - Stand
    - Heimat
    - Glaube
    - Besondere Merkmale
- Block 8, Spalte 3,4
    - "Liste der gelernten und angeborenen Fertigkeiten"
- Block 9, Spalte 3, Darstellung Tabelle, heightMode: rows
    - Tabelle Fertigkeiten
- Block 10, Spalte 3,4, Darstellung Label: Wert, anchor: bottom
    - Sehen, Nachtsicht, Hören, Riechen/Schmecken, Sechster Sinn
- Block 11, Spalte 4, Darstellung Tabelle, heightMode: rows
    - Fortsetzung Tabelle Fertigkeiten


## Seite 2 (Spielbogen).

- Block 1, Spalte 1,2,3,4, Darstellung Label: Wert
    - Name des Charakters, Grad 
    - Typ, GG, SG
- Block 2, Spalte 1,2, Darstellung Box Label: Wert
    - St, Gs, Gw, Ko, In, 
    - Zt, Au, pA, Wk, B
- Block 3, Spalte 1,2, Darstellung Label: Wert
    - Abwehr, Abwehr mit Verteidigungswaffe, Resistenz, Zaubern
- Block 4, Spalte 1, Darstellung Tabelle, heightMode: rows
    - Tabelle Fertigkeiten
- Block 5, Spalte 2, Darstellung Tabelle
    - Tabelle ungelernte Fertigkeiten
- Block 6, Spalte 2, Darstellung Tabelle, heightMode: rows
    - Fortsetzung Tabelle Fertigkeiten
- Block 7, Spalte 3,4, Darstellung Tabelle
    - Tabelle LP/AP
- Block 8, Spalte 3, Darstellung Label: Wert
    - "Bonus für"
    - "Schaden", "Angriff", "Abwehr"
    - Schaden, Angriff, Abwehr
    - "Mit Rüstung", Angriff, Abwehr
- Block 9 ,Spalte 3, Darstellung Tabelle, heightMode: rows
    - Tabelle Waffen
- Block 10, Spalte 3,4, Darstellung Label: Wert, anchor: bottom
    - Sehen, Nachtsicht, Hören, Riechen/Schmecken, Sechster Sinn
- Block 11, Spalte 4, Darstellung Label: Wert
    - "RK", rk-lp-bonus
    - , mit Rüstung", B, Gw
- Block 12 ,Spalte 4, Darstellung Tabelle, heightMode: rows
    - Fortsetzung Tabelle Waffen
- Block 13, Spalte 4, Darstellung Tabelle
    - Tabelle PP Zauber

## Seite 3 (Zauber )

- Block 1, Spalte 1,2,3,4, Darstellung Label: Wert
    - Name des Charakters, Grad 
    - Typ, Zaubern
- Block 2. Spalte 1,2, Darstellung Tabelle, heightMode: rows
    - Tabelle Zauber
- Block 3. Spalte 3,4, Darstellung Tabelle, heightMode: rows
    - Fortsetzung Tabelle Zauber
- Block 4, Spalte 3,4, Darstellung Tabelle, heightMode: rows
    - Tabelle magische Gegenstände /Ausrüstung

## Seite 4 (Ausrüstung)
- Block 1, Spalte 1,2,3,4, Darstellung Label: Wert
    - Name des Charakters, Grad 
    - Typ
- Block 2. Spalte 1,2, Darstellung Tabelle, heightMode: rows
    - Tabelle Ausrüstung (Am Körper getragen)
- Block 3. Spalte 1,2, Darstellung Tabelle, heightMode: rows
    - Tabelle Behälter 1 (Im Wagen)
- Block 4. Spalte 3,4, Darstellung Tabelle, heightMode: rows
    - Tabelle Behälter 2
- Block 5. Spalte 3,4, Darstellung Tabelle, heightMode: rows
    - Tabelle Behälter 2


## Seite 2.1 (Fortsetzung Spielbogen).

Wenn die Anzahl der Fertigkeiten oder die Anzahl der Waffen den Platz in der Liste überschreitet wird eine oder mehrere Seiten angefügt. Diese enthält:

- Block 1, Spalte 1,2,3,4, Darstellung Label: Wert
    - Name des Charakters, Grad 
    - Typ, GG, SG
- Block 2, Spalte 1, Darstellung Tabelle, heightMode: rows
    - Fortsetzung Tabelle Fertigkeiten
- Block 3, Spalte 2, Darstellung Tabelle, heightMode: rows
    - Fortsetzung Tabelle Fertigkeiten
- Block 4 ,Spalte 3, Darstellung Tabelle, heightMode: rows
    - Fortsetzung Tabelle Fertigkeiten
- Block 5 ,Spalte 4, Darstellung Tabelle, heightMode: rows
    - Fortsetzung Tabelle Waffen

## Seite 3.1 (Forsetung Zauber)

Wenn die Anzahl der Zauber den Platz in der Liste überschreitet wird eine oder mehrere identische Seiten angefügt, Block 4 wird dann auf die von Seite 3 auf Seite 3.1 verschoben.

- Block 1, Spalte 1,2,3,4, Darstellung Label: Wert
    - Name des Charakters, Grad 
    - Typ, Zaubern
- Block 2. Spalte 1,2, Darstellung Tabelle, heightMode: rows
    - Fortsetzung Tabelle Zauber
- Block 3. Spalte 3,4, Darstellung Tabelle, heightMode: rows
    - Fortsetzung Tabelle Zauber
- Block 4, Spalte 3,4, Darstellung Tabelle, heightMode: rows
    - Tabelle magische Gegenstände /Ausrüstung


## Definition der Tabellen

### Tabelle Fertigkeiten
| Fertigkeit | EW | PP |
| ---------- | -- | -- |
| Fertigkeit | EW | PP |

### Tabelle Waffen
|Waffe | EW | Schaden | Nah |
| ---- | -- | -- | -- |
|Waffe | EW | Schaden | Nah |

### Tabelle PP Zauber
| Zauber klasse | PP |
| ---------- | -- |
| Fertigkeit | PP |

### Tabelle Zauber
| AP/Prozess | Zauber | Zd/Rw | Wb/Wd | Wirkung | WirkZiel/Art |
| ---- | -------- | ---- | ---- | -------- | ---- |
| AP/Prozess | Zauber | Zd/Rw | Wb/Wd | Wirkung | WirkZiel/Art |

### Tabelle magische Gegenstände
| Gegenstande | Wirkung, Inhalt, und andere Erlä#uterungen |
| -------- | -------------------- |
| Gegenstande | Wirkung, Inhalt, und andere Erlä#uterungen |

### Tabelle Spielergebnisse

|Datum | | | | | | | | | | |
|ES    | | | | | | | | | | |
|EP    | | | | | | | | | | |
|Geld  | | | | | | | | | | |

### Tabelle LP/AP

| LP | | | | | | |
| AP | | | | | | |