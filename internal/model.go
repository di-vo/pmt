// TODO:
// - create SQLite DB and connect to it

package internal

import (
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Help   key.Binding
	Quit   key.Binding
	Add    key.Binding
	Delete key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Add, k.Delete},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	Delete: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "delete"),
	),
}

type project struct {
	id   int
	name string
}

type model struct {
	keys    keyMap
	help    help.Model
	state   string
	table   table.Model
	entries []project
}

func (m model) getRowsFromEntries() []table.Row {
	rows := make([]table.Row, 0)

	for _, v := range m.entries {
		rows = append(rows, table.Row{strconv.Itoa(v.id), v.name})
	}

	return rows
}

func InitalModel() model {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Name", Width: 30},
	}

	entries := []project{
		{id: 1, name: "Arch Setup"},
		{id: 2, name: "Project App"},
		{id: 3, name: "Project App"},
		{id: 4, name: "Project App"},
		{id: 5, name: "Project App"},
		{id: 6, name: "Project App"},
		{id: 7, name: "Project App"},
		{id: 8, name: "Project App"},
		{id: 9, name: "Project App"},
		{id: 10, name: "Project App"},
	}

	var rows []table.Row

	for _, v := range entries {
		rows = append(rows, table.Row{strconv.Itoa(v.id), v.name})
	}

	return model{
		keys:  keys,
		help:  help.New(),
		state: "overview",
		table: table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(10),
		),
		entries: entries,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
