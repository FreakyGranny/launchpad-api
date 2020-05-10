package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/vcraescu/go-paginator"
	"github.com/vcraescu/go-paginator/adapter"

	"github.com/FreakyGranny/launchpad-api/db"
)

// ProjectListResponse paginated projects
type ProjectListResponse struct {
	Results  []ProjectListEntry `json:"results"`
	NextPage int                `json:"next"`
	HasNext  bool               `json:"has_next"`
}

// ProjectListEntry light project entry
type ProjectListEntry struct {
	ID          uint           `json:"id"`
	Title       string         `json:"title"`
	SubTitle    string         `json:"subtitle"`
	Status      string         `json:"status"`
	ReleaseDate time.Time      `json:"release_date"`
	EventDate   time.Time      `json:"event_date"`
	ImageLink   string         `json:"image_link"`
	Total       uint           `json:"total"`
	Percent     uint           `json:"percent"`
	Category    db.Category    `json:"category"`
	ProjectType db.ProjectType `json:"project_type"`
}

func filterQuery(userID int, filter string, dbClient *gorm.DB) *gorm.DB {
	if filter == "open" {
		return dbClient.Where("closed = ?", false)
	}
	if filter == "owned" {
		return dbClient.Where("owner_id = ?", userID)
	}
	if filter == "contributed" {
		return dbClient.Where("id IN (?)", dbClient.Table("donations").Select("project_id").Where("user_id = ?", userID).SubQuery())
	}
	
	return dbClient
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

	allRows := client.Preload("ProjectType").Preload("Category").Model(db.Project{}).Where("published = ?", true)  
  
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
	
	projectListEntries := make([]ProjectListEntry, 0)

	for _, project := range(projects) {
		projectListEntries = append(projectListEntries, ProjectListEntry{
			ID: project.ID,
			Title: project.Title,
			SubTitle: project.SubTitle,
			Status: project.Status(),
			ReleaseDate: project.ReleaseDate,
			EventDate: project.EventDate,
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

	return c.JSON(http.StatusOK, project)
}
