package currency

import (
	"fmt"
)

// Value represents a currency value in a certain currency. The actual value
// shall be multiplied by "10eExp" to obtain the value in the base unit of the
// currency.
//
// Value should only be created via the provided constructors, to prevent
// invalid "values". Do not modify "Value" directly.
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

	if err := validateCurrency(currency); err != nil {
		return Value{}, fmt.Errorf("creating value: %w", err)
	}
	return value, nil
}

func NewValueBasedOnMinorCurrency(amount int, currency string, exponent *int) (Value, error) {

	// May change in the future, but at the moment the only currencies that
	// are allowed have the least valuable unit two orders of magnitude smaller
	// than the actual unit.
	actualExponent := -2
	if exponent != nil {
		actualExponent = *exponent
	}

	if err := validateCurrency(currency); err != nil {
		return Value{}, fmt.Errorf("creating value: %w", NewErrInvalidCurrency(currency))
	}

	return Value{
		Amount:   amount,
		Currency: currency,
		Exponent: actualExponent,
	}, nil
}

func NewUSDValue(amount int) Value {
	return Value{
		Currency: USD,
		Amount:   amount,
		Exponent: -2,
	}
}
func NewEURValue(amount int) Value {
	return Value{
		Currency: EUR,
		Amount:   amount,
		Exponent: -2,
	}
}
func NewRONValue(amount int) Value {
	return Value{
		Currency: RON,
		Amount:   amount,
		Exponent: -2,
	}
}

// validateCurrency provides a way to test whether or not a string can represent
// a valid currency.
func validateCurrency(currency string) error {
	for _, c := range validCurrencies {
		if currency == c {
			return nil
		}
	}
	return NewErrInvalidCurrency(currency)
}
