package currency

import "fmt"

type ErrInvalidCurrency struct {
	attemptedCurrency string
}

func (e ErrInvalidCurrency) Error() string {
	return fmt.Sprintf("currency %s does not exist, must be one of %s", e.attemptedCurrency, validCurrenciesString)
}

func NewErrInvalidCurrency(attemptedCurrency string) ErrInvalidCurrency {
	return ErrInvalidCurrency{
		attemptedCurrency: attemptedCurrency,
	}
}
