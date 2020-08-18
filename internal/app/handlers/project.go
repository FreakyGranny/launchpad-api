package handlers

import (
	"net/http"
	"strconv"

	"github.com/FreakyGranny/launchpad-api/internal/app/models"
	"github.com/labstack/echo/v4"
)

const (
	dateLayout     = "2006-01-02"
	dateTimeLayout = "2006-01-02 15:04:05"
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

// type projectRequest struct {
// 	Title        string `json:"title"`
// 	SubTitle     string `json:"subtitle"`
// 	ReleaseDate  string `json:"release_date"`
// 	EventDate    string `json:"event_date,omitempty"`
// 	Category     uint   `json:"category"`
// 	GoalPeople   uint   `json:"goal_people"`
// 	GoalAmount   uint   `json:"goal_amount"`
// 	ImageLink    string `json:"image_link"`
// 	Instructions string `json:"instructions"`
// 	Description  string `json:"description"`
// 	ProjectType  uint   `json:"project_type"`
// 	Published    bool   `json:"published,omitempty"`
// }

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
// @Produce  json
// @Param page query int false "Page num"
// @Param page_size query int false "Capasity of one page"
// @Param category query int false "Category ID"
// @Param project_type query int false "Project Type ID"
// @Param open query bool false "Return only open"
// @Success 200 {object} ProjectListResponse
// @Security Bearer
// @Router /project [get]
func (h *ProjectHandler) GetProjects(c echo.Context) error {
	// userID, err := getUserIDFromToken(c.Get("user"))
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, err)
	// }

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
		// strategy, err := misc.GetStrategy(project.ProjectType)
		// if err != nil {
		// 	return c.JSON(http.StatusInternalServerError, nil)
		// }
		plv := ProjectListView{
			ID:          project.ID,
			Title:       project.Title,
			SubTitle:    project.SubTitle,
			Status:      project.Status(),
			ReleaseDate: project.ReleaseDate.Format(dateLayout),
			ImageLink:   project.ImageLink,
			Total:       project.Total,
			// Percent: strategy.Percent(&project),
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

// GetSingleProject godoc
// @Summary Show a single project
// @Description Returns project by ID
// @Tags project
// @ID get-project-by-id
// @Produce  json
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

	// strategy, err := misc.GetStrategy(project.ProjectType)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, nil)
	// }

	projectResponse := ProjectDetailView{
		ID:          project.ID,
		Title:       project.Title,
		SubTitle:    project.SubTitle,
		Status:      project.Status(),
		ReleaseDate: project.ReleaseDate.Format(dateLayout),
		ImageLink:   project.ImageLink,
		Total:       project.Total,
		// Percent:      strategy.Percent(&project),
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

// // CreateProject create new project
// func CreateProject(c echo.Context) error {
// 	cpRequest := new(projectRequest)
// 	if err := c.Bind(cpRequest); err != nil {
// 		return err
// 	}

// 	userToken := c.Get("user").(*jwt.Token)
// 	claims := userToken.Claims.(jwt.MapClaims)
// 	userID := claims["id"].(float64)

// 	releaseTime, err := time.Parse(dateLayout, cpRequest.ReleaseDate)
// 	if err != nil {
// 		return err
// 	}
// 	var eventTime time.Time

// 	if cpRequest.EventDate != "" {
// 		eventTime, err = time.Parse(dateLayout, cpRequest.EventDate)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	newProject := db.Project{
// 		OwnerID: uint(userID),
// 		Title: cpRequest.Title,
// 		SubTitle: cpRequest.SubTitle,
// 		ReleaseDate: releaseTime,
// 		EventDate: eventTime,
// 		GoalPeople: cpRequest.GoalPeople,
// 		GoalAmount: cpRequest.GoalAmount,
// 		Description: cpRequest.Description,
// 		ImageLink: cpRequest.ImageLink,
// 		Instructions: cpRequest.Instructions,
// 		CategoryID: cpRequest.Category,
// 		ProjectTypeID: cpRequest.ProjectType,
// 	}
// 	client := db.GetDbClient()
// 	client.Create(&newProject)

// 	return c.JSON(http.StatusCreated, map[string]uint{
// 		"id": newProject.ID,
// 	})
// }

// // UpdateProject update single value of project
// func UpdateProject(c echo.Context) error {
// 	userToken := c.Get("user").(*jwt.Token)
// 	claims := userToken.Claims.(jwt.MapClaims)
// 	userID := claims["id"].(float64)

// 	projectParam := c.Param("id")
// 	projectID, _ := strconv.Atoi(projectParam)

// 	dbClient := db.GetDbClient()
// 	var project db.Project

// 	if err := dbClient.Preload("ProjectType").First(&project, projectID).Error; gorm.IsRecordNotFoundError(err) {
// 		return c.JSON(http.StatusNotFound, nil)
// 	}
// 	if project.Published || project.OwnerID != uint(userID) {
// 		return c.JSON(http.StatusForbidden, nil)
// 	}

// 	upRequest := new(projectRequest)
// 	if err := c.Bind(upRequest); err != nil {
// 		return err
// 	}
// 	var parseErr error
// 	releaseTime := project.ReleaseDate
// 	eventTime := project.EventDate

// 	if upRequest.ReleaseDate != "" {
// 		releaseTime, parseErr = time.Parse(dateLayout, upRequest.ReleaseDate)
// 		if parseErr != nil {
// 			return parseErr
// 		}
// 	}
// 	if upRequest.EventDate != "" {
// 		eventTime, parseErr = time.Parse(dateTimeLayout, upRequest.EventDate)
// 		if parseErr != nil {
// 			return parseErr
// 		}
// 	}
// 	dbClient.Model(&project).Updates(db.Project{
// 		Title: upRequest.Title,
// 		SubTitle: upRequest.SubTitle,
// 		Instructions: upRequest.Instructions,
// 		Description: upRequest.Description,
// 		ImageLink: upRequest.ImageLink,
// 		CategoryID: upRequest.Category,
// 		ProjectTypeID: upRequest.ProjectType,
// 		GoalAmount: upRequest.GoalAmount,
// 		GoalPeople: upRequest.GoalPeople,
// 		ReleaseDate: releaseTime,
// 		EventDate: eventTime,
// 		Published: upRequest.Published,
// 	})
// 	if upRequest.Published {
// 		ch := misc.GetUpdatePipe()
// 		ch <- project.OwnerID
// 	}
// 	strategy, err := misc.GetStrategy(project.ProjectType)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, nil)
// 	}
// 	projectResponse := ProjectDetailView{
// 		ID: project.ID,
// 		Title: project.Title,
// 		SubTitle: project.SubTitle,
// 		Status: project.Status(),
// 		ReleaseDate: project.ReleaseDate.Format(dateLayout),
// 		ImageLink: project.ImageLink,
// 		Total: project.Total,
// 		Percent: strategy.Percent(&project),
// 		Category: project.Category,
// 		ProjectType: project.ProjectType,
// 		GoalPeople: project.GoalPeople,
// 		GoalAmount: project.GoalAmount,
// 		Description: project.Description,
// 		Instructions: project.Instructions,
// 		Owner: project.Owner,
// 	}
// 	if project.EventDate.IsZero() {
// 		projectResponse.EventDate = nil
// 	} else {
// 		ed := project.EventDate.Format(dateTimeLayout)
// 		projectResponse.EventDate = &ed
// 	}

// 	return c.JSON(http.StatusOK, projectResponse)
// }

// // DeleteProject delete project
// func DeleteProject(c echo.Context) error {
// 	userToken := c.Get("user").(*jwt.Token)
// 	claims := userToken.Claims.(jwt.MapClaims)
// 	userID := claims["id"].(float64)

// 	projectParam := c.Param("id")
// 	projectID, _ := strconv.Atoi(projectParam)

// 	dbClient := db.GetDbClient()
// 	var project db.Project

// 	if err := dbClient.First(&project, projectID).Error; gorm.IsRecordNotFoundError(err) {
// 		return c.JSON(http.StatusNotFound, nil)
// 	}
// 	if project.Published || project.OwnerID != uint(userID) {
// 		return c.JSON(http.StatusForbidden, nil)
// 	}

// 	dbClient.Delete(&project)

// 	return c.JSON(http.StatusNoContent, nil)
// }
