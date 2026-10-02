package service

import (
	"GinAdmin/config"
	"GinAdmin/internal/pkg/errors"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var allowedExt = map[string]string{
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".pdf":  "application/pdf",
	".txt":  "text/plain",
	".js":   "text/plain",
	".mov":  "video/quicktime",
}

type UploadService struct{}

func NewUploadService() *UploadService {

	return &UploadService{}
}

//func (s *UploadService) SaveFile(file *multipart.FileHeader) (string, error) {
//	dir := filepath.Join(config.GetConfig().BasePath, "uploadFiles") // 拼接保存文件的目标文件夹
//	//  os.MkdirAll 递归创建目录，把目标路径上所有未存在的目录都创建
//	if err := os.MkdirAll(dir, 0o755); err != nil {
//		return "", fmt.Errorf("创建目录失败: %w", err)
//	}
//
//	src, err := file.Open() // 打开这个文件，准备进行读取或者写入，打开后必须 Close 关闭
//	if err != nil {
//		return "", fmt.Errorf("打开文件失败 %w", err)
//	}
//	defer src.Close()
//
//	ext := filepath.Ext(file.Filename)                                // 取文件拓展名，结果带 .
//	baseName := strings.TrimSuffix(filepath.Base(file.Filename), ext) // strings.TrimSuffix 从字符串末尾开始，删除指定 字符串
//	name := baseName + uuid.NewString() + ext
//	dstPath := filepath.Join(dir, name) // filepath.Join把多个路径 拼接成一个完整的路径
//
//	/*
//		Create 创建一个文件，并返回这个文件对象，以供后续操作， dstPath 文件路径
//		如果已经有改文件，则清空该文件内容，然后重新写入
//	*/
//	dst, err := os.Create(dstPath)
//	if err != nil {
//		return "", fmt.Errorf("创建文件失败： %w", err)
//	}
//	defer dst.Close()
//
//	// 把 src 的数据写入 dst 中。
//	if _, err := io.Copy(dst, src); err != nil {
//		return "", fmt.Errorf("写入文件失败, %w", err)
//	}
//
//	return filepath.Join("uploadFiles", name), nil
//}

// 检查上传的文件是否支持
func checkUploadType(filename string, head []byte) error {
	ext := strings.ToLower(filepath.Ext(filename))

	if _, ok := allowedExt[ext]; !ok {
		return errors.NewBusinessError(errors.InvalidParameter, "不支持该文件类型")
	}
	if len(head) == 0 {
		return errors.NewBusinessError(errors.InvalidParameter, "该文件内容为空")
	}

	// 单独处理 .mov 文件，因为 .mov 文件的文件头在第 4 到第8 字节，标准库识别不了
	if ext == ".mov" {
		if len(head) < 8 || string(head[4:8]) != "ftyp" {
			return errors.NewBusinessError(errors.InvalidParameter, "文件内容与文件类型不符合")
		}

		return nil
	}

	kind := http.DetectContentType(head) // 根据文件内容的开头的几个字节，去判断文件类型，比较靠谱
	if i := strings.Index(kind, ";"); i >= 0 {
		kind = kind[:i]
	}
	if kind != allowedExt[ext] {
		return errors.NewBusinessError(errors.InvalidParameter, "文件内容与文件类型不符合")
	}
	return nil

}

func (s *UploadService) SaveStream(filename string, src io.Reader) (string, error) {
	dir := filepath.Join(config.GetConfig().BasePath, "uploadFiles")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("创建目录失败 %w", err)
	}

	ext := filepath.Ext(filename)                                // 取文件拓展名，结果带 .
	baseName := strings.TrimSuffix(filepath.Base(filename), ext) // strings.TrimSuffix 从字符串末尾开始，删除指定 字符串
	name := baseName + uuid.NewString() + ext
	dstPath := filepath.Join(dir, name) // filepath.Join把多个路径 拼接成一个完整的路径

	head := make([]byte, 512)
	n, err := io.ReadFull(src, head)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return "", fmt.Errorf("读取文件头失败: %w", err)
	}
	head = head[:n]
	if err := checkUploadType(filename, head); err != nil {
		return "", err
	}
	src = io.MultiReader(bytes.NewReader(head), src)

	//	/*
	//		Create 创建一个文件，并返回这个文件对象，以供后续操作， dstPath 文件路径
	//		如果已经有改文件，则清空该文件内容，然后重新写入
	//	*/
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败 %w", err)
	}

	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		os.Remove(dstPath)
		return "", fmt.Errorf("写入文件失败 %w", err)
	}

	if err := dst.Close(); err != nil {
		os.Remove(dstPath)
		return "", fmt.Errorf("关闭文件失败 %w", err)
	}

	return filepath.Join("uploadFiles", name), nil
}
