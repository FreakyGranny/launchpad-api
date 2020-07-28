package handlers

// import (
// 	"net/http"

// 	"github.com/labstack/echo/v4"
// 	"github.com/FreakyGranny/launchpad-api/db"
// )

// // GetCategory return list of categories
// func GetCategory(c echo.Context) error {
// 	dbClient := db.GetDbClient()
// 	var categories []db.Category
// 	dbClient.Find(&categories)

// 	return c.JSON(http.StatusOK, categories)
// }
