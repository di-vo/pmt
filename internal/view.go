package internal

import (
	lip "github.com/charmbracelet/lipgloss"
)

var (
	baseStyle = lip.NewStyle().
			BorderStyle(lip.RoundedBorder()).
			BorderForeground(lip.Color("240")).
			Padding(1, 2)

	listStyle = lip.NewStyle().
			Border(lip.RoundedBorder()).
			Padding(1, 2).
			Margin(0, 1).
			Width(20)

	listTitleStyle = lip.NewStyle().
			Foreground(lip.Color("#FFFDF5")).
			Background(lip.Color("#25A065")).
			Padding(0, 1)

	spacer = lip.NewStyle().
		Width(2).
		Render
)

func (m model) View() string {
	s := ""

	switch m.state {
	case "overview":
		s = m.table.View()
	case "addingProject":
		s = m.projectTi.View() + "\n\n" + m.table.View()
	case "detailed":
		listViews := make([]string, len(m.detailLists))
		for i, v := range m.detailLists {
			listViews[i] = listStyle.Render(v.View())
		}

		s = lip.JoinHorizontal(lip.Top, listViews...)
	}

	helpView := m.help.View(m.keys)

	return baseStyle.Render(s) + "\n\n" + helpView
}
