package models

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // ...
)

//go:generate mockgen -destination=../mocks/model_user_mock.go -package=mocks . UserImpl

// UserImpl ...
type UserImpl interface {
	FindByID(id int) (*User, bool)
	Create(*User) (*User, error)
	Update(*User) (*User, error)
}

// User model
type User struct {
	ID           int     `db:"id" json:"id"`
	Username     string  `db:"username" json:"username"`
	FirstName    string  `db:"first_name" json:"first_name"`
	LastName     string  `db:"last_name" json:"last_name"`
	Avatar       string  `db:"avatar" json:"avatar"`
	Email        string  `db:"email" json:"-"`
	IsAdmin      bool    `db:"is_admin" json:"is_admin"`
	ProjectCount int     `db:"project_count" json:"project_count"`
	SuccessRate  float32 `db:"success_rate" json:"success_rate"`
}

// UserRepo ...
type UserRepo struct {
	db *sqlx.DB
}

// NewUserModel ...
func NewUserModel(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

// FindByID ...
func (r *UserRepo) FindByID(id int) (*User, bool) {
	user := &User{}
	if err := r.db.Get(user, "SELECT * FROM users where id = $1 limit 1", id); err != nil {
		return user, false
	}

	return user, true
}

// Create ...
func (r *UserRepo) Create(u *User) (*User, error) {
	_, err := r.db.NamedExec("INSERT INTO users (id, username, first_name, last_name, avatar, email) VALUES (:id, :username, :first_name, :last_name, :avatar, :email)", u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Update ...
func (r *UserRepo) Update(u *User) (*User, error) {
	_, err := r.db.NamedExec("UPDATE users SET username=:username, first_name=:first_name, last_name=:last_name, avatar=:avatar, email=:email WHERE id=:id", u)
	if err != nil {
		return nil, err
	}

	return u, nil
}
