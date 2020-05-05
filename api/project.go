package api

import (
	"net/http"
	"strconv"
	"time"

	// "github.com/labstack/gommon/log"
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

// GetProjects return list of projects
func GetProjects(c echo.Context) error {
	categoryParam := c.QueryParam("category")
	typeParam := c.QueryParam("type")
	userParam := c.QueryParam("user")

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
  
	if categoryInt != 0 {
		allRows = allRows.Where("category_id = ?", categoryInt)
	}
	if typeInt != 0 {
		allRows = allRows.Where("project_type_id = ?", typeInt)
	}
	if userInt != 0 {
		allRows = allRows.Where("owner_id = ?", userInt)
	}

	paginated := paginator.New(adapter.NewGORMAdapter(allRows), pageSizeInt)
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
