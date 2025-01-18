package character

import (
	"bamort/database"
	"bamort/models"
	"bamort/skills"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Character Handlers

Add CRUD operations for characters:
*/

func ListCharacters(c *gin.Context) {
	var characters []Char
	var listOfChars []CharList
	if err := database.DB.Find(&characters).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve characters"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&character).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create character"})
		return
	}

	c.JSON(http.StatusCreated, character)
}
func GetCharacter(c *gin.Context) {
	id := c.Param("id")
	var character Char
	err := character.FirstID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve character"})
		return
	}
	feChar := ToFeChar(&character)
	c.JSON(http.StatusOK, feChar)
}
func UpdateCharacter(c *gin.Context) {
	var character Char
	/*
		if err := c.ShouldBindJSON(&character.ID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := database.DB.Create(&character).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create character"})
			return
		}
	*/
	c.JSON(http.StatusCreated, character)
}
func DeleteCharacter(c *gin.Context) {
	id := c.Param("id")
	var character Char
	err := character.FirstID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve character"})
		return
	}
	err = character.Delete()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete character"})
		return
	}
	/*
		if err := c.ShouldBindJSON(&character.ID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := database.DB.Create(&character).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create character"})
			return
		}
	*/
	c.JSON(http.StatusCreated, character)
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
