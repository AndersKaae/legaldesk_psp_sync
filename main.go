package main

import (
	"io"
	"log"
	"net/http"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Log the webhook payload
	log.Printf("Received webhook: %s\n", string(body))

	// Respond to sender
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
