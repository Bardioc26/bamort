package maintenance

import (
	"bamort/database"
	"bamort/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type believeResponse struct {
	models.Believe
	SourceCode string `json:"source_code,omitempty"`
}

type believeUpdateRequest struct {
	Name         string `json:"name"`
	Beschreibung string `json:"beschreibung"`
	SourceID     *uint  `json:"source_id"`
	PageNumber   *int   `json:"page_number"`
}

func resolveGameSystemOrDefault(c *gin.Context) *models.GameSystem {
	gsIDStr := c.Query("game_system_id")
	gsName := c.Query("game_system")

	var gsID uint
	if gsIDStr != "" {
		id, err := strconv.ParseUint(gsIDStr, 10, 32)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid game_system_id")
			return nil
		}
		gsID = uint(id)
	}

	gs := models.GetGameSystem(gsID, gsName)
	if gs == nil {
		respondWithError(c, http.StatusBadRequest, "Invalid game system")
		return nil
	}

	return gs
}

// GetBelieves godoc
// @Summary Get beliefs
// @Description Returns list of all beliefs/religions (maintainer only)
// @Tags Maintenance
// @Produce json
// @Success 200 {array} models.Believe "List of beliefs"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/gsm-believes [get]
func GetBelieves(c *gin.Context) {
	gs := resolveGameSystemOrDefault(c)
	if gs == nil {
		return
	}

	var believes []models.Believe
	if err := database.DB.Where("game_system=? OR game_system_id=?", gs.Name, gs.ID).
		Order("name ASC").
		Find(&believes).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve believes")
		return
	}

	var sources []models.Source
	if err := database.DB.Find(&sources).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve sources")
		return
	}

	sourceMap := make(map[uint]string, len(sources))
	for _, source := range sources {
		sourceMap[source.ID] = source.Code
	}

	enhanced := make([]believeResponse, len(believes))
	for i, believe := range believes {
		enhanced[i] = believeResponse{
			Believe:    believe,
			SourceCode: sourceMap[believe.SourceID],
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"believes": enhanced,
		"sources":  sources,
	})
}

// UpdateBelieve godoc
// @Summary Update belief
// @Description Updates an existing belief/religion (maintainer only)
// @Tags Maintenance
// @Accept json
// @Produce json
// @Param id path int true "Belief ID"
// @Param belief body models.Believe true "Updated belief data"
// @Success 200 {object} models.Believe "Updated belief"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/gsm-believes/{id} [put]
func UpdateBelieve(c *gin.Context) {
	gs := resolveGameSystemOrDefault(c)
	if gs == nil {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var believe models.Believe
	if err := database.DB.First(&believe, uint(id)).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "Believe not found")
		return
	}

	var payload believeUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	name := strings.TrimSpace(payload.Name)
	if name == "" {
		respondWithError(c, http.StatusBadRequest, "name is required")
		return
	}

	believe.Name = name
	believe.Beschreibung = payload.Beschreibung
	if payload.SourceID != nil {
		believe.SourceID = *payload.SourceID
	} else {
		believe.SourceID = 0
	}
	if payload.PageNumber != nil {
		believe.PageNumber = *payload.PageNumber
	} else {
		believe.PageNumber = 0
	}

	believe.GameSystem = gs.Name
	believe.GameSystemId = gs.ID

	if err := database.DB.Save(&believe).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update believe")
		return
	}

	c.JSON(http.StatusOK, believeResponse{
		Believe:    believe,
		SourceCode: lookupSourceCode(believe.SourceID),
	})
}

func lookupSourceCode(sourceID uint) string {
	if sourceID == 0 {
		return ""
	}

	var source models.Source
	if err := database.DB.Select("code").First(&source, sourceID).Error; err != nil {
		return ""
	}

	return source.Code
}
