package middleware

import (
	"GinAdmin/internal/pkg/utils"

	req "GinAdmin/internal/pkg/request"

	"GinAdmin/internal/global"

	"github.com/gin-gonic/gin"
)

// ParseToken 全局 Token 解析中间件
// - 从 Authorization 头提取 Bearer Token
// - Token 有效：将 UserID 和 Username 注入 context
// - Token 无效或缺失：静默继续（不阻断请求）
func ParseToken() gin.HandlerFunc {

	return func(c *gin.Context) {
		accessToken, err := req.GetAccessToken(c)
		if err != nil || accessToken == "" {
			c.Next()
			return
		}

		claims, err := utils.ParseToken(accessToken)
		if err != nil || claims == nil {
			c.Next()
			return
		}


		c.Set(global.ContextKeyUID, claims.UserID)
		c.Set(global.ContextKeyAuthPrincipal, claims)
		c.Next()
	}

	
}