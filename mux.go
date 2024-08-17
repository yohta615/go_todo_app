package main

import "net/http"

// 内部実装に依存しないようhttp.NewServeMuxではなくhttp.Handlerインターフェースで返すことがポイント
func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// 静的解析のエラーを回避するため明示的に戻り値を捨てている
		w.Write([]byte(`{"status": "ok"}`))
	})
	return mux
}
