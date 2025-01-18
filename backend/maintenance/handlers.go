package maintenance

import (
	"bamort/character"
	"bamort/database"
	"bamort/equipment"
	"bamort/gsmaster"
	"bamort/importer"
	"bamort/skills"
	"bamort/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupCheck(c *gin.Context) {
	db := database.ConnectDatabase()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to DataBase"})
		return
	}
	err := database.MigrateStructure()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to automigrate database DataBase"})
		return
	}
	err = user.MigrateStructure()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to automigrate user DataBase"})
		return
	}
	err = character.MigrateStructure()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to automigratec haracter DataBase"})
		return
	}
	err = gsmaster.MigrateStructure()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to automigrate gsmaster DataBase"})
		return
	}
	err = equipment.MigrateStructure()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to automigrate equipment DataBase"})
		return
	}
	err = skills.MigrateStructure()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to automigrate skills DataBase"})
		return
	}
	err = importer.MigrateStructure()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to automigrate importer DataBase"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Setup Check OK"})
}
