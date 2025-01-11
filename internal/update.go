package internal

import (
	"strings"

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
		switch m.state {
		case "overview":
			if key.Matches(msg, m.keys.Add) {
				m.state = "addingProject"
				m.projectTi.Focus()
				m.projectTi.SetValue("")
				return m, nil
			} else if key.Matches(msg, m.keys.Delete) {
				m.state = "removingProject"

			} else if key.Matches(msg, m.keys.Enter) {
				// enter detailed view for selected item
				m.state = "detailed"
			} else if key.Matches(msg, m.keys.Help) {
				m.help.ShowAll = !m.help.ShowAll
			} else if key.Matches(msg, m.keys.Quit) {
				return m, tea.Quit
			}
		case "detailed":
			if key.Matches(msg, m.keys.Up) {
				(*m.entries[m.table.Cursor()].activeItems)[m.itemIndex].isActive = false

				m.itemIndex--

				if m.itemIndex < 0 {
					m.itemIndex = len(*m.entries[m.table.Cursor()].activeItems) - 1
				}
			} else if key.Matches(msg, m.keys.Down) {
				(*m.entries[m.table.Cursor()].activeItems)[m.itemIndex].isActive = false

				m.itemIndex++

				if m.itemIndex == len(*m.entries[m.table.Cursor()].activeItems) {
					m.itemIndex = 0
				}
			} else if key.Matches(msg, m.keys.Add) {

			} else if key.Matches(msg, m.keys.Delete) {

			} else if key.Matches(msg, m.keys.Escape) {
				m.state = "overview"
				m.table.Focus()
			} else if key.Matches(msg, m.keys.Help) {
				m.help.ShowAll = !m.help.ShowAll
			} else if key.Matches(msg, m.keys.Quit) {
				return m, tea.Quit
			} else if key.Matches(msg, m.keys.Tab) {
				(*m.entries[m.table.Cursor()].activeItems)[m.itemIndex].isActive = false

				m.listIndex++

				if m.listIndex == len(m.entries[0].itemLists) {
					m.listIndex = 0
				}

				m.itemIndex = 0

			}
		case "addingProject":
			if key.Matches(msg, m.keys.Enter) {
				// add project to slice and close textinput
				m.state = "overview"
				m.table.Focus()

				if strings.Trim(m.projectTi.Value(), " ") != "" {
					m.entries = append(m.entries, project{id: 3, name: m.projectTi.Value()})
					m.table.SetRows(m.getRowsFromEntries())
					m.table.GotoBottom()
				}
			} else if key.Matches(msg, m.keys.Escape) {
				m.state = "overview"
				m.table.Focus()
			}
		case "addingItem":
			if key.Matches(msg, m.keys.Enter) {

			} else if key.Matches(msg, m.keys.Escape) {

			}
		case "removingProject":
			if key.Matches(msg, m.keys.Confirm) {
				c := m.table.Cursor()

				if len(m.entries) > 0 {
					newEntries := make([]project, 0)
					newEntries = append(newEntries, m.entries[:c]...)
					newEntries = append(newEntries, m.entries[c+1:]...)
					m.entries = newEntries

					m.table.SetRows(m.getRowsFromEntries())
					m.table.SetCursor(c)
				}

				m.state = "overview"
			} else if key.Matches(msg, m.keys.Cancel) {
				m.state = "overview"
			}
		case "removingItem":
			// confirm, cancel
		}
	}

	// after input handling, update elements according to the current state
	switch {
	case m.state == "overview":
		m.table, cmd = m.table.Update(msg)
	case m.state == "detailed":
		m.entries[m.table.Cursor()].activeItems = &m.entries[m.table.Cursor()].itemLists[m.listIndex]
		(*m.entries[m.table.Cursor()].activeItems)[m.itemIndex].isActive = true
	case m.state == "addingProject":
		m.projectTi, cmd = m.projectTi.Update(msg)
	}
	return m, cmd
}
