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

// CustomerResponse models the customer API response.
type CustomerResponse struct {
	ActiveSubscriptions             int       `json:"active_subscriptions"`
	Address                         string    `json:"address"`
	Address2                        string    `json:"address2"`
	CancelledAmount                 int64     `json:"cancelled_amount"`
	CancelledInvoices               int       `json:"cancelled_invoices"`
	CancelledSubscriptions          int       `json:"cancelled_subscriptions"`
	City                            string    `json:"city"`
	Company                         string    `json:"company"`
	Country                         string    `json:"country"`
	Created                         time.Time `json:"created"`
	DunningAmount                   int64     `json:"dunning_amount"`
	DunningInvoices                 int       `json:"dunning_invoices"`
	Email                           string    `json:"email"`
	ExpiredSubscriptions            int       `json:"expired_subscriptions"`
	FailedAmount                    int64     `json:"failed_amount"`
	FailedInvoices                  int       `json:"failed_invoices"`
	FirstName                       string    `json:"first_name"`
	Handle                          string    `json:"handle"`
	LastName                        string    `json:"last_name"`
	NonRenewingSubscriptions        int       `json:"non_renewing_subscriptions"`
	OnHoldSubscriptions             int       `json:"on_hold_subscriptions"`
	PendingAdditionalCostAmount     int64     `json:"pending_additional_cost_amount"`
	PendingAdditionalCosts          int       `json:"pending_additional_costs"`
	PendingAmount                   int64     `json:"pending_amount"`
	PendingCreditAmount             int64     `json:"pending_credit_amount"`
	PendingCredits                  int       `json:"pending_credits"`
	PendingInvoices                 int       `json:"pending_invoices"`
	Phone                           string    `json:"phone"`
	PostalCode                      string    `json:"postal_code"`
	RefundedAmount                  int64     `json:"refunded_amount"`
	SettledAmount                   int64     `json:"settled_amount"`
	SettledInvoices                 int       `json:"settled_invoices"`
	Subscriptions                   int       `json:"subscriptions"`
	Test                            bool      `json:"test"`
	TransferredAdditionalCostAmount int64     `json:"transferred_additional_cost_amount"`
	TransferredAdditionalCosts      int       `json:"transferred_additional_costs"`
	TransferredCreditAmount         int64     `json:"transferred_credit_amount"`
	TransferredCredits              int       `json:"transferred_credits"`
	TrialActiveSubscriptions        int       `json:"trial_active_subscriptions"`
	TrialCancelledSubscriptions     int       `json:"trial_cancelled_subscriptions"`
}

// CustomerListResponse models the customer list API response
type CustomerListResponse struct {
	Customers []CustomerResponse `json:"content"` // Assuming "content" is the key for the list
	NextPage  string             `json:"next_page_token"`
	HasMore   bool               `json:"has_more"` // Assuming API returns this for completeness, though NextPage is sufficient
}

// Customer is your domain model
type Customer struct {
	ActiveSubscriptions             int
	Address                         string
	Address2                        string
	CancelledAmount                 int64
	CancelledInvoices               int
	CancelledSubscriptions          int
	City                            string
	Company                         string
	Country                         string
	Created                         time.Time
	DunningAmount                   int64
	DunningInvoices                 int
	Email                           string
	ExpiredSubscriptions            int
	FailedAmount                    int64
	FailedInvoices                  int
	FirstName                       string
	Handle                          string
	LastName                        string
	NonRenewingSubscriptions        int
	OnHoldSubscriptions             int
	PendingAdditionalCostAmount     int64
	PendingAdditionalCosts          int
	PendingAmount                   int64
	PendingCreditAmount             int64
	PendingCredits                  int
	PendingInvoices                 int
	Phone                           string
	PostalCode                      string
	RefundedAmount                  int64
	SettledAmount                   int64
	SettledInvoices                 int
	Subscriptions                   int
	Test                            bool
	TransferredAdditionalCostAmount int64
	TransferredAdditionalCosts      int
	TransferredCreditAmount         int64
	TransferredCredits              int
	TrialActiveSubscriptions        int
	TrialCancelledSubscriptions     int
}

// mapCustomer maps API response -> domain Customer
func mapCustomer(r CustomerResponse) Customer {
	return Customer{
		ActiveSubscriptions:             r.ActiveSubscriptions,
		Address:                         r.Address,
		Address2:                        r.Address2,
		CancelledAmount:                 r.CancelledAmount,
		CancelledInvoices:               r.CancelledInvoices,
		CancelledSubscriptions:          r.CancelledSubscriptions,
		City:                            r.City,
		Company:                         r.Company,
		Country:                         r.Country,
		Created:                         r.Created,
		DunningAmount:                   r.DunningAmount,
		DunningInvoices:                 r.DunningInvoices,
		Email:                           r.Email,
		ExpiredSubscriptions:            r.ExpiredSubscriptions,
		FailedAmount:                    r.FailedAmount,
		FailedInvoices:                  r.FailedInvoices,
		FirstName:                       r.FirstName,
		Handle:                          r.Handle,
		LastName:                        r.LastName,
		NonRenewingSubscriptions:        r.NonRenewingSubscriptions,
		OnHoldSubscriptions:             r.OnHoldSubscriptions,
		PendingAdditionalCostAmount:     r.PendingAdditionalCostAmount,
		PendingAdditionalCosts:          r.PendingAdditionalCosts,
		PendingAmount:                   r.PendingAmount,
		PendingCreditAmount:             r.PendingCreditAmount,
		PendingCredits:                  r.PendingCredits,
		PendingInvoices:                 r.PendingInvoices,
		Phone:                           r.Phone,
		PostalCode:                      r.PostalCode,
		RefundedAmount:                  r.RefundedAmount,
		SettledAmount:                   r.SettledAmount,
		SettledInvoices:                 r.SettledInvoices,
		Subscriptions:                   r.Subscriptions,
		Test:                            r.Test,
		TransferredAdditionalCostAmount: r.TransferredAdditionalCostAmount,
		TransferredAdditionalCosts:      r.TransferredAdditionalCosts,
		TransferredCreditAmount:         r.TransferredCreditAmount,
		TransferredCredits:              r.TransferredCredits,
		TrialActiveSubscriptions:        r.TrialActiveSubscriptions,
		TrialCancelledSubscriptions:     r.TrialCancelledSubscriptions,
	}
}

func GetCustomer(customerId string, country string) (Customer, error) {
	url := "https://api.frisbii.com/v1/customer/" + customerId

	cfg := config.LoadConfig()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Customer{}, fmt.Errorf("build request: %w", err)
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
		return Customer{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Customer{}, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return Customer{}, fmt.Errorf("customer %s not found", customerId)
	}
	if resp.StatusCode != http.StatusOK {
		return Customer{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp CustomerResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return Customer{}, fmt.Errorf("unmarshal response: %w", err)
	}

	return mapCustomer(apiResp), nil
}

func GetCustomerList(nextPage string, country string) ([]Customer, string, error) {
	cfg := config.LoadConfig()

	url := "https://api.frisbii.com/v1/list/customer?from=" + cfg.Backfill_from

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("build request: %w", err)
	}

	// Add next_page token if provided
	if nextPage != "" {
		q := req.URL.Query()
		q.Add("next_page_token", nextPage)
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

	var apiResp CustomerListResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, "", fmt.Errorf("unmarshal response: %w", err)
	}

	customers := make([]Customer, len(apiResp.Customers))
	for i, r := range apiResp.Customers {
		customers[i] = mapCustomer(r)
	}

	return customers, apiResp.NextPage, nil
}
