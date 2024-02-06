package models

import(
	"github.com/NamrithaGirish/hack4tkm/utils"
    // "net/url"
)

type User struct {
	ID uint `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`
    Mail string `json:"gmail" gorm:"not null"`
    Team string `json:"team_name" gorm:"not null"`
    Points int `json:"points" gorm:"default:0"`
}

func (user *User) Save() (*User, error) {
    err := utils.DB.Create(&user).Error
    if err != nil {
        return &User{}, err
    }
    return user, nil
}

// type Question struct{
// 	ID uint `json:"id" gorm:"primaryKey"`
// 	Question string `json:"question" gorm:"not null"`
// 	Answer string `json:"answer" gorm:"not null"`
// 	UserID uint `json:"user_id" gorm:"not null"`
// 	User User `gorm:"foreignKey:UserID"`
// }
// func (question *Question) Save() (*Question, error) {
//     err := utils.DB.Create(&question).Error
//     if err != nil {
//         return &Question{}, err
//     }
//     return question, nil
// }
// type Answer struct{
// 	ID uint `json:"id" gorm:"primaryKey"`
// 	Answer string `json:"answer" gorm:"not null"`
    
// 	QID uint `json:"question_id" gorm:"not null"`
// 	Question Question `gorm:"foreignKey:QID"`
// 	UserID uint `json:"user_id" gorm:"not null"`
// 	User User `gorm:"foreignKey:UserID"`
// }
// func (answer *Answer) Save() (*Answer, error) {
//     err := utils.DB.Create(&answer).Error
//     if err != nil {
//         return &Answer{}, err
//     }
//     return answer, nil
// }

type Comments struct{
	ID uint `json:"id" gorm:"primaryKey"`
	Comment string `json:"comment" gorm:"not null"`
    LinkedinUrl string `json:"linkedin_url"`
    Image string `gorm:"column:image" json:"image"`
	SenderID uint `json:"sender_id" gorm:"not null"`
	ReceiverID uint `json:"receiver_id" gorm:"not null"`
	Sender User `gorm:"foreignKey:SenderID"`
	Receiver User `gorm:"foreignKey:ReceiverID"`
}
func (comment *Comments) Save() (*Comments, error) {
    err := utils.DB.Create(&comment).Error
    if err != nil {
        return &Comments{}, err
    }
    return comment, nil
}
// func isImageURL(fl validator.FieldLevel) bool {
// 	value := fl.Field().String()

// 	// Perform validation logic
// 	if strings.HasSuffix(strings.ToLower(value), ".jpg") || strings.HasSuffix(strings.ToLower(value), ".jpeg") || strings.HasSuffix(strings.ToLower(value), ".png") || strings.HasSuffix(strings.ToLower(value), ".gif") {
// 		return true
// 	}

// 	return false
// }

// func isValidURL(fl validator.FieldLevel) bool {
// 	value := fl.Field().String()

// 	// Regular expression for a simple URL validation
// 	urlPattern := regexp.MustCompile(`^(http|https):\/\/[^\s/$.?#].[^\s]*$`)

// 	return urlPattern.MatchString(value)
// }
