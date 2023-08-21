package currency

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
