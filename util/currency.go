package util

// constants for supported currencies
const (
	USD = "USD"
	CAD = "CAD"
	KES = "KES"
	EUR = "EUR"
)

// IsSupportedCurrency checks whether a currency is supported or not
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, CAD, KES, EUR:
		return true
	}
	return false
}
