package app

import "github.com/FreakyGranny/launchpad-api/internal/models"

// ExtendedUser user extended with participation.
type ExtendedUser struct {
	models.User
	Participation []models.Participation `json:"participation"`
}

// ExtendedProject light project entry
type ExtendedProject struct {
	ID           int                `json:"id"`
	Title        string             `json:"title"`
	SubTitle     string             `json:"subtitle"`
	Status       string             `json:"status"`
	ReleaseDate  string             `json:"release_date"`
	EventDate    *string            `json:"event_date"`
	ImageLink    string             `json:"image_link"`
	Total        int                `json:"total"`
	Percent      int                `json:"percent"`
	Category     models.Category    `json:"category"`
	ProjectType  models.ProjectType `json:"project_type"`
	GoalPeople   int                `json:"goal_people"`
	GoalAmount   int                `json:"goal_amount"`
	Description  string             `json:"description"`
	Instructions string             `json:"instructions"`
	Owner        models.User        `json:"owner"`
}


// ShortDonation project donation without payment
type ShortDonation struct {
	ID     int         `json:"id"`
	User   models.User `json:"user"`
	Locked bool        `json:"locked"`
	Paid   bool        `json:"paid"`
}
