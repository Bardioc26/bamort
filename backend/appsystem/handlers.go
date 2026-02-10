package appsystem

import (
	"github.com/gin-gonic/gin"
)

// Versionsinfo returns version and git commit information
// Versionsinfo godoc
// @Summary Get application version
// @Description Returns the current application version information
// @Tags System
// @Produce json
// @Success 200 {object} object "Version information"
// @Router /api/version [get]
// @Router /api/public/version [get]
func Versionsinfo(c *gin.Context) {
	c.JSON(200, GetInfo())
}

// SystemInfo godoc
// @Summary Get system information
// @Description Returns system configuration and environment information
// @Tags System
// @Produce json
// @Success 200 {object} object "System information"
// @Router /api/systeminfo [get]
// @Router /api/public/systeminfo [get]
func SystemInfo(c *gin.Context) {
	c.JSON(200, GetInfo2())
}
