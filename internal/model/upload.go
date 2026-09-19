package model

import "time"

// 上传任务状态
const (
	UploadStatusUploading uint8 = 1
	UploadStatusMerging   uint8 = 2
	UploadStatusCompleted uint8 = 3
	UploadStatusFailed    uint8 = 4
	UploadStatusAborted   uint8 = 5
	UploadStatusExpired   uint8 = 6
)

// UploadTask 大文件分片上传任务
// 分片进度不落库，由临时目录里面的 .part 文件承担，避免文件与数据库状态不一致。
type UploadTask struct {
	BaseModel
	UploadId   string    `json:"upload_id" gorm:"column:upload_id;type:varchar(64);not null;default:'';uniqueIndex;"`
	UserId     uint      `json:"user_id" gorm:"column:user_id;type:int(11) unsigned;not null; default:0;index"`
	FileName   string    `json:"file_name" gorm:"column:file_name;type:varchar(255);not null;default:''"`
	FileSize   int64     `json:"file_size" gorm:"column:file_size;type:bigint unsigned;not null;default:0"`
	FileHash   string    `json:"file_hash" gorm:"column:file_hash;type:char(64);not null;default:'';index"`
	ChunkSize  int64     `json:"chunk_size" gorm:"column:chunk_size;type:int unsigned;not null;default:0"`
	ChunkTotal int       `json:"chunk_total" gorm:"column:chunk_total;type:int unsigned;not null;default:0"`
	Status     uint8     `json:"status" gorm:"column:status;type:tinyint unsigned;not null;default:1"`
	FileId     uint      `json:"file_id" gorm:"column:file_id;type:int unsigned;not null;default:0"`
	ExpiredAt  time.Time `json:"expired_at" gorm:"column:expired_at"`
}

func (UploadTask) TableName() string {
	return "upload_task"
}

type FileObject struct {
	BaseModel
	Hash        string `json:"hash" gorm:"column:hash;type:char(64);not null;default:'';uniqueIndex;"`
	Size        int64  `json:"size" gorm:"column:size;type:bigint unsigned;not null;default:0"`
	Ext         string `json:"ext" gorm:"column:ext;type:varchar(20);not null;default:''"`
	Mime        string `json:"mime" gorm:"column:mime;type:varchar(120);not null;default:''"`
	StoragePath string `json:"storage_path" gorm:"column:storage_path;type:varchar(255);not null;default:''"`
	RefCount    uint   `json:"ref_count" gorm:"column:ref_count;type:int unsigned;not null;default:0"`
}

func (FileObject) TableName() string {
	return "file_object"
}
