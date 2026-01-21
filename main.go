package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AndersKaae/legaldesk_psp_sync/api"
	"github.com/AndersKaae/legaldesk_psp_sync/config"
	"github.com/AndersKaae/legaldesk_psp_sync/database"
	"github.com/AndersKaae/legaldesk_psp_sync/handlers"
)

var backfillFlag = flag.Bool("backfill", false, "Run backfill process to populate old data")

func setupLogging() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Create a multi-writer to write to both stdout and the log file
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

// backfillCustomersForCountry fetches and saves all customers for a given country
func backfillCustomersForCountry(country string) int {
	log.Printf("Starting customer backfill for country: %s", country)
	customerCount := 0
	nextPage := ""
	for {
		customers, newNextPage, err := api.GetCustomerList(nextPage, country)
		if err != nil {
			log.Printf("Error fetching customer list for %s (nextPage: %s): %v", country, nextPage, err)
			break // Stop backfill for this country on error
		}

		for _, apiCustomer := range customers {
			dbCustomer := database.Customer{
				ActiveSubscriptions:             apiCustomer.ActiveSubscriptions,
				Address:                         apiCustomer.Address,
				Address2:                        apiCustomer.Address2,
				CancelledAmount:                 apiCustomer.CancelledAmount,
				CancelledInvoices:               apiCustomer.CancelledInvoices,
				CancelledSubscriptions:          apiCustomer.CancelledSubscriptions,
				City:                            apiCustomer.City,
				Company:                         apiCustomer.Company,
				Country:                         apiCustomer.Country,
				Created:                         apiCustomer.Created,
				DunningAmount:                   apiCustomer.DunningAmount,
				DunningInvoices:                 apiCustomer.DunningInvoices,
				Email:                           apiCustomer.Email,
				ExpiredSubscriptions:            apiCustomer.ExpiredSubscriptions,
				FailedAmount:                    apiCustomer.FailedAmount,
				FailedInvoices:                  apiCustomer.FailedInvoices,
				FirstName:                       apiCustomer.FirstName,
				Handle:                          apiCustomer.Handle,
				LastName:                        apiCustomer.LastName,
				NonRenewingSubscriptions:        apiCustomer.NonRenewingSubscriptions,
				OnHoldSubscriptions:             apiCustomer.OnHoldSubscriptions,
				PendingAdditionalCostAmount:     apiCustomer.PendingAdditionalCostAmount,
				PendingAdditionalCosts:          apiCustomer.PendingAdditionalCosts,
				PendingAmount:                   apiCustomer.PendingAmount,
				PendingCreditAmount:             apiCustomer.PendingCreditAmount,
				PendingCredits:                  apiCustomer.PendingCredits,
				PendingInvoices:                 apiCustomer.PendingInvoices,
				Phone:                           apiCustomer.Phone,
				PostalCode:                      apiCustomer.PostalCode,
				RefundedAmount:                  apiCustomer.RefundedAmount,
				SettledAmount:                   apiCustomer.SettledAmount,
				SettledInvoices:                 apiCustomer.SettledInvoices,
				Subscriptions:                   apiCustomer.Subscriptions,
				Test:                            apiCustomer.Test,
				TransferredAdditionalCostAmount: apiCustomer.TransferredAdditionalCostAmount,
				TransferredAdditionalCosts:      apiCustomer.TransferredAdditionalCosts,
				TransferredCreditAmount:         apiCustomer.TransferredCreditAmount,
				TransferredCredits:              apiCustomer.TransferredCredits,
				TrialActiveSubscriptions:        apiCustomer.TrialActiveSubscriptions,
				TrialCancelledSubscriptions:     apiCustomer.TrialCancelledSubscriptions,
			}
			if err := database.CreateOrUpdateCustomer(&dbCustomer); err != nil {
				log.Printf("Error saving backfilled customer %s to DB for %s: %v", apiCustomer.Handle, country, err)
			} else {
				log.Printf("Backfilled customer %s for %s", apiCustomer.Handle, country)
				customerCount++
			}
		}

		if newNextPage == "" {
			break // No more pages
		}
		nextPage = newNextPage
		time.Sleep(100 * time.Millisecond) // Be nice to the API
	}
	log.Printf("Finished customer backfill for country: %s", country)
	return customerCount
}

// backfillInvoicesForCountry fetches and saves all invoices for a given country
func backfillInvoicesForCountry(country string) int {
	log.Printf("Starting invoice backfill for country: %s", country)
	invoiceCount := 0
	nextPage := ""
	for {
		invoices, newNextPage, err := api.GetInvoiceList(nextPage, country)
		if err != nil {
			log.Printf("Error fetching invoice list for %s (nextPage: %s): %v", country, nextPage, err)
			break // Stop backfill for this country on error
		}

		for _, apiInvoice := range invoices {
			dbInvoice := database.Invoice{
				ID:               apiInvoice.ID,
				Handle:           apiInvoice.Handle,
				Customer:         apiInvoice.Customer,
				Currency:         apiInvoice.Currency,
				Created:          apiInvoice.Created,
				DiscountAmount:   apiInvoice.DiscountAmount,
				OrgAmount:        apiInvoice.OrgAmount,
				AmountVAT:        apiInvoice.AmountVAT,
				AmountExVAT:      apiInvoice.AmountExVAT,
				RefundedAmount:   apiInvoice.RefundedAmount,
				AuthorizedAmount: apiInvoice.AuthorizedAmount,
				Country:          apiInvoice.Country,
				Plan:             apiInvoice.Plan,
				States:           database.InvoiceStates(apiInvoice.States),
			}
			jsonInvoice, _ := json.Marshal(apiInvoice)
			if err := database.CreateOrUpdateInvoice(&dbInvoice); err != nil {
				log.Printf("API Invoice for %s: %s", apiInvoice.ID, string(jsonInvoice))
				log.Printf("Error saving backfilled invoice %s to DB for %s: %v", apiInvoice.ID, country, err)
			} else {
				log.Printf("(%s) Backfilled invoice %s for %s", apiInvoice.Created, apiInvoice.Handle, country)
				invoiceCount++
			}
		}

		if newNextPage == "" {
			break // No more pages
		}
		nextPage = newNextPage
		time.Sleep(100 * time.Millisecond) // Be nice to the API
	}
	log.Printf("Finished invoice backfill for country: %s", country)
	return invoiceCount
}

func runBackfill() {
	log.Println("Starting full backfill process...")
	totalCustomersBackfilled := 0
	totalInvoicesBackfilled := 0
	countries := []string{"DK", "SE", "NO"}
	for _, country := range countries {
		totalCustomersBackfilled += backfillCustomersForCountry(country)
		totalInvoicesBackfilled += backfillInvoicesForCountry(country)
	}
	log.Printf("Full backfill process finished. Total Customers Backfilled: %d, Total Invoices Backfilled: %d", totalCustomersBackfilled, totalInvoicesBackfilled)
}

func main() {
	flag.Parse() // Parse command-line flags

	setupLogging()

	cfg := config.LoadConfig()
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if *backfillFlag {
		go runBackfill() // Run backfill in a goroutine
	}

	http.HandleFunc("/invoices/", handlers.RequireBasicAuth(handlers.Invoices("all"), cfg.BasicAuthUser, cfg.BasicAuthPass))

	http.HandleFunc("/virtual-office-invoices/", handlers.RequireBasicAuth(handlers.Invoices("virtualOffice"), cfg.BasicAuthUser, cfg.BasicAuthPass))
	http.HandleFunc("/webhook/denmark", handlers.Webhook("DK"))
	http.HandleFunc("/webhook/sweden", handlers.Webhook("SE"))
	http.HandleFunc("/webhook/norway", handlers.Webhook("NO"))

	addr := ":6969"
	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
