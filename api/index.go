package handler

import (
	"GinAdmin/cmd/bootstrapx"
	"GinAdmin/internal/routers"
	"net/http"
	"sync"
	_ "time/tzdata" // 函数运行环境可能没有时区数据库，InitializeTimezone 需要它
)

var (
	mu     sync.Mutex
	engine http.Handler
)

func getEngine() (http.Handler, error) {
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

func Handler(w http.ResponseWriter, r *http.Request) {
	h, err := getEngine()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.ServeHTTP(w, r)
}
