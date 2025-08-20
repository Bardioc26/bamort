package importer

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	charGrp := r.Group("/importer")
	charGrp.POST("/upload", UploadFiles)
	charGrp.POST("/spells/csv", ImportSpellCSVHandler)
}
