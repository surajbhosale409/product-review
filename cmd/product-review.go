package cmd

import (
	service "product-review"

	"github.com/caarlos0/env/v6"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

var defaultConfig = &service.Config{}

func initialise() (svc *service.Service, err error) {
	log.Info("Initialising service")
	cfg := &service.Config{}

	if err = env.Parse(cfg); err != nil {
		return nil, err
	}

	if svc, err = service.NewService(cfg); err != nil {
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
