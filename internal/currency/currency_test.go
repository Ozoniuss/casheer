package currency

import "testing"

func TestCurrency(t *testing.T) {
	t.Run("Creating a currency value with valid code should not fail", func(t *testing.T) {
		_, err := NewValue(100, "EUR", 0)
		if err != nil {
			t.Error("Currency should have been created")
		}
	})
	t.Run("Creating a currency value with invalid code should fail", func(t *testing.T) {
		_, err := NewValue(100, "abcd", 0)
		if err == nil {
			t.Error("Currency should have not been created")
		}
	})
}
