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
	statement, err := DB.Prepare("CREATE TABLE IF NOT EXISTS USER (uuid CHAR(36) NOT NULL check(LENGTH(uuid) = 36) PRIMARY KEY, name CHAR(60) NOT NULL check(LENGTH(name) BETWEEN 3 AND 60), email TEXT UNIQUE NOT NULL ,password TEXT NOT NULL)")
	if err != nil {
		return err
	}
	statement.Exec()

	statement, err = DB.Prepare("CREATE TABLE IF NOT EXISTS quotes (uuid CHAR(36) NOT NULL check(LENGTH(uuid) = 36) PRIMARY KEY,title TEXT NOT NULL check(LENGTH(title) BETWEEN 1 AND 100),note TEXT NOT NULL check(LENGTH(title) BETWEEN 1 AND 500),user_uuid NOT NULL check(LENGTH(user_uuid) = 36),FOREIGN KEY (user_uuid) REFERENCES USER(uuid))")
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
