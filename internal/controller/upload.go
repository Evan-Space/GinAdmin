package controller

import (
	"GinAdmin/internal/pkg/errors"
	"GinAdmin/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const maxChunkSize = 5 * 1024 * 1024 // 最大分片大小
type UploadController struct {
	Api
	uploadService *service.UploadService
}

func NewUploadController() *UploadController {
	return &UploadController{
		uploadService: service.NewUploadService(),
	}
}

func (ctl *UploadController) Init(c *gin.Context) {
	var req struct {
		FileName string `json:"file_name"`
		FileSize int64  `json:"file_size"`
		FileHash string `json:"file_hash"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ctl.Fail(c, errors.InvalidParameter, "参数错误")
		return
	}
	result, err := ctl.uploadService.Init(ctl.GetCurrentUserID(c), req.FileName, req.FileSize, req.FileHash)
	if err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, result)
}

func (ctl *UploadController) Chunk(c *gin.Context) {
	uploadId := c.Query("upload_id")
	index, err := strconv.Atoi(c.Query("index")) // 把 url 中的字符串转换成整数。
	if err != nil {
		ctl.Fail(c, errors.InvalidParameter, "分片序号错误")
		return
	}

	uid := ctl.GetCurrentUserID(c)
	task, err := ctl.uploadService.LoadTask(uploadId, uid)
	if err != nil {
		ctl.Err(c, err)
		return
	}
	if index < 0 || index >= task.ChunkTotal {
		ctl.Fail(c, errors.InvalidParameter, "分片序号超出范围")
		return
	}

	body := http.MaxBytesReader(c.Writer, c.Request.Body, maxChunkSize)
	defer body.Close()
	if err := ctl.uploadService.SaveChunk(uploadId, index, body); err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, gin.H{
		"index": index,
	})
}

func (ctl *UploadController) Complete(c *gin.Context) {
	var req struct {
		UploadId string `json:"upload_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.UploadId == "" {
		ctl.Fail(c, errors.InvalidParameter, "参数错误")
		return
	}
	path, err := ctl.uploadService.MergeChunks(req.UploadId, ctl.GetCurrentUserID(c))
	if err != nil {
		ctl.Err(c, err)
		return
	}

	ctl.Success(c, gin.H{"path": path})
}

func (ctl *UploadController) Status(c *gin.Context) {
	result, err := ctl.uploadService.Status(c.Query("upload_id"), ctl.GetCurrentUserID(c))
	if err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, result)
}
