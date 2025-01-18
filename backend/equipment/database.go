package equipment

import "bamort/database"

func MigrateStructure() error {
	err := database.DB.AutoMigrate(
		&Ausruestung{},
		&Waffe{},
		&Container{},
		//&Transportation{},
	)
	if err != nil {
		return err
	}
	return nil
}
