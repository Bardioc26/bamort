package database

func MigrateStructure() error {
	err := DB.AutoMigrate()
	if err != nil {
		return err
	}
	return nil
}
