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

func GenerateResetToken(user_id uint) (string, error) {

	token_lifespan, err := strconv.Atoi(os.Getenv("RESET_TOKEN_MINUTE_LIFESPAN"))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	//claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(token_lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_RESET_TOKEN_SECRET")))
}

func ExtractResetToken(c *gin.Context) string {

	refershToken := c.Query("reset-token")
	if refershToken != "" {
		return refershToken
	}

	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}

func ValidateResetToken(c *gin.Context) (*jwt.Token, string, error) {

	tokenString := ExtractResetToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_RESET_TOKEN_SECRET")), nil
	})

	if err != nil {
		return nil, "", err
	}

	return token, tokenString, nil
}
