package character

import (
	"bamort/database"
	"bamort/gsmaster"
	"bamort/models"
	"bamort/skills"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Character Handlers

type LearnRequestStruct struct {
	SkillType string `json:"skillType"`
	Name      string `json:"name"`
	Stufe     int    `json:"stufe"`
}

func respondWithError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func ListCharacters(c *gin.Context) {
	var characters []Char
	var listOfChars []CharList
	if err := database.DB.Find(&characters).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve characters")
		return
	}
	for i := range characters {
		listOfChars = append(listOfChars, CharList{
			BamortBase: models.BamortBase{
				ID:   characters[i].ID,
				Name: characters[i].Name,
			},
			Rasse: characters[i].Rasse,
			Typ:   characters[i].Typ,
			Grad:  characters[i].Grad,
			Owner: "test",
		})
	}
	c.JSON(http.StatusOK, listOfChars)
}

func CreateCharacter(c *gin.Context) {
	var character Char
	if err := c.ShouldBindJSON(&character); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DB.Create(&character).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create character")
		return
	}

	c.JSON(http.StatusCreated, character)
}
func GetCharacter(c *gin.Context) {
	id := c.Param("id")
	var character Char
	err := character.FirstID(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve character")
		return
	}
	feChar := ToFeChar(&character)
	c.JSON(http.StatusOK, feChar)
}
func UpdateCharacter(c *gin.Context) {
	id := c.Param("id")
	var character Char

	// First, find the existing character
	err := character.FirstID(id)
	if err != nil {
		respondWithError(c, http.StatusNotFound, "Character not found")
		return
	}

	// Bind the updated data
	if err := c.ShouldBindJSON(&character); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Save the updated character
	if err := database.DB.Save(&character).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update character")
		return
	}

	c.JSON(http.StatusOK, character)
}
func DeleteCharacter(c *gin.Context) {
	id := c.Param("id")
	var character Char
	err := character.FirstID(id)
	if err != nil {
		respondWithError(c, http.StatusNotFound, "Character not found")
		return
	}
	err = character.Delete()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to delete character")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Character deleted successfully"})
}

// Add Fertigkeit by putting it directly to the DB
func AddFertigkeit(charID uint, fertigkeit *skills.Fertigkeit) error {
	// Set the foreign key for the new Eigenschaft
	fertigkeit.CharacterID = charID

	// Save the new Eigenschaft to the database
	if err := database.DB.Create(&fertigkeit).Error; err != nil {
		return fmt.Errorf("failed to add Eigenschaft: %w", err)
	}
	return nil
}

// Append the new Fertigkeit to the slice of the characters property
//character.Fertigkeiten = append(character.Fertigkeiten, fertigkeit)

func ToFeChar(object *Char) *FeChar {
	feC := &FeChar{
		Char: *object,
	}
	skills, innateSkills, categories := splitSkills(object.Fertigkeiten)
	feC.Fertigkeiten = skills
	feC.InnateSkills = innateSkills
	feC.CategorizedSkills = categories
	return feC
}

func splitSkills(object []skills.Fertigkeit) ([]skills.Fertigkeit, []skills.Fertigkeit, map[string][]skills.Fertigkeit) {
	var normSkills []skills.Fertigkeit
	var innateSkills []skills.Fertigkeit
	//var categories map[string][]skills.Fertigkeit
	categories := make(map[string][]skills.Fertigkeit)
	for _, skill := range object {
		gsmsk := skill.GetGsm()
		if gsmsk.Improvable {
			category := "Unkategorisiert"
			if gsmsk.ID != 0 && gsmsk.Category != "" {
				category = gsmsk.Category
			}
			normSkills = append(normSkills, skill)
			if _, exists := categories[category]; !exists {
				categories[category] = make([]skills.Fertigkeit, 0)
			}
			categories[category] = append(categories[category], skill)
		} else {
			innateSkills = append(innateSkills, skill)
		}
	}

	return normSkills, innateSkills, categories
}

func GetLearnSkillCost(c *gin.Context) {
	// Get the character ID from the request
	charID := c.Param("id")

	// Load the character from the database
	var character Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve character")
		return
	}

	// Load the skill from the request
	var s skills.Fertigkeit
	if err := c.ShouldBindJSON(&s); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	var skill gsmaster.Skill
	if err := skill.First(s.Name); err != nil {
		respondWithError(c, http.StatusBadRequest, "can not find speel in gsmaster: "+err.Error())
		return
	}

	cost, err := gsmaster.CalculateSkillLearnCost(skill.Name, character.Typ)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "error getting costs to learn skill: "+err.Error())
		return
	}

	// Return the updated character
	c.JSON(http.StatusOK, cost)
}

func GetLearnSpellCost(c *gin.Context) {
	// Get the character ID from the request
	charID := c.Param("id")

	// Load the character from the database
	var character Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve character")
		return
	}

	// Load the spell from the request
	var s skills.Zauber
	if err := c.ShouldBindJSON(&s); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	var spell gsmaster.Spell
	if err := spell.First(s.Name); err != nil {
		respondWithError(c, http.StatusBadRequest, "can not find speel in gsmaster: "+err.Error())
		return
	}
	sd := gsmaster.SpellDefinition{
		Name:   spell.Name,
		Stufe:  spell.Stufe,
		School: spell.Category,
	}

	cost, err := gsmaster.CalculateSpellLearnCost(spell.Name, character.Typ)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "error getting costs to learn spell: "+err.Error())
		return
	}

	sd.CostEP = cost
	// Return the updated character
	c.JSON(http.StatusOK, sd)
}

func GetLearnCost(c *gin.Context) {
	// Get the character ID from the request
	charID := c.Param("id")

	// Load the character from the database
	var character Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve character")
		return
	}

	// Load the spell from the request
	var s LearnRequestStruct
	if err := c.ShouldBindJSON(&s); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	if s.SkillType != "spell" && s.SkillType != "skill" {
		respondWithError(c, http.StatusBadRequest, "unknown skill type")
		return
	}
	if s.Name == "" {
		respondWithError(c, http.StatusBadRequest, "no name given")
	}
	if s.SkillType == "skill" && s.Stufe <= 6 {
		respondWithError(c, http.StatusBadRequest, "stufe must be greater than 6")
	}

	cost, err := gsmaster.CalculateLearnCost(s.SkillType, s.Name, character.Typ)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "error getting costs to learn spell: "+err.Error())
		return
	}

	// Return the updated character
	c.JSON(http.StatusOK, cost)
}

func GetSkillNextLevelCosts(c *gin.Context) {
	// Get the character ID from the request
	charID := c.Param("id")

	// Load the character from the database
	var character Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve character")
		return
	}

	for int, skill := range character.Fertigkeiten {
		lCost, err := gsmaster.CalculateSkillImprovementCost(skill.Name, character.Typ, skill.Fertigkeitswert)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "error getting costs to learn skill: "+err.Error())
			return
		}
		character.Fertigkeiten[int].LearningCost = *lCost
	}

	// Load the skill from the request
	var s skills.Fertigkeit
	if err := c.ShouldBindJSON(&s); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Return the updated character
	c.JSON(http.StatusOK, character.Fertigkeiten)
}

func GetSkillAllLevelCosts(c *gin.Context) {
	// Get the character ID from the request
	charID := c.Param("id")

	// Load the character from the database
	var character Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve character")
		return
	}
	var s LearnRequestStruct
	if err := c.ShouldBindJSON(&s); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	if s.Name == "" {
		respondWithError(c, http.StatusBadRequest, "no name given")
	}

	costArr := make([]gsmaster.LearnCost, 0)
	notfound := true
	for _, skill := range character.Fertigkeiten {
		if skill.Name == s.Name {
			for i := skill.Fertigkeitswert; i <= 20; i++ {
				lCost, err := gsmaster.CalculateSkillImprovementCost(skill.Name, character.Typ, skill.Fertigkeitswert)
				if err != nil {
					respondWithError(c, http.StatusBadRequest, "error getting costs to learn skill: "+err.Error())
					return
				}
				costArr = append(costArr, *lCost)
			}
			notfound = false
			break
		}
	}
	if notfound {
		for _, skill := range character.Waffenfertigkeiten {
			if skill.Name == s.Name {
				for i := skill.Fertigkeitswert; i <= 20; i++ {
					lCost, err := gsmaster.CalculateSkillImprovementCost(skill.Name, character.Typ, skill.Fertigkeitswert)
					if err != nil {
						respondWithError(c, http.StatusBadRequest, "error getting costs to learn skill: "+err.Error())
						return
					}
					costArr = append(costArr, *lCost)
				}
				break
			}
		}
	}

	// Return the updated character
	c.JSON(http.StatusOK, costArr)
}
