package controller

import (
	"GinAdmin/internal/pkg/errors"
	"GinAdmin/internal/service"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const chunkSize = 5 << 20

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
}

//func (ctl *UploadController) Chunk(c *gin.Context) {
//	index, err := strconv.Atoi(c.Query("index"))
//	if err != nil || index < 0 {
//		ctl.Fail(c, errors.InvalidParameter, "分片序号无效")
//		return
//	}
//
//	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, chunkSize)
//	if err = ctl.uploadService.SaveChunk(index, c.Request.Body); err != nil {
//		ctl.Fail(c, errors.ServerErr, err.Error())
//		return
//	}
//
//	ctl.Success(c, gin.H{"index": index})
//
//}

func (ctl *UploadController) Merge(c *gin.Context) {
	total, err := strconv.Atoi(c.Query("total"))
	filename := c.Query("filename")
	if err != nil || total <= 0 || filename == "" {
		ctl.Fail(c, errors.InvalidParameter, "合并参数无效")
		return
	}

	path, err := ctl.uploadService.Merge(total, filename)
	if err != nil {
		ctl.Fail(c, errors.ServerErr, err.Error())
		return
	}
	ctl.Success(c, gin.H{"path": path})
}

func (ctl *UploadController) Init(c *gin.Context) {
	var req struct {
		FileName string `json:"file_name"`
		FileSize int64  `json:"file_size"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.FileName == "" || req.FileSize <= 0 {
		ctl.Fail(c, errors.InvalidParameter, "上传的参数无效")
		return
	}
	uploadID, chunkTotal, err := ctl.uploadService.Init(ctl.GetCurrentUserID(c), req.FileName, req.FileSize)
	if err != nil {
		ctl.Fail(c, errors.ServerErr, err.Error())
		return
	}
	ctl.Success(c, gin.H{
		"upload_id":   uploadID,
		"chunk_total": chunkTotal,
	})
}

func (ctl *UploadController) Chunk(c *gin.Context) {
	index, err := strconv.Atoi(c.Query("index"))
	uploadID := c.Query("upload_id")
	if err != nil || index < 0 || uploadID == "" {
		ctl.Fail(c, errors.InvalidParameter, "分片参数无效")
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, chunkSize)
	if err = ctl.uploadService.SaveChunk(uploadID, ctl.GetCurrentUserID(c), index, c.Request.Body); err != nil {
		ctl.Fail(c, errors.ServerErr, err.Error())
		return
	}
	ctl.Success(c, gin.H{"index": index})
}

func (ctl *UploadController) Complete(c *gin.Context) {
	uploadID := c.Query("upload_id")
	if uploadID == "" {
		ctl.Fail(c, errors.InvalidParameter, "任务不存在,请求参数 uploadID 不存在")
		return
	}
	path, err := ctl.uploadService.Complete(uploadID, ctl.GetCurrentUserID(c))
	if err != nil {
		ctl.Fail(c, errors.ServerErr, err.Error())
		return
	}
	ctl.Success(c, gin.H{"path": path})
}
