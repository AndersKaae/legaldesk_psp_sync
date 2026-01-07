package main

import (
	"encoding/json"
	"github.com/AndersKaae/legaldesk_psp_sync/api"
	"github.com/AndersKaae/legaldesk_psp_sync/config"
	"github.com/AndersKaae/legaldesk_psp_sync/database"
	"io"
	"log"
	"net/http"
	"os"
)

type WebhookPayload struct {
	ID                     string `json:"id"`
	Timestamp              string `json:"timestamp"`
	Signature              string `json:"signature"`
	Invoice                string `json:"invoice"`
	Customer               string `json:"customer"`
	EventType              string `json:"event_type"`
	EventID                string `json:"event_id"`
	PaymentMethod          string `json:"payment_method"`
	PaymentMethodReference string `json:"payment_method_reference"`
	Subscription           string `json:"subscription"`
	Dispute                string `json:"dispute"`
	Transaction            string `json:"transaction"`
	CreditNote             string `json:"credit_note"`
	Credit                 string `json:"credit"`
}

func findStatus(statuses []string, target string) (string, bool) {
	for _, v := range statuses {
		if v == target {
			return v, true
		}
	}
	return "", false
}

func webhookHandler(country string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload WebhookPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Printf("EventType: %s\nReceived webhook for country %s: %+v\n",
			payload.EventType, country, payload)

		// Process the data
		invoiceStatus := []string{"invoice_created", "invoice_authorized", "invoice_settled", "invoice_failed", "invoice_refund"}
		customerStatus := []string{"customer_created", "customer_deleted", "customer_changed"}

		if _, found := findStatus(invoiceStatus, payload.EventType); found {
			log.Printf("Processing invoice event: %s\n", payload.EventType)
			invoice, err := api.GetInvoice(payload.Invoice, country)
			if err != nil {
				log.Printf("Error fetching invoice: %v", err)
				http.Error(w, "Failed to fetch invoice", http.StatusInternalServerError)
				return
			}
			log.Printf("Fetched invoice from API: %+v\n", invoice)
		} else if _, found := findStatus(customerStatus, payload.EventType); found {
			invoice, err := api.GetCustomer(payload.Customer, country)
			if err != nil {
				log.Printf("Error fetching customer: %v", err)
				http.Error(w, "Failed to fetch customer", http.StatusInternalServerError)
				return
			}
			log.Printf("Fetched customer from API: %+v\n", invoice)
		} else {
			log.Printf("Unknown event type: %s\n", payload.EventType)
		}

		// Respond quickly to the sender
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

func main() {
	// Open a file for logging
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Create a multi-writer to write to both stdout and the log file
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	cfg := config.LoadConfig()
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	http.HandleFunc("/webhook/denmark", webhookHandler("DK"))
	http.HandleFunc("/webhook/sweden", webhookHandler("SE"))
	http.HandleFunc("/webhook/norway", webhookHandler("NO"))

	addr := ":6969"
	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
