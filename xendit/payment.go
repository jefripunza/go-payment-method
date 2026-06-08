package xendit

import (
	"encoding/json"
	"fmt"
)

// ==========================================
// PAYMENT REQUEST MODELS
// ==========================================

type CardDetails struct {
	CVN                   string `json:"cvn,omitempty"`
	CardNumber            string `json:"card_number,omitempty"`
	ExpiryYear            string `json:"expiry_year,omitempty"`
	ExpiryMonth           string `json:"expiry_month,omitempty"`
	CardholderFirstName   string `json:"cardholder_first_name,omitempty"`
	CardholderLastName    string `json:"cardholder_last_name,omitempty"`
	CardholderEmail       string `json:"cardholder_email,omitempty"`
	CardholderPhoneNumber string `json:"cardholder_phone_number,omitempty"`
}

type PaymentRequestChannelProperties struct {
	MidLabel         string       `json:"mid_label,omitempty"`
	CardDetails      *CardDetails `json:"card_details,omitempty"`
	SkipThreeDS      bool         `json:"skip_three_ds,omitempty"`
	FailureReturnUrl string       `json:"failure_return_url,omitempty"`
	SuccessReturnUrl string       `json:"success_return_url,omitempty"`
	CardNumber       string       `json:"card_number,omitempty"`
}

type PaymentRequest struct {
	ReferenceId       string                           `json:"reference_id,omitempty"`
	Type              string                           `json:"type,omitempty"`
	Country           string                           `json:"country,omitempty"`
	Currency          string                           `json:"currency,omitempty"`
	RequestAmount     float64                          `json:"request_amount,omitempty"`
	CaptureMethod     string                           `json:"capture_method,omitempty"`
	ChannelCode       string                           `json:"channel_code,omitempty"`
	ChannelProperties *PaymentRequestChannelProperties `json:"channel_properties,omitempty"`
	Description       string                           `json:"description,omitempty"`
	Metadata          map[string]interface{}           `json:"metadata,omitempty"`
}

type PaymentResponse struct {
	Id                string                           `json:"id,omitempty"`
	ReferenceId       string                           `json:"reference_id,omitempty"`
	BusinessId        string                           `json:"business_id,omitempty"`
	Type              string                           `json:"type,omitempty"`
	Status            string                           `json:"status,omitempty"`
	Country           string                           `json:"country,omitempty"`
	Currency          string                           `json:"currency,omitempty"`
	RequestAmount     float64                          `json:"request_amount,omitempty"`
	CaptureMethod     string                           `json:"capture_method,omitempty"`
	ChannelCode       string                           `json:"channel_code,omitempty"`
	ChannelProperties *PaymentRequestChannelProperties `json:"channel_properties,omitempty"`
	Description       string                           `json:"description,omitempty"`
	CustomerId        string                           `json:"customer_id,omitempty"`
	Metadata          map[string]interface{}           `json:"metadata,omitempty"`
	Created           string                           `json:"created,omitempty"`
	Updated           string                           `json:"updated,omitempty"`
}

type PaymentSimulationResponse struct {
	Status string `json:"status,omitempty"`
}

// ==========================================
// PAYMENT OPTIONS MODELS
// ==========================================

type PaymentOptionsRequest struct {
	ChannelCode       string                           `json:"channel_code,omitempty"`
	Country           string                           `json:"country,omitempty"`
	Amount            float64                          `json:"amount,omitempty"`
	Currency          string                           `json:"currency,omitempty"`
	CustomerId        string                           `json:"customer_id,omitempty"`
	ChannelProperties *PaymentRequestChannelProperties `json:"channel_properties,omitempty"`
}

type PaymentOptionsResponse struct {
	Data []interface{} `json:"data,omitempty"`
}

// ==========================================
// PAYMENTS MODELS
// ==========================================

type PaymentData struct {
	Id                string                           `json:"id,omitempty"`
	PaymentRequestId  string                           `json:"payment_request_id,omitempty"`
	ReferenceId       string                           `json:"reference_id,omitempty"`
	BusinessId        string                           `json:"business_id,omitempty"`
	Amount            float64                          `json:"amount,omitempty"`
	Currency          string                           `json:"currency,omitempty"`
	Status            string                           `json:"status,omitempty"`
	ChannelCode       string                           `json:"channel_code,omitempty"`
	Country           string                           `json:"country,omitempty"`
	Description       string                           `json:"description,omitempty"`
	Created           string                           `json:"created,omitempty"`
	Updated           string                           `json:"updated,omitempty"`
	FailureCode       string                           `json:"failure_code,omitempty"`
	CaptureAmount     float64                          `json:"capture_amount,omitempty"`
	ChannelProperties *PaymentRequestChannelProperties `json:"channel_properties,omitempty"`
}

type PaymentCaptureRequest struct {
	CaptureAmount float64 `json:"capture_amount"`
}

// ==========================================
// PAYMENT TOKENS MODELS
// ==========================================

type IndividualDetail struct {
	GivenNames string `json:"given_names,omitempty"`
	Surname    string `json:"surname,omitempty"`
}

type PaymentTokenCustomer struct {
	ReferenceId      string            `json:"reference_id,omitempty"`
	Type             string            `json:"type,omitempty"`
	IndividualDetail *IndividualDetail `json:"individual_detail,omitempty"`
	Email            string            `json:"email,omitempty"`
	MobileNumber     string            `json:"mobile_number,omitempty"`
}

type PaymentTokenRequest struct {
	ReferenceId       string                           `json:"reference_id,omitempty"`
	Customer          *PaymentTokenCustomer            `json:"customer,omitempty"`
	Country           string                           `json:"country,omitempty"`
	Currency          string                           `json:"currency,omitempty"`
	ChannelCode       string                           `json:"channel_code,omitempty"`
	ChannelProperties *PaymentRequestChannelProperties `json:"channel_properties,omitempty"`
	Description       string                           `json:"description,omitempty"`
	Metadata          map[string]interface{}           `json:"metadata,omitempty"`
}

type PaymentTokenAction struct {
	Type       string `json:"type,omitempty"`
	Descriptor string `json:"descriptor,omitempty"`
	Value      string `json:"value,omitempty"`
}

type TokenDetails struct {
	AccountName         string `json:"account_name,omitempty"`
	AccountBalance      string `json:"account_balance,omitempty"`
	AccountPointBalance string `json:"account_point_balance,omitempty"`
}

type PaymentTokenResponse struct {
	Id                string                           `json:"id,omitempty"`
	ReferenceId       string                           `json:"reference_id,omitempty"`
	BusinessId        string                           `json:"business_id,omitempty"`
	CustomerId        string                           `json:"customer_id,omitempty"`
	Country           string                           `json:"country,omitempty"`
	Currency          string                           `json:"currency,omitempty"`
	ChannelCode       string                           `json:"channel_code,omitempty"`
	ChannelProperties *PaymentRequestChannelProperties `json:"channel_properties,omitempty"`
	Status            string                           `json:"status,omitempty"`
	Actions           []PaymentTokenAction             `json:"actions,omitempty"`
	TokenDetails      *TokenDetails                    `json:"token_details,omitempty"`
	Metadata          map[string]interface{}           `json:"metadata,omitempty"`
	Created           string                           `json:"created,omitempty"`
	Updated           string                           `json:"updated,omitempty"`
}

// ==========================================
// REFUNDS MODELS
// ==========================================

type RefundRequest struct {
	ReferenceId      string  `json:"reference_id,omitempty"`
	PaymentRequestId string  `json:"payment_request_id,omitempty"`
	Currency         string  `json:"currency,omitempty"`
	Amount           float64 `json:"amount,omitempty"`
	Reason           string  `json:"reason,omitempty"`
}

type RefundResponse struct {
	Id                string                 `json:"id,omitempty"`
	PaymentId         string                 `json:"payment_id,omitempty"`
	InvoiceId         string                 `json:"invoice_id,omitempty"`
	Amount            float64                `json:"amount,omitempty"`
	PaymentMethodType string                 `json:"payment_method_type,omitempty"`
	ChannelCode       string                 `json:"channel_code,omitempty"`
	Currency          string                 `json:"currency,omitempty"`
	Status            string                 `json:"status,omitempty"`
	Reason            string                 `json:"reason,omitempty"`
	ReferenceId       string                 `json:"reference_id,omitempty"`
	FailureCode       string                 `json:"failure_code,omitempty"`
	RefundFeeAmount   float64                `json:"refund_fee_amount,omitempty"`
	Created           string                 `json:"created,omitempty"`
	Updated           string                 `json:"updated,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// ==========================================
// RECURRING PLANS MODELS
// ==========================================

type RecurringSchedule struct {
	Interval                   string `json:"interval,omitempty"`
	IntervalCount              int    `json:"interval_count,omitempty"`
	TotalRecurrence            int    `json:"total_recurrence,omitempty"`
	AnchorDate                 string `json:"anchor_date,omitempty"`
	RetryInterval              string `json:"retry_interval,omitempty"`
	RetryIntervalCount         int    `json:"retry_interval_count,omitempty"`
	TotalRetry                 int    `json:"total_retry,omitempty"`
	FailedAttemptNotifications []int  `json:"failed_attempt_notifications,omitempty"`
}

type RecurringPaymentToken struct {
	PaymentTokenId string `json:"payment_token_id,omitempty"`
	Rank           int    `json:"rank,omitempty"`
}

type RecurringItem struct {
	Type          string                 `json:"type,omitempty"`
	ReferenceId   string                 `json:"reference_id,omitempty"`
	Name          string                 `json:"name,omitempty"`
	NetUnitAmount float64                `json:"net_unit_amount,omitempty"`
	Quantity      int                    `json:"quantity,omitempty"`
	Url           string                 `json:"url,omitempty"`
	Category      string                 `json:"category,omitempty"`
	Subcategory   string                 `json:"subcategory,omitempty"`
	Description   string                 `json:"description,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

type RecurringPlanRequest struct {
	ReferenceId                 string                  `json:"reference_id,omitempty"`
	CustomerId                  string                  `json:"customer_id,omitempty"`
	Currency                    string                  `json:"currency,omitempty"`
	Amount                      float64                 `json:"amount,omitempty"`
	Schedule                    *RecurringSchedule      `json:"schedule,omitempty"`
	PaymentTokens               []RecurringPaymentToken `json:"payment_tokens,omitempty"`
	ImmediatePayment            bool                    `json:"immediate_payment,omitempty"`
	FailedCycleAction           string                  `json:"failed_cycle_action,omitempty"`
	NotificationChannels        []string                `json:"notification_channels,omitempty"`
	Locale                      string                  `json:"locale,omitempty"`
	PaymentLinkForFailedAttempt bool                    `json:"payment_link_for_failed_attempt,omitempty"`
	Metadata                    map[string]interface{}  `json:"metadata,omitempty"`
	Description                 string                  `json:"description,omitempty"`
	Items                       []RecurringItem         `json:"items,omitempty"`
}

type RecurringPlanUpdateRequest struct {
	Amount                      float64                 `json:"amount,omitempty"`
	Locale                      string                  `json:"locale,omitempty"`
	NotificationChannels        []string                `json:"notification_channels,omitempty"`
	Description                 string                  `json:"description,omitempty"`
	Items                       []RecurringItem         `json:"items,omitempty"`
	PaymentLinkForFailedAttempt bool                    `json:"payment_link_for_failed_attempt,omitempty"`
	Metadata                    map[string]interface{}  `json:"metadata,omitempty"`
	PaymentTokens               []RecurringPaymentToken `json:"payment_tokens,omitempty"`
	Schedule                    *RecurringSchedule      `json:"schedule,omitempty"`
}

type RecurringPlanResponse struct {
	Id                          string                  `json:"id,omitempty"`
	ReferenceId                 string                  `json:"reference_id,omitempty"`
	CustomerId                  string                  `json:"customer_id,omitempty"`
	Currency                    string                  `json:"currency,omitempty"`
	Amount                      float64                 `json:"amount,omitempty"`
	Schedule                    *RecurringSchedule      `json:"schedule,omitempty"`
	PaymentTokens               []RecurringPaymentToken `json:"payment_tokens,omitempty"`
	ImmediatePayment            bool                    `json:"immediate_payment,omitempty"`
	FailedCycleAction           string                  `json:"failed_cycle_action,omitempty"`
	NotificationChannels        []string                `json:"notification_channels,omitempty"`
	Locale                      string                  `json:"locale,omitempty"`
	PaymentLinkForFailedAttempt bool                    `json:"payment_link_for_failed_attempt,omitempty"`
	Status                      string                  `json:"status,omitempty"`
	Metadata                    map[string]interface{}  `json:"metadata,omitempty"`
	Description                 string                  `json:"description,omitempty"`
	Items                       []RecurringItem         `json:"items,omitempty"`
	Created                     string                  `json:"created,omitempty"`
	Updated                     string                  `json:"updated,omitempty"`
}

type RecurringPlanCycle struct {
	Id                 string                 `json:"id,omitempty"`
	PlanId             string                 `json:"plan_id,omitempty"`
	ScheduledTimestamp string                 `json:"scheduled_timestamp,omitempty"`
	Currency           string                 `json:"currency,omitempty"`
	Amount             float64                `json:"amount,omitempty"`
	Status             string                 `json:"status,omitempty"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

type RecurringPlanCycleUpdateRequest struct {
	ScheduledTimestamp string                 `json:"scheduled_timestamp,omitempty"`
	Currency           string                 `json:"currency,omitempty"`
	Amount             float64                `json:"amount,omitempty"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

type RecurringPlanCycleSimulationResponse struct {
	Status string `json:"status,omitempty"`
}

// ==========================================
// WEBHOOK MODELS
// ==========================================

// Payment Webhook Event
type PaymentWebhookValue struct {
	Event      string       `json:"event,omitempty"`
	BusinessId string       `json:"business_id,omitempty"`
	Created    string       `json:"created,omitempty"`
	Data       *PaymentData `json:"data,omitempty"`
}

type PaymentWebhookWrapper struct {
	Value *PaymentWebhookValue `json:"value,omitempty"`
}

type PaymentWebhookPayload struct {
	PaymentCapture       *PaymentWebhookWrapper `json:"paymentCapture,omitempty"`
	PaymentAuthorization *PaymentWebhookWrapper `json:"paymentAuthorization,omitempty"`
	PaymentFailure       *PaymentWebhookWrapper `json:"paymentFailure,omitempty"`
}

// Token Webhook Event
type TokenWebhookValue struct {
	Event      string                `json:"event,omitempty"`
	BusinessId string                `json:"business_id,omitempty"`
	Created    string                `json:"created,omitempty"`
	Data       *PaymentTokenResponse `json:"data,omitempty"`
}

type TokenWebhookWrapper struct {
	Value *TokenWebhookValue `json:"value,omitempty"`
}

type TokenWebhookPayload struct {
	TokenActivation *TokenWebhookWrapper `json:"tokenActivation,omitempty"`
	TokenFailure    *TokenWebhookWrapper `json:"tokenFailure,omitempty"`
	TokenExpiry     *TokenWebhookWrapper `json:"tokenExpiry,omitempty"`
}

// Refund Webhook Event
type RefundWebhookValueInner struct {
	Id                string                 `json:"id,omitempty"`
	PaymentId         string                 `json:"payment_id,omitempty"`
	InvoiceId         string                 `json:"invoice_id,omitempty"`
	Amount            float64                `json:"amount,omitempty"`
	PaymentMethodType string                 `json:"payment_method_type,omitempty"`
	ChannelCode       string                 `json:"channel_code,omitempty"`
	Currency          string                 `json:"currency,omitempty"`
	Status            string                 `json:"status,omitempty"`
	Reason            string                 `json:"reason,omitempty"`
	ReferenceId       string                 `json:"reference_id,omitempty"`
	FailureCode       string                 `json:"failure_code,omitempty"`
	RefundFeeAmount   float64                `json:"refund_fee_amount,omitempty"`
	Created           string                 `json:"created,omitempty"`
	Updated           string                 `json:"updated,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

type RefundWebhookValueData struct {
	Event      string                   `json:"event,omitempty"`
	BusinessId string                   `json:"business_id,omitempty"`
	Created    string                   `json:"created,omitempty"`
	Data       *RefundWebhookValueInner `json:"data,omitempty"`
}

type RefundWebhookValue struct {
	Event      string                  `json:"event,omitempty"`
	BusinessId string                  `json:"business_id,omitempty"`
	Created    string                  `json:"created,omitempty"`
	Data       *RefundWebhookValueData `json:"data,omitempty"`
}

type RefundWebhookWrapper struct {
	Value *RefundWebhookValue `json:"value,omitempty"`
}

type RefundWebhookPayload struct {
	RefundSucceeded *RefundWebhookWrapper `json:"refundSucceeded,omitempty"`
	RefundFailed    *RefundWebhookWrapper `json:"refundFailed,omitempty"`
}

// ==========================================
// CLIENT METHODS
// ==========================================

// CreatePaymentRequest: POST /v3/payment_requests
func (x *Xendit) CreatePaymentRequest(req *PaymentRequest) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/v3/payment_requests", x.BaseUrl)
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, _, err := x.doRequest("POST", url, "2024-11-11", "", payload)
	if err != nil {
		return nil, err
	}

	var result PaymentResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPaymentRequest: GET /v3/payment_requests/{payment_request_id}
func (x *Xendit) GetPaymentRequest(paymentRequestId string) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/v3/payment_requests/%s", x.BaseUrl, paymentRequestId)
	resp, _, err := x.doRequest("GET", url, "2024-11-11", "", nil)
	if err != nil {
		return nil, err
	}

	var result PaymentResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelPaymentRequest: POST /v3/payment_requests/{payment_request_id}/cancel
func (x *Xendit) CancelPaymentRequest(paymentRequestId string) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/v3/payment_requests/%s/cancel", x.BaseUrl, paymentRequestId)
	resp, _, err := x.doRequest("POST", url, "2024-11-11", "", nil)
	if err != nil {
		return nil, err
	}

	var result PaymentResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SimulatePaymentRequest: POST /v3/payment_requests/{payment_request_id}/simulate
func (x *Xendit) SimulatePaymentRequest(paymentRequestId string, amount float64) (*PaymentSimulationResponse, error) {
	url := fmt.Sprintf("%s/v3/payment_requests/%s/simulate", x.BaseUrl, paymentRequestId)
	req := map[string]interface{}{"amount": amount}
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, _, err := x.doRequest("POST", url, "2024-11-11", "", payload)
	if err != nil {
		return nil, err
	}

	var result PaymentSimulationResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPaymentOptions: POST /v3/payment_options
func (x *Xendit) GetPaymentOptions(req *PaymentOptionsRequest, businessId string) (*PaymentOptionsResponse, error) {
	url := fmt.Sprintf("%s/v3/payment_options", x.BaseUrl)
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, _, err := x.doRequest("POST", url, "2024-11-11", businessId, payload)
	if err != nil {
		return nil, err
	}

	var result PaymentOptionsResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPayment: GET /v3/payments/{payment_id}
func (x *Xendit) GetPayment(paymentId string) (*PaymentData, error) {
	url := fmt.Sprintf("%s/v3/payments/%s", x.BaseUrl, paymentId)
	resp, _, err := x.doRequest("GET", url, "2024-11-11", "", nil)
	if err != nil {
		return nil, err
	}

	var result PaymentData
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelPayment: POST /v3/payments/{payment_id}/cancel
func (x *Xendit) CancelPayment(paymentId string) (*PaymentData, error) {
	url := fmt.Sprintf("%s/v3/payments/%s/cancel", x.BaseUrl, paymentId)
	resp, _, err := x.doRequest("POST", url, "2024-11-11", "", nil)
	if err != nil {
		return nil, err
	}

	var result PaymentData
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CapturePayment: POST /v3/payments/{payment_id}/capture
func (x *Xendit) CapturePayment(paymentId string, captureAmount float64) (*PaymentData, error) {
	url := fmt.Sprintf("%s/v3/payments/%s/capture", x.BaseUrl, paymentId)
	req := PaymentCaptureRequest{CaptureAmount: captureAmount}
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, _, err := x.doRequest("POST", url, "2024-11-11", "", payload)
	if err != nil {
		return nil, err
	}

	var result PaymentData
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreatePaymentToken: POST /v3/payment_tokens
func (x *Xendit) CreatePaymentToken(req *PaymentTokenRequest) (*PaymentTokenResponse, error) {
	url := fmt.Sprintf("%s/v3/payment_tokens", x.BaseUrl)
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, _, err := x.doRequest("POST", url, "2024-11-11", "", payload)
	if err != nil {
		return nil, err
	}

	var result PaymentTokenResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPaymentToken: GET /v3/payment_tokens/{payment_token_id}
func (x *Xendit) GetPaymentToken(paymentTokenId string) (*PaymentTokenResponse, error) {
	url := fmt.Sprintf("%s/v3/payment_tokens/%s", x.BaseUrl, paymentTokenId)
	resp, _, err := x.doRequest("GET", url, "2024-11-11", "", nil)
	if err != nil {
		return nil, err
	}

	var result PaymentTokenResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelPaymentToken: POST /v3/payment_tokens/{payment_token_id}/cancel
func (x *Xendit) CancelPaymentToken(paymentTokenId string) (*PaymentTokenResponse, error) {
	url := fmt.Sprintf("%s/v3/payment_tokens/%s/cancel", x.BaseUrl, paymentTokenId)
	resp, _, err := x.doRequest("POST", url, "2024-11-11", "", nil)
	if err != nil {
		return nil, err
	}

	var result PaymentTokenResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateRefund: POST /refunds
func (x *Xendit) CreateRefund(req *RefundRequest) (*RefundResponse, error) {
	url := fmt.Sprintf("%s/refunds", x.BaseUrl)
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, _, err := x.doRequest("POST", url, "", "", payload)
	if err != nil {
		return nil, err
	}

	var result RefundResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateRecurringPlan: POST /recurring/plans
func (x *Xendit) CreateRecurringPlan(req *RecurringPlanRequest) (*RecurringPlanResponse, error) {
	url := fmt.Sprintf("%s/recurring/plans", x.BaseUrl)
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, _, err := x.doRequest("POST", url, "2026-01-01", "", payload)
	if err != nil {
		return nil, err
	}

	var result RecurringPlanResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRecurringPlan: GET /recurring/plans/{id}
func (x *Xendit) GetRecurringPlan(id string) (*RecurringPlanResponse, error) {
	url := fmt.Sprintf("%s/recurring/plans/%s", x.BaseUrl, id)
	resp, _, err := x.doRequest("GET", url, "2026-01-01", "", nil)
	if err != nil {
		return nil, err
	}

	var result RecurringPlanResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateRecurringPlan: PATCH /recurring/plans/{id}
func (x *Xendit) UpdateRecurringPlan(id string, req *RecurringPlanUpdateRequest) (*RecurringPlanResponse, error) {
	url := fmt.Sprintf("%s/recurring/plans/%s", x.BaseUrl, id)
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, _, err := x.doRequest("PATCH", url, "2026-01-01", "", payload)
	if err != nil {
		return nil, err
	}

	var result RecurringPlanResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeactivateRecurringPlan: POST /recurring/plans/{id}/deactivate
func (x *Xendit) DeactivateRecurringPlan(id string) (*RecurringPlanResponse, error) {
	url := fmt.Sprintf("%s/recurring/plans/%s/deactivate", x.BaseUrl, id)
	resp, _, err := x.doRequest("POST", url, "2026-01-01", "", nil)
	if err != nil {
		return nil, err
	}

	var result RecurringPlanResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListRecurringPlanCycles: GET /recurring/plans/{plan_id}/cycles?limit={limit}
func (x *Xendit) ListRecurringPlanCycles(planId string, limit int) ([]RecurringPlanCycle, error) {
	url := fmt.Sprintf("%s/recurring/plans/%s/cycles?limit=%d", x.BaseUrl, planId, limit)
	resp, _, err := x.doRequest("GET", url, "2026-01-01", "", nil)
	if err != nil {
		return nil, err
	}

	var result []RecurringPlanCycle
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetRecurringPlanCycle: GET /recurring/plans/{plan_id}/cycles/{cycle_id}
func (x *Xendit) GetRecurringPlanCycle(planId string, cycleId string) (*RecurringPlanCycle, error) {
	url := fmt.Sprintf("%s/recurring/plans/%s/cycles/%s", x.BaseUrl, planId, cycleId)
	resp, _, err := x.doRequest("GET", url, "2026-01-01", "", nil)
	if err != nil {
		return nil, err
	}

	var result RecurringPlanCycle
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateRecurringPlanCycle: PATCH /recurring/plans/{plan_id}/cycles/{cycle_id}
func (x *Xendit) UpdateRecurringPlanCycle(planId string, cycleId string, req *RecurringPlanCycleUpdateRequest) (*RecurringPlanCycle, error) {
	url := fmt.Sprintf("%s/recurring/plans/%s/cycles/%s", x.BaseUrl, planId, cycleId)
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, _, err := x.doRequest("PATCH", url, "2026-01-01", "", payload)
	if err != nil {
		return nil, err
	}

	var result RecurringPlanCycle
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelRecurringPlanCycle: POST /recurring/plans/{plan_id}/cycles/{cycle_id}/cancel
func (x *Xendit) CancelRecurringPlanCycle(planId string, cycleId string) (*RecurringPlanCycle, error) {
	url := fmt.Sprintf("%s/recurring/plans/%s/cycles/%s/cancel", x.BaseUrl, planId, cycleId)
	resp, _, err := x.doRequest("POST", url, "2026-01-01", "", nil)
	if err != nil {
		return nil, err
	}

	var result RecurringPlanCycle
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ForceAttemptRecurringPlanCycle: POST /recurring/plans/{plan_id}/cycles/{cycle_id}/force_attempt
func (x *Xendit) ForceAttemptRecurringPlanCycle(planId string, cycleId string) (*RecurringPlanCycle, error) {
	url := fmt.Sprintf("%s/recurring/plans/%s/cycles/%s/force_attempt", x.BaseUrl, planId, cycleId)
	resp, _, err := x.doRequest("POST", url, "2026-01-01", "", nil)
	if err != nil {
		return nil, err
	}

	var result RecurringPlanCycle
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SimulateRecurringPlanCycle: POST /recurring/plans/{plan_id}/cycles/{cycle_id}/simulate
func (x *Xendit) SimulateRecurringPlanCycle(planId string, cycleId string, amount float64) (*RecurringPlanCycleSimulationResponse, error) {
	url := fmt.Sprintf("%s/recurring/plans/%s/cycles/%s/simulate", x.BaseUrl, planId, cycleId)
	req := map[string]interface{}{"amount": amount}
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, _, err := x.doRequest("POST", url, "2026-01-01", "", payload)
	if err != nil {
		return nil, err
	}

	var result RecurringPlanCycleSimulationResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
