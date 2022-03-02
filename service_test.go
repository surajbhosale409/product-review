package product_review

import (
	"github.com/labstack/echo/v4"
)

func createTestService() (s *Service) {
	s = &Service{
		Config: &Config{
			Username: "test",
			Password: "test",
			e:        echo.New(),
		},
	}
	return
}
