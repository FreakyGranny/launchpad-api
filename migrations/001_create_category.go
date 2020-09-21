package migrate

import (
	"github.com/FreakyGranny/launchpad-api/internal/models"
	"github.com/go-pg/migrations/v8"
	"github.com/labstack/gommon/log"
)

func init() {
	migrations.MustRegister(createCategories, rollbackCategories)
}

func createCategories(db migrations.DB) error {
	log.Info("creating table [categories]...")
	_, err := db.Exec(
		`CREATE TABLE categories (
			id bigserial NOT NULL primary key,
			alias varchar NOT NULL,
			name varchar NOT NULL);
	`)
	if err != nil {
		return err
	}

	categories := []models.Category{
		{
			Alias: "default",
			Name:  "Other",
		},
		{
			Alias: "video_games",
			Name:  "Video games",
		},
		{
			Alias: "board_games",
			Name:  "Board games",
		},
		{
			Alias: "party",
			Name:  "Party",
		},
	}
	_, err = db.Model(&categories).Insert()

	return err
}

func rollbackCategories(db migrations.DB) error {
	log.Warn("dropping table [categories]...")
	_, err := db.Exec(`DROP TABLE categories`)

	return err
}
