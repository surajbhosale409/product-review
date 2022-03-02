package product_review

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func Test_FindProductByID(t *testing.T) {

	s := createTestService()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("Get product by ID: positive", func(mt *mtest.T) {
		s.ProductCollection = mt.Coll
		expectedProduct := Product{
			ID:          101,
			Name:        "ABC",
			Description: "abc def",
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "product-reviews.products", mtest.FirstBatch, bson.D{
			{"id", expectedProduct.ID},
			{"name", expectedProduct.Name},
			{"description", expectedProduct.Description},
		}))

		productResponse, err := s.FindProductByID(expectedProduct.ID)
		assert.Nil(mt, err)
		assert.Equal(mt, true, reflect.DeepEqual(expectedProduct, productResponse))
	})

}

func Test_UpdateProductByID(t *testing.T) {
	s := createTestService()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("Update product by ID: positive", func(mt *mtest.T) {
		s.ProductCollection = mt.Coll
		rating := 4
		expectedProduct := Product{
			ID:          101,
			Name:        "ABC",
			Description: "abc def",
			Reviews: []*Review{{
				ReviewerName: "test",
				Rating:       &rating,
			}},
		}

		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"value", bson.D{
				{"id", expectedProduct.ID},
				{"name", expectedProduct.Name},
				{"description", expectedProduct.Description},
				{"reviews", expectedProduct.Reviews},
			}},
		})

		update := bson.D{{"$set", expectedProduct}}

		productResponse, err := s.UpdateProductByID(expectedProduct.ID, update)
		assert.Nil(mt, err)
		assert.Equal(mt, true, reflect.DeepEqual(expectedProduct, productResponse))
	})

}

func Test_FindProducts(t *testing.T) {
	s := createTestService()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	expectedProducts := []*Product{
		{
			ID:          101,
			Name:        "ABC",
			Description: "abc def",
		},
		{
			ID:          102,
			Name:        "PQR",
			Description: "abc def",
		},
	}

	mt.Run("Get products: positive", func(mt *mtest.T) {
		s.ProductCollection = mt.Coll

		first := mtest.CreateCursorResponse(1, "product-reviews.products", mtest.FirstBatch, bson.D{
			{"id", expectedProducts[0].ID},
			{"name", expectedProducts[0].Name},
			{"description", expectedProducts[0].Description},
		})
		second := mtest.CreateCursorResponse(1, "product-reviews.products", mtest.NextBatch, bson.D{
			{"id", expectedProducts[1].ID},
			{"name", expectedProducts[1].Name},
			{"description", expectedProducts[1].Description},
		})
		killCursors := mtest.CreateCursorResponse(0, "product-reviews.products", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		productResponse, err := s.FindProducts("", 5, 0)
		assert.Nil(mt, err)
		assert.Equal(mt, true, reflect.DeepEqual(expectedProducts, productResponse))
	})
}

func Test_GetProductsHandler(t *testing.T) {
	s := createTestService()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	expectedProducts := []Product{
		{
			ID:          101,
			Name:        "ABC",
			Description: "abc def",
		},
		{
			ID:          102,
			Name:        "PQR",
			Description: "abc def",
		},
	}

	mt.Run("Get products handler: positive", func(mt *mtest.T) {
		s.ProductCollection = mt.Coll

		first := mtest.CreateCursorResponse(1, "product-reviews.products", mtest.FirstBatch, bson.D{
			{"id", expectedProducts[0].ID},
			{"name", expectedProducts[0].Name},
			{"description", expectedProducts[0].Description},
		})
		second := mtest.CreateCursorResponse(1, "product-reviews.products", mtest.NextBatch, bson.D{
			{"id", expectedProducts[1].ID},
			{"name", expectedProducts[1].Name},
			{"description", expectedProducts[1].Description},
		})
		killCursors := mtest.CreateCursorResponse(0, "product-reviews.products", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.SetPath("/products")

		if assert.NoError(t, s.GetProductsHandler(c)) {
			productResponse := make([]Product, 0)
			json.Unmarshal(rec.Body.Bytes(), &productResponse)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, true, reflect.DeepEqual(expectedProducts, productResponse))
		}

	})

}
