package tripay

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type ChannelFee struct {
	Flat    float64     `json:"flat"`
	Percent interface{} `json:"percent"`
}

type PaymentChannel struct {
	Group         string      `json:"group"`
	Code          string      `json:"code"`
	Name          string      `json:"name"`
	Type          string      `json:"type"`
	FeeMerchant   ChannelFee  `json:"fee_merchant"`
	FeeCustomer   ChannelFee  `json:"fee_customer"`
	TotalFee      interface{} `json:"total_fee"`
	MinimumAmount float64     `json:"minimum_amount"`
	MaximumAmount float64     `json:"maximum_amount"`
	IconUrl       string      `json:"icon_url"`
	Active        bool        `json:"active"`
}

type MerchantPaymentChannelsResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    []PaymentChannel `json:"data"`
}

type FeeDetail struct {
	Flat    float64     `json:"flat"`
	Percent interface{} `json:"percent"`
	Min     interface{} `json:"min"`
	Max     interface{} `json:"max"`
}

type TotalFeeDetail struct {
	Merchant float64 `json:"merchant"`
	Customer float64 `json:"customer"`
}

type FeeCalculatorItem struct {
	Code     string         `json:"code"`
	Name     string         `json:"name"`
	Fee      FeeDetail      `json:"fee"`
	TotalFee TotalFeeDetail `json:"total_fee"`
}

type MerchantFeeCalculatorResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    []FeeCalculatorItem `json:"data"`
}

type MerchantTransaction struct {
	Reference         string      `json:"reference"`
	MerchantRef       string      `json:"merchant_ref"`
	PaymentMethod     string      `json:"payment_method"`
	PaymentMethodCode string      `json:"payment_method_code"`
	TotalAmount       float64     `json:"total_amount"`
	FeeMerchant       float64     `json:"fee_merchant"`
	AmountReceived    float64     `json:"amount_received"`
	Status            string      `json:"status"`
	PaidAt            interface{} `json:"paid_at"`
	CreatedAt         interface{} `json:"created_at"`
	UpdatedAt         interface{} `json:"updated_at"`
}

type PaginationMetadata struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	From        int `json:"from"`
	To          int `json:"to"`
}

type MerchantTransactionsResponse struct {
	Success    bool                  `json:"success"`
	Message    string                `json:"message"`
	Data       []MerchantTransaction `json:"data"`
	Pagination PaginationMetadata    `json:"pagination,omitempty"`
}

type MerchantTransactionsFilter struct {
	Page        int    `json:"page,omitempty"`
	PerPage     int    `json:"per_page,omitempty"`
	Sort        string `json:"sort,omitempty"`
	Reference   string `json:"reference,omitempty"`
	MerchantRef string `json:"merchant_ref,omitempty"`
	Method      string `json:"method,omitempty"`
	Status      string `json:"status,omitempty"`
}

// GetMerchantPaymentChannels retrieves the list of payment channels active on the merchant account.
func (t *Tripay) GetMerchantPaymentChannels() (*MerchantPaymentChannelsResponse, error) {
	respBody, _, err := t.doRequest(http.MethodGet, t.BaseUrl+"/merchant/payment-channel", nil)
	if err != nil {
		return nil, err
	}

	var result MerchantPaymentChannelsResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &result, nil
}

// CalculateMerchantFees calculates transactional fees for a given channel and amount.
func (t *Tripay) CalculateMerchantFees(code string, amount int) (*MerchantFeeCalculatorResponse, error) {
	u, err := url.Parse(t.BaseUrl + "/merchant/fee-calculator")
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %v", err)
	}

	q := u.Query()
	q.Set("code", code)
	q.Set("amount", strconv.Itoa(amount))
	u.RawQuery = q.Encode()

	respBody, _, err := t.doRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	var result MerchantFeeCalculatorResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &result, nil
}

// GetMerchantTransactions retrieves the list of transactions for the merchant.
func (t *Tripay) GetMerchantTransactions(filter MerchantTransactionsFilter) (*MerchantTransactionsResponse, error) {
	u, err := url.Parse(t.BaseUrl + "/merchant/transactions")
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %v", err)
	}

	q := u.Query()
	if filter.Page > 0 {
		q.Set("page", strconv.Itoa(filter.Page))
	}
	if filter.PerPage > 0 {
		q.Set("per_page", strconv.Itoa(filter.PerPage))
	}
	if filter.Sort != "" {
		q.Set("sort", filter.Sort)
	}
	if filter.Reference != "" {
		q.Set("reference", filter.Reference)
	}
	if filter.MerchantRef != "" {
		q.Set("merchant_ref", filter.MerchantRef)
	}
	if filter.Method != "" {
		q.Set("method", filter.Method)
	}
	if filter.Status != "" {
		q.Set("status", filter.Status)
	}
	u.RawQuery = q.Encode()

	respBody, _, err := t.doRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	var result MerchantTransactionsResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &result, nil
}
