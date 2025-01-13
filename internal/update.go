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
				m.textInput.Focus()
				m.textInput.SetValue("")
				return m, nil
			} else if key.Matches(msg, m.keys.Edit) {
				m.state = "editingProject"
				m.textInput.Focus()
				m.textInput.SetValue(m.entries[m.table.Cursor()].name)
				return m, nil
			} else if key.Matches(msg, m.keys.Delete) {
				m.state = "removingProject"
			} else if key.Matches(msg, m.keys.Enter) {
				// enter detailed view for selected item
				if len(m.entries) > 0 {
					m.state = "detailed"
				}
			} else if key.Matches(msg, m.keys.Help) {
				m.help.ShowAll = !m.help.ShowAll
			} else if key.Matches(msg, m.keys.Quit) {
				return m, tea.Quit
			}
		case "detailed":
			if key.Matches(msg, m.keys.Up) {
				m.toggleActiveItemState(false)

				m.itemIndex--

				if m.itemIndex < 0 {
					m.itemIndex = len(*m.entries[m.table.Cursor()].activeItems) - 1
				}
			} else if key.Matches(msg, m.keys.Down) {
				m.toggleActiveItemState(false)

				m.itemIndex++

				if m.itemIndex == len(*m.entries[m.table.Cursor()].activeItems) {
					m.itemIndex = 0
				}
			} else if key.Matches(msg, m.keys.Add) {
				m.state = "addingItem"
				m.textInput.Focus()
				m.textInput.SetValue("")
				m.textArea.SetValue("")
				return m, nil
			} else if key.Matches(msg, m.keys.Edit) {
				m.state = "editingItem"
				m.textInput.Focus()
				m.textInput.SetValue(m.getActiveItem().title)
				m.textArea.SetValue(m.getActiveItem().desc)
				return m, nil
			} else if key.Matches(msg, m.keys.Delete) {
				m.state = "removingItem"
			} else if key.Matches(msg, m.keys.Escape) {
				// update items in db
				activeProj := m.entries[m.table.Cursor()]

				for i, v := range activeProj.itemLists {
					for j, w := range v {
						updateItem(m.database, w, i, j)
					}
				}

				m.state = "overview"
				m.table.Focus()
			} else if key.Matches(msg, m.keys.Help) {
				m.help.ShowAll = !m.help.ShowAll
			} else if key.Matches(msg, m.keys.Quit) {
				// update items in db
				activeProj := m.entries[m.table.Cursor()]

				for i, v := range activeProj.itemLists {
					for j, w := range v {
						updateItem(m.database, w, i, j)
					}
				}

				return m, tea.Quit
			} else if key.Matches(msg, m.keys.Right) {
				m.toggleActiveItemState(false)

				m.listIndex++

				if m.listIndex == len(m.entries[0].itemLists) {
					m.listIndex = 0
				}

				m.itemIndex = 0
			} else if key.Matches(msg, m.keys.Left) {
				m.toggleActiveItemState(false)

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

				if strings.Trim(m.textInput.Value(), " ") != "" {
					var p project
					p.name = m.textInput.Value()
					p.itemLists = make([][]item, 3)

					// insert into db
					p.id = insertProject(m.database, p)

					m.entries = append(m.entries, p)
					m.table.SetRows(m.getRowsFromEntries())
					m.table.GotoBottom()
				}
			} else if key.Matches(msg, m.keys.Escape) {
				m.state = "overview"
				m.table.Focus()
			}
		case "addingItem":
			if key.Matches(msg, m.keys.Enter) && m.textInput.Focused() {
				newItem := item{title: m.textInput.Value(), desc: m.textArea.Value()}

				// insert into db
				newItem.id = insertItem(m.database, newItem, m.listIndex, m.itemIndex, m.entries[m.table.Cursor()].id)

				*m.entries[m.table.Cursor()].activeItems = append(*m.entries[m.table.Cursor()].activeItems, newItem)

				m.state = "detailed"
			} else if key.Matches(msg, m.keys.Escape) {
				m.state = "detailed"
			} else if key.Matches(msg, m.keys.Tab) {
				if m.textInput.Focused() {
					m.textInput.Blur()
					m.textArea.Focus()
				} else {
					m.textArea.Blur()
					m.textInput.Focus()
				}
			}
		case "editingProject":
			if key.Matches(msg, m.keys.Enter) {
				if strings.Trim(m.textInput.Value(), " ") != "" {
					m.entries[m.table.Cursor()].name = m.textInput.Value()
					m.table.SetRows(m.getRowsFromEntries())

					// update in db
					updateProject(m.database, m.entries[m.table.Cursor()])
				}

				m.state = "overview"
				m.table.Focus()
			} else if key.Matches(msg, m.keys.Escape) {
				m.state = "overview"
				m.table.Focus()
			}
		case "editingItem":
			if key.Matches(msg, m.keys.Enter) {
				if strings.Trim(m.textInput.Value(), " ") != "" {
					updatedItem := m.getActiveItem()
					updatedItem.title = m.textInput.Value()
					updatedItem.desc = m.textArea.Value()

					(*m.entries[m.table.Cursor()].activeItems)[m.itemIndex] = updatedItem
				}

				m.state = "detailed"
			} else if key.Matches(msg, m.keys.Escape) {
				m.state = "detailed"
			} else if key.Matches(msg, m.keys.Tab) {
				if m.textInput.Focused() {
					m.textInput.Blur()
					m.textArea.Focus()
				} else {
					m.textArea.Blur()
					m.textInput.Focus()
				}
			}
		case "removingProject":
			if key.Matches(msg, m.keys.Confirm) {
				c := m.table.Cursor()

				if len(m.entries) > 0 {
					// delete from db
					deleteProject(m.database, m.entries[c])

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
					// delete from db
					deleteItem(m.database, (*m.entries[c].activeItems)[m.itemIndex])

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

		m.toggleActiveItemState(true)
	case m.state == "addingProject" || m.state == "editingProject":
		m.textInput, cmd = m.textInput.Update(msg)
	case m.state == "addingItem" || m.state == "editingItem":
		m.textInput, cmd = m.textInput.Update(msg)
		m.textArea, cmd = m.textArea.Update(msg)
	}
	return m, cmd
}
