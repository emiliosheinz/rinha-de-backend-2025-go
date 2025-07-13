package main

import (
	"log"
	"net/http"

	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/config"
	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/database"
	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/health"
	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/payments"
	"github.com/emiliosheinz/rinha-de-backend-2025-go/package/queue"
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

	q := queue.NewRedisQueue(database.RedisClient, payments.PendingPaymentsQueueKey)
	queue.StartWorkerPool(8, q)
	defer q.Close()

	paymentsService := payments.NewPaymentsService(db)
	paymentsHandler := payments.NewPaymentsHandler(q, paymentsService)

	payments.RegisterPaymentsJobs(paymentsService)

	http.HandleFunc("/payments", paymentsHandler.HandleCreatePayment)
	http.HandleFunc("/payments-summary", paymentsHandler.HandleGetSummary)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
