package generators

import (
	"testing"

	"github.com/Notes-App/validators"
)

func TestUUIDGenerator(t *testing.T) {
	got := UUIDGenerator()
	res := validators.IsUUIDValid(got)
	if !res {
		t.Fatal(res)
	}
}
