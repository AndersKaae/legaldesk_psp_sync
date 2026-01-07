package database

import (
	"database/sql"
	"fmt"
	"github.com/AndersKaae/legaldesk_psp_sync/config"
	_ "modernc.org/sqlite"
	"time"
)

type InvoiceStates map[string]*time.Time

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

var db *sql.DB

func InitDB(cfg config.Config) error {
	var err error
	db, err = sql.Open("sqlite", cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}
	return nil
}

func createTables() error {
	customerTable := `
	CREATE TABLE IF NOT EXISTS customers (
		handle TEXT PRIMARY KEY,
		active_subscriptions INTEGER,
		address TEXT,
		address2 TEXT,
		cancelled_amount INTEGER,
		cancelled_invoices INTEGER,
		cancelled_subscriptions INTEGER,
		city TEXT,
		company TEXT,
		country TEXT,
		created DATETIME,
		dunning_amount INTEGER,
		dunning_invoices INTEGER,
		email TEXT,
		expired_subscriptions INTEGER,
		failed_amount INTEGER,
		failed_invoices INTEGER,
		first_name TEXT,
		last_name TEXT,
		non_renewing_subscriptions INTEGER,
		on_hold_subscriptions INTEGER,
		pending_additional_cost_amount INTEGER,
		pending_additional_costs INTEGER,
		pending_amount INTEGER,
		pending_credit_amount INTEGER,
		pending_credits INTEGER,
		pending_invoices INTEGER,
		phone TEXT,
		postal_code TEXT,
		refunded_amount INTEGER,
		settled_amount INTEGER,
		settled_invoices INTEGER,
		subscriptions INTEGER,
		test BOOLEAN,
		transferred_additional_cost_amount INTEGER,
		transferred_additional_costs INTEGER,
		transferred_credit_amount INTEGER,
		transferred_credits INTEGER,
		trial_active_subscriptions INTEGER,
		trial_cancelled_subscriptions INTEGER
	);`

	invoiceTable := `
	CREATE TABLE IF NOT EXISTS invoices (
		id TEXT PRIMARY KEY,
		handle TEXT,
		customer TEXT,
		currency TEXT,
		created DATETIME,
		discount_amount INTEGER,
		org_amount INTEGER,
		amount_vat INTEGER,
		amount_ex_vat INTEGER,
		refunded_amount INTEGER,
		authorized_amount INTEGER,
		country TEXT,
		states TEXT
	);`

	if _, err := db.Exec(customerTable); err != nil {
		return err
	}

	if _, err := db.Exec(invoiceTable); err != nil {
		return err
	}

	return nil
}

