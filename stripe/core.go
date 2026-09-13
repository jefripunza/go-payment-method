package stripe

import (
	"net/url"
	"strconv"
)

// ==========================================
// PAYMENT INTENT (pi_...) — model pembayaran modern (2-step)
// ==========================================

// PaymentIntent adalah niat membayar — object inti integrasi Stripe modern.
type PaymentIntent struct {
	ID                 string            `json:"id,omitempty"`
	Object             string            `json:"object,omitempty"`
	Amount             int64             `json:"amount,omitempty"`
	AmountCapturable   int64             `json:"amount_capturable,omitempty"`
	AmountReceived     int64             `json:"amount_received,omitempty"`
	CaptureMethod      string            `json:"capture_method,omitempty"`
	ClientSecret       string            `json:"client_secret,omitempty"`
	ConfirmationMethod string            `json:"confirmation_method,omitempty"`
	Currency           string            `json:"currency,omitempty"`
	Customer           string            `json:"customer,omitempty"`
	Description        string            `json:"description,omitempty"`
	Livemode           bool              `json:"livemode,omitempty"`
	Metadata           map[string]string `json:"metadata,omitempty"`
	PaymentMethod      string            `json:"payment_method,omitempty"`
	PaymentMethodTypes []string          `json:"payment_method_types,omitempty"`
	ReceiptEmail       string            `json:"receipt_email,omitempty"`
	Status             string            `json:"status,omitempty"`
	LatestCharge       string            `json:"latest_charge,omitempty"`
	Created            int64             `json:"created,omitempty"`
}

// PaymentIntentParams adalah parameter untuk POST /v1/payment_intents.
type PaymentIntentParams struct {
	Amount             int64             `json:"-"`
	Currency           string            `json:"-"`
	PaymentMethodType  string            `json:"-"`
	PaymentMethodTypes []string          `json:"-"`
	Customer           string            `json:"-"`
	ReceiptEmail       string            `json:"-"`
	Description        string            `json:"-"`
	CaptureMethod      string            `json:"-"` // AUTOMATIC | MANUAL (pre-auth)
	Confirm            bool              `json:"-"`
	PaymentMethod      string            `json:"-"` // pm_... (gunakan tanpa konfirmasi)
	Metadata           map[string]string `json:"-"`
	ReturnURL          string            `json:"-"`
	OffSession         bool              `json:"-"`
	// ManualCapture / extra field tambahan jika perlu
	Extra map[string]string `json:"-"`
}

// ToForm mengubah params menjadi form urlencoded (Stripe pakai form, bukan JSON).
func (p PaymentIntentParams) ToForm() url.Values {
	f := url.Values{}
	f.Set("amount", strconv.FormatInt(p.Amount, 10))
	f.Set("currency", p.Currency)
	if p.PaymentMethodType != "" {
		f.Add("payment_method_types[]", p.PaymentMethodType)
	}
	for _, t := range p.PaymentMethodTypes {
		f.Add("payment_method_types[]", t)
	}
	if p.Customer != "" {
		f.Set("customer", p.Customer)
	}
	if p.ReceiptEmail != "" {
		f.Set("receipt_email", p.ReceiptEmail)
	}
	if p.Description != "" {
		f.Set("description", p.Description)
	}
	if p.CaptureMethod != "" {
		f.Set("capture_method", p.CaptureMethod)
	}
	if p.Confirm {
		f.Set("confirm", "true")
	}
	if p.PaymentMethod != "" {
		f.Set("payment_method", p.PaymentMethod)
	}
	if p.ReturnURL != "" {
		f.Set("return_url", p.ReturnURL)
	}
	if p.OffSession {
		f.Set("off_session", "true")
	}
	for k, v := range p.Metadata {
		f.Set("metadata["+k+"]", v)
	}
	for k, v := range p.Extra {
		f.Set(k, v)
	}
	return f
}

// CreatePaymentIntent membuat PaymentIntent (POST /v1/payment_intents).
// Returns client_secret untuk dikirim ke frontend (Payment Element / Stripe.js).
func (c *Client) CreatePaymentIntent(params PaymentIntentParams, idempotencyKey string) (*PaymentIntent, error) {
	var out PaymentIntent
	if err := c.doPost("/v1/payment_intents", params.ToForm(), idempotencyKey, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPaymentIntent mengambil PaymentIntent berdasarkan ID (GET /v1/payment_intents/{id}).
func (c *Client) GetPaymentIntent(id string) (*PaymentIntent, error) {
	var out PaymentIntent
	if err := c.doGet("/v1/payment_intents/"+url.PathEscape(id), "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdatePaymentIntent meng-update PaymentIntent (POST /v1/payment_intents/{id}).
func (c *Client) UpdatePaymentIntent(id string, form url.Values) (*PaymentIntent, error) {
	var out PaymentIntent
	if err := c.doPost("/v1/payment_intents/"+url.PathEscape(id), form, "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ConfirmPaymentIntent mengonfirmasi PaymentIntent (POST /v1/payment_intents/{id}/confirm).
// Dipakai backend untuk confirm pakai payment_method yang sudah ada.
func (c *Client) ConfirmPaymentIntent(id string, form url.Values) (*PaymentIntent, error) {
	var out PaymentIntent
	path := "/v1/payment_intents/" + url.PathEscape(id) + "/confirm"
	if form == nil {
		form = url.Values{}
	}
	if err := c.doPost(path, form, "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelPaymentIntent membatalkan PaymentIntent (POST /v1/payment_intents/{id}/cancel).
func (c *Client) CancelPaymentIntent(id string) (*PaymentIntent, error) {
	var out PaymentIntent
	path := "/v1/payment_intents/" + url.PathEscape(id) + "/cancel"
	if err := c.doPost(path, url.Values{}, "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CapturePaymentIntent meng-capture PaymentIntent manual/pre-auth (POST /v1/payment_intents/{id}/capture).
// amountToCapture optional (0 = full).
func (c *Client) CapturePaymentIntent(id string, amountToCapture int64) (*PaymentIntent, error) {
	form := url.Values{}
	if amountToCapture > 0 {
		form.Set("amount_to_capture", strconv.FormatInt(amountToCapture, 10))
	}
	var out PaymentIntent
	path := "/v1/payment_intents/" + url.PathEscape(id) + "/capture"
	if err := c.doPost(path, form, "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ApplyCustomerBalancePaymentIntent menerapkan saldo customer (POST /v1/payment_intents/{id}/apply_customer_balance).
func (c *Client) ApplyCustomerBalancePaymentIntent(id string) (*PaymentIntent, error) {
	var out PaymentIntent
	path := "/v1/payment_intents/" + url.PathEscape(id) + "/apply_customer_balance"
	if err := c.doPost(path, url.Values{}, "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListPaymentIntents mengambil daftar PaymentIntent (GET /v1/payment_intents).
func (c *Client) ListPaymentIntents(limit int, customer string) ([]PaymentIntent, error) {
	form := url.Values{}
	if limit > 0 {
		form.Set("limit", strconv.Itoa(limit))
	}
	if customer != "" {
		form.Set("customer", customer)
	}
	var out struct {
		Data []PaymentIntent `json:"data"`
	}
	if err := c.doGet("/v1/payment_intents?"+form.Encode(), "", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// SearchPaymentIntents mencari PaymentIntent (GET /v1/payment_intents/search).
func (c *Client) SearchPaymentIntents(query string) ([]PaymentIntent, error) {
	form := url.Values{}
	form.Set("query", query)
	var out struct {
		Data []PaymentIntent `json:"data"`
	}
	if err := c.doGet("/v1/payment_intents/search?"+form.Encode(), "", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// ==========================================
// SETUP INTENT (seti_...) — simpan payment method tanpa charge
// ==========================================

// SetupIntent adalah niat menyimpan metode pembayaran (untuk recurring/one-click).
type SetupIntent struct {
	ID                 string   `json:"id,omitempty"`
	Object             string   `json:"object,omitempty"`
	ClientSecret       string   `json:"client_secret,omitempty"`
	Status             string   `json:"status,omitempty"`
	Customer           string   `json:"customer,omitempty"`
	PaymentMethod      string   `json:"payment_method,omitempty"`
	PaymentMethodTypes []string `json:"payment_method_types,omitempty"`
	Created            int64    `json:"created,omitempty"`
}

// CreateSetupIntent menyimpan payment method (POST /v1/setup_intents).
func (c *Client) CreateSetupIntent(customer, paymentMethodType string, idempotencyKey string) (*SetupIntent, error) {
	form := url.Values{}
	if customer != "" {
		form.Set("customer", customer)
	}
	if paymentMethodType != "" {
		form.Add("payment_method_types[]", paymentMethodType)
	}
	var out SetupIntent
	if err := c.doPost("/v1/setup_intents", form, idempotencyKey, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSetupIntent mengambil SetupIntent (GET /v1/setup_intents/{id}).
func (c *Client) GetSetupIntent(id string) (*SetupIntent, error) {
	var out SetupIntent
	if err := c.doGet("/v1/setup_intents/"+url.PathEscape(id), "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelSetupIntent membatalkan SetupIntent (POST /v1/setup_intents/{id}/cancel).
func (c *Client) CancelSetupIntent(id string) (*SetupIntent, error) {
	var out SetupIntent
	path := "/v1/setup_intents/" + url.PathEscape(id) + "/cancel"
	if err := c.doPost(path, url.Values{}, "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ==========================================
// PAYMENT METHOD (pm_...) — data pembayaran ter-tokenize
// ==========================================

// PaymentMethod adalah metode pembayaran Stripe.
type PaymentMethod struct {
	ID       string                 `json:"id,omitempty"`
	Object   string                 `json:"object,omitempty"`
	Type     string                 `json:"type,omitempty"`
	Customer string                 `json:"customer,omitempty"`
	Card     map[string]interface{} `json:"card,omitempty"`
	Created  int64                  `json:"created,omitempty"`
	Metadata map[string]string      `json:"metadata,omitempty"`
}

// CardParams adalah field kartu untuk membuat PaymentMethod card.
type CardParams struct {
	Number   string
	ExpMonth int
	ExpYear  int
	CVC      string
}

// CreatePaymentMethodCard membuat PaymentMethod tipe card (POST /v1/payment_methods).
// type=card, card[number] dsb.
func (c *Client) CreatePaymentMethodCard(card CardParams, billingDetails url.Values, idempotencyKey string) (*PaymentMethod, error) {
	form := url.Values{}
	form.Set("type", "card")
	form.Set("card[number]", card.Number)
	form.Set("card[exp_month]", strconv.Itoa(card.ExpMonth))
	form.Set("card[exp_year]", strconv.Itoa(card.ExpYear))
	if card.CVC != "" {
		form.Set("card[cvc]", card.CVC)
	}
	for k, v := range billingDetails {
		for _, vv := range v {
			form.Set(k, vv)
		}
	}
	var out PaymentMethod
	if err := c.doPost("/v1/payment_methods", form, idempotencyKey, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPaymentMethod mengambil PaymentMethod (GET /v1/payment_methods/{id}).
func (c *Client) GetPaymentMethod(id string) (*PaymentMethod, error) {
	var out PaymentMethod
	if err := c.doGet("/v1/payment_methods/"+url.PathEscape(id), "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AttachPaymentMethod menempelkan payment method ke customer (POST /v1/payment_methods/{id}/attach).
func (c *Client) AttachPaymentMethod(id, customerID string) (*PaymentMethod, error) {
	form := url.Values{}
	form.Set("customer", customerID)
	var out PaymentMethod
	path := "/v1/payment_methods/" + url.PathEscape(id) + "/attach"
	if err := c.doPost(path, form, "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DetachPaymentMethod melepas payment method dari customer (POST /v1/payment_methods/{id}/detach).
func (c *Client) DetachPaymentMethod(id string) (*PaymentMethod, error) {
	var out PaymentMethod
	path := "/v1/payment_methods/" + url.PathEscape(id) + "/detach"
	if err := c.doPost(path, url.Values{}, "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdatePaymentMethod meng-update PaymentMethod (POST /v1/payment_methods/{id}).
func (c *Client) UpdatePaymentMethod(id string, form url.Values) (*PaymentMethod, error) {
	var out PaymentMethod
	if err := c.doPost("/v1/payment_methods/"+url.PathEscape(id), form, "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListPaymentMethods mengambil daftar PaymentMethod (GET /v1/payment_methods).
func (c *Client) ListPaymentMethods(customer, pType string, limit int) ([]PaymentMethod, error) {
	form := url.Values{}
	if customer != "" {
		form.Set("customer", customer)
	}
	if pType != "" {
		form.Set("type", pType)
	}
	if limit > 0 {
		form.Set("limit", strconv.Itoa(limit))
	}
	var out struct {
		Data []PaymentMethod `json:"data"`
	}
	if err := c.doGet("/v1/payment_methods?"+form.Encode(), "", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// ==========================================
// CUSTOMER (cus_...) — identitas + saved payment method
// ==========================================

// Customer adalah pelanggan Stripe.
type Customer struct {
	ID            string            `json:"id,omitempty"`
	Object        string            `json:"object,omitempty"`
	Email         string            `json:"email,omitempty"`
	Name          string            `json:"name,omitempty"`
	Phone         string            `json:"phone,omitempty"`
	Description   string            `json:"description,omitempty"`
	DefaultSource string            `json:"default_source,omitempty"`
	Balance       int64             `json:"balance,omitempty"`
	Currency      string            `json:"currency,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	Created       int64             `json:"created,omitempty"`
}

// CreateCustomer membuat customer (POST /v1/customers).
func (c *Client) CreateCustomer(email, name string, metadata map[string]string, idempotencyKey string) (*Customer, error) {
	form := url.Values{}
	if email != "" {
		form.Set("email", email)
	}
	if name != "" {
		form.Set("name", name)
	}
	for k, v := range metadata {
		form.Set("metadata["+k+"]", v)
	}
	var out Customer
	if err := c.doPost("/v1/customers", form, idempotencyKey, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCustomer mengambil customer (GET /v1/customers/{id}).
func (c *Client) GetCustomer(id string) (*Customer, error) {
	var out Customer
	if err := c.doGet("/v1/customers/"+url.PathEscape(id), "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateCustomer meng-update customer (POST /v1/customers/{id}).
func (c *Client) UpdateCustomer(id string, form url.Values) (*Customer, error) {
	var out Customer
	if err := c.doPost("/v1/customers/"+url.PathEscape(id), form, "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteCustomer menghapus customer (DELETE /v1/customers/{id}).
func (c *Client) DeleteCustomer(id string) error {
	var out struct {
		Deleted bool `json:"deleted"`
	}
	if err := c.doDelete("/v1/customers/"+url.PathEscape(id), "", &out); err != nil {
		return err
	}
	return nil
}

// ListCustomers mengambil daftar customer (GET /v1/customers).
func (c *Client) ListCustomers(limit int) ([]Customer, error) {
	form := url.Values{}
	if limit > 0 {
		form.Set("limit", strconv.Itoa(limit))
	}
	var out struct {
		Data []Customer `json:"data"`
	}
	if err := c.doGet("/v1/customers?"+form.Encode(), "", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// ListCustomerPaymentMethods mengambil payment method milik customer (GET /v1/customers/{id}/payment_methods).
func (c *Client) ListCustomerPaymentMethods(customerID, pType string, limit int) ([]PaymentMethod, error) {
	form := url.Values{}
	if pType != "" {
		form.Set("type", pType)
	}
	if limit > 0 {
		form.Set("limit", strconv.Itoa(limit))
	}
	var out struct {
		Data []PaymentMethod `json:"data"`
	}
	path := "/v1/customers/" + url.PathEscape(customerID) + "/payment_methods?" + form.Encode()
	if err := c.doGet(path, "", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}
