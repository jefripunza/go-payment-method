package stripe

import (
	"net/url"
	"strconv"
)

// ==========================================
// CHECKOUT SESSION — hosted checkout page
// ==========================================

// CheckoutSession adalah sesi checkout hosted Stripe.
type CheckoutSession struct {
	ID            string                 `json:"id,omitempty"`
	Object        string                 `json:"object,omitempty"`
	ClientSecret  string                 `json:"client_secret,omitempty"`
	Customer      string                 `json:"customer,omitempty"`
	CustomerEmail string                 `json:"customer_email,omitempty"`
	Mode          string                 `json:"mode,omitempty"` // payment | subscription | setup
	PaymentStatus string                 `json:"payment_status,omitempty"`
	Status        string                 `json:"status,omitempty"`
	SuccessURL    string                 `json:"success_url,omitempty"`
	CancelURL     string                 `json:"cancel_url,omitempty"`
	PaymentIntent string                 `json:"payment_intent,omitempty"`
	Subscription  string                 `json:"subscription,omitempty"`
	TotalDetails  map[string]interface{} `json:"total_details,omitempty"`
	AmountTotal   int64                  `json:"amount_total,omitempty"`
	Currency      string                 `json:"currency,omitempty"`
	URL           string                 `json:"url,omitempty"`
	ExpiresAt     int64                  `json:"expires_at,omitempty"`
	Metadata      map[string]string      `json:"metadata,omitempty"`
	Created       int64                  `json:"created,omitempty"`
}

// CheckoutItem adalah item dalam Checkout Session.
type CheckoutItem struct {
	PriceID    string
	Quantity   int64
	Adjustable bool // tampilkan adjustable quantity
}

// CreateCheckoutSession membuat Checkout Session (POST /v1/checkout/sessions).
// mode: payment | subscription | setup. successURL/cancelURL wajib.
func (c *Client) CreateCheckoutSession(mode, successURL, cancelURL string, items []CheckoutItem, customerEmail string, idempotencyKey string) (*CheckoutSession, error) {
	form := url.Values{}
	form.Set("mode", mode)
	form.Set("success_url", successURL)
	form.Set("cancel_url", cancelURL)
	if customerEmail != "" {
		form.Set("customer_email", customerEmail)
	}
	for i, item := range items {
		if item.PriceID != "" {
			form.Set("line_items["+strconv.Itoa(i)+"][price]", item.PriceID)
		}
		if item.Quantity > 0 {
			form.Set("line_items["+strconv.Itoa(i)+"][quantity]", strconv.FormatInt(item.Quantity, 10))
		}
		if item.Adjustable {
			form.Set("line_items["+strconv.Itoa(i)+"][adjustable_quantity][enabled]", "true")
		}
	}
	var out CheckoutSession
	if err := c.doPost("/v1/checkout/sessions", form, idempotencyKey, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCheckoutSession mengambil Checkout Session (GET /v1/checkout/sessions/{id}).
func (c *Client) GetCheckoutSession(id string) (*CheckoutSession, error) {
	var out CheckoutSession
	if err := c.doGet("/v1/checkout/sessions/"+url.PathEscape(id), "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ExpireCheckoutSession meng-expire Checkout Session (POST /v1/checkout/sessions/{id}/expire).
func (c *Client) ExpireCheckoutSession(id string) (*CheckoutSession, error) {
	var out CheckoutSession
	path := "/v1/checkout/sessions/" + url.PathEscape(id) + "/expire"
	if err := c.doPost(path, url.Values{}, "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ==========================================
// SUBSCRIPTION (sub_...) — billing berulang
// ==========================================

// Subscription adalah langganan berulang Stripe.
type Subscription struct {
	ID                   string                 `json:"id,omitempty"`
	Object               string                 `json:"object,omitempty"`
	Status               string                 `json:"status,omitempty"` // active | past_due | canceled | ...
	Customer             string                 `json:"customer,omitempty"`
	CurrentPeriodStart   int64                  `json:"current_period_start,omitempty"`
	CurrentPeriodEnd     int64                  `json:"current_period_end,omitempty"`
	CancelAtPeriodEnd    bool                   `json:"cancel_at_period_end,omitempty"`
	DefaultPaymentMethod string                 `json:"default_payment_method,omitempty"`
	Items                map[string]interface{} `json:"items,omitempty"`
	Metadata             map[string]string      `json:"metadata,omitempty"`
	Created              int64                  `json:"created,omitempty"`
}

// CreateSubscription membuat subscription (POST /v1/subscriptions).
// customer required, priceID untuk item pertama.
func (c *Client) CreateSubscription(customer, priceID string, quantity int64, idempotencyKey string) (*Subscription, error) {
	form := url.Values{}
	form.Set("customer", customer)
	if priceID != "" {
		form.Set("items[0][price]", priceID)
		if quantity > 0 {
			form.Set("items[0][quantity]", strconv.FormatInt(quantity, 10))
		}
	}
	var out Subscription
	if err := c.doPost("/v1/subscriptions", form, idempotencyKey, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSubscription mengambil subscription (GET /v1/subscriptions/{id}).
func (c *Client) GetSubscription(id string) (*Subscription, error) {
	var out Subscription
	if err := c.doGet("/v1/subscriptions/"+url.PathEscape(id), "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateSubscription meng-update subscription (POST /v1/subscriptions/{id}).
func (c *Client) UpdateSubscription(id string, form url.Values) (*Subscription, error) {
	var out Subscription
	if err := c.doPost("/v1/subscriptions/"+url.PathEscape(id), form, "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelSubscription membatalkan subscription (DELETE /v1/subscriptions/{id}).
// At end of period: cancelAtPeriodEnd=true (bukan langsung). Untuk langsung, kirim invoice_now & prorate.
func (c *Client) CancelSubscription(id string) (*Subscription, error) {
	var out Subscription
	if err := c.doDelete("/v1/subscriptions/"+url.PathEscape(id), "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListSubscriptions mengambil daftar subscription (GET /v1/subscriptions).
func (c *Client) ListSubscriptions(limit int, customer string) ([]Subscription, error) {
	form := url.Values{}
	if limit > 0 {
		form.Set("limit", strconv.Itoa(limit))
	}
	if customer != "" {
		form.Set("customer", customer)
	}
	var out struct {
		Data []Subscription `json:"data"`
	}
	if err := c.doGet("/v1/subscriptions?"+form.Encode(), "", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// ==========================================
// PRODUCT & PRICE — katalog
// ==========================================

// Product adalah produk/katalog Stripe.
type Product struct {
	ID          string            `json:"id,omitempty"`
	Object      string            `json:"object,omitempty"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Active      bool              `json:"active,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Created     int64             `json:"created,omitempty"`
}

// CreateProduct membuat product (POST /v1/products).
func (c *Client) CreateProduct(name, description string, idempotencyKey string) (*Product, error) {
	form := url.Values{}
	form.Set("name", name)
	if description != "" {
		form.Set("description", description)
	}
	var out Product
	if err := c.doPost("/v1/products", form, idempotencyKey, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetProduct mengambil product (GET /v1/products/{id}).
func (c *Client) GetProduct(id string) (*Product, error) {
	var out Product
	if err := c.doGet("/v1/products/"+url.PathEscape(id), "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListProducts mengambil daftar product (GET /v1/products).
func (c *Client) ListProducts(limit int, active bool) ([]Product, error) {
	form := url.Values{}
	if limit > 0 {
		form.Set("limit", strconv.Itoa(limit))
	}
	if active {
		form.Set("active", "true")
	}
	var out struct {
		Data []Product `json:"data"`
	}
	if err := c.doGet("/v1/products?"+form.Encode(), "", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// Price adalah harga untuk product Stripe.
type Price struct {
	ID         string                 `json:"id,omitempty"`
	Object     string                 `json:"object,omitempty"`
	Active     bool                   `json:"active,omitempty"`
	Currency   string                 `json:"currency,omitempty"`
	Product    string                 `json:"product,omitempty"`
	UnitAmount int64                  `json:"unit_amount,omitempty"`
	Recurring  map[string]interface{} `json:"recurring,omitempty"`
	Type       string                 `json:"type,omitempty"` // one_time | recurring
	Metadata   map[string]string      `json:"metadata,omitempty"`
	Created    int64                  `json:"created,omitempty"`
}

// CreatePrice membuat price (POST /v1/prices).
// unitAmount dalam minor units. recurring (interval: day|week|month|year) optional.
func (c *Client) CreatePrice(productID, currency string, unitAmount int64, recurringInterval string, idempotencyKey string) (*Price, error) {
	form := url.Values{}
	form.Set("product", productID)
	form.Set("currency", currency)
	form.Set("unit_amount", strconv.FormatInt(unitAmount, 10))
	if recurringInterval != "" {
		form.Set("recurring[interval]", recurringInterval)
	}
	var out Price
	if err := c.doPost("/v1/prices", form, idempotencyKey, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPrice mengambil price (GET /v1/prices/{id}).
func (c *Client) GetPrice(id string) (*Price, error) {
	var out Price
	if err := c.doGet("/v1/prices/"+url.PathEscape(id), "", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListPrices mengambil daftar price (GET /v1/prices).
func (c *Client) ListPrices(limit int, product string) ([]Price, error) {
	form := url.Values{}
	if limit > 0 {
		form.Set("limit", strconv.Itoa(limit))
	}
	if product != "" {
		form.Set("product", product)
	}
	var out struct {
		Data []Price `json:"data"`
	}
	if err := c.doGet("/v1/prices?"+form.Encode(), "", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}
