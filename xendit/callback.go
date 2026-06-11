package xendit

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
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

func (x *Xendit) CallbackInvoiceGofiberV2(c *fiber.Ctx) error {
	fmt.Println("Callback hit")
	return nil
}

func (x *Xendit) CallbackDisbursementGofiberV2(c *fiber.Ctx) error {
	fmt.Println("Callback hit")
	return nil
}
