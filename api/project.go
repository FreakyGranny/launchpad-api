package api

import (
	"net/http"
	"strconv"
	"time"

	// "github.com/labstack/gommon/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/vcraescu/go-paginator"
	"github.com/vcraescu/go-paginator/adapter"

	"github.com/FreakyGranny/launchpad-api/db"
)

const dateLayout = "2006-01-02"

// ProjectListResponse paginated projects
type ProjectListResponse struct {
	Results  []ProjectListView  `json:"results"`
	NextPage int                `json:"next"`
	HasNext  bool               `json:"has_next"`
}

// ProjectListView light project entry
type ProjectListView struct {
	ID          uint           `json:"id"`
	Title       string         `json:"title"`
	SubTitle    string         `json:"subtitle"`
	Status      string         `json:"status"`
	ReleaseDate string         `json:"release_date"`
	EventDate   string         `json:"event_date"`
	ImageLink   string         `json:"image_link"`
	Total       uint           `json:"total"`
	Percent     uint           `json:"percent"`
	Category    db.Category    `json:"category"`
	ProjectType db.ProjectType `json:"project_type"`
}

// ProjectDetailView light project entry
type ProjectDetailView struct {
	ID            uint           `json:"id"`
	Title         string         `json:"title"`
	SubTitle      string         `json:"subtitle"`
	Status        string         `json:"status"`
	ReleaseDate   string         `json:"release_date"`
	EventDate     string         `json:"event_date"`
	ImageLink     string         `json:"image_link"`
	Total         uint           `json:"total"`
	Percent       uint           `json:"percent"`
	Category      db.Category    `json:"category"`
	ProjectType   db.ProjectType `json:"project_type"`
    GoalPeople    uint           `json:"goal_people"`
    GoalAmount    uint           `json:"goal_amount"`
    Description   string         `json:"description"`
    Instructions  string         `json:"instructions"`
    Owner         db.User        `json:"owner"`
}

type projectCreateRequest struct {
	Title        string `json:"title"`
	SubTitle     string `json:"subtitle"`
	ReleaseDate  string `json:"release_date"`
	EventDate    string `json:"event_date,omitempty"`
	Category     uint   `json:"category"`
	GoalPeople   uint   `json:"goal_people"`
	GoalAmount   uint   `json:"goal_amount"`
	ImageLink    string `json:"image_link"`
	Instructions string `json:"instructions"`
	Description  string `json:"description"`
	ProjectType  uint   `json:"project_type"`
}


func filterQuery(userID int, filter string, dbClient *gorm.DB) *gorm.DB {
	query := dbClient
	if filter == "open" {
		query = query.Where("closed = ?", false)
	}
	if filter == "owned" {
		query = query.Where("owner_id = ?", userID)
	} else {
		query = query.Where("published = ?", true)
	}
	if filter == "contributed" {
		query = query.Where("id IN (?)", dbClient.Table("donations").Select("project_id").Where("user_id = ?", userID).SubQuery())
	}
	
	return query
}

func filterQueryByCategory(categoryID int, dbClient *gorm.DB) *gorm.DB {
	if categoryID != 0 {
		return dbClient.Where("category_id = ?", categoryID)
	}

	return dbClient
}

func filterQueryByProjectType(projectType int, dbClient *gorm.DB) *gorm.DB {
	if projectType != 0 {
		return dbClient.Where("project_type_id = ?", projectType)
	}

	return dbClient
}

func filterQueryByUserID(userID int, dbClient *gorm.DB) *gorm.DB {
	if userID != 0 {
		return dbClient.Where("owner_id = ?", userID)
	}

	return dbClient
}

// GetProjects return list of projects
func GetProjects(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(float64)

	categoryParam := c.QueryParam("category")
	typeParam := c.QueryParam("type")
	userParam := c.QueryParam("user")
	filterParam := c.QueryParam("filter")

	page := c.QueryParam("page")
	pageSize := c.QueryParam("page_size")

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	categoryInt, _ := strconv.Atoi(categoryParam)
	typeInt, _ := strconv.Atoi(typeParam)
	userInt, _ := strconv.Atoi(userParam)
  
	client := db.GetDbClient()
	var projects []db.Project

	allRows := client.Preload("ProjectType").Preload("Category").Model(db.Project{})  
  
	allRows = filterQueryByCategory(categoryInt, allRows)
	allRows = filterQueryByProjectType(typeInt, allRows)
	allRows = filterQueryByUserID(userInt, allRows)
	allRows = filterQuery(int(userID), filterParam, allRows)

	paginated := paginator.New(adapter.NewGORMAdapter(allRows.Order("id desc")), pageSizeInt)
	paginated.SetPage(pageInt)
  
	if err := paginated.Results(&projects); err != nil {
			  panic(err)
			}
	next, _ := paginated.NextPage()
	
	projectListEntries := make([]ProjectListView, 0)

	for _, project := range(projects) {
		projectListEntries = append(projectListEntries, ProjectListView{
			ID: project.ID,
			Title: project.Title,
			SubTitle: project.SubTitle,
			Status: project.Status(),
			ReleaseDate: project.ReleaseDate.Format(dateLayout),
			EventDate: project.EventDate.Format(dateLayout),
			ImageLink: project.ImageLink,
			Total: project.Total,
			Percent: project.Percent(),
			Category: project.Category,
			ProjectType: project.ProjectType,
		})
	}

	return c.JSON(http.StatusOK, ProjectListResponse{
		Results: projectListEntries,
		NextPage: next,
		HasNext: paginated.HasNext(),
	})
}

// GetSingleProject return single project
func GetSingleProject(c echo.Context) error {
	projectParam := c.Param("id")
	projectID, _ := strconv.Atoi(projectParam)
  
	dbClient := db.GetDbClient()
	var project db.Project

	if err := dbClient.Preload("ProjectType").Preload("Category").Preload("Owner").First(&project, projectID).Error; gorm.IsRecordNotFoundError(err) {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusOK, ProjectDetailView{
		ID: project.ID,
		Title: project.Title,
		SubTitle: project.SubTitle,
		Status: project.Status(),
		ReleaseDate: project.ReleaseDate.Format(dateLayout),
		EventDate: project.EventDate.Format(dateLayout),
		ImageLink: project.ImageLink,
		Total: project.Total,
		Percent: project.Percent(),
		Category: project.Category,
		ProjectType: project.ProjectType,
		GoalPeople: project.GoalPeople,
		GoalAmount: project.GoalAmount,
		Description: project.Description,
		Instructions: project.Instructions,
		Owner: project.Owner,
	})
}

// CreateProject create new project
func CreateProject(c echo.Context) error {
	cpRequest := new(projectCreateRequest)
	if err := c.Bind(cpRequest); err != nil {
		return err
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(float64)

	releaseTime, err := time.Parse(dateLayout, cpRequest.ReleaseDate)
	if err != nil {
		return err
	}
	var eventTime time.Time

	if cpRequest.EventDate != "" {
		eventTime, err = time.Parse(dateLayout, cpRequest.EventDate)
		if err != nil {
			return err
		}
	}

	newProject := db.Project{
		OwnerID: uint(userID),
		Title: cpRequest.Title,
		SubTitle: cpRequest.SubTitle,
		ReleaseDate: releaseTime,
		EventDate: eventTime,
		GoalPeople: cpRequest.GoalPeople,
		GoalAmount: cpRequest.GoalAmount,
		Description: cpRequest.Description,
		ImageLink: cpRequest.ImageLink,
		Instructions: cpRequest.Instructions,
		CategoryID: cpRequest.Category,
		ProjectTypeID: cpRequest.ProjectType,
	}
	client := db.GetDbClient()
	client.Create(&newProject)

	return c.JSON(http.StatusCreated, map[string]uint{
		"id": newProject.ID,
	})
}
