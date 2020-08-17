package models

import (
	"github.com/go-pg/pg/v10"
)

//go:generate mockgen -source=$GOFILE -destination=../mocks/model_donation_mock.go -package=mocks . DonationImpl

// DonationImpl ...
type DonationImpl interface {
	GetAllByUser(id int) ([]Donation, error)
	GetAllByProject(id int) ([]Donation, error)
}

// Donation for project
type Donation struct {
	tableName struct{} `pg:"donations,alias:d"` //nolint
	ID        int      `json:"id"`
	Payment   int      `json:"payment"`
	Locked    bool     `json:"locked"`
	Paid      bool     `json:"paid"`
	User      User     `json:"-"`
	UserID    int      `json:"-"`
	Project   Project  `json:"-"`
	ProjectID int      `json:"project"`
}

// DonationRepo ...
type DonationRepo struct {
	db *pg.DB
}

// NewDonationModel ...
func NewDonationModel(db *pg.DB) *DonationRepo {
	return &DonationRepo{
		db: db,
	}
}

// GetAllByUser returns all donations associated with user
func (r *DonationRepo) GetAllByUser(id int) ([]Donation, error) {
	donations := make([]Donation, 0)
	err := r.db.Model(&donations).Where("user_id = ", id).Select()
	if err != nil {
		return nil, err
	}

	return donations, nil
}

// GetAllByProject return all donations associated with project
func (r *DonationRepo) GetAllByProject(id int) ([]Donation, error) {
	donations := make([]Donation, 0)
	err := r.db.Model(&donations).Relation("User").Where("project_id = ", id).Select()
	if err != nil {
		return nil, err
	}

	return donations, nil
}
