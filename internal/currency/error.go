package currency

import "fmt"

type ErrInvalidCurrency struct {
	attemptedCurrency string
}

func (e ErrInvalidCurrency) Error() string {
	return fmt.Sprintf("currency %s does not exist or is not supported", e.attemptedCurrency)
}

func NewErrInvalidCurrency(attemptedCurrency string) ErrInvalidCurrency {
	return ErrInvalidCurrency{
		attemptedCurrency: attemptedCurrency,
	}
}
