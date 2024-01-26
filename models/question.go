package models

import(
	"github.com/NamrithaGirish/hack4tkm/utils"
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

type Question struct{
	ID uint `json:"id" gorm:"primaryKey"`
	Question string `json:"question" gorm:"not null"`
	Answer string `json:"answer" gorm:"not null"`
	UserID uint `json:"user_id" gorm:"not null"`
	User User `gorm:"foreignKey:UserID"`
}
func (question *Question) Save() (*Question, error) {
    err := utils.DB.Create(&question).Error
    if err != nil {
        return &Question{}, err
    }
    return question, nil
}
type Answer struct{
	ID uint `json:"id" gorm:"primaryKey"`
	Answer string `json:"answer" gorm:"not null"`
    
	QID uint `json:"question_id" gorm:"not null"`
	Question Question `gorm:"foreignKey:QID"`
	UserID uint `json:"user_id" gorm:"not null"`
	User User `gorm:"foreignKey:UserID"`
}
func (answer *Answer) Save() (*Answer, error) {
    err := utils.DB.Create(&answer).Error
    if err != nil {
        return &Answer{}, err
    }
    return answer, nil
}
