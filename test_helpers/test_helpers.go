package test_helpers

import (
	"fmt"

	"github.com/Bnei-Baruch/mms-file-manager/cmd"
	"github.com/Bnei-Baruch/mms-file-manager/config"
	"github.com/Bnei-Baruch/mms-file-manager/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func SetupSpec() (db *gorm.DB) {
	cmd.RootCmd.SetArgs([]string{"db:empty", "--env=../../.env.test"})
	if err := cmd.RootCmd.Execute(); err != nil {
		panic(fmt.Sprintf("Unable to empty database %v", err))
	}

	cmd.RootCmd.SetArgs([]string{"db:migrate", "--env=../../.env.test"})
	if err := cmd.RootCmd.Execute(); err != nil {
		panic(fmt.Sprintf("Unable to migrate database %v", err))
	}

	godotenv.Load("../../.env.test")
	db = config.NewDB()
	models.New(db)
	return
}
