# Experience and Wealth Update API

Diese API-Endpunkte ermöglichen es, Erfahrungspunkte und Vermögen von Charakteren zu aktualisieren.

## Erfahrungspunkte aktualisieren

### PUT /api/characters/{id}/experience

Aktualisiert die Erfahrungspunkte eines Charakters.

**Request Body:**
```json
{
  "experience_points": 150
}
```

**Response (Success):**
```json
{
  "message": "Experience updated successfully",
  "experience_points": 150
}
```

**Validierung:**
- `experience_points` muss >= 0 sein

## Vermögen aktualisieren

### PUT /api/characters/{id}/wealth

Aktualisiert das Vermögen eines Charakters. Es können einzelne Währungen oder alle zusammen aktualisiert werden.

**Request Body (einzelne Währung):**
```json
{
  "goldstücke": 250
}
```

**Request Body (mehrere Währungen):**
```json
{
  "goldstücke": 100,
  "silberstücke": 50,
  "kupferstücke": 200
}
```

**Response (Success):**
```json
{
  "message": "Wealth updated successfully",
  "wealth": {
    "goldstücke": 100,
    "silberstücke": 50,
    "kupferstücke": 200
  }
}
```

## Frontend Integration

Die ExperianceView.vue Komponente bietet jetzt:

1. **Bearbeitbare Eingabefelder** für EP und Goldstücke
2. **Automatisches Speichern** beim Verlassen der Felder oder Enter-Taste
3. **Visueller Feedback** während der Bearbeitung
4. **Fehlerbehandlung** mit Reset bei Fehlern
5. **Event Emission** (`character-updated`) für Parent-Komponenten

### Verwendung:

```vue
<ExperianceView 
  :character="selectedCharacter" 
  @character-updated="refreshCharacter" 
/>
```

Das `character-updated` Event sollte in der Parent-Komponente behandelt werden, um die Charakter-Daten neu zu laden.
