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

// Init 初始化上传任务
func (s *UploadService) Init(userId uint, fileName string, fileSize int64, fileHash string) (*InitResult, error) {
	name := filepath.Base(fileName)                                            // 获取文件名
	ext := strings.ToLower(filepath.Ext(name))                                 // 获取文件扩展名
	if _, ok := allowedExt[ext]; !ok || fileSize <= 0 || len(fileHash) != 64 { // 校验参数，看上传的文件是否支持
		return nil, errors.NewBusinessError(errors.InvalidParameter, "参数错误")
	}

	var object model.FileObject                                           // 声明变量，用来查询数据库中是否已经存在该文件
	err1 := data.GetDB().Where("hash = ?", fileHash).First(&object).Error // 查询 FileObject 表中是否存在 hash值为 fileHash 的数据
	if err1 == nil {                                                      // 如果查到数据，说明此文件之前已经被上传过，则更新该文件的引用计数
		if err := data.GetDB().Model(&model.FileObject{}).Where("id = ?", object.ID).Update("ref_count", gorm.Expr("ref_count + 1")).Error; err != nil { // 更新该文件的引用计数
			return nil, fmt.Errorf("更新引用计数失败 %w", err)
		}
		return &InitResult{ // 返回结果，返回文件的存储路径
			Finished: true,
			Path:     object.StoragePath,
			Uploaded: []int{},
		}, nil
	}

	if !stdErrors.Is(err1, gorm.ErrRecordNotFound) { // 如果不是因为记录不存在这个错误，说明程序出错了，返回错误信息。
		return nil, fmt.Errorf("查询文件失败 %w", err1)
	}

	var task model.UploadTask // 声明变量，用来查询数据库中是否存在该文件的上传任务
	err := data.GetDB().Where(
		"user_id = ? AND file_hash = ? AND file_size = ? AND status = ? AND expired_at > ?",
		userId, fileHash, fileSize, model.UploadStatusUploading, time.Now()).First(&task).Error // 查询数据库中是否存在符合这些参数条件的任务，并且把查到的结果写入声明的 task 变量里面
	if err == nil { // 如果查到数据，说明该文件的上传任务已经存在，则返回上传任务的结果
		return s.taskResult(&task)
	}
	if !stdErrors.Is(err, gorm.ErrRecordNotFound) { // 如果不是因为记录不存在这个错误，说明程序出错了，返回错误信息。
		return nil, fmt.Errorf("查询上传任务失败 %w", err)
	}
	// 如果查不到数据，说明该文件的上传任务不存在，则按照新文件上传逻辑执行。
	chunkTotal := int((fileSize + ChunkSize - 1) / ChunkSize) // 计算文件被分片的总数
	uploadId := uuid.NewString()                              // 生成一个唯一的上传 ID

	dir := filepath.Join(config.GetConfig().BasePath, "uploadFiles", "tmp", uploadId) // 拼接服务端最终保存文件的路径
	if err := os.MkdirAll(dir, 0o755); err != nil {                                   // 创建目标文件夹，用来存放上传文件的位置
		return nil, fmt.Errorf("创建文件夹失败 %w", err)
	}
	task = model.UploadTask{ // 创建一个变量，用来存储上传文件的信息
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
	if err := data.GetDB().Create(&task).Error; err != nil { // 创建一个 sql， 将上传文件的信息写入数据库
		os.RemoveAll(dir) // 如果创建失败，则删除目标文件夹
		return nil, fmt.Errorf("创建上传任务失败 %w", err)
	}

	return s.taskResult(&task) // 返回上传任务的结果

}

func (s *UploadService) SaveChunk(uploadId string, index int, src io.Reader) error {
	if _, err := uuid.Parse(uploadId); err != nil || index < 0 { // 校验请求参数
		return errors.NewBusinessError(errors.InvalidParameter, "参数错误")
	}
	dir := filepath.Join(config.GetConfig().BasePath, "uploadFiles", "tmp", uploadId) // 拼接分片保存位置目录
	if err := os.MkdirAll(dir, 0o755); err != nil {                                   // 创建分片保存位置目录
		return fmt.Errorf("创建目录失败 %w", err)
	}

	final := filepath.Join(dir, fmt.Sprintf("%d.part", index)) // 拼接分片文件的最终名字
	tmp := final + ".tmp"                                      // 分片文件的临时名字，仅供保存临时分片文件时使用。
	dst, err := os.Create(tmp)                                 // 创建分片文件，并打开，准备向其中写入文件。
	if err != nil {
		return fmt.Errorf("创建分片失败 %w", err)
	}
	if _, err := io.Copy(dst, src); err != nil { // 将分片内容写入刚才创建的分片文件
		dst.Close()    // 关闭分片文件
		os.Remove(tmp) // 删除临时分片文件
		return fmt.Errorf("分片文件写入失败 %w", err)
	}
	if err != dst.Sync() { // 对文件调用 fsync。os.File 的写入先进内核页缓存，Sync 会把这些脏页强制落盘，并等待写完。只有返回 nil 才说明数据确实到了磁盘。
		dst.Close()    // 关闭分片文件
		os.Remove(tmp) // 删除临时分片文件
		return fmt.Errorf("分片刷盘失败 %w", err)
	}
	if err := dst.Close(); err != nil { // 关闭分片文件
		os.Remove(tmp) // 删除临时分片文件
		return fmt.Errorf("关闭分片文件失败 %w", err)
	}
	return os.Rename(tmp, final) // 分片内容写入完之后改名
}

func (s *UploadService) MergeChunks(uploadId string, userId uint) (string, error) {
	if _, err := uuid.Parse(uploadId); err != nil { // 校验请求参数
		return "", errors.NewBusinessError(errors.InvalidParameter, "参数错误")
	}

	task, err := s.LoadTask(uploadId, userId) // 加载上传任务，查看请求合并的 uploadId 任务ID 在数据库中是否存在对应的上传任务
	if err != nil {
		return "", err
	}

	chunkDir := filepath.Join(config.GetConfig().BasePath, "uploadFiles", "tmp", uploadId) // 拼接保存分片的文件夹路径

	for i := 0; i < task.ChunkTotal; i++ { // 遍历分片，查看分片是否存在
		// os.Stat 用来获取目录的元信息，如果文件不存在，则返回错误。常用来判断某个目录/文件是否存在。
		if _, err := os.Stat(filepath.Join(chunkDir, fmt.Sprintf("%d.part", i))); err != nil { // 检查分片是否存在
			return "", errors.NewBusinessError(errors.InvalidParameter, fmt.Sprintf("缺少第%d个分片", i))
		}
	}

	filename := task.FileName
	first, err := os.Open(filepath.Join(chunkDir, "0.part")) // 打开第一个分片，准备读取分片内容
	if err != nil {
		return "", errors.NewBusinessError(errors.InvalidParameter, "缺少分片")
	}
	head := make([]byte, 512)                                      // 创建一个切片，用来存储文件头
	n, err := io.ReadFull(first, head)                             // 读取文件头，并存储到 head 切片中
	first.Close()                                                  // 关闭分片文件
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF { // 如果读取文件头失败，则返回错误
		return "", fmt.Errorf("读取文件头失败 %w", err)
	}
	if err := checkUploadType(filename, head[:n]); err != nil { // 检查文件类型是否支持
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(filename))                                                   // 获取文件扩展名
	staging := filepath.Join(config.GetConfig().BasePath, "uploadFiles", "tmp", uploadId+".merging") // 拼接合并分片的临时文件路径
	dst, err := os.Create(staging)                                                                   // 创建合并分片的临时文件（最终文件的临时文件），并打开，准备向其中写入文件。
	if err != nil {
		return "", fmt.Errorf("创建文件失败 %w", err)
	}

	sums := make([]string, 0, task.ChunkTotal) // 创建一个切片，用来存储分片的 SHA256 值
	for i := 0; i < task.ChunkTotal; i++ {     // 遍历分片，读取分片内容
		part, err := os.Open(filepath.Join(chunkDir, fmt.Sprintf("%d.part", i)))
		if err != nil {
			dst.Close()
			os.Remove(staging)
			return "", errors.NewBusinessError(errors.InvalidParameter, "缺少分片")
		}
		h := sha256.New()                                   // 创建一个 SHA256 哈希对象
		_, err = io.Copy(io.MultiWriter(dst, h), part)      // 将分片内容写入合并分片的临时文件，并计算 SHA256 值 （把 part 的值，写入 dst 和 h 中）
		sums = append(sums, hex.EncodeToString(h.Sum(nil))) // 将 SHA256 值添加到 sums 切片中
		part.Close()
		if err != nil {
			dst.Close()
			os.Remove(staging)
			return "", fmt.Errorf("合并分片失败 %w", err)
		}
	} // 遍历完所有分片，所有分片的值都写入到合并分片的临时文件中，计算并存储所有分片的 SHA256 值

	if err := dst.Close(); err != nil { // 关闭合并分片的临时文件
		os.Remove(staging)
		return "", fmt.Errorf("关闭文件失败 %w", err)
	}

	gotSum := sha256.Sum256([]byte(strings.Join(sums, ""))) // 计算所有分片的 SHA256 值
	got := hex.EncodeToString(gotSum[:])                    // 将 SHA256 值转换为字符串

	if !strings.EqualFold(got, task.FileHash) { // 比较计算出的 SHA256 值与数据库中存储的 SHA256 值是否一致
		// 如果不一样，则删除合并分片的临时文件，并更新上传任务状态为 failed
		os.Remove(staging)
		data.GetDB().Model(&model.UploadTask{}).Where("upload_id = ?", uploadId).
			Update("status", model.UploadStatusFailed)
		return "", errors.NewBusinessError(errors.InvalidParameter, "文件校验失败")
	}
	// 如果一样，则继续执行，创建最终文件
	name := strings.TrimSuffix(filepath.Base(filename), ext) + uuid.NewString() + ext // 拼接最终文件的文件名
	rel := filepath.Join("uploadFiles", name)
	abs := filepath.Join(config.GetConfig().BasePath, rel)
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil { // 创建最终文件的保存目录
		os.Remove(staging)
		return "", fmt.Errorf("创建目录失败 %w", err)
	}
	if err := os.Rename(staging, abs); err != nil { // 将合并分片的临时文件按照需求重命名为最终文件
		os.Remove(staging)
		return "", fmt.Errorf("保存文件失败 %w", err)
	}
	//--------------------------------
	obj := model.FileObject{ // 创建一个变量，用来存储文件对象
		Hash:        got,
		Size:        task.FileSize,
		Ext:         ext,
		Mime:        allowedExt[ext],
		StoragePath: filepath.ToSlash(rel),
		RefCount:    1,
	}
	if err := data.GetDB().Create(&obj).Error; err != nil { // 创建一个 sql， 将文件对象写入数据库
		var existing model.FileObject
		if err2 := data.GetDB().Where("hash = ?", got).First(&existing).Error; err2 != nil { // 查询数据库中是否存在符合这些参数条件的文件对象，并且把查到的结果写入声明的 existing 变量里面
			return "", fmt.Errorf("保存文件记录失败 %w", err)
		}
		os.Remove(abs)
		data.GetDB().Model(&existing).Update("ref_count", gorm.Expr("ref_count + 1")) // 更新文件对象的引用计数
		obj = existing
	}

	// 更新上传任务状态为 completed，并设置文件对象ID
	if err := data.GetDB().
		Model(&model.UploadTask{}).
		Where("upload_id = ?", uploadId).
		Updates(map[string]any{
			"status":  model.UploadStatusCompleted,
			"file_id": obj.ID,
		}).Error; err != nil {
		return "", fmt.Errorf("更新任务状态失败 %w", err)
	}
	os.RemoveAll(chunkDir) // 删除所有的临时分片文件
	return obj.StoragePath, nil
}

// LoadTask
/*
查询数据库中符合要求的 task 数据
每次保存 分片之前都需要查一遍，用以区分不同用户上传的文件。
确保当前上传任务存在，并且状态为 uploading， 确保上传任务没有过期。
避免不同用户上传的文件出现混乱。
*/
func (s *UploadService) LoadTask(uploadId string, userId uint) (*model.UploadTask, error) {
	if _, err := uuid.Parse(uploadId); err != nil {
		return nil, errors.NewBusinessError(errors.InvalidParameter, "参数错误")
	}
	var task model.UploadTask                                                                // 声明变量，用来存储查询到的上传任务
	if err := data.GetDB().Where("upload_id = ?", uploadId).First(&task).Error; err != nil { // 查询数据库中是否存在符合这些参数条件的任务，并且把查到的结果写入声明的 task 变量里面
		return nil, errors.NewBusinessError(errors.InvalidParameter, "上传任务不存在")
	}
	if task.UserId != userId || task.Status != model.UploadStatusUploading {
		return nil, errors.NewBusinessError(errors.InvalidParameter, "权限不足")
	}
	return &task, nil
}

// taskResult
/*
taskResult 会扫描磁盘目录 uploadFiles/tmp{uploadId} 目录下，已经上传的分片序号
组装成 InitResult 结构体，返回给前端。

*/
func (s *UploadService) taskResult(task *model.UploadTask) (*InitResult, error) {
	uploaded, err := uploadedIndexes(task.UploadId) // 查询已经上传的分片序号
	if err != nil {
		return nil, err
	}
	return &InitResult{
		UploadId:   task.UploadId,
		ChunkSize:  task.ChunkSize,
		ChunkTotal: task.ChunkTotal,
		Uploaded:   uploaded,
		Finished:   task.Status == model.UploadStatusCompleted,
	}, nil
}

// uploadedIndexes
/*
查已上传的分片序号
忽略 .part.tmp 的文件，因为这些文件是写了一半， 不算一个完整的切片，不算在内。
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
