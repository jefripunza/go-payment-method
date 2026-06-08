package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"payment_method/xendit"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	xendit_dev_api_key := getEnv("XENDIT_DEV_API_KEY", "")
	if xendit_dev_api_key == "" {
		log.Println("Warning: XENDIT_DEV_API_KEY is empty. Please set it in your environment or a .env file.")
	}
	client := xendit.NewXendit(xendit_dev_api_key, "")

	// ============================================================================
	// First Balance
	// ============================================================================
	fmt.Println("1. Retrieving first balance from Xendit...")
	first_balance, err := client.GetBalance()
	if err != nil {
		log.Fatalf("Error retrieving balance: %v", err)
	}
	fmt.Printf("First Balance retrieved successfully: %.2f\n", first_balance)

	// ============================================================================
	// Make Transaction
	// ============================================================================

	// STEP 1: Create Invoice
	// ============================================================================
	fmt.Println("2.1. Creating invoice...")
	items := []xendit.InvoiceItem{
		{
			Name:     "Product A",
			Quantity: 1,
			Price:    100000,
			Category: "Electronics",
		},
		{
			Name:     "Product B",
			Quantity: 2,
			Price:    50000,
			Category: "Electronics",
		},
	}

	invoice, err := client.InvoiceCreate("INV-2024-001", "John Doe", "john.doe@example.com", items, []string{xendit.PaymentMethodQRIS}, 2000)
	if err != nil {
		log.Fatalf("Error creating invoice: %v", err)
	}
	fmt.Printf("Invoice created successfully! ID: %s, URL: %s\n", invoice.Id, invoice.InvoiceUrl)
	fmt.Printf("Total Amount: Rp %.2f\n", invoice.Amount)
	fmt.Printf("Status: %s\n", invoice.Status)

	// STEP 2: Check Invoice Status
	// ============================================================================
	fmt.Println("2.2. Checking invoice status...")
	for {
		time.Sleep(3 * time.Second)
		updatedInvoice, err := client.GetInvoice(invoice.Id)
		if err != nil {
			log.Fatalf("Error getting invoice: %v", err)
		}
		if updatedInvoice.Status == "PAID" || updatedInvoice.Status == "SETTLED" {
			fmt.Printf("Invoice Status changed to: %s. Invoice is paid!\n", updatedInvoice.Status)
			break
		}
	}

	// STEP 3: Get Available Disbursement Banks
	// ============================================================================
	// fmt.Println("2.3. Getting available disbursement banks...")
	// banks, err := client.GetDisbursementBanks()
	// if err != nil {
	// 	log.Fatalf("Error getting banks: %v", err)
	// }
	// fmt.Println("Available Banks:")
	// for _, bank := range banks {
	// 	fmt.Printf("  - %s (%s)\n", bank.Name, bank.Code)
	// }

	// if len(banks) == 0 {
	// 	log.Println("No banks available for disbursement. Skipping disbursement test.")
	// 	return
	// }

	// STEP 4: Make Disbursement
	// ============================================================================
	// fmt.Println("2.4. Creating disbursement...")
	// firstBank := banks[0]
	// disbursement, err := client.DisbursementCreate(
	// 	firstBank.Code,
	// 	"John Doe",
	// 	"123456789", // Dummy account number
	// 	100000,
	// 	"Withdrawal test",
	// )
	// if err != nil {
	// 	log.Fatalf("Error creating disbursement: %v", err)
	// }
	// fmt.Printf("Disbursement created successfully! ID: %s\n", disbursement.Id)
	// fmt.Printf("Bank: %s\n", disbursement.BankCode)
	// fmt.Printf("Amount: %.2f\n", disbursement.Amount)
	// fmt.Printf("Status: %s\n", disbursement.Status)

	// STEP 5: Get Second Balance (should be lower after disbursement)
	// ============================================================================
	fmt.Println("2.5. Checking second balance...")
	second_balance, err := client.GetBalance()
	if err != nil {
		log.Fatalf("Error retrieving second balance: %v", err)
	}
	fmt.Printf("Second Balance retrieved successfully: %.2f\n", second_balance)
	fmt.Printf("Balance difference: %.2f\n", first_balance-second_balance)
	if second_balance < first_balance {
		fmt.Println("Balance decreased after disbursement - success!")
	} else {
		fmt.Println("Warning: Balance did not decrease. Check disbursement status.")
	}
}
