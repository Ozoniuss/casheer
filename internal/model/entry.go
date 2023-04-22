package model

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Entry models an entry of a planning. Fields are automatically mapped by gorm
// to their database equivalents.
//
// TODO: custom validation message.
type Entry struct {
	BaseModel

	// Postgresql doesn't support unsigned int.
	Month         int8   `validate:"required,gte=1,lte=12"`
	Year          int16  `validate:"required,gte=2020"`
	Category      string `validate:"required" json:"category"`
	Subcategory   string `validate:"required"`
	ExpectedTotal float32
	RunningTotal  float32
	Recurring     bool

	Expenses []Expense
}

type InvalidEntryErr struct {
}

func (e InvalidEntryErr) Error() string {
	return "invalid entry"
}

// ValidEntry can be used as a scope to validate an entry before inserting it
// into the database.
func ValidEntry(e Entry) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		v := validator.New()
		if err := v.Struct(e); err != nil {
			db.AddError(InvalidEntryErr{})
		}
		return db
	}
}

// ValidEntryFields works like ValidEntry, but can be used with a specific
// subset of the entry's fields.
func ValidEntryFields(e Entry, fields []string) func(db *gorm.DB) *gorm.DB {
	// TODO: find a nicer way to update these fields
	// fiels are specified lowercase to match the table names, but in order to
	// match struct names they must be uppercased.

	// Use one single builder to avoid unnecessary allocations.
	b := strings.Builder{}
	for i := 0; i < len(fields); i++ {
		fields[i] = capitalize(fields[i], b)
	}

	return func(db *gorm.DB) *gorm.DB {
		v := validator.New()
		if err := v.StructPartial(e, fields...); err != nil {
			db.AddError(InvalidEntryErr{})
		}
		return db
	}
}

// AfterUpdate is a gorm hook that adds an error if the entry was not found
// during an update operation. This implicitly assumes that the update query
// executes with a "returning" clause that writes to an empty entry.
func (e *Entry) AfterUpdate(tx *gorm.DB) (err error) {
	if e.Id == uuid.Nil {
		err = gorm.ErrRecordNotFound
	}
	return
}

// AfterDelete is a gorm hook that adds an error if the entry was not found
// during an delete operation. This implicitly assumes that the delete query
// executes with a "returning" clause that writes to an empty entry.
func (e *Entry) AfterDelete(tx *gorm.DB) (err error) {
	if e.Id == uuid.Nil {
		err = gorm.ErrRecordNotFound
	}
	return
}

func capitalize(s string, b strings.Builder) string {
	// Builder is reused across multiple calls.
	defer b.Reset()
	b.WriteByte(s[0] - 32)
	for i := 1; i < len(s); i++ {
		b.WriteByte(s[i])
	}
	return b.String()
}
