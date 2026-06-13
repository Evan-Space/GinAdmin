package errors


const (
	SUCCESS 					= 0
	FAILURE 					= 1
	NotLogin 					= 401
	AuthorizationErr 			= 403
	NotFound 					= 404
	ServerErr 					= 500
	InvalidParameter			= 10000
	UserDoesNotExist 			= 10001
	UserDisable 				= 10002
	// 业务错误码 20000+
	UserPasswordWrong 			= 20001
	LoginFailed 				= 20002
	TokenGenerateFailed 		= 20003
)
