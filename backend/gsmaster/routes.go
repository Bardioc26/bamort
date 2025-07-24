package gsmaster

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	maintGrp := r.Group("/maintenance")
	maintGrp.GET("", GetMasterData)
	maintGrp.GET("/skills", GetMDSkills)
	maintGrp.GET("/skills/:id", GetMDSkill)
	maintGrp.PUT("/skills/:id", UpdateMDSkill)
	maintGrp.POST("/skills", AddSkill)
	maintGrp.DELETE("/skills/:id", DeleteMDSkill)

	maintGrp.GET("/weaponskills", GetMDWeaponSkills)
	maintGrp.GET("/weaponskills/:id", GetMDWeaponSkill)
	maintGrp.PUT("/weaponskills/:id", UpdateMDWeaponSkill)
	maintGrp.POST("/weaponskills", AddWeaponSkill)
	maintGrp.DELETE("/weaponskills/:id", DeleteMDWeaponSkill)

	maintGrp.GET("/spells", GetMDSpells)
	maintGrp.GET("/spells/:id", GetMDSpell)
	maintGrp.PUT("/spells/:id", UpdateMDSpell)
	maintGrp.POST("/spells", AddSpell)
	maintGrp.DELETE("/spells/:id", DeleteMDSpell)

	maintGrp.GET("/equipment", GetMDEquipments)
	maintGrp.GET("/equipment/:id", GetMDEquipment)
	maintGrp.PUT("/equipment/:id", UpdateMDEquipment)
	maintGrp.POST("/equipment", AddEquipment)
	maintGrp.DELETE("/equipment/:id", DeleteMDEquipment)

	maintGrp.GET("/weapons", GetMDWeapons)
	maintGrp.GET("/weapons/:id", GetMDWeapon)
	maintGrp.PUT("/weapons/:id", UpdateMDWeapon)
	maintGrp.POST("/weapons", AddWeapon)
	maintGrp.DELETE("/weapons/:id", DeleteMDWeapon)
}
