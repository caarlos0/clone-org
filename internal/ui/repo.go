package ui

import (
	"strings"

	cloneorg "github.com/caarlos0/clone-org"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func newRepoView(repo cloneorg.Repo, destination string) repoModel {
	var s = spinner.NewModel()
	s.Spinner = spinner.Points
	return repoModel{
		repo:        repo,
		destination: destination,
		spinner:     s,
		cloning:     true,
	}
}

type repoModel struct {
	repo        cloneorg.Repo
	destination string
	spinner     spinner.Model
	cloning     bool
	err         error
}

func (m repoModel) Init() tea.Cmd {
	return tea.Batch(cloneRepoCmd(m.repo, m.destination), spinner.Tick)
}

func (m repoModel) Update(msg tea.Msg) (repoModel, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.cloning = false
		m.err = msg.error
	case repoClonedMsg:
		m.cloning = false
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m repoModel) View() string {
	if m.err != nil {
		return faint("[failed] ") + m.repo.Name + " " + redFaintForeground(strings.TrimSpace(m.err.Error()))
	}
	if m.cloning {
		return faint("[cloning] ") + m.repo.Name + " " + boldSecondaryForeground(m.spinner.View())
	}
	return faint("[cloned] ") + m.repo.Name
}

// msgs and cmds

type repoClonedMsg struct{}

func cloneRepoCmd(repo cloneorg.Repo, destination string) tea.Cmd {
	return func() tea.Msg {
		if err := cloneorg.Clone(repo, destination); err != nil {
			return errMsg{err}
		}
		return repoClonedMsg{}
	}
}
