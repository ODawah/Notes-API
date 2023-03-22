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
	if note.UserUuid == "" {
		return nil, errors.New("user uuid is required")
	}
	if validators.IsUUIDValid(note.UserUuid) == false {
		return nil, errors.New("invalid user uuid")
	}
	uuid := generators.UUIDGenerator()
	note.UUID = uuid
	database.DB.Create(&note)
	if err != nil {
		return nil, err
	}

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
	err = database.DB.Where("title LIKE ? AND user_uuid = ?", title, userUuid).First(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func DeleteNote(UUID, userUuid string) error {
	_, err := FindUserByUUID(userUuid)
	if err != nil {
		return err
	}
	valid := validators.IsUUIDValid(UUID)
	if !valid {
		return errors.New("invalid UUID")
	}
	rows := database.DB.Where("uuid = ? AND user_uuid = ?", UUID, userUuid).Delete(&schemas.Note{})
	if rows.RowsAffected == 0 {
		return errors.New("note not found")
	}

	return nil
}

func UpdateNote(note schemas.Note) error {
	_, err := FindUserByUUID(note.UserUuid)
	if err != nil {
		return err
	}
	valid := validators.IsUUIDValid(note.UUID)
	if !valid {
		return errors.New("invalid UUID")
	}
	note.Text = strings.TrimSpace(strings.ToLower(note.Text))
	note.Title = strings.TrimSpace(strings.ToLower(note.Title))
	rows := database.DB.Model(&schemas.Note{}).Where("uuid = ? AND user_uuid = ?", note.UUID, note.UserUuid).Updates(note)
	if rows.RowsAffected == 0 {
		return errors.New("note not found")
	}

	return nil

}

func FindNotes(uuid string) ([]schemas.Note, error) {
	_, err := FindUserByUUID(uuid)
	if err != nil {
		return nil, err
	}
	var res []schemas.Note
	database.DB.Where("user_uuid = ?", uuid).Find(&res)
	return res, nil
}
