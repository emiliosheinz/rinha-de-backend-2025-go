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
	db *sql.DB
}

func NewPaymentsService(db *sql.DB) *PaymentsService {
	return &PaymentsService{
		db: db,
	}
}

// Payment processment still needs to be improved so that it takes into consideration the
// payment processors health and not simply tries the default and then the fallsback.
// additionally, we should add a retry mechanism to the queue so that each job is retried a few times
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

	defaultProcessorHealth, err := health.CheckHealth(DefaultProcessor)

	if err != nil {
		return nil, fmt.Errorf("failed to check default processor health: %v", err)
	}

	fallbackProcessorHealth, err := health.CheckHealth(FallbackProcessor)

	if err != nil {
		return nil, fmt.Errorf("failed to check fallback processor health: %v", err)
	}

	var resp *http.Response
	var processedBy string

	if !defaultProcessorHealth.Failing {
		resp, err = http.Post(config.ProcessorDefaultURL+"/payments", "application/json", bytes.NewBuffer(payload))
		processedBy = DefaultProcessor
	} else if !fallbackProcessorHealth.Failing {
		resp, err = http.Post(config.ProcessorFallbackURL+"/payments", "application/json", bytes.NewBuffer(payload))
		processedBy = FallbackProcessor
	} else {
		return nil, fmt.Errorf("non of the processors are healthy try again latter")
	}

	if err != nil || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed to process payment: %v", err)
	}

	_, err = p.db.Exec(`
		INSERT INTO payments (correlation_id, amount, processed_at, processed_by)
		VALUES ($1, $2, $3, $4)
	`, input.CorrelationID, input.Amount, requestedAt, processedBy)

	if err != nil {
		return nil, fmt.Errorf("failed to insert payment into DB: %v", err)
	}

	fmt.Print("Payment processed successfully\n")
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
		var totalAmount float64

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
		args      []any
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
