package midtrans

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// ===========================================================================
// PAYMENT LINK API — /v1/payment-links
// ===========================================================================

// PaymentLinkTransactionDetails describes the order id and gross amount.
type PaymentLinkTransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount int    `json:"gross_amount"`
}

// PaymentLinkItemDetail is a line item in a payment link.
type PaymentLinkItemDetail struct {
	ID       string `json:"id,omitempty"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
}

// PaymentLinkCustomerDetails describes the customer of a payment link.
type PaymentLinkCustomerDetails struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

// PaymentLinkExpiry controls how long the payment link stays valid.
type PaymentLinkExpiry struct {
	StartTime string `json:"start_time,omitempty"`
	Unit      string `json:"unit,omitempty"`
	Duration  int    `json:"duration,omitempty"`
}

// PaymentLinkRequest is the request body for creating a payment link.
type PaymentLinkRequest struct {
	TransactionDetails PaymentLinkTransactionDetails `json:"transaction_details"`
	ItemDetails        []PaymentLinkItemDetail       `json:"item_details,omitempty"`
	CustomerDetails    *PaymentLinkCustomerDetails   `json:"customer_details,omitempty"`
	UsageLimit         int                           `json:"usage_limit,omitempty"`
	Expiry             *PaymentLinkExpiry            `json:"expiry,omitempty"`
	EnabledPayments    []string                      `json:"enabled_payments,omitempty"`
}

// PaymentLink is a created payment link.
type PaymentLink struct {
	OrderID    string `json:"order_id"`
	PaymentURL string `json:"payment_url"`
}

// CreatePaymentLink creates a reusable payment link (POST /v1/payment-links).
func (m *Midtrans) CreatePaymentLink(req PaymentLinkRequest) (*PaymentLink, error) {
	respBody, status, err := m.doRequest(http.MethodPost, "/v1/payment-links", req, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiError(respBody, status)
	}
	var resp PaymentLink
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPaymentLink retrieves a payment link by order id (GET /v1/payment-links/{order_id}).
func (m *Midtrans) GetPaymentLink(orderID string) (map[string]interface{}, error) {
	path := "/v1/payment-links/" + url.PathEscape(orderID)
	respBody, status, err := m.doRequest(http.MethodGet, path, nil, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiError(respBody, status)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// DeletePaymentLink deactivates a payment link by order id (DELETE /v1/payment-links/{order_id}).
func (m *Midtrans) DeletePaymentLink(orderID string) error {
	path := "/v1/payment-links/" + url.PathEscape(orderID)
	respBody, status, err := m.doRequest(http.MethodDelete, path, nil, "")
	if err != nil {
		return err
	}
	if status >= 400 {
		return m.apiError(respBody, status)
	}
	return nil
}
