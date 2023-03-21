package schemas

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID     string `gorm:"primary_key" json:"uuid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
