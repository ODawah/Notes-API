package operations

import (
	"errors"
	"net/mail"

	"github.com/Notes-App/database"
	"github.com/Notes-App/encrypt"
	"github.com/Notes-App/generators"
	"github.com/Notes-App/schemas"
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
	_, err = database.DB.Exec("INSERT INTO users(uuid,email,password) VALUES (?,?,?)", uuid, user.Email, user.Password)
	user.UUID = uuid
	return &user, nil
}
