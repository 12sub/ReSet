package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/12sub/Reset/internal/cache"
	"github.com/12sub/Reset/internal/config"
	"github.com/12sub/Reset/internal/handler"
	"github.com/12sub/Reset/internal/paystack"
	"github.com/12sub/Reset/internal/python"
	"github.com/12sub/Reset/internal/repository"
	"github.com/12sub/Reset/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("cannot ping db:", err)
	}

	redisCache := cache.NewRedis(cfg.RedisAddr)
	ctx := context.Background()
	if err := redisCache.Ping(ctx); err != nil {
		log.Fatal("cannot connect to redis:", err)
	}
	log.Println("Redis connected")

	pythonClient := python.New(cfg.PythonURL)

	pc := paystack.New(cfg.PaystackSecretKey, redisCache)
	subRepo := repository.NewSubscriptionRepo(db, redisCache)
	trackedRepo := repository.NewTrackedRepo(db)

	subService := service.NewSubscriptionService(subRepo, pc)
	cancelService := service.NewCancellationService(subRepo, pc)
	trackedService := service.NewTrackedService(trackedRepo, pythonClient, pc)

	h := handler.NewSubscriptionHandler(subService, cancelService)
	wh := handler.NewWebhookHandler(subRepo, cfg.PaystackSecretKey)
	healthHandler := handler.NewHealthHandler(redisCache)
	trackedHandler := handler.NewTrackedHandler(trackedService)

	mux := http.NewServeMux()

	// Paystack flow
	mux.HandleFunc("GET /", h.Index)
	mux.HandleFunc("POST /subscriptions", h.Create)
	mux.HandleFunc("POST /subscriptions/cancel", h.Cancel)
	mux.HandleFunc("POST /webhooks/paystack", wh.Handle)

	// Universal detection & cancellation
	mux.HandleFunc("GET /detect", trackedHandler.DetectPage)
	mux.HandleFunc("POST /subscriptions/detect", trackedHandler.Detect)
	mux.HandleFunc("GET /tracked", trackedHandler.List)
	mux.HandleFunc("POST /subscriptions/tracked/cancel", trackedHandler.Cancel)

	// Health
	mux.HandleFunc("GET /health", healthHandler.Handle)

	log.Printf("ReSet server running on http://0.0.0.0:%s", cfg.Port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+cfg.Port, mux))
}