package routers

import (
	"net/http"

	"GinAdmin/global"

	"github.com/gin-gonic/gin"
)

type AuthMode = global.ApiAuthMode

// 路由构建辅助函数（减少重复代码）
func GET(path, title string, auth AuthMode, handlers ...gin.HandlerFunc) RouteDef {
	return RouteDef{Method: http.MethodGet, Path: path, Title: title, Auth: auth, Handlers: handlers}
}

func POST(path, title string, auth AuthMode, handlers ...gin.HandlerFunc) RouteDef {
	return RouteDef{Method: http.MethodPost, Path: path, Title: title, Auth: auth, Handlers: handlers}
}