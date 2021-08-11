package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// light palette: https://colorhunt.co/palette/201882
// dark palette:  https://colorhunt.co/palette/273948
var (
	primaryColor = lipgloss.AdaptiveColor{
		Light: "#1a1a2e",
		Dark:  "#f7f3e9",
	}
	secondaryColor = lipgloss.AdaptiveColor{
		Light: "#16213e",
		Dark:  "#a3d2ca",
	}
	errorColor = lipgloss.AdaptiveColor{
		Light: "#e94560",
		Dark:  "#f05945",
	}
	grayColor    = lipgloss.Color("#626262")
	midGrayColor = lipgloss.Color("#4a4a4a")

	secondaryForeground   = lipgloss.NewStyle().Foreground(secondaryColor)
	primaryForegroundBold = lipgloss.NewStyle().Bold(true).Foreground(primaryColor)
	errorFaintForeground  = lipgloss.NewStyle().Foreground(errorColor).Faint(true)
	errorForegroundPadded = lipgloss.NewStyle().Padding(4).Foreground(errorColor)
	grayForeground        = lipgloss.NewStyle().Foreground(grayColor)
	midGrayForeground     = lipgloss.NewStyle().Foreground(midGrayColor)
)

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

func errorView(action string, err error) string {
	return errorForegroundPadded.Render(fmt.Sprintf(action+": %s.\nCheck the log file for more details.", err.Error())) +
		singleOptionHelp("q", "quit")
}
