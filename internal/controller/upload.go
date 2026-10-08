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

	if err := c.ShouldBindJSON(&req); err != nil { // 校验参数
		ctl.Fail(c, errors.InvalidParameter, "参数错误")
		return
	}
	// 调用服务层，初始化上传任务
	result, err := ctl.uploadService.Init(ctl.GetCurrentUserID(c), req.FileName, req.FileSize, req.FileHash)
	if err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, result)
}

func (ctl *UploadController) Chunk(c *gin.Context) {
	uploadId := c.Query("upload_id")             // 获取上传 ID
	index, err := strconv.Atoi(c.Query("index")) // 把 url 中的 index 参数，字符串转换成整数。
	if err != nil {
		ctl.Fail(c, errors.InvalidParameter, "分片序号错误")
		return
	}

	uid := ctl.GetCurrentUserID(c)                         // 获取当前用户 ID
	task, err := ctl.uploadService.LoadTask(uploadId, uid) // 加载上传任务，查看当前数据库中是否存在上传任务
	if err != nil {
		ctl.Err(c, err)
		return
	}
	if index < 0 || index >= task.ChunkTotal { // 校验分片序号是否超出范围
		ctl.Fail(c, errors.InvalidParameter, "分片序号超出范围")
		return
	}

	body := http.MaxBytesReader(c.Writer, c.Request.Body, maxChunkSize) // 限制上传文件的大小
	defer body.Close()
	if err := ctl.uploadService.SaveChunk(uploadId, index, body); err != nil { // 保存分片
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, gin.H{
		"index": index,
	})
}

func (ctl *UploadController) Complete(c *gin.Context) {
	var req struct { // 声明变量，用来接收请求参数
		UploadId string `json:"upload_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.UploadId == "" { // 校验请求参数
		ctl.Fail(c, errors.InvalidParameter, "参数错误")
		return
	}
	path, err := ctl.uploadService.MergeChunks(req.UploadId, ctl.GetCurrentUserID(c)) // 合并分片
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

func (ctl *UploadController) UploadAbort(c *gin.Context) {
	var req struct {
		UploadId string `json:"upload_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.UploadId == "" {
		ctl.Fail(c, errors.InvalidParameter, "参数错误")
		return
	}

	if err := ctl.uploadService.UploadAbort(req.UploadId, ctl.GetCurrentUserID(c)); err != nil {
		ctl.Err(c, err)
		return
	}

	ctl.Success(c, nil)
}
