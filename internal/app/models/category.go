package models

import (
	"github.com/go-pg/pg/v10"
)

//go:generate mockgen -destination=../mocks/model_category_mock.go -package=mocks . CategoryImpl

// CategoryImpl ...
type CategoryImpl interface {
	Get(id int) (*Category, bool)
	GetAll() ([]Category, error)
}

// Category of project
type Category struct {
	tableName struct{} `pg:"categories,alias:c"`
	ID        int      `json:"id"`
	Alias     string   `json:"alias"`
	Name      string   `json:"name"`
}

// CategoryRepo ...
type CategoryRepo struct {
	db *pg.DB
}

// NewCategoryModel ...
func NewCategoryModel(db *pg.DB) *CategoryRepo {
	return &CategoryRepo{
		db: db,
	}
}

// Get ...
func (r *CategoryRepo) Get(id int) (*Category, bool) {
	category := &Category{}
	err := r.db.Model(category).Where("id = ?", id).Select()
	if err != nil {
		return category, false
	}

	return category, true
}

// GetAll ...
func (r *CategoryRepo) GetAll() ([]Category, error) {
	categories := []Category{}
	err := r.db.Model(&categories).Select()
	if err != nil {
		return categories, err
	}

	return categories, nil
}
