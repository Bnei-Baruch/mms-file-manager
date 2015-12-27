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
	"github.com/Bnei-Baruch/mms-file-manager/config"
	"github.com/Bnei-Baruch/mms-file-manager/models"
	"fmt"
)

var (
	l *log.Logger = logger.InitLogger(&logger.LogParams{LogMode: "screen", LogPrefix: "[TASK-DB] "})
)


func SetupArgs(arguments []string) (args []string, err error) {
	var flagErr flag.ErrorHandling
	flags := flag.NewFlagSet("goferFlags", flagErr)

	env := flags.String("env", ".env", "environment")
	if err = flags.Parse(arguments); err != nil {
		return
	}

	if *env == "" {
		return nil, fmt.Errorf("env should not be an empty string")
	}
	godotenv.Load(*env)
	config.CheckEnv()

	args = flags.Args()

	return
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

		args, err := SetupArgs(arguments)
		if err != nil {
			l.Println("Problem with parameters ", err)
		}

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
var DBEmptyDB = gofer.Register(gofer.Task{
	Namespace:   "db",
	Label:       "empty",
	Description: "Delete all tables from DB",
	Action: func(arguments ...string) error {

		if _, err := SetupArgs(arguments); err != nil {
			l.Println("Problem with parameters ", err)
		}

		db := config.NewDB()

		var err error
		if err = db.Exec(`
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
		`).Error; err != nil {
			l.Fatal("Could not empty database.", err)
		}

		return err
	},
})

var DBAutoMigrate = gofer.Register(gofer.Task{
	Namespace:   "db",
	Label:       "automigrate",
	Description: "Auto Migrates a database with gorm",
	Action: func(arguments ...string) error {

		if _, err := SetupArgs(arguments); err != nil {
			l.Println("Problem with parameters ", err)
		}

		db := config.NewDB()

		var err error
		if err = db.AutoMigrate(
			&models.File{},
			&models.Pattern{},
			&models.PatternPart{},
			&models.Workflow{},
			&models.ContentType{},
			&models.Line{},
		).Error; err != nil {
			l.Fatal("Could not automigrate.", err)
		}

		return err
	},
})
var Kuku = gofer.Register(gofer.Task{
	Namespace:   "kuku",
	Label:       "muku",
	Description: "Migrates a database",
	Action: func(arguments ...string) error {
		return gofer.LoadAndPerform("db:migrate", "--env=.env.test")
	},
})
var DBMigrate = gofer.Register(gofer.Task{
	Namespace:   "db",
	Label:       "migrate",
	Dependencies: []string{"db:automigrate"},
	Description: "Migrates a database",
	Action: func(arguments ...string) (err error) {

		if _, err = SetupArgs(arguments); err != nil {
			l.Fatal("Problem with parameters ", err)
		}
		path := filepath.Join(os.Getenv("GOPATH"), "/src/github.com/Bnei-Baruch/mms-file-manager/migrations")
		dbConf := goose.DBConf{
			MigrationsDir: path,
			Env: os.Getenv("ENV"),
			Driver: goose.DBDriver{
				Name: "postgres",
				OpenStr: os.Getenv("DATABASE_URL"),
				Dialect: goose.PostgresDialect{},
				Import: "github.com/lib/pq",
			},
		}

		var target int64
		target, err = goose.GetMostRecentDBVersion(dbConf.MigrationsDir)
		if err != nil {
			l.Fatal(err)
		}
		if err = goose.RunMigrations(&dbConf, dbConf.MigrationsDir, target); err != nil {
			l.Fatal(err)
		}

		return
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
