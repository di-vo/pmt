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

	listTitleStyle = lip.NewStyle().
			Foreground(lip.Color("#FFFDF5")).
			Background(lip.Color("#25A065")).
			Padding(0, 2)

	itemStyle = lip.NewStyle().
			BorderStyle(lip.RoundedBorder()).
			MarginBottom(1).
			Padding(0, 1, 1)

	itemTitleStyle = lip.NewStyle().
			Bold(true).
			Underline(true)
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

func renderList(items []item, title string, width int) string {
	s := listTitleStyle.Render(title) + "\n\n"

	for _, v := range items {
		itemString := itemTitleStyle.Render(renderElement(v.title, width)) + "\n"
		itemString += renderElement(v.desc, width)

		s += itemStyle.Render(itemString) + "\n"
	}

	return listStyle.Render(s)
}

// The main rendering function
func (m model) View() string {
	s := ""

	switch m.state {
	case "overview":
		s = m.table.View()
	case "addingProject":
		s = m.projectTi.View() + "\n\n" + m.table.View()
	case "detailed":
		sp := m.entries[m.table.Cursor()]

		listWidth := 20
		lists := []string{
			renderList(sp.todoItems, "ToDo", listWidth),
			renderList(sp.doingItems, "Doing", listWidth),
			renderList(sp.doneItems, "Done", listWidth)}

		s = lip.JoinHorizontal(lip.Top, lists...)
	}

	helpView := m.help.View(m.keys)

	return baseStyle.Render(s) + "\n\n" + helpView
}
