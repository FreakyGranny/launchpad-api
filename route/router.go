package route

import (
	"github.com/labstack/echo/v4"

	"github.com/FreakyGranny/launchpad-api/api"
)


// Init echo framework
func Init() *echo.Echo {
	e := echo.New()

	e.GET("/project", api.GetProjects)

	return e
}
