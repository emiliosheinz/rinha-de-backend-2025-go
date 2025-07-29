package payments

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/config"
	"github.com/emiliosheinz/rinha-de-backend-2025-go/internal/health"
)

type PaymentsService struct {
	db         *sql.DB
	httpClient *http.Client
}

func NewPaymentsService(db *sql.DB) *PaymentsService {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 50,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  true,
	}

	return &PaymentsService{
		db: db,
		httpClient: &http.Client{
			Timeout:   5 * time.Second,
			Transport: transport,
		},
	}
}

func (p *PaymentsService) ProcessPayment(input ProcessPaymentInput) (*ProcessPaymentOutput, error) {
	requestedAt := time.Now().UTC()
	payloadMap := map[string]any{
		"correlationId": input.CorrelationID,
		"amount":        json.Number(input.Amount.String()),
		"requestedAt":   requestedAt.Format(time.RFC3339Nano),
	}

	payload, err := json.Marshal(payloadMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	defaultProcessorHealth, err := health.CheckHealth(DefaultProcessor)
	if err != nil {
		return nil, fmt.Errorf("failed to check default processor health: %v", err)
	}

	fallbackProcessorHealth, err := health.CheckHealth(FallbackProcessor)
	if err != nil {
		return nil, fmt.Errorf("failed to check fallback processor health: %v", err)
	}

	var processor string
	var processorURL string

	if shouldUseFallbackProcessor(defaultProcessorHealth, fallbackProcessorHealth) {
		processor = FallbackProcessor
		processorURL = config.ProcessorFallbackURL
	} else {
		processor = DefaultProcessor
		processorURL = config.ProcessorDefaultURL
	}

	resp, err := p.httpClient.Post(processorURL+"/payments", "application/json", bytes.NewBuffer(payload))

	if err != nil {
		return nil, fmt.Errorf("failed to process payment: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		health.SetAsFailing(processor)
		return nil, fmt.Errorf("processor returned status %d", resp.StatusCode)
	} else {
		health.SetAsSucceeding(processor)
	}

	_, err = p.db.Exec(`
		INSERT INTO payments (correlation_id, amount, processed_at, processed_by)
		VALUES ($1, $2, $3, $4)
	`, input.CorrelationID, input.Amount, requestedAt, processor)

	if err != nil {
		return nil, fmt.Errorf("failed to insert payment into DB: %v", err)
	}

	return &ProcessPaymentOutput{}, nil
}

func (p *PaymentsService) SummarizePayments(input SummarizePaymentsInput) (*SummarizePaymentsOutput, error) {
	query, args := buildQuery(input)
	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query payments: %v", err)
	}
	defer rows.Close()

	summary := SummarizePaymentsOutput{
		Default:  Summary{},
		Fallback: Summary{},
	}

	for rows.Next() {
		var processor string
		var totalRequests int
		var totalAmount Decimal

		if err := rows.Scan(&processor, &totalRequests, &totalAmount); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		switch processor {
		case DefaultProcessor:
			summary.Default.TotalRequests = totalRequests
			summary.Default.TotalAmount = totalAmount
		case FallbackProcessor:
			summary.Fallback.TotalRequests = totalRequests
			summary.Fallback.TotalAmount = totalAmount
		}
	}

	return &summary, nil
}

func buildQuery(input SummarizePaymentsInput) (string, []any) {
	var (
		args       []any
		conditions []string
	)

	query := strings.Builder{}
	query.WriteString(`
		SELECT processed_by, COUNT(*), SUM(amount)
		FROM payments
	`)

	if input.From != nil {
		args = append(args, *input.From)
		conditions = append(conditions, fmt.Sprintf("processed_at >= $%d", len(args)))
	}

	if input.To != nil {
		args = append(args, *input.To)
		conditions = append(conditions, fmt.Sprintf("processed_at <= $%d", len(args)))
	}

	if len(conditions) > 0 {
		query.WriteString("WHERE ")
		query.WriteString(strings.Join(conditions, " AND "))
	}

	query.WriteString(" GROUP BY processed_by")

	return query.String(), args
}

func shouldUseFallbackProcessor(defaultHealth, fallbackHealth *health.HealthResponse) bool {
	// If both are failing we want to try the default processor 
	if fallbackHealth.Failing && defaultHealth.Failing {
		return false
	}

	// If only fallback is failing we want to request the default processor 
	if fallbackHealth.Failing {
		return false
	}

	// If only default is failing we want to request the fallback
	if defaultHealth.Failing {
		return true
	}

	// If both are working but the default is 25% slower than the fallback 
	// we want to request the fallback, otherwise, the default processor
	return float64(defaultHealth.MinResponseTime) > 1.25*float64(fallbackHealth.MinResponseTime)
}

func (p *PaymentsService) PurgePayments() error {
	_, err := p.db.Exec("DELETE FROM payments")
	if err != nil {
		return fmt.Errorf("failed to purge payments: %v", err)
	}
	return nil
}
