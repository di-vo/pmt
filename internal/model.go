package internal

import (
	"database/sql"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/di-vo/pmt/internal/db"
)

type keyMap struct {
	Up        key.Binding
	Down      key.Binding
	Help      key.Binding
	Quit      key.Binding
	Add       key.Binding
	Delete    key.Binding
	Edit      key.Binding
	Escape    key.Binding
	Enter     key.Binding
	Right     key.Binding
	Left      key.Binding
	Confirm   key.Binding
	Cancel    key.Binding
	Tab       key.Binding
	MoveUp    key.Binding
	MoveDown  key.Binding
	MoveLeft  key.Binding
	MoveRight key.Binding
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
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "enter"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("l", "move right"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("h", "move left"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "confirm"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "cancel"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "tab"),
	),
	MoveUp: key.NewBinding(
		key.WithKeys("ctrl+up", "ctrl+k"),
		key.WithHelp("ctrl+k", "move item up"),
	),
	MoveDown: key.NewBinding(
		key.WithKeys("ctrl+down", "ctrl+j"),
		key.WithHelp("ctrl+j", "move item down"),
	),
	MoveLeft: key.NewBinding(
		key.WithKeys("ctrl+left", "ctrl+h"),
		key.WithHelp("ctrl+h", "move item left"),
	),
	MoveRight: key.NewBinding(
		key.WithKeys("ctrl+right", "ctrl+l"),
		key.WithHelp("ctrl+l", "move item right"),
	),
}

type item struct {
	title    string
	desc     string
	isActive bool
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }

type project struct {
	id          int
	name        string
	itemLists   [][]item
	activeItems *[]item
}

type model struct {
	keys      keyMap
	help      help.Model
	state     string
	table     table.Model
	entries   []project
	textInput textinput.Model
	textArea  textarea.Model
	listIndex int
	itemIndex int
	database  *sql.DB
}

func (m model) getRowsFromEntries() []table.Row {
	rows := make([]table.Row, 0)

	for _, v := range m.entries {
		rows = append(rows, table.Row{strconv.Itoa(v.id), v.name})
	}

	return rows
}

func (m model) getActiveItem() item {
	return (*m.entries[m.table.Cursor()].activeItems)[m.itemIndex]
}

func (m model) toggleActiveItemState(isActive bool) {
	if len(*m.entries[m.table.Cursor()].activeItems) > 0 {
		(*m.entries[m.table.Cursor()].activeItems)[m.itemIndex].isActive = isActive
	}
}

func InitalModel(database *sql.DB) model {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Name", Width: 30},
	}

	entries := []project{
		{id: 1, name: "Arch Setup", itemLists: [][]item{
			{
				{title: "last todo", desc: "almost over"},
				{title: "last test", desc: "almost over"},
				{title: "last todo", desc: "almost joever"},
				{title: "last todo", desc: "this is going to be a very long description"},
			},
			{
				{title: "last todo", desc: "almost over"},
				{title: "last test", desc: "almost over"},
				{title: "last todo", desc: "almost joever"},
			},
			{
				{title: "last todo", desc: "almost over"},
				{title: "last test", desc: "almost over"},
				{title: "last todo", desc: "almost joever"},
			}},
		},
		{id: 2, name: "Project App", itemLists: [][]item{{}, {}, {}}},
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

	ta := textarea.New()
	ta.Placeholder = "Add Description..."

	db.CreateTables(database)

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
		textInput: ti,
		textArea:  ta,
		listIndex: 0,
		itemIndex: 0,
		database:  database,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
