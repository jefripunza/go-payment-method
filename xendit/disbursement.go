package xendit

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type DisbursementBank struct {
	Name            string `json:"name"`
	Code            string `json:"code"`
	CanDisburse     bool   `json:"can_disburse"`
	CanNameValidate bool   `json:"can_name_validate"`
}

type DisbursementRequestRequest struct {
	ExternalId        string  `json:"external_id"`
	Amount            float64 `json:"amount"`
	BankCode          string  `json:"bank_code"`
	AccountHolderName string  `json:"account_holder_name"`
	AccountNumber     string  `json:"account_number"`
	Description       string  `json:"description"`
}

type DisbursementCreateResponse struct {
	Id                string  `json:"id"`
	UserId            string  `json:"user_id"`
	ExternalId        string  `json:"external_id"`
	Amount            float64 `json:"amount"`
	BankCode          string  `json:"bank_code"`
	AccountHolderName string  `json:"account_holder_name"`
	Description       string  `json:"description"`
	Status            string  `json:"status"`
}

func (x *Xendit) GetDisbursementBanks(forUserId ...string) ([]DisbursementBank, error) {
	url := fmt.Sprintf("%s/available_disbursements_banks", x.BaseUrl)

	resp, _, err := x.doRequest("GET", url, "", "", nil, forUserId...)
	if err != nil {
		return nil, err
	}

	// Parse response
	banks := make([]DisbursementBank, 0)
	if err := json.Unmarshal(resp, &banks); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return banks, nil
}

func (x *Xendit) DisbursementCreate(bankCode string, bankHolderName string, bankAccountNumber string, amount float64, description string, forUserId ...string) (*DisbursementCreateResponse, error) {
	url := fmt.Sprintf("%s/disbursements", x.BaseUrl)

	eid_code := strings.ReplaceAll(bankCode, "ID_", "")
	eid_code = strings.ToUpper(eid_code)

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	req := DisbursementRequestRequest{
		ExternalId:        fmt.Sprintf("withdraw-%s-%s", eid_code, timestamp),
		Amount:            amount,
		BankCode:          bankCode,
		AccountHolderName: bankHolderName,
		AccountNumber:     bankAccountNumber,
		Description:       description,
	}

	// Convert request to JSON
	req_json, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, _, err := x.doRequest("POST", url, "", "", req_json, forUserId...)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result DisbursementCreateResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}
