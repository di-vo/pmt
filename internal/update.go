package internal

import (
	"github.com/charmbracelet/bubbles/key"
	_ "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/di-vo/pmt/lib"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	inputStates := []string{
		"addingProject",
	}
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case !lib.Contains(inputStates, m.state):
			switch {
			case key.Matches(msg, m.keys.Help):
				m.help.ShowAll = !m.help.ShowAll
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keys.Add):
				m.state = "addingProject"
				m.projectTi.Focus()
				m.projectTi.SetValue("")
				return m, nil
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
		case key.Matches(msg, m.keys.Cancel):
			m.state = "overview"
			m.table.Focus()
		case key.Matches(msg, m.keys.Enter):
			switch m.state {
			case "addingProject":
				m.state = "overview"
				m.table.Focus()

				m.entries = append(m.entries, project{id: 3, name: m.projectTi.Value()})
				m.table.SetRows(m.getRowsFromEntries())
				//m.table.SetCursor(len(m.entries))
				m.table.GotoBottom()
			}
		}
	}

	if lib.Contains(inputStates, m.state) {
		m.projectTi, cmd = m.projectTi.Update(msg)
	} else {
		m.table, cmd = m.table.Update(msg)
	}
	return m, cmd
}
