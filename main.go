package main
import (
	"github.com/gin-gonic/gin"
	// "net/http"
	"github.com/NamrithaGirish/hack4tkm/controllers"
	"github.com/NamrithaGirish/hack4tkm/utils"
 )
 
 func main() {
	gin.SetMode(gin.ReleaseMode) //optional to not get warning
	//route.SetTrustedProxies([]string{"192.168.1.2"}) //to trust only a specific value
	utils.ConnectDB()
	router := gin.Default()
	router.POST("/add-user", controllers.AddUser)
	router.GET("/profile/:id", controllers.GetUserById)
	router.GET("/all-teams", controllers.GetAllTeams)
	router.GET("/all-members/:team", controllers.GetTeamMembers)
	router.POST("/comments", controllers.AddComment)
	router.GET("/comments/:name", controllers.DisplayComments)
	router.GET("/leaderboard", controllers.Leaderboard)

	router.Run(":8080")
 
 }