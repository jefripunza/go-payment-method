package xendit

import (
	"fmt"
	"math"
)

const (
	vatRate = 0.11
)

// CalculateMarginPrice menghitung berapa margin yang perlu
// ditambahkan ke harga agar setelah dipotong fee Xendit,
// merchant tetap menerima targetNetAmount.
func CalculateMarginPrice(
	targetNetAmount int64,
	paymentMethod string,
) (int64, error) {

	net := float64(targetNetAmount)

	var gross float64

	switch paymentMethod {

	// =========================
	// Virtual Account
	// =========================
	case PaymentMethodBCA,
		PaymentMethodBNI,
		PaymentMethodBSI,
		PaymentMethodBRI,
		PaymentMethodMandiri,
		PaymentMethodPermata,
		PaymentMethodSahabatSampoerna,
		PaymentMethodBNC:

		gross = net + (4000.0 * (1 + vatRate))

	// =========================
	// Over The Counter
	// =========================
	case PaymentMethodAlfamart:

		gross = net + (5000.0 * (1 + vatRate))

	case PaymentMethodIndomaret:

		gross = net + (5500.0 * (1 + vatRate))

	// =========================
	// E-Wallet
	// =========================
	case PaymentMethodOVO,
		PaymentMethodDana,
		PaymentMethodLinkAja:

		rate := 0.015 * (1 + vatRate)
		gross = net / (1 - rate)

	case PaymentMethodJeniusPay:

		rate := 0.020 * (1 + vatRate)
		gross = net / (1 - rate)

	case PaymentMethodShopeePay:

		gross = net / (1 - 0.02)

	// =========================
	// Direct Debit
	// =========================
	case PaymentMethodDD_BRI,
		PaymentMethodDD_BCA_KlikPay:

		rate := 0.019 * (1 + vatRate)
		gross = net / (1 - rate)

	// =========================
	// Credit Card
	// =========================
	case PaymentMethodCreditCard:

		percentRate := 0.029 * (1 + vatRate)
		fixedFee := 2000.0 * (1 + vatRate)

		gross = (net + fixedFee) / (1 - percentRate)

	// =========================
	// PayLater
	// =========================
	case PaymentMethodAkulaku:

		rate := 0.017 * (1 + vatRate)
		gross = net / (1 - rate)

	case PaymentMethodKredivo:

		rate := 0.023 * (1 + vatRate)
		gross = net / (1 - rate)

	case PaymentMethodAtome:

		rate := 0.050 * (1 + vatRate)
		gross = net / (1 - rate)

	// =========================
	// QRIS
	// =========================
	case PaymentMethodQRIS:

		gross = net / (1 - 0.007)

	default:
		return 0, fmt.Errorf("unsupported payment method: %s", paymentMethod)
	}

	grossAmount := int64(math.Ceil(gross))
	margin := grossAmount - targetNetAmount

	return margin, nil
}
