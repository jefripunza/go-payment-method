package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jefripunza/go-payment-method/midtrans"
)

// Midtrans mendemonstrasikan alur Midtrans: Snap (checkout page), Core API charge
// (QRIS + VA), verifikasi webhook, dan cek status transaksi.
func Midtrans() {
	serverKey := getEnv("MIDTRANS_SERVER_KEY", "")
	clientKey := getEnv("MIDTRANS_CLIENT_KEY", "")
	isProduction := getEnv("MIDTRANS_IS_PRODUCTION", "false") == "true"
	if serverKey == "" {
		log.Println("Warning: MIDTRANS_SERVER_KEY is empty. Please set it in your environment or a .env file.")
		return
	}

	m := midtrans.NewMidtrans(isProduction, serverKey, clientKey)
	ts := fmt.Sprintf("%d", time.Now().Unix())

	// ============================================================================
	// STEP 1: Snap — buat transaction token (frontend menampilkan popup/redirect)
	// ============================================================================
	fmt.Println("1. Creating Snap transaction token...")
	snap, err := m.CreateSnapTransaction(midtrans.SnapRequest{
		TransactionDetails: midtrans.SnapTransactionDetails{
			OrderID:     "SNAP-" + ts,
			GrossAmount: 10000,
		},
		CustomerDetails: &midtrans.SnapCustomerDetails{
			FirstName: "Budi",
			LastName:  "Susanto",
			Email:     "budi.susanto@example.com",
			Phone:     "08123456789",
		},
		ItemDetails: []midtrans.SnapItemDetail{
			{ID: "a1", Price: 10000, Quantity: 1, Name: "Produk A"},
		},
	})
	if err != nil {
		log.Fatalf("Error creating Snap transaction: %v", err)
	}
	fmt.Printf("Snap token: %s\n", snap.Token)
	fmt.Printf("Redirect URL: %s\n", snap.RedirectURL)

	// ============================================================================
	// STEP 2: Core API — charge QRIS (bukan Snap)
	// ============================================================================
	fmt.Println("2. Charging via Core API (QRIS)...")
	orderID := "QRIS-" + ts
	charge, err := m.Charge(midtrans.CoreChargeRequest{
		PaymentType: "qris",
		TransactionDetails: midtrans.CoreTransactionDetails{
			OrderID:     orderID,
			GrossAmount: 15000,
		},
		AdditionalFields: map[string]interface{}{
			"qris": map[string]interface{}{"acquirer": "gopay"},
		},
	})
	if err != nil {
		log.Fatalf("Error charging QRIS: %v", err)
	}
	fmt.Printf("Transaction ID: %v, Status: %v\n", charge["transaction_id"], charge["transaction_status"])
	if qr, ok := charge["qr_string"].(string); ok {
		fmt.Printf("QR String: %s\n", qr)
	}

	// ============================================================================
	// STEP 3: Core API — charge bank transfer (VA BCA)
	// ============================================================================
	fmt.Println("3. Charging via Core API (BCA VA)...")
	vaCharge, err := m.Charge(midtrans.CoreChargeRequest{
		PaymentType: "bank_transfer",
		TransactionDetails: midtrans.CoreTransactionDetails{
			OrderID:     "VA-" + ts,
			GrossAmount: 50000,
		},
		AdditionalFields: map[string]interface{}{
			"bank_transfer": map[string]interface{}{"bank": "bca"},
		},
	})
	if err != nil {
		log.Fatalf("Error charging VA: %v", err)
	}
	if vaNumbers, ok := vaCharge["va_numbers"].([]interface{}); ok && len(vaNumbers) > 0 {
		fmt.Printf("VA Numbers: %v\n", vaNumbers[0])
	}

	// ============================================================================
	// STEP 4: Cek status transaksi (untuk verifikasi webhook / failover)
	// ============================================================================
	fmt.Println("4. Checking transaction status...")
	status, err := m.GetTransactionStatus(orderID)
	if err != nil {
		log.Fatalf("Error getting status: %v", err)
	}
	fmt.Printf("Order: %s, Transaction status: %s\n", status.OrderID, status.TransactionStatus)

	// ============================================================================
	// STEP 5: Verifikasi signature webhook (SHA512)
	// ============================================================================
	fmt.Println("5. Verifying webhook signature...")
	ok := m.VerifyNotificationSignature(status.SignatureKey, status.OrderID, status.StatusCode, status.GrossAmount)
	fmt.Printf("Signature valid: %v\n", ok)

	// ============================================================================
	// STEP 6: IRIS — cek saldo disbursement
	// ============================================================================
	fmt.Println("6. IRIS balance (disbursement)...")
	balance, err := m.IRISGetBalance()
	if err != nil {
		log.Printf("IRIS balance error (mungkin IRIS belum diaktifkan): %v", err)
	} else {
		fmt.Printf("IRIS balance: %v\n", balance)
	}
}
