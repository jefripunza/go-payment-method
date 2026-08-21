package tripay

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

// EWalletLinkRequest defines the request for linking an e-wallet (OVO/DANA) account
// to the merchant (tripay.co.id/developer E-Wallet Link API).
type EWalletLinkRequest struct {
	CustomerName  string `json:"customer_name"`
	CustomerPhone string `json:"customer_phone"`
	CallbackURL   string `json:"callback_url,omitempty"`
}

// EWalletLinkResponse is the response from linking an e-wallet account.
type EWalletLinkResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    *EWalletLinkData `json:"data,omitempty"`
}

// EWalletLinkData holds the linked e-wallet account details.
type EWalletLinkData struct {
	Reference     string `json:"reference"`
	CustomerName  string `json:"customer_name"`
	CustomerPhone string `json:"customer_phone"`
	Status        string `json:"status"`
}

// EWalletUnlinkResponse is the response from unlinking an e-wallet account.
type EWalletUnlinkResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// EWalletDetailResponse is the response from fetching an e-wallet account detail.
type EWalletDetailResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    *EWalletDetailData `json:"data,omitempty"`
}

// EWalletDetailData holds details of a linked e-wallet account.
type EWalletDetailData struct {
	Reference     string `json:"reference"`
	Status        string `json:"status"`
	MerchantRef   string `json:"merchant_ref"`
	CustomerName  string `json:"customer_name"`
	CustomerPhone string `json:"customer_phone"`
	EwalletType   string `json:"ewallet_type"`
	Note          string `json:"note"`
}

// LinkEWallet links an e-wallet (OVO/DANA) account to the merchant.
// Requires the customer phone number tied to the e-wallet.
func (t *Tripay) LinkEWallet(req EWalletLinkRequest) (*EWalletLinkResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	respBody, _, err := t.doRequest(http.MethodPost, t.BaseUrl+"/ewallet/link", body)
	if err != nil {
		return nil, err
	}
	var resp EWalletLinkResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UnlinkEWallet unlinks an e-wallet account from the merchant by its reference.
func (t *Tripay) UnlinkEWallet(reference string) (*EWalletUnlinkResponse, error) {
	body, err := json.Marshal(map[string]string{"reference": reference})
	if err != nil {
		return nil, err
	}
	respBody, _, err := t.doRequest(http.MethodPost, t.BaseUrl+"/ewallet/unlink", body)
	if err != nil {
		return nil, err
	}
	var resp EWalletUnlinkResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetEWalletDetail retrieves the status/detail of a linked e-wallet account by reference.
func (t *Tripay) GetEWalletDetail(reference string) (*EWalletDetailResponse, error) {
	u := t.BaseUrl + "/ewallet/detail"
	params := url.Values{}
	params.Set("reference", reference)
	qs := params.Encode()
	if strings.Contains(u, "?") {
		u += "&" + qs
	} else {
		u += "?" + qs
	}
	respBody, _, err := t.doRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	var resp EWalletDetailResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
