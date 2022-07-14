package AuthMiddleware

import (
	tokenUtils "acourse-auth-user-service/pkg/utils/jwt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func CanListUserPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("List User Middleware")
		log.Println(c.Request.Header.Get("X-User-Permission"))
	}
}

func CanReadUserPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Read User Middleware")
		log.Println(c.Request.Header.Get("X-User-Permission"))
	}
}

func CanUpdateUserPermission() gin.HandlerFunc {
	return func(c *gin.Context) {

		log.Println("Update User Middleware")
		log.Println(c.Request.Header.Get("X-User-Permission"))

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
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		//A User cannot modify data owned by another user
		if userId != uint(uid) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Forbidden",
			})
			c.Abort()
			return
		}
	}
}

func CanDeleteUserPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Delete User Middleware")
		log.Println(c.Request.Header.Get("X-User-Permission"))
	}
}
