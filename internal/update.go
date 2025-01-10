package internal

import (
	"github.com/charmbracelet/bubbles/key"
	_ "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Add):
			m.entries = append(m.entries, project{id: 3, name: "Added Project"})
			m.table.SetRows(m.getRowsFromEntries())

			m.table.SetCursor(len(m.entries))
		case key.Matches(msg, m.keys.Delete):
			c := m.table.Cursor()

			if len(m.entries) > 0 {
				newEntries := make([]project, 0)
				newEntries = append(newEntries, m.entries[:c]...)
				newEntries = append(newEntries, m.entries[c+1:]...)
				m.entries = newEntries

				m.table.SetRows(m.getRowsFromEntries())
				m.table.SetCursor(c)
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}
