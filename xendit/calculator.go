package xendit

import (
	"fmt"
	"math"
)

// CalculateMarginPrice menghitung berapa margin yang perlu
// ditambahkan ke harga agar setelah dipotong fee Xendit,
// merchant tetap menerima targetNetAmount.
func CalculateMarginPrice(
	targetNetAmount float64,
	paymentMethod string,
	vatRates ...float64,
) (float64, error) {

	vatRate := 0.11
	if len(vatRates) > 0 {
		vatRate = vatRates[0]
	}

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
		gross = targetNetAmount + (4000.0 * (1 + vatRate))

	// =========================
	// Over The Counter
	// =========================
	case PaymentMethodAlfamart:
		gross = targetNetAmount + (5000.0 * (1 + vatRate))

	case PaymentMethodIndomaret:
		gross = targetNetAmount + (5500.0 * (1 + vatRate))

	// =========================
	// E-Wallet
	// =========================
	case PaymentMethodOVO,
		PaymentMethodDana,
		PaymentMethodLinkAja:
		rate := 0.015 * (1 + vatRate)
		gross = targetNetAmount / (1 - rate)

	case PaymentMethodJeniusPay:
		rate := 0.020 * (1 + vatRate)
		gross = targetNetAmount / (1 - rate)

	case PaymentMethodShopeePay:
		gross = targetNetAmount / (1 - 0.02)

	case PaymentMethodDD_BRI,
		PaymentMethodDD_BCA_KlikPay:
		rate := 0.019 * (1 + vatRate)
		gross = targetNetAmount / (1 - rate)

	case PaymentMethodCreditCard:
		percentRate := 0.029 * (1 + vatRate)
		fixedFee := 2000.0 * (1 + vatRate)
		gross = (targetNetAmount + fixedFee) / (1 - percentRate)

	case PaymentMethodAkulaku:
		rate := 0.017 * (1 + vatRate)
		gross = targetNetAmount / (1 - rate)

	case PaymentMethodKredivo:
		rate := 0.023 * (1 + vatRate)
		gross = targetNetAmount / (1 - rate)

	case PaymentMethodAtome:
		rate := 0.050 * (1 + vatRate)
		gross = targetNetAmount / (1 - rate)

	case PaymentMethodQRIS:
		gross = targetNetAmount / (1 - 0.007)

	default:
		return 0, fmt.Errorf("unsupported payment method: %s", paymentMethod)
	}

	grossAmount := math.Ceil(gross)
	margin := grossAmount - targetNetAmount
	return margin, nil
}
