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
	adminUserCtrl := controller.NewAdminUserController()
	roleCtrl := controller.NewRoleController()
	menuCtrl := controller.NewMenuController()

	return RouteGroupDef{
		Prefix: "api/v1",
		Children: []RouteGroupDef{
			// ---- 公开接口（AuthNone）----
				// 不需要登录的
			{
				Routes: []RouteDef{
					POST("login", "登录", AuthNone, loginCtrl.Login),
				},
			},
				// 需要登录的
			{
				Prefix: "admin-user",
				Routes: []RouteDef{
					GET("get", "获取当前用户信息", AuthLogin, adminUserCtrl.GetUserInfo),
					POST("update-profile", "更新个人资料", AuthLogin, adminUserCtrl.UpdateProfile),
				},
			},

			// 需要权限的
			{
				Middleware: []gin.HandlerFunc{middleware.AdminAuthHandler()},
				Children: []RouteGroupDef{
					{
						Prefix: "admin-user",
						Routes: []RouteDef{
							Get("list", "用户列表", AuthPerm, adminUserCtrl.List),
							GET("list", "用户列表", AuthPerm, adminUserCtrl.List),
							GET("detail", "用户详情", AuthPerm, adminUserCtrl.Detail),
							POST("create", "新增用户", AuthLogin, adminUserCtrl.Create),
							POST("update", "更新用户", AuthLogin, adminUserCtrl.Update),
							POST("delete", "删除用户", AuthLogin, adminUserCtrl.Delete),
							POST("bind-role", "绑定角色", AuthLogin, adminUserCtrl.BindRole),
						},
					},
					{
						Prefix: "role",
						Routes: []RouteDef{
							GET("list", "角色列表", AuthLogin, roleCtrl.List),
							GET("detail", "角色详情", AuthLogin, roleCtrl.Detail),
							POST("create", "新增角色", AuthLogin, roleCtrl.Create),
							POST("update", "更新角色", AuthLogin, roleCtrl.Update),
							POST("delete", "删除角色", AuthLogin, roleCtrl.Delete),
						},
					},
		
					{
						Prefix: "menu",
						Routes: []RouteDef{
							GET("list", "菜单列表", AuthLogin, menuCtrl.List),
							GET("detail", "菜单详情", AuthLogin, menuCtrl.Detail),
							POST("create", "新增菜单", AuthLogin, menuCtrl.Create),
							POST("update", "更新菜单", AuthLogin, menuCtrl.Update),
							POST("delete", "删除菜单", AuthLogin, menuCtrl.Delete),
						},
					},
					
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
