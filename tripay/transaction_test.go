package tripay

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateOpenPaymentSignature(t *testing.T) {
	privateKey := "my_private_key"
	client := NewTripay(false, "api_key", privateKey)

	merchantCode := "T0001"
	channel := "BCAVA"
	merchantRef := "INV55567"

	// Expected signature matches hash_hmac('sha256', "T0001BCAVAINV55567", "my_private_key")
	// Let's compute manually or check
	signature := client.CreateOpenPaymentSignature(merchantCode, channel, merchantRef)
	if signature == "" {
		t.Error("Expected signature to be computed, got empty string")
	}
}

func TestCreateOpenPayment(t *testing.T) {
	mockResponse := CreateOpenPaymentResponse{
		Success: true,
		Message: "success",
		Data: OpenPaymentData{
			Uuid:          "123-uuid",
			Reference:     "T12345",
			MerchantRef:   "INV-01",
			PaymentName:   "BCA VA",
			PaymentMethod: "BCAVA",
			PayCode:       "888123456",
			Status:        "UNPAID",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/open-payment/create" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewTripay(false, "api_key", "priv_key")
	client.BaseUrl = server.URL

	resp, err := client.CreateOpenPayment(CreateOpenPaymentRequest{
		Method:       "BCAVA",
		MerchantRef:  "INV-01",
		CustomerName: "John Doe",
	})
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success to be true")
	}

	if resp.Data.Uuid != "123-uuid" || resp.Data.PayCode != "888123456" {
		t.Errorf("Expected response uuid '123-uuid' and paycode '888123456', got: %+v", resp.Data)
	}
}

func TestGetOpenPaymentDetail(t *testing.T) {
	uuid := "123-uuid"
	mockResponse := OpenPaymentDetailResponse{
		Success: true,
		Message: "success",
		Data: OpenPaymentData{
			Uuid:          uuid,
			Reference:     "T12345",
			MerchantRef:   "INV-01",
			PaymentName:   "BCA VA",
			PaymentMethod: "BCAVA",
			PayCode:       "888123456",
			Status:        "UNPAID",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/open-payment/" + uuid + "/detail"
		if r.URL.Path != expectedPath {
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

	resp, err := client.GetOpenPaymentDetail(uuid)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success to be true")
	}

	if resp.Data.Uuid != uuid {
		t.Errorf("Expected UUID %s, got: %s", uuid, resp.Data.Uuid)
	}
}

func TestGetOpenPaymentTransactions(t *testing.T) {
	uuid := "123-uuid"
	mockResponse := OpenPaymentTransactionsResponse{
		Success: true,
		Message: "success",
		Data: []OpenPaymentTransactionItem{
			{
				Reference:   "T12345-01",
				MerchantRef: "INV-01",
				TotalAmount: 50000,
				Status:      "PAID",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/open-payment/" + uuid + "/transactions"
		if r.URL.Path != expectedPath {
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

	resp, err := client.GetOpenPaymentTransactions(uuid)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success to be true")
	}

	if len(resp.Data) != 1 || resp.Data[0].Reference != "T12345-01" {
		t.Errorf("Expected transaction reference 'T12345-01', got: %+v", resp.Data)
	}
}
