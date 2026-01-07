package main

import (
	"encoding/json"
	"fmt"
	"github.com/AndersKaae/legaldesk_psp_sync/api"
	"log"
	"net/http"
)

type WebhookPayload struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Signature string `json:"signature"`
	Invoice   string `json:"invoice"`
	Customer  string `json:"customer"`
	EventType string `json:"event_type"`
	EventID   string `json:"event_id"`
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

		fmt.Println(payload.EventType)
		fmt.Printf("Received webhook for country %s: %+v\n", country, payload)

		// === Process the data ===
		invoice, err := api.GetInvoice(payload.Invoice, country)
		if err != nil {
			log.Printf("Error fetching invoice: %v", err)
			http.Error(w, "Failed to fetch invoice", http.StatusInternalServerError)
			return
		}

		fmt.Printf("Fetched invoice from API: %+v\n", invoice)

		// Respond quickly to the sender
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

func main() {
	http.HandleFunc("/webhook/denmark", webhookHandler("DK"))
	http.HandleFunc("/webhook/sweden", webhookHandler("SE"))
	http.HandleFunc("/webhook/norway", webhookHandler("NO"))

	addr := ":6969"
	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
