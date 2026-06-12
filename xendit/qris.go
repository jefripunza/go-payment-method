package xendit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	fiberV2 "github.com/gofiber/fiber/v2"
	fiberV3 "github.com/gofiber/fiber/v3"
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
)

// Currencies supported by Xendit Payment Requests
const (
	QrisCurrencyIndonesiaRupiah    = "IDR"
	QrisCurrencyPhilippinePeso     = "PHP"
	QrisCurrencyVietnamDong        = "VND"
	QrisCurrencyThaiBaht           = "THB"
	QrisCurrencySingaporeDollar    = "SGD"
	QrisCurrencyMalaysianRinggit   = "MYR"
	QrisCurrencyUnitedStatesDollar = "USD"
	QrisCurrencyHongKongDollar     = "HKD"
	QrisCurrencyAustralianDollar   = "AUD"
	QrisCurrencyPoundSterling      = "GBP"
	QrisCurrencyEuro               = "EUR"
	QrisCurrencyJapaneseYen        = "JPY"
	QrisCurrencyMexicanPeso        = "MXN"
)

type QrisCreateRequest struct {
	ReferenceId       string                 `json:"reference_id"`
	Type              string                 `json:"type"`
	Country           string                 `json:"country"`
	Currency          string                 `json:"currency"`
	RequestAmount     float64                `json:"request_amount"`
	CaptureMethod     string                 `json:"capture_method"`
	ChannelCode       string                 `json:"channel_code"`
	ChannelProperties QrisChannelProperties  `json:"channel_properties"`
	Description       string                 `json:"description,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

type QrisChannelProperties struct {
	ExpiresAt    string `json:"expires_at,omitempty"`
	QrStringType string `json:"qr_string_type,omitempty"`
}

type QrisResponse struct {
	BusinessId        string                 `json:"business_id"`
	ReferenceId       string                 `json:"reference_id"`
	PaymentRequestId  string                 `json:"payment_request_id"`
	Type              string                 `json:"type"`
	Country           string                 `json:"country"`
	Currency          string                 `json:"currency"`
	RequestAmount     float64                `json:"request_amount"`
	CaptureMethod     string                 `json:"capture_method"`
	ChannelCode       string                 `json:"channel_code"`
	ChannelProperties QrisChannelProperties  `json:"channel_properties"`
	Actions           []QrisAction           `json:"actions"`
	Status            string                 `json:"status"`
	Description       string                 `json:"description"`
	Metadata          map[string]interface{} `json:"metadata"`
	Created           string                 `json:"created"`
	Updated           string                 `json:"updated"`
}

type QrisAction struct {
	Type       string `json:"type"`
	Descriptor string `json:"descriptor"`
	Value      string `json:"value"`
}

func (x *Xendit) QrisCreate(referenceId string, amount float64, currency string, expiresAt time.Time, description string, metadata map[string]interface{}, forUserId ...string) (*QrisResponse, error) {
	url := fmt.Sprintf("%s/v3/payment_requests", x.BaseUrl)

	req := QrisCreateRequest{
		ReferenceId:   referenceId,
		Type:          "PAY",
		Country:       "ID",
		Currency:      currency,
		RequestAmount: amount,
		CaptureMethod: "AUTOMATIC",
		ChannelCode:   "QRIS", // find: https://docs.xendit.co/apidocs/create-payment-request
		ChannelProperties: QrisChannelProperties{
			ExpiresAt:    expiresAt.UTC().Format(time.RFC3339),
			QrStringType: "DYNAMIC",
		},
		Description: description,
		Metadata:    metadata,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, _, err := x.doRequest("POST", url, "2024-11-11", "", payload, forUserId...)
	if err != nil {
		return nil, err
	}

	var result QrisResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}

// QrisCallbackResponse represents the webhook payload for payment.capture (QRIS)
type QrisCallbackResponse struct {
	Created    string           `json:"created"`
	BusinessId string           `json:"business_id"`
	Event      string           `json:"event"`
	Data       QrisCallbackData `json:"data"`
	ApiVersion string           `json:"api_version"`
}

type QrisCallbackData struct {
	Type             string                 `json:"type"`
	Status           string                 `json:"status"`
	Country          string                 `json:"country"`
	Created          string                 `json:"created"`
	Updated          string                 `json:"updated"`
	Captures         []QrisCallbackCapture  `json:"captures"`
	Currency         string                 `json:"currency"`
	Metadata         map[string]interface{} `json:"metadata"`
	PaymentId        string                 `json:"payment_id"`
	BusinessId        string                 `json:"business_id"`
	ChannelCode      string                 `json:"channel_code"`
	ReferenceId      string                 `json:"reference_id"`
	CaptureMethod    string                 `json:"capture_method"`
	RequestAmount    float64                `json:"request_amount"`
	PaymentDetails   QrisPaymentDetails     `json:"payment_details"`
	PaymentRequestId string                 `json:"payment_request_id"`
}

type QrisCallbackCapture struct {
	CaptureId        string  `json:"capture_id"`
	CaptureAmount    float64 `json:"capture_amount"`
	CaptureTimestamp string  `json:"capture_timestamp"`
}

type QrisPaymentDetails struct {
	PayerName  string `json:"payer_name"`
	ReceiptId  string `json:"receipt_id"`
	IssuerName string `json:"issuer_name"`
}

func (x *Xendit) QrisCallbackGofiberV2(c *fiberV2.Ctx) (*QrisCallbackResponse, error) {
	callbackToken := c.Get("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		_ = c.Status(fiberV2.StatusUnauthorized).JSON(map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body QrisCallbackResponse
	if err := c.BodyParser(&body); err != nil {
		_ = c.Status(fiberV2.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

func (x *Xendit) QrisCallbackGofiberV3(c fiberV3.Ctx) (*QrisCallbackResponse, error) {
	callbackToken := c.Get("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		_ = c.Status(http.StatusUnauthorized).JSON(map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body QrisCallbackResponse
	if err := c.Bind().JSON(&body); err != nil {
		_ = c.Status(http.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

func (x *Xendit) QrisCallbackGin(c *gin.Context) (*QrisCallbackResponse, error) {
	callbackToken := c.GetHeader("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body QrisCallbackResponse
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

func (x *Xendit) QrisCallbackEcho(c echo.Context) (*QrisCallbackResponse, error) {
	callbackToken := c.Request().Header.Get("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		_ = c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body QrisCallbackResponse
	if err := c.Bind(&body); err != nil {
		_ = c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

func (x *Xendit) QrisCallbackHttp(w http.ResponseWriter, r *http.Request) (*QrisCallbackResponse, error) {
	callbackToken := r.Header.Get("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body QrisCallbackResponse
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

func (x *Xendit) QrisCallbackFasthttp(ctx *fasthttp.RequestCtx) (*QrisCallbackResponse, error) {
	callbackToken := string(ctx.Request.Header.Peek("X-Callback-Token"))
	if callbackToken == "" || callbackToken != x.CallbackToken {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		_ = json.NewEncoder(ctx).Encode(map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body QrisCallbackResponse
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
