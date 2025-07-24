package maintenance

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	charGrp := r.Group("/maintenance")
	charGrp.GET("/setupcheck", SetupCheck)
	charGrp.GET("/mktestdata", MakeTestdataFromLive)
}
