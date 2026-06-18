package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/GokujyouKaisennDonnburi/NatuIve_API/internal/config"
	"github.com/GokujyouKaisennDonnburi/NatuIve_API/internal/server"
)

//	@title			NatuIve API
//	@version		1.0
//	@description	NatuIve のバックエンド API
//	@BasePath		/
func main() {
	// 開発用に .env を読み込む（無ければ環境変数をそのまま使う）。
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	cfg := config.Load()

	r, err := server.NewRouter(cfg)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// SIGINT / SIGTERM を受け取るためのコンテキスト。
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// サーバーを別 goroutine で起動する。
	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// 終了シグナルを待つ。
	<-ctx.Done()
	stop()
	log.Println("shutting down...")

	// 進行中のリクエストを最大 10 秒待ってから終了する。
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped")
}
