package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	appcfg "github.com/AXONcompany/POS/internal/config"
	apphttp "github.com/AXONcompany/POS/internal/http"
)

func main() {

	cfg, err := appcfg.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM)
	<-stop

	log.Printf("server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	log.Printf("server stopped")

}
