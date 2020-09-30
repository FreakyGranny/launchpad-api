package handlers

import (
	"net/http"

	"github.com/FreakyGranny/launchpad-api/internal/app"
	"github.com/labstack/echo/v4"
)

// ProjectTypeHandler ...
type ProjectTypeHandler struct {
	app app.Application
}

// NewProjectTypeHandler ...
func NewProjectTypeHandler(a app.Application) *ProjectTypeHandler {
	return &ProjectTypeHandler{app: a}
}

// GetProjectTypes godoc
// @Summary return list of project types
// @Description Returns list of project types
// @Tags project type
// @ID get-project-types
// @Produce json
// @Success 200 {object} []models.ProjectType
// @Security Bearer
// @Router /project_type [get]
func (h *ProjectTypeHandler) GetProjectTypes(c echo.Context) error {
	projectTypes, err := h.app.GetProjectTypes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("unable to get project types"))
	}

	return c.JSON(http.StatusOK, projectTypes)
}
