package internal

import (
	"fmt"
	"strings"

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
		if !lib.Contains(inputStates, m.state) {
			switch {
			case key.Matches(msg, m.keys.Help):
				m.help.ShowAll = !m.help.ShowAll
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keys.Add):
				// Fix: still works in detail state
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
		}

		if key.Matches(msg, m.keys.Cancel) {
			m.state = "overview"
			m.table.Focus()
		}

		if key.Matches(msg, m.keys.Enter) {
			switch m.state {
			case "addingProject":
				// add project to slice and close textinput
				m.state = "overview"
				m.table.Focus()

				if strings.Trim(m.projectTi.Value(), " ") != "" {
					m.entries = append(m.entries, project{id: 3, name: m.projectTi.Value()})
					m.table.SetRows(m.getRowsFromEntries())
					m.table.GotoBottom()
				}
			case "overview":
				// enter detailed view for selected item
				fmt.Println("going into detail")
				m.state = "detailed"
			}
		}
	}

	switch {
	case lib.Contains(inputStates, m.state):
		m.projectTi, cmd = m.projectTi.Update(msg)
	case m.state == "overview":
		m.table, cmd = m.table.Update(msg)
	case m.state == "detailed":
		m.todoList, cmd = m.todoList.Update(msg)
		m.doingList, cmd = m.doingList.Update(msg)
		m.doneList, cmd = m.doneList.Update(msg)
	}
	return m, cmd
}
