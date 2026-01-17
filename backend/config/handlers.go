package config

import (
	"github.com/gin-gonic/gin"
)

// Versionsinfo returns version information
func Versionsinfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"version": GetVersion(),
	})
}
