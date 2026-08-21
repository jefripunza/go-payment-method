package xendit

import (
	"encoding/json"
	"fmt"
)

// ==========================================
// PAYOUTS V3 MODELS — money-out (domestic + cross-border)
// ==========================================
// Referensi openapi: /v3/payouts (api-version: 2025-09-01).

// PayoutV3RecipientAddress adalah alamat penerima payout.
type PayoutV3RecipientAddress struct {
	Country       string `json:"country,omitempty"`
	StreetLine1   string `json:"street_line_1,omitempty"`
	City          string `json:"city,omitempty"`
	ProvinceState string `json:"province_state,omitempty"`
	PostalCode    string `json:"postal_code,omitempty"`
}

// PayoutV3RecipientDetails adalah info tambahan penerima (opsional).
type PayoutV3RecipientDetails struct {
	PersonalMobileNumber string `json:"personal_mobile_number,omitempty"`
	Nationality          string `json:"nationality,omitempty"`
	// Field lain bisa ditambahkan sesuai kebutuhan (document_number, dsb).
}

// PayoutV3AccountDetails adalah rekening tujuan payout.
type PayoutV3AccountDetails struct {
	Currency          string `json:"currency"`
	AccountCountry    string `json:"account_country"`
	AccountHolderName string `json:"account_holder_name"`
	AccountNumber     string `json:"account_number"`
	RoutingType1      string `json:"routing_type_1,omitempty"`  // WALLET, ABA, IBAN, dsb
	RoutingValue1     string `json:"routing_value_1,omitempty"` // PH_GCASH, 021000021, dsb
	RoutingType2      string `json:"routing_type_2,omitempty"`
	RoutingValue2     string `json:"routing_value_2,omitempty"`
}

// PayoutV3Recipient adalah penerima dana payout.
type PayoutV3Recipient struct {
	Type           string                    `json:"type"` // INDIVIDUAL | BUSINESS
	GivenName      string                    `json:"given_name,omitempty"`
	Surname        string                    `json:"surname,omitempty"`
	BusinessName   string                    `json:"business_name,omitempty"`
	Relationship   string                    `json:"relationship"` // SELF | EMPLOYEE | CUSTOMER | CONTRACTOR | SUPPLIER
	Details        *PayoutV3RecipientDetails `json:"details,omitempty"`
	Address        *PayoutV3RecipientAddress `json:"address,omitempty"`
	AccountDetails *PayoutV3AccountDetails   `json:"account_details"`
}

// PayoutV3PayoutDetails adalah rincian nominal payout.
type PayoutV3PayoutDetails struct {
	SourceCurrency      string `json:"source_currency"`
	SourceAmount        string `json:"source_amount"`
	DestinationCurrency string `json:"destination_currency"`
	DestinationAmount   string `json:"destination_amount,omitempty"`
}

// PayoutV3ReceiptNotification adalah notifikasi email untuk payout.
type PayoutV3ReceiptNotification struct {
	EmailTo []string `json:"email_to,omitempty"`
	EmailCc []string `json:"email_cc,omitempty"`
}

// PayoutV3UnderlyingDocument adalah lampiran dokumen payout.
type PayoutV3UnderlyingDocument struct {
	Type string `json:"type,omitempty"` // INVOICE | CONTRACT | OTHER
	URL  string `json:"url,omitempty"`
}

// PayoutV3CreateRequest adalah body untuk POST /v3/payouts.
type PayoutV3CreateRequest struct {
	ReferenceID         string                       `json:"reference_id"`
	Recipient           *PayoutV3Recipient           `json:"recipient"`
	PayoutDetails       *PayoutV3PayoutDetails       `json:"payout_details"`
	SourceOfFund        string                       `json:"source_of_fund"` // BUSINESS_REVENUE | LOAN | INVESTMENT | OTHERS
	PurposeCode         string                       `json:"purpose_code"`   // SALARY | SHIPPING | ...
	Description         string                       `json:"description,omitempty"`
	ReceiptNotification *PayoutV3ReceiptNotification `json:"receipt_notification,omitempty"`
	Metadata            map[string]interface{}       `json:"metadata,omitempty"`
	UnderlyingDocuments []PayoutV3UnderlyingDocument `json:"underlying_documents,omitempty"`
}

// PayoutV3 adalah response payout (GET /v3/payouts/{payout_id}).
type PayoutV3 struct {
	PayoutID            string                 `json:"payout_id,omitempty"`
	Status              string                 `json:"status,omitempty"`
	ReferenceID         string                 `json:"reference_id,omitempty"`
	ProcessorReference  string                 `json:"processor_reference,omitempty"`
	Type                string                 `json:"type,omitempty"` // B2B | B2C | C2C | C2B
	SourceCurrency      string                 `json:"source_currency,omitempty"`
	SourceAmount        string                 `json:"source_amount,omitempty"`
	DestinationCurrency string                 `json:"destination_currency,omitempty"`
	DestinationAmount   string                 `json:"destination_amount,omitempty"`
	Description         string                 `json:"description,omitempty"`
	FailureCode         string                 `json:"failure_code,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
	Created             string                 `json:"created,omitempty"`
	Updated             string                 `json:"updated,omitempty"`
}

// PayoutV3Create membuat payout v3 (POST /v3/payouts).
// Wajib header api-version: 2025-09-01 + idempotency-key (cegah double payout).
func (x *Xendit) PayoutV3Create(idempotencyKey string, req *PayoutV3CreateRequest, forUserId ...string) (*PayoutV3, error) {
	url := fmt.Sprintf("%s/v3/payouts", x.BaseUrl)
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)
	if idempotencyKey != "" {
		headers["Idempotency-key"] = idempotencyKey
	}

	resp, _, err := x.doRequestWithHeaders("POST", url, "2025-09-01", "", headers, payload, forUserId...)
	if err != nil {
		return nil, err
	}

	var result PayoutV3
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PayoutV3Get mengambil status payout v3 (GET /v3/payouts/{payout_id}).
func (x *Xendit) PayoutV3Get(payoutID string, forUserId ...string) (*PayoutV3, error) {
	url := fmt.Sprintf("%s/v3/payouts/%s", x.BaseUrl, payoutID)
	resp, _, err := x.doRequest("GET", url, "2025-09-01", "", nil, forUserId...)
	if err != nil {
		return nil, err
	}

	var result PayoutV3
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
