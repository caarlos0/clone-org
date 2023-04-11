package ui

import (
	"log"
	"strings"

	cloneorg "github.com/caarlos0/clone-org"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func newRepoView(repo cloneorg.Repo, destination string) repoModel {
	s := spinner.New()
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
	return tea.Batch(cloneRepoCmd(m.repo, m.destination), m.spinner.Tick)
}

func (m repoModel) Update(msg tea.Msg) (repoModel, tea.Cmd) {
	switch msg := msg.(type) {
	case repoCloneErrMsg:
		if msg.name == m.repo.Name {
			m.cloning = false
			m.err = msg.error
			log.Println("failed to clone", m.repo.Name, m.err)
		}
	case repoClonedMsg:
		if msg.name == m.repo.Name {
			m.cloning = false
			log.Println("cloned", m.repo.Name)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m repoModel) View() string {
	if m.err != nil {
		return secondaryForeground.Render("[failed] ") +
			m.repo.Name +
			" " +
			errorFaintForeground.Render(strings.TrimSpace(m.err.Error()))
	}
	if m.cloning {
		return secondaryForeground.Render("[cloning] ") +
			m.repo.Name +
			" " +
			primaryForegroundBold.Render(m.spinner.View())
	}
	return secondaryForeground.Render("[cloned] ") +
		m.repo.Name
}

// msgs and cmds

type repoClonedMsg struct {
	name string
}

type repoCloneErrMsg struct {
	error
	name string
}

func cloneRepoCmd(repo cloneorg.Repo, destination string) tea.Cmd {
	return func() tea.Msg {
		if err := cloneorg.Clone(repo, destination); err != nil {
			return repoCloneErrMsg{err, repo.Name}
		}
		return repoClonedMsg{repo.Name}
	}
}
