package main

import (
	"fmt"

	"github.com/NamrithaGirish/hack4tkm/utils"
	"github.com/NamrithaGirish/hack4tkm/models"
)

func init() {
	utils.ConnectDB()
}

func main() {
	utils.DB.AutoMigrate(&models.User{},&models.Question{},&models.Answer{})
	fmt.Println("? Migration complete")
}

