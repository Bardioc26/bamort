package equipment

import "bamort/database"

func MigrateStructure() error {
	err := database.DB.AutoMigrate(
		&Ausruestung{},
		&Waffe{},
		&Behaeltniss{},
		&Transportation{},
	)
	if err != nil {
		return err
	}
	return nil
}
