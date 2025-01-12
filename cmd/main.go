package main

import (
	"database/sql"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/di-vo/pmt/internal"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, err := sql.Open("sqlite3", "./projectPlanner.db")
	if err != nil {
		fmt.Printf("Error opening db: %v", err)
		os.Exit(1)
	}
	defer database.Close()

	p := tea.NewProgram(internal.InitalModel(database), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There was en error trying to run the program: %v", err)
		os.Exit(1)
	}

}
