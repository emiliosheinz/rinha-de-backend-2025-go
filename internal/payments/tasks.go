package payments

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
)

const (
	TypePaymentProcessing = "payment:process"
)

func NewPaymentProcessingTask(input *ProcessPaymentInput) (*asynq.Task, error) {
	payload, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypePaymentProcessing, payload), nil
}

type PaymentProcessor struct { 
	paymentsService *PaymentsService
}

func NewPaymentsProcessor(paymentsService *PaymentsService) *PaymentProcessor {
	return &PaymentProcessor{
		paymentsService: paymentsService,
	}
}

func (processor *PaymentProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var input ProcessPaymentInput
	if err := json.Unmarshal(t.Payload(), &input); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	_, err := 	processor.paymentsService.ProcessPayment(input)
	return err 
}

