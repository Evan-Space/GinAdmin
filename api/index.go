package handler

import (
	"net/http"
	_ "time/tzdata" // 函数运行环境可能没有时区数据库，InitializeTimezone 需要它

	"GinAdmin/app"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	h, err := app.Handler()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.ServeHTTP(w, r)
}
