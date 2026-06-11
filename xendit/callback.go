package xendit

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	fiber "github.com/gofiber/fiber/v2"
	fiberV3 "github.com/gofiber/fiber/v3"
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
)

type CallbackInvoiceResponse struct {
	Id                      string      `json:"id"`
	ExternalId              string      `json:"external_id"`
	UserId                  string      `json:"user_id"`
	PaymentMethod           string      `json:"payment_method"`
	Status                  string      `json:"status"`
	MerchantName            string      `json:"merchant_name"`
	MerchantProfilePicUrl   string      `json:"merchant_profile_picture_url"`
	Amount                  float64     `json:"amount"`
	PayerEmail              string      `json:"payer_email"`
	Description             string      `json:"description"`
	ExpiryDate              string      `json:"expiry_date"`
	InvoiceUrl              string      `json:"invoice_url"`
	AvailableRetailOutlets  []string    `json:"available_retail_outlets"`
	AvailableEwallets       []string    `json:"available_ewallets"`
	AvailableQRCodes        []string    `json:"available_qr_codes"`
	AvailableDirectDebits   []string    `json:"available_direct_debits"`
	AvailablePaylaters      []string    `json:"available_paylaters"`
	ShouldExcludeCreditCard bool        `json:"should_exclude_credit_card"`
	ShouldSendEmail         bool        `json:"should_send_email"`
	Created                 string      `json:"created"`
	Updated                 string      `json:"updated"`
	Currency                string      `json:"currency"`
	Metadata                interface{} `json:"metadata"`
	Error                   string      `json:"error,omitempty"`
}

type CallbackDisbursementResponse struct {
	Id                      string   `json:"id"`
	UserId                  string   `json:"user_id"`
	ExternalId              string   `json:"external_id"`
	Amount                  float64  `json:"amount"`
	BankCode                string   `json:"bank_code"`
	AccountHolderName       string   `json:"account_holder_name"`
	DisbursementDescription string   `json:"disbursement_description"`
	FailureCode             string   `json:"failure_code"`
	IsInstant               bool     `json:"is_instant"`
	Status                  string   `json:"status"`
	Updated                 string   `json:"updated"`
	Created                 string   `json:"created"`
	EmailTo                 []string `json:"email_to"`
	EmailCc                 []string `json:"email_cc"`
	EmailBcc                []string `json:"email_bcc"`
}

func (x *Xendit) CallbackInvoiceGofiberV2(c *fiber.Ctx) (*CallbackInvoiceResponse, error) {
	callbackToken := c.Get("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		_ = c.Status(fiber.StatusUnauthorized).JSON(map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body CallbackInvoiceResponse
	if err := c.BodyParser(&body); err != nil {
		_ = c.Status(fiber.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

func (x *Xendit) CallbackDisbursementGofiberV2(c *fiber.Ctx) (*CallbackDisbursementResponse, error) {
	callbackToken := c.Get("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		_ = c.Status(fiber.StatusUnauthorized).JSON(map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body CallbackDisbursementResponse
	if err := c.BodyParser(&body); err != nil {
		_ = c.Status(fiber.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

// --- GoFiber V3 ---

func (x *Xendit) CallbackInvoiceGofiberV3(c fiberV3.Ctx) (*CallbackInvoiceResponse, error) {
	callbackToken := c.Get("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		_ = c.Status(http.StatusUnauthorized).JSON(map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body CallbackInvoiceResponse
	if err := c.Bind().JSON(&body); err != nil {
		_ = c.Status(http.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

func (x *Xendit) CallbackDisbursementGofiberV3(c fiberV3.Ctx) (*CallbackDisbursementResponse, error) {
	callbackToken := c.Get("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		_ = c.Status(http.StatusUnauthorized).JSON(map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body CallbackDisbursementResponse
	if err := c.Bind().JSON(&body); err != nil {
		_ = c.Status(http.StatusBadRequest).JSON(map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

// --- Gin ---

func (x *Xendit) CallbackInvoiceGin(c *gin.Context) (*CallbackInvoiceResponse, error) {
	callbackToken := c.GetHeader("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body CallbackInvoiceResponse
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

func (x *Xendit) CallbackDisbursementGin(c *gin.Context) (*CallbackDisbursementResponse, error) {
	callbackToken := c.GetHeader("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body CallbackDisbursementResponse
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

// --- Echo ---

func (x *Xendit) CallbackInvoiceEcho(c echo.Context) (*CallbackInvoiceResponse, error) {
	callbackToken := c.Request().Header.Get("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		_ = c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body CallbackInvoiceResponse
	if err := c.Bind(&body); err != nil {
		_ = c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

func (x *Xendit) CallbackDisbursementEcho(c echo.Context) (*CallbackDisbursementResponse, error) {
	callbackToken := c.Request().Header.Get("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		_ = c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body CallbackDisbursementResponse
	if err := c.Bind(&body); err != nil {
		_ = c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

// --- HTTP Biasa (net/http) ---

func (x *Xendit) CallbackInvoiceHttp(w http.ResponseWriter, r *http.Request) (*CallbackInvoiceResponse, error) {
	callbackToken := r.Header.Get("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body CallbackInvoiceResponse
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

func (x *Xendit) CallbackDisbursementHttp(w http.ResponseWriter, r *http.Request) (*CallbackDisbursementResponse, error) {
	callbackToken := r.Header.Get("X-Callback-Token")
	if callbackToken == "" || callbackToken != x.CallbackToken {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body CallbackDisbursementResponse
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

// --- FastHTTP ---

func (x *Xendit) CallbackInvoiceFasthttp(ctx *fasthttp.RequestCtx) (*CallbackInvoiceResponse, error) {
	callbackToken := string(ctx.Request.Header.Peek("X-Callback-Token"))
	if callbackToken == "" || callbackToken != x.CallbackToken {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		_ = json.NewEncoder(ctx).Encode(map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body CallbackInvoiceResponse
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

func (x *Xendit) CallbackDisbursementFasthttp(ctx *fasthttp.RequestCtx) (*CallbackDisbursementResponse, error) {
	callbackToken := string(ctx.Request.Header.Peek("X-Callback-Token"))
	if callbackToken == "" || callbackToken != x.CallbackToken {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		_ = json.NewEncoder(ctx).Encode(map[string]string{
			"error": "Invalid callback token",
		})
		return nil, fmt.Errorf("invalid callback token")
	}

	var body CallbackDisbursementResponse
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
