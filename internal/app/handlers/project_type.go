package handlers

// import (
// 	"net/http"

// 	"github.com/FreakyGranny/launchpad-api/db"
// 	"github.com/labstack/echo/v4"
// )

// // GetProjectType return list of project types
// func GetProjectType(c echo.Context) error {
// 	dbClient := db.GetDbClient()
// 	var projectTypes []db.ProjectType
// 	dbClient.Find(&projectTypes)

// 	return c.JSON(http.StatusOK, projectTypes)
// }
