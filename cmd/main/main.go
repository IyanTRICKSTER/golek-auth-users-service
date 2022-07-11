package main

import (
	"acourse-auth-user-service/pkg/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	//Load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	//Enable Gin Debugging Mode
	//gin.SetMode(gin.ReleaseMode)

	engine := gin.Default()

	//Registering Routes
	routes.RegisterRoutes(engine)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "80"
	}

	engine.Run(":" + port)

}
