package migrate

import (
	"github.com/FreakyGranny/launchpad-api/internal/app/models"
	"github.com/go-pg/migrations/v8"
	"github.com/labstack/gommon/log"
)

func init() {
	migrations.MustRegister(createProjectTypes, rollbackProjectTypes)
}

func createProjectTypes(db migrations.DB) error {
	log.Info("creating table [project_types]...")
	_, err := db.Exec(
		`CREATE TABLE project_types (
			id bigserial NOT NULL primary key,
			alias varchar NOT NULL,
			name varchar NOT NULL,
			options TEXT [],
			goal_by_amount boolean NOT NULL DEFAULT FALSE,
			goal_by_people boolean NOT NULL DEFAULT FALSE,
			end_by_goal_gain boolean NOT NULL DEFAULT FALSE);
	`)
	if err != nil {
		return err
	}

	categories := []models.ProjectType{
		{
			Alias: "money_fast",
			Name:  "Сampaign",
			Options: []string{
				"Partakers deposit an arbitrary amount.",
				"When the required amount is reached, the campaign stops.",
				"Campaign is considered successful when the author has marked all money transfers.",
			},
			GoalByAmount:  true,
			GoalByPeople:  false,
			EndByGoalGain: true,
		},
		{
			Alias: "money_equal",
			Name:  "Fair campaign",
			Options: []string{
				"Partakers agree to split the amount among themselves.",
				"The minimum number of partakers must be recruited.",
				"The number of partakers is not limited.",
				"Fundraising starts on the date specified by the author.",
				"Fair campaign is considered successful when the author has marked all money transfers.",
			},
			GoalByAmount:  true,
			GoalByPeople:  true,
			EndByGoalGain: false,
		},
		{
			Alias: "event_fast",
			Name:  "Event",
			Options: []string{
				"Partakers agree to participate in the event.",
				"Event is considered successful when the required number of partakers is reached.",
			},
			GoalByAmount:  false,
			GoalByPeople:  true,
			EndByGoalGain: true,
		},
		{
			Alias: "event_overflow",
			Name:  "Event+",
			Options: []string{
				"Partakers agree to participate in the event.",
				"The number of partakers is not limited.",
				"Event+ considered successful if a sufficient number of people have gathered on the event date.",
			},
			GoalByAmount:  false,
			GoalByPeople:  true,
			EndByGoalGain: false,
		},
	}

	_, err = db.Model(&categories).Insert()

	return err
}

func rollbackProjectTypes(db migrations.DB) error {
	log.Warn("dropping table [project_types]...")
	_, err := db.Exec(`DROP TABLE project_types`)

	return err
}
