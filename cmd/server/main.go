package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/ordersvc/internal/api"
	"github.com/example/ordersvc/internal/config"
	"github.com/example/ordersvc/internal/repository"
	"github.com/example/ordersvc/internal/service"
	"github.com/example/ordersvc/internal/worker"
	"github.com/example/ordersvc/pkg/cache"
	"github.com/example/ordersvc/pkg/logger"
	"github.com/example/ordersvc/pkg/queue"
)

// FLAW: god function -- main() does server setup, routing, middleware,
// worker init, graceful shutdown, health checks all in one place.
// This function is over 100 lines with multiple responsibilities.
func main() {
	cfg := config.Load()

	logr := logger.New(cfg.LogLevel)
	logr.Info("starting ordersvc", "version", "1.4.2")

	// Database setup with magic numbers
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	db.SetMaxOpenConns(25)         // FLAW: magic number
	db.SetMaxIdleConns(5)          // FLAW: magic number
	db.SetConnMaxLifetime(300 * time.Second) // FLAW: magic number

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}

	// Initialize cache
	cacheClient := cache.NewRedisCache(cfg.RedisAddr, 512)

	// Initialize queue
	memQueue := queue.NewMemoryQueue(100)

	// Initialize repositories
	orderRepo := repository.NewOrderRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	orderSvc := service.NewOrderService(orderRepo, inventoryRepo)
	inventorySvc := service.NewInventoryService(inventoryRepo)
	paymentSvc := service.NewPaymentService(cfg.PaymentGatewayURL)
	notificationSvc := service.NewNotificationService(cfg.SMTPHost, cfg.SMTPPort)

	// Initialize worker dispatcher
	dispatcher := worker.NewDispatcher(memQueue, 10)

	// Setup HTTP handlers
	handler := api.NewHandler(orderSvc, inventorySvc, paymentSvc, notificationSvc)
	router := api.NewRouter(handler)

	// Add middleware
	router.Use(api.LoggingMiddleware(logr))
	router.Use(api.AuthMiddleware(cfg.JWTSecret))
	router.Use(api.RateLimitMiddleware(1000))

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		status := map[string]string{
			"status":  "ok",
			"version": "1.4.2",
			"uptime":  fmt.Sprintf("%v", time.Since(startTime)),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	})

	// Readiness probe
	router.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "database not ready: %v", err)
			return
		}
		w.WriteHeader(200)
		fmt.Fprint(w, "ready")
	})

	// Create HTTP server with magic number timeouts
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  30 * time.Second, // FLAW: magic number
		WriteTimeout: 30 * time.Second, // FLAW: magic number
		IdleTimeout:  120 * time.Second, // FLAW: magic number
	}

	// Start background workers
	go dispatcher.Start()

	// Start cache maintenance goroutine
	go func() {
		ticker := time.NewTicker(60 * time.Second) // FLAW: magic number
		defer ticker.Stop()
		for range ticker.C {
			cacheClient.Cleanup()
		}
	}()

	// Start server in goroutine
	go func() {
		logr.Info("listening", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logr.Info("shutting down server")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced shutdown: %v", err)
	}

	// Close resources
	db.Close()
	dispatcher.Stop()

	logr.Info("server exited")

	// FLAW: dead code -- unreachable logging after server exit
	_ = cacheClient
	_ = memQueue
	_ = userRepo
	_ = notificationSvc
}

var startTime = time.Now()
