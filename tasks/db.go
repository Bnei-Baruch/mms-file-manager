package tasks

import (
	"github.com/chuckpreslar/gofer"
	"github.com/joho/godotenv"
	"bitbucket.org/liamstask/goose/lib/goose"
	"github.com/DavidHuie/gomigrate"
	_ "github.com/lib/pq"
	"database/sql"
	"errors"
	"os"
	"flag"
	"log"
	"github.com/Bnei-Baruch/mms-file-manager/services/logger"
	"time"
	"path/filepath"
)

var l *log.Logger = logger.InitLogger(&logger.LogParams{LogMode: "screen", LogPrefix: "[TASK-DB] "})


func SetupEnv(env *string) error {
	if *env != "" {
		return godotenv.Load(*env)
	}
	return nil
}

func DBConnect() (*sql.DB, error) {
	db, dbError := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	return db, dbError
}

func SetupMigrator() (*gomigrate.Migrator, error) {
	db, dbError := DBConnect()
	if dbError != nil {
		return nil, errors.New("DB connection failed")
	}
	migrator, migError := gomigrate.NewMigrator(db, gomigrate.Postgres{}, "./migrations")
	return migrator, migError
}

var DBGenerate = gofer.Register(gofer.Task{
	Namespace:   "db",
	Label:       "generate",
	Description: "Generate new migration",
	Action: func(arguments ...string) error {

		var flagErr flag.ErrorHandling
		flags := flag.NewFlagSet("goferFlags", flagErr)

		env := flags.String("env", ".env", "environment")
		if loadError := SetupEnv(env); loadError != nil {
			l.Println("env file does not exist: ", env)
		}

		flags.Parse(arguments)
		args := flags.Args()

		dbConf := goose.DBConf{
			MigrationsDir: "migrations",
			Env: "production",
			Driver: goose.DBDriver{
				Name: "postgres",
				OpenStr: "$DATABASE_URL",

			},
		}

		if len(args) < 1 {
			l.Fatal("create: migration name required")
		}

		migrationType := "go" // default to Go migrations
		if len(args) >= 2 {
			migrationType = args[1]
		}

		if err := os.MkdirAll(dbConf.MigrationsDir, 0777); err != nil {
			l.Fatal(err)
		}

		n, err := goose.CreateMigration(args[0], migrationType, dbConf.MigrationsDir, time.Now())
		if err != nil {
			l.Fatal(err)
		}

		a, e := filepath.Abs(n)
		if e != nil {
			l.Fatal(e)
		}

		l.Println("created:", a)

		return nil
	},
})
var DBMigrate = gofer.Register(gofer.Task{
	Namespace:   "db",
	Label:       "migrate",
	Description: "Migrates a database",
	Action: func(arguments ...string) error {

		env := flag.String("env", ".env", "environment")
		if loadError := SetupEnv(env); loadError != nil {
			l.Println("env file does not exist: ", env)
		}

		dbConf := goose.DBConf{
			MigrationsDir: "migrations",
			Env: "production",
			Driver: goose.DBDriver{
				Name: "postgres",
				OpenStr: "$DATABASE_URL",

			},
		}

		target, err := goose.GetMostRecentDBVersion(dbConf.MigrationsDir)
		if err != nil {
			l.Fatal(err)
		}

		if err := goose.RunMigrations(&dbConf, dbConf.MigrationsDir, target); err != nil {
			l.Fatal(err)
		}

		return nil
	},
})

var DBRollback = gofer.Register(gofer.Task{
	Namespace:   "db",
	Label:       "rollback",
	Description: "Rolls back a database",
	Action: func(arguments ...string) error {

		/*
				loadError := SetupEnv(arguments)
				if loadError != nil {
					return errors.New("env file does not exist")
				}
		*/

		migrator, migError := SetupMigrator()
		if migError != nil {
			return migError
		}

		migrateError := migrator.Rollback()
		if migrateError != nil {
			return errors.New("Migration failed")
		}

		return nil
	},
})