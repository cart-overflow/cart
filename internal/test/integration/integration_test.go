package integration

import (
	"testing"

	"github.com/cart-overflow/cart/internal/cart_test"
	"github.com/cart-overflow/cart/internal/test"
)

func TestIntegration(t *testing.T) {
	test.SetupValkey(t)
	t.Run("cart integration tests", cart_test.Integration)
}
