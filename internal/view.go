package internal

import (
	lip "github.com/charmbracelet/lipgloss"
)

var baseStyle = lip.NewStyle().
	BorderStyle(lip.RoundedBorder()).
	BorderForeground(lip.Color("240"))

func (m model) View() string {
	s := ""

	switch m.state {
	case "overview":
		s = m.table.View()
	case "addingProject":
		s = m.projectTi.View() + "\n\n" + m.table.View()
	}

	helpView := m.help.View(m.keys)

	return baseStyle.Render(s) + "\n\n" + helpView
}
