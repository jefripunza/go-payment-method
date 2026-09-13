package midtrans

import (
	"encoding/json"
	"net/http"
)

// ===========================================================================
// SNAP API — host-your-own checkout (POST /snap/v1/transactions)
// ===========================================================================

// SnapCustomerDetails describes the customer in a Snap transaction.
type SnapCustomerDetails struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

// SnapItemDetail is a single line item in a Snap transaction.
type SnapItemDetail struct {
	ID       string `json:"id,omitempty"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
}

// SnapTransactionDetails contains the order id and gross amount.
type SnapTransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount int    `json:"gross_amount"`
}

// SnapRequest is the request body for creating a Snap checkout transaction.
type SnapRequest struct {
	TransactionDetails SnapTransactionDetails `json:"transaction_details"`
	ItemDetails        []SnapItemDetail       `json:"item_details,omitempty"`
	CustomerDetails    *SnapCustomerDetails   `json:"customer_details,omitempty"`
	EnabledPayments    []string               `json:"enabled_payments,omitempty"`
	CreditCard         map[string]interface{} `json:"credit_card,omitempty"`
}

// SnapResponse is the response from creating a Snap transaction.
type SnapResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

// CreateSnapTransaction creates a Snap checkout transaction and returns the
// Snap token + redirect URL. Pakai untuk pop-up / redirect checkout.
func (m *Midtrans) CreateSnapTransaction(req SnapRequest) (*SnapResponse, error) {
	respBody, status, err := m.doRequest(http.MethodPost, "/snap/v1/transactions", req, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiError(respBody, status)
	}
	var resp SnapResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSnapToken is an alias proviiding a convenient way to obtain just the Snap token.
func (m *Midtrans) GetSnapToken(req SnapRequest) (string, error) {
	resp, err := m.CreateSnapTransaction(req)
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}
