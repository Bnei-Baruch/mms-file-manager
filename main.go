package main

import (
	"fmt"
	"os"

	"github.com/Bnei-Baruch/mms-file-manager/cmd"
	"github.com/Bnei-Baruch/mms-file-manager/config"
	"github.com/Bnei-Baruch/mms-file-manager/routes"
	"github.com/Bnei-Baruch/mms-file-manager/services/file_manager"
	"github.com/Bnei-Baruch/mms-file-manager/services/logger"
	"github.com/Bnei-Baruch/mms-file-manager/tasks"
	"github.com/joho/godotenv"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}

func stam() {
	tasks.Run([]string{"rake", "db", "migrate", "--name", "Jeremy"})
	godotenv.Load(".env")
	app := config.NewApp(".")
	routes.Setup(app)
	watchDir, targetDir := "tmp/source", "tmp/target"

	file_manager.Logger(&logger.LogParams{LogMode: "screen", LogPrefix: "[FM] "})

	fm, err := file_manager.NewFM(targetDir, file_manager.WatchPair{watchDir, "main"})
	if err != nil {
		panic(err)
	}
	defer fm.Destroy()

	app.Negroni.Run(":" + os.Getenv("PORT"))

}
