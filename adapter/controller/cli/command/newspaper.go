package command

import (
	"sort"

	"github.com/urfave/cli/v2"

	"newspaper-api/adapter/controller/cli/action"
	"newspaper-api/adapter/controller/cli/presenter"
	"newspaper-api/pkg/logger"
)

var NewspaperTitle string

func SetNewspaperCommand(app *cli.App, newspaperAction *action.NewspaperAction) {
	cliFlag := []cli.Flag{
		&cli.StringFlag{
			Name:        "newspaper_title",
			Aliases:     []string{"a"},
			Usage:       "Title for the newspaper",
			Destination: &NewspaperTitle,
		},
	}
	app.Flags = append(app.Flags, cliFlag...)

	cliCommand := []*cli.Command{
		{
			Name:    "newspaper",
			Aliases: []string{"a"},
			Usage:   "Select a newspaper",
			Subcommands: []*cli.Command{
				{
					Name:  "create",
					Usage: "Create for newspaper",
					Action: func(c *cli.Context) error {
						newspaper, err := newspaperAction.CreateNewspaper(NewspaperTitle)
						if err != nil {
							logger.Error(err.Error())
							return err
						}
						presenter.PrettyPrintStructToJson(newspaper)
						return nil
					},
				},
			},
		},
	}
	app.Commands = append(app.Commands, cliCommand...)

	sort.Sort(cli.CommandsByName(app.Commands))
	sort.Sort(cli.FlagsByName(app.Flags))
}
