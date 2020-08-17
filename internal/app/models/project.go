package models

import (
	"time"

	"github.com/go-pg/pg/v10"
)

//go:generate mockgen -source=$GOFILE -destination=../mocks/model_project_mock.go -package=mocks . ProjectImpl

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
	// GetAll() ([]Category, error)
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

// // GetAll ...
// func (r *CategoryRepo) GetAll() ([]Category, error) {
// 	categories := []Category{}
// 	err := r.db.Select(&categories, "SELECT * FROM categories order by id asc")
// 	if err != nil {
// 		return categories, err
// 	}

// 	return categories, nil
// }
