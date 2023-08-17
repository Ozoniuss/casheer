package testutils

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

// CheckNoContextErrors verifies that there are no errors are attached to the
// context. If stops and fails the test otherwise.
func CheckNoContextErrors(t *testing.T, ctx *gin.Context) {
	if len(ctx.Errors) != 0 {
		t.Fatalf("Expected to have no errors attached to the context, found %d. First error of type %v: %s\n", len(ctx.Errors), reflect.TypeOf(ctx.Errors[0]), ctx.Errors[0].Error())
	}
}

// CheckIsContextErrors verifies that the error attached to the context has the
// specified value. It stops and fails the test otherwise.
func CheckIsContextError(t *testing.T, ctx *gin.Context, target error) {
	if len(ctx.Errors) == 0 {
		t.Fatal("Expected to have errors attached to the context, found none.")
	}
	if !errors.Is(ctx.Errors[0], target) {
		t.Fatalf("Expected error to be of type %s, got %s\n", reflect.TypeOf(target), reflect.TypeOf(ctx.Errors[0]))
	}
}

// CheckCanBeContextError verifies that the error attached to the context can
// be matched to the target error. It stops and fails the test otherwise.
func CheckCanBeContextError(t *testing.T, ctx *gin.Context, target any) {
	if len(ctx.Errors) == 0 {
		t.Fatalf("Expected to have an error attached to the context.")
	}
	ok := errors.As(ctx.Errors[0], target)
	if !ok {
		t.Errorf("Expected error to match type %s, got %s\n: %v", reflect.TypeOf(target), reflect.TypeOf(ctx.Errors[0]), ctx.Errors[0])
	}
}
