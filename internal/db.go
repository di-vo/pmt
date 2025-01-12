package internal

import (
	"database/sql"
	"fmt"
	"os"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

func resetDb(database *sql.DB) {
	query := "DROP TABLE item"

	_, err := database.Exec(query)
	if err != nil {
		fmt.Printf("Error running query: %v", err)
		os.Exit(1)
	}

	query = "DROP TABLE project"

	_, err = database.Exec(query)
	if err != nil {
		fmt.Printf("Error running query: %v", err)
		os.Exit(1)
	}

	createTables(database)
}

func createTables(database *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS project (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);`

	_, err := database.Exec(query)
	if err != nil {
		fmt.Printf("Error running query: %v", err)
		os.Exit(1)
	}

	query = `
	CREATE TABLE IF NOT EXISTS item (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		desc TEXT,
		lid INTEGER,
		orderNr INTEGER,
		pid INTEGER REFERENCES project(id)
	);`

	_, err = database.Exec(query)
	if err != nil {
		fmt.Printf("Error running query: %v", err)
		os.Exit(1)
	}
}

func getProjects(database *sql.DB) []project {
	var projects []project

	rows, err := database.Query("SELECT * FROM project")
	if err != nil {
		fmt.Printf("Error starting transaction: %v", err)
		os.Exit(1)
	}
	defer rows.Close()

	for rows.Next() {
		var p project
		// hardcoded atm
		p.itemLists = make([][]item, 3)

		err = rows.Scan(&p.id, &p.name)
		if err != nil {
			fmt.Printf("1Error querying db: %v", err)
			os.Exit(1)
		}

		stmt, itemErr := database.Prepare("SELECT id, title, desc, lid, orderNr FROM item WHERE pid = ?")
		if itemErr != nil {
			fmt.Printf("2Error querying db: %v", err)
			os.Exit(1)
		}
		defer stmt.Close()

		var items *sql.Rows

		items, itemErr = stmt.Query(p.id)
		if itemErr != nil {
			fmt.Printf("3Error querying db: %v", err)
			os.Exit(1)
		}
		defer items.Close()

		listOrders := make([]map[int]item, 3)
		for i := range listOrders {
			listOrders[i] = make(map[int]item, 0)
		}

		for items.Next() {
			var it item
			var lid int
			var orderNr int

			itemErr = items.Scan(&it.id, &it.title, &it.desc, &lid, &orderNr)
			if itemErr != nil {
				fmt.Printf("4Error querying db: %v", err)
				os.Exit(1)
			}

			listOrders[lid][orderNr] = it
		}

		for i, orders := range listOrders {
			keys := make([]int, 0)

			for k := range orders {
				keys = append(keys, k)
			}

			sort.Ints(keys)

			for _, k := range keys {
				p.itemLists[i] = append(p.itemLists[i], orders[k])
			}
		}

		projects = append(projects, p)
	}

	err = rows.Err()
	if err != nil {
		fmt.Printf("5Error querying db: %v", err)
		os.Exit(1)
	}

	return projects
}
