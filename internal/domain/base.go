package domain

import (
	"errors"
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
