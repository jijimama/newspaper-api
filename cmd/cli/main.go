package main

import (
	"os"
	"sort"

	"github.com/urfave/cli/v2"

	"newspaper-api/adapter/controller/cli/action"
	"newspaper-api/adapter/controller/cli/command"
	"newspaper-api/adapter/gateway"
	"newspaper-api/infrastructure/database"
	"newspaper-api/pkg/logger"
	"newspaper-api/usecase"
)

func main() {
	db, err := database.NewDatabaseSQLFactory(database.InstanceMySQL)
	if err != nil {
		logger.Fatal(err.Error())
	}

	newspaperRepository := gateway.NewNewspaperRepository(db)
	newspaperUseCase := usecase.NewNewspaperUseCase(newspaperRepository)
	newspaperAction := action.NewNewspaperAction(newspaperUseCase)

	app := &cli.App{}
	command.SetNewspaperCommand(app, newspaperAction)
	sort.Sort(cli.CommandsByName(app.Commands))
	sort.Sort(cli.FlagsByName(app.Flags))
	err = app.Run(os.Args)
	if err != nil {
		logger.Fatal(err.Error())
	}
}
