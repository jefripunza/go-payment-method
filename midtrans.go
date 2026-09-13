package payment_method

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Midtrans is a client for the Midtrans payment gateway (Core API, Snap,
// IRIS Disbursement, Payment Link, Subscription, Card & GoPay Tokenization).
// Wajib mengirim Server Key dan gunakan HTTPS. Baca https://docs.midtrans.com.
type Midtrans struct {
	ServerKey    string
	ClientKey    string
	IsProduction bool
	// HTTPClient optionally overrides the default *http.Client (useful for tests / proxy).
	HTTPClient *http.Client
	snapURL    string
	coreURL    string
}

const (
	midtransSandboxSnapURL = "https://app.sandbox.midtrans.com"
	midtransProdSnapURL    = "https://app.midtrans.com"
	midtransSandboxCoreURL = "https://api.sandbox.midtrans.com"
	midtransProdCoreURL    = "https://api.midtrans.com"
)

// NewMidtrans creates a Midtrans client. isProduction menentukan env (sandbox vs production).
// Snap memakai host app.*, sedangkan Core API / IRIS / Payment Link / Subscription
// memakai host api.* (sesuai openapi resmi).
func NewMidtrans(isProduction bool, serverKey, clientKey string) *Midtrans {
	snapURL := midtransSandboxSnapURL
	coreURL := midtransSandboxCoreURL
	if isProduction {
		snapURL = midtransProdSnapURL
		coreURL = midtransProdCoreURL
	}
	return &Midtrans{
		ServerKey:    serverKey,
		ClientKey:    clientKey,
		IsProduction: isProduction,
		HTTPClient:   &http.Client{Timeout: 30 * time.Second},
		snapURL:      snapURL,
		coreURL:      coreURL,
	}
}

// resolveBaseURL memilih host sesuai prefix path.
func (m *Midtrans) resolveBaseURL(path string) string {
	// Snap transactions dilayani host app.* (khusus pembuatan token).
	if strings.HasPrefix(path, "/snap/") {
		return m.snapURL
	}
	// Semua path lain (Core API, IRIS, Payment Link, Subscription, tokenization)
	// dilayani host api.*.
	return m.coreURL
}

// --- low-level helpers ---

// doRequest performs an HTTP request against the Midtrans API.
// It sets Authorization using Basic auth with the Server Key (as required by Midtrans).
func (m *Midtrans) doRequest(method, path string, body interface{}, contentType string) ([]byte, int, error) {
	return m.doRequestWithAuth(method, path, body, contentType, true)
}

// doRequestNoAuth performs an HTTP request tanpa Basic Auth — dipakai endpoint
// frontend (card token) yang hanya memakai client_key sebagai query param.
func (m *Midtrans) doRequestNoAuth(method, path string, body interface{}, contentType string) ([]byte, int, error) {
	return m.doRequestWithAuth(method, path, body, contentType, false)
}

// doRequestWithAuth adalah inti request; useAuth=false menghilangkan header Basic.
func (m *Midtrans) doRequestWithAuth(method, path string, body interface{}, contentType string, useAuth bool) ([]byte, int, error) {
	var bodyReader io.Reader
	if body != nil {
		var buf bytes.Buffer
		switch v := body.(type) {
		case []byte:
			buf.Write(v)
		case string:
			buf.WriteString(v)
		default:
			if err := json.NewEncoder(&buf).Encode(body); err != nil {
				return nil, 0, err
			}
		}
		bodyReader = &buf
	}

	u := m.resolveBaseURL(path) + path

	req, err := http.NewRequest(method, u, bodyReader)
	if err != nil {
		return nil, 0, err
	}

	if contentType == "" {
		contentType = "application/json"
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	if useAuth {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(m.ServerKey+":")))
	}

	client := m.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return respBody, resp.StatusCode, nil
}

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

// ===========================================================================
// IRIS DISBURSEMENT API — /iris/api/v1
// ===========================================================================

// IRISBeneficiary represents a beneficiary (payee) in IRIS.
type IRISBeneficiary struct {
	Name      string `json:"name,omitempty"`
	Account   string `json:"account,omitempty"`
	Bank      string `json:"bank,omitempty"`
	AliasName string `json:"alias_name,omitempty"`
	Email     string `json:"email,omitempty"`
}

// IRISBeneficiaryCreateResponse is the response from creating a beneficiary.
type IRISBeneficiaryCreateResponse struct {
	Status      string           `json:"status,omitempty"`
	Message     string           `json:"message,omitempty"`
	Beneficiary *IRISBeneficiary `json:"beneficiary,omitempty"`
}

// IRISPayoutItem is a single payout request item.
type IRISPayoutItem struct {
	BeneficiaryName    string `json:"beneficiary_name"`
	BeneficiaryAccount string `json:"beneficiary_account"`
	BeneficiaryBank    string `json:"beneficiary_bank"`
	BeneficiaryEmail   string `json:"beneficiary_email,omitempty"`
	Amount             string `json:"amount"`
	Notes              string `json:"notes,omitempty"`
}

// IRISPayoutRequest is the request body for creating payouts.
type IRISPayoutRequest struct {
	Payouts []IRISPayoutItem `json:"payouts"`
}

// IRISPayoutApproval is the body for approving payouts (with OTP).
type IRISPayoutApproval struct {
	ReferenceNos []string `json:"reference_nos"`
	OTP          string   `json:"otp"`
}

// IRISPayoutRejection is the body for rejecting payouts.
type IRISPayoutRejection struct {
	ReferenceNos []string `json:"reference_nos"`
	RejectReason string   `json:"reject_reason"`
}

// Ping checks IRIS API connectivity (GET /iris/api/v1/ping).
func (m *Midtrans) IRISPing() (map[string]interface{}, error) {
	respBody, status, err := m.doRequest(http.MethodGet, "/iris/api/v1/ping", nil, "")
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

// IRISListBeneficiaries lists all beneficiaries (GET /iris/api/v1/beneficiaries).
func (m *Midtrans) IRISListBeneficiaries() ([]IRISBeneficiary, error) {
	respBody, status, err := m.doRequest(http.MethodGet, "/iris/api/v1/beneficiaries", nil, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiError(respBody, status)
	}
	var resp []IRISBeneficiary
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// IRISCreateBeneficiary creates a new beneficiary (POST /iris/api/v1/beneficiaries).
func (m *Midtrans) IRISCreateBeneficiary(b IRISBeneficiary) (*IRISBeneficiaryCreateResponse, error) {
	respBody, status, err := m.doRequest(http.MethodPost, "/iris/api/v1/beneficiaries", b, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiError(respBody, status)
	}
	var resp IRISBeneficiaryCreateResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// IRISUpdateBeneficiary updates an existing beneficiary by alias name (PATCH /iris/api/v1/beneficiaries/{alias_name}).
func (m *Midtrans) IRISUpdateBeneficiary(aliasName string, b IRISBeneficiary) (map[string]interface{}, error) {
	path := "/iris/api/v1/beneficiaries/" + url.PathEscape(aliasName)
	respBody, status, err := m.doRequest(http.MethodPatch, path, b, "")
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

// IRISCreatePayouts creates one or more payouts (POST /iris/api/v1/payouts).
func (m *Midtrans) IRISCreatePayouts(req IRISPayoutRequest) (*IRISPayoutResponse, error) {
	respBody, status, err := m.doRequest(http.MethodPost, "/iris/api/v1/payouts", req, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiError(respBody, status)
	}
	var resp IRISPayoutResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// IRISPayoutResponse wraps the payouts array in the IRIS response.
type IRISPayoutResponse struct {
	Payouts []IRISPayoutResult `json:"payouts"`
}

// IRISPayoutResult is a single payout in the IRIS response.
type IRISPayoutResult struct {
	Status          string `json:"status,omitempty"`
	ReferenceNo     string `json:"reference_no,omitempty"`
	Amount          string `json:"amount,omitempty"`
	BeneficiaryName string `json:"beneficiary_name,omitempty"`
	Bank            string `json:"bank,omitempty"`
	Account         string `json:"account,omitempty"`
}

// IRISApprovePayouts approves pending payouts (POST /iris/api/v1/payouts/approve).
func (m *Midtrans) IRISApprovePayouts(approval IRISPayoutApproval) (*IRISPayoutResponse, error) {
	respBody, status, err := m.doRequest(http.MethodPost, "/iris/api/v1/payouts/approve", approval, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiError(respBody, status)
	}
	var resp IRISPayoutResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// IRISRejectPayouts rejects pending payouts (POST /iris/api/v1/payouts/reject).
func (m *Midtrans) IRISRejectPayouts(rejection IRISPayoutRejection) (*IRISPayoutResponse, error) {
	respBody, status, err := m.doRequest(http.MethodPost, "/iris/api/v1/payouts/reject", rejection, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiError(respBody, status)
	}
	var resp IRISPayoutResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// IRISGetPayoutStatus retrieves the status of a payout by reference number
// (GET /iris/api/v1/payouts/{reference_no}).
func (m *Midtrans) IRISGetPayoutStatus(referenceNo string) (*IRISPayoutResult, error) {
	path := "/iris/api/v1/payouts/" + url.PathEscape(referenceNo)
	respBody, status, err := m.doRequest(http.MethodGet, path, nil, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiError(respBody, status)
	}
	var resp IRISPayoutResult
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// IRISGetBalance retrieves the current IRIS balance (GET /iris/api/v1/balance).
func (m *Midtrans) IRISGetBalance() (map[string]interface{}, error) {
	respBody, status, err := m.doRequest(http.MethodGet, "/iris/api/v1/balance", nil, "")
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

// IRISBankAccount is a bank account in the IRIS bank list.
type IRISBankAccount struct {
	BankCode string `json:"bank_code,omitempty"`
}

// IRISListBankAccounts lists supported payout banks (GET /iris/api/v1/bank_accounts).
func (m *Midtrans) IRISListBankAccounts() ([]IRISBankAccount, error) {
	respBody, status, err := m.doRequest(http.MethodGet, "/iris/api/v1/bank_accounts", nil, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiError(respBody, status)
	}
	var resp []IRISBankAccount
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// IRISValidateAccount validates a beneficiary bank account before payout
// (GET /iris/api/v1/account_validation?bank=&account=).
func (m *Midtrans) IRISValidateAccount(bank, account string) (map[string]interface{}, error) {
	q := url.Values{}
	q.Set("bank", bank)
	q.Set("account", account)
	respBody, status, err := m.doRequest(http.MethodGet, "/iris/api/v1/account_validation?"+q.Encode(), nil, "")
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

// ===========================================================================
// CARD TOKENIZATION API — /v2/token & /v2/card/register (frontend / client-side)
// ===========================================================================

// GetCardTokenRequest adalah param untuk mengambil card token dari FRONTEND.
// Dipanggil tanpa Basic Auth — cukup query param client_key (Client Key).
type GetCardTokenRequest struct {
	ClientKey    string // Client Key (publik). Fallback ke m.ClientKey jika kosong.
	CardNumber   string
	CardExpMonth string // MM
	CardExpYear  string // YYYY
	CardCVV      string
	GrossAmount  int    // opsional
	Currency     string // opsional, default IDR
	Secure       bool   // opsional: true = hasilkan token untuk 3DS (1-click)
}

// GetCardToken mengambil token_id kartu dari browser (GET /v2/token).
// Endpoint frontend — TIDAK pakai Basic Auth, hanya Client Key.
func (m *Midtrans) GetCardToken(req GetCardTokenRequest) (map[string]interface{}, error) {
	q := url.Values{}
	clientKey := req.ClientKey
	if clientKey == "" {
		clientKey = m.ClientKey
	}
	q.Set("client_key", clientKey)
	q.Set("card_number", req.CardNumber)
	q.Set("card_exp_month", req.CardExpMonth)
	q.Set("card_exp_year", req.CardExpYear)
	if req.CardCVV != "" {
		q.Set("card_cvv", req.CardCVV)
	}
	if req.GrossAmount > 0 {
		q.Set("gross_amount", strconv.Itoa(req.GrossAmount))
	}
	if req.Currency != "" {
		q.Set("currency", req.Currency)
	}
	if req.Secure {
		q.Set("secure", "true")
	}

	respBody, status, err := m.doRequestNoAuth(http.MethodGet, "/v2/token?"+q.Encode(), nil, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiErrorNoAuth(respBody, status)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// RegisterCardRequest adalah param untuk mendaftarkan card token (1-click/two-clicks).
type RegisterCardRequest struct {
	ClientKey    string // Client Key. Fallback ke m.ClientKey jika kosong.
	CardNumber   string
	CardExpMonth string // MM
	CardExpYear  string // YYYY
	CardCVV      string
}

// RegisterCardToken mendaftarkan kartu untuk pembayaran berulang (GET /v2/card/register).
// Endpoint frontend — TIDAK pakai Basic Auth, hanya Client Key.
func (m *Midtrans) RegisterCardToken(req RegisterCardRequest) (map[string]interface{}, error) {
	q := url.Values{}
	clientKey := req.ClientKey
	if clientKey == "" {
		clientKey = m.ClientKey
	}
	q.Set("client_key", clientKey)
	q.Set("card_number", req.CardNumber)
	q.Set("card_exp_month", req.CardExpMonth)
	q.Set("card_exp_year", req.CardExpYear)
	if req.CardCVV != "" {
		q.Set("card_cvv", req.CardCVV)
	}

	respBody, status, err := m.doRequestNoAuth(http.MethodGet, "/v2/card/register?"+q.Encode(), nil, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, m.apiErrorNoAuth(respBody, status)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

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

// ===========================================================================
// Utilities
// ===========================================================================

// apiError builds a descriptive error from a Midtrans error response body.
func (m *Midtrans) apiError(body []byte, status int) error {
	var parsed struct {
		StatusMessage   string `json:"status_message"`
		StatusMessageID string `json:"status_message_"`
	}
	_ = json.Unmarshal(body, &parsed)
	msg := parsed.StatusMessage
	if msg == "" {
		msg = string(body)
	}
	return fmt.Errorf("midtrans API error (status: %d): %s", status, msg)
}

// apiErrorNoAuth sama seperti apiError, untuk endpoint frontend (tanpa auth).
func (m *Midtrans) apiErrorNoAuth(body []byte, status int) error {
	return m.apiError(body, status)
}

// VerifyNotificationSignature memverifikasi signature_key webhook Midtrans.
// Sesuai openapi: SHA512(orderId + statusCode + grossAmount + ServerKey).
func (m *Midtrans) VerifyNotificationSignature(signatureKey, orderID, statusCode, grossAmount string) bool {
	h := sha512.New()
	h.Write([]byte(orderID + statusCode + grossAmount + m.ServerKey))
	expected := hex.EncodeToString(h.Sum(nil))
	return expected == signatureKey
}
