package routes

import (
	AuthController "acourse-auth-user-service/pkg/http/controllers/auth"
	UserController "acourse-auth-user-service/pkg/http/controllers/user"
	AuthMiddleware "acourse-auth-user-service/pkg/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(route *gin.Engine) {

	publicRoutes := route.Group("/api/auth")
	publicRoutes.POST("/register", AuthController.Register)
	publicRoutes.POST("/login", AuthController.Login)
	publicRoutes.POST("/change-password", AuthController.ChangePassword)
	publicRoutes.POST("/reset-password", AuthController.ResetPassword)

	protectedRoutes := route.Group("/api")
	protectedRoutes.GET("/auth/introspect", AuthMiddleware.IsUserAuthenticatedMiddleware(), AuthController.InstrospectToken)
	protectedRoutes.GET("/auth/token/refresh", AuthMiddleware.IsUserAllowedToRefreshTokenMiddleware(), AuthController.RefreshToken)

	userRoute := protectedRoutes.Group("/user")
	userRoute.GET("/all", AuthMiddleware.CanListUserPermission(), UserController.All)
	userRoute.GET("/current", AuthMiddleware.CanReadUserPermission(), UserController.CurrentUser)
	userRoute.PATCH("/update/:id", AuthMiddleware.CanUpdateUserPermission(), UserController.Update)
	userRoute.DELETE("/delete/:id", AuthMiddleware.CanDeleteUserPermission(), UserController.Delete)

}
