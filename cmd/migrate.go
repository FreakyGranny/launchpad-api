package cmd

import (
	"fmt"
	"os"

	"github.com/labstack/gommon/log"

	"github.com/FreakyGranny/launchpad-api/internal/config"
	"github.com/FreakyGranny/launchpad-api/internal/db"
	_ "github.com/FreakyGranny/launchpad-api/migrations" //
	"github.com/go-pg/migrations/v8"
	"github.com/spf13/cobra"
)

const usageText = `Supported commands are:
  - init - creates version info table in the database
  - up - runs all available migrations.
  - up [target] - runs available migrations up to the target one.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.
`

// NewMigrateCmd returns migrate cmd struct
func NewMigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "run migrations",
		Long:  usageText,
		Run:   Migrate,
	}
}

// Migrate ...
func Migrate(cmd *cobra.Command, args []string) {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	d, err := db.Connect(&cfg.Db)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	oldVersion, newVersion, err := migrations.Run(d, args...)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if newVersion != oldVersion {
		log.Info(fmt.Sprintf("migrated from version %d to %d", oldVersion, newVersion))
	} else {
		log.Warn(fmt.Sprintf("version is %d", oldVersion))
	}
}
