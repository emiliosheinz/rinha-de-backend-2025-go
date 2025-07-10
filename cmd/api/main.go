package main

import (
	"log"
	"net/http"

	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/payments"
)


func main() {
    http.HandleFunc("/payments", payments.HandleCreatePayment)
    http.HandleFunc("/payments-summary", payments.HandleCreatePayment)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
