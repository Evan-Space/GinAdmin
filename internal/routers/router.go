package routers

import (
	"net/http"

	"GinAdmin/internal/controller"
	"GinAdmin/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetRouters 创建引擎并注册所有路由
func SetRouters() *gin.Engine {
	engine := gin.Default()
	engine.Use(middleware.CustomRecovery())
	engine.Use(middleware.RequestMeta())
	engine.Use(middleware.ParseToken())
	engine.Use(middleware.Cors())
	engine.Use(middleware.CustomLogger())
	// 你的全局中间件
	// engine.Use(middleware.Cors(), middleware.Logger())
	RegisterRoutes(engine, AppRouteTree())
	return engine
}

// AppRouteTree 应用完整路由树
func AppRouteTree() RouteGroupDef {
	return RouteGroupDef{
		// Routes: []RouteDef{
		// 	GET("ping", "心跳", AuthNone, func(c *gin.Context) {
		// 		c.String(200, "pong")
		// 	}),
		// },
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
	deptCtrl := controller.NewDeptController()
	dashboardCtrl := controller.NewDashboardController()
	authCtrl := controller.NewAuthController()

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
			{
				Routes: []RouteDef{
					POST("auth/refresh-token", "刷新Token", AuthLogin, authCtrl.RefreshToken),
					GET("auth/check-token", "检查Token", AuthLogin, authCtrl.CheckToken),
				},
			},
			{
				Prefix: "dashboard",
				Routes: []RouteDef{
					GET("overview", "仪表盘概览", AuthLogin, dashboardCtrl.Overview),
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
							GET("list", "用户列表", AuthPerm, adminUserCtrl.List),
							GET("detail", "用户详情", AuthPerm, adminUserCtrl.Detail),
							POST("create", "新增用户", AuthPerm, adminUserCtrl.Create),
							POST("update", "更新用户", AuthPerm, adminUserCtrl.Update),
							POST("delete", "删除用户", AuthPerm, adminUserCtrl.Delete),
							POST("bind-role", "绑定角色", AuthPerm, adminUserCtrl.BindRole),
						},
					},
					{
						Prefix: "role",
						Routes: []RouteDef{
							GET("list", "角色列表", AuthPerm, roleCtrl.List),
							GET("detail", "角色详情", AuthPerm, roleCtrl.Detail),
							POST("create", "新增角色", AuthPerm, roleCtrl.Create),
							POST("update", "更新角色", AuthPerm, roleCtrl.Update),
							POST("delete", "删除角色", AuthPerm, roleCtrl.Delete),
						},
					},

					{
						Prefix: "menu",
						Routes: []RouteDef{
							GET("list", "菜单列表", AuthPerm, menuCtrl.List),
							GET("detail", "菜单详情", AuthPerm, menuCtrl.Detail),
							POST("create", "新增菜单", AuthPerm, menuCtrl.Create),
							POST("update", "更新菜单", AuthPerm, menuCtrl.Update),
							POST("delete", "删除菜单", AuthPerm, menuCtrl.Delete),
						},
					},

					{
						Prefix: "department",
						Routes: []RouteDef{
							GET("list", "部门列表", AuthPerm, deptCtrl.List),
							GET("detail", "部门详情", AuthPerm, deptCtrl.Detail),
							POST("create", "新增部门", AuthPerm, deptCtrl.Create),
							POST("update", "更新部门", AuthPerm, deptCtrl.Update),
							POST("delete", "删除部门", AuthPerm, deptCtrl.Delete),
							POST("bind-role", "部门绑定角色", AuthPerm, deptCtrl.BindRole),
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
