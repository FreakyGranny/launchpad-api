package models

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

//go:generate mockgen -source=$GOFILE -destination=../mocks/model_project_mock.go -package=mocks .

const (
	// StatusDraft draft project
	StatusDraft string = "draft"
	// StatusSuccess success project
	StatusSuccess string = "success"
	// StatusFail fail project
	StatusFail string = "fail"
	// StatusHarvest harvest project
	StatusHarvest string = "harvest"
	// StatusSearch search project
	StatusSearch string = "search"
)

// ProjectImpl ...
type ProjectImpl interface {
	Get(id int) (*Project, bool)
	GetProjectsWithPagination(f *ProjectListFilter) (ProjectPaginatorImpl, error)
	GetUserProjects(f *ProjectUserFilter) (*[]Project, error)
	Create(p *Project) error
	Update(p *Project) error
	Delete(id int, userID int) error
}

// Project model
type Project struct {
	tableName     struct{} `pg:"projects,alias:p"` //nolint
	ID            int
	Title         string
	SubTitle      string
	ReleaseDate   time.Time
	EventDate     time.Time
	GoalPeople    int `pg:",use_zero"`
	GoalAmount    int `pg:",use_zero"`
	Total         int `pg:",use_zero"`
	Description   string
	ImageLink     string
	Instructions  string
	Locked        bool `pg:",use_zero"`
	Published     bool `pg:",use_zero"`
	Closed        bool `pg:",use_zero"`
	Owner         User
	OwnerID       int
	Category      Category
	CategoryID    int
	ProjectType   ProjectType
	ProjectTypeID int
}

// Status of project
func (p *Project) Status() string {
	if !p.Published {
		return StatusDraft
	}
	if p.Closed {
		if p.Locked {
			return StatusSuccess
		}
		return StatusFail
	}
	if p.Locked {
		return StatusHarvest
	}

	return StatusSearch
}

// ProjectListFilter ...
type ProjectListFilter struct {
	Category    int
	ProjectType int
	OnlyOpen    bool
	PageSize    int
	Page        int
}

// ProjectUserFilter ...
type ProjectUserFilter struct {
	UserID      int
	Contributed bool
	Owned       bool
}

// ProjectPaginatorImpl ...
type ProjectPaginatorImpl interface {
	NextPage() (int, bool)
	Retrieve() (*[]Project, error)
}

// ProjectPaginator ...
type ProjectPaginator struct {
	EntryCount int
	PageSize   int
	Query      *orm.Query
	Page       int
	Values     *[]Project
}

// NextPage ...
func (pp *ProjectPaginator) NextPage() (int, bool) {
	if pp.Page*pp.Page < pp.EntryCount {
		return pp.Page + 1, true
	}
	return 0, false
}

// Retrieve ...
func (pp *ProjectPaginator) Retrieve() (*[]Project, error) {
	var ofst int

	if pp.Page > 1 {
		ofst = pp.PageSize * (pp.Page - 1)
	}
	err := pp.Query.Offset(ofst).Limit(pp.PageSize).Order("id DESC").Select()
	if err != nil {
		return nil, err
	}

	return pp.Values, nil
}

// ProjectRepo ...
type ProjectRepo struct {
	db *pg.DB
}

// NewProjectModel ...
func NewProjectModel(db *pg.DB) *ProjectRepo {
	return &ProjectRepo{
		db: db,
	}
}

// Get ...
func (r *ProjectRepo) Get(id int) (*Project, bool) {
	project := &Project{}
	err := r.db.Model(project).Relation("Owner").Relation("Category").Relation("ProjectType").Where("p.id = ?", id).Select()
	if err != nil {
		return project, false
	}

	return project, true
}

// GetProjectsWithPagination ...
func (r *ProjectRepo) GetProjectsWithPagination(f *ProjectListFilter) (ProjectPaginatorImpl, error) {
	projects := []Project{}
	q := r.db.Model(&projects).Relation("Category").Relation("ProjectType").Where("p.published", true)
	if f.Category != 0 {
		q = q.Where("category_id = ?", f.Category)
	}
	if f.ProjectType != 0 {
		q = q.Where("project_type_id = ?", f.ProjectType)
	}
	if f.OnlyOpen {
		q = q.Where("closed = ?", false)
	}
	x, err := q.Count()
	if err != nil {
		return nil, err
	}
	return &ProjectPaginator{
		EntryCount: x,
		Query:      q,
		Values:     &projects,
		Page:       f.Page,
		PageSize:   f.PageSize,
	}, nil
}

// GetUserProjects ...
func (r *ProjectRepo) GetUserProjects(f *ProjectUserFilter) (*[]Project, error) {
	projects := []Project{}
	q := r.db.Model(&projects).Relation("Category").Relation("ProjectType")
	if f.UserID != 0 {
		q = q.Where("owner_id = ?", f.UserID)
	}
	if !f.Owned {
		q = q.Where("p.published", true)
	}
	if f.Contributed {
		q = filterQueryByContribution(q, f.UserID)
	}
	err := q.Limit(20).Order("id DESC").Select()
	if err != nil {
		return nil, err
	}

	return &projects, nil
}

func filterQueryByContribution(q *orm.Query, userID int) *orm.Query {
	query := q
	// 	query = query.Where("id IN (?)", dbClient.Table("donations").Select("project_id").Where("user_id = ?", userID).SubQuery())

	return query
}

// Create new project
func (r *ProjectRepo) Create(p *Project) error {

	return nil
}

// Update new project
func (r *ProjectRepo) Update(p *Project) error {

	return nil
}

// Delete project by id
func (r *ProjectRepo) Delete(id int, userID int) error {
	// dbClient := db.GetDbClient()
	// var project db.Project

	// if err := dbClient.First(&project, projectID).Error; gorm.IsRecordNotFoundError(err) {
	// 	return c.JSON(http.StatusNotFound, nil)
	// }
	// if project.Published || project.OwnerID != uint(userID) {
	// 	return c.JSON(http.StatusForbidden, nil)
	// }

	// dbClient.Delete(&project)
	return nil
}
