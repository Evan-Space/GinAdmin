package errors

const (
	SUCCESS          = 0
	FAILURE          = 1
	NotLogin         = 401
	AuthorizationErr = 403
	NotFound         = 404
	ServerErr        = 500
	InvalidParameter = 10000
	UserDoesNotExist = 10001
	UserDisable      = 10002
	// 业务错误码 20000+
	UserPasswordWrong   = 20001
	LoginFailed         = 20002
	TokenGenerateFailed = 20003
	// 大文件上传相关
	UploadTaskNotFound      = 21001
	UploadChunkInvalid      = 21002
	UploadHashMismatch      = 21003
	UploadFileTooLarge      = 21004
	UploadExtNotAllowed     = 21005
	UploadTaskProcessing    = 21006
	UploadStorageErr        = 21007
	UploadChunkSizeMismatch = 21008
)
