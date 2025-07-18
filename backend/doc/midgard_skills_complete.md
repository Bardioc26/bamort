# Vollständige Midgard Fertigkeiten-Konsolidierung

## Charakterklassen (vollständig)
```go
"As": "Assassine",
"Bb": "Barbar", 
"Gl": "Glücksritter",
"Hä": "Händler",
"Kr": "Krieger",
"Sp": "Spitzbube",
"Wa": "Waldläufer",
"Ba": "Barde",
"Or": "Ordensritter",
"Dr": "Druide",
"Hx": "Hexer",
"Ma": "Magier",
"PB": "Priester Beschützer",
"PS": "Priester Streiter",
"Sc": "Schamane"
```

## Kategorien mit allen Fertigkeiten und Schwierigkeiten

### 1. Alltag
**Lernkosten (LE):**
- leicht (1 LE): Bootfahren+12, Glücksspiel+12, Klettern+12, Musizieren+12, Reiten+12, Schwimmen+12, Seilkunst+12, Wagenlenken+12
- normal (1 LE): Schreiben+8, Sprache+8
- schwer (2 LE): Erste Hilfe+8, Etikette+8
- sehr schwer (10 LE): Gerätekunde+8, Geschäftssinn+8

### 2. Freiland
**Lernkosten (LE):**
- leicht (1 LE): Überleben+8
- normal (2 LE): Naturkunde+8, Pflanzenkunde+8, Tierkunde+8
- schwer (4 LE): Schleichen+8, Spurensuche+8, Tarnen+8

### 3. Halbwelt
**Lernkosten (LE):**
- leicht (1 LE): Balancieren+12, Fälschen+12, Gaukeln+12, Glücksspiel+12, Klettern+12
- normal (2 LE): Akrobatik+8
- schwer (2 LE): Gassenwissen+8, Stehlen+8
- sehr schwer (10 LE): Betäuben+8

### 4. Kampf
**Lernkosten (LE):**
- leicht (1 LE): Geländelauf+12, Reiten+12, Reiterkampf+12
- normal (2 LE): Anführen+8, Athletik+8
- schwer (10 LE): Betäuben+8
- sehr schwer (10 LE): Beidhändiger Kampf+5, Kampf in Vollrüstung+5, Fechten+5, Scharfschießen+5

### 5. Körper
**Lernkosten (LE):**
- leicht (1 LE): Balancieren+12, Geländelauf+12, Klettern+12, Schwimmen+12
- normal (1 LE): Tauchen+8
- schwer (2 LE): Akrobatik+8, Athletik+8, Laufen+8, Meditieren+8

### 6. Sozial
**Lernkosten (LE):**
- leicht (2 LE): Anführen+8, Etikette+8, Verführen+8, Verstellen+8
- normal (2 LE): Beredsamkeit+8, Gassenwissen+8, Verhören+8
- schwer (4 LE): Menschenkenntnis+8

### 7. Unterwelt
**Lernkosten (LE):**
- leicht (2 LE): Fälschen+12, Stehlen+8
- normal (2 LE): Gassenwissen+8, Verhören+8
- schwer (4 LE): Einschüchtern+8, Menschenkenntnis+8

### 8. Waffen
**Lernkosten (LE):**
- leicht (1 LE): Armbrust+8, Bogen+8, Dolch+8, Hiebwaffen+8, Kettenwaffen+8, Wurfwaffen+8
- normal (2 LE): Fechten+5, Schilde+8, Stangenwaffen+8, Zweihandhiebwaffen+8, Zweihandflegel+8
- schwer (10 LE): Beidhändiger Kampf+5, Kampf in Vollrüstung+5, Scharfschießen+5

### 9. Wissen
**Lernkosten (LE):**
- normal (2 LE): Alchemie+8, Götter+8, Geschichte+8, Heraldik+8, Kriegsführung+8, Rechtskunde+8, Strategie+8
- schwer (4 LE): Gerätekunde+8, Pflanzenkunde+8, Tierkunde+8

## Mapping für Datenbank-Skills

### Skills ohne Category/Difficulty in JSON:
```
Hören -> Alltag/leicht (basierend auf Wahrnehmung)
Nachtsicht -> Alltag/leicht (basierend auf Wahrnehmung)
Riechen -> Alltag/leicht (basierend auf Wahrnehmung)
Sechster Sinn -> Alltag/leicht (basierend auf Wahrnehmung)
Sehen -> Alltag/leicht (basierend auf Wahrnehmung)
Wahrnehmung -> Alltag/leicht (implizit aus Tabellen)
Athletik -> Kampf/normal ODER Körper/schwer (Dopplung in Tabellen!)
Erste Hilfe -> Alltag/schwer
Geländelauf -> Kampf/leicht ODER Körper/leicht (Dopplung!)
Kampf in Vollrüstung -> Kampf/sehr_schwer ODER Waffen/schwer (Dopplung!)
Klettern -> Alltag/leicht ODER Halbwelt/leicht ODER Körper/leicht (Dreifachung!)
Landeskunde -> fehlt in Tabellen (sollte Wissen/normal sein)
Laufen -> Körper/schwer
Robustheit -> fehlt in Tabellen (sollte Körper/normal sein)
Sprache -> Alltag/normal
```

### Weapon Skills:
```
Armbrüste -> Waffen/leicht
Einhandschlagwaffen -> Waffen/leicht (als Hiebwaffen)
Schilde -> Waffen/normal
Spießwaffen -> Waffen/normal (als Stangenwaffen)
Stichwaffen -> Waffen/leicht (als Dolch)
Stielwurfwaffen -> Waffen/leicht (als Wurfwaffen)
Waffenloser Kampf -> fehlt in Tabellen (sollte Waffen/normal sein)
Zweihandschlagwaffen -> Waffen/normal (als Zweihandhiebwaffen)
```

## Probleme identifiziert:

1. **Doppelte/Dreifache Einträge**: Athletik, Geländelauf, Klettern erscheinen in mehreren Kategorien
2. **Fehlende Skills**: Landeskunde, Robustheit, Waffenloser Kampf sind nicht in den Lerntabellen
3. **Inkonsistente Namensgebung**: Hiebwaffen vs. Einhandschlagwaffen
4. **Wahrnehmungsunterkategorien**: Hören, Sehen, etc. sind separate Skills aber nicht in Lerntabellen

## Empfohlene Auflösung:

1. **Eindeutige Zuordnung** nach Kontext/Priorität
2. **Fehlende Skills** mit sinnvollen Defaults ergänzen
3. **Namens-Mapping** für Kompatibilität
4. **Wahrnehmungssubskills** als Alltag/leicht behandeln
