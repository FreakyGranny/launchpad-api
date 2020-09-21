package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/jonboulle/clockwork"
	"github.com/labstack/gommon/log"

	"github.com/FreakyGranny/launchpad-api/internal/app"
	"github.com/FreakyGranny/launchpad-api/internal/auth"
	"github.com/FreakyGranny/launchpad-api/internal/config"
	"github.com/FreakyGranny/launchpad-api/internal/db"
	"github.com/FreakyGranny/launchpad-api/internal/models"
	"github.com/FreakyGranny/launchpad-api/internal/server"
	"github.com/spf13/cobra"
)

// @title Launchpad API
// @version 1.0
// @description This is a launchpad backend.

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

// API ...
func API(cmd *cobra.Command, args []string) {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.DebugMode {
		log.SetLevel(log.DEBUG)
	}

	d, err := db.Connect(&cfg.Db)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	sModel := models.NewSystemModel(d)
	pModel := models.NewProjectModel(d)
	ptModel := models.NewProjectTypeModel(d)
	uModel := models.NewUserModel(d)
	cModel := models.NewCategoryModel(d)
	dModel := models.NewDonationModel(d)

	ctx, cancel := context.WithCancel(context.Background())
	b := app.NewBackground(sModel, pModel, uModel)
	b.Start(ctx)
	e := server.New(
		app.New(cModel, uModel, pModel, ptModel, dModel, auth.NewVk(cfg.Vk), clockwork.NewRealClock(), cfg.JWTSecret),
		[]byte(cfg.JWTSecret),
	)
	go func() {
		if err := e.Start(":1323"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	signal.Stop(terminate)
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	cancel()
	b.Wait()
}

// NewAPICmd return api command
func NewAPICmd() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "run api",
		Long:  "starts launchpad API server",
		Run:   API,
	}
}
