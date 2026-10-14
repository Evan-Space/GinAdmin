package app

import (
	"GinAdmin/cmd/bootstrapx"
	"GinAdmin/internal/routers"
	"net/http"
	"sync"
)

var (
	mu     sync.Mutex
	engine http.Handler
)

// Handler 返回初始化好的 HTTP 处理器，供 Vercel 函数入口调用
// Vercel 构建时会改写 api 包的导入路径，入口无法直接导入 internal 包，只能经由这里中转

func Handler() (http.Handler, error) {
	mu.Lock()
	defer mu.Unlock()
	if engine != nil {
		return engine, nil
	}
	if err := bootstrapx.Bootstrap(""); err != nil {
		return nil, err
	}
	engine = routers.SetRouters()
	return engine, nil
}
