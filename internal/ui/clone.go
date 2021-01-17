package ui

import (
	"fmt"

	cloneorg "github.com/caarlos0/clone-org"
	tea "github.com/charmbracelet/bubbletea"
)

func newCloneModel(repos []cloneorg.Repo, destination string) tea.Model {
	var models []repoModel
	for _, r := range repos {
		models = append(models, newRepoView(r, destination))
	}
	return cloneModel{
		repos:       models,
		destination: destination,
	}
}

// ListModel is the UI in which the user can select which forks should be
// deleted if any, and see details on each of them.
type cloneModel struct {
	repos       []repoModel
	destination string
	done        bool
}

func (m cloneModel) Init() tea.Cmd {
	var inits []tea.Cmd
	for _, r := range m.repos {
		inits = append(inits, r.Init())
	}
	return tea.Batch(inits...)
}

func (m cloneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	}
	if m.done {
		return m, nil
	}
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.done = true
	for i := range m.repos {
		m.repos[i], cmd = m.repos[i].Update(msg)
		cmds = append(cmds, cmd)
		if m.repos[i].cloning {
			m.done = false
		}
	}
	return m, tea.Batch(cmds...)
}

func (m cloneModel) View() string {
	var verb = "Cloning"
	if m.done {
		verb = "Cloned"
	}
	var s = boldSecondaryForeground(fmt.Sprintf("%s %d repositories to %s ...\n\n", verb, len(m.repos), m.destination))
	for _, r := range m.repos {
		s += r.View() + "\n"
	}

	return s + singleOptionHelp("q/esc", "quit")
}
