package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/di-vo/pmt/internal"
)

func main() {
	p := tea.NewProgram(internal.InitalModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There was en error trying to run the program: %v", err)
		os.Exit(1)
	}

}
