package payments

import (
	"fmt"
	"math/rand"
	"time"
)

type PaymentJob struct {
	Message string
}

func (p PaymentJob) Execute() error {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(1000)
	time.Sleep(time.Duration(r) * time.Millisecond)
	fmt.Println(time.Now().Format("15:04:05.99999999"), "-", p.Message)
	return nil
}
