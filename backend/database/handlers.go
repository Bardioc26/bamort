package database

import (
	"github.com/gin-gonic/gin"
)

func SetupCheck(c *gin.Context) {
	ConnectDatabase()
}
