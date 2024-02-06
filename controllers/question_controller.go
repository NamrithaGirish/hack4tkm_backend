package controllers

import (
	"github.com/NamrithaGirish/hack4tkm/models"
	"github.com/gin-gonic/gin"
    "net/http"
    "github.com/NamrithaGirish/hack4tkm/utils"
    "strconv"
	
)

func AddUser(context *gin.Context) {
    var input models.User

    if err := context.ShouldBindJSON(&input); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user := models.User{
        Name: input.Name,
        ID: input.ID,
        Mail: input.Mail,
        Team: input.Team,
    }

    savedUser, err := user.Save()

    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    context.JSON(http.StatusCreated, gin.H{"user": savedUser})
}

// func AddQuestion(context *gin.Context) {
//     var input models.Question

//     if err := context.ShouldBindJSON(&input); err != nil {
//         context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     question := models.Question{
//         Question: input.Question,
//         Answer: input.Answer,
//         UserID: input.UserID,
//         ID: input.ID,
//     }

//     savedQuestion, err := question.Save()

//     if err != nil {
//         context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     context.JSON(http.StatusCreated, gin.H{"question": savedQuestion})
// }

func GetAllTeams(context *gin.Context) {
    var teams []string
    result:=utils.DB.Model(&models.User{}).Select("team").Group("team").Find(&teams)
    if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
	}
	context.JSON(http.StatusOK, gin.H{"teams": teams})
}

func GetTeamMembers(context *gin.Context) {
    var names []string
    team:=context.Param("team")
    result:=utils.DB.Model(&models.User{}).Select("name").Where("team = ?",team).Find(&names)
    if result.Error != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
	}
	context.JSON(http.StatusOK, gin.H{"names": names})
}

func CommentEnable(receiver_id uint, sender_id uint) bool {
    var comment models.Comments
    result:=utils.DB.Where("receiver_id = ? and sender_id = ?",receiver_id,sender_id).First(&comment)
    if result.Error!=nil{
        return true
    } else {
        return false
    }
}

func AddComment(context *gin.Context) {
    
    var input models.Comments

    if err := context.ShouldBindJSON(&input); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if (CommentEnable(input.ReceiverID, input.SenderID)){
        comment := models.Comments{
            Comment: input.Comment,
            LinkedinUrl: input.LinkedinUrl,
            Image:input.Image,
            SenderID: input.SenderID,
            ReceiverID: input.ReceiverID,
            ID: input.ID,
        }
    
        savedComment, err := comment.Save()
    
        if err != nil {
            context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        var user models.User
        if err := utils.DB.First(&user, input.SenderID).Error; err != nil {
            context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
            return
        }
        user.Points = user.Points+10
        utils.DB.Save(&user)
        context.JSON(http.StatusCreated, gin.H{"comment": savedComment})
    } else{
        context.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add more than one comment"})
    }

    
}

func FindUserID(name string) uint{
    var user models.User
    utils.DB.Where("name = ?",name).First(&user)
    // fmt.Println(user.ID)
    return user.ID

}

func DisplayComments(context *gin.Context) {
    var comments []models.Comments
    name:=context.Param("name")
    result:=utils.DB.Where("receiver_id = ?",FindUserID(name)).Find(&comments)

    if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
	}

	context.JSON(http.StatusOK, comments)

}


// func CheckAnswer(context *gin.Context) {
//     var input models.Answer

//     if err := context.ShouldBindJSON(&input); err != nil {
//         context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     answer := models.Answer{
//         Answer: input.Answer,
//         QID: input.QID,
//         UserID: input.UserID,
//         ID: input.ID,
//     }
//     savedAnswer, err := answer.Save()
//     if err != nil {
//         context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }
//     var question models.Question
//     result:=utils.DB.First(&question, input.QID)
//     err = result.Error
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve question"})
// 		return
// 	} 

//     isCorrect := (input.Answer == question.Answer)
//     if isCorrect {
//         var user models.User
//         if err := utils.DB.First(&user, input.UserID).Error; err != nil {
//             context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
//             return
//         } 
//         var other_user models.User
//         utils.DB.First(&other_user, question.UserID)
//         if user.Team != other_user.Team{
//             user.Points = user.Points+10
//             utils.DB.Save(&user)
//             context.JSON(http.StatusCreated, gin.H{"updated answer": savedAnswer,"updated user":user})
//             return
//         }
        
//     }   
//     context.JSON(http.StatusCreated, gin.H{"updated answer": savedAnswer,"remark": "Wrong Answer"})
    
// }

func GetUserById(context *gin.Context) {
    user_id, err := strconv.ParseUint(context.Param("id"), 10, 64)
    if err!=nil{
        context.JSON(http.StatusBadRequest, gin.H{"error":"Incorrect user ID"})
        return
    }
    var user_profile models.User
    result := utils.DB.First(&user_profile, user_id)
    if result.Error!=nil {
        context.JSON(http.StatusBadRequest, gin.H{"error":"User not found"})
        return
    }
    context.JSON(http.StatusOK, user_profile)

}
func Leaderboard(context *gin.Context) {
	var users []models.User

	// Query the database to get the first five users with the maximum points
	err := utils.DB.Order("points desc").Limit(5).Find(&users).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, users)
}
// func GetUsersList(context *gin.Context) {
//     user_id, err := strconv.ParseUint(context.Param("id"), 10, 64)
//     if err!=nil{
//         context.JSON(http.StatusBadRequest, gin.H{"error":"Incorrect user ID"})
//         return
//     }
//     var user_profile models.User
//     utils.DB.First(&user_profile, user_id)
//     var users []models.User
//     utils.DB.Where("team <> ? and team <> 'coordinator'", user_profile.Team).Find(&users)
//     // if result.Error!=nil {
//     //     context.JSON(http.StatusBadRequest, gin.H{"error":"Users not found"})
//     //     return
//     // }
//     context.JSON(http.StatusOK, users)

// }
// func GetQuestion(context *gin.Context){
//     chosen_user_id, err := strconv.ParseUint(context.Param("id"), 10, 64)
//     if err!=nil{
//         context.JSON(http.StatusBadRequest, gin.H{"error":"Incorrect user ID"})
//         return
//     }
//     var question models.Question
//     result:=utils.DB.Where("user_id = ?",chosen_user_id).First(&question)
//     if result.RowsAffected == 0{
//         context.JSON(http.StatusBadRequest, gin.H{"error":"Question not added"})
//         return
//     }
//     context.JSON(http.StatusOK, question)
    
// }

// func GetAnswers(context *gin.Context){
//     user_id, err := strconv.ParseUint(context.Param("uid"), 10, 64)
//     question_id, err := strconv.ParseUint(context.Param("qid"), 10, 64)
//     if err!=nil{
//         context.JSON(http.StatusBadRequest, gin.H{"error":err.Error})
//         return
//     }
//     var answers []models.Answer
//     utils.DB.Where("user_id = ? AND q_id = ?", user_id, question_id).Find(&answers)
//     context.JSON(http.StatusOK, answers)
// }

//Admin functions
// func GetAllParticipants(context *gin.Context){
//     var participants []models.User
//     utils.DB.Where("team <> 'coordinator'").Find(&participants)
//     context.JSON(http.StatusOK, participants)
// }

// type RequestBody struct {
// 	Points int `json:"points"`
// }

// func PointsDecrement(context *gin.Context){
//     var requestBody RequestBody
// 	if err := context.ShouldBindJSON(&requestBody); err != nil {
// 		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//     points := requestBody.Points
//     user_id, err := strconv.ParseUint(context.Param("uid"), 10, 64)
//     if err!=nil{
//         context.JSON(http.StatusBadRequest, gin.H{"error":err.Error})
//         return
//     }
//     var user_profile models.User
//     result := utils.DB.First(&user_profile, user_id)
//     if result.Error!=nil {
//         context.JSON(http.StatusBadRequest, gin.H{"error":"User not found"})
//         return
//     }
//     user_profile.Points = user_profile.Points - points
//     utils.DB.Save(&user_profile)

//     context.JSON(http.StatusOK, user_profile)
// }