// Package midtrans menyediakan klien Go untuk Midtrans Payment Gateway Indonesia:
// Core API, Snap, IRIS Disbursement, Payment Link, Subscription,
// serta Card & GoPay Tokenization.
//
// Autentikasi: Basic Auth (username = Server Key, password kosong).
// Verifikasi webhook: signature_key = SHA512(order_id + status_code + gross_amount + ServerKey).
// Host: Snap memakai app.*, Core API/IRIS/Payment Link/Subscription memakai api.*.
package midtrans

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
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
