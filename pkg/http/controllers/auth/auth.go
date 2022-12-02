package auth

import (
	"errors"
	"golek-auth-user-service/pkg/http/requests"
	model "golek-auth-user-service/pkg/models"
	"golek-auth-user-service/pkg/notification"
	tokenUtils "golek-auth-user-service/pkg/utils/jwt"
	"log"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func Login(c *gin.Context) {

	var input requests.LoginCredentialRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pairToken, err := model.AuthenticateUser(input.Email, input.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authentication Failed"})
		return
	}

	c.JSON(http.StatusOK, pairToken)

}

func Register(c *gin.Context) {

	var input requests.RegisterCredentialRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{}

	user.Username = input.Username
	user.Password = input.Password
	user.Email = input.Email
	user.NIM = &input.NIM
	user.NIP = &input.NIP
	user.Major = input.Major

	_, err := user.CreateUser()

	if err != nil {

		mysqlErr := err.(*mysql.MySQLError)

		switch mysqlErr.Number {
		case 1062:
			c.JSON(http.StatusBadRequest, gin.H{"error": mysqlErr.Message})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration success"})
}

func IntrospectToken(c *gin.Context) {

	color.Set(color.FgYellow)

	//Extract Target Resource From API GATEWAY
	targetResource := c.Request.Header.Get("X-Target-Resource")
	log.Println("REQUEST TARGET RESOURCE:", targetResource)

	//Extract User Id From JWT Token
	userId, err := tokenUtils.ExtractAccessTokenID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//INTROSPECT USER and PERMISSIONS, determine what this User can do with Target Resource
	user, err := model.IntrospectUser(userId, targetResource)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Prepare User Permission
	var permissionCode string
	for _, v := range user.Role.Permissions {
		permissionCode += v.Code
	}

	log.Println("INTROSPECT USER ID:", user.ID)
	log.Println("INTROSPECT USERNAME:", user.Username)
	//log.Println("INTROSPECT USER MAJOR:", user.Major)
	log.Println("INTROSPECT USER ROLE:", user.Role.Name)
	log.Println("INTROSPECT USER PERMISSIONS:", permissionCode)

	//Set Custom Response Header
	c.Header("X-User-Id", strconv.Itoa(int(user.ID)))
	c.Header("X-User-Role", user.Role.Name)
	c.Header("X-User-Permission", permissionCode)
	c.Header("X-User-Name", user.Username)
	c.Header("X-User-Major", user.Major)

	c.JSON(http.StatusOK, gin.H{"message": "authorized", "user": gin.H{
		"user_id":     user.ID,
		"username":    user.Username,
		"role":        user.Role.Name,
		"permissions": permissionCode,
	}})

	color.Unset()
}

func RefreshToken(c *gin.Context) {

	//Validate Refresh Token
	refreshToken, err := tokenUtils.ValidateRefreshToken(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		if claims["user_id"] != "" {

			user_id := uint(claims["user_id"].(float64))

			_, err := model.FindUser(user_id)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			newAccessToken, err := tokenUtils.GenerateAccessToken(user_id)
			// newRefreshToken, err := tokenUtils.GenerateRefreshToken(user_id)

			c.JSON(http.StatusOK, gin.H{
				"access_token":  newAccessToken,
				"refresh_token": "",
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Cannot Indetifiy User Id In The Given Access Token",
			})
			return
		}

	}

}

func ChangePassword(c *gin.Context) {

	var input requests.ChangePasswordCredentialRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/*
		Check if user exist by email, then
		create the reset token for intended user
	*/
	user, err := model.FindUserByEmail(input.Email)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{"message": "Email is not exists"})
		return
	}

	user.IssueResetTokenUser()

	//Send Reset Token with smpt mail
	var mailError = make(chan error)
	go notification.Sendmail(mailError, user.ResetToken)

	c.JSON(http.StatusOK, gin.H{"message": "Reset Token has been sent"})
	return

}

func ResetPassword(c *gin.Context) {

	var input requests.ResetPasswordCredentialRequest

	//1. Validate Reset Password Credential Input Binding
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//2. Validate if provided password are the same
	if err := input.ValidateResetPasswordCredential(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//3. Validate If Reset Token is valid
	resetToken, tokenString, err := tokenUtils.ValidateResetToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	//4. Finally, We Can Reset The Password
	if claims, ok := resetToken.Claims.(jwt.MapClaims); ok && resetToken.Valid {

		userId := uint(claims["user_id"].(float64))

		if claims["user_id"] != "" {
			//FindUser A User and Update the password
			user, _ := model.FindUser(userId)

			/*
				But we also need to check that the reset token also exist in the user record,
				If it exists, we compare the given reset token with the token in the user record,
				if it's the same reset token, we remove it and update the password
			*/
			if user.ResetToken != tokenString {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "Reset Token is not valid",
				})
				return
			}

			user.RemoveResetTokenUser()

			err := user.UpdateUserPassword(input.Password)

			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"message": err.Error(),
				})
				return
			}

		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password has been updated",
	})
	return
}
