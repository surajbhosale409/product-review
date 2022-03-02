package product_review

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	ID              int       `json:"id" bson:"id"`
	Name            string    `json:"name" bson:"name"`
	Description     string    `json:"description" bson:"description"`
	ThumbnailImgURL string    `json:"thumbnail_img_url,omitempty" bson:"thumbnail_img_url"`
	OverallRating   int       `json:"overall_rating" bson:"overall_rating"`
	Reviews         []*Review `json:"reviews,omitempty" bson:"reviews"`
}

// FindProductByID finds and returns a product with given id
func (s *Service) FindProductByID(id int) (product Product, err error) {
	filter := bson.M{"id": id}
	err = s.ProductCollection.FindOne(context.TODO(), filter).Decode(&product)
	return
}

// UpdateProductByID updates product with given id and update params
func (s *Service) UpdateProductByID(id int, update primitive.D) (product Product, err error) {
	filter := bson.M{"id": id}
	err = s.ProductCollection.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(1),
	).Decode(&product)
	return
}

// FindProducts filters the product based on given name pattern, also supports pagination
func (s *Service) FindProducts(name string, limit, skip int) (products []*Product, err error) {
	products = make([]*Product, 0)

	var filter interface{}
	if name != "" {
		filter = bson.M{
			"name": bson.M{
				"$regex": primitive.Regex{
					Pattern: "^" + name + "$",
					Options: "i",
				},
			},
		}
	} else {
		filter = bson.D{{}}
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))

	var cur *mongo.Cursor
	if cur, err = s.ProductCollection.Find(context.TODO(), filter, findOptions); err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var elem Product
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		// TODO: remove once aggregation is handled by mongodb
		elem.OverallRating = getAverageRating(elem.Reviews)
		elem.Reviews = nil
		products = append(products, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(context.TODO())
	return
}

// GetProductsHandler serves API requests for getting products using filters and pagination options
func (s *Service) GetProductsHandler(c echo.Context) (err error) {
	var products []*Product
	queryParams := c.QueryParams()

	// pagination options - limit, skip
	limit, _ := strconv.Atoi(queryParams.Get("limit"))
	if limit == 0 {
		limit = 10
	}
	skip, _ := strconv.Atoi(queryParams.Get("skip"))

	// filters - name
	name := queryParams.Get("name")

	if products, err = s.FindProducts(name, limit, skip); err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, Error{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, products)
}
