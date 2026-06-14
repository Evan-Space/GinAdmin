package controller

import (
	"GinAdmin/global"
	"GinAdmin/internal/errors"
	"GinAdmin/internal/pkg/response"

	"github.com/gin-gonic/gin"
)


type Api struct {}


// Success 业务成功响应
func (api Api) Success(c *gin.Context, data ...any) {
	if len(data) > 0 && data[0] != nil {
		response.Success(c, data[0])
		return
	}
	response.Success(c)
}

// FailCode 业务失败响应（按错误码返回）
func (api Api) FailCode(c *gin.Context, code int, data ...any) {
	response.FailCode(c, code, data...)
}

// Fail 业务失败响应（自定义消息）
func (api Api) Fail(c *gin.Context, code int, message string, data ...any) {
	response.Fail(c, code, message, data...)
}

// Err 统一错误处理
func (api Api) Err(c *gin.Context, err error) {
	helper := errors.ErrorHelper{}
	businessError, parseErr := helper.AsBusinessError(err)
	if parseErr != nil {
		response.FailCode(c, errors.ServerErr)
		return
	}
	response.Fail(c, businessError.GetCode(), businessError.GetMessage())
}

// GetCurrentUserID 获取当前登录用户 ID
func (api Api) GetCurrentUserID(c *gin.Context) uint {
	return c.GetUint(global.ContextKeyUID)
}