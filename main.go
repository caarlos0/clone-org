package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/caarlos0/spin"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

func main() {
	app := cli.NewApp()
	app.Name = "clone-org"
	app.Usage = "Clone all repos of a github organization"
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
		client := client(c.String("token"))
		org := c.String("org")
		if org == "" {
			return cli.NewExitError("Missing organization name", 1)
		}
		destination := c.String("destination")
		if destination == "" {
			destination = "/tmp/" + org
		}
		s := spin.New("%s Finding repositories to clone...")
		s.Start()
		repos, err := findRepos(org, client)
		s.Stop()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		if err := os.Mkdir(destination, 0700); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		s = spin.New(
			fmt.Sprintf("%s Cloning %d repositories...", "%s", len(repos)),
		)
		s.Start()
		defer s.Stop()
		var wg sync.WaitGroup
		wg.Add(len(repos))
		for _, repo := range repos {
			go func(repo *github.Repository) {
				name := *repo.Name
				dest := destination + "/" + name
				url := *repo.SSHURL
				cmd := exec.Command("git", "clone", "--depth", "1", url, dest)
				if bts, err := cmd.CombinedOutput(); err != nil {
					fmt.Printf("\ngit clone failed for %s: %s", url, string(bts))
				}
				wg.Done()
			}(repo)
		}
		wg.Wait()
		return nil
	}

	app.Run(os.Args)
}

func client(token string) *github.Client {
	return github.NewClient(
		oauth2.NewClient(
			oauth2.NoContext,
			oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
		),
	)
}

func findRepos(org string, client *github.Client) (result []*github.Repository, err error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 30},
	}

	for {
		repos, resp, err := client.Repositories.ListByOrg(org, opt)
		if err != nil {
			return result, err
		}
		result = append(result, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	return result, nil
}
