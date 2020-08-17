package migrate

import (
	"github.com/go-pg/migrations/v8"
	"github.com/labstack/gommon/log"
)

func init() {
	migrations.MustRegister(createDonations, rollbackDonations)
}

func createDonations(db migrations.DB) error {
	log.Info("creating table [donations]...")
	_, err := db.Exec(
		`CREATE TABLE donations (
			id bigserial NOT NULL primary key,
			payment int NOT NULL DEFAULT 0,
			locked boolean NOT NULL DEFAULT FALSE,
			paid boolean NOT NULL DEFAULT FALSE,
			user_id int NOT NULL,
			project_id int NOT NULL);
	`)

	return err
}

func rollbackDonations(db migrations.DB) error {
	log.Warn("dropping table [categories]...")
	_, err := db.Exec(`DROP TABLE categories`)

	return err
}
