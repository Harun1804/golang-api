package routes

import (
	"galaxy/backend-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	//initialize gin
	router := gin.Default()

	// route register
	router.POST("/api/register", controllers.Register)

	return router

}