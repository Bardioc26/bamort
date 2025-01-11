package gsmaster

import "bamort/database"

func MigrateStructure() error {
	err := database.DB.AutoMigrate(
		&Skill{},
		&WeaponSkill{},
		&Spell{},
		&Equipment{},
		&Weapon{},
		&Container{},
		&Transportation{},
		&Believe{},
	)
	if err != nil {
		return err
	}
	return nil
}
