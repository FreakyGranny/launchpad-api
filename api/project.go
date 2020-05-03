package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/FreakyGranny/launchpad-api/db"
)


// GetProjects return list of projects
func GetProjects(c echo.Context) error {
	client := db.GetDbClient()
	var products []db.Product
	client.Find(&products)

	return c.JSON(http.StatusOK, products)
}
