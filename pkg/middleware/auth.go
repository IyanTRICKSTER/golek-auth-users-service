package AuthMiddleware

import (
	tokenUtils "acourse-auth-user-service/pkg/utils/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tokenUtils.ValidateAccessToken(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

func JwtAuthRefreshTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := tokenUtils.ValidateRefreshToken(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
