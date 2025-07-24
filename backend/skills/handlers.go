package skills

import (
	"github.com/gin-gonic/gin"
)

// Helper function for error responses
func respondWithError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
