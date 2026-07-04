package tripay

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCallbackHttp_ValidSignature(t *testing.T) {
	privateKey := "my_super_secret_private_key"
	client := NewTripay(false, "my_api_key", privateKey)

	payload := `{"reference":"T1234567890","merchant_ref":"INV-001","payment_method":"BRIVA","payment_method_code":"BRIVA","total_amount":100000,"fee_merchant":4250,"amount_received":95750,"status":"PAID","paid_at":"2024-01-15 10:30:00"}`

	// Generate signature
	h := hmac.New(sha256.New, []byte(privateKey))
	h.Write([]byte(payload))
	signature := hex.EncodeToString(h.Sum(nil))

	req := httptest.NewRequest(http.MethodPost, "/callback", bytes.NewBufferString(payload))
	req.Header.Set("X-Callback-Signature", signature)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	resp, err := client.CallbackHttp(w, req)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response not to be nil")
	}

	if resp.Reference != "T1234567890" {
		t.Errorf("Expected reference 'T1234567890', got: %s", resp.Reference)
	}

	if resp.MerchantRef != "INV-001" {
		t.Errorf("Expected merchant_ref 'INV-001', got: %s", resp.MerchantRef)
	}

	if resp.TotalAmount != 100000 {
		t.Errorf("Expected total_amount 100000, got: %f", resp.TotalAmount)
	}

	if resp.Status != "PAID" {
		t.Errorf("Expected status 'PAID', got: %s", resp.Status)
	}
}

func TestCallbackHttp_InvalidSignature(t *testing.T) {
	privateKey := "my_super_secret_private_key"
	client := NewTripay(false, "my_api_key", privateKey)

	payload := `{"reference":"T1234567890","merchant_ref":"INV-001","status":"PAID"}`

	req := httptest.NewRequest(http.MethodPost, "/callback", bytes.NewBufferString(payload))
	req.Header.Set("X-Callback-Signature", "wrong_signature")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	resp, err := client.CallbackHttp(w, req)
	if err == nil {
		t.Fatal("Expected error due to invalid signature, got nil")
	}

	if resp != nil {
		t.Fatalf("Expected response to be nil, got: %v", resp)
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestNewTripay_BaseUrl(t *testing.T) {
	sandboxClient := NewTripay(false, "api_key", "priv_key")
	if sandboxClient.BaseUrl != "https://tripay.co.id/api-sandbox" {
		t.Errorf("Expected sandbox BaseUrl 'https://tripay.co.id/api-sandbox', got: %s", sandboxClient.BaseUrl)
	}

	prodClient := NewTripay(true, "api_key", "priv_key")
	if prodClient.BaseUrl != "https://tripay.co.id/api" {
		t.Errorf("Expected production BaseUrl 'https://tripay.co.id/api', got: %s", prodClient.BaseUrl)
	}
}

