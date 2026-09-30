package controller

import (
	"GinAdmin/internal/pkg/errors"
	"GinAdmin/internal/service"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	Api
	uploadService *service.UploadService
}

func NewUploadController() *UploadController {
	return &UploadController{
		uploadService: service.NewUploadService(),
	}
}

func (ctl *UploadController) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		ctl.Fail(c, errors.InvalidParameter, "请选择文件")
		return
	}
	path, err := ctl.uploadService.SaveFile(file) // 调用 service 中方法，写入文件并且拿到返回值 path
	if err != nil {
		ctl.Fail(c, errors.ServerErr, err.Error())
		return
	}

	ctl.Success(c, gin.H{"path": path})
}
