package handlers

import (
	"net/http"

	"github.com/FreakyGranny/launchpad-api/internal/app/models"
	"github.com/labstack/echo/v4"
)

// ProjectTypeHandler ...
type ProjectTypeHandler struct {
	ProjectTypeModel models.ProjectTypeImpl
}

// NewProjectTypeHandler ...
func NewProjectTypeHandler(pt models.ProjectTypeImpl) *ProjectTypeHandler {
	return &ProjectTypeHandler{ProjectTypeModel: pt}
}

// GetProjectTypes godoc
// @Summary return list of project types
// @Description Returns list of categories
// @ID get-project-types
// @Produce json
// @Success 200 {object} models.ProjectType
// @Security Bearer
// @Router /project_type [get]
func (h *ProjectTypeHandler) GetProjectTypes(c echo.Context) error {
	projectTypes, err := h.ProjectTypeModel.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("unable to get project types"))
	}

	return c.JSON(http.StatusOK, projectTypes)
}
