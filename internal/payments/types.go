package payments

import "time"

type ProcessPaymentInput struct {
	CorrelationID string  `json:"correlationId"`
	Amount        Decimal `json:"amount"`
}

type ProcessPaymentOutput struct{}

type SummarizePaymentsInput struct {
	From *time.Time `json:"from,omitempty"`
	To   *time.Time `json:"to,omitempty"`
}

type Summary struct {
	TotalRequests int     `json:"totalRequests"`
	TotalAmount   Decimal `json:"totalAmount"`
}

type SummarizePaymentsOutput struct {
	Default  Summary `json:"default"`
	Fallback Summary `json:"fallback"`
}
