package service

import (
	"GinAdmin/config"
	"GinAdmin/internal/pkg/errors"
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

func (s *UploadService) SaveChunk(uploadId string, index int, src io.Reader) error {
	if _, err := uuid.Parse(uploadId); err != nil || index < 0 {
		return errors.NewBusinessError(errors.InvalidParameter, "参数错误")
	}
	dir := filepath.Join(config.GetConfig().BasePath, "uploadFiles", "tmp", uploadId) // 拼接分片保存位置目录
	if err := os.MkdirAll(dir, 0o755); err != nil {                                   // 创建目录
		return fmt.Errorf("创建目录失败 %w", err)
	}

	final := filepath.Join(dir, fmt.Sprintf("%d.part", index)) // 拼接分片文件的最终名字
	tmp := final + ".tmp"                                      // 分片文件的临时名字，仅供保存临时分片文件时使用。
	dst, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("创建分片失败 %w", err)
	}
	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		os.Remove(tmp)
		return fmt.Errorf("分片文件写入失败 %w", err)
	}
	if err := dst.Close(); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("关闭分片文件失败 %w", err)
	}
	return os.Rename(tmp, final) // 分片内容写入完之后改名
}

func (s *UploadService) MergeChunks(uploadId string, filename string, total int) (string, error) {
	if _, err := uuid.Parse(uploadId); err != nil || total < 0 { // 校验请求参数
		return "", errors.NewBusinessError(errors.InvalidParameter, "参数错误")
	}
	chunkDir := filepath.Join(config.GetConfig().BasePath, "uploadFiles", "tmp", uploadId) // 拼接保存分片的文件夹路径

	first, err := os.Open(filepath.Join(chunkDir, "0.part"))
	if err != nil {
		return "", errors.NewBusinessError(errors.InvalidParameter, "缺少分片")
	}
	head := make([]byte, 512)
	n, err := io.ReadFull(first, head)
	first.Close()
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return "", fmt.Errorf("读取文件头失败 %w", err)
	}
	if err := checkUploadType(filename, head[:n]); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(filename))
	// filepath.Base(filename) 取完整路径的最后一部分
	name := strings.TrimSuffix(filepath.Base(filename), ext) + uuid.NewString() + ext // 先去掉filepath 自带的拓展名字，然后拼接上 uuid + 文件后缀
	dstPath := filepath.Join(config.GetConfig().BasePath, "uploadFiles", name)        // 拼接最终文件的保存路径，包含最终文件的文件名
	dst, err := os.Create(dstPath)                                                    // 创建最终的文件，并打开，准备向其中写入文件。
	if err != nil {
		return "", fmt.Errorf("创建文件失败 %w", err)
	}
	defer dst.Close() // 延迟关闭

	for i := 0; i < total; i++ {
		part, err := os.Open(filepath.Join(chunkDir, fmt.Sprintf("%d.part", i))) // 打开某个分片，获取分片内容
		if err != nil {
			os.Remove(dstPath)
			return "", errors.NewBusinessError(errors.InvalidParameter, "缺少分片")
		}
		_, err = io.Copy(dst, part)
		part.Close() // 读取分片内容完成后，就立马关闭分片。
		if err != nil {
			os.Remove(dstPath)
			return "", fmt.Errorf("合并分片失败 %w", err)
		}
	}
	os.RemoveAll(chunkDir) // 写入完成后移除所有的临时分片文件
	return filepath.Join("uploadFiles", name), nil
}
