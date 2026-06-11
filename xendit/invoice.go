package xendit

import (
	"encoding/json"
	"fmt"
	"strings"
)

type InvoiceCreateRequest struct {
	ExternalId     string          `json:"external_id"`
	Amount         float64         `json:"amount"`
	PayerEmail     string          `json:"payer_email"`
	PaymentMethods []string        `json:"payment_methods"` // ["CREDIT_CARD", "BCA", "BNI", BSI, "BRI", "MANDIRI", "PERMATA", "SAHABAT_SAMPOERNA", "BNC", "ALFAMART", "INDOMARET", "OVO", "DANA", "SHOPEEPAY", "LINKAJA", "JENIUSPAY", "DD_BRI", "DD_BCA_KLIKPAY", "KREDIVO", "AKULAKU", "ATOME", "QRIS"]
	Customer       InvoiceCustomer `json:"customer"`
	Items          []InvoiceItem   `json:"items"`
	Description    string          `json:"description"`
}
type InvoiceCustomer struct {
	GivenNames string `json:"given_names"`
	Email      string `json:"email"`
}
type InvoiceItem struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type InvoiceData struct {
	ID                      string                `json:"id"`
	ExternalId              string                `json:"external_id"`
	UserId                  string                `json:"user_id"`
	Status                  string                `json:"status"`
	MerchantName            string                `json:"merchant_name"`
	MerchantProfilePicUrl   string                `json:"merchant_profile_picture_url"`
	Amount                  float64               `json:"amount"`
	PayerEmail              string                `json:"payer_email"`
	Description             string                `json:"description"`
	ExpiryDate              string                `json:"expiry_date"`
	InvoiceUrl              string                `json:"invoice_url"`
	AvailableBanks          []InvoiceBank         `json:"available_banks"`
	AvailableRetailOutlets  []InvoiceRetailOutlet `json:"available_retail_outlets"`
	AvailableEwallets       []InvoiceEWallet      `json:"available_ewallets"`
	AvailableQRCodes        []InvoiceQRCode       `json:"available_qr_codes"`
	AvailableDirectDebits   []InvoiceDirectDebit  `json:"available_direct_debits"`
	AvailablePaylaters      []InvoicePaylater     `json:"available_paylaters"`
	ShouldExcludeCreditCard bool                  `json:"should_exclude_credit_card"`
	ShouldSendEmail         bool                  `json:"should_send_email"`
	Created                 string                `json:"created"`
	Updated                 string                `json:"updated"`
	Currency                string                `json:"currency"`
	Metadata                interface{}           `json:"metadata"`
	Error                   string                `json:"error,omitempty"`
}
type InvoiceBank struct {
	BankCode          string  `json:"bank_code"`
	CollectionType    string  `json:"collection_type"`
	TransferAmount    float64 `json:"transfer_amount"`
	BankBranch        string  `json:"bank_branch"`
	AccountHolderName string  `json:"account_holder_name"`
	IdentityAmount    float64 `json:"identity_amount"`
}
type InvoiceRetailOutlet struct {
	RetailOutletName string `json:"retail_outlet_name"`
}
type InvoiceEWallet struct {
	EWalletType string `json:"ewallet_type"`
}
type InvoiceQRCode struct {
	QRCodeType string `json:"qr_code_type"`
}
type InvoiceDirectDebit struct {
	DirectDebitType string `json:"direct_debit_type"`
}
type InvoicePaylater struct {
	PaylaterType string `json:"paylater_type"`
}

func (x *Xendit) InvoiceCreate(externalId string, name string, email string, items []InvoiceItem, paymentMethods []string, margin float64, forUserId ...string) (*InvoiceData, error) {
	url := fmt.Sprintf("%s/v2/invoices", x.BaseUrl)

	amount := float64(0)
	for _, item := range items {
		amount += item.Price * float64(item.Quantity)
	}

	fixPaymentMethods := make([]string, len(paymentMethods))
	for i, paymentMethod := range paymentMethods {
		fixPaymentMethods[i] = strings.ToUpper(paymentMethod)
	}

	req := InvoiceCreateRequest{
		ExternalId:     externalId,
		Amount:         amount + margin,
		PayerEmail:     email,
		PaymentMethods: fixPaymentMethods,
		Customer: InvoiceCustomer{
			GivenNames: name,
			Email:      email,
		},
		Description: "Fee Admin: " + fmt.Sprintf("%.2f", margin),
		Items:       items,
	}

	// Convert request to JSON
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, _, err := x.doRequest("POST", url, "", "", payload, forUserId...)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result InvoiceData
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}

func (x *Xendit) GetInvoice(id string, forUserId ...string) (*InvoiceData, error) {
	url := fmt.Sprintf("%s/v2/invoices/%s", x.BaseUrl, id)

	resp, _, err := x.doRequest("GET", url, "", "", nil, forUserId...)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result InvoiceData
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}
