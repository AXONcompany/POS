package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	appcfg "github.com/AXONcompany/POS/internal/config"
	apphttp "github.com/AXONcompany/POS/internal/http"
	apppg "github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres"
)

func main() {

	cfg, err := appcfg.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	db, err := apppg.Connect(context.Background(), cfg)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}
	defer db.Close()

	router := apphttp.NewRouter(cfg)

	srv := &http.Server{
		Addr:         cfg.GetHTTPAddr(),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		log.Printf("server starting on %s (env=%s)", cfg.GetHTTPAddr(), cfg.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	stopCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-stopCtx.Done()

	log.Printf("server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	shutdownTimer := time.NewTimer(250 * time.Millisecond)
	defer shutdownTimer.Stop()
	select {
	case <-ctx.Done():
	case <-shutdownTimer.C:
	}
	log.Printf("server stopped")

}
