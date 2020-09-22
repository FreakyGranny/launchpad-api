package handlers

import (
	"net/http"
	"strconv"

	"github.com/FreakyGranny/launchpad-api/internal/app"
	"github.com/FreakyGranny/launchpad-api/internal/models"
	"github.com/labstack/echo/v4"
)

// ProjectListResponse paginated projects
type ProjectListResponse struct {
	Results  []ProjectListView `json:"results"`
	NextPage int               `json:"next"`
	HasNext  bool              `json:"has_next"`
}

// ProjectListView light project entry
type ProjectListView struct {
	ID          int                `json:"id"`
	Title       string             `json:"title"`
	SubTitle    string             `json:"subtitle"`
	Status      string             `json:"status"`
	ReleaseDate string             `json:"release_date"`
	EventDate   *string            `json:"event_date"`
	ImageLink   string             `json:"image_link"`
	Total       int                `json:"total"`
	Percent     int                `json:"percent"`
	Category    models.Category    `json:"category"`
	ProjectType models.ProjectType `json:"project_type"`
}

// ProjectModifyRequest Request for project creation
type ProjectModifyRequest struct {
	Title         string `json:"title"`
	SubTitle      string `json:"subtitle"`
	ReleaseDate   string `json:"release_date"`
	EventDate     string `json:"event_date,omitempty"`
	Category      int    `json:"category"`
	GoalPeople    int    `json:"goal_people"`
	GoalAmount    int    `json:"goal_amount"`
	ImageLink     string `json:"image_link"`
	Instructions  string `json:"instructions"`
	Description   string `json:"description"`
	ProjectType   int    `json:"project_type"`
	Published     bool   `json:"published,omitempty"`
	DropEventDate bool   `json:"drop_event_date,omitempty"`
}

// ProjectCreateResponse Response for project creation
type ProjectCreateResponse struct {
	ID int `json:"id"`
}

// ProjectHandler ...
type ProjectHandler struct {
	app          app.Application
}

// NewProjectHandler ...
func NewProjectHandler(a app.Application) *ProjectHandler {
	return &ProjectHandler{app: a}
}

// GetProjects godoc
// @Summary Returns list of projects
// @Description Returns list of projects with filters
// @Tags project
// @ID get-projects
// @Produce json
// @Param page query int false "Page num"
// @Param page_size query int false "Capasity of one page"
// @Param category query int false "Category ID"
// @Param project_type query int false "Project Type ID"
// @Param open query bool false "Return only open"
// @Success 200 {object} ProjectListResponse
// @Security Bearer
// @Router /project [get]
func (h *ProjectHandler) GetProjects(c echo.Context) error {
	categoryInt, _ := strconv.Atoi(c.QueryParam("category"))
	typeInt, _ := strconv.Atoi(c.QueryParam("type"))
	onlyOpen, _ := strconv.ParseBool(c.QueryParam("open"))

	pageInt, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || pageInt == 0 {
		pageInt = 1
	}
	pageSizeInt, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSizeInt == 0 {
		pageSizeInt = 10
	}

	projects, next, hasNext, err := h.app.GetProjectsWithPagination(categoryInt, typeInt, pageInt, pageSizeInt, onlyOpen)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, ProjectListResponse{
		Results:  projectToListView(projects),
		NextPage: next,
		HasNext:  hasNext,
	})
}

func projectToListView(projects []*app.ExtendedProject) []ProjectListView {
	projectListEntries := make([]ProjectListView, 0)
	for _, project := range projects {
		plv := ProjectListView{
			ID:          project.ID,
			Title:       project.Title,
			SubTitle:    project.SubTitle,
			Status:      project.Status,
			ReleaseDate: project.ReleaseDate,
			ImageLink:   project.ImageLink,
			Total:       project.Total,
			Percent:     project.Percent,
			Category:    project.Category,
			ProjectType: project.ProjectType,
		}
		projectListEntries = append(projectListEntries, plv)
	}

	return projectListEntries
}

// GetUserProjects godoc
// @Summary Returns list of projects associated with user
// @Description Returns list of projects associated with user with filters
// @Tags project
// @ID get-user-projects
// @Produce json
// @Param owned query bool false "Return projects where user is owner"
// @Param contributed query bool false "Return projects where user is contributor"
// @Param id path int true "User ID"
// @Success 200 {object} []ProjectListView
// @Security Bearer
// @Router /project/user/{id} [get]
func (h *ProjectHandler) GetUserProjects(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("wrong user id"))
	}
	onlyOwned, _ := strconv.ParseBool(c.QueryParam("owned"))
	onlyContributed, _ := strconv.ParseBool(c.QueryParam("contributed"))

	projects, err := h.app.GetUserProjects(userID, onlyContributed, onlyOwned)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, projectToListView(projects))
}

// GetSingleProject godoc
// @Summary Show a single project
// @Description Returns project by ID
// @Tags project
// @ID get-project-by-id
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} app.ExtendedProject
// @Security Bearer
// @Router /project/{id} [get]
func (h *ProjectHandler) GetSingleProject(c echo.Context) error {
	projectID, _ := strconv.Atoi(c.Param("id"))

	project, err := h.app.GetProject(projectID)
	switch err {
	case app.ErrProjectNotFound:
		return c.JSON(http.StatusNotFound, nil)
	case nil:
		return c.JSON(http.StatusOK, project)
	default:
		return c.JSON(http.StatusInternalServerError, err)
	}
}

// CreateProject create new project
// @Summary Create project
// @Description Create new project
// @Tags project
// @ID post-project
// @Accept json
// @Produce json
// @Param request body DonationModifyRequest true "Request body"
// @Success 200 {object} ProjectCreateResponse
// @Security Bearer
// @Router /project [post]
func (h *ProjectHandler) CreateProject(c echo.Context) error {
	cpRequest := new(ProjectModifyRequest)
	if err := c.Bind(cpRequest); err != nil {
		return err
	}
	userID, err := getUserIDFromToken(c.Get("user"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("wrong ID"))
	}
	releaseDate, err := parseDate(cpRequest.ReleaseDate)
	if err != nil || releaseDate.IsZero() {
		return c.JSON(http.StatusBadRequest, errorResponse("wrong release date"))
	}
	eventTime, err := parseDateTime(cpRequest.EventDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("wrong event date"))
	}
	id, err := h.app.CreateProject(
		userID, 
		cpRequest.GoalPeople, 
		cpRequest.GoalAmount, 
		cpRequest.Category, 
		cpRequest.ProjectType,
		cpRequest.Title,
		cpRequest.SubTitle,
		cpRequest.Description,
		cpRequest.ImageLink,
		cpRequest.Instructions,
		releaseDate,
		eventTime,
	)
	switch err {
	case nil:
		return c.JSON(http.StatusCreated, ProjectCreateResponse{ID: id})
	case models.ErrUserNotFound:
		return c.JSON(http.StatusBadRequest, err)
	default:
		return c.JSON(http.StatusInternalServerError, err)
	}
}

// UpdateProject godoc
// @Summary update single value of project
// @Description Mofidy project fields
// @Tags project
// @ID update-project
// @Accept json
// @Produce json
// @Param request body DonationModifyRequest true "Request body"
// @Success 200 {object} ProjectDetailView
// @Security Bearer
// @Router /project/{id} [patch]
func (h *ProjectHandler) UpdateProject(c echo.Context) error {
	upRequest := new(ProjectModifyRequest)
	if err := c.Bind(upRequest); err != nil {
		return err
	}
	userID, err := getUserIDFromToken(c.Get("user"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("wrong ID"))
	}
	releaseDate, err := parseDate(upRequest.ReleaseDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("wrong release date"))
	}
	eventTime, err := parseDateTime(upRequest.EventDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("wrong event date"))
	}

	projectID, _ := strconv.Atoi(c.Param("id"))
	project, err := h.app.UpdateProject(
		projectID,
		userID, 
		upRequest.GoalPeople, 
		upRequest.GoalAmount, 
		upRequest.Category, 
		upRequest.ProjectType,
		upRequest.Title,
		upRequest.SubTitle,
		upRequest.Description,
		upRequest.ImageLink,
		upRequest.Instructions,
		releaseDate,
		eventTime,
		upRequest.Published,
		upRequest.DropEventDate,
	)
	switch err {
	case app.ErrProjectNotFound:
		return c.JSON(http.StatusNotFound, errorResponse("project not found"))
	case app.ErrProjectModifyNotAllowed:
		return c.JSON(http.StatusForbidden, errorResponse("modification is not allowed"))
	case nil:
		return c.JSON(http.StatusOK, project)
	default:
		return c.JSON(http.StatusInternalServerError, errorResponse("unable to update project"))
	}
}

// DeleteProject delete project
// @Summary Delete not published project
// @Description Delete not published project
// @Tags project
// @ID delete-project
// @Param id path int true "Project ID"
// @Success 204
// @Security Bearer
// @Router /project/{id} [delete]
func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	userID, err := getUserIDFromToken(c.Get("user"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("wrong ID"))
	}
	projectID, _ := strconv.Atoi(c.Param("id"))

	err = h.app.DeleteProject(userID, projectID)
	switch err {
	case app.ErrProjectNotFound:
		return c.JSON(http.StatusNotFound, errorResponse("project not found"))
	case app.ErrProjectModifyNotAllowed:
		return c.JSON(http.StatusForbidden, errorResponse("project modify not allowed"))
	case nil:
		return c.JSON(http.StatusNoContent, nil)
	default:
		return c.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
	}
}
