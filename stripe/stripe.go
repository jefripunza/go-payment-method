// Package stripe menyediakan klien Go untuk Stripe Payment Gateway (global).
// Mengikuti openapi resmi github.com/stripe/openapi (versi GA terbaru).
//
// Autentikasi: Basic Auth (username = secret key sk_..., password kosong).
// Test vs live ditentukan oleh API key (sk_test_... vs sk_live_...), bukan URL
// — karena itu TIDAK ada argumen isProduction (boolean hanya untuk switch domain).
package stripe

import (
	"crypto/hmac"
	"crypto/sha256"
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

// BaseURL adalah base URL resmi Stripe API (v1).
const BaseURL = "https://api.stripe.com"

// Client adalah klien Stripe API.
type Client struct {
	SecretKey string // sk_test_... / sk_live_... (server-side only)
	// HTTPClient optional — override default *http.Client (untuk test/proxy).
	HTTPClient *http.Client
}

// NewClient membuat klien Stripe. secretKey menentukan env (test vs live).
func NewClient(secretKey string) *Client {
	return &Client{
		SecretKey:  secretKey,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// --- low-level request helpers ---

// doRequest melakukan request HTTP ke Stripe API (Basic Auth).
// idempotencyKey opsional — dikirim sebagai header Idempotency-Key (cegah duplikat).
func (c *Client) doRequest(method, path string, form url.Values, idempotencyKey string, out interface{}) error {
	var bodyReader io.Reader
	contentType := "application/x-www-form-urlencoded"
	if form != nil {
		bodyReader = strings.NewReader(form.Encode())
	} else {
		contentType = "application/json"
	}

	req, err := http.NewRequest(method, BaseURL+path, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.SecretKey+":")))
	if idempotencyKey != "" {
		req.Header.Set("Idempotency-Key", idempotencyKey)
	}

	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return c.apiError(respBody, resp.StatusCode)
	}

	if out != nil {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("stripe: error decoding response: %v", err)
		}
	}
	return nil
}

// doGet melakukan GET request (tanpa form).
func (c *Client) doGet(path string, idempotencyKey string, out interface{}) error {
	return c.doRequest(http.MethodGet, path, nil, idempotencyKey, out)
}

// doPost melakukan POST request dengan form.
func (c *Client) doPost(path string, form url.Values, idempotencyKey string, out interface{}) error {
	return c.doRequest(http.MethodPost, path, form, idempotencyKey, out)
}

// doDelete melakukan DELETE request.
func (c *Client) doDelete(path string, idempotencyKey string, out interface{}) error {
	return c.doRequest(http.MethodDelete, path, nil, idempotencyKey, out)
}

// apiError membangun error dari response error Stripe.
type stripeError struct {
	Error struct {
		Type    string `json:"type"`
		Code    string `json:"code"`
		Message string `json:"message"`
		Param   string `json:"param"`
	} `json:"error"`
}

func (c *Client) apiError(body []byte, status int) error {
	var se stripeError
	_ = json.Unmarshal(body, &se)
	if se.Error.Message != "" {
		return fmt.Errorf("stripe API error (status: %d, type: %s, code: %s): %s",
			status, se.Error.Type, se.Error.Code, se.Error.Message)
	}
	return fmt.Errorf("stripe API error (status: %d): %s", status, string(body))
}

// --- Webhook verification ---

// VerifyWebhookSignature memverifikasi signature webhook Stripe.
// Header Stripe-Signature berbentuk: t=<timestamp>,v1=<hex signature>
// Algoritma: signed_payload = timestamp + "." + rawBody; expected = HMAC-SHA256(signed_payload, webhookSecret).
// toleranceSeconds = toleransi waktu (default ±300 detik / 5 menit, anti-replay).
func (c *Client) VerifyWebhookSignature(payload []byte, stripeSignature string, webhookSecret string, toleranceSeconds int64) (bool, error) {
	if webhookSecret == "" {
		return false, fmt.Errorf("stripe: webhook secret kosong")
	}
	if toleranceSeconds == 0 {
		toleranceSeconds = 300
	}

	// Parse header: t=...,v1=...
	var timestamp int64
	var signature string
	for _, part := range strings.Split(stripeSignature, ",") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "t":
			ts, err := strconv.ParseInt(kv[1], 10, 64)
			if err != nil {
				return false, fmt.Errorf("stripe: invalid timestamp in signature")
			}
			timestamp = ts
		case "v1":
			signature = kv[1]
		}
	}
	if timestamp == 0 || signature == "" {
		return false, fmt.Errorf("stripe: signature header tidak lengkap")
	}

	// Cek toleransi waktu (anti-replay)
	now := time.Now().Unix()
	if now-timestamp > toleranceSeconds || timestamp-now > toleranceSeconds {
		return false, fmt.Errorf("stripe: timestamp di luar toleransi (replay?)")
	}

	signedPayload := fmt.Sprintf("%d.%s", timestamp, string(payload))
	mac := hmac.New(sha256.New, []byte(webhookSecret))
	mac.Write([]byte(signedPayload))
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature)), nil
}

// Event adalah object webhook event Stripe (evt_...).
type Event struct {
	ID              string          `json:"id,omitempty"`
	Object          string          `json:"object,omitempty"`
	APIVersion      string          `json:"api_version,omitempty"`
	Created         int64           `json:"created,omitempty"`
	Type            string          `json:"type,omitempty"`
	Data            json.RawMessage `json:"data,omitempty"`
	PendingWebhooks int             `json:"pending_webhooks,omitempty"`
	Request         json.RawMessage `json:"request,omitempty"`
	Livemode        bool            `json:"livemode,omitempty"`
}

// GetEvent mengambil event webhook berdasarkan ID (GET /v1/events/{id}).
func (c *Client) GetEvent(id string) (*Event, error) {
	var out Event
	if err := c.doGet("/v1/events/"+url.PathEscape(id), "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListEvents mengambil daftar event (GET /v1/events).
func (c *Client) ListEvents(limit int, types []string) ([]Event, error) {
	form := url.Values{}
	if limit > 0 {
		form.Set("limit", strconv.Itoa(limit))
	}
	for _, t := range types {
		form.Add("type", t)
	}
	var out struct {
		Data []Event `json:"data"`
	}
	if err := c.doGet("/v1/events?"+form.Encode(), "", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}
