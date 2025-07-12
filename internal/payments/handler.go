package payments

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emiliosheinz/rinha-de-backend-2025-go/package/queue"
)

type PaymentsHandler struct {
	jobsQueue queue.Queue
	paymentsService *PaymentsService
}

func NewPaymentsHandler(q queue.Queue, ps *PaymentsService) *PaymentsHandler {
	return &PaymentsHandler{
		jobsQueue: q,
		paymentsService: ps,
	}
}

func (ph *PaymentsHandler) HandleCreatePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input ProcessPaymentInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		fmt.Fprintf(w, "Error decoding request body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	job := &ProcessPaymentJob{
		payment: input,
		service: ph.paymentsService,
	}
	ph.jobsQueue.Enqueue(job)

	w.WriteHeader(http.StatusAccepted)
}

func (ph *PaymentsHandler) HandleGetSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "HandleGetSummary")
}
