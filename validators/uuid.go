package validators

import (
	"github.com/google/uuid"
)

func IsUUIDValid(UUID string) bool {
	_, err := uuid.Parse(UUID)
	if err != nil {
		return false
	}
	return true
}
