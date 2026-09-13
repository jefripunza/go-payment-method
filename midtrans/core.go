package midtrans

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// ===========================================================================
// CORE API — charge & transaction management
// ===========================================================================

// CoreTransactionDetails describes the order id and gross amount for Core API.
type CoreTransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount int    `json:"gross_amount"`
}

// CoreCustomerDetails describes the customer in a Core API charge.
type CoreCustomerDetails struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

// CoreItemDetail is a single line item in a Core API charge.
type CoreItemDetail struct {
	ID       string `json:"id,omitempty"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
}

// VAItem is a virtual account number returned in charge/status responses.
type VAItem struct {
	Bank     string `json:"bank,omitempty"`
	VaNumber string `json:"va_number,omitempty"`
}

// CoreChargeRequest is the request body for /v2/charge. Field PaymentType is
// wajib; tambahkan payment-specific object sesuai channel (mis. bank_transfer,
// credit_card, gopay, shopeepay) via AdditionalFields.
type CoreChargeRequest struct {
	PaymentType        string                 `json:"payment_type"`
	TransactionDetails CoreTransactionDetails `json:"transaction_details"`
	ItemDetails        []CoreItemDetail       `json:"item_details,omitempty"`
	CustomerDetails    *CoreCustomerDetails   `json:"customer_details,omitempty"`
	// AdditionalFields membawa objek payment-type spesifik (bank_transfer, credit_card, dst).
	AdditionalFields map[string]interface{} `json:"-"`
}

// MarshalJSON merges AdditionalFields into the main payload.
func (c CoreChargeRequest) MarshalJSON() ([]byte, error) {
	type alias CoreChargeRequest
	a := struct{ alias }{alias: alias(c)}
	base, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	if len(c.AdditionalFields) == 0 {
		return base, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(base, &m); err != nil {
		return nil, err
	}
	for k, v := range c.AdditionalFields {
		m[k] = v
	}
	return json.Marshal(m)
}

// TransactionAction is a hyperlink action in the charge/status response.
type TransactionAction struct {
	Name   string `json:"name,omitempty"`
	Method string `json:"method,omitempty"`
	URL    string `json:"url,omitempty"`
}

// ChargeResponse is the response from /v2/charge.
type ChargeResponse map[string]interface{}

// Charge performs a Core API charge (transaction/create) request.
func (m *Midtrans) Charge(req CoreChargeRequest) (ChargeResponse, error) {
	respBody, status, err := m.doRequest(http.MethodPost, "/v2/charge", req, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiError(respBody, status)
	}
	var resp ChargeResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// TransactionStatus describes the response from /v2/{order_id}/status.
type TransactionStatus struct {
	StatusCode        string                 `json:"status_code,omitempty"`
	StatusMessage     string                 `json:"status_message,omitempty"`
	TransactionID     string                 `json:"transaction_id,omitempty"`
	OrderID           string                 `json:"order_id,omitempty"`
	GrossAmount       string                 `json:"gross_amount,omitempty"`
	PaymentType       string                 `json:"payment_type,omitempty"`
	TransactionTime   string                 `json:"transaction_time,omitempty"`
	TransactionStatus string                 `json:"transaction_status,omitempty"`
	FraudStatus       string                 `json:"fraud_status,omitempty"`
	SignatureKey      string                 `json:"signature_key,omitempty"`
	VANumbers         []VAItem               `json:"va_numbers,omitempty"`
	Actions           []TransactionAction    `json:"actions,omitempty"`
	Store             string                 `json:"store,omitempty"`
	Raw               map[string]interface{} `json:"-"`
}

// GetTransactionStatus retrieves the status of a transaction by order id (GET /v2/{order_id}/status).
func (m *Midtrans) GetTransactionStatus(orderID string) (*TransactionStatus, error) {
	path := "/v2/" + url.PathEscape(orderID) + "/status"
	respBody, status, err := m.doRequest(http.MethodGet, path, nil, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiError(respBody, status)
	}
	var resp TransactionStatus
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Approve approves a challenging charge based on the fraud status (POST /v2/{order_id}/approve).
func (m *Midtrans) Approve(orderID string) (*TransactionStatus, error) {
	return m.transactionAction(orderID, "approve")
}

// Deny denies a challenging charge based on the fraud status (POST /v2/{order_id}/deny).
func (m *Midtrans) Deny(orderID string) (*TransactionStatus, error) {
	return m.transactionAction(orderID, "deny")
}

// Cancel cancels a transaction (POST /v2/{order_id}/cancel).
func (m *Midtrans) Cancel(orderID string) (*TransactionStatus, error) {
	return m.transactionAction(orderID, "cancel")
}

// Expire expires a transaction that is still waiting for payment (POST /v2/{order_id}/expire).
func (m *Midtrans) Expire(orderID string) (*TransactionStatus, error) {
	return m.transactionAction(orderID, "expire")
}

// transactionAction performs a POST action against /v2/{order_id}/{action}.
func (m *Midtrans) transactionAction(orderID, action string) (*TransactionStatus, error) {
	path := "/v2/" + url.PathEscape(orderID) + "/" + action
	respBody, status, err := m.doRequest(http.MethodPost, path, nil, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiError(respBody, status)
	}
	var resp TransactionStatus
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CaptureRequest is the request body for POST /v2/capture (pre-authorize capture).
type CaptureRequest struct {
	TransactionID string `json:"transaction_id"`
	GrossAmount   string `json:"gross_amount"`
}

// Capture captures a pre-authorized card transaction (POST /v2/capture).
// transaction_id diambil dari response Charge, gross_amount dalam bentuk string,
// contoh "10000.00" sesuai openapi resmi.
func (m *Midtrans) Capture(req CaptureRequest) (ChargeResponse, error) {
	respBody, status, err := m.doRequest(http.MethodPost, "/v2/capture", req, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiError(respBody, status)
	}
	var resp ChargeResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetTransactionStatusB2B retrieves status of B2B transactions for an order id
// (GET /v2/{order_id}/status/b2b).
func (m *Midtrans) GetTransactionStatusB2B(orderID string) (*TransactionStatus, error) {
	path := "/v2/" + url.PathEscape(orderID) + "/status/b2b"
	respBody, status, err := m.doRequest(http.MethodGet, path, nil, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiError(respBody, status)
	}
	var resp TransactionStatus
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RefundRequest is the request body for /v2/{order_id}/refund.
type RefundRequest struct {
	RefundKey string `json:"refund_key"`
	Amount    int    `json:"amount"`
	Reason    string `json:"reason"`
}

// Refund issues a refund for a succeeded transaction (POST /v2/{order_id}/refund).
func (m *Midtrans) Refund(orderID string, req RefundRequest) (*TransactionStatus, error) {
	path := "/v2/" + url.PathEscape(orderID) + "/refund"
	respBody, status, err := m.doRequest(http.MethodPost, path, req, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiError(respBody, status)
	}
	var resp TransactionStatus
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RefundOnlineDirect performs a direct online refund for card transactions
// (POST /v2/{order_id}/refund/online/direct). Amount optional (0 = full refund).
func (m *Midtrans) RefundOnlineDirect(orderID string, amount int) (*TransactionStatus, error) {
	path := "/v2/" + url.PathEscape(orderID) + "/refund/online/direct"
	body := map[string]int{}
	if amount > 0 {
		body["amount"] = amount
	}
	respBody, status, err := m.doRequest(http.MethodPost, path, body, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiError(respBody, status)
	}
	var resp TransactionStatus
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
