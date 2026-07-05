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

	fmt.Println("1. Retrieving payment instructions for BRI Virtual Account...")
	resp, err := client.GetPaymentInstruction(tripay.PaymentInstructionRequest{
		Code: tripay.ChannelBriVA,
	})
	if err != nil {
		log.Fatalf("Error retrieving payment instruction: %v", err)
	}

	fmt.Printf("Success retrieving instructions: %t\n", resp.Success)
	fmt.Printf("Message RAW: %+v\n", resp.Data)
	for i, data := range resp.Data {
		fmt.Printf("Method %d: %s\n", i+1, data.Title)
		for j, step := range data.Step {
			fmt.Printf("  Step %d: %s\n", j+1, step)
		}
	}
}
