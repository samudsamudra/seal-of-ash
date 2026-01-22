package main

import (
	"seal-of-ash/internal/config"
	"seal-of-ash/internal/database"
	"seal-of-ash/internal/models"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	database.Init()

	database.ActiveDB.AutoMigrate(
		&models.Transaction{},
	)

	database.ForensicDB.AutoMigrate(
		&models.ForensicAsh{},
	)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Seal of Ash is alive"})
	})

	r.Run(":" + config.Get("APP_PORT"))
}
