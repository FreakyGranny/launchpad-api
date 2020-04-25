package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/FreakyGranny/launchpad-api/db"
)


// GetProjects return list of projects
func GetProjects(c echo.Context) error {
	client := db.GetDbClient()
	// Чтение
	var product db.Product
	client.First(&product, 1) // find product with id 1
	// client.First(&product, "code = ?", "L1212") // find product with code l1212

	return c.JSON(http.StatusOK, product)
}
