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

func GetMasterData(c *gin.Context) {
	type dtaStruct struct {
		Skills          []models.Skill       `json:"skills"`
		Weaponskills    []models.WeaponSkill `json:"weaponskills"`
		Spell           []models.Spell       `json:"spells"`
		Equipment       []Equipment          `json:"equipment"`
		Weapons         []Weapon             `json:"weapons"`
		SkillCategories []string             `json:"skillcategories"`
		SpellCategories []string             `json:"spellcategories"`
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

func GetMDSkill(c *gin.Context) {
	getMDItem[models.Skill](c)
}

func UpdateMDSkill(c *gin.Context) {
	updateMDItem[models.Skill](c)
}

func AddSkill(c *gin.Context) {
	addMDItem[models.Skill](c)
}

func DeleteMDSkill(c *gin.Context) {
	deleteMDItem[models.Skill](c)
}

//

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

func GetMDWeaponSkill(c *gin.Context) {
	getMDItem[models.WeaponSkill](c)
}

func UpdateMDWeaponSkill(c *gin.Context) {
	updateMDItem[models.WeaponSkill](c)
}

func AddWeaponSkill(c *gin.Context) {
	addMDItem[models.WeaponSkill](c)
}

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

func GetMDSpells(c *gin.Context) {
	getMDItems[models.Spell](c)
}

func GetMDSpell(c *gin.Context) {
	getMDItem[models.Spell](c)
}

func UpdateMDSpell(c *gin.Context) {
	updateMDItem[models.Spell](c)
}

func AddSpell(c *gin.Context) {
	addMDItem[models.Spell](c)
}

func DeleteMDSpell(c *gin.Context) {
	deleteMDItem[models.Spell](c)
}

func GetMDEquipments(c *gin.Context) {
	getMDItems[Equipment](c)
}

func GetMDEquipment(c *gin.Context) {
	getMDItem[Equipment](c)
}

func UpdateMDEquipment(c *gin.Context) {
	updateMDItem[Equipment](c)
}

func AddEquipment(c *gin.Context) {
	addMDItem[Equipment](c)
}

func DeleteMDEquipment(c *gin.Context) {
	deleteMDItem[Equipment](c)
}

// Refactored handler functions
func GetMDWeapons(c *gin.Context) {
	getMDItems[Weapon](c)
}

func GetMDWeapon(c *gin.Context) {
	getMDItem[Weapon](c)
}

func UpdateMDWeapon(c *gin.Context) {
	updateMDItem[Weapon](c)
}

func AddWeapon(c *gin.Context) {
	addMDItem[Weapon](c)
}

func DeleteMDWeapon(c *gin.Context) {
	deleteMDItem[Weapon](c)
}
