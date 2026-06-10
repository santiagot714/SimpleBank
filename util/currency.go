package util

// Supported currency codes.
const (
	USD = "USD"
	EUR = "EUR"
	GBP = "GBP"
	CAD = "CAD"
	CHF = "CHF"
	NZD = "NZD"
	AUD = "AUD"
	COP = "COP"
)

// IsSupportedCurrency checks if the currency is supported.
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, GBP, CAD, CHF, NZD, AUD, COP:
		return true
	}
	return false
}
