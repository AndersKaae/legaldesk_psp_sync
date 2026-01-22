package handlers

import (
	"encoding/json"
	"github.com/AndersKaae/legaldesk_psp_sync/database"
	"log"
	"net/http"
	"time"
)

// JSONInvoice is the representation of an invoice for JSON responses.
type JSONInvoice struct {
	ID               string               `json:"id"`
	Handle           string               `json:"handle"`
	Customer         string               `json:"customer"`
	CustomerEmail    *string              `json:"customer_email"` // Use a pointer to handle nulls
	Currency         string               `json:"currency"`
	Created          time.Time            `json:"created"`
	DiscountAmount   int64                `json:"discount_amount"`
	OrgAmount        int64                `json:"org_amount"`
	AmountVAT        int64                `json:"amount_vat"`
	AmountExVAT      int64                `json:"amount_ex_vat"`
	RefundedAmount   int64                `json:"refunded_amount"`
	AuthorizedAmount int64                `json:"authorized_amount"`
	Country          string               `json:"country"`
	Plan             string               `json:"plan"`
	States           database.InvoiceStates `json:"states"`
}

func Invoices(filter string) http.HandlerFunc {
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

		var dbInvoices []database.Invoice

		if filter == "all" {
			dbInvoices, err = database.GetInvoicesByDateRange(from, to.Add(24*time.Hour))
			if err != nil {
				log.Printf("Failed to fetch invoices: %v", err)
				http.Error(w, "Failed to fetch invoices", http.StatusInternalServerError)
				return
			}
		} else if filter == "virtualOffice" {
			dbInvoices, err = database.GetVirtualOfficeInvoicesByDateRange(from, to.Add(24*time.Hour))
			if err != nil {
				log.Printf("Failed to fetch invoices: %v", err)
				http.Error(w, "Failed to fetch invoices", http.StatusInternalServerError)
				return
			}
		}

		// Map database invoices to JSON invoices
		jsonInvoices := make([]JSONInvoice, len(dbInvoices))
		for i, dbInvoice := range dbInvoices {
			var customerEmail *string
			if dbInvoice.CustomerEmail.Valid {
				customerEmail = &dbInvoice.CustomerEmail.String
			}

			jsonInvoices[i] = JSONInvoice{
				ID:               dbInvoice.ID,
				Handle:           dbInvoice.Handle,
				Customer:         dbInvoice.Customer,
				CustomerEmail:    customerEmail,
				Currency:         dbInvoice.Currency,
				Created:          dbInvoice.Created,
				DiscountAmount:   dbInvoice.DiscountAmount,
				OrgAmount:        dbInvoice.OrgAmount,
				AmountVAT:        dbInvoice.AmountVAT,
				AmountExVAT:      dbInvoice.AmountExVAT,
				RefundedAmount:   dbInvoice.RefundedAmount,
				AuthorizedAmount: dbInvoice.AuthorizedAmount,
				Country:          dbInvoice.Country,
				Plan:             dbInvoice.Plan,
				States:           dbInvoice.States,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(jsonInvoices); err != nil {
			log.Printf("Failed to encode invoices to JSON: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
