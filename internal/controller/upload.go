package controller

import (
	"GinAdmin/internal/pkg/errors"
	"GinAdmin/internal/service"
	stderrors "errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

const maxUploadSize = 10 * 1024 * 1024 * 1024 // 100MB
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
	rawName := c.GetHeader("X-File-Name")       // 获取文件名字
	filename, err := url.QueryUnescape(rawName) // 进行 url 解码
	if err != nil {
		ctl.Fail(c, errors.ServerErr, "文件名解词出错")
		return
	}
	if filename == "" {
		ctl.Fail(c, errors.InvalidParameter, "请选择文件")
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)
	path, err := ctl.uploadService.SaveStream(filename, c.Request.Body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if stderrors.As(err, &maxErr) {
			ctl.Fail(c, errors.InvalidParameter, fmt.Sprintf("文件过大，最大支持 %d ", maxUploadSize))
			return
		}
		ctl.Fail(c, errors.ServerErr, "服务端保存文件出错")
		return
	}
	ctl.Success(c, gin.H{"path": path})
}
