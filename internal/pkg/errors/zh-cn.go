package errors

var zhCNText = map[int]string{
	SUCCESS:             "OK",
	NotLogin:            "请先登录",
	AuthorizationErr:    "暂无权限",
	NotFound:            "资源不存在",
	ServerErr:           "服务器内部错误",
	InvalidParameter:    "参数错误",
	UserDoesNotExist:    "用户不存在",
	UserDisable:         "用户已被禁用",
	UserPasswordWrong:   "用户密码错误",
	LoginFailed:         "登录失败",
	TokenGenerateFailed: "Token 生成失败",
	// 大文件上传相关
	UploadTaskNotFound:      "上传任务不存在或已过期",
	UploadChunkInvalid:      "分片序号非法",
	UploadHashMismatch:      "文件校验失败，内容与声明不一致",
	UploadFileTooLarge:      "文件超出大小限制",
	UploadExtNotAllowed:     "不支持的文件类型",
	UploadTaskProcessing:    "文件正在处理中，请稍后重试",
	UploadStorageErr:        "文件存储失败",
	UploadChunkSizeMismatch: "分片大小不匹配",
}

// GetErrorMessage 根据错误码获取中文错误消息
func GetErrorMessage(code int) string {
	if msg, ok := zhCNText[code]; ok {
		return msg
	}
	return "未知错误"
}
