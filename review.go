package product_review

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
)

type Review struct {
	ReviewerName  string `json:"reviewer_name" bson:"reviewer_name"`
	WrittenReview string `json:"written_review" bson:"written_review"`
	Rating        *int   `json:"rating" bson:"rating"`
}

type Error struct {
	Message string `json:"message"`
}

var (
	ErrReviewerNameCannotBeEmpty = errors.New("reviewer_name cannot be empty")
	ErrRatingCannotBeEmpty       = errors.New("rating cannot be empty")
	ErrInvalidRatingValue        = errors.New("invalid rating value specified, should an integer from range 0-5")
)

func validateReview(review *Review) error {
	if review.ReviewerName == "" {
		return ErrReviewerNameCannotBeEmpty
	}

	if review.Rating == nil {
		return ErrRatingCannotBeEmpty
	}

	if *review.Rating < 0 || *review.Rating > 5 {
		return ErrInvalidRatingValue
	}

	return nil
}

func getAverageRating(reviews []*Review) int {
	sum := 0
	len := 0

	for _, review := range reviews {
		sum += *review.Rating
		len++
	}

	return sum / len
}

func (s *Service) AddReviewHandler(c echo.Context) (err error) {

	review := &Review{}
	if err = c.Bind(review); err != nil {
		return c.JSON(http.StatusBadRequest, &Error{
			Message: "Unable to parse request payload",
		})
	}

	if err = validateReview(review); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &Error{
			Message: err.Error(),
		})
	}

	var product Product

	id, _ := strconv.Atoi(c.Param("id"))
	if product, err = s.FindProductByID(id); err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, Error{
			Message: err.Error(),
		})
	}

	product.Reviews = append(product.Reviews, review)
	update := bson.D{{"$set", bson.D{{"reviews", product.Reviews}, {"overall_rating", getAverageRating(product.Reviews)}}}}

	if err = s.UpdateProductByID(id, update); err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, Error{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, review)
}
