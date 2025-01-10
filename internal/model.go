// TODO:
// - set up a table element
// - create SQLite DB and connect to it

package internal

import tea "github.com/charmbracelet/bubbletea"

type model struct {
}

func InitalModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m model) View() string {
	s := ""

	return s
}
