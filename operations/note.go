package operations

import (
	"errors"
	"strings"

	"github.com/Notes-App/database"
	"github.com/Notes-App/generators"
	"github.com/Notes-App/schemas"
	"github.com/Notes-App/validators"
)

func CreateNote(note schemas.Note) (*schemas.Note, error) {
	err := validators.IsValidText(note.Text)
	if err != nil {
		return nil, err
	}
	note.Text = strings.TrimSpace(strings.ToLower(note.Text))
	err = validators.IsValidTitle(note.Title)
	if err != nil {
		return nil, err
	}
	note.Title = strings.TrimSpace(strings.ToLower(note.Title))
	uuid := generators.UUIDGenerator()
	_, err = database.DB.Exec("INSERT INTO notes(uuid,title, note_text) VALUES (?,?,?)", uuid, note.Title, note.Text)
	if err != nil {
		return nil, err
	}
	note.UUID = uuid

	return &note, nil
}

func FindNoteByTitle(title string) (*schemas.Note, error) {
	err := validators.IsValidTitle(title)
	if err != nil {
		return nil, err
	}
	title = strings.TrimSpace(strings.ToLower(title))
	var note schemas.Note
	err = database.DB.QueryRow("SELECT uuid,title,note_text FROM notes WHERE title LIKE ?", title).Scan(&note.UUID, &note.Title, &note.Text)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func DeleteNote(UUID string) (int, error) {
	valid := validators.IsUUIDValid(UUID)
	if !valid {
		return 0, errors.New("invalid UUID")
	}
	res, err := database.DB.Exec("DELETE FROM notes WHERE uuid = ?", UUID)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return 0, errors.New("note not found")
	}
	return int(n), nil
}