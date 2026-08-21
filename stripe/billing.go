package stripe

import (
	"net/url"
	"strconv"
)

// ==========================================
// CHARGE (ch_...) — transaksi ter-charge (legacy direct charge)
// ==========================================

// Charge adalah transaksi Stripe yang sudah di-charge.
type Charge struct {
	ID             string            `json:"id,omitempty"`
	Object         string            `json:"object,omitempty"`
	Amount         int64             `json:"amount,omitempty"`
	AmountCaptured int64             `json:"amount_captured,omitempty"`
	AmountRefunded int64             `json:"amount_refunded,omitempty"`
	Currency       string            `json:"currency,omitempty"`
	Customer       string            `json:"customer,omitempty"`
	Description    string            `json:"description,omitempty"`
	PaymentMethod  string            `json:"payment_method,omitempty"`
	ReceiptEmail   string            `json:"receipt_email,omitempty"`
	Status         string            `json:"status,omitempty"`
	Captured       bool              `json:"captured,omitempty"`
	Paid           bool              `json:"paid,omitempty"`
	Refunded       bool              `json:"refunded,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	Created        int64             `json:"created,omitempty"`
}

// CreateCharge membuat charge langsung (POST /v1/charges) — legacy.
// amount dalam minor units, currency misal "usd".
func (c *Client) CreateCharge(amount int64, currency string, sourceOrMethod string, description string, idempotencyKey string) (*Charge, error) {
	form := url.Values{}
	form.Set("amount", strconv.FormatInt(amount, 10))
	form.Set("currency", currency)
	if sourceOrMethod != "" {
		// Stripe menerima source=... atau payment_method=... untuk charge langsung
		form.Set("payment_method", sourceOrMethod)
	}
	if description != "" {
		form.Set("description", description)
	}
	var out Charge
	if err := c.doPost("/v1/charges", form, idempotencyKey, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCharge mengambil charge (GET /v1/charges/{id}).
func (c *Client) GetCharge(id string) (*Charge, error) {
	var out Charge
	if err := c.doGet("/v1/charges/"+url.PathEscape(id), "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CaptureCharge meng-capture charge (POST /v1/charges/{id}/capture) — pre-auth.
// amount optional (0 = full).
func (c *Client) CaptureCharge(id string, amount int64) (*Charge, error) {
	form := url.Values{}
	if amount > 0 {
		form.Set("amount", strconv.FormatInt(amount, 10))
	}
	var out Charge
	path := "/v1/charges/" + url.PathEscape(id) + "/capture"
	if err := c.doPost(path, form, "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListCharges mengambil daftar charge (GET /v1/charges).
func (c *Client) ListCharges(limit int, customer string) ([]Charge, error) {
	form := url.Values{}
	if limit > 0 {
		form.Set("limit", strconv.Itoa(limit))
	}
	if customer != "" {
		form.Set("customer", customer)
	}
	var out struct {
		Data []Charge `json:"data"`
	}
	if err := c.doGet("/v1/charges?"+form.Encode(), "", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// ==========================================
// REFUND (re_...) — pengembalian dana
// ==========================================

// Refund adalah pengembalian dana Stripe.
type Refund struct {
	ID                 string            `json:"id,omitempty"`
	Object             string            `json:"object,omitempty"`
	Amount             int64             `json:"amount,omitempty"`
	BalanceTransaction string            `json:"balance_transaction,omitempty"`
	Charge             string            `json:"charge,omitempty"`
	Currency           string            `json:"currency,omitempty"`
	PaymentIntent      string            `json:"payment_intent,omitempty"`
	Reason             string            `json:"reason,omitempty"`
	Status             string            `json:"status,omitempty"`
	Metadata           map[string]string `json:"metadata,omitempty"`
	Created            int64             `json:"created,omitempty"`
}

// CreateRefund membuat refund (POST /v1/refunds).
// chargeOrPI: ID charge (ch_...) atau payment_intent (pi_...). amount 0 = full.
func (c *Client) CreateRefund(chargeOrPI string, amount int64, reason string, idempotencyKey string) (*Refund, error) {
	form := url.Values{}
	if amount > 0 {
		form.Set("amount", strconv.FormatInt(amount, 10))
	}
	if reason != "" {
		form.Set("reason", reason)
	}
	// Stripe menerima payment_intent=pi_... atau charge=ch_...
	if len(chargeOrPI) > 3 && chargeOrPI[:3] == "pi_" {
		form.Set("payment_intent", chargeOrPI)
	} else {
		form.Set("charge", chargeOrPI)
	}
	var out Refund
	if err := c.doPost("/v1/refunds", form, idempotencyKey, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRefund mengambil refund (GET /v1/refunds/{id}).
func (c *Client) GetRefund(id string) (*Refund, error) {
	var out Refund
	if err := c.doGet("/v1/refunds/"+url.PathEscape(id), "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelRefund membatalkan refund (POST /v1/refunds/{id}/cancel).
func (c *Client) CancelRefund(id string) (*Refund, error) {
	var out Refund
	path := "/v1/refunds/" + url.PathEscape(id) + "/cancel"
	if err := c.doPost(path, url.Values{}, "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListRefunds mengambil daftar refund (GET /v1/refunds).
func (c *Client) ListRefunds(limit int, charge string) ([]Refund, error) {
	form := url.Values{}
	if limit > 0 {
		form.Set("limit", strconv.Itoa(limit))
	}
	if charge != "" {
		form.Set("charge", charge)
	}
	var out struct {
		Data []Refund `json:"data"`
	}
	if err := c.doGet("/v1/refunds?"+form.Encode(), "", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// ==========================================
// PAYOUT (po_...) — transfer dana ke bank merchant (money-out)
// ==========================================

// Payout adalah transfer dana dari saldo Stripe ke bank.
type Payout struct {
	ID          string            `json:"id,omitempty"`
	Object      string            `json:"object,omitempty"`
	Amount      int64             `json:"amount,omitempty"`
	ArrivalDate int64             `json:"arrival_date,omitempty"`
	Currency    string            `json:"currency,omitempty"`
	Destination string            `json:"destination,omitempty"`
	Description string            `json:"description,omitempty"`
	Method      string            `json:"method,omitempty"`
	Status      string            `json:"status,omitempty"`
	Type        string            `json:"type,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Created     int64             `json:"created,omitempty"`
}

// CreatePayout membuat payout (POST /v1/payouts).
func (c *Client) CreatePayout(amount int64, currency, destination string, idempotencyKey string) (*Payout, error) {
	form := url.Values{}
	form.Set("amount", strconv.FormatInt(amount, 10))
	form.Set("currency", currency)
	if destination != "" {
		form.Set("destination", destination)
	}
	var out Payout
	if err := c.doPost("/v1/payouts", form, idempotencyKey, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPayout mengambil payout (GET /v1/payouts/{id}).
func (c *Client) GetPayout(id string) (*Payout, error) {
	var out Payout
	if err := c.doGet("/v1/payouts/"+url.PathEscape(id), "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelPayout membatalkan payout (POST /v1/payouts/{id}/cancel).
func (c *Client) CancelPayout(id string) (*Payout, error) {
	var out Payout
	path := "/v1/payouts/" + url.PathEscape(id) + "/cancel"
	if err := c.doPost(path, url.Values{}, "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListPayouts mengambil daftar payout (GET /v1/payouts).
func (c *Client) ListPayouts(limit int) ([]Payout, error) {
	form := url.Values{}
	if limit > 0 {
		form.Set("limit", strconv.Itoa(limit))
	}
	var out struct {
		Data []Payout `json:"data"`
	}
	if err := c.doGet("/v1/payouts?"+form.Encode(), "", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// ==========================================
// BALANCE — saldo Stripe
// ==========================================

// BalanceAmount adalah saldo dalam satu mata uang.
type BalanceAmount struct {
	Amount   int64  `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

// Balance adalah saldo Stripe (available + pending).
type Balance struct {
	Object    string          `json:"object,omitempty"`
	Available []BalanceAmount `json:"available,omitempty"`
	Pending   []BalanceAmount `json:"pending,omitempty"`
	Livemode  bool            `json:"livemode,omitempty"`
}

// GetBalance mengambil saldo (GET /v1/balance).
func (c *Client) GetBalance() (*Balance, error) {
	var out Balance
	if err := c.doGet("/v1/balance", "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetBalanceHistory mengambil riwayat saldo (GET /v1/balance/history).
func (c *Client) GetBalanceHistory(limit int) ([]BalanceTransaction, error) {
	form := url.Values{}
	if limit > 0 {
		form.Set("limit", strconv.Itoa(limit))
	}
	var out struct {
		Data []BalanceTransaction `json:"data"`
	}
	if err := c.doGet("/v1/balance/history?"+form.Encode(), "", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetBalanceHistoryItem mengambil satu transaksi dari riwayat saldo (GET /v1/balance/history/{id}).
func (c *Client) GetBalanceHistoryItem(id string) (*BalanceTransaction, error) {
	var out BalanceTransaction
	if err := c.doGet("/v1/balance/history/"+url.PathEscape(id), "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetBalanceTransactions mengambil daftar balance transaction (GET /v1/balance_transactions).
func (c *Client) GetBalanceTransactions(limit int) ([]BalanceTransaction, error) {
	form := url.Values{}
	if limit > 0 {
		form.Set("limit", strconv.Itoa(limit))
	}
	var out struct {
		Data []BalanceTransaction `json:"data"`
	}
	if err := c.doGet("/v1/balance_transactions?"+form.Encode(), "", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// BalanceTransaction adalah transaksi yang memengaruhi saldo.
type BalanceTransaction struct {
	ID       string `json:"id,omitempty"`
	Object   string `json:"object,omitempty"`
	Amount   int64  `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
	Type     string `json:"type,omitempty"`
	Status   string `json:"status,omitempty"`
	Created  int64  `json:"created,omitempty"`
}

// GetBalanceTransaction mengambil transaksi saldo (GET /v1/balance_transactions/{id}).
func (c *Client) GetBalanceTransaction(id string) (*BalanceTransaction, error) {
	var out BalanceTransaction
	if err := c.doGet("/v1/balance_transactions/"+url.PathEscape(id), "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}
