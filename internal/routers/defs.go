package routers

import (
	"GinAdmin/global"

	"github.com/gin-gonic/gin"
)

type AuthMode = global.ApiAuthMode

const (
	AuthNone  AuthMode = 0 // 无需登录
	AuthLogin AuthMode = 1 // 只需登录
	AuthPerm  AuthMode = 2 // 需要权限
)

// 一条路由
type RouteDef struct {
	Method   string
	Path     string
	Title    string
	Auth     AuthMode
	Handlers []gin.HandlerFunc
}

// 一个路由组
type RouteGroupDef struct {
	Prefix     string
	Middleware []gin.HandlerFunc
	Routes     []RouteDef
	Children   []RouteGroupDef
}
