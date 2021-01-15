package ui

import (
	cloneorg "github.com/caarlos0/clone-org"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// NewInitialModel creates a new InitialModel with required fields.
func NewInitialModel(token, org, destination string) tea.Model {
	var s = spinner.NewModel()
	s.Spinner = spinner.Jump

	return initialModel{
		token:       token,
		org:         org,
		destination: destination,
		spinner:     s,
		loading:     true,
	}
}

// InitialModel is the UI when the CLI starts, basically loading the repos.
type initialModel struct {
	err         error
	spinner     spinner.Model
	token       string
	org         string
	destination string
	loading     bool
}

func (m initialModel) Init() tea.Cmd {
	return tea.Batch(getReposCmd(m.token, m.org, m.destination), spinner.Tick)
}

func (m initialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.loading = false
		m.err = msg.error
		return m, nil
	case gotRepoListMsg:
		var list = newCloneModel(msg.repos, m.destination)
		return list, list.Init()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m initialModel) View() string {
	if m.loading {
		return boldPrimaryForeground(m.spinner.View()) + " Gathering repositories..." + singleOptionHelp("q", "quit")
	}
	if m.err != nil {
		return errorView("Error gathering the repositories: ", m.err)
	}
	return ""
}

// msgs and cmds

type gotRepoListMsg struct {
	repos []cloneorg.Repo
}

func getReposCmd(token, org, destination string) tea.Cmd {
	return func() tea.Msg {
		repos, err := cloneorg.AllOrgRepos(token, org)
		if err != nil {
			return errMsg{err}
		}
		if err := cloneorg.CreateDir(destination); err != nil {
			return errMsg{err}
		}
		return gotRepoListMsg{repos}
	}
}
