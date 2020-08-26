package handlers

import (
	"net/http"
	"strconv"

	"github.com/FreakyGranny/launchpad-api/internal/app/misc"
	"github.com/FreakyGranny/launchpad-api/internal/app/models"
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

// ProjectDetailView light project entry
type ProjectDetailView struct {
	ID           int                `json:"id"`
	Title        string             `json:"title"`
	SubTitle     string             `json:"subtitle"`
	Status       string             `json:"status"`
	ReleaseDate  string             `json:"release_date"`
	EventDate    *string            `json:"event_date"`
	ImageLink    string             `json:"image_link"`
	Total        int                `json:"total"`
	Percent      int                `json:"percent"`
	Category     models.Category    `json:"category"`
	ProjectType  models.ProjectType `json:"project_type"`
	GoalPeople   int                `json:"goal_people"`
	GoalAmount   int                `json:"goal_amount"`
	Description  string             `json:"description"`
	Instructions string             `json:"instructions"`
	Owner        models.User        `json:"owner"`
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
	ProjectModel models.ProjectImpl
}

// NewProjectHandler ...
func NewProjectHandler(m models.ProjectImpl) *ProjectHandler {
	return &ProjectHandler{ProjectModel: m}
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
	pFilter := &models.ProjectListFilter{
		Category:    categoryInt,
		ProjectType: typeInt,
		OnlyOpen:    onlyOpen,
		Page:        pageInt,
		PageSize:    pageSizeInt,
	}

	paginator, err := h.ProjectModel.GetProjectsWithPagination(pFilter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	next, hasNext := paginator.NextPage()

	projectListEntries := make([]ProjectListView, 0)

	projects, err := paginator.Retrieve()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	for _, project := range *projects {
		strategy, err := misc.GetStrategy(&project.ProjectType, h.ProjectModel)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
		plv := ProjectListView{
			ID:          project.ID,
			Title:       project.Title,
			SubTitle:    project.SubTitle,
			Status:      project.Status(),
			ReleaseDate: project.ReleaseDate.Format(dateLayout),
			ImageLink:   project.ImageLink,
			Total:       project.Total,
			Percent:     strategy.Percent(&project),
			Category:    project.Category,
			ProjectType: project.ProjectType,
		}
		if !project.EventDate.IsZero() {
			ed := project.EventDate.Format(dateTimeLayout)
			plv.EventDate = &ed
		}
		projectListEntries = append(projectListEntries, plv)
	}

	return c.JSON(http.StatusOK, ProjectListResponse{
		Results:  projectListEntries,
		NextPage: next,
		HasNext:  hasNext,
	})
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

	pFilter := &models.ProjectUserFilter{
		UserID:      userID,
		Contributed: onlyContributed,
		Owned:       onlyOwned,
	}

	projects, err := h.ProjectModel.GetUserProjects(pFilter)

	projectListEntries := make([]ProjectListView, 0)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	for _, project := range *projects {
		strategy, err := misc.GetStrategy(&project.ProjectType, h.ProjectModel)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
		plv := ProjectListView{
			ID:          project.ID,
			Title:       project.Title,
			SubTitle:    project.SubTitle,
			Status:      project.Status(),
			ReleaseDate: project.ReleaseDate.Format(dateLayout),
			ImageLink:   project.ImageLink,
			Total:       project.Total,
			Percent:     strategy.Percent(&project),
			Category:    project.Category,
			ProjectType: project.ProjectType,
		}
		if !project.EventDate.IsZero() {
			ed := project.EventDate.Format(dateTimeLayout)
			plv.EventDate = &ed
		}
		projectListEntries = append(projectListEntries, plv)
	}

	return c.JSON(http.StatusOK, projectListEntries)
}

// GetSingleProject godoc
// @Summary Show a single project
// @Description Returns project by ID
// @Tags project
// @ID get-project-by-id
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} ProjectDetailView
// @Security Bearer
// @Router /project/{id} [get]
func (h *ProjectHandler) GetSingleProject(c echo.Context) error {
	projectID, _ := strconv.Atoi(c.Param("id"))

	project, ok := h.ProjectModel.Get(projectID)
	if !ok {
		return c.JSON(http.StatusNotFound, nil)
	}

	strategy, err := misc.GetStrategy(&project.ProjectType, h.ProjectModel)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	projectResponse := ProjectDetailView{
		ID:           project.ID,
		Title:        project.Title,
		SubTitle:     project.SubTitle,
		Status:       project.Status(),
		ReleaseDate:  project.ReleaseDate.Format(dateLayout),
		ImageLink:    project.ImageLink,
		Total:        project.Total,
		Percent:      strategy.Percent(project),
		Category:     project.Category,
		ProjectType:  project.ProjectType,
		GoalPeople:   project.GoalPeople,
		GoalAmount:   project.GoalAmount,
		Description:  project.Description,
		Instructions: project.Instructions,
		Owner:        project.Owner,
	}

	if project.EventDate.IsZero() {
		projectResponse.EventDate = nil
	} else {
		ed := project.EventDate.Format(dateTimeLayout)
		projectResponse.EventDate = &ed
	}

	return c.JSON(http.StatusOK, projectResponse)
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

	newProject := models.Project{
		OwnerID:       userID,
		Title:         cpRequest.Title,
		SubTitle:      cpRequest.SubTitle,
		ReleaseDate:   releaseDate,
		EventDate:     eventTime,
		GoalPeople:    cpRequest.GoalPeople,
		GoalAmount:    cpRequest.GoalAmount,
		Description:   cpRequest.Description,
		ImageLink:     cpRequest.ImageLink,
		Instructions:  cpRequest.Instructions,
		CategoryID:    cpRequest.Category,
		ProjectTypeID: cpRequest.ProjectType,
		Closed:        false,
		Locked:        false,
		Published:     false,
		Total:         0,
	}
	err = h.ProjectModel.Create(&newProject)
	switch err {
	case nil:
		return c.JSON(http.StatusCreated, ProjectCreateResponse{ID: newProject.ID})
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
	projectID, _ := strconv.Atoi(c.Param("id"))
	project, ok := h.ProjectModel.Get(projectID)
	if !ok {
		return c.JSON(http.StatusNotFound, errorResponse("project not found"))
	}
	if project.Published || project.OwnerID != userID {
		return c.JSON(http.StatusForbidden, errorResponse("modification is not allowed"))
	}
	strategy, err := misc.GetStrategy(&project.ProjectType, h.ProjectModel)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("strategy not found"))
	}

	if upRequest.DropEventDate {
		err = h.ProjectModel.DropEventDate(project)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errorResponse("unable to drop event date"))
		}
	} else {
		releaseDate, err := parseDate(upRequest.ReleaseDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, errorResponse("wrong release date"))
		}
		eventTime, err := parseDateTime(upRequest.EventDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, errorResponse("wrong event date"))
		}

		project.Title = upRequest.Title
		project.SubTitle = upRequest.SubTitle
		project.Instructions = upRequest.Instructions
		project.Description = upRequest.Description
		project.ImageLink = upRequest.ImageLink
		project.CategoryID = upRequest.Category
		project.ProjectTypeID = upRequest.ProjectType
		project.GoalAmount = upRequest.GoalAmount
		project.GoalPeople = upRequest.GoalPeople
		project.ReleaseDate = releaseDate
		project.EventDate = eventTime
		project.Published = upRequest.Published

		err = h.ProjectModel.Update(project)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errorResponse("unable to update project"))
		}
	}
	projectResponse := ProjectDetailView{
		ID:           project.ID,
		Title:        project.Title,
		SubTitle:     project.SubTitle,
		Status:       project.Status(),
		ReleaseDate:  project.ReleaseDate.Format(dateLayout),
		ImageLink:    project.ImageLink,
		Total:        project.Total,
		Percent:      strategy.Percent(project),
		Category:     project.Category,
		ProjectType:  project.ProjectType,
		GoalPeople:   project.GoalPeople,
		GoalAmount:   project.GoalAmount,
		Description:  project.Description,
		Instructions: project.Instructions,
		Owner:        project.Owner,
	}
	if project.EventDate.IsZero() {
		projectResponse.EventDate = nil
	} else {
		ed := project.EventDate.Format(dateTimeLayout)
		projectResponse.EventDate = &ed
	}

	return c.JSON(http.StatusOK, projectResponse)
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
	project, ok := h.ProjectModel.Get(projectID)
	if !ok {
		return c.JSON(http.StatusNotFound, errorResponse("project not found"))
	}
	if project.Published || project.OwnerID != userID {
		return c.JSON(http.StatusForbidden, errorResponse("project modify not allowed"))
	}

	err = h.ProjectModel.Delete(project)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
	}

	return c.JSON(http.StatusNoContent, nil)
}
