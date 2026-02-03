package appsystem

import (
	"github.com/gin-gonic/gin"
)

// Versionsinfo returns version and git commit information
func Versionsinfo(c *gin.Context) {
	c.JSON(200, GetInfo())
}

func SystemInfo(c *gin.Context) {
	c.JSON(200, GetInfo2())
}
