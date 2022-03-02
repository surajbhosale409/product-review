package product_review

import (
	"testing"
)

func Test_AddRoutes(t *testing.T) {
	s := createTestService()
	t.Run("Test AddRoutes", func(t *testing.T) {
		s.AddRoutes()
	})
}
