package migrate

import (
	"github.com/go-pg/migrations/v8"
	"github.com/labstack/gommon/log"
)

func init() {
	migrations.MustRegister(createProjects, rollbackProjects)
}

func createProjects(db migrations.DB) error {
	log.Info("creating table [projects]...")
	_, err := db.Exec(
		`CREATE TABLE projects (
			id bigserial NOT NULL primary key,
			title varchar NOT NULL,
			sub_title varchar NOT NULL,
			release_date timestamptz NOT NULL,
			event_date timestamptz,
			goal_people int NOT NULL DEFAULT 0,
			goal_amount int NOT NULL DEFAULT 0,
			total int NOT NULL DEFAULT 0,
			description varchar NOT NULL,
			image_link varchar NOT NULL,
			instructions varchar NOT NULL,
			locked boolean NOT NULL DEFAULT FALSE,
			published boolean NOT NULL DEFAULT FALSE,
			closed boolean NOT NULL DEFAULT FALSE,
			owner_id int NOT NULL,
			category_id int NOT NULL,
			project_type_id int NOT NULL
		);		
	`)

	return err
}

func rollbackProjects(db migrations.DB) error {
	log.Warn("dropping table [projects]...")
	_, err := db.Exec(`DROP TABLE projects`)

	return err
}
