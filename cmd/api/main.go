package main

import (
	"log"
	"net/http"
	"time"

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
	http.HandleFunc("/purge-payments", paymentsHandler.HandlePurgePayments)
	http.HandleFunc("/health", handleHealthCheck)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func startQueueWorkersServer(paymentsService *payments.PaymentsService) {
	srv := asynq.NewServerFromRedisClient(
		database.RedisClient,
		asynq.Config{
			Concurrency:              16,
			DelayedTaskCheckInterval: time.Duration(250) * time.Millisecond,
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				return time.Duration(1<<uint(n-1)) * time.Second
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.Handle(payments.TypePaymentProcessing, payments.NewPaymentsProcessor(paymentsService))

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run worker server: %v", err)
	}
}
