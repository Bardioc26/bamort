# Migration von GetLearnCost zu GetSkillCost

## Zusammenfassung der Änderungen

### Entfernte Legacy-Funktionen:
1. **`GetLearnCost`** - Legacy-Endpunkt für einfache Lernkosten entfernt
2. **Route `/api/characters/:id/learn`** - Endpunkt entfernt
3. **Test `TestGetLearnCost`** - Durch `TestGetSkillCost` ersetzt

### Grund für die Migration:

**GetLearnCost** war ein Legacy-Endpunkt mit begrenzten Funktionen:
- Nur für einfache Lernkosten von Fertigkeiten und Zaubern
- Keine Verbesserungskosten
- Keine Praxispunkte-Unterstützung
- Keine Belohnungsoptionen
- Verwendete veraltete `LearnRequestStruct`

**GetSkillCost** ist der moderne, vollständige Endpunkt:
- Berechnet Kosten für Lernen und Verbessern
- Unterstützt Multi-Level-Berechnungen
- Berücksichtigt Praxispunkte (PP) zur Kostenreduktion
- Unterstützt Belohnungsoptionen (kostenloses Lernen, Gold-für-EP-Tausch)
- Moderne `SkillCostRequest` Struktur

### Status der verbleibenden Legacy-Endpunkte:

**Noch vorhanden (marked als Legacy):**
- `/api/characters/:id/improve` - `GetSkillNextLevelCosts`
- `/api/characters/:id/improve/skill` - `GetSkillAllLevelCosts`

Diese verwenden noch `LearnRequestStruct` und könnten langfristig auch durch `GetSkillCost` ersetzt werden.

### Nächste Schritte für Audit-Log Integration:

Da es derzeit keine direkten Learn/Improve Endpunkte gibt, die automatisch EP und Gold abziehen, sollten folgende Endpunkte erstellt werden:

1. **POST `/api/characters/:id/learn-skill`** - Fertigkeit lernen (mit Audit-Log)
2. **POST `/api/characters/:id/improve-skill`** - Fertigkeit verbessern (mit Audit-Log)
3. **POST `/api/characters/:id/learn-spell`** - Zauber lernen (mit Audit-Log)
4. **POST `/api/characters/:id/improve-spell`** - Zauber verbessern (mit Audit-Log)

Diese würden:
- Kosten mit `GetSkillCost` berechnen
- EP und Gold automatisch abziehen
- Audit-Log-Einträge mit `ReasonSkillLearning`, `ReasonSkillImprovement`, etc. erstellen
- Charakterdaten aktualisieren

### Audit-Log Gründe für das Lernsystem:

```go
const (
    ReasonSkillLearning    = "skill_learning"      // Fertigkeit lernen
    ReasonSkillImprovement = "skill_improvement"   // Fertigkeit verbessern
    ReasonSpellLearning    = "spell_learning"      // Zauber lernen
    ReasonSpellImprovement = "spell_improvement"   // Zauber verbessern
)
```

Das Audit-Log-System ist bereits vollständig implementiert und bereit für die Integration in das Lernsystem.
