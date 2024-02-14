package main
import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	// "net/http"
	"github.com/NamrithaGirish/hack4tkm/controllers"
	"github.com/NamrithaGirish/hack4tkm/utils"
 )
 
 func main() {
	gin.SetMode(gin.ReleaseMode) //optional to not get warning
	//route.SetTrustedProxies([]string{"192.168.1.2"}) //to trust only a specific value
	utils.ConnectDB()
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST","PUT"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "ngrok-skip-browser-warning"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12, // Maximum age in seconds
	}))
	// router.Use(sessions.Sessions())  //Add password check sessions

	router.POST("/add-user", controllers.AddUser)
	router.GET("/profile/:id", controllers.GetUserById)
	router.GET("/all-teams", controllers.GetAllTeams)
	router.GET("/all-members/:team", controllers.GetTeamMembers)
	router.POST("/comments", controllers.AddComment)
	router.GET("/comments/:id", controllers.DisplayComments)
	router.GET("/leaderboard", controllers.Leaderboard)
	router.GET("/delete-comment/:comment_id", controllers.DeleteComment)
	router.GET("/delete-user/:user_id", controllers.DeleteUser)
	router.GET("/login/:mail", controllers.GetUserByMail)
	router.PUT("/update/:id",controllers.UpdatePhoto)

	router.Run(":8000")
 
 }