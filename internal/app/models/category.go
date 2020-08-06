package models

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // ...
)

//go:generate mockgen -destination=../mocks/model_category_mock.go -package=mocks . CategoryImpl

// CategoryImpl ...
type CategoryImpl interface {
	Get(id int) (*Category, bool)
	GetAll() ([]Category, error)
}

// Category of project
type Category struct {
	ID    uint   `db:"id" json:"id"`
	Alias string `db:"alias" json:"alias"`
	Name  string `db:"name" json:"name"`
}

// CategoryRepo ...
type CategoryRepo struct {
	db *sqlx.DB
}

// NewCategoryModel ...
func NewCategoryModel(db *sqlx.DB) *CategoryRepo {
	return &CategoryRepo{
		db: db,
	}
}

// Get ...
func (r *CategoryRepo) Get(id int) (*Category, bool) {
	category := &Category{}
	if err := r.db.Get(category, "SELECT * FROM category where id = $1 limit 1", id); err != nil {
		return category, false
	}

	return category, true
}

// GetAll ...
func (r *CategoryRepo) GetAll() ([]Category, error) {
	categories := []Category{}
	err := r.db.Select(&categories, "SELECT * FROM categories order by id asc")
	if err != nil {
		return categories, err
	}

	return categories, nil
}
