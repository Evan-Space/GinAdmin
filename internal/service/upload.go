package service

import (
	"GinAdmin/config"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

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

func (s *UploadService) SaveChunk(index int, body io.Reader) error {
	dir := filepath.Join(config.GetConfig().BasePath, "uploadFiles", "temp")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("创建分片目录失败 %w", err)
	}

	dstPath := filepath.Join(dir, fmt.Sprintf("%d.part", index))
	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("创建分片失败 %w", err)
	}
	defer dst.Close()

	if _, err = io.CopyBuffer(dst, body, make([]byte, 32<<10)); err != nil {
		_ = os.Remove(dstPath)
		return fmt.Errorf("写入分片失败 %w", err)
	}
	return nil

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
