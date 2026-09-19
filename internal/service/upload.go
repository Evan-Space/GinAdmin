package service

import (
	"io"
	"path/filepath"
	"strings"
	"time"

	"GinAdmin/config"
	"GinAdmin/data"
	"GinAdmin/internal/model"
	"GinAdmin/internal/pkg/errors"
	"GinAdmin/internal/validator/form"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UploadService struct {
	store *uploadStore
}

func NewUploadService() *UploadService {
	return &UploadService{store: newUploadStore()}
}

// TaskResult init 与 status 的统一返回，前端照着 uploaded 补齐缺失分片即可
type TaskResult struct {
	UploadId   string `json:"upload_id"`
	ChunkSize  int64  `json:"chunk_size"`
	ChunkTotal int    `json:"chunk_total"`
	Uploaded   []int  `json:"uploaded"`
	Finished   bool   `json:"finished"`
	FileId     uint   `json:"file_id"`
	Url        string `json:"url"`
}

// CompleteResult Missing 非空表示还缺分片，前端补传后再次调用
type CompleteResult struct {
	Finished bool   `json:"finished"`
	Missing  []int  `json:"missing"`
	FileId   uint   `json:"file_id"`
	Url      string `json:"url"`
	Size     int64  `json:"size"`
}

// Init 建任务，同时承担秒传探测与断点续传探测
func (s *UploadService) Init(userId uint, params *form.InitUploadForm) (*TaskResult, error) {
	cfg := config.GetConfig()
	if params.FileSize > cfg.MaxFileSize() {
		return nil, errors.NewBusinessError(errors.UploadFileTooLarge)
	}
	ext := strings.ToLower(filepath.Ext(params.FileName))
	if !cfg.ExtAllowed(ext) {
		return nil, errors.NewBusinessError(errors.UploadExtNotAllowed)
	}
	// 组合哈希依赖分片大小，前后端不一致会导致秒传失效，直接拒绝
	if params.ChunkSize != cfg.ChunkSize() {
		return nil, errors.NewBusinessError(errors.UploadChunkSizeMismatch)
	}

	// 1. 秒传：内容相同的文件已存在，只加引用不搬数据
	var object model.FileObject
	err := data.GetDB().Where("hash = ?", params.FileHash).First(&object).Error
	if err == nil {
		data.GetDB().Model(&model.FileObject{}).Where("id = ?", object.ID).
			Update("ref_count", gorm.Expr("ref_count + 1"))
		return &TaskResult{
			Finished: true,
			FileId:   object.ID,
			Url:      object.StoragePath,
			Uploaded: []int{},
		}, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 2. 续传：命中未完成的同一文件任务，扫目录报出已传分片
	var task model.UploadTask
	err = data.GetDB().Where("user_id = ? AND file_hash = ? AND status = ? AND expired_at > ?",
		userId, params.FileHash, model.UploadStatusUploading, time.Now()).First(&task).Error
	if err == nil {
		return s.taskResult(&task), nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 3. 新任务
	task = model.UploadTask{
		UploadId:   uuid.NewString(),
		UserId:     userId,
		FileName:   params.FileName,
		FileSize:   params.FileSize,
		FileHash:   params.FileHash,
		ChunkSize:  params.ChunkSize,
		ChunkTotal: int((params.FileSize + params.ChunkSize - 1) / params.ChunkSize),
		Status:     model.UploadStatusUploading,
		ExpiredAt:  time.Now().Add(cfg.TmpTTLDuration()),
	}
	if err := s.store.prepare(task.UploadId); err != nil {
		return nil, errors.NewBusinessError(errors.UploadStorageErr)
	}
	if err := data.GetDB().Create(&task).Error; err != nil {
		return nil, err
	}
	return s.taskResult(&task), nil
}

// Status 页面刷新后恢复进度
func (s *UploadService) Status(userId uint, uploadId string) (*TaskResult, error) {
	task, err := s.findTask(userId, uploadId)
	if err != nil {
		return nil, err
	}
	return s.taskResult(task), nil
}

// SaveChunk 落盘单个分片，body 是裸二进制流，全程不进内存
func (s *UploadService) SaveChunk(userId uint, params *form.ChunkQuery, body io.Reader) error {
	task, err := s.findTask(userId, params.UploadId)
	if err != nil {
		return err
	}
	switch task.Status {
	case model.UploadStatusUploading:
	case model.UploadStatusCompleted:
		return nil // 已完成，重复上传直接忽略
	default:
		return errors.NewBusinessError(errors.UploadTaskProcessing)
	}
	if params.Index < 0 || params.Index >= task.ChunkTotal {
		return errors.NewBusinessError(errors.UploadChunkInvalid)
	}

	// 除最后一片外必须是整片，最后一片是余数；客户端多送数据会在这里被拦下
	expect := task.ChunkSize
	if params.Index == task.ChunkTotal-1 {
		expect = task.FileSize - int64(params.Index)*task.ChunkSize
	}

	_, sum, err := s.store.saveChunk(task.UploadId, params.Index, body, expect)
	if err != nil {
		return errors.NewBusinessError(errors.UploadChunkSizeMismatch, err.Error())
	}
	// 分片级校验：坏片当场发现，不用等到合并才失败
	if params.ChunkHash != "" && !strings.EqualFold(params.ChunkHash, sum) {
		s.store.removeChunk(task.UploadId, params.Index)
		return errors.NewBusinessError(errors.UploadHashMismatch)
	}
	return nil
}

// Complete 合并分片、校验内容、落库去重
func (s *UploadService) Complete(userId uint, uploadId string) (*CompleteResult, error) {
	task, err := s.findTask(userId, uploadId)
	if err != nil {
		return nil, err
	}

	// 幂等：重复调用直接返回已有结果
	if task.Status == model.UploadStatusCompleted {
		var object model.FileObject
		data.GetDB().Where("id = ?", task.FileId).First(&object)
		return &CompleteResult{
			Finished: true, FileId: task.FileId,
			Url: object.StoragePath, Size: object.Size, Missing: []int{},
		}, nil
	}
	if task.Status != model.UploadStatusUploading {
		return nil, errors.NewBusinessError(errors.UploadTaskProcessing)
	}

	// 点名。缺片不算错误，返回缺失序号让前端补传
	if missing := s.store.missingIndexes(task.UploadId, task.ChunkTotal); len(missing) > 0 {
		return &CompleteResult{Finished: false, Missing: missing}, nil
	}

	// 状态 CAS 抢占，用数据库行锁防止并发重复合并写坏目标文件，不需要额外的分布式锁
	claim := data.GetDB().Model(&model.UploadTask{}).
		Where("upload_id = ? AND status = ?", task.UploadId, model.UploadStatusUploading).
		Update("status", model.UploadStatusMerging)
	if claim.Error != nil {
		return nil, claim.Error
	}
	if claim.RowsAffected == 0 {
		return nil, errors.NewBusinessError(errors.UploadTaskProcessing)
	}

	ext := strings.ToLower(filepath.Ext(task.FileName))
	relPath := s.store.objectRelPath(task.FileHash, ext)
	composite, size, err := s.store.merge(task.UploadId, task.ChunkTotal, relPath)
	if err != nil {
		s.markFailed(task.UploadId)
		return nil, errors.NewBusinessError(errors.UploadStorageErr)
	}

	// 最后一道防线：内容哈希与总大小都必须和 init 时声明的一致
	if !strings.EqualFold(composite, task.FileHash) || size != task.FileSize {
		s.store.removeObject(relPath)
		s.markFailed(task.UploadId)
		return nil, errors.NewBusinessError(errors.UploadHashMismatch)
	}

	object := model.FileObject{
		Hash:        task.FileHash,
		Size:        size,
		Ext:         ext,
		Mime:        s.store.detectMime(relPath),
		StoragePath: relPath,
		RefCount:    1,
	}
	err = data.GetDB().Transaction(func(tx *gorm.DB) error {
		// 并发上传同一文件时，后到者复用已有记录
		if err := tx.Where("hash = ?", object.Hash).FirstOrCreate(&object).Error; err != nil {
			return err
		}
		return tx.Model(&model.UploadTask{}).Where("upload_id = ?", task.UploadId).
			Updates(map[string]any{
				"status":  model.UploadStatusCompleted,
				"file_id": object.ID,
			}).Error
	})
	if err != nil {
		s.markFailed(task.UploadId)
		return nil, err
	}

	s.store.removeTmp(task.UploadId)
	return &CompleteResult{
		Finished: true, FileId: object.ID,
		Url: relPath, Size: size, Missing: []int{},
	}, nil
}

// Abort 用户主动取消，立刻回收临时空间
func (s *UploadService) Abort(userId uint, uploadId string) error {
	task, err := s.findTask(userId, uploadId)
	if err != nil {
		return err
	}
	if task.Status == model.UploadStatusCompleted {
		return nil
	}
	if err := data.GetDB().Model(&model.UploadTask{}).Where("upload_id = ?", task.UploadId).
		Update("status", model.UploadStatusAborted).Error; err != nil {
		return err
	}
	return s.store.removeTmp(task.UploadId)
}

// CleanExpired 清理过期任务及其临时分片，由定时协程调用
func (s *UploadService) CleanExpired() int {
	var tasks []model.UploadTask
	data.GetDB().Where("status = ? AND expired_at < ?",
		model.UploadStatusUploading, time.Now()).Find(&tasks)

	for i := range tasks {
		s.store.removeTmp(tasks[i].UploadId)
		data.GetDB().Model(&model.UploadTask{}).Where("id = ?", tasks[i].ID).
			Update("status", model.UploadStatusExpired)
	}
	s.store.cleanExpiredTmp(config.GetConfig().TmpTTLDuration())
	return len(tasks)
}

// findTask 取任务并校验归属，越权访问在这里统一拦掉
func (s *UploadService) findTask(userId uint, uploadId string) (*model.UploadTask, error) {
	var task model.UploadTask
	err := data.GetDB().Where("upload_id = ?", uploadId).First(&task).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.NewBusinessError(errors.UploadTaskNotFound)
	}
	if err != nil {
		return nil, err
	}
	if task.UserId != userId {
		return nil, errors.NewBusinessError(errors.AuthorizationErr)
	}
	return &task, nil
}

func (s *UploadService) taskResult(task *model.UploadTask) *TaskResult {
	return &TaskResult{
		UploadId:   task.UploadId,
		ChunkSize:  task.ChunkSize,
		ChunkTotal: task.ChunkTotal,
		Uploaded:   s.store.uploadedIndexes(task.UploadId, task.ChunkTotal),
		Finished:   task.Status == model.UploadStatusCompleted,
		FileId:     task.FileId,
	}
}

func (s *UploadService) markFailed(uploadId string) {
	data.GetDB().Model(&model.UploadTask{}).Where("upload_id = ?", uploadId).
		Update("status", model.UploadStatusFailed)
}
