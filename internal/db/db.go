package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func CreateTables(database *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS project (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);`

	_, err := database.Exec(query)
	if err != nil {
		// TODO: think of better error handling
		os.Exit(1)
	}

	query = `
	CREATE TABLE IF NOT EXISTS item (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		desc TEXT,
		pid INTEGER REFERENCES project(id)
	);`

	_, err = database.Exec(query)
	if err != nil {
		// TODO: think of better error handling
		os.Exit(1)
	}
}
