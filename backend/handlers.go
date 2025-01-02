package main

import (
	"bamort/character"
	"bamort/database"
	"bamort/equipment"
	"bamort/gsmaster"
	"bamort/importer"
	"bamort/skills"
	"bamort/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupCheck(c *gin.Context) {
	db := database.ConnectDatabase()

	if db != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to DataBase"})
		return
	}
	err := migrateStructure()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to automigrate DataBase"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Setup Check OK"})
}

func migrateStructure() error {
	err := database.MigrateStructure()
	if err != nil {
		return err
	}
	err = character.MigrateStructure()
	if err != nil {
		return err
	}
	err = equipment.MigrateStructure()
	if err != nil {
		return err
	}
	err = gsmaster.MigrateStructure()
	if err != nil {
		return err
	}
	err = importer.MigrateStructure()
	if err != nil {
		return err
	}
	err = skills.MigrateStructure()
	if err != nil {
		return err
	}
	err = user.MigrateStructure()
	if err != nil {
		return err
	}

	return nil
}
