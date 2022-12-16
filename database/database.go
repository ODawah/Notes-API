package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Connect() error {
	var err error
	DB, err = sql.Open("sqlite3", "notes.db")
	if err != nil {
		return err
	}
	statement, err := DB.Prepare("CREATE TABLE IF NOT EXISTS notes (uuid CHAR(36) NOT NULL check(LENGTH(uuid) = 36) PRIMARY KEY,title TEXT NOT NULL check(LENGTH(title) BETWEEN 1 AND 100),note_text TEXT NOT NULL check(LENGTH(title) BETWEEN 1 AND 500))")
	if err != nil {
		return err
	}
	statement.Exec()

	return nil
}

func CleanUp() error {
	err := os.Remove("notes.db")
	if err != nil {
		return err
	}
	return nil
}
