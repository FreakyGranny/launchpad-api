package handlers

import (
	"net/http"

	"github.com/FreakyGranny/launchpad-api/internal/app/models"
	"github.com/labstack/echo/v4"
)

// CategoryHandler ...
type CategoryHandler struct {
	CategoryModel models.CategoryImpl
}

// NewCategoryHandler ...
func NewCategoryHandler(c models.CategoryImpl) *CategoryHandler {
	return &CategoryHandler{CategoryModel: c}
}

// GetCategories godoc
// @Summary Returns list of categories
// @Description Returns list of categories
// @Tags category
// @ID get-categories
// @Produce  json
// @Success 200 {object} models.Category
// @Security Bearer
// @Router /category [get]
func (h *CategoryHandler) GetCategories(c echo.Context) error {
	categories, err := h.CategoryModel.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("unable to get categories"))
	}

	return c.JSON(http.StatusOK, categories)
}
