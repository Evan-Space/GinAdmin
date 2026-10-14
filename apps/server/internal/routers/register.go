package routers

import (
	"strings"

	"github.com/gin-gonic/gin"
)



func RegisterRoutes(engine *gin.Engine, root RouteGroupDef) {
	registerGroup(&engine.RouterGroup, root)
}


func registerGroup(parent *gin.RouterGroup, group RouteGroupDef) {
	current := parent
	if group.Prefix != "" || len(group.Middleware) > 0 {
		current = parent.Group(group.Prefix, group.Middleware...)
	}
	for _, route := range group.Routes {
		current.Handle(route.Method, "/"+strings.Trim(route.Path, "/"), route.Handlers...)
	}
	for _, child := range group.Children {
		registerGroup(current, child)
	}
}