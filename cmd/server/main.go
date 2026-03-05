package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	appcfg "github.com/AXONcompany/POS/internal/config"
	apppg "github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres"
	apphttp "github.com/AXONcompany/POS/internal/infrastructure/rest"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/auth"
	httping "github.com/AXONcompany/POS/internal/infrastructure/rest/ingredient"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/order"
	httpowner "github.com/AXONcompany/POS/internal/infrastructure/rest/owner"
	httppayment "github.com/AXONcompany/POS/internal/infrastructure/rest/payment"
	httppos "github.com/AXONcompany/POS/internal/infrastructure/rest/pos"
	httpproduct "github.com/AXONcompany/POS/internal/infrastructure/rest/product"
	httpreport "github.com/AXONcompany/POS/internal/infrastructure/rest/report"
	tableHttp "github.com/AXONcompany/POS/internal/infrastructure/rest/table"
	httpuser "github.com/AXONcompany/POS/internal/infrastructure/rest/user"
	httpvenue "github.com/AXONcompany/POS/internal/infrastructure/rest/venue"
	uauth "github.com/AXONcompany/POS/internal/usecase/auth"
	uing "github.com/AXONcompany/POS/internal/usecase/ingredient"
	uorder "github.com/AXONcompany/POS/internal/usecase/order"
	uowner "github.com/AXONcompany/POS/internal/usecase/owner"
	upayment "github.com/AXONcompany/POS/internal/usecase/payment"
	upos "github.com/AXONcompany/POS/internal/usecase/pos"
	uproducts "github.com/AXONcompany/POS/internal/usecase/product"
	ureport "github.com/AXONcompany/POS/internal/usecase/report"
	tableUsecase "github.com/AXONcompany/POS/internal/usecase/table"
	uuser "github.com/AXONcompany/POS/internal/usecase/user"
	uvenue "github.com/AXONcompany/POS/internal/usecase/venue"
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
	userRepo := apppg.NewUserRepository(db)
	sessionRepo := apppg.NewSessionRepository(db)
	orderRepo := apppg.NewOrderRepository(db)
	paymentRepo := apppg.NewPaymentRepository(db)
	reportRepo := apppg.NewReportRepository(db)
	tableRepo := apppg.NewTableRepository(db)
	ownerRepo := apppg.NewOwnerRepository(db)
	venueRepo := apppg.NewVenueRepository(db)
	posRepo := apppg.NewPOSTerminalRepository(db)

	// Service / Usecase
	ingredientService := uing.NewUsecase(ingredientRepo)
	productService := uproducts.NewUsecase(productRepo, categoryRepo, recipeRepo)
	authUsecase := uauth.NewUsecase(userRepo, sessionRepo, cfg.JWTSecret, ownerRepo, venueRepo)
	orderUsecase := uorder.NewUsecase(orderRepo)
	userUsecase := uuser.NewUsecase(userRepo)
	paymentUsecase := upayment.NewUsecase(paymentRepo)
	reportUsecase := ureport.NewUsecase(reportRepo)
	tableService := tableUsecase.NewUsecase(tableRepo)
	ownerUsecase := uowner.NewUsecase(ownerRepo)
	venueUsecase := uvenue.NewUsecase(venueRepo)
	posUsecase := upos.NewUsecase(posRepo)

	// Handler
	ingredientHandler := httping.NewIngredientHandler(ingredientService)
	productHandler := httpproduct.NewHandler(productService)
	authHandler := auth.NewHandler(authUsecase)
	orderHandler := order.NewHandler(orderUsecase)
	userHandler := httpuser.NewHandler(userUsecase)
	paymentHandler := httppayment.NewHandler(paymentUsecase)
	reportHandler := httpreport.NewHandler(reportUsecase)
	tableHandler := tableHttp.NewHandler(tableService)
	ownerHandler := httpowner.NewHandler(ownerUsecase)
	venueHandler := httpvenue.NewHandler(venueUsecase)
	posHandler := httppos.NewHandler(posUsecase)

	// Router
	router := apphttp.NewRouter(cfg, ingredientHandler, productHandler, authHandler, orderHandler, tableHandler, userHandler, paymentHandler, reportHandler, ownerHandler, venueHandler, posHandler)

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
