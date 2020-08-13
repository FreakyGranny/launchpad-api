package models

import (
	"time"

	"github.com/go-pg/pg/v10"
)

//go:generate mockgen -destination=../mocks/model_project_mock.go -package=mocks . ProjectImpl

// ProjectImpl ...
type ProjectImpl interface {
	Get(id int) (*Project, bool)
	// GetAll() ([]Category, error)
}

// Project model
type Project struct {
	tableName     struct{} `pg:"projects,alias:p"`
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
