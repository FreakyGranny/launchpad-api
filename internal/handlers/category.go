package handlers

import (
	"net/http"

	"github.com/FreakyGranny/launchpad-api/internal/app"
	"github.com/labstack/echo/v4"
)

// CategoryHandler ...
type CategoryHandler struct {
	app app.Application
}

// NewCategoryHandler ...
func NewCategoryHandler(a app.Application) *CategoryHandler {
	return &CategoryHandler{app: a}
}

// GetCategories godoc
// @Summary Returns list of categories
// @Description Returns list of categories
// @Tags category
// @ID get-categories
// @Produce json
// @Success 200 {object} []models.Category
// @Security Bearer
// @Router /category [get]
func (h *CategoryHandler) GetCategories(c echo.Context) error {
	categories, err := h.app.GetCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("unable to get categories"))
	}

	return c.JSON(http.StatusOK, categories)
}
