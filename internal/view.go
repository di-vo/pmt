package internal

import (
	lip "github.com/charmbracelet/lipgloss"
)

var (
	baseStyle = lip.NewStyle().
			BorderStyle(lip.RoundedBorder()).
			BorderForeground(lip.Color("240")).
			Padding(1, 2)

	listStyle = lip.NewStyle().
			Border(lip.RoundedBorder()).
			Padding(1, 2, 0).
			Margin(0, 1)

	activeListStyle = lip.NewStyle().
			Border(lip.RoundedBorder()).
			BorderForeground(lip.Color("#458588")).
			Padding(1, 2, 0).
			Margin(0, 1)

	listTitleStyle = lip.NewStyle().
			Foreground(lip.Color("#FFFDF5")).
			Background(lip.Color("#25A065")).
			Padding(0, 2)

	activeListTitleStyle = lip.NewStyle().
				Foreground(lip.Color("#FFFDF5")).
				Background(lip.Color("#458588")).
				Padding(0, 2)

	itemStyle = lip.NewStyle().
			BorderStyle(lip.RoundedBorder()).
			MarginBottom(1).
			Padding(0, 1, 1)

	activeItemStyle = lip.NewStyle().
			BorderStyle(lip.RoundedBorder()).
			BorderForeground(lip.Color("#F186C7")).
			MarginBottom(1).
			Padding(0, 1, 1)

	itemTitleStyle = lip.NewStyle().
			Bold(true).
			Underline(true)

	deleteHintStyle = lip.NewStyle().
			Foreground(lip.Color("#FF0000")).
			Bold(true)
)

func renderElement(text string, width int) string {
	s := ""
	buf := ""
	isSpaceAtStart := false

	for i, c := range text {
		if isSpaceAtStart {
			isSpaceAtStart = false
			continue
		}

		buf += string(c)

		// handle wrapping according to width
		if len(buf) == width {
			if string(c) == " " {
				// is the last rune a space
				s += buf + "\n"
				buf = ""
			} else if i < len(text)-1 && string(text[i+1]) == " " {
				// is the next rune
				s += buf + "\n"
				isSpaceAtStart = true
				buf = ""
			} else if i > 0 && string(c) != " " && string(buf[len(buf)-2]) == " " {
				// is the last rune the first rune of a word
				s += buf[:len(buf)-1] + "\n"
				buf = string(buf[len(buf)-1])
			} else {
				// does the buffer end without the word finished
				s += buf[:len(buf)-1] + "-\n"
				buf = string(buf[len(buf)-1])
			}
		}

		if i == len(text)-1 {
			s += buf
		}
	}

	return s
}

func renderList(items []item, title string, width int, isActive bool) string {
	s := ""

	if isActive {
		s = activeListTitleStyle.Render(title) + "\n\n"
	} else {
		s = listTitleStyle.Render(title) + "\n\n"
	}

	for _, v := range items {
		itemString := itemTitleStyle.Render(renderElement(v.title, width)) + "\n"
		itemString += renderElement(v.desc, width)

		if v.isActive {
			s += activeItemStyle.Render(itemString) + "\n"
		} else {
			s += itemStyle.Render(itemString) + "\n"
		}
	}

	if isActive {
		return activeListStyle.Render(s)
	} else {
		return listStyle.Render(s)
	}
}

func renderDetailed(m model) string {
	sp := &m.entries[m.table.Cursor()]

	listWidth := 20
	lists := []string{
		renderList(sp.itemLists[0], "ToDo", listWidth, &sp.itemLists[0] == sp.activeItems),
		renderList(sp.itemLists[1], "Doing", listWidth, &sp.itemLists[1] == sp.activeItems),
		renderList(sp.itemLists[2], "Done", listWidth, &sp.itemLists[2] == sp.activeItems)}

	return lip.JoinHorizontal(lip.Top, lists...)
}

// The main rendering function
func (m model) View() string {
	s := ""

	switch m.state {
	case "overview":
		s = m.table.View()
	case "addingProject", "editingProject":
		s = m.textInput.View() + "\n\n" + m.table.View()
	case "detailed":
		s = renderDetailed(m)
	case "addingItem", "editingItem":
		s = "Tab to switch | Enter when title is focused to confirm\n"
		s += m.textInput.View() + "\n\n" + m.textArea.View() + "\n\n" + renderDetailed(m)
	case "removingProject":
		s = lip.JoinVertical(lip.Center,
			deleteHintStyle.Render("Permanently delete this Project: ")+m.entries[m.table.Cursor()].name,
			"(y)es / (c)ancel",
			" ",
			m.table.View())
	case "removingItem":
		s = lip.JoinVertical(lip.Center,
			deleteHintStyle.Render("Delete ")+(*m.entries[m.table.Cursor()].activeItems)[m.itemIndex].title,
			"(y)es / (c)ancel",
			" ",
			renderDetailed(m))
	}

	helpView := m.help.View(m.keys)

	return baseStyle.Render(s) + "\n\n" + helpView
}
