package tripay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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

type OrderItem struct {
	Sku      string `json:"sku,omitempty"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
	Subtotal int    `json:"subtotal,omitempty"`
}

type CreateClosedTransactionRequest struct {
	Method        string      `json:"method"`
	MerchantRef   string      `json:"merchant_ref"`
	Amount        int         `json:"amount"`
	CustomerName  string      `json:"customer_name"`
	CustomerEmail string      `json:"customer_email,omitempty"`
	CustomerPhone string      `json:"customer_phone,omitempty"`
	OrderItems    []OrderItem `json:"order_items"`
	CallbackUrl   string      `json:"callback_url,omitempty"`
	ReturnUrl     string      `json:"return_url,omitempty"`
	ExpiredTime   int64       `json:"expired_time,omitempty"`
	Signature     string      `json:"signature"`
}

type TransactionInstruction struct {
	Title string   `json:"title"`
	Steps []string `json:"steps"`
}

type ClosedTransactionData struct {
	Reference            string                   `json:"reference"`
	MerchantRef          string                   `json:"merchant_ref"`
	PaymentSelectionType string                   `json:"payment_selection_type,omitempty"`
	PaymentName          string                   `json:"payment_name"`
	PaymentMethod        string                   `json:"payment_method"`
	PaymentMethodCode    string                   `json:"payment_method_code"`
	CustomerName         string                   `json:"customer_name,omitempty"`
	CustomerEmail        string                   `json:"customer_email,omitempty"`
	CustomerPhone        string                   `json:"customer_phone,omitempty"`
	CallbackUrl          string                   `json:"callback_url,omitempty"`
	ReturnUrl            string                   `json:"return_url,omitempty"`
	Amount               float64                  `json:"amount"`
	TotalAmount          float64                  `json:"total_amount,omitempty"`
	FeeMerchant          float64                  `json:"fee_merchant"`
	FeeCustomer          float64                  `json:"fee_customer"`
	TotalFee             float64                  `json:"total_fee"`
	AmountReceived       float64                  `json:"amount_received"`
	PayCode              string                   `json:"pay_code"`
	PayUrl               string                   `json:"pay_url"`
	CheckoutUrl          string                   `json:"checkout_url"`
	Status               string                   `json:"status"`
	PaidAt               int64                    `json:"paid_at,omitempty"`
	ExpiredTime          int64                    `json:"expired_time"`
	OrderItems           []OrderItem              `json:"order_items"`
	Instructions         []TransactionInstruction `json:"instructions,omitempty"`
}

type CreateClosedTransactionResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Data    ClosedTransactionData `json:"data"`
}

type ClosedTransactionDetailResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Data    ClosedTransactionData `json:"data"`
}

type ClosedTransactionStatusData struct {
	Reference   string `json:"reference"`
	MerchantRef string `json:"merchant_ref"`
	Status      string `json:"status"`
}

type ClosedTransactionStatusResponse struct {
	Success bool                        `json:"success"`
	Message string                      `json:"message"`
	Data    ClosedTransactionStatusData `json:"data"`
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

	respBody, _, err := t.doRequest(http.MethodPost, t.BaseUrl+"/open-payment/create", reqBody)
	if err != nil {
		return nil, err
	}

	if debug {
		_ = os.WriteFile("open_payment.json", respBody, 0644)
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

	if debug {
		_ = os.WriteFile("open_payment_detail.json", respBody, 0644)
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

	if debug {
		_ = os.WriteFile("open_payment_transactions.json", respBody, 0644)
	}

	var result OpenPaymentTransactionsResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &result, nil
}

// CreateClosedTransactionSignature generates HMAC-SHA256 signature for Closed Payment transaction.
func (t *Tripay) CreateClosedTransactionSignature(merchantCode string, merchantRef string, amount int) string {
	h := hmac.New(sha256.New, []byte(t.PrivateKey))
	h.Write([]byte(merchantCode + merchantRef + strconv.Itoa(amount)))
	return hex.EncodeToString(h.Sum(nil))
}

// CreateClosedTransaction creates a closed payment transaction.
func (t *Tripay) CreateClosedTransaction(req CreateClosedTransactionRequest) (*CreateClosedTransactionResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	respBody, _, err := t.doRequest(http.MethodPost, t.BaseUrl+"/transaction/create", reqBody)
	if err != nil {
		return nil, err
	}

	if debug {
		_ = os.WriteFile("closed_payment.json", respBody, 0644)
	}

	var result CreateClosedTransactionResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &result, nil
}

// GetClosedTransactionDetail retrieves details of a closed transaction by its reference.
func (t *Tripay) GetClosedTransactionDetail(reference string) (*ClosedTransactionDetailResponse, error) {
	u, err := url.Parse(t.BaseUrl + "/transaction/detail")
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %v", err)
	}

	q := u.Query()
	q.Set("reference", reference)
	u.RawQuery = q.Encode()

	respBody, _, err := t.doRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if debug {
		_ = os.WriteFile("closed_payment_detail.json", respBody, 0644)
	}

	var result ClosedTransactionDetailResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &result, nil
}

// CheckClosedTransactionStatus checks the status of a closed transaction by its reference.
func (t *Tripay) CheckClosedTransactionStatus(reference string) (*ClosedTransactionStatusResponse, error) {
	u, err := url.Parse(t.BaseUrl + "/transaction/check-status")
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %v", err)
	}

	q := u.Query()
	q.Set("reference", reference)
	u.RawQuery = q.Encode()

	respBody, _, err := t.doRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if debug {
		_ = os.WriteFile("closed_payment_status.json", respBody, 0644)
	}

	var result ClosedTransactionStatusResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	// If Data fields are empty, robustly parse from Message and request parameters
	if result.Data.Reference == "" {
		result.Data.Reference = reference
	}
	if result.Data.Status == "" && result.Message != "" {
		msgUpper := strings.ToUpper(result.Message)
		if strings.Contains(msgUpper, "BELUM DIBAYAR") || strings.Contains(msgUpper, "UNPAID") {
			result.Data.Status = "UNPAID"
		} else if strings.Contains(msgUpper, "DIBAYAR") || strings.Contains(msgUpper, "PAID") {
			result.Data.Status = "PAID"
		} else if strings.Contains(msgUpper, "KADALUARSA") || strings.Contains(msgUpper, "EXPIRED") {
			result.Data.Status = "EXPIRED"
		} else if strings.Contains(msgUpper, "GAGAL") || strings.Contains(msgUpper, "FAILED") {
			result.Data.Status = "FAILED"
		}
	}

	return &result, nil
}
