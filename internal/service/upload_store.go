package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"GinAdmin/config"
)

// uploadStore 负责分片上传的全部磁盘操作
type uploadStore struct {
	root string
}

func newUploadStore() *uploadStore {
	return &uploadStore{root: config.GetConfig().StorageRoot()}
}

func (s *uploadStore) tmpDir(uploadId string) string {
	return filepath.Join(s.root, "tmp", uploadId)
}

func (s *uploadStore) chunkPath(uploadId string, index int) string {
	return filepath.Join(s.tmpDir(uploadId), strconv.Itoa(index)+".part")
}

// objectRelPath 最终文件的相对路径
// 路径只由哈希生成，不掺入用户提供的文件名，从根上排除路径穿越
func (s *uploadStore) objectRelPath(hash, ext string) string {
	return filepath.Join("files", hash[:2], hash[2:4], hash+ext)
}

func (s *uploadStore) abs(relPath string) string {
	return filepath.Join(s.root, relPath)
}

func (s *uploadStore) prepare(uploadId string) error {
	return os.MkdirAll(s.tmpDir(uploadId), 0o750)
}

// uploadedIndexes 扫描临时目录得到已完整落盘的分片序号
func (s *uploadStore) uploadedIndexes(uploadId string, chunkTotal int) []int {
	entries, err := os.ReadDir(s.tmpDir(uploadId))
	if err != nil {
		return []int{}
	}

	indexes := make([]int, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".part") {
			continue
		}
		index, err := strconv.Atoi(strings.TrimSuffix(name, ".part"))
		if err != nil || index < 0 || index >= chunkTotal {
			continue
		}
		indexes = append(indexes, index)
	}
	sort.Ints(indexes)
	return indexes
}

// missingIndexes 合并前点名，返回缺失的分片序号
func (s *uploadStore) missingIndexes(uploadId string, chunkTotal int) []int {
	uploaded := make(map[int]struct{}, chunkTotal)
	for _, index := range s.uploadedIndexes(uploadId, chunkTotal) {
		uploaded[index] = struct{}{}
	}

	missing := make([]int, 0)
	for index := 0; index < chunkTotal; index++ {
		if _, ok := uploaded[index]; !ok {
			missing = append(missing, index)
		}
	}
	return missing
}

// saveChunk 写入单个分片，返回实际字节数与该片的 sha256
// 先写 .tmp 再原子改名：磁盘上存在的 .part 一定是完整分片，断点续传才敢直接信任目录
// limit 是该分片的期望字节数，多一个字节就判定为异常数据
func (s *uploadStore) saveChunk(uploadId string, index int, src io.Reader, limit int64) (int64, string, error) {
	if err := s.prepare(uploadId); err != nil {
		return 0, "", err
	}

	final := s.chunkPath(uploadId, index)
	tmp := final + ".tmp"
	file, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o640)
	if err != nil {
		return 0, "", err
	}

	hash := sha256.New()
	buf := make([]byte, 64<<10)
	written, err := io.CopyBuffer(io.MultiWriter(file, hash), io.LimitReader(src, limit+1), buf)
	if err == nil {
		err = file.Sync()
	}
	file.Close()

	if err != nil {
		os.Remove(tmp)
		return written, "", err
	}
	if written != limit {
		os.Remove(tmp)
		return written, "", fmt.Errorf("分片字节数异常: 实际 %d, 期望 %d", written, limit)
	}
	if err := os.Rename(tmp, final); err != nil {
		os.Remove(tmp)
		return written, "", err
	}
	return written, hex.EncodeToString(hash.Sum(nil)), nil
}

func (s *uploadStore) removeChunk(uploadId string, index int) error {
	return os.Remove(s.chunkPath(uploadId, index))
}

// merge 按序拼接全部分片，返回组合哈希与总字节数
// 组合哈希 = sha256(每个分片 sha256 的十六进制串按序拼接)，与前端算法保持一致
func (s *uploadStore) merge(uploadId string, chunkTotal int, relPath string) (string, int64, error) {
	dst := s.abs(relPath)
	if err := os.MkdirAll(filepath.Dir(dst), 0o750); err != nil {
		return "", 0, err
	}

	tmp := dst + ".merging"
	file, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o640)
	if err != nil {
		return "", 0, err
	}

	fail := func(err error) (string, int64, error) {
		file.Close()
		os.Remove(tmp)
		return "", 0, err
	}

	composite := sha256.New()
	buf := make([]byte, 128<<10)
	var size int64
	for index := 0; index < chunkTotal; index++ {
		part, err := os.Open(s.chunkPath(uploadId, index))
		if err != nil {
			return fail(err)
		}
		chunkHash := sha256.New()
		written, err := io.CopyBuffer(io.MultiWriter(file, chunkHash), part, buf)
		part.Close()
		if err != nil {
			return fail(err)
		}
		size += written
		composite.Write([]byte(hex.EncodeToString(chunkHash.Sum(nil))))
	}

	if err := file.Sync(); err != nil {
		return fail(err)
	}
	file.Close()
	// 同目录改名是原子的，并发合并同一文件时后到者覆盖，内容一致所以无害
	if err := os.Rename(tmp, dst); err != nil {
		os.Remove(tmp)
		return "", 0, err
	}
	return hex.EncodeToString(composite.Sum(nil)), size, nil
}

// detectMime 读文件头判断真实类型，防止伪造扩展名
func (s *uploadStore) detectMime(relPath string) string {
	file, err := os.Open(s.abs(relPath))
	if err != nil {
		return "application/octet-stream"
	}
	defer file.Close()

	header := make([]byte, 512)
	n, _ := file.Read(header)
	return http.DetectContentType(header[:n])
}

func (s *uploadStore) removeTmp(uploadId string) error {
	return os.RemoveAll(s.tmpDir(uploadId))
}

func (s *uploadStore) removeObject(relPath string) error {
	return os.Remove(s.abs(relPath))
}

// cleanExpiredTmp 兜底清理没有任务记录的孤儿临时目录
// 目录 ModTime 会随新分片写入更新，所以活跃任务不会被误删
func (s *uploadStore) cleanExpiredTmp(ttl time.Duration) int {
	base := filepath.Join(s.root, "tmp")
	entries, err := os.ReadDir(base)
	if err != nil {
		return 0
	}

	deadline := time.Now().Add(-ttl)
	cleaned := 0
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil || info.ModTime().After(deadline) {
			continue
		}
		if os.RemoveAll(filepath.Join(base, entry.Name())) == nil {
			cleaned++
		}
	}
	return cleaned
}
