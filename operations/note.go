package operations

import (
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
