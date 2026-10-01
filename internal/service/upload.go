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
	dir := filepath.Join(config.GetConfig().BasePath, "uploadFiles") // 拼接保存文件的目标文件夹
	//  os.MkdirAll 递归创建目录，把目标路径上所有未存在的目录都创建
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	src, err := file.Open() // 打开这个文件，准备进行读取或者写入，打开后必须 Close 关闭
	if err != nil {
		return "", fmt.Errorf("打开文件失败 %w", err)
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)                                // 取文件拓展名，结果带 .
	baseName := strings.TrimSuffix(filepath.Base(file.Filename), ext) // strings.TrimSuffix 从字符串末尾开始，删除指定 字符串
	name := baseName + uuid.NewString() + ext
	dstPath := filepath.Join(dir, name) // filepath.Join把多个路径 拼接成一个完整的路径

	/*
		Create 创建一个文件，并返回这个文件对象，以供后续操作， dstPath 文件路径
		如果已经有改文件，则清空该文件内容，然后重新写入
	*/
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败： %w", err)
	}
	defer dst.Close()

	// 把 src 的数据写入 dst 中。
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("写入文件失败, %w", err)
	}

	return filepath.Join("uploadFiles", name), nil
}
