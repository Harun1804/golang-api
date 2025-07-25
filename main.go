package main

import (
	"galaxy/backend-api/config"
	"galaxy/backend-api/database"
	"galaxy/backend-api/routes"
)

func main() {
	// Load environment variables from .env file
	config.LoadEnv()

	// Initialize the database connection
	database.InitDB()

	//setup router
	app := routes.SetupRouter()
	app.Run(":" + config.GetEnv("APP_PORT", "3000"))
}