package product_review

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type Config struct {

	// Credentials for basic auth
	Username string `json:"username" env:"PR_USERNAME"`
	Password string `json:"password" env:"PR_PASSWORD"`

	// Echo HTTP Server
	e *echo.Echo

	// MongoDB config
	MongoDBURL        string `json:"mongodb_url" env:"MONGODB_URL"`
	MongoDBName       string `json:"mongodb_name" env:"MONGODB_NAME"`
	MongoDBClient     *mongo.Client
	ProductCollection *mongo.Collection
}

type Service struct {
	*Config
}

func (s *Service) setupIndexes() (err error) {

	models := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{{Key: "id", Value: bsonx.Int32(1)}},
		},
		{
			Keys: bsonx.Doc{{Key: "name", Value: bsonx.String("text")}},
		},
	}
	_, err = s.ProductCollection.Indexes().CreateMany(context.TODO(), models)
	return
}

func NewService(conf *Config) (s *Service, err error) {
	s = &Service{
		Config: conf,
	}

	// initialise http server
	s.e = echo.New()
	s.AddRoutes()

	// tempDB = make(map[string]*Product)
	// populateDB()

	if s.MongoDBClient, err = mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(s.Config.MongoDBURL),
	); err != nil {
		return nil, err
	}
	s.ProductCollection = s.MongoDBClient.Database(s.MongoDBName).Collection("products")
	err = s.setupIndexes()
	return
}

//Serve starts the halo-sync as a daemon process/service - Not being used currently
func (s *Service) Serve() {
	log.Info("Starting HTTP API server on port 80")
	log.Fatal(s.e.Start(":80"))
}
