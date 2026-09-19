package form

// InitUploadForm 初始化上传参数
// FileSize 用 gte =0 而不是 requires 这样就算空文件，也可以走通流程。
type InitUploadForm struct {
	FileName  string `json:"file_name" binging:"required,max=255"`
	FileSize  int64  `json:"file_size" binging:"gte=0"`
	FileHash  string `json:"file_hash" binding:"required,len=64",hexadecimal`
	ChunkSize int64  `json:"chunk_size" binding:"required,gt=0"`
}

// UploadIdForm 只携带 upload_id 的参数
type UploadIdForm struct {
	UploadId string `form:"upload_id" json:"upload_id" binding:"required,max=64"`
}

// ChunkQuery 分片元信息，走查询参数，请求体留给二进制流
type ChunkQuery struct {
	UploadId  string `form:"upload_id" binding:"required,max=64"`
	Index     int    `form:"index" binding:"gte=0"`
	ChunkHash string `form:"chunk_hash" binding:"omitempty,len=64,hexadecimal"`
}
