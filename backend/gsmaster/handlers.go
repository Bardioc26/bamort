package gsmaster

import (
	"bamort/database"
	"fmt"
	"net/http"
	"strconv"

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
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve character"})
		return
	}
	sk := Skill{}
	err = sk.FirstId(uint(intId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve character"})
		return
	}
	c.JSON(http.StatusOK, sk)

}

func UpdateMDSkill(c *gin.Context) {
	var sk Skill
	if err := c.ShouldBindJSON(&sk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if sk.System == "" {
		sk.System = "midgard"
	}
	fmt.Printf("UpdateMDSkill: %v\n", sk)
	if err := sk.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save skill" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, sk)
}

func AddSkill(c *gin.Context) {
	var sk Skill
	if err := c.ShouldBindJSON(&sk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := sk.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save skill" + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sk)
}

func DeleteMDSkill(c *gin.Context) {
	var sk Skill
	if err := c.ShouldBindJSON(&sk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := sk.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save skill" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, sk)
}
