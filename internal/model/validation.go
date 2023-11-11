package model

import "strings"

// ErrInvalidModel stores the data of a model validation error. The reasons
// are included to provide all issues that occured during validation.
type ErrInvalidModel struct {
	modelName string
	reasons   []string
}

func NewBaseModelError(prefix string, reasons []string) ErrInvalidModel {
	return ErrInvalidModel{
		modelName: prefix,
		reasons:   reasons,
	}
}

func (b ErrInvalidModel) Error() string {
	if len(b.reasons) == 0 {
		return "unexpected error"
	}
	sb := &strings.Builder{}
	sb.WriteString(b.modelName)
	sb.WriteString(": ")
	for _, reason := range b.reasons {
		sb.WriteString(reason)
		sb.WriteString("; ")
	}
	out := sb.String()
	return out[:len(out)-2]
}

// InvalidModelErrBuilder can be used to create detailed model validation
// errors, by adding each individual error after instantiating a builder.
type InvalidModelErrBuilder struct {
	err ErrInvalidModel
}

func NewBaseModelErrorBuilder(modelName string) *InvalidModelErrBuilder {
	return &InvalidModelErrBuilder{
		err: ErrInvalidModel{
			modelName: "invalid " + modelName,
		},
	}
}

func (b *InvalidModelErrBuilder) AddError(err string) {
	b.err.reasons = append(b.err.reasons, err)
}

func (b *InvalidModelErrBuilder) Error() error {
	if len(b.err.reasons) == 0 {
		return nil
	}
	return b.err
}
