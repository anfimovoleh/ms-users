package main

import (
	"fmt"

	"github.com/anfimovoleh/ms-users/db"

	"github.com/anfimovoleh/ms-users"
	"github.com/anfimovoleh/ms-users/config"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

type Migrator func(*db.DB, db.MigrateDir, int) (int, error)

func MigrateDB(direction string, count int, dbClient *db.DB, migrator Migrator) (int, error) {
	applied, err := migrator(dbClient, db.MigrateDir(direction), count)
	return applied, errors.Wrap(err, "failed to apply migrations")
}

func main() {
	apiConfig := config.New()
	log := apiConfig.Log()

	rootCmd := &cobra.Command{}

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "run command",
		Run: func(cmd *cobra.Command, args []string) {
			api := app.New(apiConfig)
			if err := api.Start(); err != nil {
				panic(errors.Wrap(err, "failed to start api"))
			}
		},
	}

	migrateCmd := &cobra.Command{
		Use:   "migrate [up|down|redo] [COUNT]",
		Short: "migrate schema",
		Long:  "performs a schema migration command",
		Run: func(cmd *cobra.Command, args []string) {
			log = log.WithField("service", "migration")
			var count int
			// Allow invocations with 1 or 2 args.  All other args counts are erroneous.
			if len(args) < 1 || len(args) > 2 {
				log.WithField("arguments", args).Error("wrong argument count")
				return
			}
			// If a second arg is present, parse it to an int and use it as the count
			// argument to the migration call.
			if len(args) == 2 {
				var err error
				if count, err = cast.ToIntE(args[1]); err != nil {
					log.WithError(err).Error("failed to parse count")
					return
				}
			}

			applied, err := MigrateDB(args[0], count, apiConfig.DB(), db.Migrations.Migrate)
			log = log.WithField("applied", applied)
			if err != nil {
				log.WithError(err).Error("migration failed")
				return
			}
			log.Info("migrations applied")
		},
	}

	rootCmd.AddCommand(runCmd, migrateCmd)
	if err := rootCmd.Execute(); err != nil {
		log.WithField("cobra", "read").Error(fmt.Sprintf("failed to read command %s", err.Error()))
		return
	}
}
