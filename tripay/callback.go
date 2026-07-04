package tripay

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	fiberV2 "github.com/gofiber/fiber/v2"
	fiberV3 "github.com/gofiber/fiber/v3"
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
)

type CallbackResponse struct {
	Reference         string      `json:"reference"`
	MerchantRef       string      `json:"merchant_ref"`
	PaymentMethod     string      `json:"payment_method"`
	PaymentMethodCode string      `json:"payment_method_code"`
	TotalAmount       float64     `json:"total_amount"`
	FeeMerchant       float64     `json:"fee_merchant"`
	AmountReceived    float64     `json:"amount_received"`
	Status            string      `json:"status"`
	PaidAt            interface{} `json:"paid_at"`
}

func (t *Tripay) VerifySignature(rawBody []byte, receivedSignature string) bool {
	h := hmac.New(sha256.New, []byte(t.PrivateKey))
	h.Write(rawBody)
	expectedSignature := hex.EncodeToString(h.Sum(nil))
	return expectedSignature == receivedSignature
}

// --- GoFiber V2 ---

func (t *Tripay) CallbackGofiberV2(c *fiberV2.Ctx) (*CallbackResponse, error) {
	signature := c.Get("X-Callback-Signature")
	bodyBytes := c.Body()

	if signature == "" || !t.VerifySignature(bodyBytes, signature) {
		_ = c.Status(fiberV2.StatusUnauthorized).JSON(map[string]interface{}{
			"success": false,
			"message": "Invalid signature",
		})
		return nil, fmt.Errorf("invalid signature")
	}

	var body CallbackResponse
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		_ = c.Status(fiberV2.StatusBadRequest).JSON(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

// --- GoFiber V3 ---

func (t *Tripay) CallbackGofiberV3(c fiberV3.Ctx) (*CallbackResponse, error) {
	signature := c.Get("X-Callback-Signature")
	bodyBytes := c.Body()

	if signature == "" || !t.VerifySignature(bodyBytes, signature) {
		_ = c.Status(http.StatusUnauthorized).JSON(map[string]interface{}{
			"success": false,
			"message": "Invalid signature",
		})
		return nil, fmt.Errorf("invalid signature")
	}

	var body CallbackResponse
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		_ = c.Status(http.StatusBadRequest).JSON(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

// --- Gin ---

func (t *Tripay) CallbackGin(c *gin.Context) (*CallbackResponse, error) {
	signature := c.GetHeader("X-Callback-Signature")
	if signature == "" {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Invalid signature",
		})
		return nil, fmt.Errorf("invalid signature")
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return nil, err
	}
	// Restore body
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if !t.VerifySignature(bodyBytes, signature) {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Invalid signature",
		})
		return nil, fmt.Errorf("invalid signature")
	}

	var body CallbackResponse
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

// --- Echo ---

func (t *Tripay) CallbackEcho(c echo.Context) (*CallbackResponse, error) {
	signature := c.Request().Header.Get("X-Callback-Signature")
	if signature == "" {
		_ = c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Invalid signature",
		})
		return nil, fmt.Errorf("invalid signature")
	}

	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		_ = c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return nil, err
	}
	// Restore body
	c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if !t.VerifySignature(bodyBytes, signature) {
		_ = c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Invalid signature",
		})
		return nil, fmt.Errorf("invalid signature")
	}

	var body CallbackResponse
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		_ = c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

// --- HTTP Biasa (net/http) ---

func (t *Tripay) CallbackHttp(w http.ResponseWriter, r *http.Request) (*CallbackResponse, error) {
	signature := r.Header.Get("X-Callback-Signature")
	if signature == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid signature",
		})
		return nil, fmt.Errorf("invalid signature")
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return nil, err
	}
	// Restore body
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if !t.VerifySignature(bodyBytes, signature) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid signature",
		})
		return nil, fmt.Errorf("invalid signature")
	}

	var body CallbackResponse
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}

// --- FastHTTP ---

func (t *Tripay) CallbackFasthttp(ctx *fasthttp.RequestCtx) (*CallbackResponse, error) {
	signature := string(ctx.Request.Header.Peek("X-Callback-Signature"))
	bodyBytes := ctx.PostBody()

	if signature == "" || !t.VerifySignature(bodyBytes, signature) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetContentType("application/json")
		_ = json.NewEncoder(ctx).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid signature",
		})
		return nil, fmt.Errorf("invalid signature")
	}

	var body CallbackResponse
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		_ = json.NewEncoder(ctx).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return nil, err
	}

	return &body, nil
}
