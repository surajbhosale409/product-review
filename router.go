package product_review

import (
	"github.com/labstack/echo/v4/middleware"
)

func (svc *Service) AddRoutes() {
	authenticated := svc.e.Group("/api", middleware.BasicAuth(Auth), middleware.Logger())
	authenticated.GET("/products", svc.GetProductsHandler)
}
