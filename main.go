package main

import (
	"encoding/json"
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

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload WebhookPayload

	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Printf("Received webhook: %+v\n", payload)

	// === Process the data ===
	//processWebhook(payload)

	// Respond quickly to the sender
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func main() {
	http.HandleFunc("/webhook/denmark/transactions", webhookHandler)
	http.HandleFunc("/webhook/sweden/transactions", webhookHandler)
	http.HandleFunc("/webhook/norway/transactions", webhookHandler)

	addr := ":6969"
	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
