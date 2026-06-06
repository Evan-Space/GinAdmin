package controller

import (
	"GinAdmin/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)



type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

func (ctl *UserController) List(c *gin.Context) {
	users, err := ctl.userService.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "query users failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}