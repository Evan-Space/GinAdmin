package service

import (
	"GinAdmin/config"
	"GinAdmin/data"
	"GinAdmin/internal/model"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UploadService struct{}

func NewUploadService() *UploadService {

	return &UploadService{}
}

var allowedExt = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".pdf":  true,
	".txt":  true,
	".js":   true,
	".mov":  true,
}

const uploadChunkSize = 5 << 20

func (s *UploadService) SaveFile(part *multipart.Part) (string, error) {
	ext := strings.ToLower(filepath.Ext(part.FileName()))
	if !allowedExt[ext] {
		return "", fmt.Errorf("改文件类型不支持")
	}

	dir := filepath.Join(config.GetConfig().BasePath, "uploadFiles") // 服务端存储的目标文件夹
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("创建目录失败, %w", err)
	}

	name := uuid.NewString() + ext
	dstPath := filepath.Join(dir, name)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	if _, err := io.CopyBuffer(dst, part, make([]byte, 32<<10)); err != nil {
		_ = os.Remove(dstPath)
		return "", fmt.Errorf("写入文件失败: %w", err)
	}
	return filepath.Join("uploadFiles", name), nil
}

func (s *UploadService) Merge(total int, filename string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filepath.Base(filename)))
	if !allowedExt[ext] {
		return "", fmt.Errorf("不支持改文件类型")
	}
	base := config.GetConfig().BasePath
	tmpDir := filepath.Join(base, "uploadFiles", "temp")
	name := uuid.NewString() + ext
	dstPath := filepath.Join(base, "uploadFiles", name)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败 %w", err)
	}
	defer dst.Close()

	buf := make([]byte, 32<<10)
	for i := 0; i < total; i++ {
		src, openErr := os.Open(filepath.Join(tmpDir, fmt.Sprintf("%d.part", i)))
		if openErr != nil {
			_ = os.Remove(dstPath)
			return "", fmt.Errorf("却少第 %d 片", i)
		}

		_, copyErr := io.CopyBuffer(dst, src, buf)
		_ = src.Close()
		if copyErr != nil {
			_ = os.Remove(dstPath)
			return "", fmt.Errorf("合并失败: %w", copyErr)
		}
	}
	_ = os.RemoveAll(tmpDir)
	return filepath.Join("uploadFiles", name), nil
}

func (s *UploadService) taskDir(uploadID string) string {
	return filepath.Join(config.GetConfig().BasePath, "uploadFiles", "temp", uploadID)
}

func (s *UploadService) findTask(uploadID string, userID uint) (*model.UploadTask, error) {
	if _, err := uuid.Parse(uploadID); err != nil {
		return nil, fmt.Errorf("任务不存在")
	}
	var task model.UploadTask
	if err := data.GetDB().Where("upload_id = ?", uploadID).First(&task).Error; err != nil {
		return nil, fmt.Errorf("任务不存在")
	}
	if task.UserID != userID {
		return nil, fmt.Errorf("权限不足")
	}

	return &task, nil
}

func (s *UploadService) Init(userID uint, fileName string, fileSize int64) (string, int, error) {
	ext := strings.ToLower(filepath.Ext(filepath.Base(fileName)))
	if !allowedExt[ext] {
		return "", 0, fmt.Errorf("不支持该文件类型")
	}

	if fileSize <= 0 {
		return "", 0, fmt.Errorf("文件大小无效")
	}

	uploadID := uuid.NewString()
	chunkTotal := int((fileSize + uploadChunkSize - 1) / uploadChunkSize)

	if err := os.MkdirAll(s.taskDir(uploadID), 0o755); err != nil {
		return "", 0, fmt.Errorf("创建上传任务失败， %w", err)
	}

	task := model.UploadTask{
		UploadID:   uploadID,
		UserID:     userID,
		FileName:   filepath.Base(fileName),
		FileSize:   fileSize,
		ChunkSize:  uploadChunkSize,
		ChunkTotal: chunkTotal,
		Status:     0,
		ExpiredAt:  time.Now().Add(24 * time.Hour),
	}
	if err := data.GetDB().Create(&task).Error; err != nil {
		return "", 0, fmt.Errorf("创建上传任务失败: %w", err)
	}

	return uploadID, chunkTotal, nil
}

func (s *UploadService) SaveChunk(uploadID string, userID uint, index int, body io.Reader) error {
	task, err := s.findTask(uploadID, userID)
	if err != nil {
		return err
	}
	if index < 0 || index >= task.ChunkTotal {
		return fmt.Errorf("序号无效")
	}
	dir := s.taskDir(uploadID)
	if err = os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("创建分片目录失败 %w", err)
	}
	dstPath := filepath.Join(dir, fmt.Sprintf("%d.part", index))
	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("创建文件分片失败: %w", err)
	}
	defer dst.Close()

	if _, err = io.CopyBuffer(dst, body, make([]byte, 32<<10)); err != nil {
		_ = os.Remove(dstPath)
		return fmt.Errorf("写入分片失败 %w", err)
	}
	return nil
}

func (s *UploadService) Complete(uploadID string, userID uint) (string, error) {
	task, err := s.findTask(uploadID, userID)
	if err != nil {
		return "", err
	}
	tmpDir := s.taskDir(uploadID)
	for i := 0; i < task.ChunkTotal; i++ {
		part := filepath.Join(tmpDir, fmt.Sprintf("%d.part", i))
		if _, statErr := os.Stat(part); statErr != nil {
			return "", fmt.Errorf("缺少第 %d 片", i)
		}
	}
	ext := strings.ToLower(filepath.Ext(task.FileName))
	name := uuid.NewString() + ext
	dstPath := filepath.Join(config.GetConfig().BasePath, "uploadFiles", name)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败 %w", err)
	}
	defer dst.Close()

	buf := make([]byte, 32<<10)
	for i := 0; i < task.ChunkTotal; i++ {
		src, openErr := os.Open(filepath.Join(tmpDir, fmt.Sprintf("%d.part", i)))
		if openErr != nil {
			_ = os.Remove(dstPath)
			return "", fmt.Errorf("缺少第%d 片", i)
		}
		_, copyErr := io.CopyBuffer(dst, src, buf)
		_ = src.Close()
		if copyErr != nil {
			_ = os.Remove(dstPath)
			return "", fmt.Errorf("合并失败 %w", copyErr)
		}
	}

	_ = os.RemoveAll(tmpDir)
	if err = data.GetDB().Model(&model.UploadTask{}).Where("upload_id = ?", uploadID).Update("status", 1).Error; err != nil {
		return "", fmt.Errorf("更新任务状态失败: %w", err)
	}
	return filepath.Join("uploadFiles", name), nil
}
