package ui

import (
	"strings"
)

func singleOptionHelp(k, v string) string {
	return helpView([]helpOption{
		{k, v, true},
	})
}

var separator = midGrayForeground.Render(" • ")

func helpView(options []helpOption) string {
	var lines []string

	var line []string
	for i, help := range options {
		if help.primary {
			s := grayForeground.Render(help.key) +
				" " +
				secondaryForeground.Faint(true).Render((help.help))
			line = append(line, s)
		} else {
			s := grayForeground.Render(help.key) +
				" " +
				midGrayForeground.Render(help.help)
			line = append(line, s)
		}
		// splits in rows of 3 options max
		if (i+1)%3 == 0 {
			lines = append(lines, strings.Join(line, separator))
			line = []string{}
		}
	}

	// append remainder
	lines = append(lines, strings.Join(line, separator))

	return "\n\n" + strings.Join(lines, "\n")
}

type helpOption struct {
	key, help string
	primary   bool
}
