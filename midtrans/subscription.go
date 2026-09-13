package midtrans

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// ===========================================================================
// SUBSCRIPTION API — /v1/subscriptions
// ===========================================================================

// SubscriptionSchedule defines the billing interval for a subscription.
type SubscriptionSchedule struct {
	Interval     int    `json:"interval"`
	IntervalUnit string `json:"interval_unit"`
	MaxInterval  int    `json:"max_interval,omitempty"`
	StartTime    string `json:"start_time,omitempty"`
}

// SubscriptionCustomerDetails describes the customer of a subscription.
type SubscriptionCustomerDetails struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

// SubscriptionRequest is the request body for creating a subscription.
type SubscriptionRequest struct {
	Name            string                       `json:"name"`
	Amount          string                       `json:"amount"`
	Currency        string                       `json:"currency"`
	PaymentType     string                       `json:"payment_type"`
	Token           string                       `json:"token,omitempty"`
	Schedule        *SubscriptionSchedule        `json:"schedule,omitempty"`
	Metadata        map[string]interface{}       `json:"metadata,omitempty"`
	CustomerDetails *SubscriptionCustomerDetails `json:"customer_details,omitempty"`
}

// CreateSubscription creates a subscription (POST /v1/subscriptions).
func (m *Midtrans) CreateSubscription(req SubscriptionRequest) (map[string]interface{}, error) {
	respBody, status, err := m.doRequest(http.MethodPost, "/v1/subscriptions", req, "")
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

// GetSubscription retrieves a subscription by id (GET /v1/subscriptions/{subscription_id}).
func (m *Midtrans) GetSubscription(subscriptionID string) (map[string]interface{}, error) {
	path := "/v1/subscriptions/" + url.PathEscape(subscriptionID)
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

// UpdateSubscription updates a subscription by id (PATCH /v1/subscriptions/{subscription_id}).
func (m *Midtrans) UpdateSubscription(subscriptionID string, changes map[string]interface{}) (map[string]interface{}, error) {
	path := "/v1/subscriptions/" + url.PathEscape(subscriptionID)
	respBody, status, err := m.doRequest(http.MethodPatch, path, changes, "")
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

// EnableSubscription activates a subscription (POST /v1/subscriptions/{subscription_id}/enable).
func (m *Midtrans) EnableSubscription(subscriptionID string) (map[string]interface{}, error) {
	path := "/v1/subscriptions/" + url.PathEscape(subscriptionID) + "/enable"
	return m.subscriptionAction(path)
}

// DisableSubscription deactivates a subscription (POST /v1/subscriptions/{subscription_id}/disable).
func (m *Midtrans) DisableSubscription(subscriptionID string) (map[string]interface{}, error) {
	path := "/v1/subscriptions/" + url.PathEscape(subscriptionID) + "/disable"
	return m.subscriptionAction(path)
}

// CancelSubscription cancels a subscription (POST /v1/subscriptions/{subscription_id}/cancel).
func (m *Midtrans) CancelSubscription(subscriptionID string) (map[string]interface{}, error) {
	path := "/v1/subscriptions/" + url.PathEscape(subscriptionID) + "/cancel"
	return m.subscriptionAction(path)
}

// subscriptionAction performs a POST action on a subscription.
func (m *Midtrans) subscriptionAction(path string) (map[string]interface{}, error) {
	respBody, status, err := m.doRequest(http.MethodPost, path, nil, "")
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
