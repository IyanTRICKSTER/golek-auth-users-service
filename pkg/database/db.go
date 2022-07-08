package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var conn *gorm.DB

func Connect() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	DBDriver := os.Getenv("DB_DRIVER")
	DbHost := os.Getenv("DB_HOST")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	//Open connection to the database
	if conn == nil {

		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

		conn, err = gorm.Open(mysql.Open(DBURL), &gorm.Config{})

		if err != nil {
			log.Fatal("Connection error: ", err)
		} else {
			fmt.Println("Connected to the database", DBDriver)
		}

	} else {
		fmt.Println("Database already connected")
	}

}

func GetConnection() *gorm.DB {
	return conn
}
