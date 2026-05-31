package routers

import (
	"github.com/gin-gonic/gin"
)



type RouteDef struct {
	Method   string            // HTTP 方法：GET, POST, PUT, DELETE 等
	Path     string            // 相对路径，如 "list", ":id"
	Title    string            // 路由标题，用于 API 文档
	Desc     string            // 路由描述，补充 Title 未涵盖的信息
	Auth     AuthMode          // 认证授权模式，使用 AuthModeNone/Login/Auth
	Handlers []gin.HandlerFunc // Gin 处理器链
}

