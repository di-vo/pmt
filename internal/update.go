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
				m.addTi.Focus()
				m.addTi.SetValue("")
				return m, nil
			} else if key.Matches(msg, m.keys.Edit) {
				m.state = "editingProject"
				m.addTi.Focus()
				m.addTi.SetValue(m.entries[m.table.Cursor()].name)
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
				m.state = "addingItem"
				m.addTi.Focus()
				m.addTi.SetValue("")
				m.addTa.SetValue("")
				return m, nil
			} else if key.Matches(msg, m.keys.Edit) {
				m.state = "editingItem"
			} else if key.Matches(msg, m.keys.Delete) {
				m.state = "removingItem"
			} else if key.Matches(msg, m.keys.Escape) {
				m.state = "overview"
				m.table.Focus()
			} else if key.Matches(msg, m.keys.Help) {
				m.help.ShowAll = !m.help.ShowAll
			} else if key.Matches(msg, m.keys.Quit) {
				return m, tea.Quit
			} else if key.Matches(msg, m.keys.Right) {
				(*m.entries[m.table.Cursor()].activeItems)[m.itemIndex].isActive = false

				m.listIndex++

				if m.listIndex == len(m.entries[0].itemLists) {
					m.listIndex = 0
				}

				m.itemIndex = 0
			} else if key.Matches(msg, m.keys.Left) {
				(*m.entries[m.table.Cursor()].activeItems)[m.itemIndex].isActive = false

				m.listIndex--

				if m.listIndex < 0 {
					m.listIndex = len(m.entries[0].itemLists) - 1
				}

				m.itemIndex = 0
			} else if key.Matches(msg, m.keys.MoveUp) {
				if m.itemIndex > 0 {
					tmp := (*m.entries[m.table.Cursor()].activeItems)[m.itemIndex]
					(*m.entries[m.table.Cursor()].activeItems)[m.itemIndex] = (*m.entries[m.table.Cursor()].activeItems)[m.itemIndex-1]
					(*m.entries[m.table.Cursor()].activeItems)[m.itemIndex-1] = tmp

					m.itemIndex--
				}
			} else if key.Matches(msg, m.keys.MoveDown) {
				if m.itemIndex < len(*m.entries[m.table.Cursor()].activeItems)-1 {
					tmp := (*m.entries[m.table.Cursor()].activeItems)[m.itemIndex]
					(*m.entries[m.table.Cursor()].activeItems)[m.itemIndex] = (*m.entries[m.table.Cursor()].activeItems)[m.itemIndex+1]
					(*m.entries[m.table.Cursor()].activeItems)[m.itemIndex+1] = tmp

					m.itemIndex++
				}
			} else if key.Matches(msg, m.keys.MoveLeft) {
				if m.listIndex > 0 {
					c := m.table.Cursor()

					activeItem := m.getActiveItem()

					newItems := make([]item, 0)
					newItems = append(newItems, (*m.entries[c].activeItems)[:m.itemIndex]...)
					newItems = append(newItems, (*m.entries[c].activeItems)[m.itemIndex+1:]...)
					*m.entries[c].activeItems = newItems

					m.listIndex--
					m.entries[c].itemLists[m.listIndex] = append(m.entries[c].itemLists[m.listIndex], activeItem)

					m.itemIndex = len(m.entries[c].itemLists[m.listIndex]) - 1
				}
			} else if key.Matches(msg, m.keys.MoveRight) {
				if m.listIndex < len(m.entries[m.table.Cursor()].itemLists)-1 {
					c := m.table.Cursor()

					activeItem := m.getActiveItem()

					newItems := make([]item, 0)
					newItems = append(newItems, (*m.entries[c].activeItems)[:m.itemIndex]...)
					newItems = append(newItems, (*m.entries[c].activeItems)[m.itemIndex+1:]...)
					*m.entries[c].activeItems = newItems

					m.listIndex++
					m.entries[c].itemLists[m.listIndex] = append(m.entries[c].itemLists[m.listIndex], activeItem)

					m.itemIndex = len(m.entries[c].itemLists[m.listIndex]) - 1

				}
			}
		case "addingProject":
			if key.Matches(msg, m.keys.Enter) {
				// add project to slice and close textinput
				m.state = "overview"
				m.table.Focus()

				if strings.Trim(m.addTi.Value(), " ") != "" {
					m.entries = append(m.entries, project{id: 3, name: m.addTi.Value()})
					m.table.SetRows(m.getRowsFromEntries())
					m.table.GotoBottom()
				}
			} else if key.Matches(msg, m.keys.Escape) {
				m.state = "overview"
				m.table.Focus()
			}
		case "addingItem":
			if key.Matches(msg, m.keys.Enter) && m.addTi.Focused() {
				newItem := item{title: m.addTi.Value(), desc: m.addTa.Value()}
				*m.entries[m.table.Cursor()].activeItems = append(*m.entries[m.table.Cursor()].activeItems, newItem)

				m.state = "detailed"
			} else if key.Matches(msg, m.keys.Escape) {
				m.state = "detailed"
			} else if key.Matches(msg, m.keys.Tab) {
				if m.addTi.Focused() {
					m.addTi.Blur()
					m.addTa.Focus()
				} else {
					m.addTa.Blur()
					m.addTi.Focus()
				}
			}
		case "editingProject":
			if key.Matches(msg, m.keys.Enter) {
				if strings.Trim(m.addTi.Value(), " ") != "" {
					m.entries[m.table.Cursor()].name = m.addTi.Value()
					m.table.SetRows(m.getRowsFromEntries())
				}

				m.state = "overview"
				m.table.Focus()
			} else if key.Matches(msg, m.keys.Escape) {
				m.state = "overview"
				m.table.Focus()
			}

		case "editingItem":

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
			if key.Matches(msg, m.keys.Confirm) {
				c := m.table.Cursor()

				if len(*m.entries[c].activeItems) > 0 {
					newItems := make([]item, 0)
					newItems = append(newItems, (*m.entries[c].activeItems)[:m.itemIndex]...)
					newItems = append(newItems, (*m.entries[c].activeItems)[m.itemIndex+1:]...)
					*m.entries[c].activeItems = newItems

					m.itemIndex--
				}

				m.state = "detailed"
			} else if key.Matches(msg, m.keys.Cancel) {
				m.state = "detailed"
			}
		}
	}

	// after input handling, update elements according to the current state
	switch {
	case m.state == "overview":
		m.table, cmd = m.table.Update(msg)
	case m.state == "detailed":
		m.entries[m.table.Cursor()].activeItems = &m.entries[m.table.Cursor()].itemLists[m.listIndex]
		(*m.entries[m.table.Cursor()].activeItems)[m.itemIndex].isActive = true
	case m.state == "addingProject" || m.state == "editingProject":
		m.addTi, cmd = m.addTi.Update(msg)
	case m.state == "addingItem":
		m.addTi, cmd = m.addTi.Update(msg)
		m.addTa, cmd = m.addTa.Update(msg)
	}
	return m, cmd
}
