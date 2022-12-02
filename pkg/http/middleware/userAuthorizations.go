package AuthMiddleware

import (
	"github.com/gin-gonic/gin"
	tokenUtils "golek-auth-user-service/pkg/utils/jwt"
	"net/http"
	"strconv"
	"strings"
)

func CanListUserPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		//log.Println("Hit List User Middleware")
		if !strings.Contains(c.Request.Header.Get("X-User-Permission"), "l") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Forbidden",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CanReadUserPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		//log.Println("Hit Read User Middleware")
		//log.Println(c.Request.Header.Get("X-User-Permission"))
		if !strings.Contains(c.Request.Header.Get("X-User-Permission"), "r") {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Forbidden",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CanUpdateUserPermission() gin.HandlerFunc {
	return func(c *gin.Context) {

		//log.Println("Update User Middleware")
		//log.Println(c.Request.Header.Get("X-User-Permission"))
		if !strings.Contains(c.Request.Header.Get("X-User-Permission"), "u") {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Forbidden",
			})
			c.Abort()
			return
		}

		//Get User Id in access token
		userId, err := tokenUtils.ExtractAccessTokenID(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		//Extract User Id from URL Param
		uid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		//A User cannot modify data owned by another user
		if userId != uint(uid) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Forbidden",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func CanDeleteUserPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		//log.Println("Delete User Middleware")
		//log.Println(c.Request.Header.Get("X-User-Permission"))
		if !strings.Contains(c.Request.Header.Get("X-User-Permission"), "d") {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Forbidden",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
