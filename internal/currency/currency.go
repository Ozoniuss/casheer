package currency

import (
	"fmt"
)

// Value represents a currency value in a certain currency. The actual value
// shall be multiplied by "10eExp" to obtain the value in the base unit of the
// currency.
//
// E.g. a Value of 100 with Exp = -2 and USD currency is the equivalent of 1$.
type Value struct {
	Currency string
	Amount   int
	Exponent int
}

func NewValue(amount int, currency string, exp int) (Value, error) {

	value := Value{
		Amount:   amount,
		Currency: currency,
		Exponent: exp,
	}

	if !isValidCurrency(currency) {
		return Value{}, fmt.Errorf("creating currency value: %w", NewErrInvalidCurrency(currency))
	}
	return value, nil
}

func NewValueBasedOnCurrency(amount int, currency string, exponent *int) (Value, error) {

	// May change in the future, but at the moment the only currencies that
	// are allowed have the least valuable unit two orders of magnitude smaller
	// than the actual unit.
	actualExponent := -2
	if exponent != nil {
		actualExponent = *exponent
	}

	if isValidCurrency(currency) {
		return Value{
			Amount:   amount,
			Currency: currency,
			Exponent: actualExponent,
		}, nil
	}
	return Value{}, fmt.Errorf("creating new value based on currency: %w", NewErrInvalidCurrency(currency))
}

// ISO 4217 compliant currency codes
const (
	EUR = "EUR"
	RON = "RON"
	USD = "USD"
)

func isValidCurrency(currency string) bool {
	return currency == EUR ||
		currency == RON ||
		currency == USD
}
