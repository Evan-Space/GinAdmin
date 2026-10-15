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
}

// GetErrorMessage 根据错误码获取中文错误消息
func GetErrorMessage(code int) string {
	if msg, ok := zhCNText[code]; ok {
		return msg
	}
	return "未知错误"
}