package controllers

import (
	"github.com/NamrithaGirish/hack4tkm/models"
	"github.com/gin-gonic/gin"
    "net/http"
    "github.com/NamrithaGirish/hack4tkm/utils"
    "strconv"
    "github.com/joho/godotenv"
    "context"
    // "io/ioutil"
    "fmt"
    "os"
    "strings"
	
    
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

func AddUser(context *gin.Context) {
    var input models.User

    if err := context.ShouldBindJSON(&input); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user := models.User{
        Name: input.Name,
        // ID: input.ID,
        Mail: input.Mail,
        Team: input.Team,
        Image: "https://iedcbackend.s3.us-west-1.amazonaws.com/hack4tkm/head_contact.png",
    }

    savedUser, err := user.Save()
    // // fmt.Printf("%T\n",savedUser)
    // // fmt.Println(savedUser)
    // custom_password:=strings.ReplaceAll(strings.ToLower(savedUser.Name)," ","_")
    // // fmt.Println(custom_password)
    // // fmt.Println(savedUser.ID)
    // credentials := models.Credentials{
    //     Username:input.Mail,
    //     Password:custom_password,
    //     UserID:savedUser.ID,
    // }
    // savedCredential, err := credentials.Save()

    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    context.JSON(http.StatusCreated, gin.H{"user": savedUser})
}

func GetAllTeams(context *gin.Context) {
    var teams []string
    result:=utils.DB.Model(&models.User{}).Select("team").Where("team <> 'coordinator'").Group("team").Find(&teams)
    if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
	}
	context.JSON(http.StatusOK, gin.H{"teams": teams})
}

func GetTeamMembers(context *gin.Context) {
    var names []models.User
    team:=context.Param("team")
    result:=utils.DB.Where("team = ?",team).Find(&names)
    if result.Error != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
	}
	context.JSON(http.StatusOK, gin.H{"members": names})
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

func AddComment(c *gin.Context) {

    var input models.Comments
    input.Comment = c.PostForm("comment")
    input.LinkedinUrl=c.PostForm("linkedin_url")
    // input.Image=c.PostForm("image"),
    s_id, err:=strconv.ParseUint(c.PostForm("sender_id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not start connection"})
        return
    }
    input.SenderID=uint(s_id)
    r_id, err:=strconv.ParseUint(c.PostForm("receiver_id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not start connection"})
        return
    }
    input.ReceiverID=uint(r_id)
    // id, err:=strconv.ParseUint(c.PostForm("id"), 10, 64)
    // if err != nil {
    //     c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not start connection"})
    //     return
    // }
    // input.ID=uint(id)
    // fmt.Println(c.PostForm()+"KLL")
    // fmt.Printf("%T",c.PostForm())

    // if err := c.Bind(&input); err != nil {
    //     fmt.Println(input)
    //     c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    //     return
    // }
    
    if (CommentEnable(input.ReceiverID, input.SenderID)){
        //Uploading file to s3


        godotenv.Load(".aws")
        cfg, err := config.LoadDefaultConfig(context.TODO())
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not start connection"})
            return
        }

        client := s3.NewFromConfig(cfg)
        uploader := manager.NewUploader(client)
        fmt.Println("client created")
        file, err:=c.FormFile("image")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open file"})
            return
        }
        mimeType := file.Header.Get("Content-Type")
        uploadFile, err:= file.Open()
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }        

        result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
            Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
            Key:    aws.String("hack4tkm/comments/"+strconv.FormatUint(uint64(input.SenderID), 10)+"_"+strconv.FormatUint(uint64(input.ReceiverID), 10)+".jpg"),
            Body:   uploadFile,
            ContentType: aws.String(mimeType),
        })
            
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        input.Image=result.Location
        comment := models.Comments{
            Comment: input.Comment,
            LinkedinUrl: input.LinkedinUrl,
            Image:input.Image,
            SenderID: input.SenderID,
            ReceiverID: input.ReceiverID,
        }
    
        savedComment, err := comment.Save()
    
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        var user models.User
        if err := utils.DB.First(&user, input.SenderID).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
            return
        }
        user.Points = user.Points+10
        utils.DB.Save(&user)
        c.JSON(http.StatusCreated, gin.H{"comment": savedComment})
    } else{
        c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add more than one comment"})
    }
}

func FindUserID(name string) uint{
    var user models.User
    utils.DB.Where("name = ?",name).First(&user)
    // fmt.Println(user.ID)
    return user.ID
}

//name as parameter
// func DisplayComments(context *gin.Context) {
//     var comments []models.Comments
//     name:=context.Param("name")
//     result:=utils.DB.Where("receiver_id = ?",FindUserID(name)).Find(&comments)

//     if result.Error != nil {
// 		context.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
// 	}
// 	context.JSON(http.StatusOK, comments)
// }

func DisplayComments(context *gin.Context) {
    user_id, err := strconv.ParseUint(context.Param("id"), 10, 64)
    if err!=nil{
        context.JSON(http.StatusBadRequest, gin.H{"error":"Incorrect user ID"})
        return
    }
    var comments []models.Comments
    result:=utils.DB.Where("receiver_id = ?",user_id).Find(&comments)

    if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
	}
	context.JSON(http.StatusOK, comments)
}

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
    var sent_comments []models.Comments
    result = utils.DB.Where("sender_id = ?",user_id).Find(&sent_comments)
    if result.Error!=nil {
        context.JSON(http.StatusBadRequest, gin.H{"error":"Could not fetch sent comments"})
        return
    }
    var received_comments []models.Comments
    result = utils.DB.Where("receiver_id = ?",user_id).Find(&received_comments)
    if result.Error!=nil {
        context.JSON(http.StatusBadRequest, gin.H{"error":"Could not fetch received comments"})
        return
    }
    context.JSON(http.StatusOK, gin.H{"user":user_profile,"sent":sent_comments,"received":received_comments})

}

func GetUserByMail(context *gin.Context) {
    user_mail := context.Param("mail")

    var user_profile models.User
    result := utils.DB.Where("mail = ?",user_mail).First(&user_profile)
    if result.Error!=nil {
        context.JSON(http.StatusBadRequest, gin.H{"error":"User not found"})
        return
    }
    if user_profile.Team=="coordinator" {
        context.JSON(http.StatusOK, gin.H{"coordinator":true})
        return 
    }
    // var user_credentials models.Credentials
    // result = utils.DB.Where("user_id = ?",user_profile.ID).First(&user_credentials)
    // if result.Error!=nil {
    //     context.JSON(http.StatusBadRequest, gin.H{"error":"User not found"})
    //     return
    // }
    context.JSON(http.StatusOK, gin.H{"user":user_profile})

}

func Leaderboard(context *gin.Context) {
	var users []models.User

	// Query the database to get the first five users with the maximum points
	err := utils.DB.Where("team <> 'coordinator'").Order("points desc").Limit(5).Find(&users).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, users)
}

func DeleteComment(context *gin.Context){
    comment_id, err := strconv.ParseUint(context.Param("comment_id"), 10, 64)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
        return
    }
    var comment models.Comments
    result := utils.DB.First(&comment, comment_id)
    if result.Error != nil {
        context.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
        return
    }
    if err := utils.DB.Delete(&comment).Error; err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
        return
    }

    context.JSON(http.StatusOK, comment)
}

func DeleteUser(context *gin.Context){
    user_id, err := strconv.ParseUint(context.Param("user_id"), 10, 64)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    var user models.User
    result := utils.DB.First(&user, user_id)
    if result.Error != nil {
        context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    if err := utils.DB.Delete(&user).Error; err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
        return
    }

    context.JSON(http.StatusOK, user)
}


// func Login(context *gin.Context){
//     var input models.Credentials 

//     // Extract username and password from request body
//     if err := context.ShouldBindJSON(&input); err != nil {
//         context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }
//     var user_credentials models.Credentials
//     result:=utils.DB.Where("username = ? and password = ?",input.Username,input.Password).Find(&user_credentials)
//     if result.Error != nil {
//         context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
//         return
//     }
//     context.JSON(http.StatusOK, gin.H{"message": "Login successful","user_id":user_credentials.UserID})

// }

func UpdatePhoto(c *gin.Context) {
    user_id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    var user_profile models.User
    result := utils.DB.First(&user_profile, user_id)
    if result.Error!=nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":"User not found"})
        return
    }

    godotenv.Load(".aws")
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not start connection"})
        return
    }

    client := s3.NewFromConfig(cfg)
    uploader := manager.NewUploader(client)
    fmt.Println("client created")
    file, err:=c.FormFile("image")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file from input"})
        return
    }
    mimeType := file.Header.Get("Content-Type")
    uploadFile, err:= file.Open()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }        

    resultAWS, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
        Key:    aws.String("hack4tkm/profile/"+strings.ToLower(strings.ReplaceAll(user_profile.Name," ","_"))+".jpg"),
        Body:   uploadFile,
        ContentType: aws.String(mimeType),
    })
        
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    fmt.Println(user_profile.Image)
	if user_profile.Image ==""{
        fmt.Println("Adding new Image")
        user_profile.Image = resultAWS.Location
        utils.DB.Save(&user_profile)
    }
    c.JSON(http.StatusOK, gin.H{"message": "User photo updated successfully","user":user_profile})

	// if found {
	// 	c.JSON(http.StatusOK, gin.H{"message": "User photo updated successfully"})
	// } else {
	// 	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
	// }
}