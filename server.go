package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{Handler: mux},
		l:   l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	// errgroup.Group型を使うと、戻り値にエラーが含まれるゴルーチン間の並行処理の実装を簡単に行える
	// sync.WaitGroup型では別ゴルーチン上で実行する関数から戻り値でエラーを受け取ることができない
	eg, ctx := errgroup.WithContext(ctx)
	// 別ゴルーチンでHTTPサーバーを起動する
	eg.Go(func() error {
		if err := s.srv.Serve(s.l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// チャネルからの通知（終了通知）を待機する
	<-ctx.Done()
	if err := s.srv.Shutdown(context.Background()); err != nil {
		// %+vを使うと構造体のフィールドも表示可能、%vだとvalueのみになる
		log.Printf("faled to shutdown: %+v", err)
	}
	// Goメソッドで起動した別ゴルーチンの終了を待つ。
	// eg.Goで発生したエラーを返す
	// グレースフルシャットダウンを待つ
	return eg.Wait()
}
