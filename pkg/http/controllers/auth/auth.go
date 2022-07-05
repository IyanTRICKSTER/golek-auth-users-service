package auth

import (
	"acourse-auth-user-service/pkg/http/requests"
	"acourse-auth-user-service/pkg/models"
	tokenUtils "acourse-auth-user-service/pkg/utils/jwt"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func Login(c *gin.Context) {

	var input requests.LoginCredential

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := UserModel.User{}

	user.Username = input.Username
	user.Password = input.Password

	pairToken, err := UserModel.Authenticate(user.Username, user.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pairToken)

}

func Register(c *gin.Context) {

	var input requests.RegisterCredential

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := UserModel.User{}

	user.Username = input.Username
	user.Password = input.Password
	user.Email = input.Email
	user.PhoneNumber = input.PhoneNumber

	_, err := user.Create()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration success"})
}

func CurrentUser(c *gin.Context) {

	userId, err := tokenUtils.ExtractAccessTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := UserModel.Find(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": user})
}

func RefreshToken(c *gin.Context) {

	refreshToken, err := tokenUtils.ValidateRefreshToken(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		if claims["user_id"] != "" {

			user_id := uint(claims["user_id"].(float64))

			_, err := UserModel.Find(user_id)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			newAccessToken, err := tokenUtils.GenerateAccessToken(user_id)
			newRefreshToken, err := tokenUtils.GenerateRefershToken(user_id)

			c.JSON(http.StatusOK, gin.H{
				"token":         newAccessToken,
				"refresh_token": newRefreshToken,
			})
		}

	}

}

func ChangePassword(c *gin.Context) {

	var input requests.ChangePasswordCredential

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/*
		Check if user exist by email, then
		create the reset token for intended user
	*/
	user, err := UserModel.FindByEmail(input.Email)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{"message": "Email is not exists"})
		return
	}

	user.IssueResetToken()

	c.JSON(http.StatusOK, gin.H{"message": "Reset Token has been set"})

}

func ResetPassword(c *gin.Context) {

	var input requests.ResetPasswordCredential

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
			//Find A User and Update the password
			user, _ := UserModel.Find(userId)

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

			user.RemoveResetToken()

			err := user.UpdatePassword(input.Password)

			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"message": err.Error(),
				})
				return
			}

		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password has ben changed",
	})
	return
}
