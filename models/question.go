package models

import(
	"github.com/NamrithaGirish/hack4tkm/utils"
    // "net/url"
)

type User struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement" form:"id"`
	Name string `json:"name" gorm:"not null" form:"name"`
    Image string `gorm:"column:image" json:"image" form:"image"`
    Mail string `json:"gmail" gorm:"not null" form:"gmail"`
    Team string `json:"team_name" gorm:"not null" form:"team_name"`
    Points int `json:"points" gorm:"default:0" form:"points"`
}

func (user *User) Save() (*User, error) {
    err := utils.DB.Create(&user).Error
    if err != nil {
        return &User{}, err
    }
    return user, nil
}


type Comments struct{
	ID uint `json:"id" gorm:"primaryKey;autoIncrement" form:"id"`
	Comment string `json:"comment" gorm:"not null" form:"comment"`
    LinkedinUrl string `json:"linkedin_url" form:"linkedin_url"`
    Image string `gorm:"column:image" json:"image" form:"image"`
	SenderID uint `json:"sender_id" gorm:"not null" form:"sender_id"`
	ReceiverID uint `json:"receiver_id" gorm:"not null" form:"receiver_id"`
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

// type Credentials struct{
//     ID uint `json:"id" gorm:"primaryKey;autoIncrement" form:"id"`
//     UserID uint `json:"user_id" gorm:"not null" form:"user_id"`
//     Username string `json:"username" gorm:"not null" form:"username"`
//     Password string `json:"password" gorm:"not null" form:"password"`
//     User User `gorm:"foreignKey:UserID"`

// }
// func (credential *Credentials) Save() (*Credentials, error) {
//     err := utils.DB.Create(&credential).Error
//     if err != nil {
//         return &Credentials{}, err
//     }
//     return credential, nil
// }
