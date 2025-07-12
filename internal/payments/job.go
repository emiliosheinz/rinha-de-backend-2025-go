package payments

type ProcessPaymentJob struct {
	payment ProcessPaymentInput
	service *PaymentsService
}

func (j ProcessPaymentJob) Execute() error {
	_, err := j.service.ProcessPayment(j.payment)
	return err
}
