package main

import (
	"os"

	"github.com/FreakyGranny/launchpad-api/cmd"
	"github.com/spf13/cobra"
)

//go:generate swag i -g cmd/api.go -o docs

func main() {
	var rootCmd = &cobra.Command{Use: "lpad"}
	rootCmd.AddCommand(cmd.NewAPICmd(), cmd.NewMigrateCmd())
	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}
