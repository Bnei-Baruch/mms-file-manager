package main

import (
	"github.com/Bnei-Baruch/mms-file-manager/config"
	"github.com/Bnei-Baruch/mms-file-manager/routes"
	"github.com/Bnei-Baruch/mms-file-manager/services/file_manager"
	"os"
	"github.com/Bnei-Baruch/mms-file-manager/services/logger"
)

func main() {
	app := config.NewApp(".")
	routes.Setup(app)
	watchDir, targetDir := "tmp/source", "tmp/target"

	file_manager.Logger(&logger.LogParams{LogMode: "screen", LogPrefix: "[FM] "})

	fm, err := file_manager.NewFM()
	if err != nil {
		panic(err)
	}
	defer fm.Destroy()

	fm.Watch(watchDir, targetDir)
	app.Negroni.Run(":" + os.Getenv("PORT"))

}