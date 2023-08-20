package currency

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Value represents a currency value in a certain currency. The actual value
// shall be multiplied by "10eExp" to obtain the value in the base unit of the
// currency.
//
// E.g. a Value of 100 with Exp = -2 and USD currency is the equivalent of 1$.
type Value struct {
	Currency string `validate:"required,iso4217"`
	Amount   int
	Exponent int
}

func NewValue(amount int, code string, exp int) (Value, error) {

	value := Value{
		Amount:   amount,
		Currency: code,
		Exponent: exp,
	}

	validator := validator.New()
	err := validator.Struct(value)
	if err != nil {
		return Value{}, fmt.Errorf("creating currency value: %s", err.Error())
	}
	return value, nil
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

// ISO 4217 compliant currency codes
const (
	EUR = "EUR"
	RON = "RON"
	USD = "USD"
)
