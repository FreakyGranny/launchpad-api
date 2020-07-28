package db

import (
	"fmt"

	"github.com/FreakyGranny/launchpad-api/internal/app/config"
	"github.com/jmoiron/sqlx"
)

func sslMode(sslEnable bool) string {
	if sslEnable {
		return "enable"
	}

	return "disable"
}

// Connect ...
func Connect(cfg config.PgConnection) (*sqlx.DB, error) {
	connectString := fmt.Sprintf(
		"host=%s port=%v user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.DbName,
		cfg.Password,
		sslMode(cfg.SslEnable),
	)
	db, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
