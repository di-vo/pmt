package internal

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	_ "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/di-vo/pmt/lib"
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
				c := m.table.Cursor()

				if len(m.entries) > 0 {
					newEntries := make([]project, 0)
					newEntries = append(newEntries, m.entries[:c]...)
					newEntries = append(newEntries, m.entries[c+1:]...)
					m.entries = newEntries

					m.table.SetRows(m.getRowsFromEntries())
					m.table.SetCursor(c)
				}
			} else if key.Matches(msg, m.keys.Enter) {
				// enter detailed view for selected item
				m.state = "detailed"
			} else if key.Matches(msg, m.keys.Help) {
				m.help.ShowAll = !m.help.ShowAll
			} else if key.Matches(msg, m.keys.Quit) {
				return m, tea.Quit
			}
		case "detailed":
			if key.Matches(msg, m.keys.Add) {

			} else if key.Matches(msg, m.keys.Delete) {

			} else if key.Matches(msg, m.keys.Escape) {
				m.state = "overview"
				m.table.Focus()
			} else if key.Matches(msg, m.keys.Help) {
				m.help.ShowAll = !m.help.ShowAll
				lib.WriteToLog("pressed esc")
			} else if key.Matches(msg, m.keys.Quit) {
				return m, tea.Quit
			} else if key.Matches(msg, m.keys.Tab) {
				lib.WriteToLog("pressed tab")
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
			// confirm, cancel
		case "removingItem":
			// confirm, cancel
		}
	}

	// after input handling, update elements according to the current state
	switch {
	case m.state == "overview":
		m.table, cmd = m.table.Update(msg)
	case m.state == "detailed":
		lib.WriteToLog("focusindex in update: " + strconv.Itoa(m.focusIndex))
		//m.detailLists[m.focusIndex], cmd = m.detailLists[m.focusIndex].Update(msg)
	case m.state == "addingProject":
		m.projectTi, cmd = m.projectTi.Update(msg)
	}
	return m, cmd
}
