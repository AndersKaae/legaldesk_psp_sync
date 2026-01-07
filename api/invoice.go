package api

import (
	"encoding/json"
	"fmt"
	"github.com/AndersKaae/legaldesk_psp_sync/config"
	"io"
	"log"
	"net/http"
	"time"
)

// Transaction holds minimal info for state tracking
type Transaction struct {
	State   string    `json:"state"`
	Type    string    `json:"type"`
	Created time.Time `json:"created"`
	Settled time.Time `json:"settled"` // optional
}

// InvoiceResponse models the API response fields we care about
type InvoiceResponse struct {
	ID               string        `json:"id"`
	Handle           string        `json:"handle"`
	Customer         string        `json:"customer"`
	Currency         string        `json:"currency"`
	Created          time.Time     `json:"created"`
	DiscountAmount   int64         `json:"discount_amount"`
	OrgAmount        int64         `json:"org_amount"`
	AmountVAT        int64         `json:"amount_vat"`
	AmountExVAT      int64         `json:"amount_ex_vat"`
	RefundedAmount   int64         `json:"refunded_amount"`
	AuthorizedAmount int64         `json:"authorized_amount"`
	Transactions     []Transaction `json:"transactions"`
}

// InvoiceListResponse models the invoice list API response
type InvoiceListResponse struct {
	Invoices []InvoiceResponse `json:"invoices"` // Assuming "invoices" is the key for the list
	NextPage string            `json:"next_page"`
	HasMore  bool              `json:"has_more"`
}

// InvoiceStates maps state names to timestamps (nil if not occurred)
type InvoiceStates map[string]*time.Time

// Invoice is your domain model
type Invoice struct {
	ID               string
	Handle           string
	Customer         string
	Currency         string
	Created          time.Time
	DiscountAmount   int64
	OrgAmount        int64
	AmountVAT        int64
	AmountExVAT      int64
	RefundedAmount   int64
	AuthorizedAmount int64
	Country          string
	States           InvoiceStates
}

// mapStates converts transactions to InvoiceStates
func mapStates(transactions []Transaction) InvoiceStates {
	states := make(InvoiceStates)

	for _, tx := range transactions {
		var ts *time.Time

		if !tx.Settled.IsZero() {
			ts = &tx.Settled
		} else if !tx.Created.IsZero() {
			ts = &tx.Created
		} else {
			ts = nil
		}

		states[tx.State] = ts
	}

	return states
}

// mapInvoice maps API response -> domain Invoice
func mapInvoice(r InvoiceResponse, country string) Invoice {
	return Invoice{
		ID:               r.ID,
		Handle:           r.Handle,
		Customer:         r.Customer,
		Currency:         r.Currency,
		Created:          r.Created,
		DiscountAmount:   r.DiscountAmount,
		OrgAmount:        r.OrgAmount,
		AmountVAT:        r.AmountVAT,
		AmountExVAT:      r.AmountExVAT,
		RefundedAmount:   r.RefundedAmount,
		AuthorizedAmount: r.AuthorizedAmount,
		Country:          country,
		States:           mapStates(r.Transactions),
	}
}

// GetInvoice fetches an invoice from Frisbii API and returns typed Invoice
func GetInvoice(invoiceId string, country string) (Invoice, error) {
	url := "https://api.frisbii.com/v1/invoice/" + invoiceId
	cfg := config.LoadConfig()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Invoice{}, fmt.Errorf("build request: %w", err)
	}

	switch country {
	case "DK":
		req.SetBasicAuth(cfg.Psp_api_key_dk, "")
	case "SE":
		req.SetBasicAuth(cfg.Psp_api_key_se, "")
	case "NO":
		req.SetBasicAuth(cfg.Psp_api_key_no, "")
	default:
		log.Fatal("unsupported country: " + country)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Invoice{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Invoice{}, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return Invoice{}, fmt.Errorf("invoice %s not found", invoiceId)
	}
	if resp.StatusCode != http.StatusOK {
		return Invoice{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp InvoiceResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return Invoice{}, fmt.Errorf("unmarshal response: %w", err)
	}

	// Country is implied from API usage
	return mapInvoice(apiResp, "DK"), nil
}

func GetInvoiceList(nextPage string, country string) ([]Invoice, string, error) {
	url := "https://api.frisbii.com/v1/list/invoice" // List endpoint

	cfg := config.LoadConfig()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("build request: %w", err)
	}

	// Add next_page token if provided
	if nextPage != "" {
		q := req.URL.Query()
		q.Add("next_page", nextPage)
		req.URL.RawQuery = q.Encode()
	}

	switch country {
	case "DK":
		req.SetBasicAuth(cfg.Psp_api_key_dk, "")
	case "SE":
		req.SetBasicAuth(cfg.Psp_api_key_se, "")
	case "NO":
		req.SetBasicAuth(cfg.Psp_api_key_no, "")
	default:
		return nil, "", fmt.Errorf("unsupported country: %s", country)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp InvoiceListResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, "", fmt.Errorf("unmarshal response: %w", err)
	}

	invoices := make([]Invoice, len(apiResp.Invoices))
	for i, r := range apiResp.Invoices {
		invoices[i] = mapInvoice(r, country) // Pass country here
	}

	return invoices, apiResp.NextPage, nil
}
