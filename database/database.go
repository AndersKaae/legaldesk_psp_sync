package database

import (
	"database/sql"
	"encoding/json"
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

func CreateOrUpdateCustomer(customer *Customer) error {
	query := `
	INSERT INTO customers (
		handle, active_subscriptions, address, address2, cancelled_amount, cancelled_invoices,
		cancelled_subscriptions, city, company, country, created, dunning_amount, dunning_invoices,
		email, expired_subscriptions, failed_amount, failed_invoices, first_name, last_name,
		non_renewing_subscriptions, on_hold_subscriptions, pending_additional_cost_amount,
		pending_additional_costs, pending_amount, pending_credit_amount, pending_credits,
		pending_invoices, phone, postal_code, refunded_amount, settled_amount, settled_invoices,
		subscriptions, test, transferred_additional_cost_amount, transferred_additional_costs,
		transferred_credit_amount, transferred_credits, trial_active_subscriptions,
		trial_cancelled_subscriptions
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(handle) DO UPDATE SET
		active_subscriptions = excluded.active_subscriptions,
		address = excluded.address,
		address2 = excluded.address2,
		cancelled_amount = excluded.cancelled_amount,
		cancelled_invoices = excluded.cancelled_invoices,
		cancelled_subscriptions = excluded.cancelled_subscriptions,
		city = excluded.city,
		company = excluded.company,
		country = excluded.country,
		created = excluded.created,
		dunning_amount = excluded.dunning_amount,
		dunning_invoices = excluded.dunning_invoices,
		email = excluded.email,
		expired_subscriptions = excluded.expired_subscriptions,
		failed_amount = excluded.failed_amount,
		failed_invoices = excluded.failed_invoices,
		first_name = excluded.first_name,
		last_name = excluded.last_name,
		non_renewing_subscriptions = excluded.non_renewing_subscriptions,
		on_hold_subscriptions = excluded.on_hold_subscriptions,
		pending_additional_cost_amount = excluded.pending_additional_cost_amount,
		pending_additional_costs = excluded.pending_additional_costs,
		pending_amount = excluded.pending_amount,
		pending_credit_amount = excluded.pending_credit_amount,
		pending_credits = excluded.pending_credits,
		pending_invoices = excluded.pending_invoices,
		phone = excluded.phone,
		postal_code = excluded.postal_code,
		refunded_amount = excluded.refunded_amount,
		settled_amount = excluded.settled_amount,
		settled_invoices = excluded.settled_invoices,
		subscriptions = excluded.subscriptions,
		test = excluded.test,
		transferred_additional_cost_amount = excluded.transferred_additional_cost_amount,
		transferred_additional_costs = excluded.transferred_additional_costs,
		transferred_credit_amount = excluded.transferred_credit_amount,
		transferred_credits = excluded.transferred_credits,
		trial_active_subscriptions = excluded.trial_active_subscriptions,
		trial_cancelled_subscriptions = excluded.trial_cancelled_subscriptions;
	`
	_, err := db.Exec(query,
		customer.Handle, customer.ActiveSubscriptions, customer.Address, customer.Address2, customer.CancelledAmount, customer.CancelledInvoices,
		customer.CancelledSubscriptions, customer.City, customer.Company, customer.Country, customer.Created, customer.DunningAmount, customer.DunningInvoices,
		customer.Email, customer.ExpiredSubscriptions, customer.FailedAmount, customer.FailedInvoices, customer.FirstName, customer.LastName,
		customer.NonRenewingSubscriptions, customer.OnHoldSubscriptions, customer.PendingAdditionalCostAmount,
		customer.PendingAdditionalCosts, customer.PendingAmount, customer.PendingCreditAmount, customer.PendingCredits,
		customer.PendingInvoices, customer.Phone, customer.PostalCode, customer.RefundedAmount, customer.SettledAmount, customer.SettledInvoices,
		customer.Subscriptions, customer.Test, customer.TransferredAdditionalCostAmount, customer.TransferredAdditionalCosts,
		customer.TransferredCreditAmount, customer.TransferredCredits, customer.TrialActiveSubscriptions,
		customer.TrialCancelledSubscriptions,
	)
	return err
}

func CreateOrUpdateInvoice(invoice *Invoice) error {
	statesJSON, err := json.Marshal(invoice.States)
	if err != nil {
		return fmt.Errorf("failed to marshal invoice states: %w", err)
	}

	query := `
	INSERT INTO invoices (
		id, handle, customer, currency, created, discount_amount, org_amount,
		amount_vat, amount_ex_vat, refunded_amount, authorized_amount, country, states
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		handle = excluded.handle,
		customer = excluded.customer,
		currency = excluded.currency,
		created = excluded.created,
		discount_amount = excluded.discount_amount,
		org_amount = excluded.org_amount,
		amount_vat = excluded.amount_vat,
		amount_ex_vat = excluded.amount_ex_vat,
		refunded_amount = excluded.refunded_amount,
		authorized_amount = excluded.authorized_amount,
		country = excluded.country,
		states = excluded.states;
	`
	_, err = db.Exec(query,
		invoice.ID, invoice.Handle, invoice.Customer, invoice.Currency, invoice.Created, invoice.DiscountAmount, invoice.OrgAmount,
		invoice.AmountVAT, invoice.AmountExVAT, invoice.RefundedAmount, invoice.AuthorizedAmount, invoice.Country, statesJSON,
	)
	return err
}
