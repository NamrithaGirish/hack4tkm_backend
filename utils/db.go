// utils/database.go
package utils

// var db *sql.DB

// // InitializeDB initializes the PostgreSQL database connection
// func InitializeDB() *sql.DB {
// 	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
// 		os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DB_NAME"))

// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		fmt.Println("There is an error while connecting to the database ", err)
//     	log.Fatal(err)
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		fmt.Println("There is an error while connecting to the database2 ", err)
//     	log.Fatal(err)
// 	}

// 	return db
// }

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"github.com/joho/godotenv"
	"fmt"
   )

// func ConnectDatabase()  {
//    dsn := "host=localhost user=postgres password=me@not&post dbname=hack4tkm port=5432 sslmode=disable TimeZone=Asia/Mumbai"
//    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//    return db
//    }

var DB *gorm.DB

func ConnectDB() {
	godotenv.Load(".env")
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PORT"))

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to the Database")
	}
	fmt.Println("? Connected Successfully to the Database")
}
