package payments

import "time"

type ProcessPaymentInput struct {
	CorrelationID string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
}

type ProcessPaymentOutput struct{}

type SummarizePaymentsInput struct {
	From *time.Time `json:"from,omitempty"`
	To   *time.Time `json:"to,omitempty"`
}

type Summary struct {
	TotalRequests int     `json:"totalRequests"`
	TotalAmount   float64 `json:"totalAmount"`
}

type SummarizePaymentsOutput struct {
	Default  Summary `json:"default"`
	Fallback Summary `json:"fallback"`
}
