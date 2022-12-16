package generators

import (
	"github.com/google/uuid"
)

func UUIDGenerator() string {
	UUID := uuid.New().String()

	return UUID
}
