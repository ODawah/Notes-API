package models

import (
	"errors"
	"net/mail"

	"github.com/Notes-App/encrypt"
	"github.com/Notes-App/generators"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID     string `gorm:"primary_key" json:"uuid"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
}

func CreateUser(user User) (*User, error) {
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return nil, errors.New("invalid email address")
	}
	uuid := generators.UUIDGenerator()
	if user.Password == "" {
		return nil, errors.New("empty password")
	}
	hashed, err := encrypt.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashed
	user.UUID = uuid
	DB.Create(&user)
	return &user, nil
}

func FindUser(user User) (*User, error) {
	var dbUser *User
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return nil, errors.New("invalid email address")
	}
	if user.Password == "" {
		return nil, errors.New("empty password")
	}
	err = DB.Where("email = ?", user.Email).First(&dbUser).Error
	if err != nil {
		return nil, errors.New("user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return nil, errors.New("wrong password")
	}
	return dbUser, nil
}

func FindUserByUUID(uuid string) (*User, error) {
	var dbUser User
	if !IsUUIDValid(uuid) {
		return nil, errors.New("invalid uuid")
	}
	err := DB.Where("uuid = ?", uuid).First(&dbUser).Error
	if err != nil || &dbUser == nil {
		return nil, errors.New("user not found")
	}
	return &dbUser, nil
}
