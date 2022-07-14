package routes

import (
	AuthController "acourse-auth-user-service/pkg/http/controllers/auth"
	UserController "acourse-auth-user-service/pkg/http/controllers/user"
	authMiddleware "acourse-auth-user-service/pkg/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(route *gin.Engine) {

	publicRoutes := route.Group("/api/auth")
	publicRoutes.POST("/register", AuthController.Register)
	publicRoutes.POST("/login", AuthController.Login)
	publicRoutes.POST("/change-password", AuthController.ChangePassword)
	publicRoutes.POST("/reset-password", AuthController.ResetPassword)

	refreshTokenRoute := route.Group("/api/auth")
	refreshTokenRoute.Use(authMiddleware.IsUserAllowedToRefreshTokenMiddleware())
	refreshTokenRoute.GET("/token/refresh", AuthController.RefreshToken)

	protectedRoutes := route.Group("/api/")
	protectedRoutes.Use(authMiddleware.IsUserAuthenticatedMiddleware())
	protectedRoutes.GET("/auth/introspect", AuthController.InstrospectToken)

	//User Routes, each Methods implement different authorization
	userRoute := protectedRoutes.Group("/user")

	ListUserRoute := userRoute.Use(authMiddleware.CanListUserPermission())
	ListUserRoute.GET("/all", UserController.All)

	CurrentUserRoute := userRoute.Use(authMiddleware.CanReadUserPermission())
	CurrentUserRoute.GET("/current", UserController.CurrentUser)

	UpdateUserRoute := userRoute.Use(authMiddleware.CanUpdateUserPermission())
	UpdateUserRoute.PATCH("/update/:id", UserController.Update)

	DeleteUserRoute := userRoute.Use(authMiddleware.CanDeleteUserPermission())
	DeleteUserRoute.DELETE("/delete/:id", UserController.Delete)

}
