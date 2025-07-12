package payments

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/config"
)

type ProcessPaymentInput struct {
	CorrelationID string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
}

type ProcessPaymentOutput struct{}

type PaymentsService struct {
	db *sql.DB
}

func NewPaymentsService(db *sql.DB) *PaymentsService {
	return &PaymentsService{
		db: db,
	}
}

// Payment processment still needs to be improved so that it takes into consideration the 
// payment processors health
func (p *PaymentsService) ProcessPayment(input ProcessPaymentInput) (*ProcessPaymentOutput, error) {
	requestedAt := time.Now().UTC()
	payloadMap := map[string]any{
		"correlationId": input.CorrelationID,
		"amount":        input.Amount,
		"requestedAt":   requestedAt.Format(time.RFC3339Nano),
	}

	payload, err := json.Marshal(payloadMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	processorURLs := []string{
		config.ProcessorDefaultURL,
		config.ProcessorFallbackURL,
	}

	var resp *http.Response

	for _, url := range processorURLs {
		resp, err = http.Post(url+"/payments", "application/json", bytes.NewBuffer(payload))
		if err == nil && resp.StatusCode < 400 {
			break
		}
	}

	if err != nil || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed to process payment: %v", err)
	}

	_, err = p.db.Exec(`
		INSERT INTO payments (correlation_id, amount, processed_at)
		VALUES ($1, $2, $3)
	`, input.CorrelationID, input.Amount, requestedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to insert payment into DB: %v", err)
	}

	fmt.Print("Payment processed successfully\n")
	return &ProcessPaymentOutput{}, nil
}
