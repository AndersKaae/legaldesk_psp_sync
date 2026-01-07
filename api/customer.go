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

type Customer struct {
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

func GetCustomer(customerId string, country string) (map[string]any, error) {
	url := "https://api.frisbii.com/v1/customer/" + customerId

	cfg := config.LoadConfig()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
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
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return map[string]any{}, fmt.Errorf("customer %s not found", customerId)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	// Pretty-print so you can inspect it easily
	pretty, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(pretty))

	return data, nil
}
