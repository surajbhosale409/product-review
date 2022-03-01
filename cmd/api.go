package cmd

import (
	service "product-review"

	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

var APICmd = &cobra.Command{
	Use: "serve",
	Run: executeAPICmd,
	// Args:      cobra.ExactValidArgs(1),
	Long: `Starts the REST API server for product-review engine`,
}

func executeAPICmd(cmd *cobra.Command, args []string) {
	var s *service.Service
	var err error

	if s, err = initialise(); err != nil {
		log.Fatal(err.Error())
	}
	s.Serve()
}
