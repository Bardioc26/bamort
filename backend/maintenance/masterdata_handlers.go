package maintenance

import (
	"bamort/database"
	"bamort/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetGameSystems godoc
// @Summary Get game systems
// @Description Returns list of all game systems (maintainer only)
// @Tags Maintenance
// @Produce json
// @Success 200 {array} models.GameSystem "List of game systems"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/game-systems [get]
func GetGameSystems(c *gin.Context) {
	var systems []models.GameSystem
	if err := database.DB.Order("code ASC").Find(&systems).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve game systems")
		return
	}

	c.JSON(http.StatusOK, gin.H{"game_systems": systems})
}

type gameSystemUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

// UpdateGameSystem godoc
// @Summary Update game system
// @Description Updates an existing game system (maintainer only)
// @Tags Maintenance
// @Accept json
// @Produce json
// @Param id path int true "Game system ID"
// @Param system body models.GameSystem true "Updated game system data"
// @Success 200 {object} models.GameSystem "Updated game system"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/game-systems/{id} [put]
func UpdateGameSystem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var gs models.GameSystem
	if err := database.DB.First(&gs, uint(id)).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "Game system not found")
		return
	}

	var payload gameSystemUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	name := strings.TrimSpace(payload.Name)
	if name != "" {
		gs.Name = name
	}
	gs.Description = payload.Description
	if payload.IsActive != nil {
		gs.IsActive = *payload.IsActive
	}

	if err := database.DB.Save(&gs).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update game system")
		return
	}

	c.JSON(http.StatusOK, gs)
}

// --- Sources ---

type sourceUpdateRequest struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Edition     string `json:"edition"`
	Publisher   string `json:"publisher"`
	PublishYear int    `json:"publish_year"`
	Description string `json:"description"`
	IsCore      *bool  `json:"is_core"`
	IsActive    *bool  `json:"is_active"`
}

// GetLitSources godoc
// @Summary Get literature sources
// @Description Returns list of all literature sources (maintainer only)
// @Tags Maintenance
// @Produce json
// @Success 200 {array} models.Source "List of literature sources"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/gsm-lit-sources [get]
func GetLitSources(c *gin.Context) {
	gs := resolveGameSystemOrDefault(c)
	if gs == nil {
		return
	}

	var sources []models.Source
	if err := database.DB.Where("game_system=? OR game_system_id=?", gs.Name, gs.ID).
		Order("code ASC").
		Find(&sources).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve sources")
		return
	}

	c.JSON(http.StatusOK, gin.H{"sources": sources})
}

// UpdateLitSource godoc
// @Summary Update literature source
// @Description Updates an existing literature source (maintainer only)
// @Tags Maintenance
// @Accept json
// @Produce json
// @Param id path int true "Literature source ID"
// @Param source body models.Source true "Updated literature source data"
// @Success 200 {object} models.Source "Updated literature source"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/gsm-lit-sources/{id} [put]
func UpdateLitSource(c *gin.Context) {
	gs := resolveGameSystemOrDefault(c)
	if gs == nil {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var src models.Source
	if err := database.DB.First(&src, uint(id)).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "Source not found")
		return
	}

	var payload sourceUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if name := strings.TrimSpace(payload.Name); name != "" {
		src.Name = name
	}
	src.FullName = payload.FullName
	src.Edition = payload.Edition
	src.Publisher = payload.Publisher
	src.PublishYear = payload.PublishYear
	src.Description = payload.Description
	if payload.IsCore != nil {
		src.IsCore = *payload.IsCore
	}
	if payload.IsActive != nil {
		src.IsActive = *payload.IsActive
	}
	src.GameSystem = gs.Name
	src.GameSystemId = gs.ID

	if err := database.DB.Save(&src).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update source")
		return
	}

	c.JSON(http.StatusOK, src)
}

// --- Misc lookup ---

type miscUpdateRequest struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	SourceID   *uint  `json:"source_id"`
	PageNumber *int   `json:"page_number"`
}

// GetMisc godoc
// @Summary Get miscellaneous master data
// @Description Returns list of all miscellaneous master data entries (maintainer only)
// @Tags Maintenance
// @Produce json
// @Success 200 {array} object "List of miscellaneous data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/gsm-misc [get]
func GetMisc(c *gin.Context) {
	gs := resolveGameSystemOrDefault(c)
	if gs == nil {
		return
	}

	keyFilter := strings.TrimSpace(c.Query("key"))

	var items []models.MiscLookup
	q := database.DB.Where("game_system=? OR game_system_id=?", gs.Name, gs.ID)
	if keyFilter != "" {
		q = q.Where("`key` = ?", keyFilter)
	}
	if err := q.Order("`key` ASC, value ASC").Find(&items).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve misc")
		return
	}

	c.JSON(http.StatusOK, gin.H{"misc": items})
}

// UpdateMisc godoc
// @Summary Update miscellaneous master data
// @Description Updates a miscellaneous master data entry (maintainer only)
// @Tags Maintenance
// @Accept json
// @Produce json
// @Param id path int true "Misc data ID"
// @Param data body object true "Updated miscellaneous data"
// @Success 200 {object} object "Updated miscellaneous data"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/gsm-misc/{id} [put]
func UpdateMisc(c *gin.Context) {
	gs := resolveGameSystemOrDefault(c)
	if gs == nil {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var item models.MiscLookup
	if err := database.DB.First(&item, uint(id)).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "Misc entry not found")
		return
	}

	var payload miscUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if key := strings.TrimSpace(payload.Key); key != "" {
		item.Key = key
	}
	if payload.Value != "" {
		item.Value = payload.Value
	}
	if payload.SourceID != nil {
		item.SourceID = *payload.SourceID
	}
	if payload.PageNumber != nil {
		item.PageNumber = *payload.PageNumber
	}
	item.GameSystem = gs.Name
	item.GameSystemId = gs.ID

	if err := database.DB.Save(&item).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update misc entry")
		return
	}

	c.JSON(http.StatusOK, item)
}

// --- Skill improvement cost2 ---

type skillImprovementUpdateRequest struct {
	CurrentLevel *int  `json:"current_level"`
	TERequired   *int  `json:"te_required"`
	CategoryID   *uint `json:"category_id"`
	DifficultyID *uint `json:"difficulty_id"`
}

// GetSkillImprovementCost2 godoc
// @Summary Get skill improvement costs
// @Description Returns skill improvement cost table (maintainer only)
// @Tags Maintenance
// @Produce json
// @Success 200 {array} models.SkillImprovementCost "Skill improvement costs"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/skill-improvement-cost2 [get]
func GetSkillImprovementCost2(c *gin.Context) {
	var costs []models.SkillImprovementCost
	if err := database.DB.Order("current_level ASC").Find(&costs).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve skill improvement costs")
		return
	}

	categoryNames, difficultyNames, err := loadSkillMetadata(costs)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve skill improvement costs")
		return
	}

	responses := make([]skillImprovementCostResponse, len(costs))
	for i, cost := range costs {
		responses[i] = skillImprovementCostResponse{
			SkillImprovementCost: cost,
			CategoryName:         categoryNames[cost.CategoryID],
			DifficultyName:       difficultyNames[cost.DifficultyID],
		}
	}

	c.JSON(http.StatusOK, gin.H{"costs": responses})
}

// UpdateSkillImprovementCost2 godoc
// @Summary Update skill improvement cost
// @Description Updates a skill improvement cost entry (maintainer only)
// @Tags Maintenance
// @Accept json
// @Produce json
// @Param id path int true "Cost entry ID"
// @Param cost body models.SkillImprovementCost true "Updated cost data"
// @Success 200 {object} models.SkillImprovementCost "Updated cost entry"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/skill-improvement-cost2/{id} [put]
func UpdateSkillImprovementCost2(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var cost models.SkillImprovementCost
	if err := database.DB.First(&cost, uint(id)).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "Skill improvement cost not found")
		return
	}

	var payload skillImprovementUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if payload.CurrentLevel != nil {
		cost.CurrentLevel = *payload.CurrentLevel
	}
	if payload.TERequired != nil {
		cost.TERequired = *payload.TERequired
	}
	if payload.CategoryID != nil {
		cost.CategoryID = *payload.CategoryID
	}
	if payload.DifficultyID != nil {
		cost.DifficultyID = *payload.DifficultyID
	}

	if err := database.DB.Save(&cost).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update skill improvement cost")
		return
	}

	categoryNames, difficultyNames, err := loadSkillMetadata([]models.SkillImprovementCost{cost})
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update skill improvement cost")
		return
	}

	c.JSON(http.StatusOK, skillImprovementCostResponse{
		SkillImprovementCost: cost,
		CategoryName:         categoryNames[cost.CategoryID],
		DifficultyName:       difficultyNames[cost.DifficultyID],
	})
}

type skillImprovementCostResponse struct {
	models.SkillImprovementCost
	CategoryName   string `json:"category_name,omitempty"`
	DifficultyName string `json:"difficulty_name,omitempty"`
}

func loadSkillMetadata(costs []models.SkillImprovementCost) (map[uint]string, map[uint]string, error) {
	categoryNames := make(map[uint]string)
	difficultyNames := make(map[uint]string)

	if len(costs) == 0 {
		return categoryNames, difficultyNames, nil
	}

	categoryIDs := make([]uint, 0)
	difficultyIDs := make([]uint, 0)
	seenCategories := make(map[uint]struct{})
	seenDifficulties := make(map[uint]struct{})

	for _, cost := range costs {
		if _, ok := seenCategories[cost.CategoryID]; !ok {
			seenCategories[cost.CategoryID] = struct{}{}
			categoryIDs = append(categoryIDs, cost.CategoryID)
		}
		if _, ok := seenDifficulties[cost.DifficultyID]; !ok {
			seenDifficulties[cost.DifficultyID] = struct{}{}
			difficultyIDs = append(difficultyIDs, cost.DifficultyID)
		}
	}

	if len(categoryIDs) > 0 {
		var categories []models.SkillCategory
		if err := database.DB.Where("id IN ?", categoryIDs).Find(&categories).Error; err != nil {
			return nil, nil, err
		}
		for _, category := range categories {
			categoryNames[category.ID] = category.Name
		}
	}

	if len(difficultyIDs) > 0 {
		var difficulties []models.SkillDifficulty
		if err := database.DB.Where("id IN ?", difficultyIDs).Find(&difficulties).Error; err != nil {
			return nil, nil, err
		}
		for _, difficulty := range difficulties {
			difficultyNames[difficulty.ID] = difficulty.Name
		}
	}

	return categoryNames, difficultyNames, nil
}
