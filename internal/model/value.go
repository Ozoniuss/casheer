package model

import "github.com/Ozoniuss/casheer/currency"

type Value struct {
	Currency string
	Amount   int
	Exponent int
}

// FromCurrencyValue converts the public currency Value type to the model
// representation of value.
func FromCurrencyValue(v currency.Value) Value {
	return Value{
		Currency: v.Currency,
		Amount:   v.Amount,
		Exponent: v.Exponent,
	}
}
