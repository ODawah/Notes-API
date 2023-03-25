package models

import (
	"errors"
	"strings"

	"github.com/Notes-App/generators"
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	UUID     string `gorm:"primary_key" json:"uuid"`
	Title    string `gorm:"unique" json:"title"`
	Text     string `json:"text"`
	UserUuid string `json:"user_uuid"`
	user     User   `gorm:"foreignKey:UserUuid"`
}

type GetNote struct {
	Title string `json:"title"`
}

func CreateNote(note Note) (*Note, error) {
	err := IsValidText(note.Text)
	if err != nil {
		return nil, err
	}
	note.Text = strings.TrimSpace(strings.ToLower(note.Text))
	err = IsValidTitle(note.Title)
	if err != nil {
		return nil, err
	}
	note.Title = strings.TrimSpace(strings.ToLower(note.Title))
	if note.UserUuid == "" {
		return nil, errors.New("user uuid is required")
	}
	if IsUUIDValid(note.UserUuid) == false {
		return nil, errors.New("invalid user uuid")
	}
	uuid := generators.UUIDGenerator()
	note.UUID = uuid
	DB.Create(&note)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func FindNoteByTitle(title string, userUuid string) (*Note, error) {
	_, err := FindUserByUUID(userUuid)
	if err != nil {
		return nil, err
	}
	err = IsValidTitle(title)
	if err != nil {
		return nil, err
	}
	title = strings.TrimSpace(strings.ToLower(title))
	var note Note
	err = DB.Where("title LIKE ? AND user_uuid = ?", title, userUuid).First(&note).Error
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
	valid := IsUUIDValid(UUID)
	if !valid {
		return errors.New("invalid UUID")
	}
	rows := DB.Where("uuid = ? AND user_uuid = ?", UUID, userUuid).Delete(&Note{})
	if rows.RowsAffected == 0 {
		return errors.New("note not found")
	}

	return nil
}

func UpdateNote(note Note) error {
	_, err := FindUserByUUID(note.UserUuid)
	if err != nil {
		return err
	}
	valid := IsUUIDValid(note.UUID)
	if !valid {
		return errors.New("invalid UUID")
	}
	note.Text = strings.TrimSpace(strings.ToLower(note.Text))
	note.Title = strings.TrimSpace(strings.ToLower(note.Title))
	rows := DB.Model(&Note{}).Where("uuid = ? AND user_uuid = ?", note.UUID, note.UserUuid).Updates(note)
	if rows.RowsAffected == 0 {
		return errors.New("note not found")
	}

	return nil

}

func FindNotes(uuid string) ([]Note, error) {
	_, err := FindUserByUUID(uuid)
	if err != nil {
		return nil, err
	}
	var res []Note
	DB.Where("user_uuid = ?", uuid).Find(&res)
	return res, nil
}
