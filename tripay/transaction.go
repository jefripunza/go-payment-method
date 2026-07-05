package tripay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateOpenPaymentRequest struct {
	Method        string `json:"method"`
	MerchantRef   string `json:"merchant_ref"`
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email,omitempty"`
	CustomerPhone string `json:"customer_phone,omitempty"`
	Signature     string `json:"signature"`
}

type OpenPaymentData struct {
	Uuid              string      `json:"uuid"`
	Reference         string      `json:"reference"`
	MerchantRef       string      `json:"merchant_ref"`
	PaymentName       string      `json:"payment_name"`
	PaymentMethod     string      `json:"payment_method"`
	PaymentMethodCode string      `json:"payment_method_code"`
	PayCode           string      `json:"pay_code"`
	TotalAmount       float64     `json:"total_amount"`
	Status            string      `json:"status"`
	CreatedAt         interface{} `json:"created_at"`
}

type CreateOpenPaymentResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    OpenPaymentData `json:"data"`
}

type OpenPaymentDetailResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    OpenPaymentData `json:"data"`
}

type OpenPaymentTransactionItem struct {
	Reference         string      `json:"reference"`
	MerchantRef       string      `json:"merchant_ref"`
	PaymentMethod     string      `json:"payment_method"`
	PaymentMethodCode string      `json:"payment_method_code"`
	TotalAmount       float64     `json:"total_amount"`
	FeeMerchant       float64     `json:"fee_merchant"`
	AmountReceived    float64     `json:"amount_received"`
	Status            string      `json:"status"`
	PaidAt            interface{} `json:"paid_at"`
	CreatedAt         interface{} `json:"created_at"`
}

type OpenPaymentTransactionsResponse struct {
	Success bool                         `json:"success"`
	Message string                       `json:"message"`
	Data    []OpenPaymentTransactionItem `json:"data"`
}

// CreateOpenPaymentSignature generates HMAC-SHA256 signature for Open Payment transaction.
func (t *Tripay) CreateOpenPaymentSignature(merchantCode string, channel string, merchantRef string) string {
	h := hmac.New(sha256.New, []byte(t.PrivateKey))
	h.Write([]byte(merchantCode + channel + merchantRef))
	return hex.EncodeToString(h.Sum(nil))
}

// CreateOpenPayment generates open payment code / link.
func (t *Tripay) CreateOpenPayment(req CreateOpenPaymentRequest) (*CreateOpenPaymentResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	respBody, _, err := t.doRequest(http.MethodPost, t.BaseUrl+"/transaction/open-payment/create", reqBody)
	if err != nil {
		return nil, err
	}

	var result CreateOpenPaymentResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &result, nil
}

// GetOpenPaymentDetail retrieves details of an open payment by its UUID.
func (t *Tripay) GetOpenPaymentDetail(uuid string) (*OpenPaymentDetailResponse, error) {
	respBody, _, err := t.doRequest(http.MethodGet, t.BaseUrl+"/open-payment/"+uuid+"/detail", nil)
	if err != nil {
		return nil, err
	}

	var result OpenPaymentDetailResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &result, nil
}

// GetOpenPaymentTransactions retrieves a list of payments made toward a specific Open Payment UUID.
func (t *Tripay) GetOpenPaymentTransactions(uuid string) (*OpenPaymentTransactionsResponse, error) {
	respBody, _, err := t.doRequest(http.MethodGet, t.BaseUrl+"/open-payment/"+uuid+"/transactions", nil)
	if err != nil {
		return nil, err
	}

	var result OpenPaymentTransactionsResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &result, nil
}
