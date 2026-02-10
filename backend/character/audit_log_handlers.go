package character

import (
	"bamort/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCharacterAuditLog gibt alle Audit-Log-Einträge für einen Charakter zurück
// GetCharacterAuditLog godoc
// @Summary Get character audit log
// @Description Returns audit log of all changes to a character (optionally filtered by field)
// @Tags Characters
// @Produce json
// @Param id path int true "Character ID"
// @Param field query string false "Filter by field name"
// @Success 200 {array} models.AuditLogEntry "Audit log entries"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Security BearerAuth
// @Router /api/characters/{id}/audit-log [get]
func GetCharacterAuditLog(c *gin.Context) {
	charID := c.Param("id")

	// Konvertiere String zu uint
	id, err := strconv.ParseUint(charID, 10, 32)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid character ID")
		return
	}

	// Filter für spezifisches Feld (optional)
	fieldName := c.Query("field")

	var entries []models.AuditLogEntry

	if fieldName != "" {
		entries, err = GetAuditLogForField(uint(id), fieldName)
	} else {
		entries, err = GetAuditLogForCharacter(uint(id))
	}

	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve audit log")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"character_id": uint(id),
		"entries":      entries,
	})
}

// GetAuditLogStats gibt Statistiken über Änderungen zurück
// GetAuditLogStats godoc
// @Summary Get audit log statistics
// @Description Returns statistics about character changes (total changes, by field, etc.)
// @Tags Characters
// @Produce json
// @Param id path int true "Character ID"
// @Success 200 {object} object "Audit log statistics"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Security BearerAuth
// @Router /api/characters/{id}/audit-log/stats [get]
func GetAuditLogStats(c *gin.Context) {
	charID := c.Param("id")

	// Konvertiere String zu uint
	id, err := strconv.ParseUint(charID, 10, 32)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid character ID")
		return
	}

	// Lade alle Einträge
	entries, err := GetAuditLogForCharacter(uint(id))
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve audit log")
		return
	}

	// Berechne Statistiken
	stats := map[string]interface{}{
		"total_changes":     len(entries),
		"by_field":          make(map[string]int),
		"by_reason":         make(map[string]int),
		"total_ep_spent":    0,
		"total_ep_gained":   0,
		"total_gold_spent":  0,
		"total_gold_gained": 0,
	}

	fieldStats := stats["by_field"].(map[string]int)
	reasonStats := stats["by_reason"].(map[string]int)

	for _, entry := range entries {
		// Zähle nach Feld
		fieldStats[entry.FieldName]++

		// Zähle nach Grund
		reasonStats[entry.Reason]++

		// Summen für EP und Gold
		if entry.FieldName == "experience_points" {
			if entry.Difference > 0 {
				stats["total_ep_gained"] = stats["total_ep_gained"].(int) + entry.Difference
			} else {
				stats["total_ep_spent"] = stats["total_ep_spent"].(int) + (-entry.Difference)
			}
		}

		if entry.FieldName == "gold" {
			if entry.Difference > 0 {
				stats["total_gold_gained"] = stats["total_gold_gained"].(int) + entry.Difference
			} else {
				stats["total_gold_spent"] = stats["total_gold_spent"].(int) + (-entry.Difference)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"character_id": uint(id),
		"stats":        stats,
	})
}
