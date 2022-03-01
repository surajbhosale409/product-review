package product_review

import (
	"math/rand"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/labstack/gommon/random"
)

var tempDB map[string]*Product

type Config struct {

	// HTTP Server
	e *echo.Echo
}

type Service struct {
	*Config
}

func NewService(conf *Config) (svc *Service, err error) {
	svc = &Service{
		Config: conf,
	}

	// initialise http server
	svc.e = echo.New()
	svc.AddRoutes()

	tempDB = make(map[string]*Product)
	populateDB()
	return
}

//Serve starts the halo-sync as a daemon process/service - Not being used currently
func (s *Service) Serve() {
	log.Info("Starting HTTP API server on port 80")
	log.Fatal(s.e.Start(":80"))
}

func populateDB() {
	for i := 0; i < 50; i++ {
		id := uuid.New().String()
		product := &Product{
			ID:          id,
			Name:        random.New().String(5),
			Description: random.New().String(20),
		}
		reviews := make([]Review, 2)
		for _, review := range reviews {
			review.ReviewerName = random.New().String(8)
			review.Rating = rand.Intn(5)
			product.Reviews = append(product.Reviews, &review)
		}
		tempDB[id] = product
	}
}
