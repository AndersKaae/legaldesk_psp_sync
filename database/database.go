package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/AndersKaae/legaldesk_psp_sync/config"
	"github.com/go-sql-driver/mysql"
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
	ID               string `json:"id"`
	Handle           string `json:"handle"`
	Customer         string `json:"customer"`
	Currency         string `json:"currency"`
	Created          time.Time `json:"created"`
	DiscountAmount   int64 `json:"discount_amount"`
	OrgAmount        int64 `json:"org_amount"`
	AmountVAT        int64 `json:"amount_vat"`
	AmountExVAT      int64 `json:"amount_ex_vat"`
	RefundedAmount   int64 `json:"refunded_amount"`
	AuthorizedAmount int64 `json:"authorized_amount"`
	Country          string `json:"country"`
	Plan             string `json:"plan"`
	States           InvoiceStates `json:"states"`
}

var db *sql.DB

func InitDB(cfg config.Config) error {
	if err := createDatabaseIfNotExists(cfg.DatabaseDSN); err != nil {
		return err
	}
	var err error
	db, err = sql.Open("mysql", cfg.DatabaseDSN)
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

func createDatabaseIfNotExists(dsn string) error {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse DSN: %w", err)
	}

	dbName := cfg.DBName
	cfg.DBName = ""

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("failed to open database for creation check: %w", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	return nil
}

func createTables() error {
	customerTable := `
	CREATE TABLE IF NOT EXISTS customers (
		handle VARCHAR(255) PRIMARY KEY,
		active_subscriptions INTEGER,
		address TEXT,
		address2 TEXT,
		cancelled_amount INTEGER,
		cancelled_invoices INTEGER,
		cancelled_subscriptions INTEGER,
		city VARCHAR(255),
		company VARCHAR(255),
		country VARCHAR(255),
		created DATETIME,
		dunning_amount INTEGER,
		dunning_invoices INTEGER,
		email VARCHAR(255),
		expired_subscriptions INTEGER,
		failed_amount INTEGER,
		failed_invoices INTEGER,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		non_renewing_subscriptions INTEGER,
		on_hold_subscriptions INTEGER,
		pending_additional_cost_amount INTEGER,
		pending_additional_costs INTEGER,
		pending_amount INTEGER,
		pending_credit_amount INTEGER,
		pending_credits INTEGER,
		pending_invoices INTEGER,
		phone VARCHAR(255),
		postal_code VARCHAR(255),
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
		id VARCHAR(255) PRIMARY KEY,
		handle VARCHAR(255),
		customer VARCHAR(255),
		currency VARCHAR(255),
		created DATETIME,
		discount_amount INTEGER,
		org_amount INTEGER,
		amount_vat INTEGER,
		amount_ex_vat INTEGER,
		refunded_amount INTEGER,
		authorized_amount INTEGER,
		country VARCHAR(255),
		plan VARCHAR(255),
		states JSON
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
	ON DUPLICATE KEY UPDATE
		active_subscriptions = VALUES(active_subscriptions),
		address = VALUES(address),
		address2 = VALUES(address2),
		cancelled_amount = VALUES(cancelled_amount),
		cancelled_invoices = VALUES(cancelled_invoices),
		cancelled_subscriptions = VALUES(cancelled_subscriptions),
		city = VALUES(city),
		company = VALUES(company),
		country = VALUES(country),
		created = VALUES(created),
		dunning_amount = VALUES(dunning_amount),
		dunning_invoices = VALUES(dunning_invoices),
		email = VALUES(email),
		expired_subscriptions = VALUES(expired_subscriptions),
		failed_amount = VALUES(failed_amount),
		failed_invoices = VALUES(failed_invoices),
		first_name = VALUES(first_name),
		last_name = VALUES(last_name),
		non_renewing_subscriptions = VALUES(non_renewing_subscriptions),
		on_hold_subscriptions = VALUES(on_hold_subscriptions),
		pending_additional_cost_amount = VALUES(pending_additional_cost_amount),
		pending_additional_costs = VALUES(pending_additional_costs),
		pending_amount = VALUES(pending_amount),
		pending_credit_amount = VALUES(pending_credit_amount),
		pending_credits = VALUES(pending_credits),
		pending_invoices = VALUES(pending_invoices),
		phone = VALUES(phone),
		postal_code = VALUES(postal_code),
		refunded_amount = VALUES(refunded_amount),
		settled_amount = VALUES(settled_amount),
		settled_invoices = VALUES(settled_invoices),
		subscriptions = VALUES(subscriptions),
		test = VALUES(test),
		transferred_additional_cost_amount = VALUES(transferred_additional_cost_amount),
		transferred_additional_costs = VALUES(transferred_additional_costs),
		transferred_credit_amount = VALUES(transferred_credit_amount),
		transferred_credits = VALUES(transferred_credits),
		trial_active_subscriptions = VALUES(trial_active_subscriptions),
		trial_cancelled_subscriptions = VALUES(trial_cancelled_subscriptions);
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
		amount_vat, amount_ex_vat, refunded_amount, authorized_amount, country, plan, states
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
		handle = VALUES(handle),
		customer = VALUES(customer),
		currency = VALUES(currency),
		created = VALUES(created),
		discount_amount = VALUES(discount_amount),
		org_amount = VALUES(org_amount),
		amount_vat = VALUES(amount_vat),
		amount_ex_vat = VALUES(amount_ex_vat),
		refunded_amount = VALUES(refunded_amount),
		authorized_amount = VALUES(authorized_amount),
		country = VALUES(country),
		plan = VALUES(plan),
		states = VALUES(states);
	`
	_, err = db.Exec(query,
		invoice.ID, invoice.Handle, invoice.Customer, invoice.Currency, invoice.Created, invoice.DiscountAmount, invoice.OrgAmount,
		invoice.AmountVAT, invoice.AmountExVAT, invoice.RefundedAmount, invoice.AuthorizedAmount, invoice.Country, invoice.Plan, statesJSON,
	)
	return err
}

func GetInvoicesByDateRange(from, to time.Time) ([]Invoice, error) {
	query := `
	SELECT id, handle, customer, currency, created, discount_amount, org_amount,
	amount_vat, amount_ex_vat, refunded_amount, authorized_amount, country, plan, states
	FROM invoices WHERE created >= ? AND created < ? ORDER BY created DESC
	`
	rows, err := db.Query(query, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to query invoices by date range: %w", err)
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var invoice Invoice
		var statesJSON []byte
		if err := rows.Scan(
			&invoice.ID, &invoice.Handle, &invoice.Customer, &invoice.Currency, &invoice.Created,
			&invoice.DiscountAmount, &invoice.OrgAmount, &invoice.AmountVAT, &invoice.AmountExVAT,
			&invoice.RefundedAmount, &invoice.AuthorizedAmount, &invoice.Country, &invoice.Plan, &statesJSON,
		); err != nil {
			return nil, fmt.Errorf("failed to scan invoice row: %w", err)
		}

		if err := json.Unmarshal(statesJSON, &invoice.States); err != nil {
			return nil, fmt.Errorf("failed to unmarshal invoice states: %w", err)
		}
		invoices = append(invoices, invoice)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return invoices, nil
}

func GetVirtualOfficeInvoicesByDateRange(from, to time.Time) ([]Invoice, error) {
	query := `
	SELECT id, handle, customer, currency, created, discount_amount, org_amount,
	amount_vat, amount_ex_vat, refunded_amount, authorized_amount, country, plan, states
	FROM invoices WHERE created >= ? AND created < ? AND handle LIKE "inv%" ORDER BY created DESC
	`
	rows, err := db.Query(query, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to query invoices by date range: %w", err)
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var invoice Invoice
		var statesJSON []byte
		if err := rows.Scan(
			&invoice.ID, &invoice.Handle, &invoice.Customer, &invoice.Currency, &invoice.Created,
			&invoice.DiscountAmount, &invoice.OrgAmount, &invoice.AmountVAT, &invoice.AmountExVAT,
			&invoice.RefundedAmount, &invoice.AuthorizedAmount, &invoice.Country, &invoice.Plan, &statesJSON,
		); err != nil {
			return nil, fmt.Errorf("failed to scan invoice row: %w", err)
		}

		if err := json.Unmarshal(statesJSON, &invoice.States); err != nil {
			return nil, fmt.Errorf("failed to unmarshal invoice states: %w", err)
		}
		invoices = append(invoices, invoice)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return invoices, nil
}
