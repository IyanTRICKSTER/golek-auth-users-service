package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"strings"
	"time"
)

func GenerateRefershToken(user_id uint) (string, error) {

	token_lifespan, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DAY_LIFESPAN"))

	if err != nil {
		return "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["authorized"] = true
	rtClaims["user_id"] = user_id
	rtClaims["exp"] = time.Now().Add(time.Hour * (time.Duration(token_lifespan) * 24)).Unix()

	return refreshToken.SignedString([]byte(os.Getenv("JWT_REFRESH_TOKEN_SECRET")))
}

func ValidateRefreshToken(c *gin.Context) (*jwt.Token, error) {

	tokenString := ExtractRefreshToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_REFRESH_TOKEN_SECRET")), nil
	})
	
	if err != nil {
		return nil, err
	}

	return token, nil
}

func ExtractRefreshToken(c *gin.Context) string {

	refershToken := c.Query("refresh-token")
	if refershToken != "" {
		return refershToken
	}

	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}
