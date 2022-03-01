package product_review

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Product struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ThumbnailImgURL string    `json:"thumbnail_img_url"`
	Reviews         []*Review `json:"reviews"`
}

func (s *Service) GetProductsHandler(c echo.Context) (err error) {
	var products []*Product

	for _, product := range tempDB {
		products = append(products, product)
	}

	return c.JSON(http.StatusOK, products)
}
