package main

import (
	"github.com/FreakyGranny/launchpad-api/cmd"
	"github.com/spf13/cobra"
)

//go:generate swag i -g cmd/api.go -o docs

func main() {
	var cmdAPI = &cobra.Command{
		Use:   "api",
		Short: "run api",
		Long:  "starts launchpad API server",
		Run:   cmd.API,
	}
	var cmdMigrate = &cobra.Command{
		Use:   "migrate",
		Short: "run migrations",
		Long:  "Apply databse migrations",
		Run:   cmd.Migrate,
	}
	var rootCmd = &cobra.Command{Use: "lpad"}
	rootCmd.AddCommand(cmdAPI, cmdMigrate)
	rootCmd.Execute()
}
