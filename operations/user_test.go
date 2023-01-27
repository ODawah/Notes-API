package operations

import (
	"fmt"
	"testing"

	"github.com/Notes-App/database"
	"github.com/Notes-App/schemas"
	"github.com/Notes-App/validators"
)

func TestCreateUser(t *testing.T) {
	dbErr := database.Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer database.CleanUp()
	type test struct {
		name     string
		input    schemas.User
		expected *schemas.User
		err      error
	}

	tests := []test{
		{name: "normal User", input: schemas.User{Email: "test@gmail.com", Password: "ping pong"}, expected: &schemas.User{Email: "test@gmail.com", Password: "ping pong"}, err: nil},
		{name: "empty User", input: schemas.User{}, expected: nil, err: fmt.Errorf("invalid email address")},
		{name: "invalid email", input: schemas.User{Email: "asdafs@.com", Password: "ping pong"}, expected: nil, err: fmt.Errorf("invalid email address")},
		{name: "no password", input: schemas.User{Email: "test@test.com", Password: ""}, expected: nil, err: fmt.Errorf("empty password")},
	}

	for _, tc := range tests {
		t.Log(tc.name)
		got, err := CreateUser(tc.input)
		if fmt.Sprint(err) != fmt.Sprint(tc.err) {
			t.Fatalf("got: %s    expected:%s", err, tc.err)
		}
		if got == nil && tc.err == nil {
			t.Fatalf("got: %s    expected:%s", got, tc.expected)
		}
		if got != nil && tc.expected != nil {
			if !validators.IsUUIDValid(got.UUID) {
				t.Fatalf("got invalid uuid: %s", got.UUID)
			}
			if got.Email != tc.expected.Email {
				t.Fatalf("got: %s    expected:%s", got.Email, tc.expected.Email)
			}
			if got.Password == tc.expected.Password {
				t.Fatalf("got: %s    expected:%s", got.Password, tc.expected.Password)
			}
		}
	}
}
