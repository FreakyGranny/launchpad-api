package models

import (
	"errors"

	"github.com/go-pg/pg/v10"
)

// ErrDonationAlreadyExist donation to project for user already exist
var ErrDonationAlreadyExist = errors.New("donation to this project already exist")

// ErrDonationForbidden donation to project is not allowed
var ErrDonationForbidden = errors.New("donation to project is not allowed")

// ErrDonationModifyForbidden donation editing is not allowed
var ErrDonationModifyForbidden = errors.New("donation editing is not allowed")

// ErrUserNotFound user not found
var ErrUserNotFound = errors.New("user not found")

//go:generate mockgen -source=$GOFILE -destination=../mocks/model_donation_mock.go -package=mocks DonationImpl

// DonationImpl ...
type DonationImpl interface {
	GetAllByUser(id int) ([]Donation, error)
	GetAllByProject(id int) ([]Donation, error)
	Create(d *Donation) error
	Update(d *Donation) error
	Delete(id int, userID int) error
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
func (r *DonationRepo) Delete(id int, userID int) error {
	donation := &Donation{}
	err := r.db.Model(donation).Where("d.id = ?", id).Select()
	if err != nil {
		return err
	}
	if donation.Locked || donation.UserID != userID {
		return ErrDonationModifyForbidden
	}
	_, err = r.db.Model(donation).WherePK().Delete()
	if err != nil {
		return err
	}

	return nil
}

// Update not locked donation
func (r *DonationRepo) Update(d *Donation) error {
	donation := &Donation{}
	err := r.db.Model(donation).Relation("Project").Where("d.id = ?", d.ID).Select()
	if err != nil {
		return err
	}
	if donation.Locked && donation.Project.OwnerID == d.UserID {
		_, err = r.db.Model(d).Set("paid = ?paid").Where("id = ?id").Update()
		return err
	}
	if !donation.Locked && donation.UserID == d.UserID {
		_, err = r.db.Model(d).Set("payment = ?payment").Where("id = ?id").Update()
		return err
	}

	return ErrDonationModifyForbidden
}
