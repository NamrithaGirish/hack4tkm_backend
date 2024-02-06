package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/NamrithaGirish/hack4tkm/controllers"
)

// InitializeRoutes sets up all routes for the application
func Routes(router *gin.Engine) {
	// router.POST("/add-question", controllers.AddQuestion)
	router.POST("/add-user", controllers.AddUser)
}