package payments

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/emiliosheinz/rinha-de-backend-2025-go/package/queue"
)

type PaymentsHandler struct {
	jobsQueue       queue.Queue
	paymentsService *PaymentsService
}

func NewPaymentsHandler(q queue.Queue, ps *PaymentsService) *PaymentsHandler {
	return &PaymentsHandler{
		jobsQueue:       q,
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

	ph.jobsQueue.Enqueue(ProcessPayment, input)
	w.WriteHeader(http.StatusAccepted)
}

func (ph *PaymentsHandler) HandleGetSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	from, to, err := parseFilters(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := SummarizePaymentsInput{From: from, To: to}
	summary, err := ph.paymentsService.SummarizePayments(input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(summary); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func parseFilters(r *http.Request) (*time.Time, *time.Time, error) {
	query := r.URL.Query()
	fromStr := query.Get("from")
	toStr := query.Get("to")

	var from, to *time.Time
	layout := time.RFC3339Nano

	if fromStr != "" {
		parsedFrom, err := time.Parse(layout, fromStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid 'from' parameter format")
		}
		from = &parsedFrom
	}

	if toStr != "" {
		parsedTo, err := time.Parse(layout, toStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid 'to' parameter format")
		}
		to = &parsedTo
	}

	return from, to, nil
}
