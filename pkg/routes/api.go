package routes

import (
	AuthController "acourse-auth-user-service/pkg/http/controllers/auth"
	authMiddleware "acourse-auth-user-service/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(route *gin.Engine) {

	public := route.Group("/api/auth")

	public.POST("/register", AuthController.Register)
	public.POST("/login", AuthController.Login)
	public.POST("/change-password", AuthController.ChangePassword)
	public.POST("/reset-password", AuthController.ResetPassword)

	protected := route.Group("/api/auth")
	protected.Use(authMiddleware.JwtAuthMiddleware())
	protected.GET("/current-user", AuthController.CurrentUser)

	refreshToken := route.Group("/api/auth")
	refreshToken.Use(authMiddleware.JwtAuthRefreshTokenMiddleware())
	refreshToken.GET("/token/refresh", AuthController.RefreshToken)

}
