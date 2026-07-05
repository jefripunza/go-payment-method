package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jefripunza/go-payment-method/tripay"
)

func Tripay() {
	tripay_env := getEnv("TRIPAY_ENV", "")
	tripay_api_key := getEnv("TRIPAY_API_KEY", "")
	tripay_private_key := getEnv("TRIPAY_PRIVATE_KEY", "")

	if tripay_env == "" || tripay_api_key == "" || tripay_private_key == "" {
		log.Println("Warning: TRIPAY_ENV or TRIPAY_API_KEY or TRIPAY_PRIVATE_KEY is empty. Please set them in your environment or a .env file.")
	}

	isProduction := false
	if tripay_env == "production" {
		isProduction = true
	}

	client := tripay.NewTripay(isProduction, tripay_api_key, tripay_private_key)

	// 1. Get Payment Instruction
	fmt.Println("1. Retrieving payment instructions for BRI Virtual Account...")
	resp, err := client.GetPaymentInstruction(tripay.PaymentInstructionRequest{
		Code: tripay.ChannelBriVA,
	})
	if err != nil {
		log.Printf("Error retrieving payment instruction: %v\n", err)
	} else {
		fmt.Printf("Success retrieving instructions: %t\n", resp.Success)
		fmt.Printf("Message RAW: %+v\n", resp.Data)
		for i, data := range resp.Data {
			fmt.Printf("Method %d: %s\n", i+1, data.Title)
			for j, step := range data.Step {
				fmt.Printf("  Step %d: %s\n", j+1, step)
			}
		}
	}

	// 2. Get Merchant Payment Channels
	fmt.Println("\n2. Retrieving merchant payment channels...")
	channels, err := client.GetMerchantPaymentChannels()
	if err != nil {
		log.Printf("Error retrieving merchant payment channels: %v\n", err)
	} else {
		fmt.Printf("Success: %t, Channels Count: %d\n", channels.Success, len(channels.Data))
	}

	// 3. Calculate Merchant Fees
	fmt.Println("\n3. Calculating merchant fees for BRIVA...")
	fees, err := client.CalculateMerchantFees(tripay.ChannelBriVA, 100000)
	if err != nil {
		log.Printf("Error calculating merchant fees: %v\n", err)
	} else {
		fmt.Printf("Success: %t, Fees Item Count: %d\n", fees.Success, len(fees.Data))
	}

	// 4. Get Merchant Transactions
	fmt.Println("\n4. Retrieving merchant transactions...")
	txs, err := client.GetMerchantTransactions(tripay.MerchantTransactionsFilter{
		PerPage: 5,
	})
	if err != nil {
		log.Printf("Error retrieving merchant transactions: %v\n", err)
	} else {
		fmt.Printf("Success: %t, Transactions Count: %d\n", txs.Success, len(txs.Data))
	}

	// 5. Open Payment signature generation
	merchantCode := getEnv("TRIPAY_MERCHANT_CODE", "T0001")
	merchantRef := "INV-OP-12345"
	signature := client.CreateOpenPaymentSignature(merchantCode, tripay.ChannelBcaVA, merchantRef)
	fmt.Printf("\n5. Generated Open Payment Signature: %s\n", signature)

	// 6. Create Open Payment
	fmt.Println("\n6. Creating Open Payment...")
	openPayment, err := client.CreateOpenPayment(tripay.CreateOpenPaymentRequest{
		Method:       tripay.ChannelBcaVA,
		MerchantRef:  merchantRef,
		CustomerName: "John Doe",
		Signature:    signature,
	})
	if err != nil {
		log.Printf("Error creating open payment: %v\n", err)
	} else {
		fmt.Printf("Success: %t, UUID: %s, Pay Code: %s\n", openPayment.Success, openPayment.Data.Uuid, openPayment.Data.PayCode)

		uuid := openPayment.Data.Uuid
		// 7. Get Open Payment Detail
		fmt.Println("\n7. Getting Open Payment Detail...")
		detail, err := client.GetOpenPaymentDetail(uuid)
		if err != nil {
			log.Printf("Error getting open payment detail: %v\n", err)
		} else {
			fmt.Printf("Success: %t, Status: %s, Pay Code: %s\n", detail.Success, detail.Data.Status, detail.Data.PayCode)
		}

		// 8. Get Open Payment Transactions
		fmt.Println("\n8. Getting Open Payment Transactions...")
		opTxs, err := client.GetOpenPaymentTransactions(uuid)
		if err != nil {
			log.Printf("Error getting open payment transactions: %v\n", err)
		} else {
			fmt.Printf("Success: %t, Payments Count: %d\n", opTxs.Success, len(opTxs.Data))
		}
	}

	// ============================================================================
	// Closed Payment Testing
	// ============================================================================
	merchantCode = getEnv("TRIPAY_MERCHANT_CODE", "T0001")
	closedMerchantRef := fmt.Sprintf("INV-CL-%d", 12345) // using a simple static or Unix timestamp ref
	closedAmount := 50000

	// 9. Generate Closed Payment Signature
	closedSignature := client.CreateClosedTransactionSignature(merchantCode, closedMerchantRef, closedAmount)
	fmt.Printf("\n9. Generated Closed Payment Signature: %s\n", closedSignature)

	// 10. Create Closed Transaction
	fmt.Println("\n10. Creating Closed Payment Transaction...")
	closedTx, err := client.CreateClosedTransaction(tripay.CreateClosedTransactionRequest{
		Method:        tripay.ChannelBriVA,
		MerchantRef:   closedMerchantRef,
		Amount:        closedAmount,
		CustomerName:  "John Doe",
		CustomerEmail: "johndoe@example.com",
		CustomerPhone: "081234567890",
		OrderItems: []tripay.OrderItem{
			{
				Sku:      "ITEM01",
				Name:     "Product Test",
				Price:    closedAmount,
				Quantity: 1,
			},
		},
		Signature: closedSignature,
	})
	if err != nil {
		log.Printf("Error creating closed transaction: %v\n", err)
	} else {
		fmt.Printf("Success: %t, Reference: %s, Pay Code: %s, Status: %s\n", closedTx.Success, closedTx.Data.Reference, closedTx.Data.PayCode, closedTx.Data.Status)
		ref := closedTx.Data.Reference
		checkoutUrl := closedTx.Data.CheckoutUrl
		fmt.Printf("Checkout URL: %s\n", checkoutUrl)

		// 11. Get Closed Transaction Detail (Poll until PAID)
		fmt.Println("\n11. Getting Closed Transaction Detail...")
		fmt.Printf("Please pay using the Pay Code: %s or simulate payment in the TriPay Sandbox simulator to complete this test.\n", closedTx.Data.PayCode)
		for {
			txDetail, err := client.GetClosedTransactionDetail(ref)
			if err != nil {
				log.Printf("Error getting closed transaction detail: %v\n", err)
			} else {
				// fmt.Printf("Current Status: %s, Pay Code: %s\n", txDetail.Data.Status, txDetail.Data.PayCode)
				if txDetail.Data.Status == "PAID" {
					fmt.Println("Payment successfully completed!")
					break
				}
			}
			time.Sleep(3 * time.Second)
		}

		// 12. Check Closed Transaction Status
		fmt.Println("\n12. Checking Closed Transaction Status...")
		txStatus, err := client.CheckClosedTransactionStatus(ref)
		if err != nil {
			log.Printf("Error checking closed transaction status: %v\n", err)
		} else {
			fmt.Printf("Success: %t, Reference: %s, Status: %s\n", txStatus.Success, txStatus.Data.Reference, txStatus.Data.Status)
		}
	}
}
