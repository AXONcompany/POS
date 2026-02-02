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
	httping "github.com/AXONcompany/POS/internal/http/ingredient" //http ingredient
	httpproduct "github.com/AXONcompany/POS/internal/http/product"
	apppg "github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres"
	uing "github.com/AXONcompany/POS/internal/usecase/ingredients" //usecase ingredient
	uproducts "github.com/AXONcompany/POS/internal/usecase/products"
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

	// Repository
	ingredientRepo := apppg.NewIngredientRepository(db)
	productRepo := apppg.NewProductRepository(db)
	categoryRepo := apppg.NewCategoryRepository(db)
	recipeRepo := apppg.NewRecipeRepository(db)

	// Service
	ingredientService := uing.NewIngredientService(ingredientRepo)
	productService := uproducts.NewService(productRepo, categoryRepo, recipeRepo)

	// Handler
	ingredientHandler := httping.NewIngredientHandler(ingredientService)
	productHandler := httpproduct.NewHandler(productService)

	// Router
	router := apphttp.NewRouter(cfg, ingredientHandler, productHandler)

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
