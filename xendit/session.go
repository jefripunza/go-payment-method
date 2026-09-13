package xendit

import (
	"encoding/json"
	"fmt"
)

// ==========================================
// PAYMENT SESSION MODELS (v3) — Payment Link / Components
// ==========================================

// SessionCustomer describes the customer of a payment session.
type SessionCustomer struct {
	ReferenceID  string `json:"reference_id,omitempty"`
	GivenNames   string `json:"given_names,omitempty"`
	Surname      string `json:"surname,omitempty"`
	Email        string `json:"email,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
}

// CreateSessionRequest is the request body for POST /v3/sessions.
type CreateSessionRequest struct {
	ReferenceID      string           `json:"reference_id"`
	PaymentMethod    string           `json:"payment_method"`
	Amount           int              `json:"amount,omitempty"`
	Currency         string           `json:"currency,omitempty"`
	Description      string           `json:"description,omitempty"`
	Customer         *SessionCustomer `json:"customer,omitempty"`
	SuccessReturnURL string           `json:"success_return_url,omitempty"`
	FailureReturnURL string           `json:"failure_return_url,omitempty"`
}

// SessionAction is an action available to the end-user of a session.
type SessionAction struct {
	Type       string `json:"type,omitempty"`
	URL        string `json:"url,omitempty"`
	Descriptor string `json:"descriptor,omitempty"`
}

// Session is a payment session (Payment Link / Xendit Components).
type Session struct {
	SessionID     string           `json:"session_id,omitempty"`
	PaymentMethod string           `json:"payment_method,omitempty"`
	Amount        int              `json:"amount,omitempty"`
	Currency      string           `json:"currency,omitempty"`
	Status        string           `json:"status,omitempty"`
	URL           string           `json:"url,omitempty"`
	ReferenceID   string           `json:"reference_id,omitempty"`
	Actions       []SessionAction  `json:"actions,omitempty"`
	Customer      *SessionCustomer `json:"customer,omitempty"`
	Created       string           `json:"created,omitempty"`
	Updated       string           `json:"updated,omitempty"`
}

// CreateSession membuat payment session (Payment Link / Xendit Components).
// Response `url` dipakai untuk redirect end-user ke hosted checkout.
func (x *Xendit) CreateSession(req *CreateSessionRequest) (*Session, error) {
	url := fmt.Sprintf("%s/v3/sessions", x.BaseUrl)
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, _, err := x.doRequest("POST", url, "2024-11-11", "", payload)
	if err != nil {
		return nil, err
	}

	var result Session
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSession mengambil status payment session (GET /v3/sessions/{session_id}).
func (x *Xendit) GetSession(sessionID string) (*Session, error) {
	url := fmt.Sprintf("%s/v3/sessions/%s", x.BaseUrl, sessionID)
	resp, _, err := x.doRequest("GET", url, "2024-11-11", "", nil)
	if err != nil {
		return nil, err
	}

	var result Session
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelSession membatalkan payment session (POST /v3/sessions/{session_id}/cancel).
func (x *Xendit) CancelSession(sessionID string) (*Session, error) {
	url := fmt.Sprintf("%s/v3/sessions/%s/cancel", x.BaseUrl, sessionID)
	resp, _, err := x.doRequest("POST", url, "2024-11-11", "", nil)
	if err != nil {
		return nil, err
	}

	var result Session
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ==========================================
// XENPLATFORM ACCOUNT MODELS (v3) — marketplace
// ==========================================

// AccountPublicProfile adalah profil publik dari sub-account.
type AccountPublicProfile struct {
	BusinessName string `json:"business_name,omitempty"`
	URL          string `json:"url,omitempty"`
}

// CreateAccountRequest adalah body untuk POST /v3/accounts (create sub-account).
type CreateAccountRequest struct {
	Email         string                `json:"email"`
	Type          string                `json:"type,omitempty"` // MANAGED | OWNED
	PublicProfile *AccountPublicProfile `json:"public_profile,omitempty"`
}

// Account adalah sub-account xenPlatform.
type Account struct {
	ID     string `json:"id,omitempty"`
	Email  string `json:"email,omitempty"`
	Status string `json:"status,omitempty"`
}

// CreateAccount membuat sub-account untuk platform/marketplace (POST /v3/accounts).
func (x *Xendit) CreateAccount(req *CreateAccountRequest) (*Account, error) {
	url := fmt.Sprintf("%s/v3/accounts", x.BaseUrl)
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, _, err := x.doRequest("POST", url, "2024-11-11", "", payload)
	if err != nil {
		return nil, err
	}

	var result Account
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAccount mengambil detail sub-account (GET /v3/accounts/{id}).
func (x *Xendit) GetAccount(id string) (*Account, error) {
	url := fmt.Sprintf("%s/v3/accounts/%s", x.BaseUrl, id)
	resp, _, err := x.doRequest("GET", url, "2024-11-11", "", nil)
	if err != nil {
		return nil, err
	}

	var result Account
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
