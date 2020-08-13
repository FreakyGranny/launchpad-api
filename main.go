package main

import (
	"github.com/FreakyGranny/launchpad-api/cmd"
	"github.com/spf13/cobra"
)

//go:generate swag i -g cmd/api.go -o docs

func main() {
	var rootCmd = &cobra.Command{Use: "lpad"}
	rootCmd.AddCommand(cmd.NewAPICmd(), cmd.NewMigrateCmd())
	rootCmd.Execute()
}
