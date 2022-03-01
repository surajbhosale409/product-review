package product_review

import (
	"github.com/labstack/echo/v4"
)

// Auth provides BasicAuth implementation with configured username and password
func (s *Service) Auth(username, password string, c echo.Context) (bool, error) {
	if username == s.Username && password == s.Password {
		return true, nil
	}
	return false, nil
}
