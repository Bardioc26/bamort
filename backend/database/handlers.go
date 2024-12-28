package database

import (
	"github.com/gin-gonic/gin"
)

func SetupCheck(c *gin.Context) {
	ConnectDatabase()
	/*
		err := DB.AutoMigrate(&Character{},
			&Fertigkeit{}, &Zauber{}, &Lp{},
			&Eigenschaft{}, &Merkmale{},
			&Bennies{},
			&Gestalt{},
			&Ap{}, &B{},
			&Erfahrungsschatz{},
			&MagischTransport{},
			&Transportation{},
			&MagischAusruestung{},
			&Ausruestung{},
			&MagischBehaelter{},
			&Behaeltniss{},
			&MagischWaffe{},
			&Waffe{},
			&Waffenfertigkeit{},
			&StammFertigkeit{},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to automigrate DataBase"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Setup Check OK"})
	*/
}
