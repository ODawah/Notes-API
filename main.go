package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, _ := sql.Open("sqlite3", "notes.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS authors (uuid CHAR(36) NOT NULL check(LENGTH(uuid) = 36) UNIQUE,id INTEGER PRIMARY KEY AUTOINCREMENT ,name CHAR(60) NOT NULL check(LENGTH(name) BETWEEN 3 AND 60))")
	statement.Exec()
	fmt.Println(database)
	database.Close()
	fmt.Println(database)
}
