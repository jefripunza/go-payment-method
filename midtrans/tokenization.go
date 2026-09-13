package midtrans

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

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
