package gsmaster

import (
	"bamort/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMasterData(c *gin.Context) {
	type dtaStruct struct {
		Skills          []Skill       `json:"skills"`
		Weaponskills    []WeaponSkill `json:"weaponskills"`
		Spell           []Spell       `json:"spells"`
		Equipment       []Equipment   `json:"equipment"`
		Weapons         []Weapon      `json:"weapons"`
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
	c.JSON(http.StatusOK, dta)
}

func GetMDSkills(c *gin.Context) {
	type dtaStruct struct {
		Skills          []Skill       `json:"skills"`
		Weaponskills    []WeaponSkill `json:"weaponskills"`
		Spell           []Spell       `json:"spells"`
		Equipment       []Equipment   `json:"equipment"`
		Weapons         []Weapon      `json:"weapons"`
		SkillCategories []string      `json:"skillcategories"`
	}
	var dta dtaStruct
	if err := database.DB.Find(&dta.Skills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve characters"})
		return
	}
	c.JSON(http.StatusOK, dta.Skills)
}

func GetMDSkill(c *gin.Context) {
	type dtaStruct struct {
		Skills          []Skill       `json:"skills"`
		Weaponskills    []WeaponSkill `json:"weaponskills"`
		Spell           []Spell       `json:"spells"`
		Equipment       []Equipment   `json:"equipment"`
		Weapons         []Weapon      `json:"weapons"`
		SkillCategories []string      `json:"skillcategories"`
	}
	var dta dtaStruct
	if err := database.DB.Find(&dta.Skills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve characters"})
		return
	}
	c.JSON(http.StatusOK, dta)
}

func UpdateMDSkill(c *gin.Context) {
	type dtaStruct struct {
		Skills          []Skill       `json:"skills"`
		Weaponskills    []WeaponSkill `json:"weaponskills"`
		Spell           []Spell       `json:"spells"`
		Equipment       []Equipment   `json:"equipment"`
		Weapons         []Weapon      `json:"weapons"`
		SkillCategories []string      `json:"skillcategories"`
	}
	var dta dtaStruct
	if err := database.DB.Find(&dta.Skills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve characters"})
		return
	}
	c.JSON(http.StatusOK, dta)
}

func AddSkill(c *gin.Context) {
	type dtaStruct struct {
		Skills          []Skill       `json:"skills"`
		Weaponskills    []WeaponSkill `json:"weaponskills"`
		Spell           []Spell       `json:"spells"`
		Equipment       []Equipment   `json:"equipment"`
		Weapons         []Weapon      `json:"weapons"`
		SkillCategories []string      `json:"skillcategories"`
	}
	var dta dtaStruct
	if err := database.DB.Find(&dta.Skills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve characters"})
		return
	}
	c.JSON(http.StatusOK, dta)
}

func DeleteMDSkill(c *gin.Context) {
	type dtaStruct struct {
		Skills          []Skill       `json:"skills"`
		Weaponskills    []WeaponSkill `json:"weaponskills"`
		Spell           []Spell       `json:"spells"`
		Equipment       []Equipment   `json:"equipment"`
		Weapons         []Weapon      `json:"weapons"`
		SkillCategories []string      `json:"skillcategories"`
	}
	var dta dtaStruct
	if err := database.DB.Find(&dta.Skills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve characters"})
		return
	}
	c.JSON(http.StatusOK, dta)
}
