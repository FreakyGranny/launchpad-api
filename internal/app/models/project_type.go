package models

import (
	"github.com/go-pg/pg/v10"
)

//go:generate mockgen -source=$GOFILE -destination=../mocks/model_p_type_mock.go -package=mocks . ProjectTypeImpl

// ProjectTypeImpl ...
type ProjectTypeImpl interface {
	GetAll() ([]ProjectType, error)
}

// ProjectType of project
type ProjectType struct {
	tableName     struct{} `pg:"project_types,alias:pt"` //nolint
	ID            uint     `json:"id"`
	Alias         string   `json:"alias"`
	Name          string   `json:"name"`
	Options       []string `pg:",array" json:"options"`
	GoalByPeople  bool     `json:"goal_by_people"`
	GoalByAmount  bool     `json:"goal_by_amount"`
	EndByGoalGain bool     `json:"end_by_goal_gain"`
}

// ProjectTypeRepo ...
type ProjectTypeRepo struct {
	db *pg.DB
}

// NewProjectTypeModel ...
func NewProjectTypeModel(db *pg.DB) *ProjectTypeRepo {
	return &ProjectTypeRepo{
		db: db,
	}
}

// GetAll ...
func (r *ProjectTypeRepo) GetAll() ([]ProjectType, error) {
	projectTypes := []ProjectType{}
	err := r.db.Model(&projectTypes).Select()
	if err != nil {
		return projectTypes, err
	}

	return projectTypes, nil
}
