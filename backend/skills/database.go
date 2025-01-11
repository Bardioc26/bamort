package skills

import "bamort/database"

func MigrateStructure() error {
	err := database.DB.AutoMigrate(
		&Fertigkeit{},
		&Waffenfertigkeit{},
		&Zauber{},
	)
	if err != nil {
		return err
	}
	return nil
}
