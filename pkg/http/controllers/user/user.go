package UserController

import (
	"acourse-auth-user-service/pkg/http/requests"
	model "acourse-auth-user-service/pkg/models"
	tokenUtils "acourse-auth-user-service/pkg/utils/jwt"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"strconv"
)

func All(c *gin.Context) {

	page, _ := strconv.ParseUint(c.Query("page"), 10, 32)
	limit, _ := strconv.ParseUint(c.Query("limit"), 10, 32)

	list, err := model.AllUser(uint(limit), uint(page))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, list)
	return
}

func CurrentUser(c *gin.Context) {

	log.Println(c.Request.Header.Get("X-User-Id"))
	log.Println(c.Request.Header.Get("X-User-Role"))

	userId, err := tokenUtils.ExtractAccessTokenID(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := model.FindUser(userId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func Update(c *gin.Context) {

	userId, err := strconv.Atoi(c.Param("id"))

	log.Println("USER ID >> ", userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var input requests.UpdateUserRecordCredential

	if err := c.ShouldBindJSON(&input); err != nil {
		if errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input body"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user, err := model.FindUser(uint(userId)); err != nil {
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		err := user.UpdateUser(input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "User Updated Successfully",
		})
		return
	}
}

func Delete(c *gin.Context) {

	userId, err := strconv.Atoi(c.Param("id"))

	log.Println("USER ID >> ", userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user, err := model.FindUser(uint(userId)); err != nil {
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		err := user.DeleteUser()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "User Deleted Successfully",
		})
		return
	}

}
