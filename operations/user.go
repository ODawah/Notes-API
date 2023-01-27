package operations

import (
	"errors"
	"net/mail"

	"github.com/Notes-App/database"
	"github.com/Notes-App/encrypt"
	"github.com/Notes-App/generators"
	"github.com/Notes-App/schemas"
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
	_, err = database.DB.Exec("INSERT INTO users(uuid,email,password) VALUES (?,?,?)", uuid, user.Email, user.Password)
	user.UUID = uuid
	return &user, nil
}

func FindUser(user schemas.User) (*schemas.User, error) {
	var dbUser schemas.User
	if user.Email == "" {
		return nil, errors.New("empty email")
	} else if user.Password == "" {
		return nil, errors.New("empty password")
	}
	err := database.DB.QueryRow("SELECT uuid,email,password FROM users WHERE email = ?", user.Email).Scan(&dbUser.UUID, &dbUser.Email, &dbUser.Password)
	if err != nil {
		return nil, errors.New("user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return nil, errors.New("wrong password")
	}
	return &dbUser, nil
}
