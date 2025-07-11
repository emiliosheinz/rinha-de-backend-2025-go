package payments

import (
	"fmt"
	"log"
	"net/http"

	"github.com/emiliosheinz/rinha-de-backend-2025-go/package/queue"
)

type PaymentsHandler struct {
	jobsQueue queue.Queue
}

func NewPaymentsHandler(q queue.Queue) *PaymentsHandler {
	return &PaymentsHandler{
		jobsQueue: q,
	}
}

func (ph *PaymentsHandler) HandleCreatePayment(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Method not allowed 😅", http.StatusMethodNotAllowed)
	// 	return
	// }

	ph.jobsQueue.Enqueue(PaymentJob{Message: "💵 new payment received"})
}

func (ph *PaymentsHandler) HandleGetSummary(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleGetSummary handler called")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed 😅", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "HandleGetSummary")
}
