package midtrans

import (
	"encoding/json"
	"net/http"
	"net/url"
)

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
