package migrate

import (
	"github.com/go-pg/migrations/v8"
	"github.com/labstack/gommon/log"
)

func init() {
	migrations.MustRegister(createUsers, rollbackUsers)
}

func createUsers(db migrations.DB) error {
	log.Info("creating table [users]...")
	_, err := db.Exec(
		`CREATE TABLE users (
			id bigserial NOT NULL primary key,
			username varchar NOT NULL,
			first_name varchar NOT NULL,
			last_name varchar NOT NULL,
			avatar varchar NOT NULL,
			email varchar NOT NULL unique,
			is_admin boolean NOT NULL DEFAULT FALSE,
			project_count int NOT NULL DEFAULT 0,
			success_rate numeric NOT NULL DEFAULT 0.0
		);		
	`)

	return err
}

func rollbackUsers(db migrations.DB) error {
	log.Warn("dropping table [users]...")
	_, err := db.Exec(`DROP TABLE users`)

	return err
}
