package handlers

import (
	"encoding/json"
	"github.com/AndersKaae/legaldesk_psp_sync/database"
	"log"
	"net/http"
	"time"
)

func Invoices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		fromStr := r.URL.Query().Get("from")
		toStr := r.URL.Query().Get("to")

		if fromStr == "" || toStr == "" {
			http.Error(w, "Missing 'from' or 'to' query parameter", http.StatusBadRequest)
			return
		}

		// Parse and validate dates
		const layout = "2006-01-02" // Go's reference date for YYYY-MM-DD
		from, err := time.Parse(layout, fromStr)
		if err != nil {
			http.Error(w, "'from' date is invalid, must be YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		to, err := time.Parse(layout, toStr)
		if err != nil {
			http.Error(w, "'to' date is invalid, must be YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		// Optional: check if from is before to
		if from.After(to) {
			http.Error(w, "'from' date cannot be after 'to' date", http.StatusBadRequest)
			return
		}

		invoices, err := database.GetInvoicesByDateRange(from, to)
		if err != nil {
			log.Printf("Failed to fetch invoices: %v", err)
			http.Error(w, "Failed to fetch invoices", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(invoices); err != nil {
			log.Printf("Failed to encode invoices to JSON: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
