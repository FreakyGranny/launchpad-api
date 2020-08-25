package migrate

import (
	"time"

	"github.com/FreakyGranny/launchpad-api/internal/app/models"
	"github.com/go-pg/migrations/v8"
	"github.com/labstack/gommon/log"
)

func init() {
	migrations.MustRegister(createSystem, rollbackSystem)
}

func createSystem(db migrations.DB) error {
	log.Info("creating table [system]...")
	_, err := db.Exec(
		`CREATE TABLE system (
			id integer UNIQUE default(1),
			last_check timestamptz NOT NULL
			Constraint CHK_System_singlerow CHECK (id = 1)
		);		
	`)
	if err != nil {
		return err
	}
	n := time.Now()
	settings := models.System{LastCheck: time.Date(n.Year(), n.Month(), n.Day(), 6, 0, 0, n.Nanosecond(), n.Location())}
	_, err = db.Model(&settings).Insert()

	return err
}

func rollbackSystem(db migrations.DB) error {
	log.Warn("dropping table [system]...")
	_, err := db.Exec(`DROP TABLE system`)

	return err
}
