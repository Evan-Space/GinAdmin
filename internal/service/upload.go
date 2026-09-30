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

func (s *UploadService) SaveFile(file *multipart.FileHeader) (string, error) {
	dir := filepath.Join(config.GetConfig().BasePath, "uploadFiles")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开文件失败 %w", err)
	}
	defer src.Close()

	// 取文件名字
	ext := filepath.Ext(file.Filename)
	baseName := strings.TrimSuffix(filepath.Base(file.Filename), ext)
	name := baseName + uuid.NewString() + ext
	dstPath := filepath.Join(dir, name)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败： %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("写入文件失败, %w", err)
	}

	return filepath.Join("uploadFiles", name), nil
}
