package xendit

import (
	"encoding/json"
	"fmt"
)

type TransactionLinkData struct {
	Href   string `json:"href"`
	Method string `json:"method"`
	Rel    string `json:"rel"`
}
type TransactionGetResponse struct {
	HasMore bool                  `json:"has_more"`
	Data    []TransactionData     `json:"data"`
	Links   []TransactionLinkData `json:"links"`
}
type TransactionFee struct {
	XenditFee                int    `json:"xendit_fee"`
	ValueAddedTax            int    `json:"value_added_tax"`
	XenditWithholdingTax     int    `json:"xendit_withholding_tax"`
	ThirdPartyWithholdingTax int    `json:"third_party_withholding_tax"`
	Status                   string `json:"status"`
}
type TransactionData struct {
	Id                      string         `json:"id"`
	ProductId               string         `json:"product_id"`
	Type                    string         `json:"type"`
	Status                  string         `json:"status"`
	ChannelCategory         string         `json:"channel_category"`
	ChannelCode             string         `json:"channel_code"`
	ReferenceId             string         `json:"reference_id"`
	AccountIdentifier       string         `json:"account_identifier"`
	Currency                string         `json:"currency"`
	Amount                  float64        `json:"amount"`
	NetAmount               float64        `json:"net_amount"`
	Cashflow                string         `json:"cashflow"`
	SettlementStatus        string         `json:"settlement_status"`
	EstimatedSettlementTime string         `json:"estimated_settlement_time"`
	BusinessId              string         `json:"business_id"`
	Created                 string         `json:"created"`
	Updated                 string         `json:"updated"`
	Fee                     TransactionFee `json:"fee"`
}

func (x *Xendit) GetTransaction(forUserId ...string) ([]TransactionData, error) {
	url := fmt.Sprintf("%s/transactions", x.BaseUrl)

	resp, _, err := x.doRequest("GET", url, "", "", nil, forUserId...)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result TransactionGetResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.Data, nil
}
