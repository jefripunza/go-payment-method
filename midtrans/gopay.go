package midtrans

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// ===========================================================================
// GOPAY TOKENIZATION API — /v2/pay/account
// ===========================================================================

// CreateGoPayTokenRequest is the body for creating a GoPay account token.
type CreateGoPayTokenRequest struct {
	PaymentType     string                 `json:"payment_type"`
	Gopay           map[string]interface{} `json:"gopay,omitempty"`
	CustomerDetails *CoreCustomerDetails   `json:"customer_details,omitempty"`
}

// CreateGoPayAccountToken creates a GoPay account token (POST /v2/pay/account).
func (m *Midtrans) CreateGoPayAccountToken(req CreateGoPayTokenRequest) (map[string]interface{}, error) {
	respBody, status, err := m.doRequest(http.MethodPost, "/v2/pay/account", req, "")
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

// GetGoPayAccountStatus retrieves the status of a GoPay account token
// (GET /v2/pay/account/{account_id}).
func (m *Midtrans) GetGoPayAccountStatus(accountID string) (map[string]interface{}, error) {
	path := "/v2/pay/account/" + url.PathEscape(accountID)
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
