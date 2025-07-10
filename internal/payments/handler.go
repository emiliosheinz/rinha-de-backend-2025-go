package payments

import (
	"fmt"
	"log"
	"net/http"
)

func HandleCreatePayment(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleCreatePayment handler called")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed 😅", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "HandleCreatePayment")
}

func HandleGetSummary(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleGetSummary handler called")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed 😅", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "HandleGetSummary")
}
