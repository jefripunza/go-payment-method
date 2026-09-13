package midtrans

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// ===========================================================================
// Utilities
// ===========================================================================

// apiError builds a descriptive error from a Midtrans error response body.
func (m *Midtrans) apiError(body []byte, status int) error {
	var parsed struct {
		StatusMessage   string `json:"status_message"`
		StatusMessageID string `json:"status_message_"`
	}
	_ = json.Unmarshal(body, &parsed)
	msg := parsed.StatusMessage
	if msg == "" {
		msg = string(body)
	}
	return fmt.Errorf("midtrans API error (status: %d): %s", status, msg)
}

// apiErrorNoAuth sama seperti apiError, untuk endpoint frontend (tanpa auth).
func (m *Midtrans) apiErrorNoAuth(body []byte, status int) error {
	return m.apiError(body, status)
}

// VerifyNotificationSignature memverifikasi signature_key webhook Midtrans.
// Sesuai openapi: SHA512(orderId + statusCode + grossAmount + ServerKey).
func (m *Midtrans) VerifyNotificationSignature(signatureKey, orderID, statusCode, grossAmount string) bool {
	h := sha512.New()
	h.Write([]byte(orderID + statusCode + grossAmount + m.ServerKey))
	expected := hex.EncodeToString(h.Sum(nil))
	return expected == signatureKey
}
