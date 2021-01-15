package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/caarlos0/clone-org/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli"
)

var version = "master"

func main() {
	app := cli.NewApp()
	app.Name = "clone-org"
	app.Usage = "Clone all repos of a github organization"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "org, o",
		},
		cli.StringFlag{
			Name:   "token, t",
			EnvVar: "GITHUB_TOKEN",
		},
		cli.StringFlag{
			Name: "destination, d",
		},
	}
	app.Action = func(c *cli.Context) error {
		log.SetFlags(0)
		f, err := tea.LogToFile("clone-org.log", "")
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		defer func() { _ = f.Close() }()

		token := c.String("token")
		if token == "" {
			return cli.NewExitError("missing github token", 1)
		}

		org := c.String("org")
		if org == "" {
			return cli.NewExitError("missing organization name", 1)
		}

		destination := c.String("destination")
		if destination == "" {
			destination = filepath.Join(os.TempDir(), org)
		}

		var p = tea.NewProgram(ui.NewInitialModel(token, org, destination))
		p.EnterAltScreen()
		defer p.ExitAltScreen()
		if err = p.Start(); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
