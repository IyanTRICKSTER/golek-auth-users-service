package AuthMiddleware

import (
	tokenUtils "acourse-auth-user-service/pkg/utils/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IsUserAuthenticatedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tokenUtils.ValidateAccessToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func IsUserAllowedToRefreshTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := tokenUtils.ValidateRefreshToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
