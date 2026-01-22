package main

import (
	"seal-of-ash/internal/config"
	"seal-of-ash/internal/database"
	"seal-of-ash/internal/forensic"
	"seal-of-ash/internal/handlers"
	"seal-of-ash/internal/middleware"
	"seal-of-ash/internal/models"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	database.Init()
	forensic.StartWorker()

	database.ActiveDB.AutoMigrate(
		&models.User{},
		&models.Transaction{},
	)

	database.ForensicDB.AutoMigrate(
		&models.ForensicAsh{},
	)

	database.SeedUsers()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Seal of Ash is alive"})
	})

	r.POST("/auth/login", handlers.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())

	auth.POST("/transactions", handlers.CreateTransaction)

	auth.POST("/ashes/summon",
		middleware.RequireRole("forensic"),
		handlers.SummonAshes,
	)

	auth.GET("/ashes/verify",
		middleware.RequireRole("forensic"),
		handlers.VerifyAshChain,
	)

	r.Run(":" + config.Get("APP_PORT"))
}
