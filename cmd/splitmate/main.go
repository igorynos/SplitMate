package main

import (
	"context"
	"github.com/igorynos/SplitMate/internal/config"
	"github.com/igorynos/SplitMate/internal/repository"
	"github.com/igorynos/SplitMate/internal/transport/httpapi"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	repo := repository.NewPostgres(db)
	if err = repo.Migrate(ctx); err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: (&httpapi.Server{Repo: repo}).Handler(), ReadHeaderTimeout: 5 * time.Second}
	go func() {
		log.Printf("SplitMate listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdown)
}
