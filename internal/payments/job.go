package payments

import (
	"encoding/json"
	"fmt"

	"github.com/emiliosheinz/rinha-de-backend-2025-go/package/queue"
)

const ProcessPayment = "process-payment"

type ProcessPaymentJob struct {
	service *PaymentsService
}

func (j ProcessPaymentJob) Execute(data []byte) error {
	var input ProcessPaymentInput
	if err := json.Unmarshal(data, &input); err != nil {
		return fmt.Errorf("invalid job payload: %w", err)
	}
	_, err := j.service.ProcessPayment(input)
	return err
}

func RegisterPaymentsJobs(paymentsService *PaymentsService) {
	queue.RegisterJob(ProcessPayment, &ProcessPaymentJob{service: paymentsService})
}
