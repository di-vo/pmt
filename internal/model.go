// TODO:
// - read up on list

package internal

import (
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Help   key.Binding
	Quit   key.Binding
	Add    key.Binding
	Delete key.Binding
	Escape key.Binding
	Enter  key.Binding
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
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "enter"),
	),
}

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type project struct {
	id         int
	name       string
	todoItems  []item
	doingItems []item
	doneItems  []item
}

type model struct {
	keys      keyMap
	help      help.Model
	state     string
	table     table.Model
	entries   []project
	projectTi textinput.Model
	todoList  list.Model
	doingList list.Model
	doneList  list.Model
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

	ti := textinput.New()
	ti.Placeholder = "ESC to cancel..."
	ti.CharLimit = 50
	ti.Width = 20

	todoList := list.New([]list.Item{item{title: "last todo", desc: "almost over"}, item{title: "last test", desc: "almost over"}, item{title: "last todo", desc: "almost joever"}}, list.NewDefaultDelegate(), 30, 10)
	todoList.Title = "ToDo"

	doingList := list.New([]list.Item{item{title: "despairge", desc: "yep"}}, list.NewDefaultDelegate(), 30, 10)
	doingList.Title = "Doing"

	doneList := list.New([]list.Item{item{title: "first task", desc: "lets go"}}, list.NewDefaultDelegate(), 30, 10)
	doneList.Title = "Done"

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
		entries:   entries,
		projectTi: ti,
		todoList:  todoList,
		doingList: doingList,
		doneList:  doneList,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
