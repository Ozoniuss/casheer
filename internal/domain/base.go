package domain

import (
	"errors"
	"strings"
	"time"
)

type BaseModel struct {
	Id        int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Month int

func NewMonth(month int) (Month, error) {
	if month < 1 || month > 12 {
		return 0, ErrInvalidMonthNumber
	}
	return Month(month), nil
}

var ErrInvalidMonthNumber = errors.New("month must be between 1 and 12")

type errorWithUnderlyingError struct {
	underlying []error
}

func (e errorWithUnderlyingError) Error() string {
	errstrs := make([]string, 0, len(e.underlying))
	for _, err := range e.underlying {
		errstrs = append(errstrs, err.Error())
	}
	return strings.Join(errstrs, ";")
}

func (e errorWithUnderlyingError) Unwrap() []error {
	return e.underlying
}
