package main

import (
	"galaxy/backend-api/config"
	"galaxy/backend-api/database"
	"galaxy/backend-api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables from .env file
	config.LoadEnv()

	// Initialize the database connection
	database.InitDB()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	//setup router
	routes.SetupRouter(router)
	router.Run(":" + config.GetEnv("APP_PORT", "3000"))
}
