package main

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"

	cloneorg "github.com/caarlos0/clone-org"
	"github.com/caarlos0/spin"
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
		fmt.Printf("Destination: %v\n", destination)

		s := spin.New("%s Finding repositories to clone...")
		s.Start()
		repos, err := cloneorg.AllOrgRepos(token, org)
		s.Stop()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		if err := cloneorg.CreateDir(destination); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		fmt.Printf("Cloning %v repositories:\n", len(repos))
		var g errgroup.Group
		for _, repo := range repos {
			repo := repo
			g.Go(func() error {
				err := cloneorg.Clone(repo, destination)
				if err != nil {
					fmt.Printf("\033[33m[failed] %s\033[0m: %s", repo.Name, err)
				} else {
					fmt.Printf("\033[32m[cloned] %s\033[0m\n", repo.Name)
				}
				return nil
			})
		}
		return g.Wait()
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
