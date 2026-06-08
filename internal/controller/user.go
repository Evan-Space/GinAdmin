package controller

import (
	"GinAdmin/internal/service"
	"net/http"

	"GinAdmin/internal/pkg/response"

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
		response.Resp().Fail(c, http.StatusInternalServerError, "query users failed")
		return
	}

	response.Ok(c, users)
}