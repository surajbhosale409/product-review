package main

import (
	"product-review/cmd"
)

func main() {
	// Sub Commands
	cmd.RootCmd.AddCommand(cmd.APICmd)
	cmd.RootCmd.Execute()
}
