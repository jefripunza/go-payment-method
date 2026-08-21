# go-payment-method

Library Go untuk integrasi payment gateway Indonesia: **Tripay**, **Midtrans**, **Xendit** (+ placeholder Duitku).

Semua provider diimplementasi penuh terhadap dokumentasi/OpenAPI resmi masing-masing:
- **Tripay** — 13/13 endpoint resmi (tripay.co.id/developer)
- **Midtrans** — 33/33 endpoint resmi (7 API: Snap, Core, Card Tokenization, GoPay Tokenization, IRIS, Payment Link, Subscription)
- **Xendit** — invoice, balance, disbursement, payout, qris, transaction, dll.

## Install

```bash
go get github.com/jefripunza/go-payment-method
```

## Tripay

```go
import "github.com/jefripunza/go-payment-method/tripay"

t := tripay.NewTripay(false, apiKey, privateKey) // sandbox; true = production
```

### Closed Payment
```go
// Buat signature lalu request transaksi
sig := t.CreateClosedTransactionSignature("T0001", "INV001", 10000)

req := tripay.CreateClosedTransactionRequest{
    Method: "BRIVA", MerchantRef: "INV001", Amount: 10000,
    CustomerName: "Budi", OrderItems: []tripay.OrderItem{{Name: "Produk", Price: 10000, Quantity: 1}},
    Signature: sig,
}
resp, err := t.CreateClosedTransaction(req) // dapat reference + pay_code
detail, err := t.GetClosedTransactionDetail("T0001000000000000001")
status, err := t.CheckClosedTransactionStatus("T0001000000000000001")
```

### Open Payment
```go
sig := t.CreateOpenPaymentSignature("T0001", "BCAVA", "INV002")
resp, err := t.CreateOpenPayment(tripay.CreateOpenPaymentRequest{
    Method: "BCAVA", MerchantRef: "INV002", CustomerName: "Budi", Signature: sig,
})
detail, err := t.GetOpenPaymentDetail(uuid)
txs, err := t.GetOpenPaymentTransactions(uuid)
```

### Merchant & lainnya
```go
channels, _ := t.GetMerchantPaymentChannels()
fees, _ := t.CalculateMerchantFees("BRIVA", 10000)
txs, _ := t.GetMerchantTransactions(tripay.MerchantTransactionsFilter{})
instr, _ := t.GetPaymentInstruction(tripay.PaymentInstructionRequest{Code: "BRIVA"})

// E-Wallet (OVO/DANA)
linkResp, _ := t.LinkEWallet(tripay.EWalletLinkRequest{CustomerName: "Budi", CustomerPhone: "08123456789"})
detailEwallet, _ := t.GetEWalletDetail(reference)
unlinkResp, _ := t.UnlinkEWallet(reference)
```

### Callback
```go
// Verifikasi signature callback (X-Callback-Signature)
ok := t.VerifySignature(rawBody, receivedSignature)

// Handler siap pakai per framework:
// t.CallbackHttp(w, r) | t.CallbackGin(c) | t.CallbackEcho(c)
// t.CallbackGofiberV2(c) | t.CallbackGofiberV3(c) | t.CallbackFasthttp(ctx)
```

## Midtrans

```go
import pm "github.com/jefripunza/go-payment-method"

m := pm.NewMidtrans(serverKey, clientKey, false) // sandbox; true = production
```

### Snap (hosted checkout)
```go
snapResp, err := m.CreateSnapTransaction(pm.SnapRequest{
    TransactionDetails: pm.SnapTransactionDetails{OrderID: "ORDER-1", GrossAmount: 10000},
    ItemDetails: []pm.SnapItemDetail{{ID: "P1", Price: 10000, Quantity: 1, Name: "Produk"}},
    CustomerDetails: &pm.SnapCustomerDetails{FirstName: "Budi", Email: "budi@mail.com"},
})
// snapResp.Token, snapResp.RedirectURL
```

### Core API (charge & transaksi)
```go
chargeResp, err := m.Charge(pm.CoreChargeRequest{
    PaymentType: "bank_transfer",
    TransactionDetails: pm.CoreTransactionDetails{OrderID: "ORDER-1", GrossAmount: 10000},
    AdditionalFields: map[string]interface{}{
        "bank_transfer": map[string]interface{}{"bank": "bca"},
    },
})
status, err := m.GetTransactionStatus("ORDER-1")
approved, err := m.Approve("ORDER-1")
denied, err := m.Deny("ORDER-1")
cancelled, err := m.Cancel("ORDER-1")
expired, err := m.Expire("ORDER-1")
refunded, err := m.Refund("ORDER-1", pm.RefundRequest{RefundKey: "k1", Amount: 10000, Reason: "retur"})
refundDirect, err := m.RefundOnlineDirect("ORDER-1", 10000)
```

### IRIS Disbursement
```go
beneficiaries, _ := m.IRISListBeneficiaries()
created, _ := m.IRISCreateBeneficiary(pm.IRISBeneficiary{Name: "Budi", Account: "123456", Bank: "bca", AliasName: "budi"})
payouts, _ := m.IRISCreatePayouts(pm.IRISPayoutRequest{Payouts: []pm.IRISPayoutItem{{
    BeneficiaryName: "Budi", BeneficiaryAccount: "123456", BeneficiaryBank: "bca", Amount: "10000"}}})
approved, _ := m.IRISApprovePayouts(pm.IRISPayoutApproval{ReferenceNos: []string{"R1"}, OTP: "123456"})
rejected, _ := m.IRISRejectPayouts(pm.IRISPayoutRejection{ReferenceNos: []string{"R1"}, RejectReason: "duplicate"})
balance, _ := m.IRISGetBalance()
banks, _ := m.IRISListBankAccounts()
valid, _ := m.IRISValidateAccount("bca", "123456")
payoutStatus, _ := m.IRISGetPayoutStatus("reference-no")
ping, _ := m.IRISPing()
```

### Payment Link
```go
link, _ := m.CreatePaymentLink(pm.PaymentLinkRequest{
    TransactionDetails: pm.PaymentLinkTransactionDetails{OrderID: "PL-1", GrossAmount: 50000},
})
detail, _ := m.GetPaymentLink("PL-1")
err := m.DeletePaymentLink("PL-1")
```

### Subscription
```go
sub, _ := m.CreateSubscription(pm.SubscriptionRequest{
    Name: "Langganan", Amount: "10000", Currency: "IDR", PaymentType: "credit_card",
    Token: "card-token", Schedule: &pm.SubscriptionSchedule{Interval: 1, IntervalUnit: "month"},
})
detail, _ := m.GetSubscription("sub-id")
updated, _ := m.UpdateSubscription("sub-id", map[string]interface{}{"amount": "20000"})
enabled, _ := m.EnableSubscription("sub-id")
disabled, _ := m.DisableSubscription("sub-id")
cancelled, _ := m.CancelSubscription("sub-id")
```

### Card & GoPay Tokenization
```go
cardToken, _ := m.RegisterCardToken("4811111111111114", 12, 2025, "")
gopayToken, _ := m.CreateGoPayAccountToken(pm.CreateGoPayTokenRequest{
    PaymentType: "gopay", Gopay: map[string]interface{}{"enable_callback": true},
})
gopayStatus, _ := m.GetGoPayAccountStatus("account-id")
```

### Verifikasi notifikasi webhook
```go
ok := m.VerifyNotificationSignature(signatureKey, orderID, statusCode, grossAmount)
```

## Xendit

```go
import "github.com/jefripunza/go-payment-method/xendit"

x := xendit.NewXendit(apiKey, callbackToken)
invoice, err := x.CreateInvoice(xendit.CreateInvoiceRequest{...})
```

Lihat `example/` untuk contoh lengkap.

## Test

```bash
go test ./...
```
