package middleware

import (
	"GinAdmin/global"
	"GinAdmin/internal/errors"
	"GinAdmin/internal/pkg/response"
	"fmt"

	"github.com/gin-gonic/gin"
)

// AdminAuthHandler Casbin 权限校验中间件
// 依赖 ParseToken 中间件预先写入 uid 到 context
// 只对路由树中声明 AuthPerm 的路由生效
func AdminAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取当前用户 ID (ParseToken 中间件写入)
		uid := c.GetUint(global.ContextKeyUID)
		if uid == 0 {
			response.FailCode(c, errors.NotLogin)
			c.Abort()
			return
		}
		// 2. 超级管理直接放行
		if uid == global.SuperAdminId {
			c.Next()
			return
		}

		// 3. 获取 Casbin Enforcer
		enforcer, err := casbinx.GetEnforcer()
		if err != nil {
			response.FailCode(c, errors.ServerErr)
			c.Abort()
			return
		}

		// 4. 执行权限检查
		userKey := fmt.Sprintf("%s%s%d", global.CasbinAdminUserPrefix, global.CasbinSeparator, uid)
		path := c.Request.URL.Path
		method := c.Request.Method

		ok, err := enforcer.Enforce(userKey, path, method)
		if err != nil {
			response.FailCode(c, errors.ServerErr)
			c.Abort()
			return
		}

		if !ok {
            response.FailCode(c, errors.AuthorizationErr)
            c.Abort()
            return
        }

        c.Next()
	}
}