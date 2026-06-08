package request

import (
	"strings"

	"github.com/gin-gonic/gin"
)

/**
*  GetQueryParams 提取当前请求的查询参数。
**/
func GetQueryParams(c *gin.Context) map[string]any {
	if c == nil {
		return map[string]any{}
	}

	query := c.Request.URL.Query()

	var queryMap = make(map[string]any, len(query))
	for k := range query {
		queryMap[k] = query.Get(k)
	}

	return queryMap
}

func GetAccessToken(c *gin.Context) string {
	if c == nil {
		return ""
	}
	authorization := c.GetHeader("Authorization")

	return strings.TrimPrefix(authorization, "Bearer ")
}