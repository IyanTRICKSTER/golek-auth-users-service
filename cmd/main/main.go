package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	AuthMiddleware "golek-auth-user-service/pkg/http/middleware"
	"golek-auth-user-service/pkg/routes"
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
	engine.Use(cors.New(AuthMiddleware.CORSConfig()))

	//Registering Routes
	routes.RegisterRoutes(engine)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "80"
	}

	engine.Run(":" + port)

}
