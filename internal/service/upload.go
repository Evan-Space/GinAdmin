package service

import (
	"GinAdmin/config"
	"GinAdmin/data"
	"GinAdmin/internal/model"
	"GinAdmin/internal/pkg/errors"
	"crypto/sha256"
	"encoding/hex"
	stdErrors "errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

const ChunkSize int64 = 5 * 1024 * 1024
const uploadTaskTTL = 7 * 24 * time.Hour // 未完成上传任务保留 7 天

type InitResult struct {
	UploadId   string `json:"upload_id"`
	ChunkSize  int64  `json:"chunk_size"`
	ChunkTotal int    `json:"chunk_total"`
	Uploaded   []int  `json:"uploaded"`
	Finished   bool   `json:"finished"`
	Path       string `json:"path"`
}
type UploadService struct{}

func NewUploadService() *UploadService {

	return &UploadService{}
}

func (s *UploadService) Init(userId uint, fileName string, fileSize int64, fileHash string) (*InitResult, error) {
	name := filepath.Base(fileName)
	ext := strings.ToLower(filepath.Ext(name))
	if _, ok := allowedExt[ext]; !ok || fileSize <= 0 || len(fileHash) != 64 {
		return nil, errors.NewBusinessError(errors.InvalidParameter, "参数错误")
	}

	var object model.FileObject
	err1 := data.GetDB().Where("hash = ?", fileHash).First(&object).Error
	if err1 == nil {
		if err := data.GetDB().Model(&model.FileObject{}).Where("id = ?", object.ID).Update("ref_count", gorm.Expr("ref_count + 1")).Error; err != nil {
			return nil, fmt.Errorf("更新引用计数失败 %w", err)
		}
		return &InitResult{
			Finished: true,
			Path:     object.StoragePath,
			Uploaded: []int{},
		}, nil
	}

	if !stdErrors.Is(err1, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询文件失败 %w", err1)
	}

	var task model.UploadTask
	err := data.GetDB().Where(
		"user_id = ? AND file_hash = ? AND file_size = ? AND status = ? AND expired_at > ?",
		userId, name, fileHash, model.UploadStatusUploading, time.Now()).First(&task).Error
	if err == nil {
		return s.taskResult(&task)
	}
	if !stdErrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询上传任务失败 %w", err)
	}

	chunkTotal := int((fileSize + ChunkSize - 1) / ChunkSize)
	uploadId := uuid.NewString()

	dir := filepath.Join(config.GetConfig().BasePath, "uploadFiles", "tmp", uploadId)
	if err := os.MkdirAll(dir, 0o755); err != nil { // 创建目标文件夹，用来存放上传文件的位置
		return nil, fmt.Errorf("创建文件夹失败 %w", err)
	}
	task = model.UploadTask{
		UploadId:   uploadId,
		UserId:     userId,
		FileName:   filepath.Base(fileName),
		FileSize:   fileSize,
		FileHash:   fileHash,
		ChunkSize:  ChunkSize,
		ChunkTotal: chunkTotal,
		Status:     model.UploadStatusUploading,
		ExpiredAt:  time.Now().Add(uploadTaskTTL),
	}
	if err := data.GetDB().Create(&task).Error; err != nil {
		os.RemoveAll(dir)
		return nil, fmt.Errorf("创建上传任务失败 %w", err)
	}

	return s.taskResult(&task)

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
	if err != dst.Sync() {
		//对文件调用 fsync。os.File 的写入先进内核页缓存，Sync 会把这些脏页强制落盘，并等待写完。只有返回 nil 才说明数据确实到了磁盘。
		dst.Close()
		os.Remove(tmp)
		return fmt.Errorf("分片刷盘失败 %w", err)
	}
	if err := dst.Close(); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("关闭分片文件失败 %w", err)
	}
	return os.Rename(tmp, final) // 分片内容写入完之后改名
}

func (s *UploadService) MergeChunks(uploadId string, userId uint) (string, error) {
	if _, err := uuid.Parse(uploadId); err != nil { // 校验请求参数
		return "", errors.NewBusinessError(errors.InvalidParameter, "参数错误")
	}

	task, err := s.LoadTask(uploadId, userId)
	if err != nil {
		return "", err
	}

	chunkDir := filepath.Join(config.GetConfig().BasePath, "uploadFiles", "tmp", uploadId) // 拼接保存分片的文件夹路径

	for i := 0; i < task.ChunkTotal; i++ {
		if _, err := os.Stat(filepath.Join(chunkDir, fmt.Sprintf("%d.part", i))); err != nil {
			return "", errors.NewBusinessError(errors.InvalidParameter, fmt.Sprintf("缺少第%d个分片", i))
		}
	}

	filename := task.FileName
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
	staging := filepath.Join(config.GetConfig().BasePath, "uploadFiles", "tmp", uploadId+".merging")
	dst, err := os.Create(staging)
	if err != nil {
		return "", fmt.Errorf("创建文件失败 %w", err)
	}

	sums := make([]string, 0, task.ChunkTotal)
	for i := 0; i < task.ChunkTotal; i++ {
		part, err := os.Open(filepath.Join(chunkDir, fmt.Sprintf("%d.part", i)))
		if err != nil {
			dst.Close()
			os.Remove(staging)
			return "", errors.NewBusinessError(errors.InvalidParameter, "缺少分片")
		}
		h := sha256.New()
		_, err = io.Copy(io.MultiWriter(dst, h), part)
		sums = append(sums, hex.EncodeToString(h.Sum(nil)))
		part.Close()
		if err != nil {
			dst.Close()
			os.Remove(staging)
			return "", fmt.Errorf("合并分片失败 %w", err)
		}
	}

	if err := dst.Close(); err != nil {
		os.Remove(staging)
		return "", fmt.Errorf("关闭文件失败 %w", err)
	}

	gotSum := sha256.Sum256([]byte(strings.Join(sums, "")))
	got := hex.EncodeToString(gotSum[:])

	if !strings.EqualFold(got, task.FileHash) {
		os.Remove(staging)
		data.GetDB().Model(&model.UploadTask{}).Where("upload_id = ?", uploadId).
			Update("status", model.UploadStatusFailed)
		return "", errors.NewBusinessError(errors.InvalidParameter, "文件校验失败")
	}
	name := strings.TrimSuffix(filepath.Base(filename), ext) + uuid.NewString() + ext
	rel := filepath.Join("uploadFiles", name)
	abs := filepath.Join(config.GetConfig().BasePath, rel)
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		os.Remove(staging)
		return "", fmt.Errorf("创建目录失败 %w", err)
	}
	if err := os.Rename(staging, abs); err != nil {
		os.Remove(staging)
		return "", fmt.Errorf("保存文件失败 %w", err)
	}
	obj := model.FileObject{
		Hash:        got,
		Size:        task.FileSize,
		Ext:         ext,
		Mime:        allowedExt[ext],
		StoragePath: filepath.ToSlash(rel),
		RefCount:    1,
	}
	if err := data.GetDB().Create(&obj).Error; err != nil {
		var existing model.FileObject
		if err2 := data.GetDB().Where("hash = ?", got).First(&existing).Error; err2 != nil {
			return "", fmt.Errorf("保存文件记录失败 %w", err)
		}
		os.Remove(abs)
		data.GetDB().Model(&existing).Update("ref_count", gorm.Expr("ref_count + 1"))
		obj = existing
	}
	if err := data.GetDB().Model(&model.UploadTask{}).Where("upload_id = ?", uploadId).
		Updates(map[string]any{
			"status":  model.UploadStatusCompleted,
			"file_id": obj.ID,
		}).Error; err != nil {
		return "", fmt.Errorf("更新任务状态失败 %w", err)
	}
	os.RemoveAll(chunkDir)
	return obj.StoragePath, nil
	// filepath.Base(filename) 取完整路径的最后一部分
	//name := strings.TrimSuffix(filepath.Base(filename), ext) + uuid.NewString() + ext // 先去掉filepath 自带的拓展名字，然后拼接上 uuid + 文件后缀
	//dstPath := filepath.Join(config.GetConfig().BasePath, "uploadFiles", name)        // 拼接最终文件的保存路径，包含最终文件的文件名
	//dst, err := os.Create(dstPath)                                                    // 创建最终的文件，并打开，准备向其中写入文件。
	//if err != nil {
	//	return "", fmt.Errorf("创建文件失败 %w", err)
	//}
	//defer dst.Close() // 延迟关闭
	//
	//for i := 0; i < task.ChunkTotal; i++ {
	//	part, err := os.Open(filepath.Join(chunkDir, fmt.Sprintf("%d.part", i))) // 打开某个分片，获取分片内容
	//	if err != nil {
	//		os.Remove(dstPath)
	//		return "", errors.NewBusinessError(errors.InvalidParameter, "缺少分片")
	//	}
	//	_, err = io.Copy(dst, part)
	//	part.Close() // 读取分片内容完成后，就立马关闭分片。
	//	if err != nil {
	//		os.Remove(dstPath)
	//		return "", fmt.Errorf("合并分片失败 %w", err)
	//	}
	//}
	//os.RemoveAll(chunkDir) // 写入完成后移除所有的临时分片文件
	//return filepath.Join("uploadFiles", name), nil
}

// LoadTask
/*
查询数据库中符合要求的 task 数据
*/
func (s *UploadService) LoadTask(uploadId string, userId uint) (*model.UploadTask, error) {
	if _, err := uuid.Parse(uploadId); err != nil {
		return nil, errors.NewBusinessError(errors.InvalidParameter, "参数错误")
	}
	var task model.UploadTask
	if err := data.GetDB().Where("upload_id = ?", uploadId).First(&task).Error; err != nil {
		return nil, errors.NewBusinessError(errors.InvalidParameter, "上传任务不存在")
	}
	if task.UserId != userId || task.Status != model.UploadStatusUploading {
		return nil, errors.NewBusinessError(errors.InvalidParameter, "权限不足")
	}
	return &task, nil
}

func (s *UploadService) taskResult(task *model.UploadTask) (*InitResult, error) {
	uploaded, err := uploadedIndexes(task.UploadId)
	if err != nil {
		return nil, err
	}
	return &InitResult{
		UploadId:   task.UploadId,
		ChunkSize:  task.ChunkSize,
		ChunkTotal: task.ChunkTotal,
		Uploaded:   uploaded,
		Finished:   task.Status == model.UploadStatusCompleted,
		//Finished:   false,
	}, nil
}

// uploadedIndexes
/*
查已上传的分片序号
*/
func uploadedIndexes(uploadId string) ([]int, error) {
	dir := filepath.Join(config.GetConfig().BasePath, "uploadFiles", "tmp", uploadId)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []int{}, nil
		}
		return nil, fmt.Errorf("读取分片目录失败 %w", err)
	}
	indexes := make([]int, 0)
	for _, entry := range entries {
		name := entry.Name()
		// .part.tmp 的文件是写了一半， 不算一个完整的切片，不算在内
		if entry.IsDir() || !strings.HasSuffix(name, ".part") {
			continue
		}

		index, err := strconv.Atoi(strings.TrimSuffix(name, ".part"))
		if err != nil || index < 0 {
			continue
		}
		indexes = append(indexes, index)
	}
	sort.Ints(indexes)
	return indexes, nil
}

func (s *UploadService) Status(uploadId string, userId uint) (*InitResult, error) {
	task, err := s.LoadTask(uploadId, userId)
	if err != nil {
		return nil, err
	}
	return s.taskResult(task)
}
