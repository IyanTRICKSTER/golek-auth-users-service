package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var conn *gorm.DB

func Connect() {

	if conn == nil {
		err := godotenv.Load(".env")

		if err != nil {
			log.Fatalf("Error loading .env file")
		}

		Dbdriver := os.Getenv("DB_DRIVER")
		DbHost := os.Getenv("DB_HOST")
		DbUser := os.Getenv("DB_USER")
		DbPassword := os.Getenv("DB_PASSWORD")
		DbName := os.Getenv("DB_NAME")
		DbPort := os.Getenv("DB_PORT")

		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

		conn, err = gorm.Open(Dbdriver, DBURL)

		if err != nil {
			fmt.Println("Cannot connect to database ", Dbdriver)
			log.Fatal("connection error:", err)
		} else {
			fmt.Println("Connected to the database ", Dbdriver)
		}
	} else {
		fmt.Println("Database already connected")
	}

}

func GetConnection() *gorm.DB {
	return conn
}
