package main

import (
    "encoding/csv"
    "fmt"
    "os"
	"github.com/NamrithaGirish/hack4tkm/models"
	"github.com/NamrithaGirish/hack4tkm/utils"
	"strings"
)

func init() {
	utils.ConnectDB()
}

func main() {
    filePath := "./csv/list.csv" // Specify the CSV filename here

    file, err := os.Open(filePath)
    if err != nil {
		fmt.Println("cant open file")
        return
    }
    defer file.Close()

    reader := csv.NewReader(file)
	fmt.Println("Read from csv")

    // Read CSV records
    records, err := reader.ReadAll()
    if err != nil {
		fmt.Println("cannot read from file")
        return
    }
	fmt.Println("read from file...records initialized")

    // Iterate through each record
    for _, record := range records {
        name := record[1]
        mail := record[2]
        team := record[0]
		fmt.Println(name,record[1],record[2])

        user := models.User{
			Name: name,
			Mail: mail,
			Team: strings.TrimSpace(team),
		}
		fmt.Printf("%T",name)
        // Save User to the database
        savedUser, err := user.Save()
        if err != nil {
            fmt.Printf("Error saving user %s: %s\n", savedUser.Name, err)
            // If an error occurs, you may choose to handle it according to your application's logic
            // For example, you can log the error and continue processing the next record
        } else {
            fmt.Printf("User %s saved successfully.\n", savedUser.Name)
        }
    }
}

