package models

import (
	"time"

	"github.com/go-pg/pg/v10"
)

//go:generate mockgen -source=$GOFILE -destination=../mocks/model_system_mock.go -package=mocks SystemImpl

// SystemImpl ...
type SystemImpl interface {
	Get() (*System, error)
	Update(s *System) error
}

// System global system variables
type System struct {
	tableName struct{} `pg:"system,alias:s"` //nolint
	ID        int
	LastCheck time.Time
}

// SystemRepo ...
type SystemRepo struct {
	db *pg.DB
}

// NewSystemModel ...
func NewSystemModel(db *pg.DB) *SystemRepo {
	return &SystemRepo{
		db: db,
	}
}

// Get ...
func (r *SystemRepo) Get() (*System, error) {
	system := &System{}
	err := r.db.Model(system).Select()

	return system, err
}

// Update ...
func (r *SystemRepo) Update(s *System) error {
	_, err := r.db.Model(s).WherePK().Update()

	return err
}
