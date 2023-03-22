package operations

import (
	"errors"
	"net/mail"

	"github.com/Notes-App/database"
	"github.com/Notes-App/encrypt"
	"github.com/Notes-App/generators"
	"github.com/Notes-App/schemas"
	"github.com/Notes-App/validators"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user schemas.User) (*schemas.User, error) {
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
	database.DB.Create(&user)
	return &user, nil
}

func FindUser(user schemas.User) (*schemas.User, error) {
	var dbUser *schemas.User
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return nil, errors.New("invalid email address")
	}
	if user.Password == "" {
		return nil, errors.New("empty password")
	}
	err = database.DB.Where("email = ?", user.Email).First(&dbUser).Error
	if err != nil {
		return nil, errors.New("user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return nil, errors.New("wrong password")
	}
	return dbUser, nil
}

func FindUserByUUID(uuid string) (*schemas.User, error) {
	var dbUser schemas.User
	if !validators.IsUUIDValid(uuid) {
		return nil, errors.New("invalid uuid")
	}
	err := database.DB.Where("uuid = ?", uuid).First(&dbUser).Error
	if err != nil || &dbUser == nil {
		return nil, errors.New("user not found")
	}
	return &dbUser, nil
}
