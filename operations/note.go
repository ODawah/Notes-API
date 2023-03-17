package operations

import (
	"errors"
	"fmt"
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
	if note.UserUuid == "" {
		return nil, errors.New("user uuid is required")
	}
	if validators.IsUUIDValid(note.UserUuid) == false {
		return nil, errors.New("invalid user uuid")
	}
	uuid := generators.UUIDGenerator()
	_, err = database.DB.Exec("INSERT INTO notes(uuid,title, note_text, user_uuid) VALUES (?,?,?,?)", uuid, note.Title, note.Text, note.UserUuid)
	if err != nil {
		return nil, err
	}
	note.UUID = uuid

	return &note, nil
}

func FindNoteByTitle(title string, userUuid string) (*schemas.Note, error) {
	_, err := FindUserByUUID(userUuid)
	if err != nil {
		return nil, err
	}
	err = validators.IsValidTitle(title)
	if err != nil {
		return nil, err
	}
	title = strings.TrimSpace(strings.ToLower(title))
	var note schemas.Note
	err = database.DB.QueryRow("SELECT * FROM notes WHERE user_uuid = ? AND (title LIKE ?)", userUuid, title).Scan(&note.UUID, &note.Title, &note.Text, &note.UserUuid)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func DeleteNote(UUID, userUuid string) (int, error) {
	_, err := FindUserByUUID(userUuid)
	if err != nil {
		return 0, err
	}
	valid := validators.IsUUIDValid(UUID)
	if !valid {
		return 0, errors.New("invalid UUID")
	}
	res, err := database.DB.Exec("DELETE FROM notes WHERE uuid = ? AND user_uuid = ?", UUID, userUuid)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return 0, errors.New("note not found")
	}
	return int(n), nil
}

func UpdateNote(UUID string, note schemas.Note) (int, error) {
	_, err := FindUserByUUID(note.UserUuid)
	if err != nil {
		return 0, err
	}
	var query string
	valid := validators.IsUUIDValid(UUID)
	if !valid {
		return 0, errors.New("invalid UUID")
	}
	note.Text = strings.TrimSpace(strings.ToLower(note.Text))
	note.Title = strings.TrimSpace(strings.ToLower(note.Title))
	if note.Title == "" && note.Text != "" {
		query = fmt.Sprintf("UPDATE notes set note_text = \"%s\" WHERE uuid = \"%s\" AND user_uuid = \"%s\" ", note.Text, UUID, note.UserUuid)
	} else if note.Text == "" && note.Title != "" {
		query = fmt.Sprintf("UPDATE notes set title = \"%s\" WHERE uuid = \"%s\" AND user_uuid = \"%s\"", note.Title, UUID, note.UserUuid)
	} else if note.Text != "" && note.Title != "" {
		query = fmt.Sprintf("UPDATE notes set title = \"%s\", note_text = \"%s\" WHERE uuid = \"%s\" AND user_uuid = \"%s\"", note.Title, note.Text, UUID, note.UserUuid)
	}
	if query == "" {
		return 0, errors.New("empty updates")
	}
	res, err := database.DB.Exec(query)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return 0, errors.New("note not found")
	}
	return int(n), nil

}

func FindNotes(uuid string) (*schemas.AllNotes, error) {
	_, err := FindUserByUUID(uuid)
	if err != nil {
		return nil, err
	}
	rows, err := database.DB.Query("SELECT * FROM notes WHERE user_uuid = ?", uuid)
	if err != nil {
		return nil, err
	}
	var res schemas.AllNotes
	for rows.Next() {
		var note schemas.Note
		err = rows.Scan(&note.UUID, &note.Title, &note.Text, &note.UserUuid)
		if err != nil {
			return nil, err
		}
		res.Notes = append(res.Notes, note)
	}
	return &res, nil
}
