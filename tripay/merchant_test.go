package tripay

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMerchantPaymentChannels(t *testing.T) {
	mockResponse := MerchantPaymentChannelsResponse{
		Success: true,
		Message: "success",
		Data: []PaymentChannel{
			{
				Group: "Virtual Account",
				Code:  "BRIVA",
				Name:  "BRI Virtual Account",
				Type:  "closed",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/merchant/payment-channel" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewTripay(false, "api_key", "priv_key")
	client.BaseUrl = server.URL

	resp, err := client.GetMerchantPaymentChannels()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success to be true")
	}

	if len(resp.Data) != 1 || resp.Data[0].Code != "BRIVA" {
		t.Errorf("Expected data containing code BRIVA, got: %+v", resp.Data)
	}
}

func TestCalculateMerchantFees(t *testing.T) {
	mockResponse := MerchantFeeCalculatorResponse{
		Success: true,
		Message: "success",
		Data: []FeeCalculatorItem{
			{
				Code: "QRIS",
				Name: "QRIS",
				Fee: FeeDetail{
					Flat:    750,
					Percent: 0.7,
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/merchant/fee-calculator" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.URL.Query().Get("code") != "QRIS" || r.URL.Query().Get("amount") != "10000" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewTripay(false, "api_key", "priv_key")
	client.BaseUrl = server.URL

	resp, err := client.CalculateMerchantFees("QRIS", 10000)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success to be true")
	}

	if len(resp.Data) != 1 || resp.Data[0].Code != "QRIS" {
		t.Errorf("Expected QRIS item in calculator response, got: %+v", resp.Data)
	}
}

func TestGetMerchantTransactions(t *testing.T) {
	mockResponse := MerchantTransactionsResponse{
		Success: true,
		Message: "success",
		Data: []MerchantTransaction{
			{
				Reference:   "T12345",
				MerchantRef: "INV-01",
				Status:      "PAID",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/merchant/transactions" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.URL.Query().Get("status") != "PAID" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewTripay(false, "api_key", "priv_key")
	client.BaseUrl = server.URL

	resp, err := client.GetMerchantTransactions(MerchantTransactionsFilter{
		Status: "PAID",
	})
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success to be true")
	}

	if len(resp.Data) != 1 || resp.Data[0].Reference != "T12345" {
		t.Errorf("Expected transaction with reference T12345, got: %+v", resp.Data)
	}
}
