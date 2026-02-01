package gsmaster

import "bamort/models"

func GetGameSystem(id int, name string) *models.GameSystem {
	gs := &models.GameSystem{}
	if id == 0 && name == "" {
		gs.GetDefault()
		return gs
	}
	if id == 0 && name != "" {
		gs.FirstByCode(name)
		if gs.ID == 0 {
			gs.FirstByName(name)
		}
		return gs
	}
	gs.FirstByID(uint(id))
	return gs
}
