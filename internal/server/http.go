package server

import (
	_ "github.com/FreakyGranny/launchpad-api/docs" // openAPI
	"github.com/FreakyGranny/launchpad-api/internal/app"
	"github.com/FreakyGranny/launchpad-api/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// New returns new echo server.
func New(a app.Application, jwtSecret []byte) *echo.Echo {
	e := echo.New()

	e.GET("/docs/*", echoSwagger.WrapHandler)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))

	JWTmiddleware := middleware.JWT(jwtSecret)

	hc := handlers.NewCategoryHandler(a)
	c := e.Group("/category")
	c.Use(JWTmiddleware)
	c.GET("", hc.GetCategories)

	ha := handlers.NewAuthHandler(a)
	e.POST("/login", ha.Login)
	e.OPTIONS("/login", ha.Login)

	hu := handlers.NewUserHandler(a)
	u := e.Group("/user")
	u.Use(JWTmiddleware)
	u.GET("", hu.GetCurrentUser)
	u.GET("/:id", hu.GetUser)

	hpt := handlers.NewProjectTypeHandler(a)
	pt := e.Group("/project_type")
	pt.Use(JWTmiddleware)
	pt.GET("", hpt.GetProjectTypes)

	hp := handlers.NewProjectHandler(a)
	p := e.Group("/project")
	p.Use(JWTmiddleware)
	p.GET("", hp.GetProjects)
	p.GET("/user/:id", hp.GetUserProjects)
	p.GET("/:id", hp.GetSingleProject)
	p.POST("", hp.CreateProject)
	p.PATCH("/:id", hp.UpdateProject)
	p.DELETE("/:id", hp.DeleteProject)

	hd := handlers.NewDonationHandler(a)
	dg := e.Group("/donation")
	dg.Use(JWTmiddleware)
	dg.GET("", hd.GetUserDonations)
	// dg.GET("/project/:id", hd.GetProjectDonations)
	// dg.POST("", hd.CreateDonation)
	// dg.DELETE("/:id", hd.DeleteDonation)
	// dg.PATCH("/:id", hd.UpdateDonation)

	return e
}
