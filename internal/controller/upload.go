package controller

import (
	"GinAdmin/config"
	"GinAdmin/internal/service"
	"GinAdmin/internal/validator"
	"GinAdmin/internal/validator/form"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	Api
	svc *service.UploadService
}

func NewUploadController() *UploadController {
	return &UploadController{
		svc: service.NewUploadService(),
	}
}

/*
init 初始化上传，返回已经分片清单，或者秒传的结果
*/
func (ctl *UploadController) Init(c *gin.Context) {
	params := &form.InitUploadForm{}
	if err := validator.CheckPostParams(c, params); err != nil {
		return
	}

	result, err := ctl.svc.Init(ctl.GetCurrentUserID(c), params)
	if err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, result)
}

/*
Status 查询上传进度，用于页面刷新后恢复
*/
func (ctl *UploadController) Status(c *gin.Context) {
	params := &form.UploadIdForm{}
	if err := validator.CheckQueryParams(c, params); err != nil {
		return
	}
	result, err := ctl.svc.Status(ctl.GetCurrentUserID(c), params.UploadId)
	if err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, result)
}

// chunk 上传单个分片
// 无信息走查参数，请求体是裸二进制，避免 multipart 解析把整片读进内存
func (ctl *UploadController) Chunk(c *gin.Context) {
	params := &form.ChunkQuery{}
	if err := validator.CheckQueryParams(c, params); err != nil {
		return
	}

	// 请求体硬上限，超出一个字即切断连接
	body := http.MaxBytesReader(c.Writer, c.Request.Body, config.GetConfig().ChunkSize()+1)
	defer body.Close()

	if err := ctl.svc.SaveChunk(ctl.GetCurrentUserID(c), params, body); err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, gin.H{"index": params.Index})
}

// Complete 合并分片
func (ctl *UploadController) Complete(c *gin.Context) {
	params := &form.UploadIdForm{}
	if err := validator.CheckQueryParams(c, params); err != nil {
		return
	}
	result, err := ctl.svc.Complete(ctl.GetCurrentUserID(c), params.UploadId)
	if err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, result)
}

// Abort 取消上传
func (ctl *UploadController) Abort(c *gin.Context) {
	params := &form.UploadIdForm{}
	if err := validator.CheckQueryParams(c, params); err != nil {
		return
	}
	if err := ctl.svc.Abort(ctl.GetCurrentUserID(c), params.UploadId); err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, nil)
}
