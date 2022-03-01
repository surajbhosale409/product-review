package cmd

import (
	"errors"
	service "product-review"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

func initialise() (s *service.Service, err error) {
	log.Info("Initialising service")
	cfg := &service.Config{}

	if err = godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	if err = env.Parse(cfg); err != nil {
		return nil, err
	}

	if cfg.Username == "" || cfg.Password == "" {
		return nil, errors.New("ENV variables PR_USERNAME, PR_PASSWORD must be configured for basic auth")
	}

	if cfg.MongoDBURL == "" {
		cfg.MongoDBURL = "mongodb://localhost:27017"
	}

	if cfg.MongoDBName == "" {
		cfg.MongoDBName = "product-review"
	}

	if s, err = service.NewService(cfg); err != nil {
		return nil, err
	}
	return
}

var RootCmd = &cobra.Command{
	Use: "product-review",
	Run: executeRootCmd,
	Long: `Product review service
`,
}

func executeRootCmd(cmd *cobra.Command, args []string) {
	cmd.Help()
}
