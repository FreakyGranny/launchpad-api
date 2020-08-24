package models

import (
	"github.com/go-pg/pg/v10"
)

//go:generate mockgen -source=$GOFILE -destination=../mocks/model_donation_mock.go -package=mocks DonationImpl

// DonationImpl ...
type DonationImpl interface {
	Get(id int) (*Donation, bool)
	GetAllByUser(id int) ([]Donation, error)
	GetAllByProject(id int) ([]Donation, error)
	Create(d *Donation) error
	Update(d *Donation) error
	Delete(d *Donation) error
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

// Get donation
func (r *DonationRepo) Get(id int) (*Donation, bool) {
	donation := &Donation{}
	err := r.db.Model(donation).Relation("Project").Where("d.id = ?", id).Select()
	if err != nil {
		return nil, false
	}

	return donation, true
}

// GetAllByUser returns all donations associated with user
func (r *DonationRepo) GetAllByUser(id int) ([]Donation, error) {
	donations := make([]Donation, 0)
	err := r.db.Model(&donations).Where("user_id = ?", id).Select()
	if err != nil {
		return nil, err
	}

	return donations, nil
}

// GetAllByProject return all donations associated with project
func (r *DonationRepo) GetAllByProject(id int) ([]Donation, error) {
	donations := make([]Donation, 0)
	err := r.db.Model(&donations).Relation("User").Where("d.project_id = ?", id).Select()
	if err != nil {
		return nil, err
	}

	return donations, nil
}

// Create a new donation
func (r *DonationRepo) Create(d *Donation) error {
	count, err := r.db.Model((*Donation)(nil)).Where("d.project_id = ? AND d.user_id = ?", d.ProjectID, d.UserID).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrDonationAlreadyExist
	}
	count, err = r.db.Model((*User)(nil)).Where("u.id = ?", d.UserID).Count()
	if err != nil {
		return err
	}
	if count != 1 {
		return ErrUserNotFound
	}
	project := &Project{}
	err = r.db.Model(project).Where("p.id = ?", d.ProjectID).Select()
	if err != nil {
		return err
	}
	if project.Closed || project.Locked || !project.Published {
		return ErrDonationForbidden
	}
	_, err = r.db.Model(d).Insert()
	if err != nil {
		return err
	}

	return nil
}

// Delete not locked donation
func (r *DonationRepo) Delete(d *Donation) error {
	_, err := r.db.Model(d).WherePK().Delete()

	return err
}

// Update donation
func (r *DonationRepo) Update(d *Donation) error {
	_, err := r.db.Model(d).WherePK().UpdateNotZero()

	return err

}
