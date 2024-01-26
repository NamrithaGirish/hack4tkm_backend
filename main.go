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
	router.POST("/add-question", controllers.AddQuestion)
	router.POST("/add-user", controllers.AddUser)
	router.POST("/add-answer", controllers.CheckAnswer)
	router.GET("/profile/:id", controllers.GetUserById)
	router.GET("/other-users/:id", controllers.GetUsersList)
	router.GET("/question/:id", controllers.GetQuestion)
	router.GET("/answers/:uid/:qid", controllers.GetAnswers)


	router.Run(":8080")
	
 
 }