package ui

import (
	"fmt"

	cloneorg "github.com/caarlos0/clone-org"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func newCloneModel(repos []cloneorg.Repo, org, destination string, tui bool, w, h int) tea.Model {
	var models []repoModel
	for _, r := range repos {
		models = append(models, newRepoView(r, destination))
	}

	m := cloneModel{
		repos:       models,
		org:         org,
		destination: destination,
		tui:         tui,
	}
	margin := lipgloss.Height(footer) + lipgloss.Height(header(m))
	vp := viewport.New(w, h-margin)
	vp.YPosition = 1
	m.viewport = vp
	return m
}

// ListModel is the UI in which the user can select which forks should be
// deleted if any, and see details on each of them.
type cloneModel struct {
	repos       []repoModel
	viewport    viewport.Model
	org         string
	destination string
	done        bool
	tui         bool
}

func (m cloneModel) Init() tea.Cmd {
	inits := []tea.Cmd{m.viewport.Init()}
	for _, r := range m.repos {
		inits = append(inits, r.Init())
	}
	return tea.Batch(inits...)
}

func (m cloneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Height = msg.Height - 2
		m.viewport.Width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	}
	if m.done {
		if !m.tui {
			return m, tea.Quit
		}
		return m, nil
	}
	m.done = true

	for i := range m.repos {
		m.repos[i], cmd = m.repos[i].Update(msg)
		cmds = append(cmds, cmd)
		if m.repos[i].cloning {
			m.done = false
		}
	}

	var content string
	for _, r := range m.repos {
		if !r.cloning {
			continue
		}
		content += "\n" + r.View()
	}
	for _, r := range m.repos {
		if r.cloning {
			continue
		}
		content += "\n" + r.View()
	}
	m.viewport.SetContent(content)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m cloneModel) View() string {
	return header(m) + "\n" + m.viewport.View() + "\n" + footer
}

var footer = singleOptionHelp("q/esc", "quit")

func header(m cloneModel) string {
	verb := "Cloning"
	if m.done {
		verb = "Cloned"
	}

	return secondaryForeground.Render(fmt.Sprintf(
		"%s %d repositories from %s to %s ...",
		verb,
		len(m.repos),
		m.org,
		m.destination,
	))
}
