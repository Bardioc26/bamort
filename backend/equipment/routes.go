package equipment

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	// Equipment (Ausrüstung) routes
	equipGrp := r.Group("/equipment")
	equipGrp.POST("", CreateAusruestung)
	equipGrp.GET("/character/:character_id", ListAusruestung)
	equipGrp.PUT("/:ausruestung_id", UpdateAusruestung)
	equipGrp.DELETE("/:ausruestung_id", DeleteAusruestung)

	// Weapon (Waffen) routes
	weaponGrp := r.Group("/weapons")
	weaponGrp.POST("", CreateWaffe)
	weaponGrp.GET("/character/:character_id", ListWaffen)
	weaponGrp.PUT("/:waffe_id", UpdateWaffe)
	weaponGrp.DELETE("/:waffe_id", DeleteWaffe)
}
