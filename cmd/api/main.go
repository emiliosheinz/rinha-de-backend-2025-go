package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/config"
	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/payments"
	"github.com/emiliosheinz/rinha-de-backend-2025-go/package/queue"
)

func main() {
	config.Init()

	q := queue.NewInMemoryQueue(100)

	wg := &sync.WaitGroup{}
	queue.NewWorkerPool(5, q.GetJobs(), wg)

	paymentsHandler := payments.NewPaymentsHandler(q)

	http.HandleFunc("/payments", paymentsHandler.HandleCreatePayment)
	http.HandleFunc("/payments-summary", paymentsHandler.HandleGetSummary)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
