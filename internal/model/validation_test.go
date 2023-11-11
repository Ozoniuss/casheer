package model

import (
	"errors"
	"reflect"
	"testing"

	cmp "github.com/google/go-cmp/cmp"
)

func TestErrInvalidModel(t *testing.T) {
	tests := []struct {
		name           string
		err            ErrInvalidModel
		expectedOutput string
	}{
		{
			name: "reasons array is nil",
			err: ErrInvalidModel{
				modelName: "prefix",
				reasons:   nil,
			},
			expectedOutput: "unexpected error",
		},
		{
			name: "reasons array is empty",
			err: ErrInvalidModel{
				modelName: "prefix",
				reasons:   []string{},
			},
			expectedOutput: "unexpected error",
		},
		{
			name: "one reason",
			err: ErrInvalidModel{
				modelName: "prefix",
				reasons:   []string{"reason1"},
			},
			expectedOutput: "prefix: reason1",
		},
		{
			name: "two reasons",
			err: ErrInvalidModel{
				modelName: "prefix",
				reasons:   []string{"reason1", "reason2"},
			},
			expectedOutput: "prefix: reason1; reason2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.expectedOutput, tt.err.Error()); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestInvalidModelErrBuilder(t *testing.T) {
	t.Run("A builder with no reasons returns a nil error", func(t *testing.T) {
		b := NewBaseModelErrorBuilder("model")
		err := b.Error()
		if err != nil {
			t.Fatalf("expected nil error, got %s", reflect.TypeOf(err))
		}
	})
	t.Run("A builder with no reasons returns a InvalidModel errol", func(t *testing.T) {
		b := NewBaseModelErrorBuilder("model")
		b.AddError("ugly model")
		err := b.Error()
		if err == nil {
			t.Fatalf("got nil error")
		}
		var errInvalidModel ErrInvalidModel
		if !errors.As(err, &errInvalidModel) {
			t.Fatalf("builder return invalid error type; got %s, expected %s", reflect.TypeOf(err), reflect.TypeOf(errInvalidModel))
		}
	})
}
