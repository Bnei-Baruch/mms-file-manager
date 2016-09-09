package tasks

import (
	"github.com/Bnei-Baruch/mms-file-manager/config"
	"github.com/Bnei-Baruch/mms-file-manager/models"
	"github.com/urfave/cli"
	"github.com/joho/godotenv"
)

var rake = register()

func Run(args []string) {
	// set args for examples sake
	//os.Args = []string{"rake", "db", "migrate", "--name", "Jeremy"}
	rake.Run(args)
}

func register() *cli.App {
	app := cli.NewApp()
	app.Name = "rake"
	app.Commands = []cli.Command{
		{
			Name:        "db",
			Usage:       "use it to see a description",
			Description: "This is how we describe hello the function",
			Subcommands: []cli.Command{
				{
					Name:        "migrate",
					Usage:       "sends a greeting in english",
					Description: "greets someone in english",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "env",
							Value: ".env",
							Usage: "Name of the person to greet",
						},
					},
					Action: func(c *cli.Context) {
						var env = c.String("env")
						godotenv.Load(env)
						config.CheckEnv()

						//return automigrate()
					},
				},
			},
		},
	}

	return app
}
func automigrate() error {

	db := config.NewDB()
	defer db.Close()

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
	}

	return err
}
