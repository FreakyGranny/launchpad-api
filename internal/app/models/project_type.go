package models

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

//go:generate mockgen -destination=../mocks/model_p_type_mock.go -package=mocks . ProjectTypeImpl

// ProjectTypeImpl ...
type ProjectTypeImpl interface {
	GetAll() ([]ProjectType, error)
}

// ProjectType of project
type ProjectType struct {
	ID            uint           `db:"id" json:"id"`
	Alias         string         `db:"alias" json:"alias"`
	Name          string         `db:"name" json:"name"`
	Options       pq.StringArray `db:"options" json:"options"`
	GoalByPeople  bool           `db:"goal_by_people" json:"goal_by_people"`
	GoalByAmount  bool           `db:"goal_by_amount" json:"goal_by_amount"`
	EndByGoalGain bool           `db:"end_by_goal_gain" json:"end_by_goal_gain"`
}

// ProjectTypeRepo ...
type ProjectTypeRepo struct {
	db *sqlx.DB
}

// NewProjectTypeModel ...
func NewProjectTypeModel(db *sqlx.DB) *ProjectTypeRepo {
	return &ProjectTypeRepo{
		db: db,
	}
}

// GetAll ...
func (r *ProjectTypeRepo) GetAll() ([]ProjectType, error) {
	projectTypes := []ProjectType{}
	err := r.db.Select(&projectTypes, "SELECT * FROM project_types order by id asc")
	if err != nil {
		return projectTypes, err
	}

	return projectTypes, nil
}
