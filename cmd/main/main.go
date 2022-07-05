package main

import (
	"acourse-auth-user-service/pkg/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	routes.RegisterRoutes(engine)
	engine.Run(":8080")

}
