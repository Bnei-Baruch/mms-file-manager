package cmd

import (
	"os"
	"path/filepath"
	"time"

	"github.com/Bnei-Baruch/mms-file-manager/config"
	"github.com/Bnei-Baruch/mms-file-manager/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/pressly/goose"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(dbGenerate)
	RootCmd.AddCommand(dbEmpty)
	RootCmd.AddCommand(dbMigrate)
	dbGenerate.PersistentFlags().StringVar(&migrationType, "type", "go", "migration type (go|sql)")
}

func RegisterDb(gdb *gorm.DB) {
	db = gdb
}

func LoadEnv(loadDB bool) {
	godotenv.Load(Env)
	config.CheckEnv()

	if loadDB {
		db = config.NewDB()
	}
}

var db *gorm.DB
var migrationType string

var dbGenerate = &cobra.Command{
	Use:   "db:generate migration-name",
	Short: "Generate new migration file",
	Run: func(cmd *cobra.Command, args []string) {

		if db == nil {
			LoadEnv(false)
		}

		/*
			dbConf := goose.DBConf{
				MigrationsDir: "migrations",
				Driver: goose.DBDriver{
					Name: "postgres",
					DSN:  "$DATABASE_URL",
				},
			}
		*/

		if len(args) < 1 {
			l.Fatal("create: migration name required")
		}

		if err := os.MkdirAll("migrations", 0777); err != nil {
			l.Fatal(err)
		}

		n, err := goose.CreateMigration(args[0], migrationType, "migrations", time.Now())
		if err != nil {
			l.Fatal(err)
		}

		a, e := filepath.Abs(n)
		if e != nil {
			l.Fatal(e)
		}

		l.Println("created:", a)
	},
}

var dbEmpty = &cobra.Command{
	Use:   "db:empty",
	Short: "Empty the database",
	Run: func(cmd *cobra.Command, args []string) {
		if db == nil {
			LoadEnv(true)
		}

		if err := db.Exec(`
			DROP SCHEMA public CASCADE;
			CREATE SCHEMA public;
		`).Error; err != nil {
			l.Fatal("Could not empty database.", err)
		}
	},
}

var dbMigrate = &cobra.Command{
	Use:   "db:migrate",
	Short: "Migrate the databse to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		if db == nil {
			LoadEnv(true)
		}

		var err error
		if err = db.AutoMigrate(
			&models.Line{},
			&models.ContentType{},
			&models.File{},
			&models.Pattern{},
			&models.PatternPart{},
			&models.Workflow{},
		).Error; err != nil {
			l.Fatal("Could not automigrate.", err)
			return
		}

		path := filepath.Join(os.Getenv("GOPATH"), "/src/github.com/Bnei-Baruch/mms-file-manager/migrations")
		/*
			dbConf := goose.DBConf{
				MigrationsDir: path,
				Driver: goose.DBDriver{
					Name:    "postgres",
					DSN:     os.Getenv("DATABASE_URL"),
					Dialect: goose.PostgresDialect{},
					Import:  "github.com/lib/pq",
				},
			}
		*/

		var target int64
		target, err = goose.GetMostRecentDBVersion(path)
		if err != nil {
			l.Fatal(err)
		}
		if err = goose.RunMigrations(db.DB(), path, target); err != nil {
			l.Fatal(err)
		}
	},
}
