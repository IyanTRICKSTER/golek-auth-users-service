package routes

import (
	AuthController "acourse-auth-user-service/pkg/http/controllers/auth"
	RoleController "acourse-auth-user-service/pkg/http/controllers/role"
	UserController "acourse-auth-user-service/pkg/http/controllers/user"
	authMiddleware "acourse-auth-user-service/pkg/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(route *gin.Engine) {

	public := route.Group("/api/auth")

	public.POST("/register", AuthController.Register)
	public.POST("/login", AuthController.Login)
	public.POST("/change-password", AuthController.ChangePassword)
	public.POST("/reset-password", AuthController.ResetPassword)

	protected := route.Group("/api/auth")
	protected.Use(authMiddleware.IsUserAuthenticatedMiddleware())
	protected.GET("/user/all", UserController.All)
	protected.GET("/user/current", AuthController.CurrentUser)
	protected.PATCH("/user/update/:id", UserController.Update)
	protected.DELETE("/user/delete/:id", UserController.Delete)
	protected.GET("/role", RoleController.Find)

	refreshToken := route.Group("/api/auth")
	refreshToken.Use(authMiddleware.IsUserAllowedToRefreshTokenMiddleware())
	refreshToken.GET("/token/refresh", AuthController.RefreshToken)

}
