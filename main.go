package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/yohta615/go_todo_app/config"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// 環境変数にあるデータを取得
	ctg, err := config.New()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", ctg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", ctg.Port, err)
	}

	// urlは確認でのみ使用
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	mux := NewMux()
	s := NewServer(l, mux)
	return s.Run(ctx)
}
