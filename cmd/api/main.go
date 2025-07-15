package main

import (
	"log"
	"net/http"

	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/config"
	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/database"
	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/health"
	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/payments"
	"github.com/hibiken/asynq"
)

func main() {
	config.Init()
	database.InitRedis()
	defer database.RedisClient.Close()
	db, err := database.ConnectPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	health.NewHealthManager().Start()

	queue := asynq.NewClientFromRedisClient(database.RedisClient)
	defer queue.Close()

	paymentsService := payments.NewPaymentsService(db)
	paymentsHandler := payments.NewPaymentsHandler(queue, paymentsService)

	go startQueueWorkersServer(paymentsService)

	http.HandleFunc("/payments", paymentsHandler.HandleCreatePayment)
	http.HandleFunc("/payments-summary", paymentsHandler.HandleGetSummary)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startQueueWorkersServer(paymentsService *payments.PaymentsService) {
	srv := asynq.NewServerFromRedisClient(
		database.RedisClient,
		asynq.Config{
			Concurrency: 16,
		},
	)

	mux := asynq.NewServeMux()
	mux.Handle(payments.TypePaymentProcessing, payments.NewPaymentsProcessor(paymentsService))

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run worker server: %v", err)
	}
}
