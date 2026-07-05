package main

import (
	"fmt"
	"log"

	"github.com/jefripunza/go-payment-method/tripay"
)

func Tripay() {
	tripay_api_key := getEnv("TRIPAY_API_KEY", "")
	tripay_private_key := getEnv("TRIPAY_PRIVATE_KEY", "")

	if tripay_api_key == "" || tripay_private_key == "" {
		log.Println("Warning: TRIPAY_API_KEY or TRIPAY_PRIVATE_KEY is empty. Please set them in your environment or a .env file.")
	}

	client := tripay.NewTripay(false, tripay_api_key, tripay_private_key)

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
}
