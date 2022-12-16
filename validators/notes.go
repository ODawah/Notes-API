package validators

import (
	"errors"
)

func IsValidText(text string) error {
	if text == "" {
		return errors.New("empty note text")
	}
	if len(text) > 500 {
		return errors.New("large note text")
	}
	return nil
}

func IsValidTitle(title string) error {
	if title == "" {
		return errors.New("empty note title")
	}
	if len(title) > 100 {
		return errors.New("large note title")
	}
	return nil
}
