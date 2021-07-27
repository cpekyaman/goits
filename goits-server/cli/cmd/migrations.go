package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/cpekyaman/goits/config"
	"github.com/cpekyaman/goits/framework/orm/db"
	migrate "github.com/rubenv/sql-migrate"
)

type dbConfig struct {
	Migrations migrationConfig `mapstructure:"migrations"`
}

type migrationConfig struct {
	Platform      string `mapstructure:"platform"`
	SchemaName    string `mapstructure:"schema"`
	TableName     string `mapstructure:"table"`
	MigrationsDir string `mapstructure:"dir"`
	Dialect       string `mapstructure:"dialect"`
}

var mconf dbConfig
var MigrateCommand *cobra.Command

func init() {
	MigrateCommand = &cobra.Command{
		Use:     "migrate",
		Aliases: []string{"dbmigrate", "dbm"},
		Short:   "run db migrations",
		Long:    "applies db migrations to the database",
	}
	up := &cobra.Command{
		Use:   "up",
		Short: "apply migrations",
		Long:  "apply migrations for schema upgrade",
		Run: func(cmd *cobra.Command, args []string) {
			dbMigrate(migrate.Up)
		},
	}
	down := &cobra.Command{
		Use:   "down",
		Short: "rollback migrations",
		Long:  "apply migrations for schema downgrade",
		Run: func(cmd *cobra.Command, args []string) {
			dbMigrate(migrate.Down)
		},
	}
	MigrateCommand.AddCommand(up, down)
}

// dbMigrate runs the db migration scripts for new version or rolling back to previous version.
func dbMigrate(direction migrate.MigrationDirection) {
	config.ReadInto("db", &mconf)

	migrations := &migrate.FileMigrationSource{
		Dir: mconf.Migrations.MigrationsDir + "/" + mconf.Migrations.Platform,
	}

	ms := migrate.MigrationSet{
		TableName:  mconf.Migrations.TableName,
		SchemaName: mconf.Migrations.SchemaName,
	}

	log.Printf("Applying migrations for %v\n", direction)
	n, err := ms.Exec(db.DB().DB, mconf.Migrations.Dialect, migrations, direction)
	if err != nil {
		log.Fatal("Could not execute migrations : ", err)
	}
	log.Printf("Applied %d migrations !\n", n)
}
