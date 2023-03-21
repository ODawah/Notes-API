package schemas

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	UUID     string `gorm:"primary_key" json:"uuid"`
	Title    string `json:"title"`
	Text     string `json:"text"`
	UserUuid string `json:"user_uuid"`
	user     User   `gorm:"foreignKey:UserUuid"`
}

type AllNotes struct {
	Notes []Note
}

type GetNote struct {
	Title string `json:"title"`
}
