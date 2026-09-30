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
	"js":    true,
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
