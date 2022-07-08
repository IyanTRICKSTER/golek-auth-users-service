package UserController

import (
	"acourse-auth-user-service/pkg/http/requests"
	model "acourse-auth-user-service/pkg/models"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

func All(c *gin.Context) {

	page, _ := strconv.ParseUint(c.Query("page"), 10, 32)
	limit, _ := strconv.ParseUint(c.Query("limit"), 10, 32)

	list, err := model.AllUser(uint(limit), uint(page))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, list)
	return
}

func Update(c *gin.Context) {

	userId, err := strconv.Atoi(c.Param("id"))

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
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
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

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user, err := model.FindUser(uint(userId)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
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
