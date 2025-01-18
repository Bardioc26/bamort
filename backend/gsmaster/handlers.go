package gsmaster

import (
	"bamort/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Helper functions
func handleError(c *gin.Context, status int, message string) {
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
		handleError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}
	item := new(T)
	if getter, ok := (interface{})(item).(FirstIdGetter); ok {
		if err := getter.FirstId(id); err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to retrieve item")
			return
		}
		c.JSON(http.StatusOK, item)
	}

}

// Generic get all items handler
func getMDItems[T any](c *gin.Context) {
	var items []T
	if err := database.DB.Find(&items).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to retrieve items")
		return
	}
	c.JSON(http.StatusOK, items)
}

// Generic add handler
func addMDItem[T any](c *gin.Context) {
	item := new(T)

	if err := c.ShouldBindJSON(item); err != nil {
		handleError(c, http.StatusBadRequest, err.Error())
		return
	}

	if creator, ok := (interface{})(item).(Creator); ok {
		if err := creator.Create(); err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to create item")
			return
		}
		c.JSON(http.StatusCreated, item)
	}
}

// Generic update handler
func updateMDItem[T any](c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	item := new(T)
	if getter, ok := (interface{})(item).(FirstIdGetter); ok {
		if err := getter.FirstId(id); err != nil {
			handleError(c, http.StatusNotFound, "Item not found")
			return
		}

		if err := c.ShouldBindJSON(item); err != nil {
			handleError(c, http.StatusBadRequest, "Invalid input data")
			return
		}

		if saver, ok := (interface{})(item).(Saver); ok {
			if err := saver.Save(); err != nil {
				handleError(c, http.StatusInternalServerError, "Failed to update item")
				return
			}
			c.JSON(http.StatusOK, item)
		}
	}
}

// Generic delete handler
func deleteMDItem[T any](c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	item := new(T)
	if getter, ok := (interface{})(item).(FirstIdGetter); ok {
		if err := getter.FirstId(id); err != nil {
			handleError(c, http.StatusNotFound, "Item not found")
			return
		}

		if deleter, ok := (interface{})(item).(Deleter); ok {
			if err := deleter.Delete(); err != nil {
				handleError(c, http.StatusInternalServerError, "Failed to delete item")
				return
			}
			c.JSON(http.StatusNoContent, nil)
		}
	}
}

func GetMasterData(c *gin.Context) {
	type dtaStruct struct {
		Skills          []Skill       `json:"skills"`
		Weaponskills    []WeaponSkill `json:"weaponskills"`
		Spell           []Spell       `json:"spells"`
		Equipment       []Equipment   `json:"equipment"`
		Weapons         []Weapon      `json:"weapons"`
		SkillCategories []string      `json:"skillcategories"`
		SpellCategories []string      `json:"spellcategories"`
	}
	var dta dtaStruct
	var err error
	var ski Skill
	var spe Spell
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
		Skills          []Skill       `json:"skills"`
		Weaponskills    []WeaponSkill `json:"weaponskills"`
		SkillCategories []string      `json:"skillcategories"`
	}
	var dta dtaStruct
	var err error
	var ski Skill
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
	getMDItem[Skill](c)
}

func UpdateMDSkill(c *gin.Context) {
	updateMDItem[Skill](c)
}

func AddSkill(c *gin.Context) {
	addMDItem[Skill](c)
}

func DeleteMDSkill(c *gin.Context) {
	deleteMDItem[Skill](c)
}

func GetMDSkillCategories(c *gin.Context) {
	var ski Skill
	skillCategories, err := ski.GetSkillCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve SkillCategories" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, skillCategories)
}

func GetMDSpells(c *gin.Context) {
	getMDItems[Spell](c)
}

func GetMDSpell(c *gin.Context) {
	getMDItem[Spell](c)
}

func UpdateMDSpell(c *gin.Context) {
	updateMDItem[Spell](c)
}

func AddSpell(c *gin.Context) {
	addMDItem[Spell](c)
}

func DeleteMDSpell(c *gin.Context) {
	deleteMDItem[Spell](c)
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
