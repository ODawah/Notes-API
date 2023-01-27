package database

import (
	"database/sql"
	"errors"
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
	statement, err := DB.Prepare("CREATE TABLE IF NOT EXISTS notes (uuid CHAR(36) NOT NULL check(LENGTH(uuid) = 36) PRIMARY KEY,title TEXT NOT NULL check(LENGTH(title) BETWEEN 1 AND 100),note_text TEXT NOT NULL check(LENGTH(note_text) BETWEEN 1 AND 500), user_uuid check(LENGTH(user_uuid) = 36),FOREIGN KEY (user_uuid) REFERENCES users(uuid))")
	if err != nil {
		return errors.New("failed to prepare statement")
	}
	statement.Exec()

	statement, err = DB.Prepare("CREATE TABLE IF NOT EXISTS users (uuid CHAR(36) NOT NULL check(LENGTH(uuid) = 36) PRIMARY KEY,email TEXT NOT NULL check(LENGTH(email) BETWEEN 1 AND 100),password TEXT NOT NULL check(LENGTH(password) BETWEEN 1 AND 1000))")
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
