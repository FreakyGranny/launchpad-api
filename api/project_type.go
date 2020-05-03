package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/FreakyGranny/launchpad-api/db"
)

// GetProjectType return list of project types
func GetProjectType(c echo.Context) error {
	dbClient := db.GetDbClient()
	var projectTypes []db.ProjectType
	dbClient.Find(&projectTypes)

	return c.JSON(http.StatusOK, projectTypes)
}
