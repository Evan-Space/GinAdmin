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

func (ctl *UploadController) Chunk(c *gin.Context) {
	uploadId := c.Query("upload_id")
	index, err := strconv.Atoi(c.Query("index")) // 把 url 中的字符串转换成整数。
	if err != nil {
		ctl.Fail(c, errors.InvalidParameter, "分片序号错误")
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
		UploadId   string `json:"upload_id"`
		FileName   string `json:"file_name"`
		ChunkTotal int    `json:"chunk_total"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ctl.Fail(c, errors.InvalidParameter, "参数错误")
		return
	}
	path, err := ctl.uploadService.MergeChunks(req.UploadId, req.FileName, req.ChunkTotal)
	if err != nil {
		ctl.Err(c, err)
		return
	}

	ctl.Success(c, gin.H{"path": path})
}
