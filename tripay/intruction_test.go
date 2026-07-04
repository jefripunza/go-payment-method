package tripay

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPaymentInstruction(t *testing.T) {


	// We can mock the HTTP request using a local test server
	mockResponse := PaymentInstructionResponse{
		Success: true,
		Message: "success",
		Data: []InstructionStep{
			{
				Title: "ATM BNI",
				Step: []string{
					"Masukkan kartu ATM BNI dan PIN Anda.",
					"Pilih menu 'Menu Lainnya' > 'Transfer' > 'Virtual Account Billing'.",
				},
			},
		},
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer mock_api_key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.URL.Query().Get("code") != "BNIVA" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer testServer.Close()

	client := NewTripay(false, "mock_api_key", "mock_private_key")
	client.BaseUrl = testServer.URL // point client to test server

	req := PaymentInstructionRequest{
		Code: "BNIVA",
	}

	resp, err := client.GetPaymentInstruction(req)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success to be true, got false")
	}

	if len(resp.Data) != 1 {
		t.Fatalf("Expected 1 instruction step, got %d", len(resp.Data))
	}

	if resp.Data[0].Title != "ATM BNI" {
		t.Errorf("Expected Title 'ATM BNI', got: %s", resp.Data[0].Title)
	}
}
