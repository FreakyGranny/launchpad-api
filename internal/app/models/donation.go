package models

import (
	"github.com/go-pg/pg/v10"
)

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

// Create a new donation
func (r *DonationRepo) Create(d *Donation) error {
	// NEED IMPLEMENTATION

	// if err := dbClient.First(&project, request.ProjectID).Error; gorm.IsRecordNotFoundError(err) {
	// 	return c.JSON(http.StatusBadRequest, nil)
	// }
	// if project.Closed || project.Locked || !project.Published {
	// 	return c.JSON(http.StatusBadRequest, nil)
	// }

	// var donationCount uint
	// dbClient.Model(&db.Donation{}).Where("project_id = ? AND user_id = ?", request.ProjectID, uint(userID)).Count(&donationCount)
	// if donationCount > 0 {
	// 	return c.JSON(http.StatusForbidden, nil)
	// }
	return nil
}

// Delete not locked donation
func (r *DonationRepo) Delete(id int, userID int) error {
	// NEED IMPLEMENTATION

	// dbClient := db.GetDbClient()
	// var donation db.Donation

	// if err := dbClient.Preload("User").First(&donation, donationID).Error; gorm.IsRecordNotFoundError(err) {
	// 	return c.JSON(http.StatusNotFound, nil)
	// }
	// if donation.Locked || donation.User.ID != uint(userID) {
	// 	return c.JSON(http.StatusForbidden, nil)
	// }
	// dbClient.Delete(&donation)

	return nil
}

// Update not locked donation
func (r *DonationRepo) Update(d *Donation) error {
	// NEED IMPLEMENTATION

	// dbClient := db.GetDbClient()
	// var donation db.Donation

	// if err := dbClient.Preload("Project").Preload("User").First(&donation, donationID).Error; gorm.IsRecordNotFoundError(err) {
	// 	return c.JSON(http.StatusNotFound, nil)
	// }

	// if donation.Locked && donation.Project.OwnerID == uint(userID) {
	// 	dbClient.Model(&donation).Update("paid", request.Paid)
	// 	ch := misc.GetHarvestPipe()
	// 	ch <- donation.ProjectID

	// 	return c.JSON(http.StatusOK, donation)
	// }
	// if !donation.Locked && donation.User.ID == uint(userID) {
	// 	dbClient.Model(&donation).Update("payment", request.Payment)
	// 	ch := misc.GetRecalcPipe()
	// 	ch <- donation.ProjectID

	// 	return c.JSON(http.StatusOK, donation)
	// }

	return nil
}
