package main

import (
	"galaxy/backend-api/config"
	"galaxy/backend-api/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables from .env file
	config.LoadEnv()

	// Initialize the database connection
	database.InitDB()

	app := gin.Default()
	app.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	app.Run(":" + config.GetEnv("APP_PORT", "3000"))
}