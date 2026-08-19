package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/12sub/Reset/internal/handler"
	"github.com/12sub/Reset/internal/paystack"
	"github.com/12sub/Reset/internal/repository"
	"github.com/12sub/Reset/internal/service"
)

func main() {
	godotenv.Load()

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize layers
	pc := paystack.New(os.Getenv("PAYSTACK_SECRET_KEY"))
	subRepo := repository.NewSubscriptionRepo(db)
	subService := service.NewSubscriptionService(subRepo, pc)
	cancelService := service.NewCancellationService(subRepo, pc)

	subHandler := handler.NewSubscriptionHandler(subService, cancelService)
	webhookHandler := handler.NewWebhookHandler(subRepo)

	// Routes
	mux := http.NewServeMux()
	mux.HandleFunc("POST /subscriptions", subHandler.Create)
	mux.HandleFunc("POST /subscriptions/cancel", subHandler.Cancel)
	mux.HandleFunc("GET /subscriptions", subHandler.Get)
	mux.HandleFunc("POST /webhooks/paystack", webhookHandler.Handle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("ReSet server running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}