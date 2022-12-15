package database

import (
	"os"
	"testing"
)

func TestCleanUp(t *testing.T) {
	err := Connect()
	if err != nil {
		t.Fatal(err)
	}

	err = CleanUp()
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat("notes.db")
	if err == nil {
		t.Log(err)
	}
}
