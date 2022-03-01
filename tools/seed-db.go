package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/labstack/gommon/random"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MongoDBURL  = "mongodb://localhost:27017"
	MongoDBName = "product-review"
)

type Product struct {
	ID              int       `json:"id" bson:"id"`
	Name            string    `json:"name" bson:"name"`
	Description     string    `json:"description" bson:"description"`
	ThumbnailImgURL string    `json:"thumbnail_img_url" bson:"thumbnail_img_url"`
	OverallRating   int       `json:"overall_rating" bson:"overall_rating"`
	Reviews         []*Review `json:"reviews" bson:"reviews"`
}

type Review struct {
	ReviewerName  string `json:"reviewer_name" bson:"reviewer_name"`
	WrittenReview string `json:"written_review" bson:"written_review"`
	Rating        *int   `json:"rating" bson:"rating"`
}

func populateDB() {

	mongoDBClient, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(MongoDBURL),
	)
	if err != nil {
		log.Error(err)
	}
	productCollection := mongoDBClient.Database(MongoDBName).Collection("products")

	var products []interface{}
	for i := 0; i < 50; i++ {
		id := uuid.New().ID()
		product := &Product{
			ID:          int(id),
			Name:        random.New().String(5),
			Description: random.New().String(20),
		}

		products = append(products, product)
	}
	productCollection.InsertMany(context.TODO(), products)
}

func main() {
	populateDB()
}
