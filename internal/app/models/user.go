package models

import (
	"github.com/go-pg/pg/v10"
)

//go:generate mockgen -source=$GOFILE -destination=../mocks/model_user_mock.go -package=mocks . UserImpl

// UserImpl ...
type UserImpl interface {
	FindByID(id int) (*User, bool)
	Create(*User) (*User, error)
	Update(*User) (*User, error)
	GetParticipation(id int) ([]Participation, error)
}

// User model
type User struct {
	tableName    struct{} `pg:"users,alias:u"` //nolint
	ID           int      `json:"id"`
	Username     string   `json:"username"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Avatar       string   `json:"avatar"`
	Email        string   `json:"-"`
	IsAdmin      bool     `json:"-"`
	ProjectCount int      `json:"project_count"`
	SuccessRate  float32  `json:"success_rate"`
}

// UserRepo ...
type UserRepo struct {
	db *pg.DB
}

// NewUserModel ...
func NewUserModel(db *pg.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

// FindByID ...
func (r *UserRepo) FindByID(id int) (*User, bool) {
	user := &User{}
	err := r.db.Model(user).Where("id = ?", id).Select()
	if err != nil {
		return user, false
	}

	return user, true
}

// Create ...
func (r *UserRepo) Create(u *User) (*User, error) {
	_, err := r.db.Model(u).Insert()
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Update ...
func (r *UserRepo) Update(u *User) (*User, error) {
	_, err := r.db.Model(u).WherePK().UpdateNotZero()
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Participation ...
type Participation struct {
	Cnt           uint `json:"count"`
	ProjectTypeID uint `json:"id"`
}

// GetParticipation ...
func (r *UserRepo) GetParticipation(id int) ([]Participation, error) {
	pts := make([]Participation, 0)
	err := r.db.Model((*Donation)(nil)).
		ColumnExpr("count(d.id) AS cnt").
		ColumnExpr("p.project_type_id").
		Join("JOIN projects as p ON d.project_id = p.id").
		Group("p.project_type_id").
		Where("d.user_id = ?", id).
		Select(&pts)
	if err != nil {
		return nil, err
	}

	return pts, nil
}
