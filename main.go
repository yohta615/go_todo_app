package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func main() {
	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}
	// os.Args[0]は実行内容が入るため、[1]が最初の引数となる
	p := os.Args[1]
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen port %s: %v", p, err)
	}
	if err := run(context.Background(), l); err != nil {
		log.Printf("failed to terminate server: %v", err)
	}
}

func run(ctx context.Context, l net.Listener) error {
	// 実務ではhttp.Serverを経由してHTTPサーバーを起動するのが定番
	s := &http.Server{
		// 引数で受け取ったnet.Listenerを利用するので、
		// Addrフィールドは指定しない
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	// errgroup.Group型を使うと、戻り値にエラーが含まれるゴルーチン間の並行処理の実装を簡単に行える
	// sync.WaitGroup型では別ゴルーチン上で実行する関数から戻り値でエラーを受け取ることができない
	eg, ctx := errgroup.WithContext(ctx)
	// 別ゴルーチンでHTTPサーバーを起動する
	eg.Go(func() error {
		if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})
	// チャネルからの通知（終了通知）を待機する
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		// %+vを使うと構造体のフィールドも表示可能、%vだとvalueのみになる
		log.Printf("faled to shutdown: %+v", err)
	}
	// Goメソッドで起動した別ゴルーチンの終了を待つ。
	// eg.Goで発生したエラーを返す
	return eg.Wait()
}
