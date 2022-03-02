package product_review

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Auth(t *testing.T) {
	s := createTestService()

	testData := []struct {
		name       string
		username   string
		password   string
		shouldPass bool
	}{
		{
			name:       "valid creds: positive",
			username:   "test",
			password:   "test",
			shouldPass: true,
		},
		{
			name:       "invalid creds: negative",
			username:   "test",
			password:   "abcd",
			shouldPass: false,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			status, err := s.Auth(test.username, test.password, c)
			assert.Nil(t, err)
			assert.Equal(t, test.shouldPass, status)
		})
	}
}
