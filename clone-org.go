// Package cloneorg contains useful functions to find and clone a github
// organization repositories.
package cloneorg

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

// Repo represents the repository data we need.
type Repo struct {
	Name string
	URL  string
}

// ErrClone happens when a git clone fails.
var ErrClone = errors.New("git clone failed")

// ErrCreateDir happens when we fail to create the target directory.
var ErrCreateDir = errors.New("failed to create directory")

// Clone a given repository into a given destination.
func Clone(repo Repo, destination string) error {
	// nolint: gosec
	var cmd = exec.Command(
		"git", "clone", "--depth", "1", repo.URL,
		filepath.Join(destination, repo.Name),
	)
	if bts, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %v: %v", ErrClone, repo.Name, string(bts))
	}
	return nil
}

// AllOrgRepos finds all repositories of a given organization.
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

const pageSize = 30

func findRepos(ctx context.Context, client *github.Client, org string) (result []*github.Repository, err error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: pageSize},
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

// CreateDir creates the directory if it does not exists.
func CreateDir(directory string) error {
	stat, err := os.Stat(directory)
	directoryDoesNotExists := err != nil

	if directoryDoesNotExists {
		err := os.MkdirAll(directory, 0700)
		if err != nil {
			return fmt.Errorf("%w: %s: %s", ErrCreateDir, directory, err.Error())
		}

		return nil
	}

	if stat.IsDir() {
		return nil
	}

	return fmt.Errorf("%w: %s is a file", ErrCreateDir, directory)
}
