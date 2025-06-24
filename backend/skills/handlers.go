package skills

import (
	"bamort/database"

	"github.com/gin-gonic/gin"
)

// MigrateStructure migrates the skills database structure
func MigrateStructure() error {
	return database.DB.AutoMigrate(
		&Fertigkeit{},
		&Zauber{},
		&Waffenfertigkeit{},
	)
}

// Helper function for error responses
func respondWithError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
