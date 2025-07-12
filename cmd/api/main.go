package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/config"
	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/database"
	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/payments"
	"github.com/emiliosheinz/rinha-de-backend-2025-go/package/queue"
)

func main() {
	config.Init()
	db, err := database.Connect();

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	q := queue.NewInMemoryQueue(100)

	wg := &sync.WaitGroup{}
	queue.NewWorkerPool(5, q.GetJobs(), wg)

	paymentsService := payments.NewPaymentsService(db)
	paymentsHandler := payments.NewPaymentsHandler(q, paymentsService)

	http.HandleFunc("/payments", paymentsHandler.HandleCreatePayment)
	http.HandleFunc("/payments-summary", paymentsHandler.HandleGetSummary)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
