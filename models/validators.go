package models

import (
	"errors"

	"github.com/google/uuid"
)

func IsValidText(text string) error {
	if text == "" {
		return errors.New("empty note text")
	}
	if len(text) > 500 {
		return errors.New("long note text")
	}
	return nil
}

func IsValidTitle(title string) error {
	if title == "" {
		return errors.New("empty note title")
	}
	if len(title) > 100 {
		return errors.New("long note title")
	}
	return nil
}

func IsUUIDValid(UUID string) bool {
	_, err := uuid.Parse(UUID)
	if err != nil {
		return false
	}
	return true
}
