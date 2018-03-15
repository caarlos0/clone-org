// Package cloneorg contains useful functions to find and clone a github
// organization repositories.
package cloneorg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Repo represents the repository data we need.
type Repo struct {
	Name string
	URL  string
}

// Clone a given repository into a given destination
func Clone(repo Repo, destination string) error {
	var cmd = exec.Command(
		"git", "clone", "--depth", "1", repo.URL,
		filepath.Join(destination, repo.Name),
	)
	if bts, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone failed for %v: %v", repo.Name, string(bts))
	}
	return nil
}

// AllOrgRepos finds all repositories of a given organization
func AllOrgRepos(token, org string) (repos []Repo, err error) {
	var ctx = context.Background()
	var client = github.NewClient(oauth2.NewClient(
		ctx,
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
	))
	result, err := findRepos(ctx, client, org)
	if err != nil {
		return
	}
	for _, repo := range result {
		repos = append(repos, Repo{
			Name: *repo.Name,
			URL:  *repo.SSHURL,
		})
	}
	return
}

func findRepos(ctx context.Context, client *github.Client, org string) (result []*github.Repository, err error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 30},
	}
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, opt)
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

// CreateDir creates the directory if it does not exists
func CreateDir(directory string) error {
	stat, err := os.Stat(directory)
	directoryDoesNotExists := err != nil

	if directoryDoesNotExists {
		err := os.MkdirAll(directory, 0700)
		if err != nil {
			return fmt.Errorf("couldn't create directory: %v", err)
		}

		return nil
	}

	if stat.IsDir() {
		return nil
	}

	return fmt.Errorf("directory provided is a file: %v", directory)
}
