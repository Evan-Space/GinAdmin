package controller

import (
	"GinAdmin/internal/pkg/response"
	"GinAdmin/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wannanbigpig/gin-layout/internal/dto"
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
	var req dto.PageReq
	if !request.BingQuery(c, &req) {
		return
	}
	req.Normalize()
	
	users, err := ctl.userService.ListUsers()
	if err != nil {
		response.Resp().Fail(c, http.StatusInternalServerError, "query users failed")
		return
	}

	response.Ok(c, users)
}

func (c *UserController) Info(ctx *gin.Context) {                                                                                                                                                                  
	ctx.JSON(200, gin.H{"name": "test"})
}                                                                                                                                                                                                                  
				