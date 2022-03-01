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
	ThumbnailImgURL string    `json:"thumbnail_img_url" bson:"thumbnail_img_url"`
	OverallRating   int       `json:"overall_rating" bson:"overall_rating"`
	Reviews         []*Review `json:"reviews,omitempty" bson:"reviews"`
}

func (s *Service) FindProductByID(id int) (product Product, err error) {
	filter := bson.M{"id": id}
	err = s.ProductCollection.FindOne(context.TODO(), filter).Decode(&product)
	return
}

func (s *Service) UpdateProductByID(id int, update primitive.D) (err error) {
	filter := bson.M{"id": id}
	_, err = s.ProductCollection.UpdateOne(context.TODO(), filter, update)
	return
}

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
		// create a value into which the single document can be decoded
		var elem Product
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		elem.Reviews = nil
		products = append(products, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	// Close the cursor once finished
	cur.Close(context.TODO())
	return
}

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
