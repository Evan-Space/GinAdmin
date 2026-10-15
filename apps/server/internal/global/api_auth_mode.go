package global

type ApiAuthMode uint8

const (
	ApiAuthModeNone ApiAuthMode = iota // ApiAuthModeNone 无需登录，无需 API 权限校验。
	ApiAuthModeLogin                  // ApiAuthModeLogin 需要登录，但无需 API 权限校验。
	ApiAuthModeAuth                   // ApiAuthModeAuth 需要登录且需要 API 权限校验。
)

/*
RequiresLogin 返回该模式是否要求用户先登录。
**/
func (m ApiAuthMode) RequiresLogin() bool {
	return m != ApiAuthModeNone
}


// Label 返回该模式的人类可读名称。
func (m ApiAuthMode) Label() string {
	switch m {
	case ApiAuthModeNone:
		return "无需登录"
	case ApiAuthModeLogin:
		return "需要登录"
	case ApiAuthModeAuth:
		return "需要登录和API权限"
	default:
		return "-"
	}
}


// RequiresAPIPermission 返回该模式是否要求 API 权限校验
func (m ApiAuthMode) RequiresAPIPermission() bool {
	return m == ApiAuthModeAuth
}