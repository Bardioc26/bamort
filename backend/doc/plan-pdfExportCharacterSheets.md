## Plan: PDF-Export für Charakterbögen

**Ziel:** Charakterbögen als mehrseitige PDF-Datei exportieren mit HTML-Templates und automatischem Seitenumbruch für Listen variabler Länge.

**Tech-Stack:** Go Backend (Gin), Vue 3 Frontend, GORM Datenbank, chromedp, Go html/template

**Architektur-Entscheidung:** Statt low-level PDF-Generierung nutzen wir HTML-Templates + chromedp für Browser-basiertes Rendering. Dies ermöglicht:
- Moderne CSS-Layouts (Flexbox, Grid, @page media queries)
- Einfachere Template-Wartung (HTML/CSS statt Koordinaten-basiert)
- Perfekte Schrift-/Unicode-Unterstützung
- WYSIWYG-Entwicklung (Templates im Browser testbar)

---

### Step 1: HTML-Templates vorbereiten & integrieren

**Template-Struktur:**
```
backend/
  templates/
    Default_A4_Quer/
      page1_stats.html          # Seite 1: Statistikseite
      page2_play.html           # Seite 2: Spielbogen  
      page3_spell.html          # Seite 3: Zauberseite
      page4_equip.html          # Seite 4: Ausrüstung
      shared/
        export_format_a4_quer.css  # Gemeinsames Stylesheet
        headerimg.png              # Statische Dekorationsgrafik
```

**Template-Konvertierung (HTML → Go Templates):**
- Bestehende HTML-Dateien (`char_p1_stats_a4_quer.html` etc.) dienen als Basis
- Beispieldaten werden durch Go template-Syntax ersetzt:
  - `Bjarnfinnur Haberdson` → `{{.Character.Name}}`
  - `18` (Grad) → `{{.Character.Grade}}`
  - `79` (St) → `{{.Character.Attributes.St}}`
  - Tabellen-Loops: `{{range .Character.Skills}}...{{end}}`

**CSS-Handling:**
- Externes Stylesheet `export_format_a4_quer.css` bleibt erhalten
- Templates linken auf Stylesheet via `<link rel="stylesheet" href="export_format_a4_quer.css">`
- Beim Rendering wird CSS-Datei aus `templates/Default_A4_Quer/shared/` geladen
- Chromedp wartet auf vollständiges Laden aller Stylesheets vor PDF-Generierung

**Bild-Handling:**
- **Charakter-Icons:** 
  - Aus Datenbank geladen (Character.IconData als Blob)
  - Als base64-kodiertes Data-URI in Template eingefügt:
    - `<img src="data:image/png;base64,{{.Character.IconBase64}}">`
- **Statische Bilder (headerimg.png):**
  - Aus Filesystem geladen (`templates/Default_A4_Quer/shared/headerimg.png`)
  - Als base64-kodierte Data-URIs in Template eingebettet oder als relativer Pfad

**Print-Media CSS:**
- Bestehende `@page { size: A4 landscape; margin: 1cm; }` bleibt erhalten
- Chromedp nutzt print media emulation für korrektes Seiten-Layout

---

### Step 2: Chromedp-Integration aufbauen

**Package-Struktur:**
```
backend/
  pdfrender/
    chromedp.go       # Chromedp-Wrapper
    templates.go      # Template-Loader & Renderer
    pagination.go     # Listen-Pagination-Logik
    viewmodel.go      # Domain → Template Mapping
```

**Chromedp-Wrapper (`chromedp.go`):**
- Funktion `RenderHTMLToPDF(htmlPages []string) ([]byte, error)`
  - Input: Array von gerenderten HTML-Strings (eine pro Seite)
  - Output: PDF als byte array
- Chromedp-Konfiguration:
  - Headless-Modus
  - Disable-GPU (für Server-Umgebung)
  - Print-to-PDF mit Optionen:
    - `Landscape: true` für A4 quer
    - `PrintBackground: true` für CSS-Hintergründe
    - `PreferCSSPageSize: true` für `@page`-Regeln
  - Timeout-Handling (z.B. 30s max)

**Template-Loader (`templates.go`):**
- Funktion `LoadTemplates() (*template.Template, error)`
  - Lädt alle `.html`-Dateien aus `templates/Default_A4_Quer/`
  - Parse mit Go's `html/template`
  - Validierung: Alle erwarteten Templates vorhanden
  - Caching: Templates einmalig beim Start laden
- Funktion `RenderTemplate(name string, data interface{}) (string, error)`
  - Rendert einzelnes Template mit Daten
  - Fügt base64-kodierte statische Bilder ein
  - Error-Handling für fehlende Template-Variablen

**Helper-Funktionen:**
- `ImageToBase64DataURI(imageData []byte, mimeType string) string`
  - Konvertiert Bild-Bytes zu Data-URI
- `LoadStaticImage(path string) (string, error)`
  - Lädt statisches Bild und konvertiert zu base64

---

### Step 3: Pagination-Logik implementieren

**Problem:** Listen (Skills, Weapons, Spells, Equipment) können Template-Kapazität überschreiten.

**Lösung:** Blockweise Pagination mit Fortsetzungsseiten

**Konzept:**
- Jede Liste hat eine `maxRows` pro Block/Seite
- Beispiel Skills Seite 1:
  - Spalte 3: max ~32 Zeilen
  - Spalte 4: max ~32 Zeilen
  - Total Seite 1: ~64 Skills
- Bei Überschreitung: Fortsetzungsseite mit gleicher Struktur

**Implementation (`pagination.go`):**

```go
type PageData struct {
    Character CharacterData
    Skills    []Skill      // nur die Skills für diese Seite
    Weapons   []Weapon     // nur die Weapons für diese Seite
    // ... weitere Listen
}

type PaginationConfig struct {
    Page1SkillsMax    int  // 64 (2x32)
    Page2SkillsMax    int  // je nach Layout
    Page2WeaponsMax   int
    Page3SpellsMax    int
    // ... für alle Listen
}

func PaginateCharacterData(char Character, config PaginationConfig) []PageData
```

**Pagination-Regeln (basierend auf Templates):**
- **Seite 1 (Stats):**
  - Skills: 2 Spalten à ~32 Zeilen = 64 Skills max
  - Overflow geht zu Seite 2
  
- **Seite 2 (Play):**
  - Skills (gelernte): Linke Spalte, ~24 Zeilen
  - Skills (ungelernte): Kleine Tabelle, ~15 Zeilen
  - Skills (Sprachen/etc.): Weitere Tabelle, ~11 Zeilen
  - Weapons: Große Tabelle rechts, ~30 Zeilen
  - Bei Overflow: Seite 2.1 mit gleicher Struktur
  
- **Seite 3 (Spells):**
  - Zauber: 2 Spalten à ~12 Zeilen = 24 Zauber
  - Magische Gegenstände: ~5 Zeilen
  - Bei Overflow: Seite 3.1
  
- **Seite 4 (Equipment):**
  - Ausrüstungs-Sektionen: Variable Zeilenanzahl pro Sektion
  - Bei Overflow: Seite 4.1

**Fortsetzungsseiten:**
- Gleiche Template-Struktur wie Hauptseite
- Header zeigt "Fortsetzung" an (optional)
- Datum bleibt gleich

---

### Step 4: View-Model Mapping

**Domain → Template Transformation (`viewmodel.go`):**

```go
type CharacterSheetViewModel struct {
    // Basis-Daten
    Character struct {
        Name           string
        Player         string
        Type           string
        Grade          int
        Birthdate      string
        Age            int
        Height         int
        Weight         int
        IconBase64     string  // base64-kodiertes Bild
        // ...
    }
    
    // Attribute
    Attributes struct {
        St, Gs, Gw, Ko, In, Zt, Au, PA, Wk, B int
    }
    
    // Abgeleitete Werte
    DerivedValues struct {
        LP, AP, GG, SG        int
        Abwehr, Resistenz     string
        Zaubern               string
        // Boni
        AusdauerBonus        int
        SchadenBonus         int
        AngriffBonus         int
        // ...
    }
    
    // Listen (bereits paginiert für diese Seite)
    Skills         []SkillViewModel
    Weapons        []WeaponViewModel
    Spells         []SpellViewModel
    MagicItems     []MagicItemViewModel
    Equipment      []EquipmentViewModel
    GameResults    []GameResultViewModel
    
    // Meta
    Date           string
    PageNumber     int
    IsContinuation bool
}
```

**Mapping-Funktion:**
```go
func MapCharacterToViewModel(
    char *models.Character, 
    pageType string,
    pageNumber int,
) (*CharacterSheetViewModel, error)
```

**Spezial-Handling:**
- **Fertigkeiten gruppieren:**
  - Seite 1: Alle Skills mit Grundwerten (EW aus Grundfertigkeit)
  - Seite 2: Nur gelernte Skills mit Spielwerten (PP, aktuelle EW)
  - Ungelernte Skills mit Modifikatoren
  
- **Zauber formatieren:**
  - Mehrzeilige Zellen: AP/Prozess als zwei Zeilen in einer Zelle
  - `<hr>` für visuelle Trennung in HTML
  
- **Ausrüstung nach Container:**
  - "Am Körper getragen"
  - "Becher, Holz" (Container 1)
  - "Wasserschlauch" (Container 2)
  - etc.

---

### Step 5: API-Endpoint implementieren

**Endpoint:**
```
GET /api/characters/{id}/pdf?format=a4_quer
```

**Handler-Flow:**
1. **Authentifizierung:** User-ID aus JWT/Session
2. **Autorisierung:** Prüfen ob User Zugriff auf Character hat
3. **Daten laden:**
   ```go
   char := LoadCharacterWithRelations(id)  // Preload Skills, Spells, Equipment, etc.
   ```
4. **Icon vorbereiten:**
   ```go
   if char.IconData != nil {
       iconBase64 := base64.StdEncoding.EncodeToString(char.IconData)
       char.IconBase64 = "data:image/png;base64," + iconBase64
   }
   ```
5. **Pagination:**
   ```go
   pages := PaginateCharacterData(char, GetPaginationConfig())
   ```
6. **Template-Rendering:**
   ```go
   htmlPages := make([]string, 0, 4)
   for i, pageData := range pages {
       var templateName string
       switch pageData.Type {
           case "stats": templateName = "page1_stats.html"
           case "play":  templateName = "page2_play.html"
           case "spell": templateName = "page3_spell.html"
           case "equip": templateName = "page4_equip.html"
       }
       html := RenderTemplate(templateName, pageData)
       htmlPages = append(htmlPages, html)
   }
   ```
7. **PDF-Generierung:**
   ```go
   pdfBytes, err := RenderHTMLToPDF(htmlPages)
   ```
8. **Response:**
   ```go
   w.Header().Set("Content-Type", "application/pdf")
   w.Header().Set("Content-Disposition", 
       fmt.Sprintf("attachment; filename=\"%s_charakterbogen.pdf\"", 
       sanitizeFilename(char.Name)))
   w.Write(pdfBytes)
   ```

**Error-Handling:**
- Character nicht gefunden: 404
- Keine Berechtigung: 403
- Template-Fehler: 500 mit Logging
- Chromedp-Fehler: 500 mit Logging
- Timeout: 504

**Route registrieren:**
- In `backend/character/routes.go`:
  ```go
  charGroup.GET("/:id/pdf", CharacterPDFHandler)
  ```

---

### Step 6: Frontend-Integration

**Vue Component: Character Detail View**

```vue
<template>
  <div class="character-actions">
    <button @click="downloadPDF" :disabled="downloading">
      <i class="icon-pdf"></i>
      {{ downloading ? 'Generiere PDF...' : 'Charakterbogen (PDF)' }}
    </button>
  </div>
</template>

<script setup>
const downloading = ref(false)

async function downloadPDF() {
  downloading.value = true
  try {
    const response = await fetch(
      `/api/characters/${characterId}/pdf?format=a4_quer`,
      { headers: { 'Authorization': `Bearer ${token}` } }
    )
    
    if (!response.ok) throw new Error('PDF generation failed')
    
    const blob = await response.blob()
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${characterName}_charakterbogen.pdf`
    a.click()
    window.URL.revokeObjectURL(url)
  } catch (error) {
    // Error-Toast anzeigen
  } finally {
    downloading.value = false
  }
}
</script>
```

**Format-Auswahl (optional, später):**
```vue
<select v-model="selectedFormat">
  <option value="a4_quer">A4 Querformat</option>
  <option value="a4_hoch">A4 Hochformat</option>
</select>
```

---

### Step 7: Deployment & Dependencies

**Go Dependencies:**
```bash
go get github.com/chromedp/chromedp
```

**System Requirements:**
- Chrome/Chromium installiert auf Server
  - Ubuntu: `apt-get install chromium-browser`
  - Alpine (Docker): `apk add chromium`
- Ausreichend Speicher für Chrome-Prozesse (~100-200MB pro Instanz)

**Docker Integration:**
```dockerfile
FROM golang:1.21-alpine

# Chrome installieren
RUN apk add --no-cache chromium

# Fonts für PDF (optional, verbessert Rendering)
RUN apk add --no-cache font-dejavu

# ... Rest des Dockerfiles
```

**Environment Variables:**
```bash
CHROME_BIN=/usr/bin/chromium-browser
PDF_TIMEOUT=30s
```

---

### LEGACY REFERENCE (old approach - for reference only):

The following sections describe the original plan using YAML templates and coordinate-based rendering. This approach has been replaced by the HTML template + chromedp approach described above.

---

### OLD Step 3: API-Endpoint implementieren

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

**Font-Konfiguration:**

- Zentral im `pdf`-Package:
  - Registrierung der verwendeten Fonts (liberation sans als Arial Replacement).
- Sicherstellen, dass:
  - Umlaute (ä, ö, ü, ß) und Sonderzeichen korrekt dargestellt werden,
  - gleiche Fontkonfiguration für alle Template seitendamit die Zeilenhöhe konstant bleibt

---

## Implementierungs-Reihenfolge (New chromedp Approach)

1. **HTML Templates vorbereiten (Step 1):**
   - HTML-Templates nach `backend/templates/Default_A4_Quer/` kopieren
   - CSS nach `backend/templates/Default_A4_Quer/shared/` kopieren
   - Statische Bilder nach `backend/templates/Default_A4_Quer/shared/images/` kopieren
   - Beispieldaten durch Go-Template-Syntax ersetzen (z.B. `{{.Character.Name}}`)
   - Base64-Konvertierung für Bilder testen

2. **Chromedp Integration (Step 2):**
   - `backend/pdfrender/chromedp.go` implementieren:
     - `RenderHTMLToPDF()` mit Context-Management
     - `PrintToPDF()` mit A4-Landscape-Settings
   - `backend/pdfrender/templates.go` implementieren:
     - `LoadTemplates()` mit Parse-Fehlerbehandlung
     - `RenderTemplate()` mit ViewModel-Binding
     - CSS/Bild-Einbettung testen
   - Unit-Test: Einfaches HTML → PDF generieren

3. **Pagination-Logik (Step 3):**
   - `backend/pdfrender/pagination.go` implementieren:
     - `PaginationConfig` struct mit max rows
     - `PaginateCharacterData()` mit Liste-Splitting
   - Regeln implementieren (basierend auf HTML-Template-Kapazitäten):
     - Seite 1: Skills max 64
     - Seite 2: Weapons max 30, Skills various
     - Seite 3: Spells max 24, Magic Items max 5
     - Seite 4: Equipment nach Container
   - Unit-Tests mit Overflow-Szenarien

4. **View-Model-Mapping (Step 4):**
   - `backend/pdfrender/viewmodel.go` implementieren:
     - `CharacterSheetViewModel` struct
     - `MapCharacterToViewModel()` mit Domain → Template Transformation
   - Spezialbehandlung:
     - Skills gruppieren (Grundwerte/Spielwerte/ungelernt)
     - Zauber formatieren (AP/Prozess mehrzeilig)
     - Ausrüstung nach Container gruppieren
   - Unit-Tests mit Mock-Character-Daten

5. **API-Endpoint (Step 5):**
   - Handler `CharacterPDFHandler` in `backend/character/` implementieren:
     - Auth/Authz
     - Character + Relations laden
     - Icon als base64 vorbereiten
     - Pagination aufrufen
     - Templates rendern
     - Chromedp PDF-Generierung
     - Stream als PDF
   - Route registrieren in `backend/character/routes.go`
   - Error-Handling (404/403/500/504)

6. **Frontend-Integration (Step 6):**
   - Vue-Component anpassen:
     - Download-Button mit Loading-State
     - Fetch-Call zu `/api/characters/:id/pdf`
     - Blob-Download-Handling
     - Error-Toast bei Fehlern
   - Optional: Format-Auswahl (A4-Quer/-Hoch)

7. **Deployment (Step 7):**
   - Dockerfile anpassen:
     - Chromium installieren
     - Fonts installieren (dejavu)
   - Environment-Variables konfigurieren:
     - `CHROME_BIN`, `PDF_TIMEOUT`
   - Docker-Compose testen

8. **Feintuning & Tests:**
   - Layout in verschiedenen PDF-Readern testen
   - `maxRows` justieren bei Bedarf
   - Performance-Tests (Generierungs-Dauer messen)
   - Testing-Checkliste durchgehen

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

## Testing-Checkliste (Chromedp Approach)

### Unit-Tests
- [ ] **Template Loading:**
  - [ ] Alle 4 Templates laden ohne Fehler
  - [ ] CSS-Datei korrekt eingebettet
  - [ ] Statische Bilder als base64 geladen
  - [ ] Fehlerbehandlung bei fehlendem Template

- [ ] **Pagination:**
  - [ ] Skills: < 64 → 1 Seite, > 64 → 2 Seiten
  - [ ] Weapons: < 30 → 1 Block, > 30 → Continuation Page
  - [ ] Spells: < 24 → 1 Seite, > 24 → 2 Seiten
  - [ ] Equipment: Multiple Container richtig gruppiert

- [ ] **View-Model Mapping:**
  - [ ] Character-Basics korrekt gemappt
  - [ ] Attributes alle vorhanden
  - [ ] Derived Values berechnet
  - [ ] Skills nach Kategorie gruppiert
  - [ ] Zauber mit mehrzeiligen Feldern (AP/Prozess)
  - [ ] Equipment nach Container sortiert
  - [ ] Character-Icon als base64 Data-URI

- [ ] **Chromedp Integration:**
  - [ ] Einfaches HTML → PDF erfolgreich
  - [ ] A4-Landscape korrekt (297mm x 210mm)
  - [ ] Print-CSS-Regeln angewendet (@page)
  - [ ] Timeout-Handling funktioniert

### Integration-Tests
- [ ] **End-to-End PDF-Generierung:**
  - [ ] Charakter mit minimalen Daten (leere Listen)
  - [ ] Charakter mit wenigen Skills (< 64) → 1 Seite
  - [ ] Charakter mit vielen Skills (> 64) → Multi-Page mit Continuation
  - [ ] Charakter mit allen Listen gefüllt → 4+ Seiten
  - [ ] Character mit DB-Icon → Icon korrekt im PDF
  - [ ] Character ohne Icon → Platzhalter oder leer

### API-Tests
- [ ] **Endpoint-Funktionalität:**
  - [ ] `GET /api/characters/:id/pdf` → 200 + PDF
  - [ ] Character nicht gefunden → 404
  - [ ] Nicht authentifiziert → 401
  - [ ] Keine Berechtigung (fremder Character) → 403
  - [ ] Timeout/Chromedp-Fehler → 500 mit Log
  - [ ] Content-Type `application/pdf` korrekt
  - [ ] Content-Disposition mit Filename korrekt
  - [ ] Filename sanitized (keine Sonderzeichen-Probleme)

### Frontend-Tests
- [ ] **Download-Funktion:**
  - [ ] Button klickbar
  - [ ] Loading-State während Generierung
  - [ ] Download startet automatisch
  - [ ] Filename korrekt (z.B. `Gandalf_charakterbogen.pdf`)
  - [ ] Error-Toast bei 404/500

### PDF-Qualitäts-Tests
- [ ] **Rendering-Qualität:**
  - [ ] PDF öffnet in Adobe Reader
  - [ ] PDF öffnet in Chrome/Firefox PDF-Viewer
  - [ ] PDF öffnet in Evince (Linux)
  - [ ] Layout identisch zu HTML-Preview in Browser
  - [ ] Umlaute korrekt (ä, ö, ü, ß)
  - [ ] Tabellen-Borders korrekt
  - [ ] Schriftgrößen lesbar
  - [ ] Keine abgeschnittenen Inhalte
  - [ ] Page-Breaks an richtigen Stellen

### Performance-Tests
- [ ] **Generierungs-Dauer:**
  - [ ] 1-seitiger Character: < 2 Sekunden
  - [ ] 4-seitiger Character: < 5 Sekunden
  - [ ] 10+ Seiten: < 10 Sekunden
  - [ ] Memory-Leak-Test (100 PDFs generieren)

### Deployment-Tests
- [ ] **Docker-Integration:**
  - [ ] Chromium in Container verfügbar
  - [ ] Fonts installiert (dejavu)
  - [ ] Environment-Variables gesetzt
  - [ ] PDF-Generierung in Container funktioniert
  - [ ] Keine Permission-Probleme mit Chrome-Binary

---

## Dateien zum Erstellen/Ändern (Chromedp Approach)

### Zu erstellende Dateien:

**Backend:**
- `backend/pdfrender/chromedp.go` (Chromedp-Wrapper, PDF-Generierung)
- `backend/pdfrender/templates.go` (Template-Loader, Renderer)
- `backend/pdfrender/pagination.go` (Pagination-Logik, Config)
- `backend/pdfrender/viewmodel.go` (ViewModel-Structs, Mapping-Funktionen)
- `backend/character/pdf_handler.go` (API-Endpoint-Handler)

**Templates:**
- `backend/templates/Default_A4_Quer/page1_stats.html` (Seite 1: Stats)
- `backend/templates/Default_A4_Quer/page2_play.html` (Seite 2: Play)
- `backend/templates/Default_A4_Quer/page3_spell.html` (Seite 3: Spells)
- `backend/templates/Default_A4_Quer/page4_equip.html` (Seite 4: Equipment)
- `backend/templates/Default_A4_Quer/shared/export_format_a4_quer.css` (Shared CSS)
- `backend/templates/Default_A4_Quer/shared/images/` (Statische Bilder)

**Tests:**
- `backend/pdfrender/chromedp_test.go`
- `backend/pdfrender/pagination_test.go`
- `backend/pdfrender/viewmodel_test.go`
- `backend/character/pdf_handler_test.go`

### Zu ändernde Dateien:

**Backend:**
- `backend/character/routes.go` (Route für PDF-Endpoint hinzufügen)
- `backend/go.mod` (chromedp-Dependency)

**Frontend:**
- `frontend/src/views/CharacterDetail.vue` (oder ähnlich - Download-Button)
- `frontend/src/services/characterService.js` (API-Call für PDF)

**Docker:**
- `docker/Dockerfile.backend` (Chromium + Fonts installieren)
- `docker/docker-compose.yml` (Environment-Variables)

**Dokumentation:**
- `backend/doc/plan-pdfExportCharacterSheets.md` (dieses Dokument)
- `README.md` (PDF-Export-Feature dokumentieren)
