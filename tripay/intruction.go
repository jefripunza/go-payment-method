package tripay

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type InstructionStep struct {
	Title string   `json:"title"`
	Step  []string `json:"step"`
}

type PaymentInstructionResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    []InstructionStep `json:"data"`
}

type PaymentInstructionRequest struct {
	Code      string `json:"code"`
	PayCode   string `json:"pay_code,omitempty"`
	Amount    int    `json:"amount,omitempty"`
	AllowHtml int    `json:"allow_html,omitempty"`
}

// GetPaymentInstruction retrieves payment instructions for a specific channel.
func (t *Tripay) GetPaymentInstruction(req PaymentInstructionRequest) (*PaymentInstructionResponse, error) {
	u, err := url.Parse(t.BaseUrl + "/payment/instruction")
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %v", err)
	}

	q := u.Query()
	q.Set("code", req.Code)
	if req.PayCode != "" {
		q.Set("pay_code", req.PayCode)
	}
	if req.Amount > 0 {
		q.Set("amount", strconv.Itoa(req.Amount))
	}
	if req.AllowHtml == 0 || req.AllowHtml == 1 {
		q.Set("allow_html", strconv.Itoa(req.AllowHtml))
	}

	u.RawQuery = q.Encode()

	respBody, _, err := t.doRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	var result PaymentInstructionResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &result, nil
}
