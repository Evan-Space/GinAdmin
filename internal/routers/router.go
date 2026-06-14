package routers

import (
	"net/http"

	"GinAdmin/internal/controller"
	"GinAdmin/middleware"

	"github.com/gin-gonic/gin"
)

// SetRouters 创建引擎并注册所有路由
func SetRouters() *gin.Engine {
	engine := gin.Default()
	engine.Use(middleware.RequestMeta())
	engine.Use(middleware.ParseToken())
	engine.Use(middleware.Cors())
	// 你的全局中间件
	// engine.Use(middleware.Cors(), middleware.Logger())
	RegisterRoutes(engine, AppRouteTree())
	return engine
}


// AppRouteTree 应用完整路由树
func AppRouteTree() RouteGroupDef {
	return RouteGroupDef{
		Routes: []RouteDef{
			GET("ping", "心跳", AuthNone, func(c *gin.Context) {
				c.String(200, "pong")
			}),
		},
		Children: []RouteGroupDef{
			AdminRouteTree(),
		},
	}
}


 // AdminRouteTree 后台路由
 func AdminRouteTree() RouteGroupDef {
	loginCtrl := controller.NewLoginController()
	userCtrl  := controller.NewUserController()

	return RouteGroupDef{
		Prefix: "api/v1",
		Children: []RouteGroupDef{
			// 不需要登录的
			{
				Routes: []RouteDef{
					POST("login", "登录", AuthNone, loginCtrl.Login),
				},
			},
			// 需要登录的
			{
				// Middleware: []gin.HandlerFunc{AuthMiddleware()},  // 你的登录中间件
				Routes: []RouteDef{
					GET("user/info", "用户信息", AuthLogin, userCtrl.Info),
					GET("user/list", "用户列表", AuthPerm,  userCtrl.List),
				},
			},
		},
	}
}



func GET(path, title string, auth AuthMode, handlers ...gin.HandlerFunc) RouteDef {
	return RouteDef{Method: http.MethodGet, Path: path, Title: title, Auth: auth, Handlers: handlers}
}

func POST(path, title string, auth AuthMode, handlers ...gin.HandlerFunc) RouteDef {
	return RouteDef{Method: http.MethodPost, Path: path, Title: title, Auth: auth, Handlers: handlers}
}