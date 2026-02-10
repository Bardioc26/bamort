package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Helper functions
func respondWithError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func parseID(c *gin.Context) (uint, error) {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return uint(intID), nil
}

func resolveGameSystem(c *gin.Context) (*models.GameSystem, bool) {
	gsIDStr := c.Query("game_system_id")
	var gsID uint
	if gsIDStr != "" {
		id, err := strconv.ParseUint(gsIDStr, 10, 32)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid game_system_id")
			return nil, false
		}
		gsID = uint(id)
	}

	gsName := c.Query("game_system")
	gs := models.GetGameSystem(gsID, gsName)
	if gs == nil {
		respondWithError(c, http.StatusBadRequest, "Invalid game system")
		return nil, false
	}

	return gs, true
}

type Creator interface {
	Create() error
}

type Saver interface {
	Save() error
}

type FirstIdGetter interface {
	FirstId(uint) error
}

// Add interface
type Deleter interface {
	Delete() error
}

// Generic get single item handler
func getMDItem[T any](c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}
	item := new(T)
	if getter, ok := (interface{})(item).(FirstIdGetter); ok {
		if err := getter.FirstId(id); err != nil {
			respondWithError(c, http.StatusNotFound, "Item not found")
			return
		}
		c.JSON(http.StatusOK, item)
	} else {
		respondWithError(c, http.StatusInternalServerError, "Item type does not support ID lookup")
	}
}

// Generic get all items handler
func getMDItems[T any](c *gin.Context) {
	var items []T
	if err := database.DB.Find(&items).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve items")
		return
	}
	c.JSON(http.StatusOK, items)
}

// Generic add handler
func addMDItem[T any](c *gin.Context) {
	item := new(T)

	if err := c.ShouldBindJSON(item); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if creator, ok := (interface{})(item).(Creator); ok {
		if err := creator.Create(); err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to create item: "+err.Error())
			return
		}
		c.JSON(http.StatusCreated, item)
	} else {
		respondWithError(c, http.StatusInternalServerError, "Item type does not support creation")
	}
}

// Generic update handler
func updateMDItem[T any](c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	item := new(T)
	if getter, ok := (interface{})(item).(FirstIdGetter); ok {
		if err := getter.FirstId(id); err != nil {
			respondWithError(c, http.StatusNotFound, "Item not found")
			return
		}

		if err := c.ShouldBindJSON(item); err != nil {
			respondWithError(c, http.StatusBadRequest, "Invalid input data")
			return
		}

		if saver, ok := (interface{})(item).(Saver); ok {
			if err := saver.Save(); err != nil {
				respondWithError(c, http.StatusInternalServerError, "Failed to update item: "+err.Error())
				return
			}
			c.JSON(http.StatusOK, item)
		} else {
			respondWithError(c, http.StatusInternalServerError, "Item type does not support saving")
		}
	} else {
		respondWithError(c, http.StatusInternalServerError, "Item type does not support ID lookup")
	}
}

// Generic delete handler
func deleteMDItem[T any](c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	item := new(T)
	if getter, ok := (interface{})(item).(FirstIdGetter); ok {
		if err := getter.FirstId(id); err != nil {
			respondWithError(c, http.StatusNotFound, "Item not found")
			return
		}

		if deleter, ok := (interface{})(item).(Deleter); ok {
			if err := deleter.Delete(); err != nil {
				respondWithError(c, http.StatusInternalServerError, "Failed to delete item: "+err.Error())
				return
			}
			c.JSON(http.StatusNoContent, nil)
		} else {
			respondWithError(c, http.StatusInternalServerError, "Item type does not support deletion")
		}
	} else {
		respondWithError(c, http.StatusInternalServerError, "Item type does not support ID lookup")
	}
}

// GetMasterData godoc
// @Summary Get all master data
// @Description Returns a summary of all game system master data (skills, spells, equipment, weapons)
// @Tags Master Data
// @Produce json
// @Success 200 {object} object "Master data summary"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /api/maintenance [get]
func GetMasterData(c *gin.Context) {
	type dtaStruct struct {
		Skills          []models.Skill       `json:"skills"`
		Weaponskills    []models.WeaponSkill `json:"weaponskills"`
		Spell           []models.Spell       `json:"spells"`
		Equipment       []models.Equipment   `json:"equipment"`
		Weapons         []models.Weapon      `json:"weapons"`
		SkillCategories []string             `json:"skillcategories"`
		SpellCategories []string             `json:"spellcategories"`
		Sources         []models.Source      `json:"sources"`
	}
	var dta dtaStruct
	var err error
	var ski models.Skill
	var spe models.Spell
	if err := database.DB.Find(&dta.Skills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Skills"})
		return
	}
	if err := database.DB.Find(&dta.Weaponskills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Weaponskills"})
		return
	}
	if err := database.DB.Find(&dta.Spell).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Spell"})
		return
	}
	if err := database.DB.Find(&dta.Equipment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Equipment"})
		return
	}
	if err := database.DB.Find(&dta.Weapons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Weapons"})
		return
	}
	if err := database.DB.Find(&dta.Sources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Sources"})
		return
	}
	dta.SkillCategories, err = ski.GetSkillCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve SkillCategories" + err.Error()})
		return
	}
	dta.SpellCategories, err = spe.GetSpellCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve SpellCategories" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, dta)
}

// GetMDSkills godoc
// @Summary Get all skills
// @Description Returns list of all game system skills
// @Tags Master Data
// @Produce json
// @Success 200 {array} models.SkFertigkeit "List of skills"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /api/maintenance/skills [get]
func GetMDSkills(c *gin.Context) {
	type dtaStruct struct {
		Skills          []models.Skill       `json:"skills"`
		Weaponskills    []models.WeaponSkill `json:"weaponskills"`
		SkillCategories []string             `json:"skillcategories"`
	}
	var dta dtaStruct
	var err error
	var ski models.Skill
	if err := database.DB.Find(&dta.Skills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Skills"})
		return
	}
	if err := database.DB.Find(&dta.Weaponskills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Weaponskills"})
		return
	}
	dta.SkillCategories, err = ski.GetSkillCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve SkillCategories" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, dta)
}

// GetMDSkill godoc
// @Summary Get skill by ID
// @Description Returns a specific skill by ID
// @Tags Master Data
// @Produce json
// @Param id path int true "Skill ID"
// @Success 200 {object} models.SkFertigkeit "Skill data"
// @Failure 404 {object} map[string]string "Skill not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /api/maintenance/skills/{id} [get]
func GetMDSkill(c *gin.Context) {
	getMDItem[models.Skill](c)
}

// UpdateMDSkill godoc
// @Summary Update skill
// @Description Updates an existing skill (maintainer only)
// @Tags Master Data Admin
// @Accept json
// @Produce json
// @Param id path int true "Skill ID"
// @Param skill body models.SkFertigkeit true "Updated skill data"
// @Success 200 {object} models.SkFertigkeit "Updated skill"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/skills/{id} [put]
func UpdateMDSkill(c *gin.Context) {
	updateMDItem[models.Skill](c)
}

// AddSkill godoc
// @Summary Add new skill
// @Description Creates a new skill in the game system (maintainer only)
// @Tags Master Data Admin
// @Accept json
// @Produce json
// @Param skill body models.SkFertigkeit true "New skill data"
// @Success 201 {object} models.SkFertigkeit "Created skill"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/skills [post]
func AddSkill(c *gin.Context) {
	addMDItem[models.Skill](c)
}

// DeleteMDSkill godoc
// @Summary Delete skill
// @Description Deletes a skill from the game system (maintainer only)
// @Tags Master Data Admin
// @Produce json
// @Param id path int true "Skill ID"
// @Success 200 {object} map[string]string "Skill deleted successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Failure 404 {object} map[string]string "Skill not found"
// @Security BearerAuth
// @Router /api/maintenance/skills/{id} [delete]
func DeleteMDSkill(c *gin.Context) {
	deleteMDItem[models.Skill](c)
}

//

// GetMDWeaponSkills godoc
// @Summary Get all weapon skills
// @Description Returns list of all weapon skills
// @Tags Master Data
// @Produce json
// @Success 200 {array} models.SkWaffenfertigkeit "List of weapon skills"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /api/maintenance/weaponskills [get]
func GetMDWeaponSkills(c *gin.Context) {
	type dtaStruct struct {
		Weaponskills []models.WeaponSkill `json:"weaponskills"`
	}
	var dta dtaStruct
	if err := database.DB.Find(&dta.Weaponskills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Weaponskills"})
		return
	}
	c.JSON(http.StatusOK, dta)
}

// GetMDWeaponSkill godoc
// @Summary Get weapon skill by ID
// @Description Returns a specific weapon skill by ID
// @Tags Master Data
// @Produce json
// @Param id path int true "Weapon skill ID"
// @Success 200 {object} models.SkWaffenfertigkeit "Weapon skill data"
// @Failure 404 {object} map[string]string "Weapon skill not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /api/maintenance/weaponskills/{id} [get]
func GetMDWeaponSkill(c *gin.Context) {
	getMDItem[models.WeaponSkill](c)
}

// UpdateMDWeaponSkill godoc
// @Summary Update weapon skill
// @Description Updates an existing weapon skill (maintainer only)
// @Tags Master Data Admin
// @Accept json
// @Produce json
// @Param id path int true "Weapon skill ID"
// @Param weaponskill body models.SkWaffenfertigkeit true "Updated weapon skill data"
// @Success 200 {object} models.SkWaffenfertigkeit "Updated weapon skill"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/weaponskills/{id} [put]
func UpdateMDWeaponSkill(c *gin.Context) {
	updateMDItem[models.WeaponSkill](c)
}

// AddWeaponSkill godoc
// @Summary Add new weapon skill
// @Description Creates a new weapon skill in the game system (maintainer only)
// @Tags Master Data Admin
// @Accept json
// @Produce json
// @Param weaponskill body models.SkWaffenfertigkeit true "New weapon skill data"
// @Success 201 {object} models.SkWaffenfertigkeit "Created weapon skill"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/weaponskills [post]
func AddWeaponSkill(c *gin.Context) {
	addMDItem[models.WeaponSkill](c)
}

// DeleteMDWeaponSkill godoc
// @Summary Delete weapon skill
// @Description Deletes a weapon skill from the game system (maintainer only)
// @Tags Master Data Admin
// @Produce json
// @Param id path int true "Weapon skill ID"
// @Success 200 {object} map[string]string "Weapon skill deleted successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Failure 404 {object} map[string]string "Weapon skill not found"
// @Security BearerAuth
// @Router /api/maintenance/weaponskills/{id} [delete]
func DeleteMDWeaponSkill(c *gin.Context) {
	deleteMDItem[models.WeaponSkill](c)
}

//

func GetMDSkillCategories(c *gin.Context) {
	var ski models.Skill
	skillCategories, err := ski.GetSkillCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve SkillCategories" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, skillCategories)
}

// GetMDSpells godoc
// @Summary Get all spells
// @Description Returns list of all game system spells
// @Tags Master Data
// @Produce json
// @Success 200 {array} models.Spell "List of spells"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /api/maintenance/spells [get]
func GetMDSpells(c *gin.Context) {
	getMDItems[models.Spell](c)
}

// GetMDSpell godoc
// @Summary Get spell by ID
// @Description Returns a specific spell by ID
// @Tags Master Data
// @Produce json
// @Param id path int true "Spell ID"
// @Success 200 {object} models.Spell "Spell data"
// @Failure 404 {object} map[string]string "Spell not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /api/maintenance/spells/{id} [get]
func GetMDSpell(c *gin.Context) {
	getMDItem[models.Spell](c)
}

// UpdateMDSpell godoc
// @Summary Update spell
// @Description Updates an existing spell (maintainer only)
// @Tags Master Data Admin
// @Accept json
// @Produce json
// @Param id path int true "Spell ID"
// @Param spell body models.Spell true "Updated spell data"
// @Success 200 {object} models.Spell "Updated spell"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/spells/{id} [put]
func UpdateMDSpell(c *gin.Context) {
	updateMDItem[models.Spell](c)
}

// AddSpell godoc
// @Summary Add new spell
// @Description Creates a new spell in the game system (maintainer only)
// @Tags Master Data Admin
// @Accept json
// @Produce json
// @Param spell body models.Spell true "New spell data"
// @Success 201 {object} models.Spell "Created spell"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/spells [post]
func AddSpell(c *gin.Context) {
	addMDItem[models.Spell](c)
}

// DeleteMDSpell godoc
// @Summary Delete spell
// @Description Deletes a spell from the game system (maintainer only)
// @Tags Master Data Admin
// @Produce json
// @Param id path int true "Spell ID"
// @Success 200 {object} map[string]string "Spell deleted successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Failure 404 {object} map[string]string "Spell not found"
// @Security BearerAuth
// @Router /api/maintenance/spells/{id} [delete]
func DeleteMDSpell(c *gin.Context) {
	deleteMDItem[models.Spell](c)
}

// GetMDEquipments godoc
// @Summary Get all equipment
// @Description Returns list of all game system equipment items
// @Tags Master Data
// @Produce json
// @Success 200 {array} models.Equipment "List of equipment"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /api/maintenance/equipment [get]
func GetMDEquipments(c *gin.Context) {
	gs, ok := resolveGameSystem(c)
	if !ok {
		return
	}

	var equipments []models.Equipment
	if err := database.DB.Where("game_system=? OR game_system_id=?", gs.Name, gs.ID).Find(&equipments).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve items")
		return
	}

	c.JSON(http.StatusOK, equipments)
}

// GetMDEquipment godoc
// @Summary Get equipment by ID
// @Description Returns a specific equipment item by ID
// @Tags Master Data
// @Produce json
// @Param id path int true "Equipment ID"
// @Success 200 {object} models.Equipment "Equipment data"
// @Failure 404 {object} map[string]string "Equipment not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /api/maintenance/equipment/{id} [get]
func GetMDEquipment(c *gin.Context) {
	gs, ok := resolveGameSystem(c)
	if !ok {
		return
	}

	id, err := parseID(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	equipment := &models.Equipment{GameSystem: gs.Name, GameSystemId: gs.ID}
	if err := equipment.FirstId(id); err != nil {
		respondWithError(c, http.StatusNotFound, "Item not found")
		return
	}

	c.JSON(http.StatusOK, equipment)
}

// UpdateMDEquipment godoc
// @Summary Update equipment item
// @Description Updates an existing equipment item (maintainer only)
// @Tags Master Data Admin
// @Accept json
// @Produce json
// @Param id path int true "Equipment ID"
// @Param equipment body models.Equipment true "Updated equipment data"
// @Success 200 {object} models.Equipment "Updated equipment"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/equipment/{id} [put]
func UpdateMDEquipment(c *gin.Context) {
	gs, ok := resolveGameSystem(c)
	if !ok {
		return
	}

	id, err := parseID(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	equipment := &models.Equipment{GameSystem: gs.Name, GameSystemId: gs.ID}
	if err := equipment.FirstId(id); err != nil {
		respondWithError(c, http.StatusNotFound, "Item not found")
		return
	}

	if err := c.ShouldBindJSON(equipment); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	equipment.ID = id
	equipment.GameSystem = gs.Name
	equipment.GameSystemId = gs.ID

	if err := equipment.Save(); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update item: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, equipment)
}

// AddEquipment godoc
// @Summary Add new equipment
// @Description Creates a new equipment item in the game system (maintainer only)
// @Tags Master Data Admin
// @Accept json
// @Produce json
// @Param equipment body models.Equipment true "New equipment data"
// @Success 201 {object} models.Equipment "Created equipment"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/equipment [post]
func AddEquipment(c *gin.Context) {
	gs, ok := resolveGameSystem(c)
	if !ok {
		return
	}

	equipment := &models.Equipment{GameSystem: gs.Name, GameSystemId: gs.ID}
	if err := c.ShouldBindJSON(equipment); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if equipment.GameSystemId == 0 && equipment.GameSystem == "" {
		equipment.GameSystem = gs.Name
		equipment.GameSystemId = gs.ID
	}

	if err := equipment.Create(); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create item: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, equipment)
}

// DeleteMDEquipment godoc
// @Summary Delete equipment
// @Description Deletes an equipment item from the game system (maintainer only)
// @Tags Master Data Admin
// @Produce json
// @Param id path int true "Equipment ID"
// @Success 200 {object} map[string]string "Equipment deleted successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Failure 404 {object} map[string]string "Equipment not found"
// @Security BearerAuth
// @Router /api/maintenance/equipment/{id} [delete]
func DeleteMDEquipment(c *gin.Context) {
	gs, ok := resolveGameSystem(c)
	if !ok {
		return
	}

	id, err := parseID(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if err := database.DB.Where("(game_system=? OR game_system_id=?) AND id = ?", gs.Name, gs.ID, id).Delete(&models.Equipment{}).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to delete item")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Refactored handler functions
// GetMDWeapons godoc
// @Summary Get all weapons
// @Description Returns list of all game system weapons
// @Tags Master Data
// @Produce json
// @Success 200 {array} models.Weapon "List of weapons"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /api/maintenance/weapons [get]
func GetMDWeapons(c *gin.Context) {
	gs, ok := resolveGameSystem(c)
	if !ok {
		return
	}

	var weapons []models.Weapon
	if err := database.DB.Where("game_system=? OR game_system_id=?", gs.Name, gs.ID).Find(&weapons).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve items")
		return
	}

	c.JSON(http.StatusOK, weapons)
}

// GetMDWeapon godoc
// @Summary Get weapon by ID
// @Description Returns a specific weapon by ID
// @Tags Master Data
// @Produce json
// @Param id path int true "Weapon ID"
// @Success 200 {object} models.Weapon "Weapon data"
// @Failure 404 {object} map[string]string "Weapon not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /api/maintenance/weapons/{id} [get]
func GetMDWeapon(c *gin.Context) {
	gs, ok := resolveGameSystem(c)
	if !ok {
		return
	}

	id, err := parseID(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	weapon := &models.Weapon{Equipment: models.Equipment{GameSystem: gs.Name, GameSystemId: gs.ID}}
	if err := weapon.FirstId(id); err != nil {
		respondWithError(c, http.StatusNotFound, "Item not found")
		return
	}

	c.JSON(http.StatusOK, weapon)
}

// UpdateMDWeapon godoc
// @Summary Update weapon
// @Description Updates an existing weapon (maintainer only)
// @Tags Master Data Admin
// @Accept json
// @Produce json
// @Param id path int true "Weapon ID"
// @Param weapon body models.Weapon true "Updated weapon data"
// @Success 200 {object} models.Weapon "Updated weapon"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/weapons/{id} [put]
func UpdateMDWeapon(c *gin.Context) {
	gs, ok := resolveGameSystem(c)
	if !ok {
		return
	}

	id, err := parseID(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	weapon := &models.Weapon{Equipment: models.Equipment{GameSystem: gs.Name, GameSystemId: gs.ID}}
	if err := weapon.FirstId(id); err != nil {
		respondWithError(c, http.StatusNotFound, "Item not found")
		return
	}

	if err := c.ShouldBindJSON(weapon); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	weapon.ID = id
	weapon.GameSystem = gs.Name
	weapon.GameSystemId = gs.ID

	if err := weapon.Save(); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update item: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, weapon)
}

// AddWeapon godoc
// @Summary Add new weapon
// @Description Creates a new weapon in the game system (maintainer only)
// @Tags Master Data Admin
// @Accept json
// @Produce json
// @Param weapon body models.Weapon true "New weapon data"
// @Success 201 {object} models.Weapon "Created weapon"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Security BearerAuth
// @Router /api/maintenance/weapons [post]
func AddWeapon(c *gin.Context) {
	gs, ok := resolveGameSystem(c)
	if !ok {
		return
	}

	weapon := &models.Weapon{Equipment: models.Equipment{GameSystem: gs.Name, GameSystemId: gs.ID}}
	if err := c.ShouldBindJSON(weapon); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if weapon.GameSystemId == 0 && weapon.GameSystem == "" {
		weapon.GameSystem = gs.Name
		weapon.GameSystemId = gs.ID
	}

	if err := weapon.Create(); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create item: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, weapon)
}

// DeleteMDWeapon godoc
// @Summary Delete weapon
// @Description Deletes a weapon from the game system (maintainer only)
// @Tags Master Data Admin
// @Produce json
// @Param id path int true "Weapon ID"
// @Success 200 {object} map[string]string "Weapon deleted successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - maintainer access required"
// @Failure 404 {object} map[string]string "Weapon not found"
// @Security BearerAuth
// @Router /api/maintenance/weapons/{id} [delete]
func DeleteMDWeapon(c *gin.Context) {
	gs, ok := resolveGameSystem(c)
	if !ok {
		return
	}

	id, err := parseID(c)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if err := database.DB.Where("(game_system=? OR game_system_id=?) AND id = ?", gs.Name, gs.ID, id).Delete(&models.Weapon{}).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to delete item")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
