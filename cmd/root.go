package cmd

import (
	"log"
	"os"

	"github.com/Bnei-Baruch/mms-file-manager/config"
	"github.com/Bnei-Baruch/mms-file-manager/routes"
	"github.com/Bnei-Baruch/mms-file-manager/services/file_manager"
	"github.com/Bnei-Baruch/mms-file-manager/services/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var Env string

var (
	l *log.Logger = logger.InitLogger(&logger.LogParams{LogMode: "screen", LogPrefix: "[CLI] "})
)

var RootCmd = &cobra.Command{
	Use:   "mms-file-manager",
	Short: "mms file manager is application to manage media processes",
	Run: func(cmd *cobra.Command, args []string) {
		godotenv.Load(Env)
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

	},
}

func init() {
	//RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	//RootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
	RootCmd.PersistentFlags().StringVar(&Env, "env", ".env", "environment variables file")
	//RootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
	//RootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
}
