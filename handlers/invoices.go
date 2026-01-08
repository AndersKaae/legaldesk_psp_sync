package handlers

import (
	"fmt"
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

		fmt.Fprintf(w, "Fetching invoices from %s to %s", from.Format(layout), to.Format(layout))
	}
}
