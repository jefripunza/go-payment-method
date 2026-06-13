package xendit

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	fiberV2 "github.com/gofiber/fiber/v2"
	fiberV3 "github.com/gofiber/fiber/v3"
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
)

// PayoutChannelProperties represents the properties of the payout channel
type PayoutChannelProperties struct {
	AccountNumber     string `json:"account_number,omitempty"`
	AccountHolderName string `json:"account_holder_name,omitempty"`
}

// PayoutCreateRequest represents the request body for creating a payout
type PayoutCreateRequest struct {
	ReferenceId       string                  `json:"reference_id"`
	ChannelCode       string                  `json:"channel_code"`
	ChannelProperties PayoutChannelProperties `json:"channel_properties"`
	Amount            float64                 `json:"amount"`
	Description       string                  `json:"description,omitempty"`
	Currency          string                  `json:"currency"`
}

// PayoutResponse represents the response body for a payout request
type PayoutResponse struct {
	Id                   string                  `json:"id"`
	Amount               float64                 `json:"amount"`
	ChannelCode          string                  `json:"channel_code"`
	Currency             string                  `json:"currency"`
	Status               string                  `json:"status"`
	Description          string                  `json:"description"`
	ReferenceId          string                  `json:"reference_id"`
	Created              string                  `json:"created"`
	Updated              string                  `json:"updated"`
	EstimatedArrivalTime string                  `json:"estimated_arrival_time"`
	BusinessId           string                  `json:"business_id"`
	ChannelProperties    PayoutChannelProperties `json:"channel_properties"`
}

// PayoutChannelAmountLimits represents the limit properties of a channel
type PayoutChannelAmountLimits struct {
	Minimum          float64 `json:"minimum"`
	Maximum          float64 `json:"maximum"`
	MinimumIncrement float64 `json:"minimum_increment"`
}

// PayoutChannelResponse represents a payout channel resource
type PayoutChannelResponse struct {
	ChannelCode     string                    `json:"channel_code"`
	ChannelCategory string                    `json:"channel_category"`
	Currency        string                    `json:"currency"`
	ChannelName     string                    `json:"channel_name"`
	AmountLimits    PayoutChannelAmountLimits `json:"amount_limits"`
}

// PayoutCallbackResponse represents the webhook payload for payout callbacks
type PayoutCallbackResponse struct {
	Created    string             `json:"created"`
	BusinessId string             `json:"business_id"`
	Event      string             `json:"event"`
	Data       PayoutCallbackData `json:"data"`
	ApiVersion string             `json:"api_version"`
}

// PayoutCallbackData represents the data inside the payout callback
type PayoutCallbackData struct {
	Id                   string  `json:"id"`
	Amount               float64 `json:"amount"`
	Status               string  `json:"status"`
	Created              string  `json:"created"`
	Updated              string  `json:"updated"`
	Currency             string  `json:"currency"`
	Description          string  `json:"description"`
	ChannelCode          string  `json:"channel_code"`
	ReferenceId          string  `json:"reference_id"`
	AccountNumber        string  `json:"account_number"`
	IdempotencyKey       string  `json:"idempotency_key"`
	ChannelCategory      string  `json:"channel_category"`
	AccountHolderName    string  `json:"account_holder_name"`
	ConnectorReference   string  `json:"connector_reference"`
	EstimatedArrivalTime string  `json:"estimated_arrival_time"`
}

// PayoutCreate creates a payout with the given parameters and Idempotency key
func (x *Xendit) PayoutCreate(
	idempotencyKey string,
	referenceId string,
	channelCode string,
	accountNumber string,
	accountHolderName string,
	amount float64,
	currency string,
	description string,
	forUserId ...string,
) (*PayoutResponse, error) {
	url := fmt.Sprintf("%s/v2/payouts", x.BaseUrl)

	req := PayoutCreateRequest{
		ReferenceId: referenceId,
		ChannelCode: channelCode,
		ChannelProperties: PayoutChannelProperties{
			AccountNumber:     accountNumber,
			AccountHolderName: accountHolderName,
		},
		Amount:      amount,
		Description: description,
		Currency:    currency,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	headers := make(map[string]string)
	if idempotencyKey != "" {
		headers["Idempotency-key"] = idempotencyKey
	}

	resp, _, err := x.doRequestWithHeaders("POST", url, "", "", headers, payload, forUserId...)
	if err != nil {
		return nil, err
	}

	var result PayoutResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}

// PayoutGet retrieves a single payout by its ID
func (x *Xendit) PayoutGet(id string, forUserId ...string) (*PayoutResponse, error) {
	url := fmt.Sprintf("%s/v2/payouts/%s", x.BaseUrl, id)

	resp, _, err := x.doRequest("GET", url, "", "", nil, forUserId...)
	if err != nil {
		return nil, err
	}

	var result PayoutResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}

// PayoutGetByReferenceId retrieves a list of payouts by their reference ID
func (x *Xendit) PayoutGetByReferenceId(referenceId string, forUserId ...string) ([]PayoutResponse, error) {
	url := fmt.Sprintf("%s/v2/payouts/?reference_id=%s", x.BaseUrl, referenceId)

	resp, _, err := x.doRequest("GET", url, "", "", nil, forUserId...)
	if err != nil {
		return nil, err
	}

	var result []PayoutResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result, nil
}

// PayoutCancel cancels a payout by its ID
func (x *Xendit) PayoutCancel(id string, forUserId ...string) (*PayoutResponse, error) {
	url := fmt.Sprintf("%s/v2/payouts/%s/cancel", x.BaseUrl, id)

	resp, _, err := x.doRequest("POST", url, "", "", nil, forUserId...)
	if err != nil {
		return nil, err
	}

	var result PayoutResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}

// GetPayoutChannels retrieves the available payout channels, optionally filtering by currency
func (x *Xendit) GetPayoutChannels(currency string, forUserId ...string) ([]PayoutChannelResponse, error) {
	url := fmt.Sprintf("%s/payouts_channels", x.BaseUrl)
	if currency != "" {
		url = fmt.Sprintf("%s?currency=%s", url, currency)
	}

	resp, _, err := x.doRequest("GET", url, "", "", nil, forUserId...)
	if err != nil {
		return nil, err
	}

	var result []PayoutChannelResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result, nil
}

// --- Callback Handlers for Payouts ---

func (x *Xendit) PayoutCallbackGofiberV2(c *fiberV2.Ctx) (*PayoutCallbackResponse, error) {
	callbackToken := c.Get("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		_ = c.Status(fiberV2.StatusUnauthorized).JSON(map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body PayoutCallbackResponse
	if err := c.BodyParser(&body); err != nil {
		_ = c.Status(fiberV2.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

func (x *Xendit) PayoutCallbackGofiberV3(c fiberV3.Ctx) (*PayoutCallbackResponse, error) {
	callbackToken := c.Get("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		_ = c.Status(http.StatusUnauthorized).JSON(map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body PayoutCallbackResponse
	if err := c.Bind().JSON(&body); err != nil {
		_ = c.Status(http.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

func (x *Xendit) PayoutCallbackGin(c *gin.Context) (*PayoutCallbackResponse, error) {
	callbackToken := c.GetHeader("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body PayoutCallbackResponse
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

func (x *Xendit) PayoutCallbackEcho(c echo.Context) (*PayoutCallbackResponse, error) {
	callbackToken := c.Request().Header.Get("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		_ = c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body PayoutCallbackResponse
	if err := c.Bind(&body); err != nil {
		_ = c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

func (x *Xendit) PayoutCallbackHttp(w http.ResponseWriter, r *http.Request) (*PayoutCallbackResponse, error) {
	callbackToken := r.Header.Get("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body PayoutCallbackResponse
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

func (x *Xendit) PayoutCallbackFasthttp(ctx *fasthttp.RequestCtx) (*PayoutCallbackResponse, error) {
	callbackToken := string(ctx.Request.Header.Peek("X-Callback-Token"))
	if callbackToken == "" || callbackToken != x.CallbackToken {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		_ = json.NewEncoder(ctx).Encode(map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body PayoutCallbackResponse
	if err := json.Unmarshal(ctx.PostBody(), &body); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		_ = json.NewEncoder(ctx).Encode(map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}
