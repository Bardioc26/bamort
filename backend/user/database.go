package user

import "bamort/database"

func MigrateStructure() error {
	err := database.DB.AutoMigrate(
		&User{},
	)
	if err != nil {
		return err
	}
	return nil
}
