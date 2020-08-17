package db

import (
	"fmt"

	"github.com/FreakyGranny/launchpad-api/internal/app/config"
	"github.com/go-pg/pg/v10"
)

// Connect ...
func Connect(cfg *config.PgConnection) (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%v", cfg.Host, cfg.Port),
		User:     cfg.Username,
		Password: cfg.Password,
		Database: cfg.DbName,
	})
	if err := db.Ping(db.Context()); err != nil {
		return nil, err
	}
	return db, nil
}
