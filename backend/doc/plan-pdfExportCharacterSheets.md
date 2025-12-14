## Plan: PDF-Export für Charakterbögen

**Ziel:** Charakterbögen als mehrseitige PDF-Datei exportieren mit flexiblen Templates (A4 Hochformat, A4 Querformat) und automatischem Seitenumbruch für Listen variabler Länge.

**Tech-Stack:** Go Backend (Gin), Vue 3 Frontend, GORM Datenbank

---

### Step 1: PDF-Bibliothek integrieren

**Details:**

- **Anforderungen aus dem Layout-/Template-Design:**
  - Absolut positionierbare Inhalte (Koordinaten in mm/pt).
  - Unterstützung für mehrere Seiten (Seitenumbrüche, Seitengröße, Ränder).
  - UTF‑8/Unicode für deutsche Texte inkl. Umlaute (DejaVu Sans als eingebettete Schrift).
  - Zeichen primitiver Formen (Linien, Rechtecke) für Tabellenraster/Boxen.
  - Möglichkeit, Text an beliebigen Positionen mit einheitlicher Schrift und Größe auszugeben.

- **Technische Auswahl (geplant):**
  - Einsatz einer reinen Go-PDF-Bibliothek (z.B. `gofpdf` o.ä.), die:
    - Seitengröße in mm/pt setzen kann,
    - Schriften laden/registrieren kann (TrueType, DejaVu Sans),
    - Textausgabe an absoluten Koordinaten unterstützt.

- **Integrationskonzept:**
  - Neues Package `backend/pdf`:
    - Kapselt die PDF-Bibliothek.
    - Stellt einfache Operationen bereit:
      - `NewDocument(pageSizeMm, marginsMm)`.
      - `AddPage()`.
      - `SetFont(name, sizePt, bold/normal)`.
      - `DrawText(xMm, yMm, text)`.
      - `DrawLine(x1Mm, y1Mm, x2Mm, y2Mm)`.
      - Hilfsfunktionen für Tabellenzellen (Text mit optional mehreren Zeilen).
    - Versteckt Einheitenkonvertierung (`ptToMm`) und Bibliotheksdetails.
  - Fonts:
    - DejaVu Sans (normal, ggf. bold) als eingebettete TTF-Fonts.
    - Registrierung beim Start im `pdf`-Package.
  - Minimaler Prototyp:
    - Funktion, die ein einfaches Test-PDF erzeugt (eine Seite, etwas Text), um Encoding/Fonts zu validieren.

---

### Step 2: Template-System aufbauen

**Details:**
- Layout- und Blockstruktur wird in `backend/doc/Characterbogen_Aufteilung.md` beschrieben und dient als Kanon.
- Ein Template (z.B. `backend/doc/template_a4_quer.yaml`) besteht aus:
  - **Globaler Seitendefinition**:
    - Seitengröße (z.B. A4 quer 297×210 mm)
    - Seitenränder (oben/unten/links/rechts)
    - Spaltenraster (Anzahl, Gutter; Breite wird abgeleitet)
    - Font-Konfiguration (Name, Größen, `lineHeightFactor`, Spacing)
  - **Seitentypen**:
    - Seite 1: Statistikseite
    - Seite 2: Spielbogen
    - Seite 3: Zauberseite
    - Seite 4: Ausrüstung
    - Seite 2.1: Fortsetzung Spielbogen (Fertigkeiten/Waffen)
    - Seite 3.1: Fortsetzung Zauber (Zauber/Magische Gegenstände)
  - **Blöcken pro Seite**:
    - Jeder Block hat:
      - `blockId` (z.B. 1–13 wie in `Characterbogen_Aufteilung.md`)
      - `pageType` (z.B. `page1_stats`, `page2_play`, `page2_play_cont`, `page3_spells`, …)
      - Spaltenstart/-spanne (`colStart`, `colSpan`), bezogen auf 4 gleich breite Spalten mit Gutter
      - `blockType`:
        - `labelValue` (Label: Wert)
        - `boxLabelValue` (Label: Wert, mit Rahmen)
        - `table`
      - `boxTypeId` (01=header, 02=stats, 03=table) → liefert Fontgröße, Style (`boxed`, `nobr`, `shorten`)
      - `anchor`: `"top"` (default) oder `"bottom"`
      - `heightMode`:
        - `"fixed"` (default): Höhe ergibt sich aus bekannter Zeilenanzahl (z.B. Spielergebnisse, LP/AP)
        - `"rows"`: Anzahl Zeilen wird dynamisch aus verfügbarem Platz berechnet (Fertigkeiten, Waffen, Zauber, Ausrüstung)
      - später: `binding` (welche Daten ins Feld/Tabelle kommen)
  - **Boxtypen**:
    - global definiert (Fontgröße, Style, Verhalten in Tabellen):
      - `header` (id: 01): Label:Wert, `f-Size: default`
      - `stats` (id: 02): Label:Wert, `f-Size: large`, `style: boxed`
      - `table` (id: 03): Tabellenzellen, `f-Size: default`, `style: nobr, shorten` (kein Umbruch, Text wird ggf. gekürzt)

- Alle dynamischen Listen (Fertigkeiten, Waffen, Zauber, Ausrüstung) werden ausschließlich in Blöcken mit `heightMode: rows` gerendert; alle anderen Tabellen haben eine feste, aus dem Layout bekannte Zeilenanzahl.

**Section-Renderer:**
- Aufgabe des Section-Renderers:
  - Aus Template + Characterdaten → konkrete Seiten mit absoluten Koordinaten (mm/pt) erzeugen.
  - Pro Seite:
    - Spaltenkoordinaten berechnen
    - Blöcke positionieren (Top-/Bottom-Anker)
    - Für `heightMode: rows` die `maxRows` aus verfügbaren Höhen bestimmen.

- Globale Layoutwerte (aus `Characterbogen_Aufteilung.md`):
  - Seitengröße A4 quer: 297×210 mm
  - Ränder: oben/unten/links/rechts je 6 mm
  - Spaltenraster:
    - 4 Spalten, alle gleich breit
    - Gutter: 4 mm
  - Font:
    - Name: DejaVu Sans
    - `size`:
      - `default`: 16 pt
      - `small`: 14 pt
      - `large`: 18 pt
    - `lineHeightFactor`:
      - `default`: 1.2
      - `small`: 1.1
      - `large`: 1.3
    - `spacing`:
      - `vertical`: 0.5 mm
      - `horizontal`: 0.5 mm

- Zeilenhöhe-Berechnung:
  - Konstante: `ptToMm` (z.B. 1 pt ≈ 0.3528 mm)
  - Für jede Font-Kategorie wird eine Zeilenhöhe berechnet:
    - `rowHeight_mm = fontSize_pt * lineHeightFactor * ptToMm + spacing.vertical`
  - Tabellen mit `style: nobr, shorten` werden so gerendert, dass jede Tabellenzeile **immer einzeilig** ist:
    - Kein automatischer Zeilenumbruch
    - Text ggf. hart gekürzt (`shorten`), um die konstante Zeilenhöhe zu garantieren.

- Mehrzeilige Tabellenzellen (Zauber-Tabelle, Spalten „AP/Prozess“ und „Zd/Rw“):
  - Einige Tabellenspalten bestehen logisch aus **zwei Werten**, die in derselben Zelle untereinander stehen.
    - Beispiel Zauber:
      - Spalte „AP/Prozess“: erste Zeile = `ap`, zweite Zeile = `prozess`
      - Spalte „Zd/Rw“: erste Zeile = `zd`, zweite Zeile = `rw`
  - Im Template kann eine Spalte optional als „multiline“ markiert werden:
    - Schema (pro Spalte):
      - `multiline.linesPerCell`: Anzahl der Textzeilen in dieser Zelle (hier: `2`)
      - `multiline.lines[]`: Liste von Bindings innerhalb einer Zelle, z.B.:
        - `lines[0].binding = "ap"`
        - `lines[1].binding = "prozess"`
  - Für Tabellen, die solche Spalten verwenden (Zauber-Tabellen), gilt:
    - Die effektive Zeilenhöhe dieser Tabelle wird so berechnet:
      - `rowHeight_mm_multiline = multiline.linesPerCell * (fontSize_pt * lineHeightFactor * ptToMm + spacing.vertical)`
    - Der Section-Renderer verwendet für diese Tabelle `rowHeight_mm_multiline` zur Berechnung von `maxRows`.
  - In der Implementierung:
    - Beim Rendern einer Zelle mit `multiline`:
      - Text nach den in `multiline.lines[]` angegebenen Bindings aufteilen,
      - jede Zeile übereinander zeichnen (gleiche X-Position, Y-Versatz pro Zeile),
      - keine automatischen Zeilenumbrüche jenseits der zwei konfigurierten Zeilen.

- Block-Flow und Anker-Regeln:
  - Wenn `anchor` für einen Block nicht angegeben ist → default `"top"`.
  - Wenn `heightMode` nicht angegeben ist → default `"fixed"`.
  - „Flow“:
    - Wenn nicht anders festgelegt, laufen die Blöcke „von links oben nach rechts unten“:
      - Reihenfolge in der Seitendefinition = Flow-Reihenfolge.
    - „Anders festgelegt“ bedeutet:
      - `anchor = "bottom"` → Block ist von unten verankert.
  - In jedem Spaltenbereich pro Seite:
    1. **Spaltenbereich berechnen**:
       - `colTopY = pageTopMargin` (bzw. Ende der Kopfblöcke, wenn Seite das so vorsieht)
       - `colBottomY = pageHeight - bottomMargin`
    2. **Bottom-Blöcke platzieren**:
       - Alle Blöcke mit `anchor="bottom"` in dieser Spaltenkombination ermitteln.
       - Für jeden:
         - Höhe bestimmen:
           - `height = fixedHeight` (für Label/Box-Blöcke mit bekannter Zeilenanzahl)
         - Position:
           - `blockY = currentBottomY - height`
           - `currentBottomY = blockY - verticalGap`
       - Ergebnis: reduzierter nutzbarer Bereich:
         - `maxContentBottomY = currentBottomY`
    3. **Top-Blöcke platzieren**:
       - Start: `cursorY = colTopY`
       - Alle `anchor="top"`-Blöcke in definierter Flow-Reihenfolge:
         - Für Blöcke mit `heightMode="fixed"`:
           - Höhe aus bekannter Struktur (Anzahl Label:Wert-Zeilen oder Tabellenzeilen).
           - `blockY = cursorY`
           - `cursorY += height + verticalGap`
         - Für Tabellen mit `heightMode="rows"`:
           - Verfügbarer vertikaler Platz:
             - `availableHeight = maxContentBottomY - cursorY`
           - `rowHeight` für diesen Block aus Font/Boxtyp
           - `maxRows = floor(availableHeight / rowHeight)`
           - `blockY = cursorY`
           - `cursorY += maxRows * rowHeight + verticalGap`
           - `maxRows` wird später in der Pagination-Logik verwendet, um „wie viele Elemente dieser Liste passen in diesen Block?“ zu beantworten.

- Spezieller Fall „Block unten, Liste oben füllt Rest“:
  - Beispiel Seite 1 (Statistikseite), Spalten 3–4:
    - Block 10 (Sinne) hat `anchor="bottom"`, feste Zeilenanzahl → feste Höhe.
    - Blöcke 9/11 (Fertigkeiten-Tabellen) haben `heightMode="rows"`, `anchor="top"`.
  - Renderer:
    - Platziert Block 10 von unten.
    - Berechnet aus dem verbleibenden vertikalen Bereich die `maxRows` für 9 (Spalte 3) und 11 (Spalte 4).
    - Diese `maxRows`-Werte liefern die Kapazität der ersten Statistikseite für Fertigkeiten. Alles darüber geht in spätere Seiten (Seite 2/2.1).

---

### Step 3: API-Endpoint implementieren

**Details:**

- **Endpoint-Design:**
  - HTTP-GET:
    - `GET /api/characters/{id}/pdf?templateId=a4_quer_char_sheet_v1`
  - Antwort:
    - `Content-Type: application/pdf`
    - `Content-Disposition: attachment; filename="<charname>_charbogen.pdf"`

- **Handler-Ablauf:**
  1. Authentifizierung/Autorisierung prüfen.
  2. Charakter inkl. aller benötigten Relationen laden:
     - Basis-/Statistikwerte,
     - Spielwerte,
     - Fertigkeiten,
     - Waffen,
     - Zauber,
     - magische Gegenstände,
     - Ausrüstung,
     - Spielergebnisse.
  3. Template wählen:
     - `templateId` aus Query (`template_a4_quer.yaml` als Default).
     - Template aus YAML in interne Structs laden/validieren.
  4. Domain-Daten in ein View-Model überführen, das zu den Template-Bindings passt (z.B. `character.attributes.st.base`, Zauberfelder `ap`, `prozess`, `zd`, `rw`, `wb`, `wd` usw.).
  5. Section-Renderer + Pagination-Engine aufrufen:
     - `pages := RenderCharacterSheet(template, viewModel)`.
  6. PDF-Dokument mit dem `pdf`-Package erzeugen:
     - Für jede Seite:
       - Seite hinzufügen,
       - alle Blöcke mit ihren Koordinaten und Daten rendern.
  7. PDF an den Client streamen.

- **Fehlerfälle:**
  - Charakter nicht gefunden → 404.
  - Template nicht gefunden/ungültig → 500, Logging, ggf. Fallback auf Default-Template.
  - PDF-Erzeugung fehlgeschlagen → 500, Logging.

**Route registrieren:**

- In der bestehenden Gin-Router-Konfiguration:
  - z.B. in `cmd/backend/main.go`:
    - Gruppe `/api/characters` → `GET("/:id/pdf", CharacterPDFHandler)`.

---

### Step 4: Listen-Rendering mit Pagination

**Details:**

- **Datenaufbereitung:**
  - Aus dem Domänenmodell werden pro Charakter logisch getrennte Listen erzeugt:
    - `statsSkills` (Fertigkeiten-Grundwerte für Seite 1),
    - `playSkills` (Fertigkeiten-Spielwerte für Seite 2/2.x),
    - `weapons`,
    - `spells`,
    - `magicItems`,
    - `equipmentOnBody`, `equipmentContainer1`, `equipmentContainer2`, …,
    - `gameResults`.
  - Diese Listen sind unabhängig von Seiten und Blöcken.

- **Template-Auswertung:**
  - Template definiert:
    - Seiten (`pages[]` mit `pageType`),
    - Blöcke (`blocks[]` pro Seite),
    - für `heightMode="rows"`: Tabellen, deren Zeilenanzahl limitiert ist.
  - Der Section-Renderer berechnet:
    - Für jeden `heightMode="rows"`-Block:
      - `rowHeight` (inkl. `multiline.linesPerCell` bei Zauber),
      - Nutzhöhe des Blocks (unter Berücksichtigung von Top-/Bottom-Ankern),
      - `maxRows` = maximale Anzahl Datensätze für diesen Block.

- **Pagination-Engine:**
  - Pro Liste (Fertigkeiten, Waffen, Zauber, Ausrüstung) wird eine Blockfüllreihenfolge definiert, basierend auf `Characterbogen_Aufteilung.md`:
    - Fertigkeiten:
      - Grundwerte: Seite 1, Block 9 → Block 11.
      - Spielwerte: Seite 2, Block 4 → Block 6 → Seite 2.1 (Block 2/3/4) → weitere 2.x-Seiten.
    - Waffen:
      - Seite 2, Block 9 → Block 12 → Seite 2.1 (Block 5) → weitere 2.x-Seiten.
    - Zauber:
      - Seite 3, Block 2 → Block 3 → 3.x (Block 2/3) → weitere 3.x-Seiten.
    - Magische Gegenstände:
      - Seite 3, Block 4 → 3.x (Block 4) → weitere 3.x-Seiten.
    - Ausrüstung:
      - Seite 4, Blöcke 2–5 → ggf. 4.x mit gleichen Blöcken.
  - Algorithmus:
    1. Für jede Liste die zugehörigen Blöcke (Slots) in Reihenfolge sammeln (mit `pageType`, `blockId`, `maxRows`).
    2. Über die Einträge der Liste iterieren und in diese Slots schreiben, bis `maxRows` erreicht ist.
    3. Wenn alle Slots eines Seitentyps voll sind und noch Einträge existieren:
       - Neue Instanz des Fortsetzungs-Seitentyps erzeugen (z.B. `page2_play_cont`).
       - Die Slots dieser neuen Seite an die Füllreihenfolge anhängen.
  - Ergebnis:
    - Eine Liste konkreter Seiteninstanzen, bei denen jeder `heightMode="rows"`-Block bereits die passenden Datenzeilen (Schnitt der Liste) enthält.

- **Schnittstelle Renderer/Pagination:**
  - Renderer erwartet:
    - Seiten in Reihenfolge mit:
      - `pageType`
      - Liste von Blöcken mit:
        - Layout-Infos (Koordinaten, Breite, Typ),
        - bereits paginierte Daten (Zeilen für Tabellen).

---

### Step 5: Frontend-Integration

**Details:**

- **Download im Charakter-UI (Vue 3):**
  - In der Detailansicht eines Charakters:
    - Button „Charakterbogen (PDF)“.
    - Klick:
      - HTTP-Request auf `GET /api/characters/{id}/pdf?templateId=a4_quer_char_sheet_v1`.
      - Browser-Auslösung eines Downloads (`Content-Disposition`).

- **Template-Auswahl (optional, später):**
  - Backend-Endpunkt `GET /api/pdf-templates`:
    - Liefert Liste der verfügbaren Templates (`id`, `name`, Ausrichtung).
  - Frontend:
    - Dropdown neben Download-Button (z.B. „A4 quer“, „A4 hoch“).
    - Ausgewählte `templateId` wird in den PDF-Request aufgenommen.

- **Fehler-Handling im UI:**
  - Bei HTTP-Fehlern:
    - Kurze Meldung („PDF-Export fehlgeschlagen“).
  - Kein Client-seitiges Layout:
    - Frontend kennt nur `templateId`, die komplette Renderlogik liegt im Backend.

---

### Step 6: Template-Konfiguration

**Details:**

- **Ablage & Format:**
  - Templates als YAML-Dateien im Repo:
    - z.B. `backend/doc/template_a4_quer.yaml` als Referenztemplate für A4 quer.
  - Struktur orientiert sich an:
    - `template.Id`, `name`,
    - `layout` (pageSizeMm, marginsMm, columns, font, ptToMm),
    - `boxTypes`,
    - `pages[]` mit `blocks[]`.

- **Laden und Validieren:**
  - Beim Start des Backends:
    - Template-Dateien aus einem definierten Verzeichnis laden.
    - Validierung:
      - Pflichtfelder vorhanden (Id, layout, pages, blocks).
      - `boxTypeId` verweist auf existierenden `boxTypes`-Eintrag.
      - Spaltenbereiche sind konsistent (keine negativen `colStart/colSpan`, max 4 Spalten).
  - Gültige Templates werden in einer Registry gehalten:
    - Lookup per `templateId`.

- **MVP-Einschränkungen:**
  - Templates nur als Files im Repo pflegen (kein UI zum Hochladen/Ändern).
  - Änderungen am Layout erfolgen durch Anpassen der YAML-Dateien und Deployment.

**In-Code-Konstanten (KISS-Prinzip):**

- `defaultTemplateId = "a4_quer_char_sheet_v1"`.
- Pfad zum Template-Verzeichnis (z.B. `./backend/doc` oder eigener Ordner).
- Fallback-Logik:
  - Wenn angefragtes Template nicht gefunden wird:
    - Auf `defaultTemplateId` zurückfallen (sofern gültig).

**Font-Konfiguration:**

- Zentral im `pdf`-Package:
  - Registrierung der verwendeten Fonts (DejaVu Sans normal/bold).
  - Mappings:
    - BoxType → Fontgröße (`default/small/large`).
- Sicherstellen, dass:
  - Umlaute (ä, ö, ü, ß) und Sonderzeichen korrekt dargestellt werden,
  - gleiche Fontkonfiguration wie in Template/Plan (`Characterbogen_Aufteilung.md`).

---

## Implementierungs-Reihenfolge

1. **PDF-Bibliothek integrieren (Step 1):**
   - `pdf`-Package bauen, Test-PDF generieren.
2. **Template-Lader und -Validator (Step 2 & 6):**
   - YAML-Struktur für Templates definieren.
   - `template_a4_quer.yaml` laden/prüfen.
3. **Section-Renderer (Layout-Math, Step 2):**
   - Spalten-/Koordinatenberechnung, Top-/Bottom-Anker, `rowHeight`-Berechnung.
4. **Pagination-Engine (Step 4):**
   - Blockfüllreihenfolgen gemäß `Characterbogen_Aufteilung.md` implementieren.
   - Unit-Tests mit künstlichen Listen (viele Skills, viele Zauber, etc.).
5. **View-Model-Mapping:**
   - Domänenstrukturen → Bindings im Template (Attribute, Fertigkeiten, Waffen, Zauber …).
6. **PDF-Render-Pipeline verbinden:**
   - ViewModel → Pagination → Section-Renderer → `pdf`-Package.
7. **API-Endpoint & Routing (Step 3):**
   - Handler implementieren, PDF streamen.
8. **Frontend-Integration (Step 5):**
   - Download-Button, optional Template-Auswahl.
9. **Feintuning & Tests:**
   - Layout gegen Originalbogen prüfen.
   - `maxRows` ggf. feinjustieren.
   - Tests aus der Testing-Checkliste durchgehen.

---

## Further Considerations

### Mehrsprachigkeit
- **Aktuell:** Templates fest auf Deutsch
- **Später:** Feldnamen aus i18n-Config laden
- **KISS:** Start mit Deutsch, i18n bei Bedarf nachrüsten

### Template-Customization
- **Option A (komplex):** Benutzer eigene Templates hochladen lassen
- **Option B (KISS):** 2–3 vordefinierte Templates fest im Code
- **Empfehlung:** Option B für MVP, Option A als Feature später

### Caching
- **On-the-fly (Start):** PDF bei jedem Request neu generieren
- **Mit Cache (später):** PDF cachen, Invalidierung bei Character-Update
- **KISS:** Start ohne Cache, Performance-Optimierung bei Bedarf

### Fehlerbehandlung
- Charakter nicht gefunden → 404
- PDF-Generierung fehlgeschlagen → 500 mit Fehlerlog
- Ungültiges Template → Fallback auf Default-Template

### Logging
- PDF-Export-Anfragen loggen (Character-ID, Template, User)
- Generierungs-Dauer messen für Performance-Monitoring
- Fehler detailliert loggen für Debugging

---

## Testing-Checkliste

- [ ] Charakter mit wenigen Skills (1 Seite)
- [ ] Charakter mit vielen Skills (Seitenumbruch)
- [ ] Charakter mit allen Listen gefüllt (Multi-Page)
- [ ] A4-Quer-Template (`template_a4_quer.yaml`)
- [ ] (später) A4-Hoch-Template
- [ ] Ungültiges Template (Fallback zu Default)
- [ ] Charakter nicht gefunden (404)
- [ ] Nicht authentifiziert (401)
- [ ] Download-Filename korrekt
- [ ] PDF öffnet korrekt in verschiedenen PDF-Readern
- [ ] Umlaute (ä, ö, ü, ß) korrekt dargestellt
- ...

---

## Dateien zum Erstellen/Ändern

- `backend/doc/Characterbogen_Aufteilung.md` (Layout-Referenz)
- `backend/doc/template_a4_quer.yaml` (konkretes Template)
- `backend/pdf/*` (PDF-Wrapper)
- `backend/pdfrender/*` oder ähnliches (Section-Renderer + Pagination)
- API-Handler im Backend (z.B. `backend/api/character_pdf.go`)
- Frontend-Komponenten für Download/Template-Auswahl
