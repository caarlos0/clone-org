package main

import (
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/google/go-github/github"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

func main() {
	app := cli.NewApp()
	app.Name = "clone-org"
	app.Usage = "Clone all repos of a github's user or organization"
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
		destination := c.String("destination")
		if destination == "" {
			destination = "/tmp/" + org
		}
		repos, err := findRepos(org, client)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		if err := os.Mkdir(destination, 0700); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		var wg sync.WaitGroup
		wg.Add(len(repos))
		for _, repo := range repos {
			go func(repo *github.Repository) {
				name := *repo.Name
				dest := destination + "/" + name
				url := *repo.SSHURL
				log.Println("Cloning", name, "into", dest)
				cmd := exec.Command("git", "clone", "--depth", "1", url, dest)
				if bts, err := cmd.CombinedOutput(); err != nil {
					log.Println("git clone failed for", url, string(bts))
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
