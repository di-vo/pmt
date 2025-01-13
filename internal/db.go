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

func insertProject(database *sql.DB, p project) int {
	ctx, err := database.Begin()
	if err != nil {
		fmt.Printf("Error opening db connection: %v", err)
		os.Exit(1)
	}

	stmt, err := ctx.Prepare("INSERT INTO project (name) VALUES (?)")
	if err != nil {
		fmt.Printf("Error preparing statement: %v", err)
		os.Exit(1)
	}
	defer stmt.Close()

	res, err := stmt.Exec(p.name)
	if err != nil {
		fmt.Printf("Error executing statement: %v", err)
		os.Exit(1)
	}

	err = ctx.Commit()
	if err != nil {
		fmt.Printf("Error commiting change: %v", err)
		os.Exit(1)
	}

	newId, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("Error retrieving id: %v", err)
		os.Exit(1)
	}

	return int(newId)
}

func insertItem(database *sql.DB, it item, lid int, orderNr int, pid int) int {
	ctx, err := database.Begin()
	if err != nil {
		fmt.Printf("Error opening db connection: %v", err)
		os.Exit(1)
	}

	stmt, err := ctx.Prepare("INSERT INTO item (title, desc, lid, orderNr, pid) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Printf("Error preparing statement: %v", err)
		os.Exit(1)
	}
	defer stmt.Close()

	res, err := stmt.Exec(it.title, it.desc, lid, orderNr, pid)
	if err != nil {
		fmt.Printf("Error executing statement: %v", err)
		os.Exit(1)
	}

	err = ctx.Commit()
	if err != nil {
		fmt.Printf("Error commiting change: %v", err)
		os.Exit(1)
	}

	newId, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("Error retrieving id: %v", err)
		os.Exit(1)
	}

	return int(newId)
}

func deleteProject(database *sql.DB, p project) {
	ctx, err := database.Begin()
	if err != nil {
		fmt.Printf("Error opening db connection: %v", err)
		os.Exit(1)
	}

	stmt, err := ctx.Prepare("DELETE FROM item WHERE pid = ?")
	if err != nil {
		fmt.Printf("Error preparing statement: %v", err)
		os.Exit(1)
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.id)
	if err != nil {
		fmt.Printf("Error executing statement: %v", err)
		os.Exit(1)
	}

	stmt, err = ctx.Prepare("DELETE FROM project WHERE id = ?")
	if err != nil {
		fmt.Printf("Error preparing statement: %v", err)
		os.Exit(1)
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.id)
	if err != nil {
		fmt.Printf("Error executing statement: %v", err)
		os.Exit(1)
	}

	err = ctx.Commit()
	if err != nil {
		fmt.Printf("Error commiting change: %v", err)
		os.Exit(1)
	}
}

func deleteItem(database *sql.DB, it item) {
	ctx, err := database.Begin()
	if err != nil {
		fmt.Printf("Error opening db connection: %v", err)
		os.Exit(1)
	}

	stmt, err := ctx.Prepare("DELETE FROM item WHERE id = ?")
	if err != nil {
		fmt.Printf("Error preparing statement: %v", err)
		os.Exit(1)
	}
	defer stmt.Close()

	_, err = stmt.Exec(it.id)
	if err != nil {
		fmt.Printf("Error executing statement: %v", err)
		os.Exit(1)
	}

	err = ctx.Commit()
	if err != nil {
		fmt.Printf("Error commiting change: %v", err)
		os.Exit(1)
	}
}

func updateProject(database *sql.DB, p project) {
	ctx, err := database.Begin()
	if err != nil {
		fmt.Printf("Error opening db connection: %v", err)
		os.Exit(1)
	}

	stmt, err := ctx.Prepare("UPDATE project SET name = ? WHERE id = ?")
	if err != nil {
		fmt.Printf("Error preparing statement: %v", err)
		os.Exit(1)
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.name, p.id)
	if err != nil {
		fmt.Printf("Error executing statement: %v", err)
		os.Exit(1)
	}

	err = ctx.Commit()
	if err != nil {
		fmt.Printf("Error commiting change: %v", err)
		os.Exit(1)
	}
}

func updateItem(database *sql.DB, it item, lid int, orderNr int) {
	ctx, err := database.Begin()
	if err != nil {
		fmt.Printf("Error opening db connection: %v", err)
		os.Exit(1)
	}

	stmt, err := ctx.Prepare("UPDATE item SET title = ?, desc = ?, lid = ?, orderNr = ? WHERE id = ?")
	if err != nil {
		fmt.Printf("Error preparing statement: %v", err)
		os.Exit(1)
	}
	defer stmt.Close()

	_, err = stmt.Exec(it.title, it.desc, lid, orderNr, it.id)
	if err != nil {
		fmt.Printf("Error executing statement: %v", err)
		os.Exit(1)
	}

	err = ctx.Commit()
	if err != nil {
		fmt.Printf("Error commiting change: %v", err)
		os.Exit(1)
	}
}
