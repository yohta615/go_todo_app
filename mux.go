package main

import "net/http"

// 内部実装に依存しないようhttp.NewServeMuxではなくhttp.Handlerインターフェースで返すことがポイント
func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(`{"status": "ok"}`))
	})
	return mux
}
