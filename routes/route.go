package routes

import (
	"galaxy/backend-api/controllers"
	"galaxy/backend-api/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	//initialize gin
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
	}))


	api := router.Group("/api")

	// Setup routes for different modules
	setupAuthRoutes(api)
	setupUserRoutes(api)

	return router
}

// setupAuthRoutes handles authentication related routes
func setupAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	auth.POST("/register", controllers.Register)
	auth.POST("/login", controllers.Login)
}

// setupUserRoutes handles user related routes
func setupUserRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	users.Use(middlewares.AuthMiddleware())
	// Add user routes here when needed
	users.GET("/", controllers.GetUsers)
	users.POST("/", controllers.CreateUser)
	users.GET("/:id", controllers.GetUser)
	users.PUT("/:id", controllers.UpdateUser)
	users.DELETE("/:id", controllers.DeleteUser)
}