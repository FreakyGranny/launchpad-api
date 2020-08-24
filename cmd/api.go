package cmd

import (
	"github.com/jonboulle/clockwork"
	"github.com/labstack/gommon/log"

	_ "github.com/FreakyGranny/launchpad-api/docs" // openAPI
	"github.com/FreakyGranny/launchpad-api/internal/app/auth"
	"github.com/FreakyGranny/launchpad-api/internal/app/config"
	"github.com/FreakyGranny/launchpad-api/internal/app/db"
	"github.com/FreakyGranny/launchpad-api/internal/app/handlers"
	"github.com/FreakyGranny/launchpad-api/internal/app/misc"
	"github.com/FreakyGranny/launchpad-api/internal/app/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	echoSwagger "github.com/swaggo/echo-swagger"
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
	
	pModel := models.NewProjectModel(d)
	b := misc.NewBackground(pModel)

	go b.RecalcProject()
	// go misc.UpdateUser()
	// go misc.HarvestCheck()

	e := echo.New()
	e.GET("/docs/*", echoSwagger.WrapHandler)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))

	JWTmiddleware := middleware.JWT([]byte(cfg.JWTSecret))

	ha := handlers.NewAuthHandler(
		cfg.JWTSecret,
		models.NewUserModel(d),
		auth.NewVk(cfg.Vk),
		clockwork.NewRealClock(),
	)

	e.POST("/login", ha.Login)
	e.OPTIONS("/login", ha.Login)

	hu := handlers.NewUserHandler(models.NewUserModel(d))
	u := e.Group("/user")
	u.Use(JWTmiddleware)
	u.GET("", hu.GetCurrentUser)
	u.GET("/:id", hu.GetUser)

	hc := handlers.NewCategoryHandler(models.NewCategoryModel(d))
	c := e.Group("/category")
	c.Use(JWTmiddleware)
	c.GET("", hc.GetCategories)

	hpt := handlers.NewProjectTypeHandler(models.NewProjectTypeModel(d))
	pt := e.Group("/project_type")
	pt.Use(JWTmiddleware)
	pt.GET("", hpt.GetProjectTypes)

	hp := handlers.NewProjectHandler(pModel)
	p := e.Group("/project")
	p.Use(JWTmiddleware)
	p.GET("", hp.GetProjects)
	p.GET("/user/:id", hp.GetUserProjects)
	p.GET("/:id", hp.GetSingleProject)
	p.POST("", hp.CreateProject)
	p.PATCH("/:id", hp.UpdateProject)
	p.DELETE("/:id", hp.DeleteProject)

	hd := handlers.NewDonationHandler(models.NewDonationModel(d), b.GetRecalcPipe())
	dg := e.Group("/donation")
	dg.Use(JWTmiddleware)
	dg.GET("", hd.GetUserDonations)
	dg.GET("/project/:id", hd.GetProjectDonations)
	dg.POST("", hd.CreateDonation)
	dg.DELETE("/:id", hd.DeleteDonation)
	dg.PATCH("/:id", hd.UpdateDonation)

	log.Fatal(e.Start(":1323"))
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
