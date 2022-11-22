package RoleController

import (
	"github.com/gin-gonic/gin"
	model "golek-auth-user-service/pkg/models"
	"net/http"
)

func Find(c *gin.Context) {
	role, err := model.FindRole(1)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": role,
	})
}
