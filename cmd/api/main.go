package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/12sub/Reset/internal/cache"
	"github.com/12sub/Reset/internal/config"
	"github.com/12sub/Reset/internal/handler"
	"github.com/12sub/Reset/internal/paystack"
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

	redisCache := cache.NewRedis(os.Getenv("REDIS_URL"))
	ctx := context.Background()
	if err := redisCache.Ping(ctx); err != nil {
		log.Fatal("cannot connect to redis:", err)
	}
	log.Println("Redis connected")

	pc := paystack.New(cfg.PaystackSecretKey, redisCache)
	subRepo := repository.NewSubscriptionRepo(db, redisCache)
	subService := service.NewSubscriptionService(subRepo, pc)
	cancelService := service.NewCancellationService(subRepo, pc)

	h := handler.NewSubscriptionHandler(subService, cancelService)
	wh := handler.NewWebhookHandler(subRepo, cfg.PaystackSecretKey)
	healthHandler := handler.NewHealthHandler(redisCache)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.Index)
	mux.HandleFunc("POST /subscriptions", h.Create)
	mux.HandleFunc("POST /subscriptions/cancel", h.Cancel)
	mux.HandleFunc("POST /webhooks/paystack", wh.Handle)
	mux.HandleFunc("GET /health", healthHandler.Handle)
	log.Printf("ReSet server running on http://localhost:%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}