package controller

import (
	"GinAdmin/internal/pkg/errors"
	"GinAdmin/internal/service"
	"io"
	"net/http"

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

const maxUploadSize = 3 << 30 // 最大3GB，方便后面用 2GB 文件验收内存

func (ctl *UploadController) UploadFile(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	reader, err := c.Request.MultipartReader()
	if err != nil {
		ctl.Fail(c, errors.InvalidParameter, "请选择文件")
		return
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			ctl.Fail(c, errors.ServerErr, "读取文件内容失败")
			return
		}

		if part.FormName() != "file" {
			part.Close()
			continue
		}

		path, err := ctl.uploadService.SaveFile(part)
		part.Close()
		if err != nil {
			ctl.Fail(c, errors.ServerErr, err.Error())
			return
		}
		ctl.Success(c, gin.H{"path": path})
		return
	}
	ctl.Fail(c, errors.InvalidParameter, "请选择文件")

	//file, err := c.FormFile("file")
	//if err != nil {
	//	ctl.Fail(c, errors.InvalidParameter, "请选择文件")
	//	return
	//}
	//path, err := ctl.uploadService.SaveFile(file) // 调用 service 中方法，写入文件并且拿到返回值 path
	//if err != nil {
	//	ctl.Fail(c, errors.ServerErr, err.Error())
	//	return
	//}
	//
	//ctl.Success(c, gin.H{"path": path})
}
