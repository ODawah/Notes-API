package generators

import (
	"testing"

	"github.com/Notes-App/models"
)

func TestUUIDGenerator(t *testing.T) {
	got := UUIDGenerator()
	res := models.IsUUIDValid(got)
	if !res {
		t.Fatal(res)
	}
}
