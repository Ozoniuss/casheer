package currency

// ISO 4217 compliant currency codes
const (
	EUR = "EUR"
	RON = "RON"
	USD = "USD"
	BTC = "BTC"
)

var validCurrencies = []string{EUR, RON, USD, BTC}

// Honestly it's not really worth doing anything better, I won't be having
// many currencies. This is just for error purposes anyway.
var validCurrenciesString = "[EUR, RON, USD, BTC]"
