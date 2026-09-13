package tripay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
)

// EWalletWalletType adalah jenis e-wallet yang didukung untuk link/unlink/detail.
// Saat ini hanya DANA yang tersedia.
const EWalletWalletType = "DANA"

// EWalletLinkRequest mendefinisikan request untuk menautkan akun e-wallet
// ke merchant (POST /ewallet/link). Sesuai openapi: wallet_type, mobile_phone, signature.
type EWalletLinkRequest struct {
	WalletType  string `json:"wallet_type"`
	MobilePhone string `json:"mobile_phone"`
	// Signature diisi hasil CreateEWalletSignature(merchantCode, walletType, mobilePhone).
	// WAJIB pada mode production.
	Signature string `json:"signature,omitempty"`
}

// EWalletLinkResponse adalah response dari menautkan akun e-wallet.
type EWalletLinkResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// EWalletUnlinkRequest mendefinisikan request untuk memutus akun e-wallet
// (POST /ewallet/unlink). Sesuai openapi: wallet_type, mobile_phone, signature.
type EWalletUnlinkRequest struct {
	WalletType  string `json:"wallet_type"`
	MobilePhone string `json:"mobile_phone"`
	// Signature wajib pada mode production.
	Signature string `json:"signature,omitempty"`
}

// EWalletUnlinkResponse adalah response dari memutus akun e-wallet.
type EWalletUnlinkResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// EWalletDetailResponse adalah response dari mengambil info akun e-wallet
// (GET /ewallet/detail). Termasuk saldo.
type EWalletDetailResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message,omitempty"`
	Data    *EWalletDetailData `json:"data,omitempty"`
}

// EWalletDetailData memuat informasi akun e-wallet yang terhubung.
type EWalletDetailData struct {
	WalletType  string `json:"wallet_type,omitempty"`
	MobilePhone string `json:"mobile_phone,omitempty"`
	Balance     string `json:"balance,omitempty"`
	Currency    string `json:"currency,omitempty"`
}

// CreateEWalletSignature menghasilkan HMAC-SHA256 untuk operasi e-wallet.
// Format: HMAC-SHA256(privateKey, merchantCode + walletType + mobilePhone).
func (t *Tripay) CreateEWalletSignature(merchantCode, walletType, mobilePhone string) string {
	h := hmac.New(sha256.New, []byte(t.PrivateKey))
	h.Write([]byte(merchantCode + walletType + mobilePhone))
	return hex.EncodeToString(h.Sum(nil))
}

// LinkEWallet menautkan akun e-wallet (DANA) ke merchant (POST /ewallet/link).
// Production wajib menyertakan Signature; sandbox biasanya tidak wajib.
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

// UnlinkEWallet memutus akun e-wallet yang sudah ditautkan (POST /ewallet/unlink).
// Production wajib menyertakan Signature.
func (t *Tripay) UnlinkEWallet(req EWalletUnlinkRequest) (*EWalletUnlinkResponse, error) {
	body, err := json.Marshal(req)
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

// GetEWalletDetail mengambil informasi akun e-wallet yang terhubung (termasuk saldo),
// berdasarkan wallet_type dan mobile_phone (GET /ewallet/detail). Production saja.
func (t *Tripay) GetEWalletDetail(walletType, mobilePhone string) (*EWalletDetailResponse, error) {
	u, err := url.Parse(t.BaseUrl + "/ewallet/detail")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("wallet_type", walletType)
	q.Set("mobile_phone", mobilePhone)
	u.RawQuery = q.Encode()

	respBody, _, err := t.doRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	var resp EWalletDetailResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
