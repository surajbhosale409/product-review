package product_review

import (
	"github.com/labstack/echo/v4/middleware"
)

func (s *Service) AddRoutes() {
	authenticated := s.e.Group("/api", middleware.BasicAuth(s.Auth), middleware.Logger())
	authenticated.GET("/products", s.GetProductsHandler)
	authenticated.POST("/products/:id/reviews", s.AddReviewHandler)
}
