package product_review

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func Test_validateReview(t *testing.T) {

	validRating := 4
	invalidRating := 10

	testData := []struct {
		name       string
		review     *Review
		shouldPass bool
	}{
		{
			name: "valid review: positive",
			review: &Review{
				ReviewerName: "ABC",
				Rating:       &validRating,
			},
			shouldPass: true,
		},
		{
			name: "invalid review, rating out of range: negative",
			review: &Review{
				ReviewerName: "ABC",
				Rating:       &invalidRating,
			},
			shouldPass: false,
		},
		{
			name: "invalid review, reviewer name empty: negative",
			review: &Review{
				ReviewerName: "",
				Rating:       &validRating,
			},
			shouldPass: false,
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			err := validateReview(test.review)
			assert.Equal(t, test.shouldPass, err == nil)
		})
	}
}

func Test_getAverageRating(t *testing.T) {

	rating1 := 2
	rating2 := 4
	expectedAverage := 3
	reviews := []*Review{
		{
			ReviewerName: "ABC",
			Rating:       &rating1,
		},
		{
			ReviewerName: "ABC",
			Rating:       &rating2,
		},
	}

	t.Run("get average rating: positive", func(t *testing.T) {
		assert.Equal(t, expectedAverage, getAverageRating(reviews))
	})
}

func Test_AddReviewHandler(t *testing.T) {
	s := createTestService()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

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

	mt.Run("Add product review: positive", func(mt *mtest.T) {
		s.ProductCollection = mt.Coll

		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"value", bson.D{
				{"id", expectedProduct.ID},
				{"name", expectedProduct.Name},
				{"description", expectedProduct.Description},
				{"reviews", expectedProduct.Reviews},
			}},
		})
		reqBody, _ := json.Marshal(expectedProduct.Reviews[0])
		// reqBody := `{"reviewer_name":"test","rating":4}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.e.NewContext(req, rec)
		c.SetPath("/products/:id/reviews")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(expectedProduct.ID))

		if assert.NoError(t, s.AddReviewHandler(c)) {
			reviewResponse := &Review{}
			json.Unmarshal(rec.Body.Bytes(), reviewResponse)
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, true, reflect.DeepEqual(expectedProduct.Reviews[0], reviewResponse))
		}

	})
}
