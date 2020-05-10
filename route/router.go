package route

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/FreakyGranny/launchpad-api/api"
)


// Init echo framework
func Init() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))

	e.POST("/login", api.Login)
	e.OPTIONS("/login", api.Login)

	p := e.Group("/project")
	p.Use(middleware.JWT([]byte("secret")))
	p.GET("", api.GetProjects)
	p.GET("/:id", api.GetSingleProject)

	u := e.Group("/user")
	u.Use(middleware.JWT([]byte("secret")))
	u.GET("", api.GetUsers)
	u.GET("/:id", api.GetUser)

	c := e.Group("/category")
	c.Use(middleware.JWT([]byte("secret")))
	c.GET("", api.GetCategory)

	pt := e.Group("/project_type")
	pt.Use(middleware.JWT([]byte("secret")))
	pt.GET("", api.GetProjectType)

	d := e.Group("/donation")
	d.Use(middleware.JWT([]byte("secret")))
	d.GET("", api.GetDonation)

	return e
}
